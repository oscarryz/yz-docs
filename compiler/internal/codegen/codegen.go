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
	sb    strings.Builder
	level int // current indentation level (tabs)
}

func (g *generator) write(s string)                    { g.sb.WriteString(s) }
func (g *generator) writef(f string, a ...any)         { fmt.Fprintf(&g.sb, f, a...) }
func (g *generator) nl()                               { g.sb.WriteByte('\n') }
func (g *generator) ind() string                       { return strings.Repeat("\t", g.level) }
func (g *generator) line(s string)                     { g.write(g.ind()); g.write(s); g.nl() }
func (g *generator) linef(f string, a ...any)          { g.write(g.ind()); g.writef(f, a...); g.nl() }

// sub returns a new generator pre-set to the given indentation level,
// used for emitting multi-line closures embedded in expressions.
func (g *generator) sub(level int) *generator { return &generator{level: level} }

// ---------------------------------------------------------------------------
// File
// ---------------------------------------------------------------------------

func (g *generator) emitFile(f *ir.File) {
	g.writef("package %s\n\n", f.PkgName)
	g.write("import std \"yz/runtime/yzrt\"\n")
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
	case *ir.SingletonDecl:
		g.emitSingletonDecl(decl)
	case *ir.FuncDecl:
		g.emitFuncDecl(decl)
	}
}

func (g *generator) emitStructDecl(sd *ir.StructDecl) {
	g.linef("type %s struct {", sd.Name)
	g.level++
	for _, f := range sd.Fields {
		g.linef("%s %s", f.Name, f.Type)
	}
	g.level--
	g.line("}")
	g.nl()

	// Constructor (only when there are fields).
	if len(sd.Fields) > 0 {
		var params []string
		for _, f := range sd.Fields {
			params = append(params, f.Name+" "+f.Type)
		}
		g.linef("func New%s(%s) *%s {", sd.Name, strings.Join(params, ", "), sd.Name)
		g.level++
		g.linef("return &%s{", sd.Name)
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
	// Private struct type.
	g.linef("type %s struct {", sd.TypeName)
	g.level++
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

	// Package-level singleton var.
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

func (g *generator) emitMethodDecl(md *ir.MethodDecl) {
	params := joinParams(md.Params)
	result := formatResults(md.Results)
	g.linef("func (%s %s) %s(%s)%s {", md.RecvName, md.RecvType, md.Name, params, result)
	g.level++
	g.emitBodyStmts(md.Body, len(md.Results) > 0)
	g.level--
	g.line("}")
}

func (g *generator) emitFuncDecl(fd *ir.FuncDecl) {
	params := joinParams(fd.Params)
	result := formatResults(fd.Results)
	g.linef("func %s(%s)%s {", fd.Name, params, result)
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
	case *ir.ForStmt:
		g.linef("for %s.GoBool() {", g.expr(st.Cond))
		g.level++
		g.emitStmts(st.Body)
		g.level--
		g.line("}")
	case *ir.IfStmt:
		g.linef("if %s.GoBool() {", g.expr(st.Cond))
		g.level++
		g.emitStmts(st.Then)
		g.level--
		if len(st.Else) > 0 {
			g.line("} else {")
			g.level++
			g.emitStmts(st.Else)
			g.level--
		}
		g.line("}")
	case *ir.WaitStmt:
		g.linef("%s.Wait()", st.GroupVar)
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

// emitThunk generates std.Go(func() T { body }) or std.NewThunk(...).
func (g *generator) emitThunk(th *ir.ThunkExpr) string {
	fn := "std.Go"
	if !th.Spawn {
		fn = "std.NewThunk"
	}
	var sb strings.Builder
	sb.WriteString(fn)
	sb.WriteString("(func() ")
	sb.WriteString(th.ResultType)
	sb.WriteString(" {\n")

	// Inner body at level+1.
	inner := g.sub(g.level + 1)
	inner.emitBodyStmts(th.Body, true)
	sb.WriteString(inner.sb.String())

	sb.WriteString(g.ind()) // closing brace at current indent
	sb.WriteString("})")
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

// emitSpawn generates groupVar.Go(func() any { body }).
func (g *generator) emitSpawn(s *ir.SpawnExpr) string {
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
