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
	l := &lowerer{analyzer: a}
	return l.lowerFile(sf, pkgName)
}

// ---------------------------------------------------------------------------
// Lowerer state
// ---------------------------------------------------------------------------

type lowerer struct {
	analyzer *sema.Analyzer

	// When lowering a method body, these describe the receiver.
	recvName   string
	recvFields map[string]bool // fields accessible via receiver
}

// ---------------------------------------------------------------------------
// File
// ---------------------------------------------------------------------------

func (l *lowerer) lowerFile(sf *ast.SourceFile, pkgName string) *File {
	f := &File{PkgName: pkgName}
	for _, node := range sf.Stmts {
		if d := l.lowerTopLevel(node); d != nil {
			f.Decls = append(f.Decls, d)
		}
	}
	return f
}

func (l *lowerer) lowerTopLevel(node ast.Node) Decl {
	switch n := node.(type) {
	case *ast.ShortDecl:
		return l.lowerTopShortDecl(n)
	case *ast.BocWithSig:
		return l.lowerBocWithSig(n, "", nil)
	default:
		return nil
	}
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
	sd := &SingletonDecl{TypeName: typeName, VarName: name}

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
		}
	}
	return sd
}

// collectFieldNames returns the set of field names (non-boc ShortDecls and TypedDecls
// without values) in the boc literal, for use in method body lowering.
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
		Name:     name,
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
			expr := l.lowerExpr(e)
			if isLast {
				inner = append(inner, &ReturnStmt{Value: expr})
			} else {
				inner = append(inner, &ExprStmt{Expr: expr})
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
		typ := l.goType(l.analyzer.ExprType(val))
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
	sd := &StructDecl{Name: name}
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
			for i, n := range e.Names {
				var initExpr Expr
				if i < len(e.Values) {
					initExpr = l.lowerExpr(e.Values[i])
				}
				typ := l.goType(l.analyzer.ExprType(l.valueAt(e.Values, i)))
				sd.Fields = append(sd.Fields, &FieldSpec{Name: n.Name, Type: typ, Init: initExpr})
			}
		}
	}
	return sd
}

// ---------------------------------------------------------------------------
// main boc → FuncDecl
// ---------------------------------------------------------------------------

func (l *lowerer) lowerMainBoc(b *ast.BocLiteral) *FuncDecl {
	fn := &FuncDecl{Name: "main"}

	// Separate boc-method-call ExprStmts from other statements.
	// Boc calls are run concurrently through a BocGroup; the group is waited
	// on before the remaining statements execute (structured concurrency).
	var bocCalls []ast.Expr
	var otherStmts []ast.Node
	for _, elem := range b.Elements {
		if expr, ok := elem.(ast.Expr); ok && l.isBocMethodCall(expr) {
			bocCalls = append(bocCalls, expr)
		} else {
			otherStmts = append(otherStmts, elem)
		}
	}

	if len(bocCalls) > 0 {
		const bgVar = "_bg"
		// _bg := &std.BocGroup{}
		fn.Body = append(fn.Body, &DeclStmt{Name: bgVar, Init: &NewGroupExpr{}})
		// _bg.Go(func() any { return call.Force() })
		for _, call := range bocCalls {
			callExpr := l.lowerExpr(call)
			fn.Body = append(fn.Body, &ExprStmt{
				Expr: &SpawnExpr{
					GroupVar: bgVar,
					Body:     []Stmt{&ReturnStmt{Value: &ForceExpr{Thunk: callExpr}}},
				},
			})
		}
		// _bg.Wait()
		fn.Body = append(fn.Body, &WaitStmt{GroupVar: bgVar})
	}

	for _, elem := range otherStmts {
		fn.Body = append(fn.Body, l.lowerMainStmt(elem)...)
	}
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
				typ = l.goType(l.analyzer.ExprType(e.Values[i]))
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
		return []Stmt{&ExprStmt{Expr: l.lowerExpr(e)}}
	default:
		return nil
	}
}

// ---------------------------------------------------------------------------
// BocWithSig
// ---------------------------------------------------------------------------

func (l *lowerer) lowerBocWithSig(bws *ast.BocWithSig, recvType string, parentFields map[string]bool) Decl {
	// BocWithSig at top level without a receiver is an independent singleton.
	// Full BocWithSig support (as methods with explicit signatures) will be
	// expanded in Phase 5 codegen.
	if bws.Body == nil {
		return nil
	}
	// Treat as a singleton boc for now.
	return l.lowerSingletonBoc(bws.Name.Name, bws.Body)
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
	case *ast.Ident:
		if expr.Name == "true" {
			return &BoolLit{Val: true}
		}
		if expr.Name == "false" {
			return &BoolLit{Val: false}
		}
		return l.lowerName(expr.Name)
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
	case *ast.ArrayLiteral:
		return l.lowerArrayLit(expr)
	case *ast.DictLiteral:
		return l.lowerDictLit(expr)
	default:
		return &UnitLit{}
	}
}

// lowerName resolves an identifier: if it's a receiver field, emit self.field.
func (l *lowerer) lowerName(name string) Expr {
	if l.recvName != "" && l.recvFields[name] {
		return &FieldAccess{Object: &Ident{Name: l.recvName}, Field: name}
	}
	return &Ident{Name: name}
}

// builtinGoName maps a Yz builtin name to its Go runtime equivalent.
// Builtins are emitted as direct calls (not goroutine-wrapped).
var builtinGoName = map[string]string{
	"print": "std.Print",
	"while": "std.While",
	"info":  "std.Info",
}

func (l *lowerer) lowerCall(c *ast.CallExpr) Expr {
	// Check for known builtins first — they emit as direct std.Xxx calls.
	if id, ok := c.Callee.(*ast.Ident); ok {
		if goName, isBuiltin := builtinGoName[id.Name]; isBuiltin {
			return l.lowerBuiltinCall(goName, c)
		}
	}

	callee := l.lowerExpr(c.Callee)
	var args []Expr
	for _, arg := range c.Args {
		args = append(args, l.lowerExpr(arg.Value))
	}

	// For a field-access callee like counter.increment or counter.value():
	// the method itself already wraps its body in std.Go and returns *Thunk[T].
	// We just emit the direct MethodCall — no extra wrapping needed.
	if fa, ok := callee.(*FieldAccess); ok {
		calleeType := l.analyzer.ExprType(c.Callee)
		if _, isBocType := calleeType.(*sema.BocType); isBocType || l.isSingletonBoc(fa.Object) {
			// Boc method call: emit as MethodCall; result is *Thunk[T].
			return &MethodCall{Recv: fa.Object, Method: fa.Field, Args: args}
		}
		// Plain struct field/method access.
		return &MethodCall{Recv: fa.Object, Method: fa.Field, Args: args}
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

// isBocMethodCall reports whether an AST expression is a call on a singleton
// boc method (i.e. `counter.value()`) which returns a *Thunk in Go.
func (l *lowerer) isBocMethodCall(e ast.Expr) bool {
	c, ok := e.(*ast.CallExpr)
	if !ok {
		return false
	}
	mem, ok := c.Callee.(*ast.MemberExpr)
	if !ok {
		return false
	}
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
	// An anonymous boc literal used as an expression — emit as a closure.
	semType := l.analyzer.ExprType(b)
	resultType := "std.Unit"
	if bt, ok := semType.(*sema.BocType); ok && len(bt.Returns) > 0 {
		resultType = l.goType(bt.Returns[0])
	}
	body := l.lowerBocBody(b, resultType)
	return &ClosureExpr{ResultType: resultType, Body: body}
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
	// Emit std.NewDict[K,V]().Set(k,v).Set(k,v)... chain.
	// For simplicity, represent as a FuncCall "std.NewDict" with alternating k,v args.
	// The codegen will expand this into Set calls.
	var args []Expr
	for _, entry := range d.Entries {
		args = append(args, l.lowerExpr(entry.Key), l.lowerExpr(entry.Value))
	}
	return &FuncCall{
		Func: &Ident{Name: "std.NewDictLit"},
		Args: args,
	}
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

// goTypeFromTypeExpr converts an ast.TypeExpr to a Go type string.
func (l *lowerer) goTypeFromTypeExpr(te ast.TypeExpr) string {
	if te == nil {
		return "any"
	}
	switch t := te.(type) {
	case *ast.SimpleTypeExpr:
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
