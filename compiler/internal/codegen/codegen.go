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
	for _, f := range sd.Fields {
		if f.Embedded {
			g.linef("%s", f.Name) // Go embedding: just the type name
		} else {
			g.linef("%s %s", f.Name, f.Type)
		}
	}
	g.level--
	g.line("}")
	g.nl()

	// Constructor (skipped for type-only declarations: Name #(params)).
	if !sd.NoConstructor {
		// Build param list: expand embedded sub-fields inline.
		var params []string
		for _, f := range sd.Fields {
			if f.Embedded {
				for _, sf := range f.EmbeddedFields {
					params = append(params, sf.Name+" "+sf.Type)
				}
			} else {
				params = append(params, f.Name+" "+f.Type)
			}
		}
		g.linef("func New%s%s(%s) *%s%s {", sd.Name, typeConstraints, strings.Join(params, ", "), sd.Name, typeArgs)
		g.level++
		g.linef("return &%s%s{", sd.Name, typeArgs)
		g.level++
		for _, f := range sd.Fields {
			if f.Embedded {
				var subArgs []string
				for _, sf := range f.EmbeddedFields {
					subArgs = append(subArgs, sf.Name)
				}
				g.linef("%s: *New%s(%s),", f.Name, f.Name, strings.Join(subArgs, ", "))
			} else {
				g.linef("%s: %s,", f.Name, f.Name)
			}
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
	case *ir.ForStmt:
		g.linef("for %s.GoBool() {", g.expr(st.Cond))
		g.level++
		g.emitStmts(st.Body)
		g.level--
		g.line("}")
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
