// Package codegen emits Go source code from an IR file.
package codegen

import (
	"fmt"
	"strings"

	"yz/internal/ir"
)

// Generate converts an ir.File to a Go source string.
func Generate(f *ir.File) string {
	zero := 0
	g := &generator{thunkCount: &zero}
	g.emitFile(f)
	return g.sb.String()
}

// ---------------------------------------------------------------------------
// Generator state
// ---------------------------------------------------------------------------

type generator struct {
	sb         strings.Builder
	level      int              // current indentation level (tabs)
	heldCowns  map[string]bool // non-nil when inside a ScheduleMulti body
	thunkCount *int            // shared counter for unique _thN intermediate thunk vars
}

func (g *generator) write(s string)                    { g.sb.WriteString(s) }
func (g *generator) writef(f string, a ...any)         { fmt.Fprintf(&g.sb, f, a...) }
func (g *generator) nl()                               { g.sb.WriteByte('\n') }
func (g *generator) ind() string                       { return strings.Repeat("\t", g.level) }
func (g *generator) line(s string)                     { g.write(g.ind()); g.write(s); g.nl() }
func (g *generator) linef(f string, a ...any)          { g.write(g.ind()); g.writef(f, a...); g.nl() }

// sub returns a new generator pre-set to the given indentation level,
// used for emitting multi-line closures embedded in expressions.
// heldCowns and thunkCount are inherited so inner closures share the
// same counter and held-cown context as the enclosing body.
func (g *generator) sub(level int) *generator {
	return &generator{level: level, heldCowns: g.heldCowns, thunkCount: g.thunkCount}
}

// freshThunkVar returns a unique _thN name for an intermediate thunk variable.
func (g *generator) freshThunkVar() string {
	n := fmt.Sprintf("_th%d", *g.thunkCount)
	*g.thunkCount++
	return n
}

// ---------------------------------------------------------------------------
// File
// ---------------------------------------------------------------------------

func (g *generator) emitFile(f *ir.File) {
	g.writef("package %s\n\n", f.PkgName)
	g.write("import std \"yz/runtime/rt\"\n")
	if len(f.Imports) > 0 {
		for _, imp := range f.Imports {
			g.writef("import %q\n", imp)
		}
	}
	g.nl()

	for i, decl := range f.Decls {
		if i > 0 {
			g.nl()
		}
		g.emitDecl(decl)
	}
}

// ---------------------------------------------------------------------------
// Declarations
// ---------------------------------------------------------------------------

func (g *generator) emitDecl(d ir.Decl) {
	switch decl := d.(type) {
	case *ir.StructDecl:
		g.emitStructDecl(decl)
	case *ir.InterfaceDecl:
		g.emitInterfaceDecl(decl)
	case *ir.SingletonDecl:
		g.emitSingletonDecl(decl)
	case *ir.FuncDecl:
		g.emitFuncDecl(decl)
	case *ir.TypeAliasDecl:
		g.linef("type %s = %s", decl.Name, decl.Target)
		g.nl()
	}
}

func (g *generator) emitInterfaceDecl(id *ir.InterfaceDecl) {
	g.linef("type %s interface {", id.Name)
	g.level++
	for _, m := range id.Methods {
		var params []string
		for _, p := range m.Params {
			params = append(params, p.Name+" "+p.Type)
		}
		g.linef("%s(%s) *std.Thunk[%s]", m.Name, strings.Join(params, ", "), m.ResultType)
	}
	g.level--
	g.line("}")
	g.nl()
}

func (g *generator) emitStructDecl(sd *ir.StructDecl) {
	// Variant (sum) type: emit discriminant enum + flat struct + per-variant constructors.
	if sd.IsVariant {
		g.emitVariantDecl(sd)
		return
	}

	// Plain result struct for multi-return bocs (YZC-0012): no Cown, no constructor, no String().
	if sd.IsResultType {
		g.linef("type %s struct {", sd.Name)
		g.level++
		for _, f := range sd.Fields {
			g.linef("%s %s", f.Name, f.Type)
		}
		g.level--
		g.line("}")
		g.nl()
		return
	}

	// Build type parameter strings for generic structs.
	// typeConstraints: "[T any]", "[T Talker]", or "[T interface{A;B}]" for declarations;
	// typeArgs: "[T]" for references.
	typeConstraints := ""
	typeArgs := ""
	if len(sd.TypeParams) > 0 {
		constraintParts := buildTypeParamConstraints(sd.TypeParams, sd.ExplicitConstraints, sd.TypeConstraints)
		typeConstraints = "[" + strings.Join(constraintParts, ", ") + "]"
		typeArgs = "[" + strings.Join(sd.TypeParams, ", ") + "]"
	}

	g.linef("type %s%s struct {", sd.Name, typeConstraints)
	g.level++
	g.line("std.Cown")
	for _, f := range sd.Fields {
		g.linef("%s %s", f.Name, f.Type)
	}
	g.level--
	g.line("}")
	g.nl()

	// Constructor (skipped for type-only declarations: Name #(params)).
	if !sd.NoConstructor {
		var params []string
		for _, f := range sd.Fields {
			params = append(params, f.Name+" "+f.Type)
		}
		g.linef("func New%s%s(%s) *%s%s {", sd.Name, typeConstraints, strings.Join(params, ", "), sd.Name, typeArgs)
		g.level++
		g.linef("return &%s%s{", sd.Name, typeArgs)
		g.level++
		for _, f := range sd.Fields {
			g.linef("%s: %s,", f.Name, f.Name)
		}
		g.level--
		g.line("}")
		g.level--
		g.line("}")
	}

	// fmt.Stringer: homoiconic String() method for backtick interpolation and print.
	if !sd.IsVariant {
		g.nl()
		g.linef("func (self *%s%s) String() string {", sd.Name, typeArgs)
		g.level++
		// Build map: type-param name → first field that uses it.
		firstFieldForParam := map[string]string{}
		for _, tp := range sd.TypeParams {
			for _, f := range sd.Fields {
				if f.Type == tp {
					firstFieldForParam[tp] = f.Name
					break
				}
			}
		}
		hasTypeParams := len(sd.TypeParams) > 0 && len(firstFieldForParam) > 0
		if len(sd.Fields) == 0 && !hasTypeParams {
			g.linef("return %q", sd.Name+"()")
		} else if hasTypeParams {
			// Generic: Name(TypeA, TypeB, field: val, ...)
			typeNames := "std.YzTypeName(self." + firstFieldForParam[sd.TypeParams[0]] + ")"
			for _, tp := range sd.TypeParams[1:] {
				if fn, ok := firstFieldForParam[tp]; ok {
					typeNames += " + \", \" + std.YzTypeName(self." + fn + ")"
				}
			}
			if len(sd.Fields) == 0 {
				g.linef("return %q + %s + \")\"", sd.Name+"(", typeNames)
			} else {
				result := fmt.Sprintf("%q", sd.Name+"(") + " + " + typeNames +
					" + \", \" + " + fmt.Sprintf("%q", sd.Fields[0].Name+": ") +
					" + std.StringifyRepr(self." + sd.Fields[0].Name + ")"
				for _, f := range sd.Fields[1:] {
					result += " + " + fmt.Sprintf("%q", ", "+f.Name+": ") +
						" + std.StringifyRepr(self." + f.Name + ")"
				}
				result += " + \")\""
				g.linef("return %s", result)
			}
		} else {
			// Non-generic: Name(field: val, ...)
			result := fmt.Sprintf("%q", sd.Name+"("+sd.Fields[0].Name+": ") +
				" + std.StringifyRepr(self." + sd.Fields[0].Name + ")"
			for _, f := range sd.Fields[1:] {
				result += " + " + fmt.Sprintf("%q", ", "+f.Name+": ") +
					" + std.StringifyRepr(self." + f.Name + ")"
			}
			result += " + \")\""
			g.linef("return %s", result)
		}
		g.level--
		g.line("}")
	}

	for _, m := range sd.Methods {
		g.nl()
		g.emitMethodDecl(m)
	}
}

func (g *generator) emitSingletonDecl(sd *ir.SingletonDecl) {
	// Private struct type. Every singleton boc embeds std.Cown so method
	// bodies can serialize through it via std.Schedule(&self.Cown, ...).
	g.linef("type %s struct {", sd.TypeName)
	g.level++
	g.line("std.Cown")
	for _, f := range sd.Fields {
		g.linef("%s %s", f.Name, f.Type)
	}
	g.level--
	g.line("}")
	g.nl()

	// Homoiconic String() for backtick interpolation.
	g.linef("func (self *%s) String() string {", sd.TypeName)
	g.level++
	if len(sd.Fields) == 0 && len(sd.Methods) == 0 {
		g.linef("return %q", "{ }")
	} else {
		first := true
		result := "\"{ \""
		for _, f := range sd.Fields {
			if !first {
				result += " + \"; \""
			}
			result += " + " + fmt.Sprintf("%q", f.Name+": ") + " + std.StringifyRepr(self."+f.Name+")"
			first = false
		}
		for _, m := range sd.Methods {
			mYz := strings.ToLower(m.Name[:1]) + m.Name[1:]
			if !first {
				result += " + \"; \""
			}
			result += " + " + fmt.Sprintf("%q", mYz+": {}")
			first = false
		}
		result += " + \" }\""
		g.linef("return %s", result)
	}
	g.level--
	g.line("}")
	g.nl()

	// Methods.
	for _, m := range sd.Methods {
		g.emitMethodDecl(m)
		g.nl()
	}

	// Package-level singleton var (only emitted when VarName is non-empty;
	// local boc structs lifted to package level have VarName == "").
	if sd.VarName != "" {
		if len(sd.Fields) > 0 {
			g.writef("var %s = &%s{\n", sd.VarName, sd.TypeName)
			for _, f := range sd.Fields {
				if f.Init != nil {
					g.writef("\t%s: %s,\n", f.Name, g.expr(f.Init))
				}
			}
			g.write("}\n")
		} else {
			g.writef("var %s = &%s{}\n", sd.VarName, sd.TypeName)
		}
	}
}

func (g *generator) emitMethodDecl(md *ir.MethodDecl) {
	params := joinParams(md.Params)
	result := formatResults(md.Results)

	// For single-cown methods (no WaitStmt, no ExtraCowns) emit a lowercase sync
	// body function alongside the exported async method. Callers that already hold
	// the receiver's cown (reentrant context) call the sync version directly instead
	// of going through std.Schedule, which would deadlock or break ordering.
	if th := extractSingleCownThunk(md.Body); th != nil {
		syncName := strings.ToLower(md.Name[:1]) + md.Name[1:]
		g.linef("func (%s %s) %s(%s) %s {", md.RecvName, md.RecvType, syncName, params, th.ResultType)
		g.level++
		g.emitBodyStmts(th.Body, true)
		g.level--
		g.line("}")
		g.nl()
		// Async method delegates to the sync body.
		names := make([]string, len(md.Params))
		for i, p := range md.Params {
			names[i] = p.Name
		}
		g.linef("func (%s %s) %s(%s)%s {", md.RecvName, md.RecvType, md.Name, params, result)
		g.level++
		g.writef("%sreturn std.Schedule(%s, func() %s {\n", g.ind(), th.RecvCown, th.ResultType)
		g.level++
		g.linef("return %s.%s(%s)", md.RecvName, syncName, strings.Join(names, ", "))
		g.level--
		g.line("})")
		g.level--
		g.line("}")
		return
	}

	g.linef("func (%s %s) %s(%s)%s {", md.RecvName, md.RecvType, md.Name, params, result)
	g.level++
	g.emitBodyStmts(md.Body, len(md.Results) > 0)
	g.level--
	g.line("}")
}

func (g *generator) emitFuncDecl(fd *ir.FuncDecl) {
	params := joinParams(fd.Params)
	result := formatResults(fd.Results)
	typeParams := formatTypeParams(fd.TypeParams)
	g.linef("func %s%s(%s)%s {", fd.Name, typeParams, params, result)
	g.level++
	g.emitStmts(fd.Body)
	g.level--
	g.line("}")
}

// ---------------------------------------------------------------------------
// Statements
// ---------------------------------------------------------------------------

// emitBodyStmts emits a method/closure body.
// When hasResult is true, the last ExprStmt is prefixed with "return".
// Declared variables never referenced elsewhere in the body get a `_ = name`
// suppressor so Go's "declared and not used" check does not fire (YZC-0007).
func (g *generator) emitBodyStmts(stmts []ir.Stmt, hasResult bool) {
	used := usedNames(stmts)
	for i, s := range stmts {
		isLast := i == len(stmts)-1
		if hasResult && isLast {
			if es, ok := s.(*ir.ExprStmt); ok {
				g.write(g.ind())
				g.write("return ")
				g.write(g.expr(es.Expr))
				g.nl()
				continue
			}
		}
		g.emitStmt(s)
		if ds, ok := s.(*ir.DeclStmt); ok && !used[ds.Name] {
			g.linef("_ = %s", ds.Name)
		}
	}
}

func (g *generator) emitStmts(stmts []ir.Stmt) {
	for _, s := range stmts {
		g.emitStmt(s)
	}
}

func (g *generator) emitStmt(s ir.Stmt) {
	switch st := s.(type) {
	case *ir.DeclStmt:
		if st.Type == "" {
			g.linef("%s := %s", st.Name, g.expr(st.Init))
		} else if st.Init == nil {
			g.linef("var %s %s", st.Name, st.Type)
		} else {
			g.linef("var %s %s = %s", st.Name, st.Type, g.expr(st.Init))
		}
	case *ir.AssignStmt:
		g.linef("%s = %s", g.expr(st.Target), g.expr(st.Value))
	case *ir.ReturnStmt:
		if st.Value == nil {
			g.line("return std.TheUnit")
		} else {
			g.linef("return %s", g.expr(st.Value))
		}
	case *ir.ExprStmt:
		if sp, ok := st.Expr.(*ir.SpawnExpr); ok {
			g.emitSpawnStmt(sp)
		} else {
			g.linef("%s", g.expr(st.Expr))
		}
	case *ir.IfStmt:
		g.emitIfStmt(st)
	case *ir.WaitStmt:
		g.linef("%s.Wait()", st.GroupVar)
	case *ir.SwitchStmt:
		g.emitSwitchStmt(st)
	}
}

// ---------------------------------------------------------------------------
// Expressions
// ---------------------------------------------------------------------------

// expr returns the Go source for an IR expression.
// Multi-line expressions (ThunkExpr, ClosureExpr) use a sub-generator so
// their inner lines are indented relative to the current level.
func (g *generator) expr(e ir.Expr) string {
	if e == nil {
		return "std.TheUnit"
	}
	switch ex := e.(type) {
	case *ir.IntLit:
		return fmt.Sprintf("std.NewInt(%d)", ex.Val)
	case *ir.DecimalLit:
		return fmt.Sprintf("std.NewDecimal(%g)", ex.Val)
	case *ir.StringLit:
		return fmt.Sprintf("std.NewString(%q)", ex.Val)
	case *ir.BoolLit:
		if ex.Val {
			return "std.NewBool(true)"
		}
		return "std.NewBool(false)"
	case *ir.UnitLit:
		return "std.TheUnit"
	case *ir.Ident:
		return ex.Name
	case *ir.FieldAccess:
		return g.expr(ex.Object) + "." + ex.Field
	case *ir.IndexExpr:
		return g.expr(ex.Object) + ".At(" + g.expr(ex.Index) + ")"
	case *ir.MethodCall:
		return fmt.Sprintf("%s.%s(%s)", g.expr(ex.Recv), ex.Method, g.exprList(ex.Args))
	case *ir.FuncCall:
		if len(ex.TypeArgs) > 0 {
			return fmt.Sprintf("%s[%s](%s)", g.expr(ex.Func), strings.Join(ex.TypeArgs, ", "), g.exprList(ex.Args))
		}
		return fmt.Sprintf("%s(%s)", g.expr(ex.Func), g.exprList(ex.Args))
	case *ir.ThunkExpr:
		return g.emitThunk(ex)
	case *ir.ForceExpr:
		return g.expr(ex.Thunk) + ".Force()"
	case *ir.ClosureExpr:
		return g.emitClosure(ex)
	case *ir.SpawnExpr:
		return g.emitSpawn(ex)
	case *ir.NewGroupExpr:
		return "&std.BocGroup{}"
	case *ir.NewStructExpr:
		return "&" + ex.TypeName + "{}"
	case *ir.MatchExpr:
		return g.emitMatchIIFE(ex)
	case *ir.SwitchExpr:
		return g.emitSwitchIIFE(ex)
	case *ir.VariantTestExpr:
		field := ex.FieldName
		if field == "" {
			field = "_variant"
		}
		return fmt.Sprintf("std.NewBool(%s.%s == %s)", g.expr(ex.Subject), field, ex.ConstName)
	case *ir.StructLitExpr:
		var parts []string
		for _, f := range ex.Fields {
			parts = append(parts, f.Name+": "+g.expr(f.Value))
		}
		return ex.TypeName + "{" + strings.Join(parts, ", ") + "}"
	default:
		return "/* ? */"
	}
}

func (g *generator) exprList(exprs []ir.Expr) string {
	parts := make([]string, len(exprs))
	for i, e := range exprs {
		parts[i] = g.expr(e)
	}
	return strings.Join(parts, ", ")
}

// splitDecl is a variable declaration split across the outer and inner scopes of
// the split-BocGroup pattern: `var name type` goes in the outer scope (visible
// after Schedule/ScheduleMulti) and `name = init` goes inside the closure.
type splitDecl struct {
	name string
	typ  string
	init ir.Expr
}

// partitionWaitBody buckets pre-Wait statements into the three categories used
// by the split-BocGroup emission pattern.
func partitionWaitBody(preWait []ir.Stmt) (hoisted []ir.Stmt, splits []splitDecl, immediate []ir.Stmt) {
	for _, s := range preWait {
		if ds, ok := s.(*ir.DeclStmt); ok {
			if _, isGroup := ds.Init.(*ir.NewGroupExpr); isGroup {
				hoisted = append(hoisted, s)
				continue
			}
			// GoStore target: Type set, Init nil — var declared for SpawnExpr{StoreVar}.
			// Hoist to outer scope so the name is visible after Schedule completes.
			if ds.Type != "" && ds.Init == nil {
				hoisted = append(hoisted, s)
				continue
			}
			// Regular var with explicit type: split so the name is visible after Schedule.
			if ds.Type != "" {
				splits = append(splits, splitDecl{ds.Name, ds.Type, ds.Init})
				continue
			}
		}
		immediate = append(immediate, s)
	}
	return
}

// emitImmediateBody emits body into g inside a ScheduleMulti closure.
// heldCowns is the set of cown addresses held by the enclosing ScheduleMulti.
//
// Two cases for SpawnExpr:
//   - Held cown (receiver's cown is in heldCowns): emit ScheduleAsSuccessor so
//     the sub-boc runs as the immediate successor of the current holder — before
//     any externally-waiting behaviours. Preserves spawn-order happens-before.
//   - Non-held cown: hoist the thunk expression before _bg0.Go so cown
//     registration happens in source order (BOC spawn-order guarantee).
//
// IfStmt branches are handled recursively so both cases apply inside conditionals.
func (g *generator) emitImmediateBody(body []ir.Stmt, heldCowns map[string]bool) {
	hoistIdx := 0
	for _, s := range body {
		switch st := s.(type) {
		case *ir.ExprStmt:
			if sp, ok := st.Expr.(*ir.SpawnExpr); ok {
				thunk := sp.Thunk
				if mc, ok := thunk.(*ir.MethodCall); ok && len(heldCowns) > 0 {
					recvStr := simpleExprStr(mc.Recv)
					// Recursive calls (callee FQN == enclosing BocDecl FQN) go to the
					// tail queue so external callers can interleave between iterations.
					// Only non-recursive calls on held cowns use ScheduleAsSuccessor.
					if recvStr != "" && heldCowns["&"+recvStr+".Cown"] && !mc.IsRecursive {
						// Held cown — schedule as successor to preserve spawn-order
						// happens-before without re-acquiring the cown (already held).
						tv := fmt.Sprintf("_st%d", hoistIdx)
						hoistIdx++
						syncName := strings.ToLower(mc.Method[:1]) + mc.Method[1:]
						g.linef("%s := std.ScheduleAsSuccessor(&%s.Cown, func() std.Unit {", tv, recvStr)
						g.level++
						g.linef("return %s.%s(%s)", recvStr, syncName, g.exprList(mc.Args))
						g.level--
						g.linef("})")
						if sp.StoreVar != "" {
							g.linef("%s.Add(func() { %s = %s.Force() })", sp.GroupVar, sp.StoreVar, tv)
						} else {
							g.linef("%s.Add(func() { %s.Force() })", sp.GroupVar, tv)
						}
						continue
					}
				}
				// Non-held cown (or not a method call) — hoist and register.
				if _, isIdent := thunk.(*ir.Ident); !isIdent {
					tv := fmt.Sprintf("_st%d", hoistIdx)
					hoistIdx++
					g.linef("%s := %s", tv, g.expr(thunk))
					if sp.StoreVar != "" {
						if sp.StoreAnyType != "" {
							g.linef("%s.Add(func() { %s = %s.Force().(%s) })", sp.GroupVar, sp.StoreVar, tv, sp.StoreAnyType)
						} else {
							g.linef("%s.Add(func() { %s = %s.Force() })", sp.GroupVar, sp.StoreVar, tv)
						}
					} else {
						g.linef("%s.Add(func() { %s.Force() })", sp.GroupVar, tv)
					}
					continue
				}
			}
			g.emitStmt(s)
		case *ir.IfStmt:
			g.linef("if %s.GoBool() {", g.expr(st.Cond))
			g.level++
			g.emitImmediateBody(st.Then, heldCowns)
			g.level--
			g.emitImmediateElse(st.Else, heldCowns)
		default:
			g.emitStmt(s)
		}
	}
}

// emitImmediateElse emits the else portion inside an emitImmediateBody context.
func (g *generator) emitImmediateElse(elseStmts []ir.Stmt, heldCowns map[string]bool) {
	if len(elseStmts) == 0 {
		g.line("}")
		return
	}
	if len(elseStmts) == 1 {
		if nested, ok := elseStmts[0].(*ir.IfStmt); ok {
			if bl, ok := nested.Cond.(*ir.BoolLit); ok && bl.Val {
				g.line("} else {")
				g.level++
				g.emitImmediateBody(nested.Then, heldCowns)
				g.level--
				g.line("}")
				return
			}
			g.writef("%s} else if %s.GoBool() {\n", g.ind(), g.expr(nested.Cond))
			g.level++
			g.emitImmediateBody(nested.Then, heldCowns)
			g.level--
			g.emitImmediateElse(nested.Else, heldCowns)
			return
		}
	}
	g.line("} else {")
	g.level++
	g.emitImmediateBody(elseStmts, heldCowns)
	g.level--
	g.line("}")
}

// emitThunk generates std.Go(func() T { body }) or std.NewThunk(...).
// When th.RecvCown is non-empty the body is serialized through the singleton's
// cown using std.Schedule. If the body also contains a WaitStmt (BocGroup pattern),
// the split-BocGroup pattern is used to avoid re-entrancy deadlocks: BocGroup
// declarations are hoisted outside the Schedule closure, and BocGroup.Wait() plus
// any subsequent statements run after the cown is released.
func (g *generator) emitThunk(th *ir.ThunkExpr) string {
	if th.RecvCown == "" {
		fn := "std.Go"
		if !th.Spawn {
			fn = "std.NewThunk"
		}
		var sb strings.Builder
		sb.WriteString(fn)
		sb.WriteString("(func() ")
		sb.WriteString(th.ResultType)
		sb.WriteString(" {\n")
		inner := g.sub(g.level + 1)
		// Inside a ScheduleMulti body (g.heldCowns set), synchronous Unit closures
		// inherit the held-cown context so that boc calls on held cowns emit
		// ScheduleAsSuccessor rather than Schedule. This is the correct behaviour
		// when ? branches are lowered as closure arguments instead of IfStmt nodes.
		if g.heldCowns != nil && !th.Spawn && th.ResultType == "std.Unit" {
			inner.emitImmediateBody(th.Body, g.heldCowns)
			inner.line("return std.TheUnit")
		} else {
			inner.emitBodyStmts(th.Body, true)
		}
		sb.WriteString(inner.sb.String())
		sb.WriteString(g.ind())
		sb.WriteString("})")
		return sb.String()
	}

	// Multi-cown: atomically acquire self + extra cowns via ScheduleMulti.
	if len(th.ExtraCowns) > 0 {
		// If the body contains a WaitStmt (from Phase E.1 implicit BocGroup), use
		// the IIFE split-BocGroup pattern so SpawnExprs establish their cown queue
		// positions eagerly while ScheduleMulti holds the cowns.
		if waitIdx := thunkFindWaitIdx(th.Body); waitIdx >= 0 {
			return g.emitScheduleMultiSplit(th, waitIdx)
		}
		// Clean ScheduleMulti path: no sub-boc goroutines, body accesses fields directly.
		cowns := append([]string{th.RecvCown}, th.ExtraCowns...)
		var sb strings.Builder
		sb.WriteString("std.ScheduleMulti([]*std.Cown{")
		sb.WriteString(strings.Join(cowns, ", "))
		sb.WriteString("}, func() ")
		sb.WriteString(th.ResultType)
		sb.WriteString(" {\n")
		inner := g.sub(g.level + 1)
		inner.emitBodyStmts(th.Body, true)
		sb.WriteString(inner.sb.String())
		sb.WriteString(g.ind())
		sb.WriteString("})")
		return sb.String()
	}

	// Single-cown method body.
	waitIdx := thunkFindWaitIdx(th.Body)

	if waitIdx == -1 {
		// Simple method: no BocGroup — use Schedule directly.
		var sb strings.Builder
		sb.WriteString("std.Schedule(")
		sb.WriteString(th.RecvCown)
		sb.WriteString(", func() ")
		sb.WriteString(th.ResultType)
		sb.WriteString(" {\n")
		inner := g.sub(g.level + 1)
		inner.emitBodyStmts(th.Body, true)
		sb.WriteString(inner.sb.String())
		sb.WriteString(g.ind())
		sb.WriteString("})")
		return sb.String()
	}

	// Split-BocGroup pattern: cown released before waiting for children.
	//
	//   std.NewThunk(func() T {
	//       _bg0 := &std.BocGroup{}      // hoisted BocGroup decl
	//       std.Schedule(cown, func() std.Unit {
	//           _bg0.Go(...)              // SpawnExprs stay inside Schedule
	//           return std.TheUnit
	//       }).Force()                   // releases cown
	//       _bg0.Wait()                  // then wait for children
	//       [post-Wait statements]
	//   })
	preWait := th.Body[:waitIdx]
	postWait := th.Body[waitIdx:] // WaitStmt and everything after

	hoistedDecls, splitDecls, immediateBody := partitionWaitBody(preWait)

	var sb strings.Builder
	sb.WriteString("std.NewThunk(func() ")
	sb.WriteString(th.ResultType)
	sb.WriteString(" {\n")

	outer := g.sub(g.level + 1)

	for _, s := range hoistedDecls {
		outer.emitStmt(s)
	}
	for _, sd := range splitDecls {
		outer.linef("var %s %s", sd.name, sd.typ)
	}

	outer.write(outer.ind())
	outer.write("std.Schedule(")
	outer.write(th.RecvCown)
	outer.write(", func() std.Unit {\n")
	schedInner := outer.sub(outer.level + 1)
	for _, sd := range splitDecls {
		schedInner.linef("%s = %s", sd.name, schedInner.expr(sd.init))
	}
	schedInner.emitImmediateBody(immediateBody, nil)
	schedInner.line("return std.TheUnit")
	outer.write(schedInner.sb.String())
	outer.write(outer.ind())
	outer.write("}).Force()\n")

	outer.emitBodyStmts(postWait, true)

	sb.WriteString(outer.sb.String())
	sb.WriteString(g.ind())
	sb.WriteString("})")
	return sb.String()
}

// thunkFindWaitIdx returns the index of the first WaitStmt in body, or -1.
func thunkFindWaitIdx(body []ir.Stmt) int {
	for i, s := range body {
		if _, ok := s.(*ir.WaitStmt); ok {
			return i
		}
	}
	return -1
}

// extractSingleCownThunk returns the ThunkExpr if body is a single ExprStmt
// wrapping a single-cown ThunkExpr (RecvCown set, ExtraCowns empty, no WaitStmt).
// Returns nil for multi-cown, spawn-only, or split-BocGroup method shapes.
func extractSingleCownThunk(body []ir.Stmt) *ir.ThunkExpr {
	if len(body) != 1 {
		return nil
	}
	es, ok := body[0].(*ir.ExprStmt)
	if !ok {
		return nil
	}
	th, ok := es.Expr.(*ir.ThunkExpr)
	if !ok || th.RecvCown == "" || len(th.ExtraCowns) != 0 {
		return nil
	}
	if thunkFindWaitIdx(th.Body) >= 0 {
		return nil // split-BocGroup pattern — not safe to call inline
	}
	return th
}

// simpleExprStr returns a dotted identifier string for simple receiver expressions
// (Ident or field-access chains), used to match receiver names against heldCowns.
// Returns "" for complex expressions.
func simpleExprStr(e ir.Expr) string {
	switch ex := e.(type) {
	case *ir.Ident:
		return ex.Name
	case *ir.FieldAccess:
		s := simpleExprStr(ex.Object)
		if s == "" {
			return ""
		}
		return s + "." + ex.Field
	default:
		return ""
	}
}

// emitScheduleMultiSplit handles multi-cown method bodies that contain an
// implicit BocGroup (from Phase E.1). ScheduleMulti is called EAGERLY (outside
// NewThunk) so that cown queue positions are established immediately when Call()
// is invoked, preserving BOC spawn-order determinism. SpawnExprs in the body
// hoist their inner thunk expressions so sub-boc scheduling also happens eagerly
// while the cowns are still held. NewThunk wraps only the waiting phase.
//
//	func() *std.Thunk[T] {
//	    _bg0 := &std.BocGroup{}          // hoisted
//	    _sched := std.ScheduleMulti(...)  // EAGER: registered at Call() time
//	    return std.NewThunk(func() T {
//	        _sched.Force()               // waits for ScheduleMulti + releases cowns
//	        _bg0.Wait()                  // then wait for goroutines
//	        [post-Wait stmts]
//	    })
//	}()
func (g *generator) emitScheduleMultiSplit(th *ir.ThunkExpr, waitIdx int) string {
	preWait := th.Body[:waitIdx]
	postWait := th.Body[waitIdx:] // WaitStmt and everything after

	hoistedDecls, splitDecls, immediateBody := partitionWaitBody(preWait)

	cowns := append([]string{th.RecvCown}, th.ExtraCowns...)
	thunkType := "*std.Thunk[" + th.ResultType + "]"

	// Build the held-cown set for reentrant inline detection.
	heldCowns := map[string]bool{th.RecvCown: true}
	for _, c := range th.ExtraCowns {
		heldCowns[c] = true
		// "&self.foo.Cown" → also register the canonical form "&self.foo.Cown"
		// that emitImmediateBody will look up via simpleExprStr.
		if strings.HasPrefix(c, "&") && strings.HasSuffix(c, ".Cown") {
			inner := c[1 : len(c)-5]
			if !strings.Contains(inner, ".") {
				heldCowns["&self."+inner+".Cown"] = true
			}
		}
	}

	var sb strings.Builder
	sb.WriteString("func() ")
	sb.WriteString(thunkType)
	sb.WriteString(" {\n")

	outer := g.sub(g.level + 1)

	for _, s := range hoistedDecls {
		outer.emitStmt(s)
	}
	for _, sd := range splitDecls {
		outer.linef("var %s %s", sd.name, sd.typ)
	}

	outer.write(outer.ind())
	outer.write("_sched := std.ScheduleMulti([]*std.Cown{")
	outer.write(strings.Join(cowns, ", "))
	outer.write("}, func() std.Unit {\n")
	schedInner := outer.sub(outer.level + 1)
	schedInner.heldCowns = heldCowns // propagate into closure arguments emitted within this body
	for _, sd := range splitDecls {
		schedInner.linef("%s = %s", sd.name, schedInner.expr(sd.init))
	}
	schedInner.emitImmediateBody(immediateBody, heldCowns)
	schedInner.line("return std.TheUnit")
	outer.write(schedInner.sb.String())
	outer.write(outer.ind())
	outer.write("})\n")

	outer.write(outer.ind())
	outer.write("return std.NewThunk(func() ")
	outer.write(th.ResultType)
	outer.write(" {\n")
	innerThunk := outer.sub(outer.level + 1)
	innerThunk.line("_sched.Force()")
	innerThunk.emitBodyStmts(postWait, true)
	outer.write(innerThunk.sb.String())
	outer.write(outer.ind())
	outer.write("})\n")

	sb.WriteString(outer.sb.String())
	sb.WriteString(g.ind())
	sb.WriteString("}()")
	return sb.String()
}

// emitClosure generates func(params) T { body }.
func (g *generator) emitClosure(c *ir.ClosureExpr) string {
	params := joinParams(c.Params)
	hasResult := c.ResultType != "" && c.ResultType != "std.Unit"

	var sb strings.Builder
	sb.WriteString("func(")
	sb.WriteString(params)
	sb.WriteString(") ")
	sb.WriteString(c.ResultType)
	sb.WriteString(" {\n")

	inner := g.sub(g.level + 1)
	inner.emitBodyStmts(c.Body, hasResult)
	sb.WriteString(inner.sb.String())

	sb.WriteString(g.ind())
	sb.WriteString("}")
	return sb.String()
}

// emitSpawnStmt emits a SpawnExpr as one or two statements:
//   _thN := Thunk
//   GroupVar.Add(func() { StoreVar = _thN.Force() })   // StoreVar case
//   GroupVar.Add(func() { _thN.Force() })               // no-StoreVar case
// The intermediate _thN variable ensures the thunk (and its goroutine) is
// created at registration time, so all goroutines run concurrently before Wait.
func (g *generator) emitSpawnStmt(s *ir.SpawnExpr) {
	tv := g.freshThunkVar()
	g.linef("%s := %s", tv, g.expr(s.Thunk))
	if s.StoreVar != "" {
		if s.StoreAnyType != "" {
			g.linef("%s.Add(func() { %s = %s.Force().(%s) })", s.GroupVar, s.StoreVar, tv, s.StoreAnyType)
		} else {
			g.linef("%s.Add(func() { %s = %s.Force() })", s.GroupVar, s.StoreVar, tv)
		}
	} else {
		g.linef("%s.Add(func() { %s.Force() })", s.GroupVar, tv)
	}
}

// emitSpawn is the expression-position fallback for SpawnExpr (used in
// emitImmediateBody hoisting). Emits the thunk expression only; callers
// must emit the Add line separately.
func (g *generator) emitSpawn(s *ir.SpawnExpr) string {
	return g.expr(s.Thunk)
}

// emitIfStmt emits an if/else-if/else chain. When the else body is a single
// nested IfStmt it is rendered as "} else if ..." for readability.
func (g *generator) emitIfStmt(st *ir.IfStmt) {
	g.linef("if %s.GoBool() {", g.expr(st.Cond))
	g.level++
	g.emitStmts(st.Then)
	g.level--
	g.emitElse(st.Else)
}

// emitElse emits the else portion of an if statement.
// If the else is a single IfStmt, it chains as "} else if ...".
func (g *generator) emitElse(elseStmts []ir.Stmt) {
	if len(elseStmts) == 0 {
		g.line("}")
		return
	}
	if len(elseStmts) == 1 {
		if nested, ok := elseStmts[0].(*ir.IfStmt); ok {
			// BoolLit{true} cond means "always true" default arm — emit as plain else.
			if bl, ok := nested.Cond.(*ir.BoolLit); ok && bl.Val {
				g.line("} else {")
				g.level++
				g.emitStmts(nested.Then)
				g.level--
				g.line("}")
				return
			}
			g.writef("%s} else if %s.GoBool() {\n", g.ind(), g.expr(nested.Cond))
			g.level++
			g.emitStmts(nested.Then)
			g.level--
			g.emitElse(nested.Else)
			return
		}
	}
	g.line("} else {")
	g.level++
	g.emitStmts(elseStmts)
	g.level--
	g.line("}")
}

// emitMatchIIFE emits a match expression as an immediately-invoked closure:
//
//	func() ResultType {
//	    if arm[0].Cond.GoBool() { arm[0].Body... }
//	    else if arm[1].Cond.GoBool() { arm[1].Body... }
//	    else { arm[n].Body... }
//	    return std.TheUnit // if no default
//	}()
func (g *generator) emitMatchIIFE(me *ir.MatchExpr) string {
	var sb strings.Builder
	sb.WriteString("func() ")
	sb.WriteString(me.ResultType)
	sb.WriteString(" {\n")

	inner := g.sub(g.level + 1)
	inner.emitMatchArms(me.Arms)

	// If no default arm, emit a fallthrough return.
	hasDefault := false
	for _, arm := range me.Arms {
		if arm.Cond == nil {
			hasDefault = true
			break
		}
	}
	if !hasDefault {
		inner.linef("return %s", zeroValueOf(me.ResultType))
	}

	sb.WriteString(inner.sb.String())
	sb.WriteString(g.ind())
	sb.WriteString("}()")
	return sb.String()
}

// zeroValueOf returns a Go zero value expression for the given type string.
func zeroValueOf(typ string) string {
	switch typ {
	case "std.Int":
		return "std.NewInt(0)"
	case "std.Decimal":
		return "std.NewDecimal(0)"
	case "std.String":
		return `std.NewString("")`
	case "std.Bool":
		return "std.NewBool(false)"
	default:
		return "std.TheUnit"
	}
}

// emitMatchArms emits the if/else if/else chain for match arms.
func (g *generator) emitMatchArms(arms []*ir.MatchArm) {
	first := true
	for _, arm := range arms {
		if arm.Cond == nil {
			// Default arm — emit as else { body }
			g.line("} else {")
		} else if first {
			g.linef("if %s.GoBool() {", g.expr(arm.Cond))
			first = false
		} else {
			g.linef("} else if %s.GoBool() {", g.expr(arm.Cond))
		}
		g.level++
		g.emitStmts(arm.Body)
		g.level--
	}
	if len(arms) > 0 {
		g.line("}")
	}
}

// ---------------------------------------------------------------------------
// Variant (sum) type emission
// ---------------------------------------------------------------------------

// emitVariantDecl emits a Go sum type: discriminant enum, flat struct, and
// per-variant constructors. For `Pet: { Cat(name String), Dog(name String) }`:
//
//	type _PetVariant int
//	const ( _PetCat _PetVariant = iota; _PetDog )
//	type Pet struct { _variant _PetVariant; name std.String; ... }
//	func NewPetCat(name std.String) *Pet { return &Pet{_variant: _PetCat, ...} }
func (g *generator) emitVariantDecl(sd *ir.StructDecl) {
	discType := "_" + sd.Name + "Variant"

	// Build type parameter strings for generic variants.
	// typeConstraints: "[V any]", "[V Talker]", or "[V interface{A;B}]" for declarations;
	// typeArgs: "[V]" for references.
	typeConstraints := ""
	typeArgs := ""
	if len(sd.TypeParams) > 0 {
		constraintParts := buildTypeParamConstraints(sd.TypeParams, sd.ExplicitConstraints, sd.TypeConstraints)
		typeConstraints = "[" + strings.Join(constraintParts, ", ") + "]"
		typeArgs = "[" + strings.Join(sd.TypeParams, ", ") + "]"
	}

	// Discriminant enum type (no type params — it's just an int alias).
	g.linef("type %s int", discType)
	g.nl()

	// Constants block.
	g.line("const (")
	g.level++
	for i, vc := range sd.Variants {
		constName := "_" + sd.Name + vc.Name
		if i == 0 {
			g.linef("%s %s = iota", constName, discType)
		} else {
			g.linef("%s", constName)
		}
	}
	g.level--
	g.line(")")
	g.nl()

	// Flat struct (with optional type params).
	g.linef("type %s%s struct {", sd.Name, typeConstraints)
	g.level++
	g.linef("_variant %s", discType)
	for _, f := range sd.Fields {
		g.linef("%s %s", f.Name, f.Type)
	}
	g.level--
	g.line("}")
	g.nl()

	// Per-variant constructors (with optional type params).
	for _, vc := range sd.Variants {
		constName := "_" + sd.Name + vc.Name
		var params []string
		for _, f := range vc.Fields {
			params = append(params, f.Name+" "+f.Type)
		}
		g.linef("func New%s%s%s(%s) *%s%s {", sd.Name, vc.Name, typeConstraints, strings.Join(params, ", "), sd.Name, typeArgs)
		g.level++
		g.linef("return &%s%s{", sd.Name, typeArgs)
		g.level++
		g.linef("_variant: %s,", constName)
		for _, f := range vc.Fields {
			g.linef("%s: %s,", f.Name, f.Name)
		}
		g.level--
		g.line("}")
		g.level--
		g.line("}")
		g.nl()
	}

	// Homoiconic String() for backtick interpolation.
	g.linef("func (self *%s%s) String() string {", sd.Name, typeArgs)
	g.level++
	g.line("switch self._variant {")
	for _, vc := range sd.Variants {
		constName := "_" + sd.Name + vc.Name
		g.linef("case %s:", constName)
		g.level++
		if len(vc.Fields) == 0 {
			g.linef("return %q", sd.Name+"."+vc.Name+"()")
		} else {
			result := fmt.Sprintf("%q", sd.Name+"."+vc.Name+"("+vc.Fields[0].Name+": ") +
				" + std.StringifyRepr(self." + vc.Fields[0].Name + ")"
			for _, f := range vc.Fields[1:] {
				result += " + " + fmt.Sprintf("%q", ", "+f.Name+": ") + " + std.StringifyRepr(self." + f.Name + ")"
			}
			result += " + \")\""
			g.linef("return %s", result)
		}
		g.level--
	}
	g.line("}")
	g.linef("return %q", sd.Name+"(?)")
	g.level--
	g.line("}")
}

// emitSwitchStmt emits a Go switch on the discriminant field.
func (g *generator) emitSwitchStmt(sw *ir.SwitchStmt) {
	field := sw.FieldName
	if field == "" {
		field = "_variant"
	}
	g.linef("switch %s.%s {", g.expr(sw.Subject), field)
	for _, c := range sw.Cases {
		g.linef("case %s:", c.ConstName)
		g.level++
		g.emitStmts(c.Body)
		g.level--
	}
	g.line("}")
}

// emitSwitchIIFE emits a discriminant match as an immediately-invoked closure.
func (g *generator) emitSwitchIIFE(sw *ir.SwitchExpr) string {
	var sb strings.Builder
	sb.WriteString("func() ")
	sb.WriteString(sw.ResultType)
	sb.WriteString(" {\n")

	inner := g.sub(g.level + 1)
	swField := sw.FieldName
	if swField == "" {
		swField = "_variant"
	}
	inner.linef("switch %s.%s {", inner.expr(sw.Subject), swField)
	for _, c := range sw.Cases {
		inner.linef("case %s:", c.ConstName)
		inner.level++
		inner.emitBodyStmts(c.Body, true)
		inner.level--
	}
	inner.line("}")
	inner.linef("return %s", zeroValueOf(sw.ResultType))

	sb.WriteString(inner.sb.String())
	sb.WriteString(g.ind())
	sb.WriteString("}()")
	return sb.String()
}

// ---------------------------------------------------------------------------
// Unused-variable analysis (YZC-0007)
// ---------------------------------------------------------------------------

// usedNames returns the set of variable names that are read within stmts.
// Plain-Ident assignment targets (writes) are not counted as reads.
// SpawnExpr.GroupVar, SpawnExpr.StoreVar, and WaitStmt.GroupVar are string
// references that also count as reads.
func usedNames(stmts []ir.Stmt) map[string]bool {
	seen := map[string]bool{}
	for _, s := range stmts {
		collectUsedStmt(s, seen)
	}
	return seen
}

func collectUsedStmt(s ir.Stmt, seen map[string]bool) {
	switch st := s.(type) {
	case *ir.DeclStmt:
		collectUsedExpr(st.Init, seen)
	case *ir.AssignStmt:
		// Ident target is a pure write — do not mark as used.
		// FieldAccess/IndexExpr objects are reads of their base variable.
		switch tgt := st.Target.(type) {
		case *ir.FieldAccess:
			collectUsedExpr(tgt.Object, seen)
		case *ir.IndexExpr:
			collectUsedExpr(tgt.Object, seen)
			collectUsedExpr(tgt.Index, seen)
		}
		collectUsedExpr(st.Value, seen)
	case *ir.ReturnStmt:
		collectUsedExpr(st.Value, seen)
	case *ir.ExprStmt:
		collectUsedExpr(st.Expr, seen)
	case *ir.IfStmt:
		collectUsedExpr(st.Cond, seen)
		for _, sub := range st.Then {
			collectUsedStmt(sub, seen)
		}
		for _, sub := range st.Else {
			collectUsedStmt(sub, seen)
		}
	case *ir.WaitStmt:
		seen[st.GroupVar] = true
	case *ir.SwitchStmt:
		collectUsedExpr(st.Subject, seen)
		for _, c := range st.Cases {
			for _, sub := range c.Body {
				collectUsedStmt(sub, seen)
			}
		}
	}
}

func collectUsedExpr(e ir.Expr, seen map[string]bool) {
	if e == nil {
		return
	}
	switch ex := e.(type) {
	case *ir.Ident:
		seen[ex.Name] = true
	case *ir.MethodCall:
		collectUsedExpr(ex.Recv, seen)
		for _, a := range ex.Args {
			collectUsedExpr(a, seen)
		}
	case *ir.FuncCall:
		collectUsedExpr(ex.Func, seen)
		for _, a := range ex.Args {
			collectUsedExpr(a, seen)
		}
	case *ir.FieldAccess:
		collectUsedExpr(ex.Object, seen)
	case *ir.IndexExpr:
		collectUsedExpr(ex.Object, seen)
		collectUsedExpr(ex.Index, seen)
	case *ir.ThunkExpr:
		for _, s := range ex.Body {
			collectUsedStmt(s, seen)
		}
	case *ir.ForceExpr:
		collectUsedExpr(ex.Thunk, seen)
	case *ir.ClosureExpr:
		for _, s := range ex.Body {
			collectUsedStmt(s, seen)
		}
	case *ir.SpawnExpr:
		seen[ex.GroupVar] = true
		if ex.StoreVar != "" {
			seen[ex.StoreVar] = true
		}
		collectUsedExpr(ex.Thunk, seen)
	case *ir.MatchExpr:
		for _, arm := range ex.Arms {
			collectUsedExpr(arm.Cond, seen)
			for _, s := range arm.Body {
				collectUsedStmt(s, seen)
			}
		}
	case *ir.SwitchExpr:
		collectUsedExpr(ex.Subject, seen)
		for _, c := range ex.Cases {
			for _, s := range c.Body {
				collectUsedStmt(s, seen)
			}
		}
	case *ir.VariantTestExpr:
		collectUsedExpr(ex.Subject, seen)
	}
}

// ---------------------------------------------------------------------------
// Formatting helpers
// ---------------------------------------------------------------------------

func joinParams(params []*ir.ParamSpec) string {
	parts := make([]string, len(params))
	for i, p := range params {
		parts[i] = p.Name + " " + p.Type
	}
	return strings.Join(parts, ", ")
}

func formatTypeParams(tps []string) string {
	if len(tps) == 0 {
		return ""
	}
	parts := make([]string, len(tps))
	for i, tp := range tps {
		parts[i] = tp + " any"
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

// buildTypeParamConstraints returns the per-param constraint string for a
// Go generic parameter list. It checks explicit constraints first (source-declared
// interface names like "Talker"), falling back to inferred method-signature
// constraints, and finally to "any".
//
//   - 0 explicit, 0 inferred → "T any"
//   - 1 explicit              → "T Talker"
//   - 2+ explicit             → "T interface{ Talker; Serializable }"
//   - 0 explicit, n inferred  → "T interface{ MethodSig; ... }" (existing behaviour)
func buildTypeParamConstraints(typeParams []string, explicit, inferred map[string][]string) []string {
	parts := make([]string, len(typeParams))
	for i, tp := range typeParams {
		if names, ok := explicit[tp]; ok && len(names) > 0 {
			switch len(names) {
			case 1:
				parts[i] = tp + " " + names[0]
			default:
				parts[i] = tp + " interface{ " + strings.Join(names, "; ") + " }"
			}
		} else if sigs, ok := inferred[tp]; ok && len(sigs) > 0 {
			parts[i] = tp + " interface{ " + strings.Join(sigs, "; ") + " }"
		} else {
			parts[i] = tp + " any"
		}
	}
	return parts
}

func formatResults(results []string) string {
	switch len(results) {
	case 0:
		return ""
	case 1:
		return " " + results[0]
	default:
		return " (" + strings.Join(results, ", ") + ")"
	}
}
