package sema

import (
	"fmt"
	"strings"
	"unicode"

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
	// Updated before each BocDecl method body is analyzed.
	activeContext string

	// fieldInit is the definite-assignment state for the current boc body.
	// nil when outside a boc body. Only locally-constructed struct variables
	// (ShortDecl `b : Bar(...)`) are tracked; parameters are always considered
	// initialized.
	fieldInit *FieldInitState

	// expectedType is the type expected at the current expression position.
	// Set by analyzeTypedDecl, analyzeAssignment, and analyzeCall before
	// recursing into a RHS/argument expression; used by analyzeCall to
	// disambiguate ambiguous variant constructor calls.
	expectedType Type

	// currentSigParams maps parameter label → resolved Type for the sig params
	// processed so far in resolveBocSigParams. Allows later params to reference
	// earlier ones in path-dependent type expressions like `n g.Node` where `g`
	// is a preceding parameter. nil when outside a sig resolution.
	currentSigParams map[string]Type

	// inAnnotation is true while analyzing an annotation body. Used to reject
	// all string interpolation forms (${ } and backtick) inside annotations.
	inAnnotation bool
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
	// First pass: pre-register all top-level type names as empty *StructType
	// stubs so that forward and mutually-recursive references resolve to a
	// stable pointer. analyzeStructBoc reuses these pointers and fills them in,
	// so any field that captured the pointer during this pass sees the completed
	// type once the second pass finishes.
	for _, node := range sf.Stmts {
		if name, ok := topLevelTypeName(node); ok {
			a.fileScope.Define(&Symbol{Name: name, Type: &StructType{Name: name}, FQN: name})
		}
	}

	// Second pass: full analysis.
	for _, node := range sf.Stmts {
		a.analyzeNode(node)
		a.lastExpr = node
	}
	if len(a.errors) > 0 {
		return a.errors
	}
	return nil
}

// topLevelTypeName returns the declared name if node is a top-level uppercase
// struct declaration (ShortDecl with a BocLiteral RHS, or a BocDecl).
func topLevelTypeName(node ast.Node) (string, bool) {
	switch n := node.(type) {
	case *ast.ShortDecl:
		if len(n.Names) == 1 && n.Names[0].TokType == token.TYPE_IDENT &&
			len(n.Values) == 1 {
			if _, ok := n.Values[0].(*ast.BocLiteral); ok {
				return n.Names[0].Name, true
			}
		}
	case *ast.BocDecl:
		if n.Name.TokType == token.TYPE_IDENT {
			return n.Name.Name, true
		}
	}
	return "", false
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

// FindInterfaceWithMethod returns the name of an interface in the file scope
// that has a method with the given name. Returns "" if not found.
// Used by the lowerer to synthesise explicit constraints for bare type params
// whose constraints are inferred from method-param usage (YZC-0071).
func (a *Analyzer) FindInterfaceWithMethod(methodName string) string {
	for _, sym := range a.fileScope.syms {
		st, ok := sym.Type.(*StructType)
		if !ok || !st.IsInterface {
			continue
		}
		for _, f := range st.Fields {
			if f.Name == methodName {
				return st.Name
			}
		}
	}
	return ""
}

// findInterfaceMethodReturnType returns the return type of methodName from
// the first matching interface in the file scope, or Unknown if not found.
// Used in analyzeCall to give concrete return types to calls on generic params.
func (a *Analyzer) findInterfaceMethodReturnType(methodName string) Type {
	for _, sym := range a.fileScope.syms {
		st, ok := sym.Type.(*StructType)
		if !ok || !st.IsInterface {
			continue
		}
		for _, f := range st.Fields {
			if f.Name == methodName {
				if bt, ok := f.Type.(*BocType); ok && len(bt.Returns) > 0 {
					return bt.Returns[0]
				}
				return TypUnit
			}
		}
	}
	return Unknown
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
	case *ast.BocDecl:
		return a.analyzeBocDeclNode(node)
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
	case *ast.Annotation:
		a.analyzeAnnotationBody(node)
		return TypUnit
	case ast.Expr:
		return a.analyzeExpr(node)
	default:
		return Unknown
	}
}

// propagateConstructorArgInits propagates sub-field init state for struct-typed
// constructor arguments. When varName is constructed with a struct-valued arg
// (e.g. w : Wrapper(i)), this ensures that w.inner.field accesses are
// considered initialized when i.field is known to be initialized.
func (a *Analyzer) propagateConstructorArgInits(varName string, st *StructType, call *ast.CallExpr) {
	hasNamed := false
	for _, arg := range call.Args {
		if arg.Label != "" {
			hasNamed = true
			break
		}
	}
	var required []StructField
	for _, f := range st.Fields {
		if f.IsTypeField {
			continue
		}
		if !f.HasDefault {
			if _, isMethod := f.Type.(*BocType); !isMethod {
				required = append(required, f)
			}
		}
	}
	for j, arg := range call.Args {
		var field *StructField
		if hasNamed {
			for k := range required {
				if required[k].Name == arg.Label {
					field = &required[k]
					break
				}
			}
		} else if j < len(required) {
			f := required[j]
			field = &f
		}
		if field == nil {
			continue
		}
		innerSt, ok := field.Type.(*StructType)
		if !ok || innerSt.IsSingleton || innerSt.IsInterface || innerSt.IsVariant {
			continue
		}
		if rhsId, ok := arg.Value.(*ast.Ident); ok {
			if _, tracked := a.fieldInit.locals[rhsId.Name]; tracked {
				a.fieldInit.propagateInner(varName, field.Name, a.fieldInit, rhsId.Name)
			} else {
				// Untracked (parameter) — all sub-fields are initialized.
				innerFI := newFieldInitState()
				innerFI.addLocalVar("_")
				for _, f := range innerSt.Fields {
					if !f.HasDefault {
						if _, isMethod := f.Type.(*BocType); !isMethod {
							innerFI.locals["_"][f.Name] = true
						}
					}
				}
				a.fieldInit.propagateInner(varName, field.Name, innerFI, "_")
			}
		} else if innerCall, ok := arg.Value.(*ast.CallExpr); ok {
			innerFI := newFieldInitState()
			initLocalVar(innerFI, "_", innerSt, innerCall)
			a.fieldInit.propagateInner(varName, field.Name, innerFI, "_")
		}
	}
}

// ---------------------------------------------------------------------------
// Short declaration
// ---------------------------------------------------------------------------

func (a *Analyzer) analyzeShortDecl(d *ast.ShortDecl) Type {
	// Special case: single name + BocLiteral value — annotation handled in analyzeBocDecl.
	if len(d.Names) == 1 && len(d.Values) == 1 {
		name := d.Names[0]
		if bocLit, ok := d.Values[0].(*ast.BocLiteral); ok {
			return a.analyzeBocDecl(name, bocLit, d)
		}
	}
	// For non-boc short decls, analyze the annotation if present.
	if d.Annotation != nil {
		a.analyzeAnnotationBody(d.Annotation)
	}

	// Multi-name with single RHS: multi-return call expansion (YZC-0012).
	if len(d.Names) > 1 && len(d.Values) == 1 {
		valType := a.analyzeExpr(d.Values[0])
		var returnTypes []Type
		if tt, ok := valType.(*TupleType); ok {
			returnTypes = tt.Types
		} else {
			returnTypes = []Type{valType}
		}
		for i, name := range d.Names {
			var typ Type = Unknown
			if i < len(returnTypes) {
				typ = returnTypes[i]
			}
			fqn := a.currentFQN(name.Name)
			a.define(&Symbol{Name: name.Name, Type: typ, FQN: fqn, Node: d})
		}
		return TypUnit
	}

	// General case: analyze RHS then bind each name.
	var valTypes []Type
	for _, v := range d.Values {
		valTypes = append(valTypes, a.analyzeExpr(v))
	}
	for i, valTyp := range valTypes {
		if valTyp == TypUnit && i < len(d.Names) {
			a.errorf(d.Names[i].Pos, "YZC-0003: expression returns nothing, cannot assign to %s", d.Names[i].Name)
		}
	}
	for i, name := range d.Names {
		var typ Type = Unknown
		if i < len(valTypes) {
			typ = valTypes[i]
		}
		fqn := a.currentFQN(name.Name)
		a.define(&Symbol{Name: name.Name, Type: typ, FQN: fqn, Node: d})
	}
	// Track field-initialization state for locally-constructed struct variables.
	if a.fieldInit != nil {
		for i, name := range d.Names {
			if i >= len(valTypes) || i >= len(d.Values) {
				break
			}
			st, ok := valTypes[i].(*StructType)
			if !ok || st.IsSingleton || st.IsInterface || st.IsVariant {
				continue
			}
			// Aliasing (YZC-0055): c : b — clone field-init state from source var.
			// If source is untracked (parameter, always initialized), leave the
			// alias untracked too so isAssigned returns true for it.
			if id, ok := d.Values[i].(*ast.Ident); ok {
				if srcFields, tracked := a.fieldInit.locals[id.Name]; tracked {
					a.fieldInit.addLocalVar(name.Name)
					for f, assigned := range srcFields {
						if assigned {
							a.fieldInit.locals[name.Name][f] = true
						}
					}
				}
				continue
			}
			call, ok := d.Values[i].(*ast.CallExpr)
			if !ok {
				continue
			}
			// Only track field-init state for direct constructor calls (e.g. Bar()).
			// Method calls or boc calls returning a struct (e.g. r.resolve()) return
			// a fully initialized value; leave the variable untracked so isAssigned
			// returns true for all fields without explicit tracking.
			calleeIdent, isDirectCall := call.Callee.(*ast.Ident)
			if !isDirectCall || calleeIdent.Name != st.Name {
				continue
			}
			initLocalVar(a.fieldInit, name.Name, st, call)
			// Propagate sub-field init state for struct-typed args so that
			// chained access (w.inner.field) does not false-positive. This
			// mirrors the same propagation done in analyzeAssignment.
			a.propagateConstructorArgInits(name.Name, st, call)
		}
	}
	return TypUnit
}

// bocLitHasParams reports whether a boc literal has any TypedDecl params
// (nil-value TypedDecls that define the closure's input signature).
func bocLitHasParams(bocLit *ast.BocLiteral) bool {
	for _, elem := range bocLit.Elements {
		if td, ok := elem.(*ast.TypedDecl); ok && td.Value == nil {
			return true
		}
	}
	return false
}

// hasInnerBocsOrMethods reports whether a boc literal contains any inner
// body-form bocs (ShortDecl with BocLiteral value) or BocDecl methods.
// These require StructType recording for correct field-access type-checking.
func hasInnerBocsOrMethods(bocLit *ast.BocLiteral) bool {
	for _, elem := range bocLit.Elements {
		switch e := elem.(type) {
		case *ast.BocDecl:
			if e.Body != nil {
				return true
			}
		case *ast.ShortDecl:
			if len(e.Values) == 1 {
				if _, ok := e.Values[0].(*ast.BocLiteral); ok {
					return true
				}
			}
		}
	}
	return false
}

// analyzeBocDecl handles `name: { ... }` for both lowercase (boc) and
// uppercase (struct type) names.
func (a *Analyzer) analyzeBocDecl(name *ast.Ident, bocLit *ast.BocLiteral, decl ast.Node) Type {
	if bocLit.Annotation != nil {
		a.analyzeAnnotationBody(bocLit.Annotation)
	}
	fqn := a.currentFQN(name.Name)
	prevFQN := a.pushFQN(name.Name)

	isLower := name.TokType != token.TYPE_IDENT && name.TokType != token.GENERIC_IDENT
	hasStructure := isLower && hasInnerBocsOrMethods(bocLit)

	// For lowercase boc definitions, pre-register the symbol in the current
	// (outer) scope before pushing the inner body scope, so that recursive
	// self-calls inside the body can resolve the boc's own name via the
	// parent scope chain.
	if isLower {
		if hasStructure {
			// Pre-register with an approximate StructType (no fields yet — they're
			// not analyzed until the body scope is entered below).
			a.define(&Symbol{
				Name: name.Name,
				Type: &StructType{Name: name.Name, IsSingleton: true, Returns: []Type{TypUnit}},
				FQN:  fqn,
				Node: decl,
			})
		} else {
			preParams := a.collectParams(bocLit.Elements)
			a.define(&Symbol{
				Name: name.Name,
				Type: &BocType{Params: preParams, Returns: []Type{TypUnit}},
				FQN:  fqn,
				Node: decl,
			})
		}
	}

	prev := a.pushScope()

	var typ Type
	if name.TokType == token.TYPE_IDENT || name.TokType == token.GENERIC_IDENT {
		// Uppercase (multi-char TYPE_IDENT or single-letter GENERIC_IDENT): struct type definition.
		st, _ := a.analyzeStructBoc(name.Name, bocLit)
		typ = st
	} else if hasStructure {
		// Lowercase with inner bocs or methods: use analyzeStructBoc for correct
		// field recording and FQN-aware analysis. The returned lastExprTypes give
		// the call return type (what `counter()` produces).
		st, returns := a.analyzeStructBoc(name.Name, bocLit)
		if len(returns) == 0 {
			returns = []Type{TypUnit}
		}
		st.IsSingleton = true
		st.Returns = returns
		typ = st
	} else {
		// Simple lowercase boc (no inner structure): original BocType path.
		prevFI := a.fieldInit
		a.fieldInit = newFieldInitState()
		bt := a.analyzeBocBody(bocLit.Elements)
		a.fieldInit = prevFI
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
			newSym := &Symbol{
				Name:           vc.Name,
				Type:           &BocType{Params: params, Returns: []Type{st}},
				Node:           decl,
				ParentTypeName: st.Name,
			}
			// Detect collision with another variant constructor of the same name.
			if existing := a.currentScope.LookupLocal(vc.Name); existing != nil &&
				(existing.ParentTypeName != "" || len(existing.Alternatives) > 0) {
				if len(existing.Alternatives) == 0 {
					// First collision: seed the alternatives list with the original symbol.
					clone := *existing
					existing.Alternatives = []*Symbol{&clone}
					existing.Type = Unknown
					existing.ParentTypeName = ""
				}
				existing.Alternatives = append(existing.Alternatives, newSym)
			} else {
				a.define(newSym)
			}
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
	if d.Annotation != nil {
		a.analyzeAnnotationBody(d.Annotation)
	}
	typ := a.resolveTypeExpr(d.Type)
	if d.Value != nil {
		prev := a.expectedType
		a.expectedType = typ
		valTyp := a.analyzeExpr(d.Value)
		a.expectedType = prev
		if valTyp != Unknown && !valTyp.IsCompatibleWith(typ) {
			a.errorf(d.Pos, "type mismatch: %s is not compatible with %s", displayType(valTyp), displayType(typ))
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

// memberPath walks a (possibly nested) MemberExpr and returns the root
// variable name and the dotted field path.
// b.inner.field → ("b", "inner.field")
// b.field       → ("b", "field")
// Returns ("", "") if the chain is not rooted at a simple identifier.
func memberPath(expr ast.Expr) (rootVar, dotPath string) {
	var parts []string
	cur := expr
	for {
		switch e := cur.(type) {
		case *ast.MemberExpr:
			parts = append([]string{e.Member.Name}, parts...)
			cur = e.Object
		case *ast.Ident:
			return e.Name, strings.Join(parts, ".")
		default:
			return "", ""
		}
	}
}

func (a *Analyzer) analyzeAssignment(asgn *ast.Assignment) Type {
	// Determine the expected type before analyzing values so that ambiguous
	// variant constructor calls (YZC-0065) can be resolved by type context.
	var targetType Type
	if asgn.Target != nil {
		targetType = a.resolveTargetType(asgn.Target)
	} else if len(asgn.Names) == 1 {
		if sym := a.currentScope.Lookup(asgn.Names[0].Name); sym != nil {
			targetType = sym.Type
		}
	}

	prev := a.expectedType
	a.expectedType = targetType
	var valTypes []Type
	for _, v := range asgn.Values {
		valTypes = append(valTypes, a.analyzeExpr(v))
	}
	a.expectedType = prev

	if asgn.Target != nil {
		// Track field assignment for definite-assignment analysis.
		if a.fieldInit != nil {
			if mem, ok := asgn.Target.(*ast.MemberExpr); ok {
				if varName, path := memberPath(mem); varName != "" {
					a.fieldInit.markAssigned(varName, path)
					// For a direct field assignment (no dots in path) whose RHS is a
					// struct value, propagate the inner struct's init state so that
					// nested reads like b.inner.field don't false-positive. (YZC-0054)
					if !strings.Contains(path, ".") && len(valTypes) > 0 {
						if innerSt, ok := valTypes[0].(*StructType); ok &&
							!innerSt.IsSingleton && !innerSt.IsVariant && !innerSt.IsInterface {
							var innerFI *FieldInitState
							if call, ok := asgn.Values[0].(*ast.CallExpr); ok {
								innerFI = newFieldInitState()
								initLocalVar(innerFI, "_", innerSt, call)
							} else if rhsId, ok := asgn.Values[0].(*ast.Ident); ok {
								if _, tracked := a.fieldInit.locals[rhsId.Name]; tracked {
									innerFI = a.fieldInit // propagateInner reads from rhsId.Name
									a.fieldInit.propagateInner(varName, path, innerFI, rhsId.Name)
									innerFI = nil // already propagated
								} else {
									// Untracked source (parameter) — all fields initialized.
									innerFI = newFieldInitState()
									innerFI.addLocalVar("_")
									for _, f := range innerSt.Fields {
										if !f.HasDefault {
											if _, isMethod := f.Type.(*BocType); !isMethod {
												innerFI.locals["_"][f.Name] = true
											}
										}
									}
								}
							}
							if innerFI != nil {
								a.fieldInit.propagateInner(varName, path, innerFI, "_")
							}
						}
					}
				}
			}
		}
		if len(valTypes) > 0 && targetType != Unknown {
			if !valTypes[0].IsCompatibleWith(targetType) {
				a.errorf(asgn.Pos, "assignment: %s is not compatible with %s", displayType(valTypes[0]), displayType(targetType))
			}
		}
	} else {
		// Expand TupleType for multi-name assignment (YZC-0012).
		effectiveValTypes := valTypes
		if len(asgn.Names) > 1 && len(valTypes) == 1 {
			if tt, ok := valTypes[0].(*TupleType); ok {
				effectiveValTypes = tt.Types
			}
		}
		for i, name := range asgn.Names {
			sym := a.currentScope.Lookup(name.Name)
			if sym == nil {
				a.errorf(name.Pos, "undefined: %s", name.Name)
				continue
			}
			if i < len(effectiveValTypes) && sym.Type != Unknown {
				if !effectiveValTypes[i].IsCompatibleWith(sym.Type) {
					a.errorf(name.Pos, "assignment to %s: %s not compatible with %s",
						name.Name, displayType(effectiveValTypes[i]), displayType(sym.Type))
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
// BocDecl
// ---------------------------------------------------------------------------

func (a *Analyzer) analyzeBocDeclNode(bd *ast.BocDecl) Type {
	// Resolve all params from the signature.
	// In body-only form (= { body }), all params are inputs (no isReturn logic).
	allParams := a.resolveBocSigParams(bd.Sig, bd.BodyOnly)

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
	if bd.Body == nil && (bd.Name.TokType == token.TYPE_IDENT || bd.Name.TokType == token.GENERIC_IDENT) {
		// If every input param is a BocType, this is an interface declaration.
		allBoc := len(inputParams) > 0
		for _, p := range inputParams {
			if _, isBoc := p.Type.(*BocType); !isBoc {
				allBoc = false
				break
			}
		}
		st := &StructType{Name: bd.Name.Name, IsInterface: allBoc}
		for _, p := range inputParams {
			if p.Label != "" {
				st.Fields = append(st.Fields, StructField{Name: p.Label, Type: p.Type})
			}
		}
		fqn := a.currentFQN(bd.Name.Name)
		sym := &Symbol{Name: bd.Name.Name, Type: st, FQN: fqn, Node: bd}
		a.define(sym)
		a.setType(bd, st)
		return st
	}

	var returns []Type
	if bd.Body != nil {
		// Pre-register the symbol before analyzing the body so recursive calls
		// inside the body can resolve the function's own name.
		// Use sig-declared return types if available; otherwise default to Unit.
		// The final symbol (with inferred return type) is re-registered after analysis.
		preReturns := explicitReturns
		if len(preReturns) == 0 {
			preReturns = []Type{TypUnit}
		}
		preFQN := a.currentFQN(bd.Name.Name)
		preSym := &Symbol{
			Name: bd.Name.Name,
			Type: &BocType{Params: inputParams, Returns: preReturns},
			FQN:  preFQN,
			Node: bd,
		}
		a.define(preSym)

		prev := a.pushScope()
		prevFQN := a.pushFQN(bd.Name.Name)
		prevFI := a.fieldInit
		a.fieldInit = newFieldInitState()
		if bd.BodyOnly {
			// Body-only form: extract named params from body's initial TypedDecls.
			// The body redeclares its own params; sig provides types (and names
			// when labeled) for validation.
			bodyParams, n := a.extractBodyParams(bd.Body.Elements, inputParams)
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
					shorthandParams := a.resolveBocSigParams(bd.Sig, false)
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
					if len(bd.Body.Elements) > 0 {
						a.errorf(bd.Body.Elements[0].Position(),
							"expected parameter declaration (name Type), got %T", bd.Body.Elements[0])
					}
				}
			}
			for _, p := range inputParams {
				if p.Label != "" {
					a.currentScope.Define(&Symbol{Name: p.Label, Type: p.Type})
				}
			}
			// Analyze body starting after the param declarations.
			bodyElems := bd.Body.Elements
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
			returns = a.analyzeBocBody(bd.Body.Elements)
		}
		a.fieldInit = prevFI
		a.popFQN(prevFQN)
		a.popScope(prev)
	}

	// Explicit return types in the signature override inferred returns.
	// Check that the body's inferred return type is compatible with the declaration.
	if len(explicitReturns) > 0 && bd.Body != nil {
		var bodyReturn Type = TypUnit
		if len(returns) > 0 {
			bodyReturn = returns[0]
		}
		declared := explicitReturns[0]
		_, bodyIsUnknown := bodyReturn.(*UnknownType)
		_, declIsUnknown := declared.(*UnknownType)
		_, declIsPDT := declared.(*PathDependentType)
		if !bodyIsUnknown && !declIsUnknown && !declIsPDT && !bodyReturn.IsCompatibleWith(declared) {
			a.errorf(bd.Name.Pos, "YZC-0035: boc body returns %s but declared output is %s",
				displayType(bodyReturn), displayType(declared))
		}
		returns = explicitReturns
	}
	// For abstract (body-less) method declarations (interface members), use the
	// signature's explicit return types rather than defaulting to Unit.
	if len(returns) == 0 && bd.Body == nil && len(explicitReturns) > 0 {
		returns = explicitReturns
	}
	if len(returns) == 0 {
		returns = []Type{TypUnit}
	}
	// When a method body's sole inferred return type is Unknown (e.g. a call on
	// a generic type param: value.hola()) and there is no explicit return type,
	// treat the method as returning Unit. Unknown here means "can't resolve" not
	// "returns a concrete non-Unit type", and the method is called for side effects.
	if len(returns) == 1 && len(explicitReturns) == 0 && bd.Body != nil {
		if _, isUnknown := returns[0].(*UnknownType); isUnknown {
			returns = []Type{TypUnit}
		}
	}

	bocType := &BocType{Params: inputParams, Returns: returns}
	fqn := a.currentFQN(bd.Name.Name)
	sym := &Symbol{Name: bd.Name.Name, Type: bocType, FQN: fqn, Node: bd}
	a.define(sym)
	a.setType(bd, bocType)
	return bocType
}

// resolveBocSigParams resolves the params of a BocTypeExpr signature.
// When bodyOnly is true (the `= { body }` form), all params are treated as
// inputs: unlabeled types are anonymous inputs, not return-type annotations.
func (a *Analyzer) resolveBocSigParams(sig *ast.BocTypeExpr, bodyOnly bool) []BocParam {
	// Populate currentSigParams so that later params can reference earlier ones
	// in path-dependent type expressions (e.g. `n g.Node` where `g` precedes `n`).
	prevSigParams := a.currentSigParams
	a.currentSigParams = make(map[string]Type)
	defer func() { a.currentSigParams = prevSigParams }()

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
		// Register input param for subsequent path-dependent lookups.
		if p.Label != "" && !isReturn {
			a.currentSigParams[p.Label] = typ
		}
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
			// Accumulate consecutive non-Unit trailing expressions — each is a
			// return value. A Unit expression (side effect) resets the sequence.
			// Multiple trailing non-Unit exprs → multi-return boc (YZC-0012).
			if t == TypUnit {
				lastExprTypes = nil
			} else {
				lastExprTypes = append(lastExprTypes, t)
			}
		case *ast.ReturnStmt:
			lastExprTypes = []Type{t}
		default:
			// Statements don't contribute to return type.
			lastExprTypes = nil
		}
	}
	return lastExprTypes
}

// analyzeBranchBody analyzes a BocLiteral as a control-flow branch (ConditionalExpr
// arm or match arm). Unlike the BocLiteral case in analyzeExpr, this does NOT
// save/restore fieldInit — the caller manages cloning and merging.
func (a *Analyzer) analyzeBranchBody(boc *ast.BocLiteral) Type {
	prev := a.pushScope()
	bodyReturns := a.analyzeBocBody(boc.Elements)
	params := a.collectParams(boc.Elements)
	a.popScope(prev)
	if len(bodyReturns) == 0 {
		bodyReturns = []Type{TypUnit}
	}
	bt := &BocType{Params: params, Returns: bodyReturns}
	a.setType(boc, bt)
	return bt
}

// ---------------------------------------------------------------------------
// Struct type analysis (uppercase boc)
// ---------------------------------------------------------------------------

// analyzeStructBoc analyzes a boc literal as a struct type, returning the
// struct type and the last-expression types (body return types).
// It is used for both uppercase struct declarations and lowercase singleton
// bocs that have inner structure (inner bocs or BocDecl methods).
func (a *Analyzer) analyzeStructBoc(name string, b *ast.BocLiteral) (*StructType, []Type) {
	// Reuse the stub pre-registered by AnalyzeFile so that forward references
	// (fields declared before their type is analyzed) capture a stable pointer
	// that gets filled in here.
	st := &StructType{Name: name}
	if sym := a.fileScope.LookupLocal(name); sym != nil {
		if stub, ok := sym.Type.(*StructType); ok {
			st = stub
		}
	}
	fieldSet := make(map[string]bool)
	var lastExprTypes []Type
	hasBocBody := false // true when any BocDecl element has a method body

	// Pre-scan: detect whether this is a generic struct so we can activate
	// constraint collection before processing method bodies.
	isGeneric := false
	for _, elem := range b.Elements {
		if id, ok := elem.(*ast.Ident); ok && id.TokType == token.GENERIC_IDENT {
			isGeneric = true
			break
		}
		if _, ok := elem.(*ast.TypeParamDecl); ok {
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
			lastExprTypes = nil

		case *ast.ShortDecl:
			// Use analyzeShortDecl so inner body-form bocs get proper setType()
			// registration (needed by lowerMethod to query ExprType on ShortDecl).
			// A lowercase ShortDecl with a BocLiteral value is a method body (YZC-0082).
			if len(e.Names) == 1 && len(e.Values) == 1 {
				if _, isBocLit := e.Values[0].(*ast.BocLiteral); isBocLit && !isUppercaseName(e.Names[0].Name) {
					hasBocBody = true
				}
			}
			a.analyzeShortDecl(e)
			// Collect the resulting field(s) into the struct.
			for _, n := range e.Names {
				sym := a.currentScope.LookupLocal(n.Name)
				if sym == nil {
					continue
				}
				if fieldSet[n.Name] {
					a.errorf(e.Pos, "duplicate field %q in %s", n.Name, name)
					continue
				}
				fieldSet[n.Name] = true
				// Step 6 (YZC-0066): detect type-alias bindings.
				// `Node: User` — uppercase LHS + bare ident RHS resolving to a
				// non-singleton, non-interface StructType → compile-time type alias.
				// Also marks nested type definitions: `Window: { size Int }` inside a
				// singleton/struct (YZC-0081).
				isTypeAlias := false
				if isUppercaseName(n.Name) && len(e.Values) == 1 {
					if _, ok := e.Values[0].(*ast.Ident); ok {
						if st2, ok2 := sym.Type.(*StructType); ok2 && !st2.IsSingleton && !st2.IsInterface && !st2.IsVariant {
							isTypeAlias = true
						}
					}
					if _, ok := e.Values[0].(*ast.BocLiteral); ok {
						if st2, ok2 := sym.Type.(*StructType); ok2 && !st2.IsSingleton && !st2.IsInterface && !st2.IsVariant {
							isTypeAlias = true
						}
					}
				}
				var defaultExpr ast.Expr
				if len(e.Values) == 1 && !isTypeAlias {
					defaultExpr = e.Values[0]
				}
				st.Fields = append(st.Fields, StructField{Name: n.Name, Type: sym.Type, HasDefault: true, DefaultExpr: defaultExpr, IsTypeField: isTypeAlias})
			}
			lastExprTypes = nil

		case *ast.TypeParamDecl:
			if e.Name.TokType == token.TYPE_IDENT {
				// YZC-0074: `Node Sizer` — named associated type with a named-type bound.
				var bound Type
				if len(e.Constraints) == 1 {
					bound = a.resolveTypeExpr(e.Constraints[0])
				}
				if !fieldSet[e.Name.Name] {
					fieldSet[e.Name.Name] = true
					st.Fields = append(st.Fields, StructField{Name: e.Name.Name, Type: TypMeta, IsTypeField: true, Bound: bound})
				}
				lastExprTypes = nil
				continue
			}
			// Constrained type param declaration: `V Talker`, `T A B`, or `V #(m #(T))`.
			a.registerTypeParam(st, &fieldSet, e.Name.Name)
			a.storeExplicitConstraints(st, e.Name.Name, e.Constraints)
			if e.InlineConstraint != nil {
				a.storeInlineConstraint(st, name, e.Name.Name, e.InlineConstraint)
			}
			gt := &GenericType{Name: e.Name.Name}
			a.currentScope.Define(&Symbol{Name: e.Name.Name, Type: gt, Node: e.Name})
			lastExprTypes = nil

		case *ast.Ident:
			// Generic type param declaration (T, E inside type boc body).
			// Register as GenericType in current scope and record on the struct.
			if e.TokType == token.GENERIC_IDENT {
				a.registerTypeParam(st, &fieldSet, e.Name)
			}
			gt := &GenericType{Name: e.Name}
			a.currentScope.Define(&Symbol{Name: e.Name, Type: gt, Node: e})
			lastExprTypes = nil

		case *ast.BocDecl:
			// `TypeName #(...)` with no body: abstract associated type field.
			// Empty sig (#()) = unconstrained; non-empty sig (#(method #(T))) = constrained bound.
			// Equivalent to an associated type in Rust/Swift.
			if e.Name.TokType == token.TYPE_IDENT && e.Body == nil && e.Sig != nil {
				var bound Type
				if len(e.Sig.Params) > 0 {
					// YZC-0074: inline constraint — build an anonymous interface from the sig params.
					bound = a.buildAssocTypeBound(name, e.Name.Name, e.Sig)
				}
				if !fieldSet[e.Name.Name] {
					fieldSet[e.Name.Name] = true
					st.Fields = append(st.Fields, StructField{Name: e.Name.Name, Type: TypMeta, IsTypeField: true, Bound: bound})
				}
				lastExprTypes = nil
				continue
			}
			// Track whether any method has a body (used in YZC-0067 interface check).
			if e.Body != nil {
				hasBocBody = true
			}
			// Set the constraint attribution context before analyzing the method body
			// so that any T-method calls recorded during analysis are tagged correctly.
			if a.activeConstraints != nil {
				a.activeContext = name + "." + e.Name.Name
			}
			typ := a.analyzeBocDeclNode(e)
			if !fieldSet[e.Name.Name] {
				fieldSet[e.Name.Name] = true
				st.Fields = append(st.Fields, StructField{Name: e.Name.Name, Type: typ})
			}
			lastExprTypes = nil

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
			lastExprTypes = nil

		default:
			// Expressions and other nodes — track last expression type for return type inference.
			t := a.analyzeNode(elem)
			if _, ok := elem.(ast.Expr); ok {
				lastExprTypes = []Type{t}
			} else if _, ok2 := elem.(*ast.ReturnStmt); ok2 {
				lastExprTypes = []Type{t}
			} else {
				lastExprTypes = nil
			}
		}
	}

	// Step 5 (YZC-0066): auto-collect implicit generic type vars into TypeParams.
	// If T appears in a field type but was never declared on its own line (no
	// IsTypeField entry), it still needs to be in TypeParams so the emitted Go
	// struct is generic (e.g. Box: { value T } → Box[T any]).
	{
		inTypeParams := make(map[string]bool, len(st.TypeParams))
		for _, tp := range st.TypeParams {
			inTypeParams[tp] = true
		}
		// Process fields in declaration order for deterministic TypeParams ordering.
		for _, f := range st.Fields {
			if f.IsTypeField {
				continue
			}
			discovered := make(map[string]bool)
			collectGenericNames(f.Type, discovered)
			for gname := range discovered {
				if !inTypeParams[gname] {
					st.TypeParams = append(st.TypeParams, gname)
					inTypeParams[gname] = true
				}
			}
		}
	}

	// YZC-0073: synthesize anonymous interface constraints for user-defined method
	// calls on generic type params when no named interface is in scope.
	// E.g. value.hola() where V has no named Holer → synthesize _StructVConstraint.
	if isGeneric && a.activeConstraints != nil {
		a.synthesizeConstraints(st, name)
	}

	// Freeze the inferred constraints into the struct type and restore outer state.
	if isGeneric {
		if len(a.activeConstraints) > 0 {
			st.TypeConstraints = a.activeConstraints
		}
		a.activeConstraints = prevConstraints
		a.activeContext = prevContext
	}

	// YZC-0067: a struct whose fields are only abstract type fields (MetaType)
	// and/or BocType method declarations — no concrete data — is a Go interface.
	// A struct with any method body is NOT an interface (it has an implementation).
	if !st.IsVariant && !st.IsSingleton {
		isIface := !hasBocBody
		if isIface {
			for _, f := range st.Fields {
				if f.IsTypeField {
					if _, isMeta := f.Type.(*MetaType); !isMeta {
						isIface = false
						break
					}
					continue
				}
				if _, isBoc := f.Type.(*BocType); !isBoc {
					isIface = false
					break
				}
			}
		}
		st.IsInterface = isIface
	}

	return st, lastExprTypes
}

func containsString(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
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
			if part.IsExpr && a.inAnnotation {
				a.errorf(part.Expr.Position(), "YZC-0025: string interpolation is not allowed inside annotations")
				continue
			}
			if part.IsExpr {
				partType := a.analyzeExpr(part.Expr)
				// YZC-0046: ${} requires to_str; backtick form accepts any type.
				if !part.IsDebug && !typeHasToStr(partType) {
					a.errorf(part.Expr.Position(), "YZC-0046: %s requires to implement to_str #(String)", displayType(partType))
				}
			}
		}
		t = TypString
	case *ast.ConditionalExpr:
		a.analyzeExpr(expr.Cond)
		var trueType Type
		if a.fieldInit != nil {
			// Branch analysis: clone state for each branch, intersect after.
			preState := a.fieldInit
			a.fieldInit = preState.clone()
			if trueBoc, ok := expr.TrueCase.(*ast.BocLiteral); ok {
				trueType = a.analyzeBranchBody(trueBoc)
			} else {
				trueType = a.analyzeExpr(expr.TrueCase)
			}
			trueAfter := a.fieldInit
			a.fieldInit = preState.clone()
			if falseBoc, ok := expr.FalseCase.(*ast.BocLiteral); ok {
				a.analyzeBranchBody(falseBoc)
			} else {
				a.analyzeExpr(expr.FalseCase)
			}
			falseAfter := a.fieldInit
			trueAfter.intersect(falseAfter)
			a.fieldInit = trueAfter
		} else {
			trueType = a.analyzeExpr(expr.TrueCase)
			a.analyzeExpr(expr.FalseCase)
		}
		// The ? operator calls the branch; the result is the branch's return value,
		// not the branch boc itself. Unwrap one BocType level.
		if bt, ok := trueType.(*BocType); ok && len(bt.Returns) == 1 {
			t = bt.Returns[0]
		} else {
			t = trueType
		}
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
		outerFI := a.fieldInit
		if hasInnerBocsOrMethods(expr) && !bocLitHasParams(expr) {
			// Anonymous boc literal with inner methods and no params: type as anonymous StructType.
			// If the literal has TypedDecl params it is a closure regardless of inner named bocs.
			prev := a.pushScope()
			if outerFI != nil {
				a.fieldInit = outerFI.clone()
			}
			st, _ := a.analyzeStructBoc("_anonBoc", expr)
			a.popScope(prev)
			a.fieldInit = outerFI
			st.IsSingleton = true
			t = st
		} else {
			prev := a.pushScope()
			// Closures get an isolated copy of the field-init state so that
			// assignments inside the closure don't affect the outer scope.
			if outerFI != nil {
				a.fieldInit = outerFI.clone()
			}
			bodyReturns := a.analyzeBocBody(expr.Elements)
			params := a.collectParams(expr.Elements)
			a.popScope(prev)
			a.fieldInit = outerFI // restore outer state
			if len(bodyReturns) == 0 {
				bodyReturns = []Type{TypUnit}
			}
			t = &BocType{Params: params, Returns: bodyReturns}
		}
	case *ast.ArrayLiteral:
		t = a.analyzeArrayLiteral(expr)
	case *ast.DictLiteral:
		t = a.analyzeDictLiteral(expr)
	case *ast.MatchExpr:
		t = a.analyzeMatch(expr)
	case *ast.InfixMatchExpr:
		t = a.analyzeInfixMatch(expr)
	case *ast.Annotation:
		a.analyzeAnnotationBody(expr)
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
		a.errorf(u.Pos, "unary '-' not defined for type %s", displayType(operandType))
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
				// Analyze args first so ParamTypes are available for constraint synthesis.
				var paramTypes []Type
				for _, arg := range c.Args {
					paramTypes = append(paramTypes, a.analyzeExpr(arg.Value))
				}
				a.activeConstraints[gt.Name] = append(a.activeConstraints[gt.Name], &GenericConstraint{
					TypeParam:  gt.Name,
					MethodName: memExpr.Member.Name,
					ParamTypes: paramTypes,
					Line:       memExpr.Member.Pos.Line,
					Col:        memExpr.Member.Pos.Col,
					Context:    a.activeContext,
				})
				// Try to infer the return type from a matching interface method signature.
				// Default to Unit when no named interface is found — user-defined methods
				// called for side effects (discarded result) return Unit (YZC-0071/0073).
				retType := a.findInterfaceMethodReturnType(memExpr.Member.Name)
				a.setType(c.Callee, &BocType{Returns: []Type{retType}})
				return retType
			}
		}
	}

	calleeType := a.analyzeExpr(c.Callee)

	// Ambiguous variant constructor (YZC-0065): callee is an identifier that maps
	// to multiple variant constructors with the same name. Use expected type to pick.
	if calleeType == Unknown {
		if id, ok := c.Callee.(*ast.Ident); ok {
			if sym := a.currentScope.Lookup(id.Name); sym != nil && len(sym.Alternatives) > 0 {
				calleeType = a.disambiguateConstructor(c, id.Name, sym.Alternatives)
			}
		}
	}

	// Collect arg types — needed for generic constraint checking at instantiation.
	// Per-argument expectedType lets nested calls resolve their own ambiguities.
	var argTypes []Type
	var formalParams []BocParam
	if bt, ok := calleeType.(*BocType); ok {
		formalParams = bt.Params
	}
	for i, arg := range c.Args {
		prev := a.expectedType
		if i < len(formalParams) {
			a.expectedType = formalParams[i].Type
		} else {
			a.expectedType = nil
		}
		argTypes = append(argTypes, a.analyzeExpr(arg.Value))
		a.expectedType = prev
	}
	// Boc-boundary check (YZC-0053): struct-typed args must have all required
	// fields definitely assigned before crossing the call boundary.
	if a.fieldInit != nil {
		for _, arg := range c.Args {
			id, ok := arg.Value.(*ast.Ident)
			if !ok {
				continue
			}
			sym := a.currentScope.Lookup(id.Name)
			if sym == nil {
				continue
			}
			st, ok := sym.Type.(*StructType)
			if !ok || st.IsSingleton {
				continue
			}
			for _, f := range st.Fields {
				if f.IsTypeField {
					continue
				}
				if f.HasDefault {
					continue
				}
				if _, isMethod := f.Type.(*BocType); isMethod {
					continue
				}
				if !a.fieldInit.isAssigned(id.Name, f.Name) {
					a.errorf(arg.Value.Position(), "YZC-0034: field %s of %s not initialized before call", f.Name, id.Name)
				}
			}
		}
	}
	switch bt := calleeType.(type) {
	case *BocType:
		if len(bt.Returns) == 0 {
			return TypUnit
		}
		// Build a label→argType map for path-dependent resolution.
		labelToArgType := make(map[string]Type)
		for i, param := range bt.Params {
			if param.Label != "" && i < len(argTypes) {
				labelToArgType[param.Label] = argTypes[i]
			}
		}
		// General argument type checking against formal params.
		for i, param := range bt.Params {
			if param.IsReturn || i >= len(argTypes) {
				continue
			}
			argT := argTypes[i]
			paramT := param.Type
			if paramT == Unknown || argT == Unknown {
				continue
			}
			switch paramT.(type) {
			case *GenericType, *MetaType, *PathDependentType:
				continue
			}
			if !argT.IsCompatibleWith(paramT) {
				pos := c.Callee.Position()
				if i < len(c.Args) {
					pos = c.Args[i].Value.Position()
				}
				label := param.Label
				if label == "" {
					label = fmt.Sprintf("%d", i)
				}
				a.errorf(pos, "argument %s: %s is not compatible with %s",
					label, displayType(argT), displayType(paramT))
			}
		}
		// YZC-0030: resolve PathDependentType params and type-check their args.
		for i, param := range bt.Params {
			pdt, ok := param.Type.(*PathDependentType)
			if !ok || i >= len(argTypes) {
				continue
			}
			objType, found := labelToArgType[pdt.Param]
			if !found {
				continue
			}
			if st, ok := objType.(*StructType); ok {
				for _, f := range st.Fields {
					if f.Name == pdt.Member && f.IsTypeField {
						if _, isMeta := f.Type.(*MetaType); !isMeta {
							// Concrete alias (Node: User): check arg is compatible.
							if !argTypes[i].IsCompatibleWith(f.Type) {
								a.errorf(c.Args[i].Value.Position(),
									"YZC-0030: argument type %s is not compatible with %s.%s (expected %s)",
									argTypes[i].typeName(), pdt.Param, pdt.Member, f.Type.typeName())
							}
						} else {
							// Abstract type field (MetaType): g has an abstract type, so g.Node
							// is structurally equivalent to its bound. Check same-path PDTs for
							// cross-root mismatches; for concrete values, check against the bound.
							if argPdt, ok := argTypes[i].(*PathDependentType); ok {
								if argPdt.Param != pdt.Param || argPdt.Member != pdt.Member {
									a.errorf(c.Args[i].Value.Position(),
										"YZC-0030: argument type %s is not compatible with %s.%s",
										argPdt.typeName(), pdt.Param, pdt.Member)
								}
							} else if argTypes[i] != Unknown && f.Bound != nil {
								// YZC-0079: g.Node ≡ its bound in a structural type system.
								// Any value satisfying the bound is a valid g.Node.
								if !argTypes[i].IsCompatibleWith(f.Bound) {
									a.errorf(c.Args[i].Value.Position(),
										"YZC-0079: argument type %s does not satisfy %s.%s bound %s",
										displayType(argTypes[i]), pdt.Param, pdt.Member, displayType(f.Bound))
								}
							}
						}
					}
				}
			}
		}
		// YZC-0074: verify associated type bounds when a concrete graph type is passed
		// as a formal interface param that declares bounded associated types.
		for i, param := range bt.Params {
			if param.IsReturn || i >= len(argTypes) {
				continue
			}
			paramSt, ok := param.Type.(*StructType)
			if !ok {
				continue
			}
			argSt, ok := argTypes[i].(*StructType)
			if !ok {
				continue
			}
			for _, pf := range paramSt.Fields {
				if !pf.IsTypeField || pf.Bound == nil {
					continue
				}
				// Find the matching type field in the actual arg struct.
				for _, af := range argSt.Fields {
					if af.Name != pf.Name || !af.IsTypeField {
						continue
					}
					if _, isMeta := af.Type.(*MetaType); !isMeta {
						// Concrete binding: verify it satisfies the bound.
						if !af.Type.IsCompatibleWith(pf.Bound) {
							pos := c.Callee.Position()
							if i < len(c.Args) {
								pos = c.Args[i].Value.Position()
							}
							a.errorf(pos,
								"YZC-0074: %s.%s (type %s) does not satisfy the required bound %s",
								argSt.Name, af.Name, af.Type.typeName(), pf.Bound.typeName())
						}
					}
					break
				}
			}
		}
		// Phase C (YZC-0066): unify formal params against actual arg types to
		// resolve any generic type variables and path-dependent types in the return.
		bindings := make(map[string]Type)
		for i, param := range bt.Params {
			if i < len(argTypes) {
				unifyTypes(param.Type, argTypes[i], bindings)
			}
		}
		// YZC-0030: add path-dependent bindings ("g.Node" → concrete type)
		// from input params typed as PathDependentType (e.g. n g.Node).
		for i, param := range bt.Params {
			pdt, ok := param.Type.(*PathDependentType)
			if !ok || i >= len(argTypes) {
				continue
			}
			objType, found := labelToArgType[pdt.Param]
			if !found {
				continue
			}
			if st, ok := objType.(*StructType); ok {
				for _, f := range st.Fields {
					if f.Name == pdt.Member && f.IsTypeField {
						if _, isMeta := f.Type.(*MetaType); !isMeta {
							bindings[pdt.Param+"."+pdt.Member] = f.Type
						}
					}
				}
			}
		}
		// Also build bindings from path-dependent RETURN types so that
		// `makeNode #(g Graph, g.Node)` called with a concrete Graph resolves
		// the return type to the concrete Node type (e.g. *User).
		for _, ret := range bt.Returns {
			pdt, ok := ret.(*PathDependentType)
			if !ok {
				continue
			}
			objType, found := labelToArgType[pdt.Param]
			if !found {
				continue
			}
			if st, ok := objType.(*StructType); ok {
				for _, f := range st.Fields {
					if f.Name == pdt.Member && f.IsTypeField {
						if _, isMeta := f.Type.(*MetaType); !isMeta {
							bindings[pdt.Param+"."+pdt.Member] = f.Type
						}
					}
				}
			}
		}
		retType := substituteType(bt.Returns[0], bindings)
		// When the return type is a partially-instantiated generic (e.g. Result[Int, E])
		// and the call site has an expectedType that fully instantiates it, fill in the
		// remaining unbound type params so the lowerer can emit explicit Go type args.
		if git, ok := retType.(*GenericInstType); ok && a.expectedType != nil {
			if expGit, ok := a.expectedType.(*GenericInstType); ok && expGit.Name == git.Name && len(expGit.TypeArgs) == len(git.TypeArgs) {
				filled := make([]Type, len(git.TypeArgs))
				changed := false
				for i, ta := range git.TypeArgs {
					if _, isUnbound := ta.(*GenericType); isUnbound {
						if _, expUnbound := expGit.TypeArgs[i].(*GenericType); !expUnbound {
							filled[i] = expGit.TypeArgs[i]
							changed = true
							continue
						}
					}
					filled[i] = ta
				}
				if changed {
					retType = &GenericInstType{Name: git.Name, TypeArgs: filled}
				}
			}
		}
		if len(bt.Returns) == 1 {
			return retType
		}
		// Multi-return: return TupleType of all substituted return types (YZC-0012).
		allRets := make([]Type, len(bt.Returns))
		for i, r := range bt.Returns {
			allRets[i] = substituteType(r, bindings)
		}
		return &TupleType{Types: allRets}
	case *StructType:
		if bt.IsSingleton {
			// Calling a singleton boc runs the body; return type is the body's last expr.
			if len(bt.Returns) > 0 {
				return bt.Returns[0]
			}
			return TypUnit
		}
		// Uppercase constructor call: verify generic constraints, return the struct type.
		if len(bt.TypeParams) > 0 && len(bt.TypeConstraints) > 0 {
			a.checkGenericConstraints(c.Callee.Position(), bt, argTypes)
		}
		// For generic structs, infer concrete type args from constructor args so that
		// passing the result to a generic HOF can unify the type variables.
		if len(bt.TypeParams) > 0 {
			// Collect data fields (skip IsTypeField and method fields).
			var dataFields []StructField
			labeledFields := make(map[string]StructField)
			for _, f := range bt.Fields {
				if f.IsTypeField {
					continue
				}
				if _, isMethod := f.Type.(*BocType); isMethod {
					continue
				}
				dataFields = append(dataFields, f)
				labeledFields[f.Name] = f
			}
			bindings := make(map[string]Type)
			dataIdx := 0
			for i, arg := range c.Args {
				var ftype Type
				if arg.Label != "" {
					if f, ok := labeledFields[arg.Label]; ok {
						ftype = f.Type
					}
				} else if dataIdx < len(dataFields) {
					ftype = dataFields[dataIdx].Type
					dataIdx++
				}
				if ftype != nil && i < len(argTypes) {
					unifyTypes(ftype, argTypes[i], bindings)
				}
			}
			typeArgs := make([]Type, len(bt.TypeParams))
			for j, tp := range bt.TypeParams {
				if bound, ok := bindings[tp]; ok {
					typeArgs[j] = bound
				} else {
					typeArgs[j] = &GenericType{Name: tp}
				}
			}
			return &GenericInstType{Name: bt.Name, TypeArgs: typeArgs}
		}
		return bt // constructor call (non-generic struct)
	case *BuiltinType:
		return bt // direct type value used as function
	case *GenericInstType:
		// Generic instantiation alias used as constructor: StringBox(value:"hello")
		// where StringBox : Box(String). The return type is the same instantiation.
		return bt
	}
	return Unknown
}

// disambiguateConstructor resolves an ambiguous constructor call using the
// current expectedType. Returns the BocType of the chosen constructor, or
// Unknown after emitting an error when no match is found.
func (a *Analyzer) disambiguateConstructor(c *ast.CallExpr, name string, alts []*Symbol) Type {
	// Try to match against expectedType.
	if a.expectedType != nil {
		if expSt, ok := a.expectedType.(*StructType); ok {
			for _, alt := range alts {
				if bt, ok := alt.Type.(*BocType); ok && len(bt.Returns) == 1 {
					if retSt, ok := bt.Returns[0].(*StructType); ok && retSt.Name == expSt.Name {
						a.setType(c, retSt)
						return bt
					}
				}
			}
		}
	}
	// Still ambiguous — build error listing parent type names.
	var parents []string
	seen := map[string]bool{}
	for _, alt := range alts {
		if alt.ParentTypeName != "" && !seen[alt.ParentTypeName] {
			parents = append(parents, alt.ParentTypeName)
			seen[alt.ParentTypeName] = true
		}
	}
	a.errorf(c.Callee.Position(), "YZC-0065: %s is defined in %s; add a type annotation or use %s.%s(...)",
		name, strings.Join(parents, " and "), parents[0], name)
	return Unknown
}

func (a *Analyzer) analyzeMember(m *ast.MemberExpr) Type {
	objType := a.analyzeExpr(m.Object)
	// Definite-assignment check: required fields of locally-constructed structs
	// must be assigned before they are read. Handles nested paths (YZC-0054).
	if a.fieldInit != nil {
		if varName, path := memberPath(m); varName != "" {
			if sym := a.currentScope.Lookup(varName); sym != nil {
				if _, ok := sym.Type.(*StructType); ok {
					// Only check non-method fields (methods are never tracked).
					skipCheck := false
					if st, ok := objType.(*StructType); ok {
						for _, f := range st.Fields {
							if f.Name == m.Member.Name {
								if _, isBoc := f.Type.(*BocType); isBoc {
									skipCheck = true // method field
								}
								if f.IsTypeField {
									skipCheck = true // compile-time type field; always available
								}
								break
							}
						}
					}
					if !skipCheck && !a.fieldInit.isAssigned(varName, path) {
						a.errorf(m.Member.Pos, "YZC-0034: field %s used before initialization", m.Member.Name)
					}
				}
			}
		}
	}
	return a.fieldType(objType, m.Member.Name, m.Pos)
}

func (a *Analyzer) fieldType(objType Type, fieldName string, pos ast.Pos) Type {
	switch ot := objType.(type) {
	case *GenericInstType:
		// e.g. b: Box(value:42) has type GenericInstType{Box,[Int]}.
		// Look up the base struct and substitute type params to get the concrete field type.
		sym := a.currentScope.Lookup(ot.Name)
		if sym == nil {
			return Unknown
		}
		st, ok := sym.Type.(*StructType)
		if !ok {
			return Unknown
		}
		subst := make(map[string]Type)
		for i, tp := range st.TypeParams {
			if i < len(ot.TypeArgs) {
				subst[tp] = ot.TypeArgs[i]
			}
		}
		for _, f := range st.Fields {
			if f.Name == fieldName {
				return substituteType(f.Type, subst)
			}
		}
		a.errorf(pos, "type %s has no field %q", displayType(objType), fieldName)
		return Unknown
	case *StructType:
		for _, f := range ot.Fields {
			if f.Name == fieldName {
				return f.Type
			}
		}
		// Qualified variant constructor: Shape.Circle — look up by constructor name.
		if ot.IsVariant {
			for _, vc := range ot.Variants {
				if vc.Name == fieldName {
					params := make([]BocParam, len(vc.Fields))
					for j, f := range vc.Fields {
						params[j] = BocParam{Label: f.Name, Type: f.Type}
					}
					return &BocType{Params: params, Returns: []Type{ot}}
				}
			}
		}
		a.errorf(pos, "type %s has no field %q", displayType(objType), fieldName)
		return Unknown
	case *BuiltinType:
		if methods, ok := builtinMethods[ot.name]; ok {
			if ret, ok := methods[fieldName]; ok {
				return ret
			}
		}
		a.errorf(pos, "type %s has no method %q", displayType(objType), fieldName)
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
	case *PathDependentType:
		// YZC-0074: `node.label()` where node has type `g.Node` — resolve via bound.
		// Look up the param `g` in scope to find Graph's Node field and its bound.
		var paramType Type
		if sym := a.currentScope.Lookup(ot.Param); sym != nil {
			paramType = sym.Type
		} else if a.currentSigParams != nil {
			paramType = a.currentSigParams[ot.Param]
		}
		if paramType != nil {
			if st, ok := paramType.(*StructType); ok {
				for _, f := range st.Fields {
					if f.Name == ot.Member && f.IsTypeField && f.Bound != nil {
						return a.fieldType(f.Bound, fieldName, pos)
					}
				}
			}
		}
		a.errorf(pos, "type %s has no field %q (associated type has no bound)", ot.typeName(), fieldName)
		return Unknown
	case *OptionType:
		// .value → the inner type; .to_str → String
		switch fieldName {
		case "value":
			return ot.Inner
		case "to_str":
			return &BocType{Returns: []Type{TypString}}
		}
		return Unknown
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
		return &OptionType{Inner: ot.Val}
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

// isOptionVariantCondition reports whether id is a valid Option variant name.
func isOptionVariantCondition(id *ast.Ident) bool {
	return id.Name == "Some" || id.Name == "None"
}

func (a *Analyzer) analyzeMatch(m *ast.MatchExpr) Type {
	var subjIsOption bool
	if m.Subject != nil {
		subjType := a.analyzeExpr(m.Subject)
		_, subjIsOption = subjType.(*OptionType)
	}
	var returnType Type = Unknown

	if a.fieldInit == nil {
		for _, arm := range m.Arms {
			if arm.Condition != nil {
				// Skip condition lookup for built-in Option variant names.
				if id, ok := arm.Condition.(*ast.Ident); ok && subjIsOption && isOptionVariantCondition(id) {
					// valid — no lookup needed
				} else {
					a.analyzeExpr(arm.Condition)
				}
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

	// With fieldInit tracking: analyze each arm with a clone of the pre-match
	// state, then intersect all arm post-states (all arms cover all paths since
	// a default arm is always present in conditional-match form).
	preState := a.fieldInit
	var afterStates []*FieldInitState

	for _, arm := range m.Arms {
		if arm.Condition != nil {
			if id, ok := arm.Condition.(*ast.Ident); ok && subjIsOption && isOptionVariantCondition(id) {
				// built-in Option variant — no lookup needed
			} else {
				a.analyzeExpr(arm.Condition)
			}
		}
		var armType Type = TypUnit
		a.fieldInit = preState.clone()
		prev := a.pushScope()
		for _, elem := range arm.Body {
			t := a.analyzeNode(elem)
			if _, ok := elem.(ast.Expr); ok {
				armType = t
			}
		}
		a.popScope(prev)
		afterStates = append(afterStates, a.fieldInit)
		if returnType == Unknown {
			returnType = armType
		}
	}

	if len(afterStates) > 0 {
		merged := afterStates[0]
		for _, s := range afterStates[1:] {
			merged.intersect(s)
		}
		a.fieldInit = merged
	} else {
		a.fieldInit = preState
	}

	return returnType
}

func (a *Analyzer) analyzeInfixMatch(m *ast.InfixMatchExpr) Type {
	subjType := a.analyzeExpr(m.Subject)

	// Built-in Option type: accept Some/None as constructors.
	if _, ok := subjType.(*OptionType); ok {
		if !isOptionVariantCondition(m.Constructor) {
			a.errorfLen(m.Constructor.Pos, len(m.Constructor.Name), "%s is not a constructor of Option", m.Constructor.Name)
			return Unknown
		}
		if m.Body == nil {
			return TypBool
		}
		var bodyType Type = TypUnit
		prev := a.pushScope()
		for _, elem := range m.Body.Elements {
			t := a.analyzeNode(elem)
			if _, ok := elem.(ast.Expr); ok {
				bodyType = t
			}
		}
		a.popScope(prev)
		if m.ElseBody != nil {
			prev2 := a.pushScope()
			for _, elem := range m.ElseBody.Elements {
				a.analyzeNode(elem)
			}
			a.popScope(prev2)
		}
		return bodyType
	}

	st, ok := subjType.(*StructType)
	if !ok || !st.IsVariant {
		a.errorf(m.Pos, "left side of 'match' must be a variant type, got %s", displayType(subjType))
		return Unknown
	}
	found := false
	for _, vc := range st.Variants {
		if vc.Name == m.Constructor.Name {
			found = true
			break
		}
	}
	if !found {
		a.errorfLen(m.Constructor.Pos, len(m.Constructor.Name), "%s is not a constructor of %s", m.Constructor.Name, st.Name)
		return Unknown
	}
	if m.Body == nil {
		return TypBool
	}
	var bodyType Type = TypUnit
	prev := a.pushScope()
	for _, elem := range m.Body.Elements {
		t := a.analyzeNode(elem)
		if _, ok := elem.(ast.Expr); ok {
			bodyType = t
		}
	}
	a.popScope(prev)
	if m.ElseBody != nil {
		prev2 := a.pushScope()
		for _, elem := range m.ElseBody.Elements {
			a.analyzeNode(elem)
		}
		a.popScope(prev2)
	}
	return bodyType
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
			if len(t.TypeArgs) > 0 {
				// Generic application: Box(A) → GenericInstType{Name:"Box", TypeArgs:[GenericType{A}]}
				args := make([]Type, len(t.TypeArgs))
				for i, arg := range t.TypeArgs {
					args[i] = a.resolveTypeExpr(arg)
				}
				return &GenericInstType{Name: t.Name, TypeArgs: args}
			}
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
	case *ast.MemberTypeExpr:
		// Path-dependent type: `g.Node` — look up g's struct type and find its
		// type field named Node. First check the normal scope (for struct fields),
		// then check currentSigParams (for preceding sig params like `g` in
		// `#(g Graph, n g.Node)`).
		var objType Type
		if sym := a.currentScope.Lookup(t.Object); sym != nil {
			objType = sym.Type
		} else if a.currentSigParams != nil {
			if typ, ok := a.currentSigParams[t.Object]; ok {
				objType = typ
			}
		}
		if objType == nil {
			a.errorf(t.Pos, "undefined: %s", t.Object)
			return Unknown
		}
		st, ok := objType.(*StructType)
		if !ok {
			a.errorf(t.Pos, "%s is not a struct type (got %s)", t.Object, objType.typeName())
			return Unknown
		}
		for _, f := range st.Fields {
			if f.Name == t.Member && f.IsTypeField {
				if _, isMeta := f.Type.(*MetaType); isMeta {
					// Abstract type field: return a PathDependentType so call sites
					// can resolve the concrete type from the actual argument.
					return &PathDependentType{Param: t.Object, Member: t.Member}
				}
				return f.Type
			}
		}
		a.errorf(t.Pos, "%s has no type field named %s", t.Object, t.Member)
		return Unknown
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
			if f.Name == methodName || NonWordMethodName(f.Name) == methodName {
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
		if field.IsTypeField {
			continue // compile-time only; not a constructor value parameter
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

// ---------------------------------------------------------------------------
// Phase C — call-site type-variable unification (YZC-0066)
// ---------------------------------------------------------------------------

// unifyTypes builds type-variable bindings by structurally matching formal
// against actual. Bindings for already-bound variables are not overwritten.
func unifyTypes(formal, actual Type, bindings map[string]Type) {
	switch f := formal.(type) {
	case *GenericType:
		if _, already := bindings[f.Name]; !already {
			if _, isGeneric := actual.(*GenericType); !isGeneric {
				bindings[f.Name] = actual
			}
		}
	case *ArrayType:
		if a, ok := actual.(*ArrayType); ok {
			unifyTypes(f.Elem, a.Elem, bindings)
		}
	case *DictType:
		if a, ok := actual.(*DictType); ok {
			unifyTypes(f.Key, a.Key, bindings)
			unifyTypes(f.Val, a.Val, bindings)
		}
	case *OptionType:
		if a, ok := actual.(*OptionType); ok {
			unifyTypes(f.Inner, a.Inner, bindings)
		}
	case *BocType:
		if a, ok := actual.(*BocType); ok {
			for i := range f.Params {
				if i < len(a.Params) {
					unifyTypes(f.Params[i].Type, a.Params[i].Type, bindings)
				}
			}
			for i := range f.Returns {
				if i < len(a.Returns) {
					unifyTypes(f.Returns[i], a.Returns[i], bindings)
				}
			}
		}
	case *GenericInstType:
		if a, ok := actual.(*GenericInstType); ok && a.Name == f.Name {
			for i := range f.TypeArgs {
				if i < len(a.TypeArgs) {
					unifyTypes(f.TypeArgs[i], a.TypeArgs[i], bindings)
				}
			}
		}
	}
}

// substituteType replaces all GenericType occurrences in t using bindings.
func substituteType(t Type, bindings map[string]Type) Type {
	if len(bindings) == 0 {
		return t
	}
	switch tt := t.(type) {
	case *GenericType:
		if bound, ok := bindings[tt.Name]; ok {
			return bound
		}
	case *ArrayType:
		return &ArrayType{Elem: substituteType(tt.Elem, bindings)}
	case *DictType:
		return &DictType{
			Key: substituteType(tt.Key, bindings),
			Val: substituteType(tt.Val, bindings),
		}
	case *OptionType:
		return &OptionType{Inner: substituteType(tt.Inner, bindings)}
	case *BocType:
		params := make([]BocParam, len(tt.Params))
		for i, p := range tt.Params {
			params[i] = BocParam{Label: p.Label, Type: substituteType(p.Type, bindings), HasDefault: p.HasDefault, IsReturn: p.IsReturn}
		}
		returns := make([]Type, len(tt.Returns))
		for i, r := range tt.Returns {
			returns[i] = substituteType(r, bindings)
		}
		return &BocType{Params: params, Returns: returns}
	case *GenericInstType:
		args := make([]Type, len(tt.TypeArgs))
		for i, arg := range tt.TypeArgs {
			args[i] = substituteType(arg, bindings)
		}
		return &GenericInstType{Name: tt.Name, TypeArgs: args}
	case *StructType:
		// A bare generic struct name used as return type (e.g. `Pair` in `makePair #(a K, b V, Pair)`).
		// Substitute TypeParams to produce a concrete GenericInstType.
		if len(tt.TypeParams) > 0 {
			typeArgs := make([]Type, len(tt.TypeParams))
			for i, tp := range tt.TypeParams {
				if bound, ok := bindings[tp]; ok {
					typeArgs[i] = bound
				} else {
					typeArgs[i] = &GenericType{Name: tp}
				}
			}
			return &GenericInstType{Name: tt.Name, TypeArgs: typeArgs}
		}
	case *PathDependentType:
		if bound, ok := bindings[tt.Param+"."+tt.Member]; ok {
			return bound
		}
	}
	return t
}

func isUppercaseName(name string) bool {
	if name == "" {
		return false
	}
	return unicode.IsUpper(rune(name[0]))
}

// registerTypeParam adds the type param name to st.TypeParams and creates an
// IsTypeField entry so member access and constructor matching work correctly.
func (a *Analyzer) registerTypeParam(st *StructType, fieldSet *map[string]bool, name string) {
	st.TypeParams = append(st.TypeParams, name)
	if !(*fieldSet)[name] {
		(*fieldSet)[name] = true
		st.Fields = append(st.Fields, StructField{Name: name, Type: TypMeta, IsTypeField: true})
	}
}

// storeExplicitConstraints resolves the TypeExpr constraints and stores their
// interface names in st.ExplicitConstraints[paramName].
func (a *Analyzer) storeExplicitConstraints(st *StructType, paramName string, constraints []ast.TypeExpr) {
	if len(constraints) == 0 {
		return
	}
	if st.ExplicitConstraints == nil {
		st.ExplicitConstraints = make(map[string][]string)
	}
	for _, c := range constraints {
		if ste, ok := c.(*ast.SimpleTypeExpr); ok {
			st.ExplicitConstraints[paramName] = append(st.ExplicitConstraints[paramName], ste.Name)
		}
	}
}

// storeInlineConstraint synthesises an anonymous interface StructType from an
// inline BocTypeExpr constraint (e.g. V #(describe #(String))). It registers
// the interface at file scope under a generated name (_StructParamConstraint) so
// the lowerer can find it, then stores that name as the explicit constraint.
func (a *Analyzer) storeInlineConstraint(st *StructType, structName, paramName string, inline *ast.BocTypeExpr) {
	syntheticName := "_" + structName + paramName + "Constraint"
	params := a.resolveBocSigParams(inline, false)

	iface := &StructType{Name: syntheticName, IsInterface: true}
	for _, p := range params {
		if !p.IsReturn && p.Label != "" {
			iface.Fields = append(iface.Fields, StructField{Name: p.Label, Type: p.Type})
		}
	}

	// Register at file scope so the lowerer can find it via LookupInFile.
	a.fileScope.Define(&Symbol{Name: syntheticName, Type: iface})

	if st.ExplicitConstraints == nil {
		st.ExplicitConstraints = make(map[string][]string)
	}
	st.ExplicitConstraints[paramName] = append(st.ExplicitConstraints[paramName], syntheticName)
}

// buildAssocTypeBound synthesises an anonymous interface StructType from an
// inline BocTypeExpr constraint on an abstract associated type field
// (e.g. Node #(label #(String))). The interface is registered at file scope
// under a generated name (_StructFieldBound) so the lowerer can emit it, and
// its StructType is stored in StructField.Bound.
func (a *Analyzer) buildAssocTypeBound(structName, fieldName string, sig *ast.BocTypeExpr) *StructType {
	syntheticName := "_" + structName + fieldName + "Bound"
	params := a.resolveBocSigParams(sig, false)
	iface := &StructType{Name: syntheticName, IsInterface: true}
	for _, p := range params {
		if !p.IsReturn && p.Label != "" {
			iface.Fields = append(iface.Fields, StructField{Name: p.Label, Type: p.Type})
		}
	}
	a.fileScope.Define(&Symbol{Name: syntheticName, Type: iface})
	return iface
}

// builtinOpMethods is the set of method names that correspond to Yz built-in
// operators and std-library methods. These are handled by constraintGoSigs in
// the lowerer (builtinConstraintSig) and must NOT be synthesized as anonymous
// interface constraints. Must stay in sync with lower.go:builtinConstraintSig.
var builtinOpMethods = map[string]bool{
	"to_string": true, "to_str": true, "length": true,
	"plus": true, "minus": true, "star": true, "slash": true, "percent": true,
	"lt": true, "gt": true, "lteq": true, "gteq": true,
	"eqeq": true, "neq": true, "ampamp": true, "pipepipe": true,
}

// synthesizeConstraints creates anonymous interface constraints for type params
// whose inferred method-call constraints (activeConstraints) cannot be satisfied
// by any named interface already in scope (YZC-0073). For each such constraint,
// a synthetic interface `_StructTPConstraint` is registered at file scope and
// added to st.ExplicitConstraints so the lowerer emits it as a Go interface.
func (a *Analyzer) synthesizeConstraints(st *StructType, structName string) {
	for tp, constraints := range a.activeConstraints {
		// Skip type params that already have explicit constraints.
		if _, hasExplicit := st.ExplicitConstraints[tp]; hasExplicit {
			continue
		}
		// Collect unique user-defined methods that need synthesis.
		type methodSig struct {
			paramTypes []Type
		}
		methodsToSynth := make(map[string]methodSig)
		for _, c := range constraints {
			if builtinOpMethods[c.MethodName] {
				continue
			}
			if a.FindInterfaceWithMethod(c.MethodName) != "" {
				continue // a named interface covers this method; lowerer handles it
			}
			if _, already := methodsToSynth[c.MethodName]; !already {
				methodsToSynth[c.MethodName] = methodSig{paramTypes: c.ParamTypes}
			}
		}
		if len(methodsToSynth) == 0 {
			continue
		}
		// Build the synthetic interface.
		syntheticName := "_" + structName + tp + "Constraint"
		iface := &StructType{Name: syntheticName, IsInterface: true}
		for methodName, sig := range methodsToSynth {
			params := make([]BocParam, len(sig.paramTypes))
			for i, t := range sig.paramTypes {
				params[i] = BocParam{Type: t}
			}
			iface.Fields = append(iface.Fields, StructField{
				Name: methodName,
				Type: &BocType{Params: params, Returns: []Type{TypUnit}},
			})
		}
		a.fileScope.Define(&Symbol{Name: syntheticName, Type: iface})
		if st.ExplicitConstraints == nil {
			st.ExplicitConstraints = make(map[string][]string)
		}
		st.ExplicitConstraints[tp] = append(st.ExplicitConstraints[tp], syntheticName)
	}
}

// collectGenericNames adds all GenericType names found in t (recursively) to
// names if not already present. Used by Step 5 (implicit T auto-collect).
func collectGenericNames(t Type, names map[string]bool) {
	switch tt := t.(type) {
	case *GenericType:
		names[tt.Name] = true
	case *ArrayType:
		collectGenericNames(tt.Elem, names)
	case *DictType:
		collectGenericNames(tt.Key, names)
		collectGenericNames(tt.Val, names)
	case *OptionType:
		collectGenericNames(tt.Inner, names)
	case *BocType:
		for _, p := range tt.Params {
			collectGenericNames(p.Type, names)
		}
		for _, r := range tt.Returns {
			collectGenericNames(r, names)
		}
	case *GenericInstType:
		for _, arg := range tt.TypeArgs {
			collectGenericNames(arg, names)
		}
	}
}

// analyzeAnnotationBody type-checks the body of an annotation boc.
// Runs in an isolated child scope so declarations don't leak to the enclosing boc.
// String interpolation (${}) is rejected — annotations are compile-time only.
func (a *Analyzer) analyzeAnnotationBody(ann *ast.Annotation) {
	if ann.Body == nil {
		return
	}
	prev := a.pushScope()
	prevInAnn := a.inAnnotation
	a.inAnnotation = true
	defer func() {
		a.inAnnotation = prevInAnn
		a.popScope(prev)
	}()
	for _, elem := range ann.Body.Elements {
		a.analyzeNode(elem)
	}
}
