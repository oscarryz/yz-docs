// Package codegen emits Go source code from an IR file.
package codegen

import (
	"fmt"
	"strings"

	"yz/internal/ir"
)

// Generate converts an ir.File to a Go source string.
func Generate(f *ir.File) string {
	g := &generator{}
	g.emitFile(f)
	return g.sb.String()
}

// ---------------------------------------------------------------------------
// Generator state
// ---------------------------------------------------------------------------

type generator struct {
	sb        strings.Builder
	level     int              // current indentation level (tabs)
	heldCowns map[string]bool // non-nil when inside a ScheduleMulti body
}

func (g *generator) write(s string)                    { g.sb.WriteString(s) }
func (g *generator) writef(f string, a ...any)         { fmt.Fprintf(&g.sb, f, a...) }
func (g *generator) nl()                               { g.sb.WriteByte('\n') }
func (g *generator) ind() string                       { return strings.Repeat("\t", g.level) }
func (g *generator) line(s string)                     { g.write(g.ind()); g.write(s); g.nl() }
func (g *generator) linef(f string, a ...any)          { g.write(g.ind()); g.writef(f, a...); g.nl() }

// sub returns a new generator pre-set to the given indentation level,
// used for emitting multi-line closures embedded in expressions.
// heldCowns is inherited so that closure bodies emitted via sub() see the
// same held-cown context as the enclosing ScheduleMulti body.
func (g *generator) sub(level int) *generator { return &generator{level: level, heldCowns: g.heldCowns} }

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
		g.linef("%s(%s) %s", m.Name, strings.Join(params, ", "), ir.ThunkOrScalar(m.ResultType))
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

	// Build type parameter strings for generic structs.
	// typeConstraints: "[T any]" or "[T interface{...}]" for declarations; typeArgs: "[T]" for references.
	typeConstraints := ""
	typeArgs := ""
	if len(sd.TypeParams) > 0 {
		var constraintParts []string
		for _, tp := range sd.TypeParams {
			if sigs, ok := sd.TypeConstraints[tp]; ok && len(sigs) > 0 {
				constraintParts = append(constraintParts, tp+" interface{ "+strings.Join(sigs, "; ")+" }")
			} else {
				constraintParts = append(constraintParts, tp+" any")
			}
		}
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
		if th.LazyWrap != "" {
			g.writef("%sreturn %s(std.Schedule(%s, func() %s {\n", g.ind(), th.LazyWrap, th.RecvCown, th.ResultType)
		} else {
			g.writef("%sreturn std.Schedule(%s, func() %s {\n", g.ind(), th.RecvCown, th.ResultType)
		}
		g.level++
		g.linef("return %s.%s(%s)", md.RecvName, syncName, strings.Join(names, ", "))
		g.level--
		if th.LazyWrap != "" {
			g.line("}))")
		} else {
			g.line("})")
		}
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
func (g *generator) emitBodyStmts(stmts []ir.Stmt, hasResult bool) {
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
		g.linef("%s", g.expr(st.Expr))
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
			if ds.IsThunk {
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
				if inner, ok := spawnForceInner(sp); ok {
					if mc, ok := inner.(*ir.MethodCall); ok && len(heldCowns) > 0 {
						recvStr := simpleExprStr(mc.Recv)
						if recvStr != "" && heldCowns["&"+recvStr+".Cown"] {
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
							g.linef("%s.Go(func() any {", sp.GroupVar)
							g.linef("\treturn %s.Force()", tv)
							g.linef("})")
							continue
						}
					}
					// Non-held cown — hoist and register goroutine.
					tv := fmt.Sprintf("_st%d", hoistIdx)
					hoistIdx++
					g.linef("%s := %s", tv, g.expr(inner))
					g.linef("%s.Go(func() any {", sp.GroupVar)
					g.linef("\treturn %s.Force()", tv)
					g.linef("})")
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

// wrapLazy applies th.LazyWrap wrapping if set: "std.LazyInt(inner)".
func wrapLazy(lazyWrap, inner string) string {
	if lazyWrap == "" {
		return inner
	}
	return lazyWrap + "(" + inner + ")"
}

// emitThunk generates std.Go(func() T { body }) or std.NewThunk(...).
// When th.RecvCown is non-empty the body is serialized through the singleton's
// cown using std.Schedule. If the body also contains a WaitStmt (BocGroup pattern),
// the split-BocGroup pattern is used to avoid re-entrancy deadlocks: BocGroup
// declarations are hoisted outside the Schedule closure, and BocGroup.Wait() plus
// any subsequent statements run after the cown is released.
// When th.LazyWrap is set, the entire emitted expression is wrapped in that call
// (e.g. "std.LazyInt(std.Schedule(...))") to produce a lazy scalar value.
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
		return wrapLazy(th.LazyWrap, sb.String())
	}

	// Multi-cown: atomically acquire self + extra cowns via ScheduleMulti.
	if len(th.ExtraCowns) > 0 {
		// If the body contains an inline-forced thunkVar (partial-reentrant: the
		// sub-boc needs its own cown plus a held cown), use ScheduleFlatten.
		if thunkFindInlineThunkVar(th.Body) >= 0 {
			return wrapLazy(th.LazyWrap, g.emitScheduleFlatten(th))
		}
		// If the body contains a WaitStmt (from Phase E.1 implicit BocGroup), use
		// the IIFE split-BocGroup pattern so SpawnExprs establish their cown queue
		// positions eagerly while ScheduleMulti holds the cowns.
		if waitIdx := thunkFindWaitIdx(th.Body); waitIdx >= 0 {
			return wrapLazy(th.LazyWrap, g.emitScheduleMultiSplit(th, waitIdx))
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
		return wrapLazy(th.LazyWrap, sb.String())
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
		return wrapLazy(th.LazyWrap, sb.String())
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
	return wrapLazy(th.LazyWrap, sb.String())
}

// spawnForceInner returns the inner thunk expression from a SpawnExpr whose
// body is exactly [ReturnStmt{ForceExpr{expr}}] where expr is not a bare Ident
// (i.e., not already a pre-computed thunk var). Returns (expr, true) on match.
func spawnForceInner(sp *ir.SpawnExpr) (ir.Expr, bool) {
	if len(sp.Body) != 1 {
		return nil, false
	}
	rs, ok := sp.Body[0].(*ir.ReturnStmt)
	if !ok || rs.Value == nil {
		return nil, false
	}
	fe, ok := rs.Value.(*ir.ForceExpr)
	if !ok {
		return nil, false
	}
	if _, isIdent := fe.Thunk.(*ir.Ident); isIdent {
		return nil, false // already a thunk var; no hoisting needed
	}
	return fe.Thunk, true
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

// thunkFindInlineThunkVar returns the index of the first DeclStmt{IsThunk:true,
// IsScalarThunk:false} in body (before any WaitStmt), or -1. Such a declaration
// means the thunk will be forced inline inside the cown closure, which deadlocks
// when the sub-boc needs the same cowns. Use emitScheduleFlatten in that case.
// Scalar thunks (IsScalarThunk:true) are never forced inline and don't deadlock.
func thunkFindInlineThunkVar(body []ir.Stmt) int {
	for i, s := range body {
		if _, ok := s.(*ir.WaitStmt); ok {
			break // past the WaitStmt the cown is already released — safe
		}
		if ds, ok := s.(*ir.DeclStmt); ok && ds.IsThunk && !ds.IsScalarThunk {
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

// emitScheduleFlatten emits the cown-suspension pattern for a multi-cown method
// body that contains an inline-forced thunkVar. The body is split at the first
// thunkVar declaration:
//
//   Phase 1 (inside ScheduleFlatten's protected fn): everything up to and
//   including the thunkVar DeclStmt — sub-boc registered while cowns held.
//
//   Phase 2 (inside returned NewThunk): force the thunkVar (cowns released),
//   then reacquire cowns via ScheduleMulti for the remainder of the body.
//
// If phase2 itself contains another inline thunkVar, the pattern recurses.
func (g *generator) emitScheduleFlatten(th *ir.ThunkExpr) string {
	cowns := append([]string{th.RecvCown}, th.ExtraCowns...)
	cownList := "[]*std.Cown{" + strings.Join(cowns, ", ") + "}"

	splitIdx := thunkFindInlineThunkVar(th.Body)
	phase1 := th.Body[:splitIdx+1] // up to and including the thunkVar DeclStmt
	phase2 := th.Body[splitIdx+1:] // continuation

	// The thunkVar name — used to force it in phase2.
	thunkVarName := th.Body[splitIdx].(*ir.DeclStmt).Name

	var sb strings.Builder
	sb.WriteString("std.ScheduleFlatten(")
	sb.WriteString(cownList)
	sb.WriteString(", func() *std.Thunk[")
	sb.WriteString(th.ResultType)
	sb.WriteString("] {\n")

	p1 := g.sub(g.level + 1)
	p1.emitBodyStmts(phase1, false) // false: don't auto-return last expr
	p1.write(p1.ind())
	p1.write("return std.NewThunk(func() ")
	p1.write(th.ResultType)
	p1.write(" {\n")

	p2 := p1.sub(p1.level + 1)
	// Force the thunkVar — safe here, cowns are released.
	p2.linef("%s := %s.Force()", thunkVarName, thunkVarName)
	// Reacquire cowns for the continuation.
	if len(phase2) > 0 {
		// Remove ForceExpr wrappers around thunkVarName in phase2: after the
		// forcing above, the variable is already the concrete type, not a thunk.
		phase2clean := stripThunkForce(phase2, thunkVarName)
		inner2 := &ir.ThunkExpr{
			ResultType: th.ResultType,
			Body:       phase2clean,
			Spawn:      false,
			RecvCown:   th.RecvCown,
			ExtraCowns: th.ExtraCowns,
		}
		p2.linef("return %s.Force()", p2.expr(inner2))
	}

	p1.write(p2.sb.String())
	p1.write(p1.ind())
	p1.write("})\n")

	sb.WriteString(p1.sb.String())
	sb.WriteString(g.ind())
	sb.WriteString("})")
	return sb.String()
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

// stripThunkForce walks body and removes ForceExpr wrappers around the named
// identifier. Used in emitScheduleFlatten so that phase2 statements (which the
// lowerer emitted with auto-Force on the thunk var) do not double-force a
// variable that has already been forced before being captured by phase2.
func stripThunkForce(body []ir.Stmt, name string) []ir.Stmt {
	result := make([]ir.Stmt, len(body))
	for i, s := range body {
		result[i] = stripForceStmt(s, name)
	}
	return result
}

func stripForceStmt(s ir.Stmt, name string) ir.Stmt {
	switch st := s.(type) {
	case *ir.ReturnStmt:
		return &ir.ReturnStmt{Value: stripForceExprNode(st.Value, name)}
	case *ir.ExprStmt:
		return &ir.ExprStmt{Expr: stripForceExprNode(st.Expr, name)}
	case *ir.DeclStmt:
		return &ir.DeclStmt{Name: st.Name, Type: st.Type, IsThunk: st.IsThunk, Init: stripForceExprNode(st.Init, name)}
	case *ir.AssignStmt:
		return &ir.AssignStmt{Target: stripForceExprNode(st.Target, name), Value: stripForceExprNode(st.Value, name)}
	default:
		return s
	}
}

func stripForceExprNode(e ir.Expr, name string) ir.Expr {
	if e == nil {
		return nil
	}
	switch ex := e.(type) {
	case *ir.ForceExpr:
		if id, ok := ex.Thunk.(*ir.Ident); ok && id.Name == name {
			return id
		}
		return &ir.ForceExpr{Thunk: stripForceExprNode(ex.Thunk, name)}
	case *ir.FieldAccess:
		return &ir.FieldAccess{Object: stripForceExprNode(ex.Object, name), Field: ex.Field}
	case *ir.MethodCall:
		args := make([]ir.Expr, len(ex.Args))
		for i, a := range ex.Args {
			args[i] = stripForceExprNode(a, name)
		}
		return &ir.MethodCall{Recv: stripForceExprNode(ex.Recv, name), Method: ex.Method, Args: args}
	case *ir.FuncCall:
		args := make([]ir.Expr, len(ex.Args))
		for i, a := range ex.Args {
			args[i] = stripForceExprNode(a, name)
		}
		return &ir.FuncCall{Func: stripForceExprNode(ex.Func, name), Args: args}
	case *ir.IndexExpr:
		return &ir.IndexExpr{Object: stripForceExprNode(ex.Object, name), Index: stripForceExprNode(ex.Index, name)}
	default:
		return e
	}
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

// emitSpawn generates groupVar.Go(func() any { body }).
func (g *generator) emitSpawn(s *ir.SpawnExpr) string {
	// Scalar lazy types implement Waitable; use GoWait instead of Go+Force.
	if s.IsScalar && len(s.Body) == 1 {
		if rs, ok := s.Body[0].(*ir.ReturnStmt); ok && rs.Value != nil {
			return s.GroupVar + ".GoWait(" + g.expr(rs.Value) + ")"
		}
	}

	var sb strings.Builder
	sb.WriteString(s.GroupVar)
	sb.WriteString(".Go(func() any {\n")

	inner := g.sub(g.level + 1)
	inner.emitBodyStmts(s.Body, true)
	sb.WriteString(inner.sb.String())

	sb.WriteString(g.ind())
	sb.WriteString("})")
	return sb.String()
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
	// typeConstraints: "[V any]" or "[V interface{...}]" for declarations; typeArgs: "[V]" for references.
	typeConstraints := ""
	typeArgs := ""
	if len(sd.TypeParams) > 0 {
		var constraintParts []string
		for _, tp := range sd.TypeParams {
			if sigs, ok := sd.TypeConstraints[tp]; ok && len(sigs) > 0 {
				constraintParts = append(constraintParts, tp+" interface{ "+strings.Join(sigs, "; ")+" }")
			} else {
				constraintParts = append(constraintParts, tp+" any")
			}
		}
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
}

// emitSwitchStmt emits a Go switch on the discriminant field.
func (g *generator) emitSwitchStmt(sw *ir.SwitchStmt) {
	g.linef("switch %s._variant {", g.expr(sw.Subject))
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
	inner.linef("switch %s._variant {", inner.expr(sw.Subject))
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
	constraints := make([]string, len(tps))
	for i, tp := range tps {
		constraints[i] = tp + " any"
	}
	return "[" + strings.Join(constraints, ", ") + "]"
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
