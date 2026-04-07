package ir

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"yz/internal/ast"
	"yz/internal/sema"
	"yz/internal/token"
)

// ---------------------------------------------------------------------------
// Public entry point
// ---------------------------------------------------------------------------

// Lower converts an analyzed source file into an IR File.
// pkgName is the Go package name to emit (usually "main" for the entry file).
func Lower(sf *ast.SourceFile, a *sema.Analyzer, pkgName string) *File {
	l := &lowerer{analyzer: a, thunkVars: make(map[string]bool)}
	return l.lowerFile(sf, pkgName)
}

// ---------------------------------------------------------------------------
// Lowerer state
// ---------------------------------------------------------------------------

type lowerer struct {
	analyzer *sema.Analyzer
	irFile   *File // current output file — used to accumulate extra imports

	// When lowering a method body, these describe the receiver.
	recvName   string
	recvFields map[string]bool // fields accessible via receiver

	// thunkVars tracks local variables that hold *Thunk[T] values (declared
	// via a: bocCall(...)). When referenced as plain values, these are
	// auto-forced to make the Yz "a is the value" semantics transparent.
	thunkVars map[string]bool
}

// ---------------------------------------------------------------------------
// File
// ---------------------------------------------------------------------------

func (l *lowerer) lowerFile(sf *ast.SourceFile, pkgName string) *File {
	f := &File{PkgName: pkgName}
	l.irFile = f
	for _, node := range sf.Stmts {
		if d := l.lowerTopLevel(node); d != nil {
			f.Decls = append(f.Decls, d)
		}
	}
	return f
}

// addImport adds importPath to the ir.File's import list, deduplicating.
func (l *lowerer) addImport(importPath string) {
	if l.irFile == nil {
		return
	}
	for _, imp := range l.irFile.Imports {
		if imp == importPath {
			return
		}
	}
	l.irFile.Imports = append(l.irFile.Imports, importPath)
}

func (l *lowerer) lowerTopLevel(node ast.Node) Decl {
	switch n := node.(type) {
	case *ast.ShortDecl:
		return l.lowerTopShortDecl(n)
	case *ast.BocWithSig:
		// Uppercase name + no body → type-only declaration: emit struct without constructor.
		if n.Body == nil && isUppercase(n.Name.Name) {
			return l.lowerTypeOnlyDecl(n)
		}
		return l.lowerBocWithSig(n, "", nil)
	case *ast.Assignment:
		return l.lowerTopAssignment(n)
	default:
		return nil
	}
}

// lowerTopAssignment handles the declare-then-assign pattern at top level:
//
//	greet #(name String)       ← declares greet as a BocType (no body)
//	greet = { name String; … } ← assigns the body, emitted as a FuncDecl
//
// Only Assignment nodes whose target is a known BocType symbol and whose
// value is a BocLiteral are lowered; all others are silently ignored.
func (l *lowerer) lowerTopAssignment(asgn *ast.Assignment) Decl {
	if asgn.Target == nil || len(asgn.Values) != 1 {
		return nil
	}
	id, ok := asgn.Target.(*ast.Ident)
	if !ok {
		return nil
	}
	bocLit, ok := asgn.Values[0].(*ast.BocLiteral)
	if !ok {
		return nil
	}
	sym := l.analyzer.LookupInFile(id.Name)
	if sym == nil {
		return nil
	}
	bt, ok := sym.Type.(*sema.BocType)
	if !ok {
		return nil
	}

	// Use AST sig for return type to preserve generic TypeArgs; also collect type params.
	var typeParams []string
	resultType := "std.Unit"
	if bws, ok := sym.Node.(*ast.BocWithSig); ok {
		resultType = l.getResultTypeFromSig(bws.Sig, bt, bws.BodyOnly)
		typeParams = collectSigTypeParams(bws.Sig)
	} else if len(bt.Returns) > 0 {
		resultType = l.goType(bt.Returns[0])
	}
	thunkResult := "*std.Thunk[" + resultType + "]"

	// Collect params from the body's leading TypedDecls (same as lowerMethod).
	var params []*ParamSpec
	for _, elem := range bocLit.Elements {
		if td, ok := elem.(*ast.TypedDecl); ok && td.Value == nil {
			params = append(params, &ParamSpec{
				Name: td.Name.Name,
				Type: l.goTypeFromTypeExpr(td.Type),
			})
		}
	}

	prev := l.setReceiver("", nil)
	bocBodyStmts := l.lowerBocBody(bocLit, resultType)
	l.restoreReceiver(prev)

	var funcBody []Stmt
	if len(bocBodyStmts) == 1 {
		if es, ok := bocBodyStmts[0].(*ExprStmt); ok {
			funcBody = []Stmt{&ReturnStmt{Value: es.Expr}}
		}
	}
	if funcBody == nil {
		funcBody = bocBodyStmts
	}

	return &FuncDecl{
		Name:       id.Name,
		TypeParams: typeParams,
		Params:     params,
		Results:    []string{thunkResult},
		Body:       funcBody,
	}
}

// lowerTypeOnlyDecl lowers `Name #(params)` (no body) into a Go struct
// declaration without a constructor. The struct type exists so it can be used
// as a parameter/variable type; instances cannot be created until a body is
// attached (structural typing is enforced by sema).
func (l *lowerer) lowerTypeOnlyDecl(bws *ast.BocWithSig) Decl {
	semType := l.analyzer.ExprType(bws)
	st, ok := semType.(*sema.StructType)
	if !ok {
		return nil
	}
	if st.IsInterface {
		// All fields are BocTypes → emit a Go interface.
		id := &InterfaceDecl{Name: bws.Name.Name}
		for _, f := range st.Fields {
			bt, _ := f.Type.(*sema.BocType)
			resultType := "std.Unit"
			if bt != nil && len(bt.Returns) == 1 {
				resultType = l.goType(bt.Returns[0])
			}
			var params []*ParamSpec
			if bt != nil {
				for _, p := range bt.Params {
					if !p.IsReturn && p.Label != "" {
						params = append(params, &ParamSpec{
							Name: p.Label,
							Type: l.goType(p.Type),
						})
					}
				}
			}
			id.Methods = append(id.Methods, &InterfaceMethod{
				Name:       f.Name,
				Params:     params,
				ResultType: resultType,
			})
		}
		return id
	}
	// Mixed or data-only: separate data fields from BocType fields.
	var dataFields, bocFields []sema.StructField
	for _, f := range st.Fields {
		if _, isBoc := f.Type.(*sema.BocType); isBoc {
			bocFields = append(bocFields, f)
		} else {
			dataFields = append(dataFields, f)
		}
	}

	// Data fields only (no BocType params) → plain struct, no constructor.
	if len(bocFields) == 0 {
		sd := &StructDecl{Name: bws.Name.Name, NoConstructor: true}
		for _, f := range dataFields {
			sd.Fields = append(sd.Fields, &FieldSpec{Name: f.Name, Type: l.goType(f.Type)})
		}
		return sd
	}

	// Mixed: data fields + BocType fields → struct with constructor + method wrappers.
	sd := &StructDecl{Name: bws.Name.Name}
	for _, f := range dataFields {
		sd.Fields = append(sd.Fields, &FieldSpec{Name: f.Name, Type: l.goType(f.Type)})
	}
	for _, f := range bocFields {
		bt := f.Type.(*sema.BocType)
		sd.Fields = append(sd.Fields, &FieldSpec{Name: f.Name, Type: l.bocFuncType(bt)})
		sd.Methods = append(sd.Methods, l.bocFieldMethod(bws.Name.Name, f.Name, bt))
	}
	return sd
}

// bocFuncType returns the Go function type string for a BocType field.
// e.g. BocType{Params:[], Returns:[Unit]} → "func() *std.Thunk[std.Unit]"
func (l *lowerer) bocFuncType(bt *sema.BocType) string {
	var paramTypes []string
	for _, p := range bt.Params {
		if !p.IsReturn {
			paramTypes = append(paramTypes, l.goType(p.Type))
		}
	}
	resultType := "std.Unit"
	if len(bt.Returns) > 0 {
		resultType = l.goType(bt.Returns[0])
	}
	return fmt.Sprintf("func(%s) *std.Thunk[%s]", strings.Join(paramTypes, ", "), resultType)
}

// bocFieldMethod generates a method wrapper that delegates to a function field.
// e.g. field "greet" → func (self *Name) Greet() *std.Thunk[std.Unit] { return self.greet() }
func (l *lowerer) bocFieldMethod(typeName, fieldName string, bt *sema.BocType) *MethodDecl {
	var params []*ParamSpec
	var args []Expr
	for _, p := range bt.Params {
		if !p.IsReturn && p.Label != "" {
			params = append(params, &ParamSpec{Name: p.Label, Type: l.goType(p.Type)})
			args = append(args, &Ident{Name: p.Label})
		}
	}
	resultType := "std.Unit"
	if len(bt.Returns) > 0 {
		resultType = l.goType(bt.Returns[0])
	}
	return &MethodDecl{
		RecvType: "*" + typeName,
		RecvName: "self",
		Name:     capitalize(fieldName),
		Params:   params,
		Results:  []string{"*std.Thunk[" + resultType + "]"},
		Body: []Stmt{
			&ReturnStmt{Value: &FuncCall{
				Func: &FieldAccess{Object: &Ident{Name: "self"}, Field: fieldName},
				Args: args,
			}},
		},
	}
}

// fillDefaults fills in lowered default expressions for omitted optional params
// in a BocWithSig call. sigParams is the raw AST param list from the signature.
// loweredArgs is the already-lowered positional args from the call site.
// Any trailing params that have a Default expression are filled in automatically.
func (l *lowerer) fillDefaults(loweredArgs []Expr, sigParams []*ast.BocParam) []Expr {
	// Collect named input params (same filter as lowerBocWithSig).
	var inputs []*ast.BocParam
	for _, p := range sigParams {
		if p.Variant == nil && p.Label != "" {
			inputs = append(inputs, p)
		}
	}
	if len(loweredArgs) >= len(inputs) {
		return loweredArgs // all args already provided
	}
	result := make([]Expr, len(inputs))
	copy(result, loweredArgs)
	for i := len(loweredArgs); i < len(inputs); i++ {
		if inputs[i].Default != nil {
			result[i] = l.lowerExpr(inputs[i].Default)
		}
	}
	return result
}

// ---------------------------------------------------------------------------
// Top-level ShortDecl dispatch
// ---------------------------------------------------------------------------

func (l *lowerer) lowerTopShortDecl(d *ast.ShortDecl) Decl {
	if len(d.Names) != 1 || len(d.Values) != 1 {
		return nil // multi-assign at top level — not yet supported
	}
	name := d.Names[0]
	bocLit, isBoc := d.Values[0].(*ast.BocLiteral)
	if !isBoc {
		// A simple `x: expr` at top level — treat as a singleton with one field.
		return l.lowerSimpleTopDecl(name.Name, d.Values[0])
	}

	// name #(params) { body } is BocWithSig; plain name: { body } lands here.
	if isUppercase(name.Name) {
		return l.lowerStructBoc(name.Name, bocLit)
	}
	if name.Name == "main" {
		return l.lowerMainBoc(bocLit)
	}
	return l.lowerSingletonBoc(name.Name, bocLit)
}

// lowerSimpleTopDecl wraps a scalar top-level assignment in a singleton with
// a single field (e.g. `x: 42` → var x = &_xBoc{x: NewInt(42)}).
func (l *lowerer) lowerSimpleTopDecl(name string, val ast.Expr) Decl {
	expr := l.lowerExpr(val)
	typ := l.goType(l.analyzer.ExprType(val))
	return &SingletonDecl{
		TypeName: "_" + name + "Boc",
		VarName:  name,
		Fields:   []*FieldSpec{{Name: name, Type: typ, Init: expr}},
	}
}

// ---------------------------------------------------------------------------
// Singleton boc
// ---------------------------------------------------------------------------

func (l *lowerer) lowerSingletonBoc(name string, b *ast.BocLiteral) *SingletonDecl {
	typeName := "_" + name + "Boc"
	sd := &SingletonDecl{TypeName: typeName, VarName: capitalize(name)}

	// Collect field names so methods can detect field access.
	fieldNames := l.collectFieldNames(b)

	for _, elem := range b.Elements {
		switch e := elem.(type) {
		case *ast.ShortDecl:
			if len(e.Names) == 1 && len(e.Values) == 1 {
				inner, isInnerBoc := e.Values[0].(*ast.BocLiteral)
				if isInnerBoc && !isUppercase(e.Names[0].Name) {
					// Inner boc → method. The BocType is on the ShortDecl, not the BocLiteral.
					bocSemType := l.analyzer.ExprType(e)
					m := l.lowerMethod(e.Names[0].Name, "*"+typeName, inner, fieldNames, bocSemType)
					sd.Methods = append(sd.Methods, m)
					continue
				}
			}
			// Otherwise it's a field.
			for i, n := range e.Names {
				var initExpr Expr
				if i < len(e.Values) {
					initExpr = l.lowerExpr(e.Values[i])
				}
				typ := l.goType(l.analyzer.ExprType(l.valueAt(e.Values, i)))
				sd.Fields = append(sd.Fields, &FieldSpec{Name: n.Name, Type: typ, Init: initExpr})
			}

		case *ast.TypedDecl:
			// Typed field without value (parameter-style).
			typ := l.goTypeFromTypeExpr(e.Type)
			var initExpr Expr
			if e.Value != nil {
				initExpr = l.lowerExpr(e.Value)
			}
			sd.Fields = append(sd.Fields, &FieldSpec{Name: e.Name.Name, Type: typ, Init: initExpr})

		case *ast.BocWithSig:
			if e.Body != nil {
				m := l.lowerBocWithSigAsMethod(e, "*"+typeName, fieldNames)
				sd.Methods = append(sd.Methods, m)
			}
		}
	}
	return sd
}

// collectFieldNames returns the set of field names (non-boc ShortDecls and TypedDecls
// without values) in the boc literal, for use in method body lowering.
// MixStmt fields are also included so methods can reference them as self.field.
func (l *lowerer) collectFieldNames(b *ast.BocLiteral) map[string]bool {
	fields := map[string]bool{}
	for _, elem := range b.Elements {
		switch e := elem.(type) {
		case *ast.ShortDecl:
			if len(e.Names) == 1 && len(e.Values) == 1 {
				if _, isBoc := e.Values[0].(*ast.BocLiteral); isBoc {
					continue // method, not field
				}
			}
			for _, n := range e.Names {
				fields[n.Name] = true
			}
		case *ast.TypedDecl:
			if e.Value == nil {
				fields[e.Name.Name] = true
			}
		case *ast.MixStmt:
			sym := l.analyzer.LookupInFile(e.Name.Name)
			if sym == nil {
				break
			}
			if mixedSt, ok := sym.Type.(*sema.StructType); ok {
				for _, f := range mixedSt.Fields {
					if _, isBoc := f.Type.(*sema.BocType); isBoc {
						continue // method, not a data field
					}
					fields[f.Name] = true
				}
			}
		}
	}
	return fields
}

// ---------------------------------------------------------------------------
// Method lowering
// ---------------------------------------------------------------------------

func (l *lowerer) lowerMethod(name, recvType string, b *ast.BocLiteral, parentFields map[string]bool, semType sema.Type) *MethodDecl {
	// Collect params (TypedDecl with no value inside the body).
	var params []*ParamSpec
	for _, elem := range b.Elements {
		if td, ok := elem.(*ast.TypedDecl); ok && td.Value == nil {
			params = append(params, &ParamSpec{
				Name: td.Name.Name,
				Type: l.goTypeFromTypeExpr(td.Type),
			})
		}
	}

	// Infer return type from the provided sema type (set on the ShortDecl node).
	resultType := "std.Unit"
	if bt, ok := semType.(*sema.BocType); ok && len(bt.Returns) > 0 {
		resultType = l.goType(bt.Returns[0])
	}
	thunkResult := "*std.Thunk[" + resultType + "]"

	// Lower method body with receiver context.
	prev := l.setReceiver("self", parentFields)
	body := l.lowerBocBody(b, resultType)
	l.restoreReceiver(prev)

	return &MethodDecl{
		RecvType: recvType,
		RecvName: "self",
		Name:     capitalize(name),
		Params:   params,
		Results:  []string{thunkResult},
		Body:     body,
	}
}

// lowerBocBody lowers the contents of a boc into a single ThunkExpr statement.
// All boc method bodies become `return std.Go(func() ResultType { ... })`.
func (l *lowerer) lowerBocBody(b *ast.BocLiteral, resultType string) []Stmt {
	var inner []Stmt
	elems := b.Elements
	for i, elem := range elems {
		isLast := i == len(elems)-1
		switch e := elem.(type) {
		case *ast.TypedDecl:
			if e.Value == nil {
				continue // param — already collected
			}
			inner = append(inner, &DeclStmt{
				Name: e.Name.Name,
				Init: l.lowerExpr(e.Value),
			})
		case *ast.ShortDecl:
			if isLast && len(e.Names) == 1 && len(e.Values) == 1 {
				// Last short decl in body — could be an expression-result.
				inner = append(inner, l.lowerBodyShortDecl(e, true, resultType))
			} else {
				inner = append(inner, l.lowerBodyShortDecl(e, false, resultType))
			}
		case *ast.Assignment:
			inner = append(inner, l.lowerAssignment(e))
		case *ast.ReturnStmt:
			var val Expr
			if e.Value != nil {
				val = l.lowerExpr(e.Value)
			}
			inner = append(inner, &ReturnStmt{Value: val})
		case ast.Expr:
			if fs, ok := l.tryLowerWhile(e); ok {
				inner = append(inner, fs)
				if isLast {
					inner = append(inner, &ReturnStmt{Value: &UnitLit{}})
				}
			} else if is, ok := l.tryLowerConditional(e); ok {
				inner = append(inner, is)
				if isLast {
					inner = append(inner, &ReturnStmt{Value: &UnitLit{}})
				}
			} else if !isLast {
				// Match in non-last position: lower as statement (no return value needed).
				if ms, ok := l.tryLowerMatch(e); ok {
					inner = append(inner, ms)
				} else {
					inner = append(inner, &ExprStmt{Expr: l.lowerExpr(e)})
				}
			} else {
				expr := l.lowerExpr(e)
				if isLast {
					// If the last expression is a boc method call it returns *Thunk[T],
					// but the closure expects T — force the thunk.
					if l.isBocMethodCall(e) {
						expr = &ForceExpr{Thunk: expr}
					}
					inner = append(inner, &ReturnStmt{Value: expr})
				} else {
					inner = append(inner, &ExprStmt{Expr: expr})
				}
			}
		default:
			// Other statements — skip for now.
		}
	}

	// Ensure there's always a return.
	if len(inner) == 0 || !isReturnStmt(inner[len(inner)-1]) {
		inner = append(inner, &ReturnStmt{Value: &UnitLit{}})
	}

	thunk := &ThunkExpr{
		ResultType: resultType,
		Body:       inner,
		Spawn:      true,
	}
	return []Stmt{&ExprStmt{Expr: thunk}}
}

func isReturnStmt(s Stmt) bool {
	_, ok := s.(*ReturnStmt)
	return ok
}

func (l *lowerer) lowerBodyShortDecl(d *ast.ShortDecl, isLast bool, resultType string) Stmt {
	if len(d.Names) == 1 && len(d.Values) == 1 {
		name := d.Names[0]
		val := d.Values[0]
		expr := l.lowerExpr(val)
		typ := l.goTypeForVar(l.analyzer.ExprType(val))
		if isLast {
			// Last expr in body: declare AND return if it's an expression.
			// But if it's a boc, make it a method.
			if _, isBoc := val.(*ast.BocLiteral); !isBoc {
				return &DeclStmt{Name: name.Name, Type: typ, Init: expr}
			}
		}
		return &DeclStmt{Name: name.Name, Type: typ, Init: expr}
	}
	// Multi-name: just declare each.
	// (simplified — multi-assign in method body not fully supported yet)
	if len(d.Names) > 0 && len(d.Values) > 0 {
		return &DeclStmt{
			Name: d.Names[0].Name,
			Init: l.lowerExpr(d.Values[0]),
		}
	}
	return &ExprStmt{Expr: &UnitLit{}}
}

func (l *lowerer) lowerAssignment(asgn *ast.Assignment) Stmt {
	var target Expr
	if asgn.Target != nil {
		target = l.lowerExpr(asgn.Target)
	} else if len(asgn.Names) > 0 {
		target = l.lowerName(asgn.Names[0].Name)
	}
	var val Expr
	if len(asgn.Values) > 0 {
		val = l.lowerExpr(asgn.Values[0])
	}
	return &AssignStmt{Target: target, Value: val}
}

// ---------------------------------------------------------------------------
// Struct type boc (uppercase)
// ---------------------------------------------------------------------------

func (l *lowerer) lowerStructBoc(name string, b *ast.BocLiteral) *StructDecl {
	// Detect variant (sum) type: body consists entirely of VariantDef elements.
	if l.isVariantBoc(b) {
		return l.lowerVariantBoc(name, b)
	}

	sd := &StructDecl{Name: name}
	fieldNames := l.collectFieldNames(b)
	for _, elem := range b.Elements {
		switch e := elem.(type) {
		case *ast.TypedDecl:
			typ := l.goTypeFromTypeExpr(e.Type)
			var init Expr
			if e.Value != nil {
				init = l.lowerExpr(e.Value)
			}
			sd.Fields = append(sd.Fields, &FieldSpec{Name: e.Name.Name, Type: typ, Init: init})
		case *ast.ShortDecl:
			if len(e.Names) == 1 && len(e.Values) == 1 {
				inner, isInnerBoc := e.Values[0].(*ast.BocLiteral)
				if isInnerBoc && !isUppercase(e.Names[0].Name) {
					bocSemType := l.analyzer.ExprType(e)
					m := l.lowerMethod(e.Names[0].Name, "*"+name, inner, fieldNames, bocSemType)
					sd.Methods = append(sd.Methods, m)
					continue
				}
			}
			for i, n := range e.Names {
				var initExpr Expr
				if i < len(e.Values) {
					initExpr = l.lowerExpr(e.Values[i])
				}
				typ := l.goType(l.analyzer.ExprType(l.valueAt(e.Values, i)))
				sd.Fields = append(sd.Fields, &FieldSpec{Name: n.Name, Type: typ, Init: initExpr})
			}
		case *ast.MixStmt:
			sym := l.analyzer.LookupInFile(e.Name.Name)
			if sym == nil {
				break
			}
			mixedSt, ok := sym.Type.(*sema.StructType)
			if !ok {
				break
			}
			var subFields []*FieldSpec
			for _, f := range mixedSt.Fields {
				if _, isBoc := f.Type.(*sema.BocType); isBoc {
					continue // method, not a data field
				}
				subFields = append(subFields, &FieldSpec{Name: f.Name, Type: l.goType(f.Type)})
			}
			sd.Fields = append(sd.Fields, &FieldSpec{
				Name:           e.Name.Name,
				Embedded:       true,
				EmbeddedFields: subFields,
			})

		case *ast.BocWithSig:
			if e.Body != nil {
				m := l.lowerBocWithSigAsMethod(e, "*"+name, fieldNames)
				sd.Methods = append(sd.Methods, m)
			}
		}
	}
	return sd
}

// ---------------------------------------------------------------------------
// Variant (sum) type boc
// ---------------------------------------------------------------------------

// isVariantBoc returns true when all non-empty elements in a boc body are
// VariantDefs or generic type param idents — indicating a sum type declaration.
func (l *lowerer) isVariantBoc(b *ast.BocLiteral) bool {
	hasVariant := false
	for _, elem := range b.Elements {
		if _, ok := elem.(*ast.VariantDef); ok {
			hasVariant = true
		} else if id, ok := elem.(*ast.Ident); ok && id.TokType == token.GENERIC_IDENT {
			// Generic type parameter (e.g., V in Option: { V; Some(value V); None() })
			continue
		} else {
			return false // mixed with non-variant content
		}
	}
	return hasVariant
}

// lowerVariantBoc lowers `Pet: { Cat(name String), Dog(name String) }` into a
// StructDecl with IsVariant=true, a merged flat field list, and one
// IRVariantCase per constructor.
func (l *lowerer) lowerVariantBoc(name string, b *ast.BocLiteral) *StructDecl {
	sd := &StructDecl{Name: name, IsVariant: true}
	fieldSet := map[string]bool{}

	for _, elem := range b.Elements {
		// Collect generic type params (e.g., V in Option: { V; Some(value V); ... }).
		if id, ok := elem.(*ast.Ident); ok && id.TokType == token.GENERIC_IDENT {
			sd.TypeParams = append(sd.TypeParams, id.Name)
			continue
		}
		vd, ok := elem.(*ast.VariantDef)
		if !ok {
			continue
		}
		vc := &IRVariantCase{Name: vd.Name}
		for _, p := range vd.Params {
			if p.Label == "" || p.Type == nil {
				continue
			}
			typ := l.goTypeFromTypeExpr(p.Type)
			vc.Fields = append(vc.Fields, &FieldSpec{Name: p.Label, Type: typ})
			if !fieldSet[p.Label] {
				fieldSet[p.Label] = true
				sd.Fields = append(sd.Fields, &FieldSpec{Name: p.Label, Type: typ})
			}
		}
		sd.Variants = append(sd.Variants, vc)
	}
	return sd
}

// ---------------------------------------------------------------------------
// main boc → FuncDecl
// ---------------------------------------------------------------------------

func (l *lowerer) lowerMainBoc(b *ast.BocLiteral) *FuncDecl {
	fn := &FuncDecl{Name: "main"}

	// Process statements in order, flushing concurrent boc-call groups
	// whenever a non-boc statement appears. This preserves ordering so that
	// local variables declared before a boc call are visible inside it.
	var pendingBocCalls []ast.Expr
	bgIdx := 0

	flushGroup := func() {
		if len(pendingBocCalls) == 0 {
			return
		}
		bgVar := fmt.Sprintf("_bg%d", bgIdx)
		bgIdx++
		fn.Body = append(fn.Body, &DeclStmt{Name: bgVar, Init: &NewGroupExpr{}})
		for _, call := range pendingBocCalls {
			callExpr := l.lowerExpr(call)
			fn.Body = append(fn.Body, &ExprStmt{
				Expr: &SpawnExpr{
					GroupVar: bgVar,
					Body:     []Stmt{&ReturnStmt{Value: &ForceExpr{Thunk: callExpr}}},
				},
			})
		}
		fn.Body = append(fn.Body, &WaitStmt{GroupVar: bgVar})
		pendingBocCalls = nil
	}

	for _, elem := range b.Elements {
		if expr, ok := elem.(ast.Expr); ok && l.isBocMethodCall(expr) {
			pendingBocCalls = append(pendingBocCalls, expr)
		} else {
			flushGroup()
			fn.Body = append(fn.Body, l.lowerMainStmt(elem)...)
		}
	}
	flushGroup()
	return fn
}

func (l *lowerer) lowerMainStmt(node ast.Node) []Stmt {
	switch e := node.(type) {
	case *ast.ShortDecl:
		var stmts []Stmt
		for i, n := range e.Names {
			var initExpr Expr
			if i < len(e.Values) {
				initExpr = l.lowerExpr(e.Values[i])
			}
			typ := ""
			if i < len(e.Values) {
				val := e.Values[i]
				if l.isBocMethodCall(val) {
					// RHS is a boc call — variable holds *Thunk[T].
					// Use := inference and mark for auto-forcing on use.
					l.thunkVars[n.Name] = true
				} else {
					typ = l.goTypeForVar(l.analyzer.ExprType(val))
				}
			}
			stmts = append(stmts, &DeclStmt{Name: n.Name, Type: typ, Init: initExpr})
		}
		return stmts
	case *ast.Assignment:
		return []Stmt{l.lowerAssignment(e)}
	case *ast.ReturnStmt:
		var val Expr
		if e.Value != nil {
			val = l.lowerExpr(e.Value)
		}
		return []Stmt{&ReturnStmt{Value: val}}
	case ast.Expr:
		if fs, ok := l.tryLowerWhile(e); ok {
			return []Stmt{fs}
		}
		if is, ok := l.tryLowerConditional(e); ok {
			return []Stmt{is}
		}
		if ms, ok := l.tryLowerMatch(e); ok {
			return []Stmt{ms}
		}
		return []Stmt{&ExprStmt{Expr: l.lowerExpr(e)}}
	default:
		return nil
	}
}

// ---------------------------------------------------------------------------
// BocWithSig
// ---------------------------------------------------------------------------

// lowerBocWithSigAsMethod lowers a BocWithSig that appears inside a singleton
// or struct boc body as a receiver method. Params come from the signature;
// the body is lowered with receiver context so field names resolve to self.xxx.
func (l *lowerer) lowerBocWithSigAsMethod(bws *ast.BocWithSig, recvType string, parentFields map[string]bool) *MethodDecl {
	bocSemType := l.analyzer.ExprType(bws)
	bt, _ := bocSemType.(*sema.BocType)

	// Use AST sig for return type to preserve generic TypeArgs.
	// Methods are always in shorthand form (BodyOnly=false).
	resultType := l.getResultTypeFromSig(bws.Sig, bt, false)

	// Build params using sema for ordering/filtering, AST types for generic TypeArgs.
	params := l.sigParams(bws.Sig, bt)

	prev := l.setReceiver("self", parentFields)
	body := l.lowerBocBody(bws.Body, resultType)
	l.restoreReceiver(prev)

	return &MethodDecl{
		RecvType: recvType,
		RecvName: "self",
		Name:     capitalize(bws.Name.Name),
		Params:   params,
		Results:  []string{"*std.Thunk[" + resultType + "]"},
		Body:     body,
	}
}

func (l *lowerer) lowerBocWithSig(bws *ast.BocWithSig, recvType string, parentFields map[string]bool) Decl {
	if bws.Body == nil {
		return nil
	}
	// Get the sema-inferred BocType (carries Params and Returns).
	bocSemType := l.analyzer.ExprType(bws)
	bt, _ := bocSemType.(*sema.BocType)

	// Use AST sig for return type to preserve generic TypeArgs (e.g. Option(V)).
	resultType := l.getResultTypeFromSig(bws.Sig, bt, bws.BodyOnly)
	thunkResult := "*std.Thunk[" + resultType + "]"

	// Collect generic type params from the sig (e.g. V in #(value V, V)).
	typeParams := collectSigTypeParams(bws.Sig)

	// Build Go params from the sig (shorthand) or body leading TypedDecls (body-only).
	var params []*ParamSpec
	if bws.BodyOnly && bws.Body != nil {
		// Body-only form (`name #(sig) = { body }`): body redeclares params as TypedDecls.
		for _, elem := range bws.Body.Elements {
			td, ok := elem.(*ast.TypedDecl)
			if !ok || td.Value != nil {
				break
			}
			params = append(params, &ParamSpec{
				Name: td.Name.Name,
				Type: l.goTypeFromTypeExpr(td.Type),
			})
		}
	} else {
		// Shorthand form (`name #(sig) { body }`): use sema params for ordering/filtering
		// but prefer AST types (which preserve generic TypeArgs like Option(V)).
		params = l.sigParams(bws.Sig, bt)
	}

	// Lower body as a boc (produces [ExprStmt{ThunkExpr}]).
	// Params are regular Go variables — no receiver context needed.
	prev := l.setReceiver("", nil)
	bocBodyStmts := l.lowerBocBody(bws.Body, resultType)
	l.restoreReceiver(prev)

	// lowerBocBody returns [ExprStmt{ThunkExpr}]; promote to ReturnStmt.
	var funcBody []Stmt
	if len(bocBodyStmts) == 1 {
		if es, ok := bocBodyStmts[0].(*ExprStmt); ok {
			funcBody = []Stmt{&ReturnStmt{Value: es.Expr}}
		}
	}
	if funcBody == nil {
		funcBody = bocBodyStmts
	}

	return &FuncDecl{
		Name:       bws.Name.Name,
		TypeParams: typeParams,
		Params:     params,
		Results:    []string{thunkResult},
		Body:       funcBody,
	}
}

// ---------------------------------------------------------------------------
// Expression lowering
// ---------------------------------------------------------------------------

func (l *lowerer) lowerExpr(e ast.Expr) Expr {
	if e == nil {
		return &UnitLit{}
	}
	switch expr := e.(type) {
	case *ast.IntLit:
		v, _ := strconv.ParseInt(expr.Value, 0, 64)
		return &IntLit{Val: v}
	case *ast.DecimalLit:
		v, _ := strconv.ParseFloat(expr.Value, 64)
		return &DecimalLit{Val: v}
	case *ast.StringLit:
		// Strip surrounding quotes and unescape.
		return &StringLit{Val: unquoteString(expr.Value)}
	case *ast.InterpolatedStringExpr:
		return l.lowerInterpString(expr)
	case *ast.Ident:
		if expr.Name == "true" {
			return &BoolLit{Val: true}
		}
		if expr.Name == "false" {
			return &BoolLit{Val: false}
		}
		name := l.lowerName(expr.Name)
		if l.thunkVars[expr.Name] {
			return &ForceExpr{Thunk: name}
		}
		return name
	case *ast.UnaryExpr:
		operand := l.lowerExpr(expr.Operand)
		return &MethodCall{Recv: operand, Method: "Neg", Args: nil}
	case *ast.BinaryExpr:
		left := l.lowerExpr(expr.Left)
		right := l.lowerExpr(expr.Right)
		method := goMethodName(sema.NonWordMethodName(expr.Op))
		return &MethodCall{Recv: left, Method: method, Args: []Expr{right}}
	case *ast.CallExpr:
		return l.lowerCall(expr)
	case *ast.MemberExpr:
		obj := l.lowerExpr(expr.Object)
		return &FieldAccess{Object: obj, Field: expr.Member.Name}
	case *ast.IndexExpr:
		obj := l.lowerExpr(expr.Object)
		idx := l.lowerExpr(expr.Index)
		return &IndexExpr{Object: obj, Index: idx}
	case *ast.GroupExpr:
		return l.lowerExpr(expr.Expr)
	case *ast.BocLiteral:
		return l.lowerBocLitExpr(expr)
	case *ast.ConditionalExpr:
		return l.lowerConditionalExpr(expr)
	case *ast.MatchExpr:
		return l.lowerMatchExpr(expr)
	case *ast.ArrayLiteral:
		return l.lowerArrayLit(expr)
	case *ast.DictLiteral:
		return l.lowerDictLit(expr)
	default:
		return &UnitLit{}
	}
}

// lowerName resolves an identifier: if it's a receiver field, emit self.field.
// Singleton boc vars are capitalized so they are exported for cross-package access.
// BocWithSig functions remain lowercase (they become Go functions, not vars).
func (l *lowerer) lowerName(name string) Expr {
	if l.recvName != "" && l.recvFields[name] {
		return &FieldAccess{Object: &Ident{Name: l.recvName}, Field: name}
	}
	// Capitalize singleton boc var references, but NOT BocWithSig functions.
	sym := l.analyzer.LookupInFile(name)
	if sym != nil {
		if _, isBoc := sym.Type.(*sema.BocType); isBoc {
			if _, isBWS := sym.Node.(*ast.BocWithSig); !isBWS {
				return &Ident{Name: capitalize(name)}
			}
		}
	}
	return &Ident{Name: name}
}

// builtinGoName maps a Yz builtin name to its Go runtime equivalent.
// Builtins are emitted as direct calls (not goroutine-wrapped).
var builtinGoName = map[string]string{
	"print": "std.Print",
	"info":  "std.Info",
}

// builtinSingleton maps a Yz built-in singleton name to its Go runtime
// receiver expression (e.g. "http" → "std.Http"). Method calls on these
// singletons are emitted as MethodCall with the runtime receiver and a
// Go-exported method name (first letter uppercased).
var builtinSingleton = map[string]string{
	"http": "std.Http",
}

// ---------------------------------------------------------------------------
// FQN (cross-package) call resolution
// ---------------------------------------------------------------------------

// tryLowerFQNCall detects and lowers cross-package FQN calls such as
// house.front.Host("Alice") → front.NewHost(std.NewString("Alice")).
// Returns (expr, true) if the callee is a known package FQN.
func (l *lowerer) tryLowerFQNCall(c *ast.CallExpr) (Expr, bool) {
	pkg, symName, ok := l.fqnCalleePackage(c.Callee)
	if !ok {
		return nil, false
	}
	exportedSym, ok := pkg.Exports[symName]
	if !ok {
		return nil, false
	}
	var args []Expr
	for _, arg := range c.Args {
		args = append(args, l.lowerExpr(arg.Value))
	}
	l.addImport(pkg.ImportPath)
	// Struct constructor: Host → pkg.NewHost(args)
	if _, isStruct := exportedSym.Type.(*sema.StructType); isStruct {
		return &FuncCall{Func: &Ident{Name: pkg.PkgAlias + ".New" + symName}, Args: args}, true
	}
	// Singleton boc or BocWithSig: pkg.func(args)
	return &FuncCall{Func: &Ident{Name: pkg.PkgAlias + "." + symName}, Args: args}, true
}

// fqnCalleePackage resolves a callee expression to (PackageType, symbolName).
// The callee must be a MemberExpr whose object chain resolves to a PackageType.
func (l *lowerer) fqnCalleePackage(callee ast.Expr) (*sema.PackageType, string, bool) {
	mem, ok := callee.(*ast.MemberExpr)
	if !ok {
		return nil, "", false
	}
	pkg := l.resolveExprToPackage(mem.Object)
	if pkg == nil {
		return nil, "", false
	}
	return pkg, mem.Member.Name, true
}

// resolveExprToPackage walks an AST expression and returns the PackageType it
// resolves to, if any.
func (l *lowerer) resolveExprToPackage(expr ast.Expr) *sema.PackageType {
	switch e := expr.(type) {
	case *ast.Ident:
		sym := l.analyzer.LookupInFile(e.Name)
		if sym == nil {
			return nil
		}
		if pt, ok := sym.Type.(*sema.PackageType); ok {
			return pt
		}
		return nil
	case *ast.MemberExpr:
		ns := l.resolveExprToNamespace(e.Object)
		if ns == nil {
			return nil
		}
		child, ok := ns.Children[e.Member.Name]
		if !ok {
			return nil
		}
		if pt, ok := child.Type.(*sema.PackageType); ok {
			return pt
		}
		return nil
	}
	return nil
}

// resolveExprToNamespace walks an AST expression and returns the NamespaceType
// it resolves to, if any.
func (l *lowerer) resolveExprToNamespace(expr ast.Expr) *sema.NamespaceType {
	id, ok := expr.(*ast.Ident)
	if !ok {
		return nil
	}
	sym := l.analyzer.LookupInFile(id.Name)
	if sym == nil {
		return nil
	}
	if ns, ok := sym.Type.(*sema.NamespaceType); ok {
		return ns
	}
	return nil
}

// tryLowerCrossPackageSingletonMethod detects and lowers calls of the form
// pkg.singleton.method(args) → pkg.Singleton.Method(args).
// The singleton must be exported from the package as a BocType symbol.
func (l *lowerer) tryLowerCrossPackageSingletonMethod(c *ast.CallExpr) (Expr, bool) {
	outer, ok := c.Callee.(*ast.MemberExpr)
	if !ok {
		return nil, false
	}
	inner, ok := outer.Object.(*ast.MemberExpr)
	if !ok {
		return nil, false
	}
	pkg := l.resolveExprToPackage(inner.Object)
	if pkg == nil {
		return nil, false
	}
	singletonName := inner.Member.Name
	exportedSym, ok := pkg.Exports[singletonName]
	if !ok {
		return nil, false
	}
	if _, isBoc := exportedSym.Type.(*sema.BocType); !isBoc {
		return nil, false
	}
	var args []Expr
	for _, arg := range c.Args {
		args = append(args, l.lowerExpr(arg.Value))
	}
	l.addImport(pkg.ImportPath)
	recv := &FieldAccess{Object: &Ident{Name: pkg.PkgAlias}, Field: capitalize(singletonName)}
	return &MethodCall{Recv: recv, Method: capitalize(outer.Member.Name), Args: args}, true
}

func (l *lowerer) lowerCall(c *ast.CallExpr) Expr {
	// Check for known builtins first — they emit as direct std.Xxx calls.
	if id, ok := c.Callee.(*ast.Ident); ok {
		if goName, isBuiltin := builtinGoName[id.Name]; isBuiltin {
			return l.lowerBuiltinCall(goName, c)
		}
	}

	// Cross-package singleton method call: pkg.counter.increment() → pkg.Counter.Increment()
	if result, ok := l.tryLowerCrossPackageSingletonMethod(c); ok {
		return result
	}

	// Cross-package FQN call: house.front.Host("Alice") → front.NewHost(...)
	if result, ok := l.tryLowerFQNCall(c); ok {
		return result
	}

	// array.map(boc) → std.ArrayMap(array, closure).
	// Go methods cannot introduce new type params, so Map must be a package-level function.
	if mem, ok := c.Callee.(*ast.MemberExpr); ok && mem.Member.Name == "map" {
		if _, isArray := l.analyzer.ExprType(mem.Object).(*sema.ArrayType); isArray {
			recv := l.lowerExpr(mem.Object)
			mapArgs := []Expr{recv}
			for _, arg := range c.Args {
				mapArgs = append(mapArgs, l.lowerExpr(arg.Value))
			}
			return &FuncCall{Func: &Ident{Name: "std.ArrayMap"}, Args: mapArgs}
		}
	}

	callee := l.lowerExpr(c.Callee)
	var args []Expr
	for _, arg := range c.Args {
		args = append(args, l.lowerExpr(arg.Value))
	}

	// Check for built-in singleton method calls: http.get("url") → std.Http.Get(url).
	// The runtime method already returns *Thunk[T], so no extra wrapping is needed.
	if fa, ok := callee.(*FieldAccess); ok {
		if recv, isIdent := fa.Object.(*Ident); isIdent {
			if goRef, isSingleton := builtinSingleton[recv.Name]; isSingleton {
				method := fa.Field
				if len(method) > 0 {
					method = strings.ToUpper(method[:1]) + method[1:]
				}
				return &MethodCall{Recv: &Ident{Name: goRef}, Method: method, Args: args}
			}
		}
	}

	// For a field-access callee like counter.increment or counter.value():
	// the method itself already wraps its body in std.Go and returns *Thunk[T].
	// We just emit the direct MethodCall — no extra wrapping needed.
	// User-defined method names are capitalized to be exported (cross-package safe).
	if fa, ok := callee.(*FieldAccess); ok {
		return &MethodCall{Recv: fa.Object, Method: capitalize(fa.Field), Args: args}
	}

	if id, ok := c.Callee.(*ast.Ident); ok {
		sym := l.analyzer.LookupInFile(id.Name)
		if sym != nil {
			// Variant constructor call: Cat("Whiskers", 9) → NewPetCat(...)
			if sym.ParentTypeName != "" {
				return &FuncCall{Func: &Ident{Name: "New" + sym.ParentTypeName + id.Name}, Args: args}
			}
			// Struct type constructor call: Named("Alice") → NewNamed(args)
			if _, isStruct := sym.Type.(*sema.StructType); isStruct {
				return &FuncCall{Func: &Ident{Name: "New" + id.Name}, Args: args}
			}
			// BocWithSig functions already return *Thunk[T] — emit a plain call.
			// Inject default argument values for any omitted optional params.
			if bws, isBWS := sym.Node.(*ast.BocWithSig); isBWS {
				args = l.fillDefaults(args, bws.Sig.Params)
				return &FuncCall{Func: callee, Args: args}
			}
		}
	}

	// Other boc identifier calls — goroutine-launched thunks.
	calleeType := l.analyzer.ExprType(c.Callee)
	if bt, ok := calleeType.(*sema.BocType); ok {
		resultType := "std.Unit"
		if len(bt.Returns) > 0 {
			resultType = l.goType(bt.Returns[0])
		}
		callExpr := &FuncCall{Func: callee, Args: args}
		return &ThunkExpr{
			ResultType: resultType,
			Body:       []Stmt{&ReturnStmt{Value: callExpr}},
			Spawn:      true,
		}
	}

	return &FuncCall{Func: callee, Args: args}
}

// isBocMethodCall reports whether an AST expression is a call that returns a
// *Thunk in Go — either a singleton boc method call (counter.value()) or a
// direct BocWithSig function call (greet("Alice")).
func (l *lowerer) isBocMethodCall(e ast.Expr) bool {
	c, ok := e.(*ast.CallExpr)
	if !ok {
		return false
	}
	// Boc/struct method call: obj.method()
	if mem, ok := c.Callee.(*ast.MemberExpr); ok {
		// Cross-package singleton method call: pkg.singleton.method()
		if innerMem, ok := mem.Object.(*ast.MemberExpr); ok {
			if pkg := l.resolveExprToPackage(innerMem.Object); pkg != nil {
				if exportedSym, ok := pkg.Exports[innerMem.Member.Name]; ok {
					if _, isBoc := exportedSym.Type.(*sema.BocType); isBoc {
						return true
					}
				}
			}
		}
		// Struct instance method call (n.hi() where n is *Named) → returns *Thunk.
		if _, isStruct := l.analyzer.ExprType(mem.Object).(*sema.StructType); isStruct {
			return true
		}
		// Singleton boc method call (counter.increment()) → returns *Thunk.
		objIdent, ok := mem.Object.(*ast.Ident)
		if !ok {
			return false
		}
		sym := l.analyzer.LookupInFile(objIdent.Name)
		if sym == nil {
			return false
		}
		_, isBoc := sym.Type.(*sema.BocType)
		return isBoc
	}
	// Direct BocWithSig call: greet("Alice")
	if id, ok := c.Callee.(*ast.Ident); ok {
		sym := l.analyzer.LookupInFile(id.Name)
		if sym != nil {
			if _, isBWS := sym.Node.(*ast.BocWithSig); isBWS {
				return true
			}
		}
	}
	return false
}

// isSingletonBoc reports whether expr is an Ident that refers to a
// singleton boc (a lowercase boc defined at file scope).
func (l *lowerer) isSingletonBoc(expr Expr) bool {
	id, ok := expr.(*Ident)
	if !ok {
		return false
	}
	sym := l.analyzer.LookupInFile(id.Name)
	if sym == nil {
		return false
	}
	_, isBoc := sym.Type.(*sema.BocType)
	return isBoc
}

// tryLowerWhile detects a while({cond},{body}) call and lowers it to a
// native ForStmt rather than a runtime call. Returns (stmt, true) on match.
func (l *lowerer) tryLowerWhile(e ast.Expr) (Stmt, bool) {
	call, ok := e.(*ast.CallExpr)
	if !ok {
		return nil, false
	}
	id, ok := call.Callee.(*ast.Ident)
	if !ok || id.Name != "while" {
		return nil, false
	}
	if len(call.Args) != 2 {
		return nil, false
	}
	condBoc, ok := call.Args[0].Value.(*ast.BocLiteral)
	if !ok {
		return nil, false
	}
	bodyBoc, ok := call.Args[1].Value.(*ast.BocLiteral)
	if !ok {
		return nil, false
	}
	return &ForStmt{
		Cond: l.lowerBocAsExpr(condBoc),
		Body: l.lowerBocAsStmts(bodyBoc),
	}, true
}

// tryLowerConditional detects a ConditionalExpr and lowers it to an IfStmt.
// Returns (stmt, true) on match.
func (l *lowerer) tryLowerConditional(e ast.Expr) (Stmt, bool) {
	cond, ok := e.(*ast.ConditionalExpr)
	if !ok {
		return nil, false
	}
	condExpr := l.lowerExpr(cond.Cond)
	var thenStmts, elseStmts []Stmt
	if b, ok := cond.TrueCase.(*ast.BocLiteral); ok {
		thenStmts = l.lowerBocAsStmts(b)
	} else {
		thenStmts = []Stmt{&ExprStmt{Expr: l.lowerExpr(cond.TrueCase)}}
	}
	if b, ok := cond.FalseCase.(*ast.BocLiteral); ok {
		elseStmts = l.lowerBocAsStmts(b)
	} else {
		elseStmts = []Stmt{&ExprStmt{Expr: l.lowerExpr(cond.FalseCase)}}
	}
	return &IfStmt{Cond: condExpr, Then: thenStmts, Else: elseStmts}, true
}

// lowerConditionalExpr lowers a ConditionalExpr used in expression position
// as an immediately-invoked closure (IIFE) with an if/else inside.
func (l *lowerer) lowerConditionalExpr(cond *ast.ConditionalExpr) Expr {
	// Determine result type from the true-case boc.
	semType := l.analyzer.ExprType(cond.TrueCase)
	resultType := "std.Unit"
	if bt, ok := semType.(*sema.BocType); ok && len(bt.Returns) > 0 {
		resultType = l.goType(bt.Returns[0])
	}

	condExpr := l.lowerExpr(cond.Cond)

	var trueBody []Stmt
	if tc, ok := cond.TrueCase.(*ast.BocLiteral); ok {
		trueBody = l.lowerMatchArmBody(tc.Elements, resultType)
	}

	var falseBody []Stmt
	if fc, ok := cond.FalseCase.(*ast.BocLiteral); ok {
		falseBody = l.lowerMatchArmBody(fc.Elements, resultType)
	}

	return &MatchExpr{
		ResultType: resultType,
		Arms: []*MatchArm{
			{Cond: condExpr, Body: trueBody},
			{Cond: nil, Body: falseBody},
		},
	}
}

// ---------------------------------------------------------------------------
// Match expression lowering
// ---------------------------------------------------------------------------

// lowerMatchExpr lowers a MatchExpr for expression position (IIFE).
func (l *lowerer) lowerMatchExpr(m *ast.MatchExpr) Expr {
	// Discriminant match in expression position → SwitchExpr IIFE.
	if m.Subject != nil {
		if sw, ok := l.tryLowerDiscriminantMatchExpr(m); ok {
			return sw
		}
	}

	semType := l.analyzer.ExprType(m)
	resultType := l.goType(semType)

	var arms []*MatchArm
	for _, arm := range m.Arms {
		var cond Expr
		if arm.Condition != nil {
			cond = l.lowerExpr(arm.Condition)
		}
		body := l.lowerMatchArmBody(arm.Body, resultType)
		arms = append(arms, &MatchArm{Cond: cond, Body: body})
	}
	return &MatchExpr{ResultType: resultType, Arms: arms}
}

// tryLowerMatch detects a MatchExpr and lowers it to a statement.
// Discriminant match (Subject != nil) → SwitchStmt.
// Condition match (Subject == nil) → IfStmt chain.
func (l *lowerer) tryLowerMatch(e ast.Expr) (Stmt, bool) {
	m, ok := e.(*ast.MatchExpr)
	if !ok {
		return nil, false
	}
	if len(m.Arms) == 0 {
		return nil, false
	}

	// Discriminant match: `match expr { Variant => body }`
	if m.Subject != nil {
		if sw, ok := l.tryLowerDiscriminantMatch(m); ok {
			return sw, true
		}
	}

	// Build chain from last arm to first.
	var pendingElse []Stmt // else body for the arm above
	var topIf *IfStmt
	for i := len(m.Arms) - 1; i >= 0; i-- {
		arm := m.Arms[i]
		body := l.lowerBocAsStmts2(arm.Body)
		if arm.Condition == nil {
			// Default arm — becomes the else body of the arm above.
			pendingElse = body
		} else {
			cond := l.lowerExpr(arm.Condition)
			st := &IfStmt{Cond: cond, Then: body, Else: pendingElse}
			pendingElse = []Stmt{st}
			topIf = st
		}
	}
	if topIf == nil {
		// All arms were default — just emit the body as statements directly.
		// Wrap in a trivially true if so we still return a Stmt.
		return &IfStmt{Cond: &BoolLit{Val: true}, Then: pendingElse}, true
	}
	return topIf, true
}

// ---------------------------------------------------------------------------
// Discriminant (variant) match lowering
// ---------------------------------------------------------------------------

// tryLowerDiscriminantMatch lowers `match subject { Variant => body }` to
// a SwitchStmt. Returns (stmt, true) when all arms have TYPE_IDENT conditions
// and the subject is a variant struct type.
func (l *lowerer) tryLowerDiscriminantMatch(m *ast.MatchExpr) (Stmt, bool) {
	typeName, ok := l.variantDiscriminantType(m.Subject)
	if !ok {
		return nil, false
	}
	subject := l.lowerExpr(m.Subject)
	sw := &SwitchStmt{Subject: subject, TypeName: typeName}
	for _, arm := range m.Arms {
		variantName, ok := l.armVariantName(arm)
		if !ok {
			return nil, false
		}
		constName := "_" + l.variantStructName(m.Subject) + variantName
		body := l.lowerBocAsStmts2(arm.Body)
		sw.Cases = append(sw.Cases, &SwitchCase{ConstName: constName, Body: body})
	}
	return sw, true
}

// tryLowerDiscriminantMatchExpr lowers discriminant match in expression position.
func (l *lowerer) tryLowerDiscriminantMatchExpr(m *ast.MatchExpr) (Expr, bool) {
	typeName, ok := l.variantDiscriminantType(m.Subject)
	if !ok {
		return nil, false
	}
	subject := l.lowerExpr(m.Subject)
	semType := l.analyzer.ExprType(m)
	resultType := l.goType(semType)
	sw := &SwitchExpr{Subject: subject, ResultType: resultType}
	_ = typeName
	for _, arm := range m.Arms {
		variantName, ok := l.armVariantName(arm)
		if !ok {
			return nil, false
		}
		constName := "_" + l.variantStructName(m.Subject) + variantName
		body := l.lowerMatchArmBody(arm.Body, resultType)
		sw.Cases = append(sw.Cases, &SwitchCase{ConstName: constName, Body: body})
	}
	return sw, true
}

// variantDiscriminantType returns the Go discriminant type name for the subject
// expression (e.g. "_PetVariant") if the subject is a variant struct, else "".
func (l *lowerer) variantDiscriminantType(subject ast.Expr) (string, bool) {
	semType := l.analyzer.ExprType(subject)
	st, ok := semType.(*sema.StructType)
	if !ok || !st.IsVariant || st.Name == "" {
		return "", false
	}
	return "_" + st.Name + "Variant", true
}

// variantStructName returns the struct name from a subject expression.
func (l *lowerer) variantStructName(subject ast.Expr) string {
	semType := l.analyzer.ExprType(subject)
	if st, ok := semType.(*sema.StructType); ok {
		return st.Name
	}
	return ""
}

// armVariantName extracts the variant name from a match arm condition.
// The condition must be a TYPE_IDENT (e.g. Cat, Dog).
func (l *lowerer) armVariantName(arm *ast.ConditionalBoc) (string, bool) {
	if arm.Condition == nil {
		return "", false
	}
	id, ok := arm.Condition.(*ast.Ident)
	if !ok {
		return "", false
	}
	if id.TokType != token.TYPE_IDENT {
		return "", false
	}
	return id.Name, true
}

// lowerMatchArmBody lowers a match arm's body elements for expression-position
// match (IIFE). The last expression becomes a ReturnStmt.
func (l *lowerer) lowerMatchArmBody(elements []ast.Node, resultType string) []Stmt {
	var stmts []Stmt
	for i, elem := range elements {
		isLast := i == len(elements)-1
		switch e := elem.(type) {
		case *ast.Assignment:
			stmts = append(stmts, l.lowerAssignment(e))
		case *ast.ShortDecl:
			stmts = append(stmts, l.lowerBodyShortDecl(e, isLast, resultType))
		case *ast.ReturnStmt:
			var val Expr
			if e.Value != nil {
				val = l.lowerExpr(e.Value)
			}
			stmts = append(stmts, &ReturnStmt{Value: val})
		case ast.Expr:
			expr := l.lowerExpr(e)
			if isLast {
				stmts = append(stmts, &ReturnStmt{Value: expr})
			} else {
				stmts = append(stmts, &ExprStmt{Expr: expr})
			}
		}
	}
	if len(stmts) == 0 || !isReturnStmt(stmts[len(stmts)-1]) {
		stmts = append(stmts, &ReturnStmt{Value: &UnitLit{}})
	}
	return stmts
}

// lowerBocAsStmts2 lowers a slice of ast.Node elements (a match arm body) as
// flat statements with no return-wrapping. Used for statement-position match.
func (l *lowerer) lowerBocAsStmts2(elements []ast.Node) []Stmt {
	var stmts []Stmt
	for _, elem := range elements {
		switch e := elem.(type) {
		case *ast.Assignment:
			stmts = append(stmts, l.lowerAssignment(e))
		case *ast.ShortDecl:
			stmts = append(stmts, l.lowerBodyShortDecl(e, false, "std.Unit"))
		case ast.Expr:
			if fs, ok := l.tryLowerWhile(e); ok {
				stmts = append(stmts, fs)
			} else if is, ok := l.tryLowerConditional(e); ok {
				stmts = append(stmts, is)
			} else if ms, ok := l.tryLowerMatch(e); ok {
				stmts = append(stmts, ms)
			} else {
				stmts = append(stmts, &ExprStmt{Expr: l.lowerExpr(e)})
			}
		}
	}
	return stmts
}

// lowerBocAsExpr extracts and lowers the primary expression from a boc
// literal, for use as a for-loop condition.
func (l *lowerer) lowerBocAsExpr(b *ast.BocLiteral) Expr {
	for _, elem := range b.Elements {
		if e, ok := elem.(ast.Expr); ok {
			return l.lowerExpr(e)
		}
	}
	return &BoolLit{Val: true} // fallback: infinite loop
}

// lowerBocAsStmts lowers boc elements as a flat list of statements without
// goroutine wrapping. Used for while loop bodies.
func (l *lowerer) lowerBocAsStmts(b *ast.BocLiteral) []Stmt {
	var stmts []Stmt
	for _, elem := range b.Elements {
		switch e := elem.(type) {
		case *ast.Assignment:
			stmts = append(stmts, l.lowerAssignment(e))
		case *ast.ShortDecl:
			stmts = append(stmts, l.lowerBodyShortDecl(e, false, "std.Unit"))
		case ast.Expr:
			if fs, ok := l.tryLowerWhile(e); ok {
				stmts = append(stmts, fs)
			} else if is, ok := l.tryLowerConditional(e); ok {
				stmts = append(stmts, is)
			} else if ms, ok := l.tryLowerMatch(e); ok {
				stmts = append(stmts, ms)
			} else {
				stmts = append(stmts, &ExprStmt{Expr: l.lowerExpr(e)})
			}
		}
	}
	return stmts
}

// lowerBuiltinCall emits a direct std.Xxx(...) call.
// For while: the boc arguments are emitted as closures, not thunks.
// For print/info: arguments are passed directly (thunks are force-materialized if needed).
func (l *lowerer) lowerBuiltinCall(goName string, c *ast.CallExpr) Expr {
	var args []Expr
	for _, arg := range c.Args {
		val := arg.Value
		if _, isBoc := val.(*ast.BocLiteral); isBoc {
			// Boc literal arg to a builtin → closure (not a goroutine-thunk).
			args = append(args, l.lowerExpr(val))
		} else {
			// For print/info: if the arg is a boc method call (returns *Thunk),
			// force it so the actual value is passed, not the thunk pointer.
			expr := l.lowerExpr(val)
			if l.isBocMethodCall(val) {
				expr = &ForceExpr{Thunk: expr}
			}
			args = append(args, expr)
		}
	}
	return &FuncCall{Func: &Ident{Name: goName}, Args: args}
}

func (l *lowerer) lowerBocLitExpr(b *ast.BocLiteral) Expr {
	// Extract TypedDecl params (nil-value TypedDecls in the body).
	var params []*ParamSpec
	for _, elem := range b.Elements {
		if td, ok := elem.(*ast.TypedDecl); ok && td.Value == nil {
			params = append(params, &ParamSpec{
				Name: td.Name.Name,
				Type: l.goTypeFromTypeExpr(td.Type),
			})
		}
	}
	semType := l.analyzer.ExprType(b)
	resultType := "std.Unit"
	if bt, ok := semType.(*sema.BocType); ok && len(bt.Returns) > 0 {
		resultType = l.goType(bt.Returns[0])
	}
	body := l.lowerClosureBody(b.Elements, resultType)
	return &ClosureExpr{Params: params, ResultType: resultType, Body: body}
}

// lowerClosureBody lowers boc elements as a synchronous closure body.
// Unlike lowerBocBody, it does NOT wrap in a ThunkExpr.
// TypedDecl with nil Value (params) are skipped — already captured in ClosureExpr.Params.
func (l *lowerer) lowerClosureBody(elements []ast.Node, resultType string) []Stmt {
	var stmts []Stmt
	for i, elem := range elements {
		isLast := i == len(elements)-1
		switch e := elem.(type) {
		case *ast.TypedDecl:
			if e.Value == nil {
				continue // param — already in ClosureExpr.Params
			}
			stmts = append(stmts, &DeclStmt{
				Name: e.Name.Name,
				Init: l.lowerExpr(e.Value),
			})
		case *ast.ShortDecl:
			stmts = append(stmts, l.lowerBodyShortDecl(e, isLast, resultType))
		case *ast.Assignment:
			stmts = append(stmts, l.lowerAssignment(e))
		case *ast.ReturnStmt:
			var val Expr
			if e.Value != nil {
				val = l.lowerExpr(e.Value)
			}
			stmts = append(stmts, &ReturnStmt{Value: val})
		case ast.Expr:
			if fs, ok := l.tryLowerWhile(e); ok {
				stmts = append(stmts, fs)
				if isLast {
					stmts = append(stmts, &ReturnStmt{Value: &UnitLit{}})
				}
			} else if is, ok := l.tryLowerConditional(e); ok {
				stmts = append(stmts, is)
				if isLast {
					stmts = append(stmts, &ReturnStmt{Value: &UnitLit{}})
				}
			} else if isLast {
				expr := l.lowerExpr(e)
				if l.isBocMethodCall(e) {
					expr = &ForceExpr{Thunk: expr}
				}
				stmts = append(stmts, &ReturnStmt{Value: expr})
			} else {
				if ms, ok := l.tryLowerMatch(e); ok {
					stmts = append(stmts, ms)
				} else {
					stmts = append(stmts, &ExprStmt{Expr: l.lowerExpr(e)})
				}
			}
		}
	}
	if len(stmts) == 0 || !isReturnStmt(stmts[len(stmts)-1]) {
		stmts = append(stmts, &ReturnStmt{Value: &UnitLit{}})
	}
	return stmts
}

// lowerInterpString lowers an InterpolatedStringExpr to a chain of Plus calls:
//   "Hello, `name`!" → std.NewString("Hello, ").Plus(std.NewString(std.Stringify(name))).Plus(std.NewString("!"))
func (l *lowerer) lowerInterpString(e *ast.InterpolatedStringExpr) Expr {
	var result Expr
	for _, part := range e.Parts {
		var node Expr
		if part.IsExpr {
			inner := l.lowerExpr(part.Expr)
			// Force thunks (boc method calls return *Thunk[T]).
			if l.isBocMethodCall(part.Expr) {
				inner = &ForceExpr{Thunk: inner}
			}
			// std.NewString(std.Stringify(inner))
			node = &FuncCall{
				Func: &Ident{Name: "std.NewString"},
				Args: []Expr{&FuncCall{
					Func: &Ident{Name: "std.Stringify"},
					Args: []Expr{inner},
				}},
			}
		} else {
			node = &StringLit{Val: unquoteString(`"` + part.Text + `"`)}
		}
		if result == nil {
			result = node
		} else {
			result = &MethodCall{Recv: result, Method: "Plus", Args: []Expr{node}}
		}
	}
	if result == nil {
		result = &StringLit{Val: ""}
	}
	return result
}

func (l *lowerer) lowerArrayLit(arr *ast.ArrayLiteral) Expr {
	// Emit std.NewArray(elem0, elem1, ...) — represented as FuncCall.
	var args []Expr
	for _, el := range arr.Elements {
		args = append(args, l.lowerExpr(el))
	}
	return &FuncCall{
		Func: &Ident{Name: "std.NewArray"},
		Args: args,
	}
}

func (l *lowerer) lowerDictLit(d *ast.DictLiteral) Expr {
	// Determine key and value Go types from the sema type.
	keyType, valType := "any", "any"
	if dt, ok := l.analyzer.ExprType(d).(*sema.DictType); ok {
		keyType = l.goType(dt.Key)
		valType = l.goType(dt.Val)
	} else if d.KeyType != nil {
		// Empty-type form [K:V]() — use the explicit type expressions.
		keyType = l.goTypeFromTypeExpr(d.KeyType)
		valType = l.goTypeFromTypeExpr(d.ValType)
	}

	// Base: std.NewDict[K, V]()
	var result Expr = &FuncCall{
		Func: &Ident{Name: "std.NewDict[" + keyType + ", " + valType + "]"},
	}

	// Chain .Set(k, v) for each entry.
	for _, entry := range d.Entries {
		result = &MethodCall{
			Recv:   result,
			Method: "Set",
			Args:   []Expr{l.lowerExpr(entry.Key), l.lowerExpr(entry.Value)},
		}
	}
	return result
}

// ---------------------------------------------------------------------------
// Receiver context helpers
// ---------------------------------------------------------------------------

type receiverState struct {
	name   string
	fields map[string]bool
}

func (l *lowerer) setReceiver(name string, fields map[string]bool) receiverState {
	prev := receiverState{name: l.recvName, fields: l.recvFields}
	l.recvName = name
	l.recvFields = fields
	return prev
}

func (l *lowerer) restoreReceiver(prev receiverState) {
	l.recvName = prev.name
	l.recvFields = prev.fields
}

// ---------------------------------------------------------------------------
// Generic type-param helpers
// ---------------------------------------------------------------------------

// sigParams builds the Go parameter list for a BocWithSig (shorthand or method form).
// It uses sema BocType for ordering and filtering (handles shortdecl/default params
// that have no explicit AST type), but prefers AST type expressions when they carry
// generic TypeArgs (e.g. Option(V) → *Option[V]).
func (l *lowerer) sigParams(sig *ast.BocTypeExpr, bt *sema.BocType) []*ParamSpec {
	// Index AST params by label for fast lookup.
	astParamByLabel := map[string]ast.TypeExpr{}
	if sig != nil {
		for _, p := range sig.Params {
			if p.Label != "" && p.Type != nil {
				astParamByLabel[p.Label] = p.Type
			}
		}
	}
	var params []*ParamSpec
	if bt != nil {
		for _, p := range bt.Params {
			if p.IsReturn || p.Label == "" {
				continue
			}
			goTyp := l.goType(p.Type) // default: sema-resolved type
			if astTE, ok := astParamByLabel[p.Label]; ok {
				goTyp = l.goTypeFromTypeExpr(astTE) // prefer AST (preserves TypeArgs)
			}
			params = append(params, &ParamSpec{Name: p.Label, Type: goTyp})
		}
	}
	return params
}

// collectSigTypeParams scans a BocTypeExpr's param type expressions and
// returns the names of all GENERIC_IDENTs (single-letter type params) used.
// Order is the order of first appearance; duplicates are deduplicated.
func collectSigTypeParams(sig *ast.BocTypeExpr) []string {
	if sig == nil {
		return nil
	}
	seen := map[string]bool{}
	var result []string
	for _, p := range sig.Params {
		if p.Type != nil {
			collectGenericIdentsFromType(p.Type, seen, &result)
		}
	}
	return result
}

// collectGenericIdentsFromType recursively collects GENERIC_IDENT names from a TypeExpr.
func collectGenericIdentsFromType(te ast.TypeExpr, seen map[string]bool, result *[]string) {
	if te == nil {
		return
	}
	switch t := te.(type) {
	case *ast.SimpleTypeExpr:
		if t.TokType == token.GENERIC_IDENT && !seen[t.Name] {
			seen[t.Name] = true
			*result = append(*result, t.Name)
		}
		for _, arg := range t.TypeArgs {
			collectGenericIdentsFromType(arg, seen, result)
		}
	case *ast.ArrayTypeExpr:
		collectGenericIdentsFromType(t.ElemType, seen, result)
	case *ast.DictTypeExpr:
		collectGenericIdentsFromType(t.KeyType, seen, result)
		collectGenericIdentsFromType(t.ValType, seen, result)
	}
}

// getResultTypeFromSig extracts the return Go type string from a BocWithSig.
// It prefers the explicit return type annotation in the AST sig (which preserves
// generic TypeArgs) over the sema-inferred return type.
// bodyOnly must be true for the `name #(sig) = { body }` form, where ALL sig
// params are inputs and there are no return-type annotations.
func (l *lowerer) getResultTypeFromSig(sig *ast.BocTypeExpr, bt *sema.BocType, bodyOnly bool) string {
	if sig != nil && !bodyOnly {
		for _, p := range sig.Params {
			if p.Label == "" && p.Type != nil {
				return l.goTypeFromTypeExpr(p.Type)
			}
		}
	}
	if bt != nil && len(bt.Returns) > 0 {
		return l.goType(bt.Returns[0])
	}
	return "std.Unit"
}

// ---------------------------------------------------------------------------
// Type helpers
// ---------------------------------------------------------------------------

// goType converts a sema.Type to a Go type string.
func (l *lowerer) goType(t sema.Type) string {
	if t == nil {
		return "any"
	}
	switch tt := t.(type) {
	case *sema.BuiltinType:
		switch tt.String() {
		case "Int":
			return "std.Int"
		case "Decimal":
			return "std.Decimal"
		case "String":
			return "std.String"
		case "Bool":
			return "std.Bool"
		case "Unit":
			return "std.Unit"
		}
	case *sema.BocType:
		if len(tt.Returns) == 1 {
			inner := l.goType(tt.Returns[0])
			return fmt.Sprintf("*std.Thunk[%s]", inner)
		}
		return "*std.Thunk[any]"
	case *sema.StructType:
		if tt.Name != "" {
			if tt.IsInterface {
				return tt.Name // Go interfaces are already reference types
			}
			return "*" + tt.Name
		}
		return "any"
	case *sema.ArrayType:
		return fmt.Sprintf("std.Array[%s]", l.goType(tt.Elem))
	case *sema.DictType:
		return fmt.Sprintf("std.Dict[%s, %s]", l.goType(tt.Key), l.goType(tt.Val))
	case *sema.ThunkType:
		return fmt.Sprintf("*std.Thunk[%s]", l.goType(tt.Inner))
	case *sema.GenericType:
		return tt.Name
	case *sema.UnknownType:
		return "any"
	}
	return "any"
}

// goTypeForVar returns the Go type string for a variable declaration.
// Returns "" (triggering := inference) when the type cannot be expressed
// precisely — e.g. an ArrayType whose element type is unknown (map result),
// or a generic struct whose type args are resolved by Go.
func (l *lowerer) goTypeForVar(t sema.Type) string {
	if at, ok := t.(*sema.ArrayType); ok {
		if _, isUnknown := at.Elem.(*sema.UnknownType); isUnknown {
			return ""
		}
	}
	if st, ok := t.(*sema.StructType); ok && len(st.TypeParams) > 0 {
		return ""
	}
	return l.goType(t)
}

// goTypeFromTypeExpr converts an ast.TypeExpr to a Go type string.
func (l *lowerer) goTypeFromTypeExpr(te ast.TypeExpr) string {
	if te == nil {
		return "any"
	}
	switch t := te.(type) {
	case *ast.SimpleTypeExpr:
		if t.TokType == token.GENERIC_IDENT {
			return t.Name // generic type param (e.g., V) — no pointer
		}
		// Generic application: Option(T) → *Option[T], Option(String) → *Option[std.String]
		if len(t.TypeArgs) > 0 {
			var args []string
			for _, arg := range t.TypeArgs {
				args = append(args, l.goTypeFromTypeExpr(arg))
			}
			return "*" + t.Name + "[" + strings.Join(args, ", ") + "]"
		}
		switch t.Name {
		case "Int":
			return "std.Int"
		case "Decimal":
			return "std.Decimal"
		case "String":
			return "std.String"
		case "Bool":
			return "std.Bool"
		case "Unit":
			return "std.Unit"
		default:
			// Go interfaces are already reference types — no pointer needed.
			sym := l.analyzer.LookupInFile(t.Name)
			if sym != nil {
				if st, ok := sym.Type.(*sema.StructType); ok && st.IsInterface {
					return t.Name
				}
			}
			return "*" + t.Name
		}
	case *ast.ArrayTypeExpr:
		return fmt.Sprintf("std.Array[%s]", l.goTypeFromTypeExpr(t.ElemType))
	case *ast.DictTypeExpr:
		return fmt.Sprintf("std.Dict[%s, %s]",
			l.goTypeFromTypeExpr(t.KeyType),
			l.goTypeFromTypeExpr(t.ValType))
	case *ast.BocTypeExpr:
		return "any" // simplified — boc type exprs as function types handled in codegen
	}
	return "any"
}

// goMethodName converts a sema symbol-style method name (e.g. "plus") to the
// exported Go method name on the yzrt type (e.g. "Plus").
func goMethodName(symbolName string) string {
	if symbolName == "" {
		return ""
	}
	return strings.ToUpper(symbolName[:1]) + symbolName[1:]
}

// ---------------------------------------------------------------------------
// Utility helpers
// ---------------------------------------------------------------------------

func isUppercase(name string) bool {
	if name == "" {
		return false
	}
	return unicode.IsUpper(rune(name[0]))
}

// capitalize uppercases the first letter of name, leaving the rest unchanged.
// Used to produce exported Go identifiers from Yz (lowercase-first) names.
func capitalize(name string) string {
	if name == "" {
		return name
	}
	return strings.ToUpper(name[:1]) + name[1:]
}

func (l *lowerer) valueAt(vals []ast.Expr, i int) ast.Expr {
	if i < len(vals) {
		return vals[i]
	}
	return nil
}

// isMainIdent checks whether an identifier token is the literal string "main".
func isMainIdent(name *ast.Ident) bool {
	return name.TokType == token.IDENT && name.Name == "main"
}

var _ = isMainIdent // suppress unused warning

// unquoteString strips the surrounding quote characters from a raw string
// literal value (e.g. `"hello"` → `hello`, `'world'` → `world`).
// Escape sequences are preserved as-is for the codegen to re-emit.
func unquoteString(raw string) string {
	if len(raw) < 2 {
		return raw
	}
	// Both " and ' delimiters are valid in Yz.
	if (raw[0] == '"' && raw[len(raw)-1] == '"') ||
		(raw[0] == '\'' && raw[len(raw)-1] == '\'') {
		inner := raw[1 : len(raw)-1]
		// Unescape basic sequences.
		inner = strings.ReplaceAll(inner, `\"`, `"`)
		inner = strings.ReplaceAll(inner, `\'`, `'`)
		inner = strings.ReplaceAll(inner, `\\`, `\`)
		inner = strings.ReplaceAll(inner, `\n`, "\n")
		inner = strings.ReplaceAll(inner, `\t`, "\t")
		inner = strings.ReplaceAll(inner, `\r`, "\r")
		return inner
	}
	return raw
}
