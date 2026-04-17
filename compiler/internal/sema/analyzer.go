package sema

import (
	"fmt"
	"strings"

	"yz/internal/ast"
	"yz/internal/token"
)

// ---------------------------------------------------------------------------
// SemaError
// ---------------------------------------------------------------------------

// SemaError is a single semantic error with source location.
type SemaError struct {
	Msg  string
	Line int
	Col  int
	Len  int // byte length to underline (0 → 1 caret)
}

func (e *SemaError) Error() string {
	return fmt.Sprintf("sema error at L%d:C%d: %s", e.Line, e.Col, e.Msg)
}

// SemaErrors collects multiple errors and implements the error interface.
type SemaErrors []*SemaError

func (es SemaErrors) Error() string {
	msgs := make([]string, len(es))
	for i, e := range es {
		msgs[i] = e.Error()
	}
	return strings.Join(msgs, "\n")
}

// ---------------------------------------------------------------------------
// Analyzer
// ---------------------------------------------------------------------------

// Analyzer performs semantic analysis over one SourceFile.
type Analyzer struct {
	// types maps every ast.Node to its inferred/checked Type.
	types map[ast.Node]Type

	// fileScope is the top-level scope for the analyzed file.
	fileScope *Scope

	// currentScope is the active scope during traversal.
	currentScope *Scope

	// fqnPrefix is the dot-joined path to the current lexical position.
	fqnPrefix string

	// errors collected during analysis.
	errors SemaErrors

	// lastExpr is the most recently analyzed top-level node (for tests).
	lastExpr ast.Node

	// activeConstraints collects inferred method requirements on generic type
	// params while analyzing a generic struct body. Non-nil only during that
	// analysis; maps type-param name (e.g. "T") → list of constraints.
	activeConstraints map[string][]*GenericConstraint

	// activeContext is "StructName.methodName" for constraint attribution.
	// Updated before each BocWithSig method body is analyzed.
	activeContext string
}

// NewAnalyzer creates a fresh Analyzer with built-in symbols pre-loaded.
func NewAnalyzer() *Analyzer {
	builtin := newBuiltinScope()
	file := newScope(builtin)
	return &Analyzer{
		types:        make(map[ast.Node]Type),
		fileScope:    file,
		currentScope: file,
	}
}

// AnalyzeFile performs semantic analysis on the given SourceFile.
func (a *Analyzer) AnalyzeFile(sf *ast.SourceFile) error {
	for _, node := range sf.Stmts {
		a.analyzeNode(node)
		a.lastExpr = node
	}
	if len(a.errors) > 0 {
		return a.errors
	}
	return nil
}

// ---------------------------------------------------------------------------
// Public result accessors
// ---------------------------------------------------------------------------

func (a *Analyzer) ExprType(n ast.Node) Type {
	if n == nil {
		return Unknown
	}
	if t, ok := a.types[n]; ok {
		return t
	}
	return Unknown
}

func (a *Analyzer) LookupInFile(name string) *Symbol {
	return a.fileScope.Lookup(name)
}

func (a *Analyzer) LastExpr() ast.Node { return a.lastExpr }

// ExportedSymbols returns a snapshot of all symbols defined at file scope.
func (a *Analyzer) ExportedSymbols() map[string]*Symbol {
	result := make(map[string]*Symbol, len(a.fileScope.syms))
	for name, sym := range a.fileScope.syms {
		result[name] = sym
	}
	return result
}

// RegisterPackage registers the exports of a compiled sub-package under the
// FQN namespace tree. relDir is like "house/front", pkgAlias is "front",
// importPath is "yzapp/house/front".
func (a *Analyzer) RegisterPackage(relDir, pkgAlias, importPath string, exports map[string]*Symbol) {
	pkgType := &PackageType{PkgAlias: pkgAlias, ImportPath: importPath, Exports: exports}
	pkgSym := &Symbol{Name: pkgAlias, Type: pkgType}

	parts := strings.Split(relDir, "/")
	if len(parts) == 1 {
		a.fileScope.Define(&Symbol{Name: parts[0], Type: pkgType})
		return
	}
	// Build or extend the namespace tree rooted at parts[0].
	rootName := parts[0]
	existing := a.fileScope.LookupLocal(rootName)
	var ns *NamespaceType
	if existing != nil {
		if existingNs, ok := existing.Type.(*NamespaceType); ok {
			ns = existingNs
		} else {
			ns = &NamespaceType{Children: make(map[string]*Symbol)}
		}
	} else {
		ns = &NamespaceType{Children: make(map[string]*Symbol)}
	}
	// Recursively insert remaining parts into the namespace tree.
	insertNamespace(ns, parts[1:], pkgSym)
	a.fileScope.Define(&Symbol{Name: rootName, Type: ns})
}

// insertNamespace inserts pkgSym at the leaf of a namespace path.
func insertNamespace(ns *NamespaceType, parts []string, pkgSym *Symbol) {
	if len(parts) == 1 {
		ns.Children[parts[0]] = pkgSym
		return
	}
	childName := parts[0]
	childSym, ok := ns.Children[childName]
	var childNs *NamespaceType
	if ok {
		if existingNs, ok2 := childSym.Type.(*NamespaceType); ok2 {
			childNs = existingNs
		} else {
			childNs = &NamespaceType{Children: make(map[string]*Symbol)}
		}
	} else {
		childNs = &NamespaceType{Children: make(map[string]*Symbol)}
	}
	insertNamespace(childNs, parts[1:], pkgSym)
	ns.Children[childName] = &Symbol{Name: childName, Type: childNs}
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

func (a *Analyzer) errorf(pos ast.Pos, format string, args ...any) {
	a.errors = append(a.errors, &SemaError{
		Msg:  fmt.Sprintf(format, args...),
		Line: pos.Line,
		Col:  pos.Col,
	})
}

func (a *Analyzer) errorfLen(pos ast.Pos, length int, format string, args ...any) {
	a.errors = append(a.errors, &SemaError{
		Msg:  fmt.Sprintf(format, args...),
		Line: pos.Line,
		Col:  pos.Col,
		Len:  length,
	})
}

func (a *Analyzer) setType(n ast.Node, t Type) { a.types[n] = t }

func (a *Analyzer) pushScope() *Scope {
	s := newScope(a.currentScope)
	prev := a.currentScope
	a.currentScope = s
	return prev
}

func (a *Analyzer) popScope(prev *Scope) { a.currentScope = prev }

func (a *Analyzer) pushFQN(name string) string {
	prev := a.fqnPrefix
	if a.fqnPrefix == "" {
		a.fqnPrefix = name
	} else {
		a.fqnPrefix = a.fqnPrefix + "." + name
	}
	return prev
}

func (a *Analyzer) popFQN(prev string) { a.fqnPrefix = prev }

func (a *Analyzer) currentFQN(name string) string {
	if a.fqnPrefix == "" {
		return name
	}
	return a.fqnPrefix + "." + name
}

// define registers a symbol in the current scope and, when at file scope,
// also in the file scope for external lookup.
func (a *Analyzer) define(sym *Symbol) {
	a.currentScope.Define(sym)
	if a.currentScope == a.fileScope {
		a.fileScope.Define(sym)
	}
}

// ---------------------------------------------------------------------------
// Node dispatch
// ---------------------------------------------------------------------------

func (a *Analyzer) analyzeNode(n ast.Node) Type {
	switch node := n.(type) {
	case *ast.ShortDecl:
		return a.analyzeShortDecl(node)
	case *ast.TypedDecl:
		return a.analyzeTypedDecl(node)
	case *ast.Assignment:
		return a.analyzeAssignment(node)
	case *ast.BocWithSig:
		return a.analyzeBocWithSig(node)
	case *ast.VariantDef:
		return a.analyzeVariantDef(node)
	case *ast.ReturnStmt:
		var t Type = TypUnit
		if node.Value != nil {
			t = a.analyzeExpr(node.Value)
		}
		a.setType(node, t)
		return t
	case *ast.BreakStmt, *ast.ContinueStmt:
		return TypUnit
	case *ast.MixStmt:
		return TypUnit // handled inside analyzeStructBoc
	case *ast.InfoString:
		return TypUnit
	case ast.Expr:
		return a.analyzeExpr(node)
	default:
		return Unknown
	}
}

// ---------------------------------------------------------------------------
// Short declaration
// ---------------------------------------------------------------------------

func (a *Analyzer) analyzeShortDecl(d *ast.ShortDecl) Type {
	// Special case: single name + BocLiteral value.
	if len(d.Names) == 1 && len(d.Values) == 1 {
		name := d.Names[0]
		if bocLit, ok := d.Values[0].(*ast.BocLiteral); ok {
			return a.analyzeBocDecl(name, bocLit, d)
		}
	}

	// General case: analyze RHS then bind each name.
	var valTypes []Type
	for _, v := range d.Values {
		valTypes = append(valTypes, a.analyzeExpr(v))
	}
	for i, name := range d.Names {
		var typ Type = Unknown
		if i < len(valTypes) {
			typ = valTypes[i]
		}
		fqn := a.currentFQN(name.Name)
		a.define(&Symbol{Name: name.Name, Type: typ, FQN: fqn, Node: d})
	}
	return TypUnit
}

// analyzeBocDecl handles `name: { ... }` for both lowercase (boc) and
// uppercase (struct type) names.
func (a *Analyzer) analyzeBocDecl(name *ast.Ident, bocLit *ast.BocLiteral, decl ast.Node) Type {
	fqn := a.currentFQN(name.Name)
	prevFQN := a.pushFQN(name.Name)
	prev := a.pushScope()

	var typ Type
	if name.TokType == token.TYPE_IDENT || name.TokType == token.GENERIC_IDENT {
		// Uppercase (multi-char TYPE_IDENT or single-letter GENERIC_IDENT): struct type definition.
		st := a.analyzeStructBoc(name.Name, bocLit)
		typ = st
	} else {
		// Lowercase: boc definition.
		bt := a.analyzeBocBody(bocLit.Elements)
		params := a.collectParams(bocLit.Elements)
		returns := bt
		if len(returns) == 0 {
			returns = []Type{TypUnit}
		}
		typ = &BocType{Params: params, Returns: returns}
	}

	a.popScope(prev)
	a.popFQN(prevFQN)

	sym := &Symbol{Name: name.Name, Type: typ, FQN: fqn, Node: decl}
	a.define(sym)

	// Register variant constructors in the outer scope so they can be called.
	if st, ok := typ.(*StructType); ok && st.IsVariant {
		for i := range st.Variants {
			vc := &st.Variants[i]
			params := make([]BocParam, len(vc.Fields))
			for j, f := range vc.Fields {
				params[j] = BocParam{Label: f.Name, Type: f.Type}
			}
			a.define(&Symbol{
				Name:           vc.Name,
				Type:           &BocType{Params: params, Returns: []Type{st}},
				Node:           decl,
				ParentTypeName: st.Name,
			})
		}
	}

	a.setType(decl, typ)
	return typ
}

// collectParams scans boc elements for uninitialized TypedDecls — these are
// the boc's input parameters.
func (a *Analyzer) collectParams(elements []ast.Node) []BocParam {
	var params []BocParam
	for _, elem := range elements {
		if td, ok := elem.(*ast.TypedDecl); ok && td.Value == nil {
			typ := a.resolveTypeExpr(td.Type)
			params = append(params, BocParam{Label: td.Name.Name, Type: typ})
		}
	}
	return params
}

// ---------------------------------------------------------------------------
// Typed declaration
// ---------------------------------------------------------------------------

func (a *Analyzer) analyzeTypedDecl(d *ast.TypedDecl) Type {
	typ := a.resolveTypeExpr(d.Type)
	if d.Value != nil {
		valTyp := a.analyzeExpr(d.Value)
		if valTyp != Unknown && !valTyp.IsCompatibleWith(typ) {
			a.errorf(d.Pos, "type mismatch: %v is not compatible with %v", valTyp, typ)
		}
	}
	fqn := a.currentFQN(d.Name.Name)
	a.define(&Symbol{Name: d.Name.Name, Type: typ, FQN: fqn, Node: d})
	a.setType(d, typ)
	return typ
}

// ---------------------------------------------------------------------------
// Assignment
// ---------------------------------------------------------------------------

func (a *Analyzer) analyzeAssignment(asgn *ast.Assignment) Type {
	var valTypes []Type
	for _, v := range asgn.Values {
		valTypes = append(valTypes, a.analyzeExpr(v))
	}
	if asgn.Target != nil {
		targetType := a.resolveTargetType(asgn.Target)
		if len(valTypes) > 0 && targetType != Unknown {
			if !valTypes[0].IsCompatibleWith(targetType) {
				a.errorf(asgn.Pos, "assignment: %v is not compatible with %v", valTypes[0], targetType)
			}
		}
	} else {
		for i, name := range asgn.Names {
			sym := a.currentScope.Lookup(name.Name)
			if sym == nil {
				a.errorf(name.Pos, "undefined: %s", name.Name)
				continue
			}
			if i < len(valTypes) && sym.Type != Unknown {
				if !valTypes[i].IsCompatibleWith(sym.Type) {
					a.errorf(name.Pos, "assignment to %s: %v not compatible with %v",
						name.Name, valTypes[i], sym.Type)
				}
			}
		}
	}
	return TypUnit
}

func (a *Analyzer) resolveTargetType(target ast.Expr) Type {
	switch t := target.(type) {
	case *ast.Ident:
		sym := a.currentScope.Lookup(t.Name)
		if sym == nil {
			a.errorfLen(t.Pos, len(t.Name), "undefined: %s", t.Name)
			return Unknown
		}
		return sym.Type
	case *ast.MemberExpr:
		objType := a.analyzeExpr(t.Object)
		return a.fieldType(objType, t.Member.Name, t.Pos)
	case *ast.IndexExpr:
		objType := a.analyzeExpr(t.Object)
		if at, ok := objType.(*ArrayType); ok {
			return at.Elem
		}
		if dt, ok := objType.(*DictType); ok {
			return dt.Val
		}
		return Unknown
	}
	return Unknown
}

// ---------------------------------------------------------------------------
// BocWithSig
// ---------------------------------------------------------------------------

func (a *Analyzer) analyzeBocWithSig(bws *ast.BocWithSig) Type {
	// Resolve all params from the signature.
	// In body-only form (= { body }), all params are inputs (no isReturn logic).
	allParams := a.resolveBocSigParams(bws.Sig, bws.BodyOnly)

	// Separate input params from anonymous return-type entries.
	var inputParams []BocParam
	var explicitReturns []Type
	for _, p := range allParams {
		if p.IsReturn {
			explicitReturns = append(explicitReturns, p.Type)
		} else {
			inputParams = append(inputParams, p)
		}
	}

	// Uppercase name + no body → structural type declaration (interface-style).
	// `Name #(name String, age Int)` registers Name as a StructType with those
	// fields. No implementation is generated; any structurally compatible boc
	// literal satisfies the type.
	if bws.Body == nil && (bws.Name.TokType == token.TYPE_IDENT || bws.Name.TokType == token.GENERIC_IDENT) {
		// If every input param is a BocType, this is an interface declaration.
		allBoc := len(inputParams) > 0
		for _, p := range inputParams {
			if _, isBoc := p.Type.(*BocType); !isBoc {
				allBoc = false
				break
			}
		}
		st := &StructType{Name: bws.Name.Name, IsInterface: allBoc}
		for _, p := range inputParams {
			if p.Label != "" {
				st.Fields = append(st.Fields, StructField{Name: p.Label, Type: p.Type})
			}
		}
		fqn := a.currentFQN(bws.Name.Name)
		sym := &Symbol{Name: bws.Name.Name, Type: st, FQN: fqn, Node: bws}
		a.define(sym)
		a.setType(bws, st)
		return st
	}

	var returns []Type
	if bws.Body != nil {
		prev := a.pushScope()
		prevFQN := a.pushFQN(bws.Name.Name)
		if bws.BodyOnly {
			// Body-only form: extract named params from body's initial TypedDecls.
			// The body redeclares its own params; sig provides types (and names
			// when labeled) for validation.
			bodyParams, n := a.extractBodyParams(bws.Body.Elements, inputParams)
			if bodyParams != nil {
				inputParams = bodyParams
			} else if n == 0 && len(inputParams) > 0 {
				// Body has no TypedDecl params. Decide based on whether sig params are labeled.
				allAnonymous := true
				for _, p := range inputParams {
					if p.Label != "" {
						allAnonymous = false
						break
					}
				}
				if allAnonymous {
					// All sig params are anonymous types — re-interpret as shorthand form:
					// trailing unlabeled type = return type, no input params.
					// This makes `foo #(String) = { "hello" }` identical to `foo #(String) { "hello" }`.
					shorthandParams := a.resolveBocSigParams(bws.Sig, false)
					inputParams = nil
					explicitReturns = nil
					for _, p := range shorthandParams {
						if p.IsReturn {
							explicitReturns = append(explicitReturns, p.Type)
						} else {
							inputParams = append(inputParams, p)
						}
					}
				} else {
					// Labeled params must be redeclared in body — report the original error.
					if len(bws.Body.Elements) > 0 {
						a.errorf(bws.Body.Elements[0].Position(),
							"expected parameter declaration (name Type), got %T", bws.Body.Elements[0])
					}
				}
			}
			for _, p := range inputParams {
				if p.Label != "" {
					a.currentScope.Define(&Symbol{Name: p.Label, Type: p.Type})
				}
			}
			// Analyze body starting after the param declarations.
			bodyElems := bws.Body.Elements
			if n > 0 && n <= len(bodyElems) {
				bodyElems = bodyElems[n:]
			}
			returns = a.analyzeBocBody(bodyElems)
		} else {
			// Shorthand form: inject params from sig directly into body scope.
			for _, p := range inputParams {
				if p.Label != "" {
					a.currentScope.Define(&Symbol{Name: p.Label, Type: p.Type})
				}
			}
			returns = a.analyzeBocBody(bws.Body.Elements)
		}
		a.popFQN(prevFQN)
		a.popScope(prev)
	}

	// Explicit return types in the signature override inferred returns.
	if len(explicitReturns) > 0 {
		returns = explicitReturns
	}
	if len(returns) == 0 {
		returns = []Type{TypUnit}
	}

	bocType := &BocType{Params: inputParams, Returns: returns}
	fqn := a.currentFQN(bws.Name.Name)
	sym := &Symbol{Name: bws.Name.Name, Type: bocType, FQN: fqn, Node: bws}
	a.define(sym)
	a.setType(bws, bocType)
	return bocType
}

// resolveBocSigParams resolves the params of a BocTypeExpr signature.
// When bodyOnly is true (the `= { body }` form), all params are treated as
// inputs: unlabeled types are anonymous inputs, not return-type annotations.
func (a *Analyzer) resolveBocSigParams(sig *ast.BocTypeExpr, bodyOnly bool) []BocParam {
	var params []BocParam
	for _, p := range sig.Params {
		if p.Variant != nil {
			continue
		}
		var typ Type
		if p.Type != nil {
			typ = a.resolveTypeExpr(p.Type)
		} else if p.Default != nil {
			// ShortDecl-style param (name : expr): infer type from default value.
			typ = a.analyzeExpr(p.Default)
		}
		// In body-only form, all params are inputs; otherwise unlabeled = return.
		isReturn := !bodyOnly && p.Label == "" && typ != nil
		params = append(params, BocParam{
			Label:      p.Label,
			Type:       typ,
			HasDefault: p.Default != nil,
			IsReturn:   isReturn,
		})
	}
	return params
}

// extractBodyParams reads the first len(sigParams) body elements as parameter
// declarations (TypedDecl). For each:
//   - If the sig param has a label, the body TypedDecl name must match.
//   - If the sig param is unlabeled, any name is accepted.
//   - Types must be compatible.
//
// Returns the matched params (names taken from the body) and the count
// consumed, or (nil, 0) if an error was reported.
func (a *Analyzer) extractBodyParams(elements []ast.Node, sigParams []BocParam) ([]BocParam, int) {
	n := len(sigParams)
	if n == 0 {
		return nil, 0
	}
	if n > len(elements) {
		pos := ast.Pos{Line: 1, Col: 1}
		if len(elements) > 0 {
			pos = elements[0].Position()
		}
		a.errorf(pos, "body has only %d statement(s) but signature expects %d param(s)", len(elements), n)
		return nil, 0
	}
	result := make([]BocParam, n)
	for i, sp := range sigParams {
		elem := elements[i]
		td, ok := elem.(*ast.TypedDecl)
		if !ok {
			// Body element is not a TypedDecl — no params found; caller handles fallback.
			return nil, 0
		}
		bodyType := a.resolveTypeExpr(td.Type)
		// Name check: if sig param is named, body must use the same name.
		if sp.Label != "" && td.Name.Name != sp.Label {
			a.errorfLen(td.Name.Pos, len(td.Name.Name),
				"param name mismatch: body declares %q but signature requires %q",
				td.Name.Name, sp.Label)
			return nil, 0
		}
		// Type check.
		if !bodyType.IsCompatibleWith(sp.Type) {
			a.errorf(td.Position(),
				"param type mismatch: body has %v, signature expects %v", bodyType, sp.Type)
			return nil, 0
		}
		result[i] = BocParam{
			Label:      td.Name.Name,
			Type:       bodyType,
			HasDefault: td.Value != nil,
		}
	}
	return result, n
}

// ---------------------------------------------------------------------------
// Boc body analysis
// ---------------------------------------------------------------------------

// analyzeBocBody analyzes elements of a boc body and returns the types of
// the trailing expression(s) — these are the boc's return values.
func (a *Analyzer) analyzeBocBody(elements []ast.Node) []Type {
	var lastExprTypes []Type
	for _, elem := range elements {
		t := a.analyzeNode(elem)
		switch elem.(type) {
		case ast.Expr:
			lastExprTypes = []Type{t}
		case *ast.ReturnStmt:
			lastExprTypes = []Type{t}
		default:
			// Statements don't contribute to return type.
			lastExprTypes = nil
		}
	}
	return lastExprTypes
}

// ---------------------------------------------------------------------------
// Struct type analysis (uppercase boc)
// ---------------------------------------------------------------------------

func (a *Analyzer) analyzeStructBoc(name string, b *ast.BocLiteral) *StructType {
	st := &StructType{Name: name}
	fieldSet := make(map[string]bool)

	// Pre-scan: detect whether this is a generic struct so we can activate
	// constraint collection before processing method bodies.
	isGeneric := false
	for _, elem := range b.Elements {
		if id, ok := elem.(*ast.Ident); ok && id.TokType == token.GENERIC_IDENT {
			isGeneric = true
			break
		}
	}

	// Save outer constraint state and activate fresh collection for generic structs.
	prevConstraints := a.activeConstraints
	prevContext := a.activeContext
	if isGeneric {
		a.activeConstraints = make(map[string][]*GenericConstraint)
	}

	for _, elem := range b.Elements {
		switch e := elem.(type) {
		case *ast.TypedDecl:
			typ := a.analyzeTypedDecl(e)
			if fieldSet[e.Name.Name] {
				a.errorf(e.Pos, "duplicate field %q in %s", e.Name.Name, name)
				continue
			}
			fieldSet[e.Name.Name] = true
			st.Fields = append(st.Fields, StructField{Name: e.Name.Name, Type: typ})

		case *ast.ShortDecl:
			for i, n := range e.Names {
				var typ Type
				if i < len(e.Values) {
					typ = a.analyzeExpr(e.Values[i])
				}
				if fieldSet[n.Name] {
					a.errorf(e.Pos, "duplicate field %q in %s", n.Name, name)
					continue
				}
				fieldSet[n.Name] = true
				st.Fields = append(st.Fields, StructField{Name: n.Name, Type: typ})
				a.currentScope.Define(&Symbol{Name: n.Name, Type: typ, Node: e})
			}

		case *ast.MixStmt:
			a.applyMix(e, st, fieldSet, name)

		case *ast.Ident:
			// Generic type param declaration (T, E inside type boc body).
			// Register as GenericType in current scope and record on the struct.
			if e.TokType == token.GENERIC_IDENT {
				st.TypeParams = append(st.TypeParams, e.Name)
			}
			gt := &GenericType{Name: e.Name}
			a.currentScope.Define(&Symbol{Name: e.Name, Type: gt, Node: e})

		case *ast.BocWithSig:
			// Set the constraint attribution context before analyzing the method body
			// so that any T-method calls recorded during analysis are tagged correctly.
			if a.activeConstraints != nil {
				a.activeContext = name + "." + e.Name.Name
			}
			typ := a.analyzeBocWithSig(e)
			if !fieldSet[e.Name.Name] {
				fieldSet[e.Name.Name] = true
				st.Fields = append(st.Fields, StructField{Name: e.Name.Name, Type: typ})
			}

		case *ast.VariantDef:
			vc := a.collectVariantCase(e)
			st.IsVariant = true
			st.Variants = append(st.Variants, vc)
			// Merge variant fields into the parent flat struct (deduplicated).
			for _, f := range vc.Fields {
				if !fieldSet[f.Name] {
					fieldSet[f.Name] = true
					st.Fields = append(st.Fields, f)
					a.currentScope.Define(&Symbol{Name: f.Name, Type: f.Type})
				}
			}

		default:
			a.analyzeNode(elem)
		}
	}

	// Freeze the inferred constraints into the struct type and restore outer state.
	if isGeneric {
		if len(a.activeConstraints) > 0 {
			st.TypeConstraints = a.activeConstraints
		}
		a.activeConstraints = prevConstraints
		a.activeContext = prevContext
	}

	return st
}

// collectVariantCase resolves a VariantDef into a VariantCase for the parent struct.
func (a *Analyzer) collectVariantCase(v *ast.VariantDef) VariantCase {
	vc := VariantCase{Name: v.Name}
	for _, p := range v.Params {
		if p.Label != "" && p.Type != nil {
			vc.Fields = append(vc.Fields, StructField{
				Name: p.Label,
				Type: a.resolveTypeExpr(p.Type),
			})
		}
	}
	return vc
}

func (a *Analyzer) applyMix(ms *ast.MixStmt, st *StructType, fieldSet map[string]bool, hostName string) {
	mixedSym := a.currentScope.Lookup(ms.Name.Name)
	if mixedSym == nil {
		a.errorfLen(ms.Name.Pos, len(ms.Name.Name), "undefined: %s", ms.Name.Name)
		return
	}
	mixedSt, ok := mixedSym.Type.(*StructType)
	if !ok {
		a.errorf(ms.Name.Pos, "mix: %s is not a struct type", ms.Name.Name)
		return
	}
	for _, f := range mixedSt.Fields {
		if fieldSet[f.Name] {
			a.errorf(ms.Pos, "mix conflict: field %q already defined (mix of %s into %s)",
				f.Name, ms.Name.Name, hostName)
			continue
		}
		fieldSet[f.Name] = true
		st.Fields = append(st.Fields, f)
		a.currentScope.Define(&Symbol{Name: f.Name, Type: f.Type})
	}
}

func (a *Analyzer) analyzeVariantDef(v *ast.VariantDef) Type {
	vc := a.collectVariantCase(v)
	variantType := &StructType{Name: v.Name, Fields: vc.Fields}
	a.currentScope.Define(&Symbol{Name: v.Name, Type: variantType, Node: v})
	return variantType
}

// ---------------------------------------------------------------------------
// Expression analysis
// ---------------------------------------------------------------------------

func (a *Analyzer) analyzeExpr(e ast.Expr) Type {
	if e == nil {
		return Unknown
	}
	var t Type
	switch expr := e.(type) {
	case *ast.IntLit:
		t = TypInt
	case *ast.DecimalLit:
		t = TypDecimal
	case *ast.StringLit:
		t = TypString
	case *ast.InterpolatedStringExpr:
		for _, part := range expr.Parts {
			if part.IsExpr {
				a.analyzeExpr(part.Expr)
			}
		}
		t = TypString
	case *ast.ConditionalExpr:
		a.analyzeExpr(expr.Cond)
		trueType := a.analyzeExpr(expr.TrueCase)
		a.analyzeExpr(expr.FalseCase)
		t = trueType
	case *ast.Ident:
		t = a.analyzeIdent(expr)
	case *ast.UnaryExpr:
		t = a.analyzeUnary(expr)
	case *ast.BinaryExpr:
		t = a.analyzeBinary(expr)
	case *ast.CallExpr:
		t = a.analyzeCall(expr)
	case *ast.MemberExpr:
		t = a.analyzeMember(expr)
	case *ast.IndexExpr:
		t = a.analyzeIndex(expr)
	case *ast.GroupExpr:
		t = a.analyzeExpr(expr.Expr)
	case *ast.BocLiteral:
		prev := a.pushScope()
		bodyReturns := a.analyzeBocBody(expr.Elements)
		params := a.collectParams(expr.Elements)
		a.popScope(prev)
		if len(bodyReturns) == 0 {
			bodyReturns = []Type{TypUnit}
		}
		t = &BocType{Params: params, Returns: bodyReturns}
	case *ast.ArrayLiteral:
		t = a.analyzeArrayLiteral(expr)
	case *ast.DictLiteral:
		t = a.analyzeDictLiteral(expr)
	case *ast.MatchExpr:
		t = a.analyzeMatch(expr)
	case *ast.InfoString:
		t = TypUnit
	default:
		t = Unknown
	}
	a.setType(e, t)
	return t
}

func (a *Analyzer) analyzeIdent(id *ast.Ident) Type {
	sym := a.currentScope.Lookup(id.Name)
	if sym == nil {
		a.errorfLen(id.Pos, len(id.Name), "undefined: %s", id.Name)
		return Unknown
	}
	return sym.Type
}

func (a *Analyzer) analyzeUnary(u *ast.UnaryExpr) Type {
	operandType := a.analyzeExpr(u.Operand)
	switch operandType {
	case TypInt:
		return TypInt
	case TypDecimal:
		return TypDecimal
	case Unknown:
		return Unknown
	default:
		a.errorf(u.Pos, "unary '-' not defined for type %v", operandType)
		return Unknown
	}
}

func (a *Analyzer) analyzeBinary(b *ast.BinaryExpr) Type {
	leftType := a.analyzeExpr(b.Left)
	a.analyzeExpr(b.Right)
	methodName := NonWordMethodName(b.Op)
	// When inside a generic struct body, a binary operator on a T-typed value
	// reveals that T must support that operator's method (e.g. == → eqeq).
	if gt, ok := leftType.(*GenericType); ok && a.activeConstraints != nil {
		a.activeConstraints[gt.Name] = append(a.activeConstraints[gt.Name], &GenericConstraint{
			TypeParam:  gt.Name,
			MethodName: methodName,
			Line:       b.Pos.Line,
			Col:        b.Pos.Col,
			Context:    a.activeContext,
		})
	}
	return a.methodReturnType(leftType, methodName, b.Pos)
}

func (a *Analyzer) analyzeCall(c *ast.CallExpr) Type {
	// When inside a generic struct body, intercept method calls on T-typed values.
	// These calls reveal what methods the type parameter T must support (constraints).
	if a.activeConstraints != nil {
		if memExpr, ok := c.Callee.(*ast.MemberExpr); ok {
			objType := a.analyzeExpr(memExpr.Object)
			if gt, ok := objType.(*GenericType); ok {
				a.activeConstraints[gt.Name] = append(a.activeConstraints[gt.Name], &GenericConstraint{
					TypeParam:  gt.Name,
					MethodName: memExpr.Member.Name,
					Line:       memExpr.Member.Pos.Line,
					Col:        memExpr.Member.Pos.Col,
					Context:    a.activeContext,
				})
				// Analyze args for their side-effects (scope bindings, nested analysis).
				for _, arg := range c.Args {
					a.analyzeExpr(arg.Value)
				}
				// Tag the callee and call as producing Unknown so downstream analysis continues.
				a.setType(c.Callee, &BocType{Returns: []Type{Unknown}})
				return Unknown
			}
		}
	}

	calleeType := a.analyzeExpr(c.Callee)
	// Collect arg types — needed for generic constraint checking at instantiation.
	var argTypes []Type
	for _, arg := range c.Args {
		argTypes = append(argTypes, a.analyzeExpr(arg.Value))
	}
	switch bt := calleeType.(type) {
	case *BocType:
		switch len(bt.Returns) {
		case 0:
			return TypUnit
		case 1:
			return bt.Returns[0]
		default:
			return Unknown // multi-return
		}
	case *StructType:
		// At a generic struct constructor call, verify that the concrete type(s)
		// bound to each type param satisfy all inferred constraints.
		if len(bt.TypeParams) > 0 && len(bt.TypeConstraints) > 0 {
			a.checkGenericConstraints(c.Callee.Position(), bt, argTypes)
		}
		return bt // constructor call
	case *BuiltinType:
		return bt // direct type value used as function
	}
	return Unknown
}

func (a *Analyzer) analyzeMember(m *ast.MemberExpr) Type {
	objType := a.analyzeExpr(m.Object)
	return a.fieldType(objType, m.Member.Name, m.Pos)
}

func (a *Analyzer) fieldType(objType Type, fieldName string, pos ast.Pos) Type {
	switch ot := objType.(type) {
	case *StructType:
		for _, f := range ot.Fields {
			if f.Name == fieldName {
				return f.Type
			}
		}
		a.errorf(pos, "type %v has no field %q", objType, fieldName)
		return Unknown
	case *BuiltinType:
		if methods, ok := builtinMethods[ot.name]; ok {
			if ret, ok := methods[fieldName]; ok {
				return ret
			}
		}
		a.errorf(pos, "type %v has no method %q", objType, fieldName)
		return Unknown
	case *BocType:
		// Accessing a method defined inside a boc — look up in scope.
		sym := a.currentScope.Lookup(fieldName)
		if sym != nil {
			return sym.Type
		}
		return Unknown
	case *NamespaceType:
		if child, ok := ot.Children[fieldName]; ok {
			return child.Type
		}
		return Unknown
	case *PackageType:
		if exp, ok := ot.Exports[fieldName]; ok {
			return exp.Type
		}
		return Unknown
	case *ArrayType:
		// HOF methods and indexed access on arrays.
		switch fieldName {
		case "filter":
			return &BocType{Returns: []Type{ot}}
		case "each":
			return &BocType{Returns: []Type{TypUnit}}
		case "map":
			return &BocType{Returns: []Type{&ArrayType{Elem: Unknown}}}
		case "any", "all":
			return &BocType{Returns: []Type{TypBool}}
		case "length":
			return &BocType{Returns: []Type{TypInt}}
		case "is_empty":
			return &BocType{Returns: []Type{TypBool}}
		case "at":
			return &BocType{Returns: []Type{ot.Elem}}
		case "append":
			return &BocType{Returns: []Type{ot}}
		}
		return Unknown // extensible — no error for unknown array methods
	case *UnknownType:
		return Unknown
	default:
		return Unknown
	}
}

func (a *Analyzer) analyzeIndex(idx *ast.IndexExpr) Type {
	objType := a.analyzeExpr(idx.Object)
	a.analyzeExpr(idx.Index)
	switch ot := objType.(type) {
	case *ArrayType:
		return ot.Elem
	case *DictType:
		return ot.Val
	}
	return Unknown
}

func (a *Analyzer) analyzeArrayLiteral(arr *ast.ArrayLiteral) Type {
	if arr.ElemType != nil {
		return &ArrayType{Elem: a.resolveTypeExpr(arr.ElemType)}
	}
	var elemType Type = Unknown
	for _, el := range arr.Elements {
		t := a.analyzeExpr(el)
		if elemType == Unknown {
			elemType = t
		}
	}
	return &ArrayType{Elem: elemType}
}

func (a *Analyzer) analyzeDictLiteral(d *ast.DictLiteral) Type {
	if d.KeyType != nil {
		return &DictType{
			Key: a.resolveTypeExpr(d.KeyType),
			Val: a.resolveTypeExpr(d.ValType),
		}
	}
	var keyType, valType Type = Unknown, Unknown
	for _, entry := range d.Entries {
		k := a.analyzeExpr(entry.Key)
		v := a.analyzeExpr(entry.Value)
		if keyType == Unknown {
			keyType = k
		}
		if valType == Unknown {
			valType = v
		}
	}
	return &DictType{Key: keyType, Val: valType}
}

func (a *Analyzer) analyzeMatch(m *ast.MatchExpr) Type {
	if m.Subject != nil {
		a.analyzeExpr(m.Subject)
	}
	var returnType Type = Unknown
	for _, arm := range m.Arms {
		if arm.Condition != nil {
			a.analyzeExpr(arm.Condition)
		}
		var armType Type = TypUnit
		prev := a.pushScope()
		for _, elem := range arm.Body {
			t := a.analyzeNode(elem)
			if _, ok := elem.(ast.Expr); ok {
				armType = t
			}
		}
		a.popScope(prev)
		if returnType == Unknown {
			returnType = armType
		}
	}
	return returnType
}

// ---------------------------------------------------------------------------
// Type resolution
// ---------------------------------------------------------------------------

func (a *Analyzer) resolveTypeExpr(te ast.TypeExpr) Type {
	if te == nil {
		return Unknown
	}
	switch t := te.(type) {
	case *ast.SimpleTypeExpr:
		// Single-letter uppercase: always a generic type parameter, never a scope lookup.
		if t.TokType == token.GENERIC_IDENT {
			return &GenericType{Name: t.Name}
		}
		sym := a.currentScope.Lookup(t.Name)
		if sym != nil {
			return sym.Type
		}
		a.errorfLen(t.Pos, len(t.Name), "undefined type: %s", t.Name)
		return Unknown
	case *ast.ArrayTypeExpr:
		return &ArrayType{Elem: a.resolveTypeExpr(t.ElemType)}
	case *ast.DictTypeExpr:
		return &DictType{
			Key: a.resolveTypeExpr(t.KeyType),
			Val: a.resolveTypeExpr(t.ValType),
		}
	case *ast.BocTypeExpr:
		params := a.resolveBocSigParams(t, false)
		var inputParams []BocParam
		var returns []Type
		for _, p := range params {
			if p.IsReturn {
				returns = append(returns, p.Type)
			} else {
				inputParams = append(inputParams, p)
			}
		}
		return &BocType{Params: inputParams, Returns: returns}
	}
	return Unknown
}

func (a *Analyzer) methodReturnType(receiverType Type, methodName string, pos ast.Pos) Type {
	switch rt := receiverType.(type) {
	case *BuiltinType:
		if methods, ok := builtinMethods[rt.name]; ok {
			if ret, ok := methods[methodName]; ok {
				return ret
			}
		}
		return Unknown
	case *StructType:
		for _, f := range rt.Fields {
			if f.Name == methodName {
				if bt, ok := f.Type.(*BocType); ok && len(bt.Returns) == 1 {
					return bt.Returns[0]
				}
				return f.Type
			}
		}
		return Unknown
	case *UnknownType, *GenericType:
		return Unknown
	default:
		return Unknown
	}
}

// checkGenericConstraints verifies that the concrete types bound to generic type
// params at a constructor call site satisfy all inferred constraints.
// It reports ALL missing methods in a single error, not one at a time.
func (a *Analyzer) checkGenericConstraints(callPos ast.Pos, st *StructType, argTypes []Type) {
	// Build typeParam → concreteType bindings by pairing constructor args with
	// data fields (skipping BocType fields, which are not constructor parameters).
	bindings := make(map[string]Type)
	argIdx := 0
	for _, field := range st.Fields {
		if _, isBoc := field.Type.(*BocType); isBoc {
			continue // method fields are not constructor params
		}
		if argIdx >= len(argTypes) {
			break
		}
		if gt, ok := field.Type.(*GenericType); ok {
			if _, alreadyBound := bindings[gt.Name]; !alreadyBound {
				bindings[gt.Name] = argTypes[argIdx]
			}
		}
		argIdx++
	}

	// For each type param with constraints, check that the concrete type has
	// all required methods.
	var violations []string
	for _, typeParam := range st.TypeParams { // iterate in declaration order for determinism
		constraints, hasConstraints := st.TypeConstraints[typeParam]
		if !hasConstraints {
			continue
		}
		concreteType, bound := bindings[typeParam]
		if !bound {
			continue // can't check if we don't know the concrete type
		}
		// Collect missing methods, deduplicating by name.
		seen := make(map[string]bool)
		var missing []string
		for _, c := range constraints {
			if seen[c.MethodName] {
				continue
			}
			seen[c.MethodName] = true
			if !a.typeHasMethod(concreteType, c.MethodName) {
				missing = append(missing, fmt.Sprintf("  %s [used in %s]", c.MethodName, c.Context))
			}
		}
		if len(missing) > 0 {
			violations = append(violations, fmt.Sprintf(
				"%s is missing methods required by %s:\n%s",
				concreteType.typeName(), typeParam, strings.Join(missing, "\n")))
		}
	}

	if len(violations) > 0 {
		a.errorf(callPos, "type constraint violation for %s:\n%s",
			st.Name, strings.Join(violations, "\n"))
	}
}

// typeHasMethod reports whether typ exposes a method or field named methodName.
func (a *Analyzer) typeHasMethod(typ Type, methodName string) bool {
	switch t := typ.(type) {
	case *StructType:
		for _, f := range t.Fields {
			if f.Name == methodName {
				return true
			}
		}
		return false
	case *BuiltinType:
		if methods, ok := builtinMethods[t.name]; ok {
			_, has := methods[methodName]
			return has
		}
		return false
	case *GenericType, *UnknownType:
		return true // can't check; assume satisfied to avoid cascading errors
	}
	return false
}
