package parser

import (
	"testing"

	"yz/internal/ast"
	"yz/internal/token"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func parse(t *testing.T, src string) *ast.SourceFile {
	t.Helper()
	p := New([]byte(src))
	sf, err := p.ParseFile()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	return sf
}

func parseErr(t *testing.T, src string) error {
	t.Helper()
	p := New([]byte(src))
	_, err := p.ParseFile()
	return err
}

// stmt returns the n-th top-level node, asserting it is a Stmt or Expr.
func stmt(t *testing.T, sf *ast.SourceFile, n int) ast.Node {
	t.Helper()
	if n >= len(sf.Stmts) {
		t.Fatalf("want stmt[%d], only %d stmts", n, len(sf.Stmts))
	}
	return sf.Stmts[n]
}

func asShortDecl(t *testing.T, n ast.Node) *ast.ShortDecl {
	t.Helper()
	d, ok := n.(*ast.ShortDecl)
	if !ok {
		t.Fatalf("expected *ast.ShortDecl, got %T", n)
	}
	return d
}

func asTypedDecl(t *testing.T, n ast.Node) *ast.TypedDecl {
	t.Helper()
	d, ok := n.(*ast.TypedDecl)
	if !ok {
		t.Fatalf("expected *ast.TypedDecl, got %T", n)
	}
	return d
}

func asAssignment(t *testing.T, n ast.Node) *ast.Assignment {
	t.Helper()
	a, ok := n.(*ast.Assignment)
	if !ok {
		t.Fatalf("expected *ast.Assignment, got %T", n)
	}
	return a
}

func asBinaryExpr(t *testing.T, e ast.Expr) *ast.BinaryExpr {
	t.Helper()
	b, ok := e.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", e)
	}
	return b
}

func asCallExpr(t *testing.T, n ast.Node) *ast.CallExpr {
	t.Helper()
	var e ast.Expr
	switch v := n.(type) {
	case ast.Expr:
		e = v
	default:
		t.Fatalf("expected expr node, got %T", n)
	}
	c, ok := e.(*ast.CallExpr)
	if !ok {
		t.Fatalf("expected *ast.CallExpr, got %T", e)
	}
	return c
}

func asIdent(t *testing.T, e ast.Expr) *ast.Ident {
	t.Helper()
	id, ok := e.(*ast.Ident)
	if !ok {
		t.Fatalf("expected *ast.Ident, got %T", e)
	}
	return id
}

func asIntLit(t *testing.T, e ast.Expr) *ast.IntLit {
	t.Helper()
	l, ok := e.(*ast.IntLit)
	if !ok {
		t.Fatalf("expected *ast.IntLit, got %T", e)
	}
	return l
}

func asStringLit(t *testing.T, e ast.Expr) *ast.StringLit {
	t.Helper()
	l, ok := e.(*ast.StringLit)
	if !ok {
		t.Fatalf("expected *ast.StringLit, got %T", e)
	}
	return l
}

func asBocLiteral(t *testing.T, e ast.Expr) *ast.BocLiteral {
	t.Helper()
	b, ok := e.(*ast.BocLiteral)
	if !ok {
		t.Fatalf("expected *ast.BocLiteral, got %T", e)
	}
	return b
}

// ---------------------------------------------------------------------------
// 01 — Literals
// ---------------------------------------------------------------------------

func TestParseIntLiteral(t *testing.T) {
	sf := parse(t, "42")
	n := stmt(t, sf, 0)
	l := asIntLit(t, n.(ast.Expr))
	if l.Value != "42" {
		t.Errorf("got %q, want %q", l.Value, "42")
	}
}

func TestParseDecimalLiteral(t *testing.T) {
	sf := parse(t, "3.14")
	n := stmt(t, sf, 0)
	l, ok := n.(*ast.DecimalLit)
	if !ok {
		t.Fatalf("expected *ast.DecimalLit, got %T", n)
	}
	if l.Value != "3.14" {
		t.Errorf("got %q, want %q", l.Value, "3.14")
	}
}

func TestParseStringLiteral(t *testing.T) {
	sf := parse(t, `"hello"`)
	n := stmt(t, sf, 0)
	l := asStringLit(t, n.(ast.Expr))
	if l.Value != `"hello"` {
		t.Errorf("got %q, want %q", l.Value, `"hello"`)
	}
}

// ---------------------------------------------------------------------------
// 02 — Declarations
// ---------------------------------------------------------------------------

func TestParseShortDecl(t *testing.T) {
	sf := parse(t, `name: "Alice"`)
	d := asShortDecl(t, stmt(t, sf, 0))
	if len(d.Names) != 1 || d.Names[0].Name != "name" {
		t.Errorf("unexpected names: %v", d.Names)
	}
	if len(d.Values) != 1 {
		t.Fatalf("expected 1 value, got %d", len(d.Values))
	}
	asStringLit(t, d.Values[0])
}

func TestParseTypedDecl(t *testing.T) {
	sf := parse(t, "age Int")
	d := asTypedDecl(t, stmt(t, sf, 0))
	if d.Name.Name != "age" {
		t.Errorf("got name %q, want %q", d.Name.Name, "age")
	}
	st, ok := d.Type.(*ast.SimpleTypeExpr)
	if !ok {
		t.Fatalf("expected SimpleTypeExpr, got %T", d.Type)
	}
	if st.Name != "Int" {
		t.Errorf("got type %q, want Int", st.Name)
	}
	if d.Value != nil {
		t.Errorf("expected no initializer")
	}
}

func TestParseTypedDeclWithInit(t *testing.T) {
	sf := parse(t, "age Int = 30")
	d := asTypedDecl(t, stmt(t, sf, 0))
	if d.Value == nil {
		t.Fatal("expected initializer")
	}
	asIntLit(t, d.Value)
}

func TestParseMultipleShortDecl(t *testing.T) {
	// a, b: swap("x", "y")
	sf := parse(t, `a, b: swap("x", "y")`)
	d := asShortDecl(t, stmt(t, sf, 0))
	if len(d.Names) != 2 {
		t.Fatalf("expected 2 names, got %d", len(d.Names))
	}
	if d.Names[0].Name != "a" || d.Names[1].Name != "b" {
		t.Errorf("names: %v", d.Names)
	}
}

// ---------------------------------------------------------------------------
// 03 — Assignment
// ---------------------------------------------------------------------------

func TestParseAssignment(t *testing.T) {
	sf := parse(t, `name = "Bob"`)
	a := asAssignment(t, stmt(t, sf, 0))
	id := asIdent(t, a.Target)
	if id.Name != "name" {
		t.Errorf("got %q, want name", id.Name)
	}
}

func TestParseMultiAssignment(t *testing.T) {
	sf := parse(t, `a, b = swap("x", "y")`)
	a := asAssignment(t, stmt(t, sf, 0))
	if len(a.Names) != 2 {
		t.Fatalf("expected 2 names, got %d", len(a.Names))
	}
	if a.Names[0].Name != "a" || a.Names[1].Name != "b" {
		t.Errorf("names: %v", a.Names)
	}
}

func TestParseMemberAssignment(t *testing.T) {
	sf := parse(t, "person.name = \"Bob\"")
	a := asAssignment(t, stmt(t, sf, 0))
	m, ok := a.Target.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected MemberExpr target, got %T", a.Target)
	}
	if m.Member.Name != "name" {
		t.Errorf("got member %q, want name", m.Member.Name)
	}
}

// ---------------------------------------------------------------------------
// 04 — Expressions: binary, unary, grouping
// ---------------------------------------------------------------------------

func TestParseBinaryExprLeftToRight(t *testing.T) {
	// 1 + 2 * 3 should parse as (1 + 2) * 3 (left-to-right, no precedence)
	sf := parse(t, "1 + 2 * 3")
	n := stmt(t, sf, 0)
	outer := asBinaryExpr(t, n.(ast.Expr))
	if outer.Op != "*" {
		t.Errorf("outer op: got %q, want *", outer.Op)
	}
	inner := asBinaryExpr(t, outer.Left)
	if inner.Op != "+" {
		t.Errorf("inner op: got %q, want +", inner.Op)
	}
	asIntLit(t, inner.Left)  // 1
	asIntLit(t, inner.Right) // 2
	asIntLit(t, outer.Right) // 3
}

func TestParseGroupedExpr(t *testing.T) {
	// 1 + (2 * 3): parens override left-to-right
	sf := parse(t, "1 + (2 * 3)")
	n := stmt(t, sf, 0)
	outer := asBinaryExpr(t, n.(ast.Expr))
	if outer.Op != "+" {
		t.Errorf("outer op: got %q, want +", outer.Op)
	}
	// Right should be a GroupExpr containing 2 * 3
	grp, ok := outer.Right.(*ast.GroupExpr)
	if !ok {
		t.Fatalf("expected GroupExpr on right, got %T", outer.Right)
	}
	inner := asBinaryExpr(t, grp.Expr)
	if inner.Op != "*" {
		t.Errorf("inner op: got %q, want *", inner.Op)
	}
}

func TestParseUnaryNeg(t *testing.T) {
	sf := parse(t, "-42")
	n := stmt(t, sf, 0)
	u, ok := n.(*ast.UnaryExpr)
	if !ok {
		t.Fatalf("expected UnaryExpr, got %T", n)
	}
	if u.Op != "-" {
		t.Errorf("got op %q, want -", u.Op)
	}
	asIntLit(t, u.Operand)
}

// ---------------------------------------------------------------------------
// 05 — Call expressions
// ---------------------------------------------------------------------------

func TestParseCallNoArgs(t *testing.T) {
	sf := parse(t, "greet()")
	c := asCallExpr(t, stmt(t, sf, 0))
	id := asIdent(t, c.Callee)
	if id.Name != "greet" {
		t.Errorf("callee: got %q, want greet", id.Name)
	}
	if len(c.Args) != 0 {
		t.Errorf("expected no args, got %d", len(c.Args))
	}
}

func TestParseCallPositionalArg(t *testing.T) {
	sf := parse(t, `greet("Alice")`)
	c := asCallExpr(t, stmt(t, sf, 0))
	if len(c.Args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(c.Args))
	}
	if c.Args[0].Label != "" {
		t.Errorf("expected positional arg, got label %q", c.Args[0].Label)
	}
	asStringLit(t, c.Args[0].Value)
}

func TestParseCallNamedArg(t *testing.T) {
	sf := parse(t, `greet(name: "Alice")`)
	c := asCallExpr(t, stmt(t, sf, 0))
	if len(c.Args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(c.Args))
	}
	if c.Args[0].Label != "name" {
		t.Errorf("expected label %q, got %q", "name", c.Args[0].Label)
	}
}

func TestParseCallMultiArg(t *testing.T) {
	sf := parse(t, `add(1, 2)`)
	c := asCallExpr(t, stmt(t, sf, 0))
	if len(c.Args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(c.Args))
	}
}

func TestParseMemberCall(t *testing.T) {
	sf := parse(t, "person.greet()")
	c := asCallExpr(t, stmt(t, sf, 0))
	m, ok := c.Callee.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected MemberExpr callee, got %T", c.Callee)
	}
	if m.Member.Name != "greet" {
		t.Errorf("method: got %q, want greet", m.Member.Name)
	}
}

func TestParseChainedCalls(t *testing.T) {
	// 1.to(10).each({})
	sf := parse(t, "1.to(10).each({})")
	// Should be: CallExpr{ Callee: MemberExpr{ Object: CallExpr{ ... } } }
	outer := asCallExpr(t, stmt(t, sf, 0))
	m, ok := outer.Callee.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("outer callee: expected MemberExpr, got %T", outer.Callee)
	}
	if m.Member.Name != "each" {
		t.Errorf("outer method: got %q, want each", m.Member.Name)
	}
	inner := asCallExpr(t, m.Object)
	im, ok := inner.Callee.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("inner callee: expected MemberExpr, got %T", inner.Callee)
	}
	if im.Member.Name != "to" {
		t.Errorf("inner method: got %q, want to", im.Member.Name)
	}
}

// ---------------------------------------------------------------------------
// 06 — Boc literals
// ---------------------------------------------------------------------------

func TestParseEmptyBoc(t *testing.T) {
	sf := parse(t, "{}")
	n := stmt(t, sf, 0)
	b := asBocLiteral(t, n.(ast.Expr))
	if len(b.Elements) != 0 {
		t.Errorf("expected empty boc, got %d elements", len(b.Elements))
	}
}

func TestParseBocWithStatements(t *testing.T) {
	sf := parse(t, `{
    name: "Alice"
    age: 30
}`)
	n := stmt(t, sf, 0)
	b := asBocLiteral(t, n.(ast.Expr))
	if len(b.Elements) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(b.Elements))
	}
	asShortDecl(t, b.Elements[0])
	asShortDecl(t, b.Elements[1])
}

func TestParseShortDeclWithBoc(t *testing.T) {
	sf := parse(t, `counter: {
    count: 0
}`)
	d := asShortDecl(t, stmt(t, sf, 0))
	if d.Names[0].Name != "counter" {
		t.Errorf("got %q, want counter", d.Names[0].Name)
	}
	asBocLiteral(t, d.Values[0])
}

// ---------------------------------------------------------------------------
// 07 — Boc with signature
// ---------------------------------------------------------------------------

func TestParseBocWithSigOnly(t *testing.T) {
	// greet #(name String)  — signature only, no body
	sf := parse(t, "greet #(name String)")
	n := stmt(t, sf, 0)
	bws, ok := n.(*ast.BocWithSig)
	if !ok {
		t.Fatalf("expected *ast.BocWithSig, got %T", n)
	}
	if bws.Name.Name != "greet" {
		t.Errorf("name: got %q, want greet", bws.Name.Name)
	}
	if len(bws.Sig.Params) != 1 {
		t.Fatalf("expected 1 param, got %d", len(bws.Sig.Params))
	}
	if bws.Sig.Params[0].Label != "name" {
		t.Errorf("param label: got %q, want name", bws.Sig.Params[0].Label)
	}
	if bws.Body != nil {
		t.Errorf("expected no body")
	}
}

func TestParseBocWithSigAndBody(t *testing.T) {
	// greet #(name String) { name }
	// Body does NOT need to redeclare name.
	sf := parse(t, "greet #(name String) { name }")
	n := stmt(t, sf, 0)
	bws, ok := n.(*ast.BocWithSig)
	if !ok {
		t.Fatalf("expected *ast.BocWithSig, got %T", n)
	}
	if bws.Body == nil {
		t.Fatal("expected body")
	}
	if len(bws.Body.Elements) != 1 {
		t.Fatalf("expected 1 body element, got %d", len(bws.Body.Elements))
	}
	id := asIdent(t, bws.Body.Elements[0].(ast.Expr))
	if id.Name != "name" {
		t.Errorf("body element: got %q, want name", id.Name)
	}
}

func TestParseBocWithSigAssignForm(t *testing.T) {
	// greet #(name String) = { name String }
	sf := parse(t, "greet #(name String) = { name String }")
	n := stmt(t, sf, 0)
	bws, ok := n.(*ast.BocWithSig)
	if !ok {
		t.Fatalf("expected *ast.BocWithSig, got %T", n)
	}
	if bws.Body == nil {
		t.Fatal("expected body")
	}
}

// ---------------------------------------------------------------------------
// 08 — Type declarations
// ---------------------------------------------------------------------------

func TestParseTypeDecl(t *testing.T) {
	sf := parse(t, `Person: {
    name String
    age Int
}`)
	d := asShortDecl(t, stmt(t, sf, 0))
	if d.Names[0].Name != "Person" {
		t.Errorf("got %q, want Person", d.Names[0].Name)
	}
	if d.Names[0].TokType != token.TYPE_IDENT {
		t.Errorf("expected TYPE_IDENT, got %v", d.Names[0].TokType)
	}
}

func TestParseTypeDeclWithVariants(t *testing.T) {
	sf := parse(t, `Result: {
    T, E,
    Ok(value T),
    Err(error E)
}`)
	d := asShortDecl(t, stmt(t, sf, 0))
	if d.Names[0].Name != "Result" {
		t.Errorf("got %q, want Result", d.Names[0].Name)
	}
}

// ---------------------------------------------------------------------------
// 09 — Arrays and Dicts
// ---------------------------------------------------------------------------

func TestParseArrayLiteral(t *testing.T) {
	sf := parse(t, "[1, 2, 3]")
	n := stmt(t, sf, 0)
	arr, ok := n.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("expected *ast.ArrayLiteral, got %T", n)
	}
	if len(arr.Elements) != 3 {
		t.Errorf("expected 3 elements, got %d", len(arr.Elements))
	}
}

func TestParseEmptyArrayLiteral(t *testing.T) {
	sf := parse(t, "[Int]()")
	n := stmt(t, sf, 0)
	arr, ok := n.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("expected *ast.ArrayLiteral, got %T", n)
	}
	if arr.ElemType == nil {
		t.Errorf("expected ElemType for empty array")
	}
	if len(arr.Elements) != 0 {
		t.Errorf("expected no elements for empty array")
	}
}

func TestParseDictLiteral(t *testing.T) {
	sf := parse(t, `["Alice": 30, "Bob": 25]`)
	n := stmt(t, sf, 0)
	d, ok := n.(*ast.DictLiteral)
	if !ok {
		t.Fatalf("expected *ast.DictLiteral, got %T", n)
	}
	if len(d.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(d.Entries))
	}
}

func TestParseEmptyDictLiteral(t *testing.T) {
	sf := parse(t, "[String:Int]()")
	n := stmt(t, sf, 0)
	d, ok := n.(*ast.DictLiteral)
	if !ok {
		t.Fatalf("expected *ast.DictLiteral, got %T", n)
	}
	if d.KeyType == nil || d.ValType == nil {
		t.Errorf("expected KeyType and ValType for empty dict")
	}
}

// ---------------------------------------------------------------------------
// 10 — Match expressions
// ---------------------------------------------------------------------------

func TestParseConditionMatch(t *testing.T) {
	sf := parse(t, `match {
    score >= 90 => "A"
}, {
    "B"
}`)
	n := stmt(t, sf, 0)
	m, ok := n.(*ast.MatchExpr)
	if !ok {
		t.Fatalf("expected *ast.MatchExpr, got %T", n)
	}
	if m.Subject != nil {
		t.Errorf("expected no subject for condition match")
	}
	if len(m.Arms) != 2 {
		t.Fatalf("expected 2 arms, got %d", len(m.Arms))
	}
	if m.Arms[0].Condition == nil {
		t.Errorf("first arm should have condition")
	}
	if m.Arms[1].Condition != nil {
		t.Errorf("second arm should have no condition (default)")
	}
}

func TestParseVariantMatch(t *testing.T) {
	sf := parse(t, `match response {
    Success => print("ok")
}, {
    Failure => print("err")
}`)
	n := stmt(t, sf, 0)
	m, ok := n.(*ast.MatchExpr)
	if !ok {
		t.Fatalf("expected *ast.MatchExpr, got %T", n)
	}
	if m.Subject == nil {
		t.Errorf("expected subject for variant match")
	}
	if len(m.Arms) != 2 {
		t.Fatalf("expected 2 arms, got %d", len(m.Arms))
	}
}

// ---------------------------------------------------------------------------
// 11 — Keyword statements
// ---------------------------------------------------------------------------

func TestParseReturn(t *testing.T) {
	sf := parse(t, "return 42")
	n := stmt(t, sf, 0)
	r, ok := n.(*ast.ReturnStmt)
	if !ok {
		t.Fatalf("expected *ast.ReturnStmt, got %T", n)
	}
	asIntLit(t, r.Value)
}

func TestParseReturnEmpty(t *testing.T) {
	sf := parse(t, "return\n")
	n := stmt(t, sf, 0)
	r, ok := n.(*ast.ReturnStmt)
	if !ok {
		t.Fatalf("expected *ast.ReturnStmt, got %T", n)
	}
	if r.Value != nil {
		t.Errorf("expected nil return value for bare return")
	}
}

func TestParseBreak(t *testing.T) {
	sf := parse(t, "break")
	n := stmt(t, sf, 0)
	if _, ok := n.(*ast.BreakStmt); !ok {
		t.Fatalf("expected *ast.BreakStmt, got %T", n)
	}
}

func TestParseContinue(t *testing.T) {
	sf := parse(t, "continue")
	n := stmt(t, sf, 0)
	if _, ok := n.(*ast.ContinueStmt); !ok {
		t.Fatalf("expected *ast.ContinueStmt, got %T", n)
	}
}

func TestParseMix(t *testing.T) {
	sf := parse(t, "mix Named")
	n := stmt(t, sf, 0)
	m, ok := n.(*ast.MixStmt)
	if !ok {
		t.Fatalf("expected *ast.MixStmt, got %T", n)
	}
	if m.Name.Name != "Named" {
		t.Errorf("mix name: got %q, want Named", m.Name.Name)
	}
}

// ---------------------------------------------------------------------------
// 12 — Index access
// ---------------------------------------------------------------------------

func TestParseIndexExpr(t *testing.T) {
	sf := parse(t, "array[0]")
	n := stmt(t, sf, 0)
	idx, ok := n.(*ast.IndexExpr)
	if !ok {
		t.Fatalf("expected *ast.IndexExpr, got %T", n)
	}
	asIdent(t, idx.Object)
	asIntLit(t, idx.Index)
}

// ---------------------------------------------------------------------------
// 13 — Conditional (Bool.?)
// ---------------------------------------------------------------------------

func TestParseConditional(t *testing.T) {
	// x == 0 ? { "zero" }, { "nonzero" }
	// Parsed as:  (x == 0) ? { "zero" }, { "nonzero" }
	// Which is:  BinaryExpr{ Left: BinaryExpr{x, ==, 0}, Op: ?, Right: { "zero" } }
	// But the second boc is passed as a separate argument in the invocation...
	// Actually per grammar: Expression = UnaryExpr { non_word UnaryExpr }
	// So: (x == 0) is left, then ? is non_word, then { "zero" } is next UnaryExpr
	// But { "zero" }, { "nonzero" } — the comma is not a non_word...
	// The ? method takes two boc args via comma inside the call. But in Yz syntax
	// `flag ? { a }, { b }` is actually `flag.?({a}, {b})` conceptually.
	// Per grammar, the comma terminates the expression and the second boc is
	// the next top-level statement... but that doesn't make semantic sense.
	// The resolution: after parsing `flag ? { a }`, if we see `,` followed by
	// a `{`, we treat it as additional arguments to the binary non-word invocation.
	// This is the special comma-continuation rule for multi-arg non-word methods.
	sf := parse(t, `x == 0 ? { "zero" }, { "nonzero" }`)
	n := stmt(t, sf, 0)
	// Outer: BinaryExpr with op "?"
	outer := asBinaryExpr(t, n.(ast.Expr))
	if outer.Op != "?" {
		t.Fatalf("outer op: got %q, want ?", outer.Op)
	}
	// Left: x == 0
	inner := asBinaryExpr(t, outer.Left)
	if inner.Op != "==" {
		t.Errorf("inner op: got %q, want ==", inner.Op)
	}
}

// ---------------------------------------------------------------------------
// 14 — Multi-line programs
// ---------------------------------------------------------------------------

func TestParseMultiLineProgram(t *testing.T) {
	src := `name: "Alice"
age: 30
print(name)`
	sf := parse(t, src)
	if len(sf.Stmts) != 3 {
		t.Fatalf("expected 3 stmts, got %d", len(sf.Stmts))
	}
	asShortDecl(t, sf.Stmts[0])
	asShortDecl(t, sf.Stmts[1])
	asCallExpr(t, sf.Stmts[2])
}

// ---------------------------------------------------------------------------
// 15 — Info strings
// ---------------------------------------------------------------------------

func TestParseInfoString(t *testing.T) {
	// Info string immediately before a declaration
	sf := parse(t, `"A counter"
counter: 0`)
	// The info string becomes part of the AST; the short decl follows it.
	if len(sf.Stmts) < 1 {
		t.Fatal("expected at least 1 stmt")
	}
	// The short decl should have the info string attached.
	d := asShortDecl(t, sf.Stmts[0])
	// Info string is accessible via the... actually info strings are separate
	// nodes OR attached to the next decl. We'll verify there are 2 nodes OR
	// the decl has an InfoString attached — either design is fine.
	// For now: verify the decl exists.
	if d.Names[0].Name != "counter" {
		t.Errorf("got %q, want counter", d.Names[0].Name)
	}
}

// ---------------------------------------------------------------------------
// 16 — Realistic programs
// ---------------------------------------------------------------------------

func TestParseCounter(t *testing.T) {
	src := `counter: {
    count: 0
    increment: { count = count + 1 }
    value: { count }
}`
	sf := parse(t, src)
	if len(sf.Stmts) != 1 {
		t.Fatalf("expected 1 top-level stmt, got %d", len(sf.Stmts))
	}
	d := asShortDecl(t, sf.Stmts[0])
	if d.Names[0].Name != "counter" {
		t.Errorf("got %q, want counter", d.Names[0].Name)
	}
	body := asBocLiteral(t, d.Values[0])
	if len(body.Elements) != 3 {
		t.Fatalf("expected 3 body elements, got %d: %v", len(body.Elements), body.Elements)
	}
}

func TestParseConditionalExpression(t *testing.T) {
	src := `grade: match {
    score >= 90 => "A"
}, {
    score >= 80 => "B"
}, {
    "C"
}`
	sf := parse(t, src)
	d := asShortDecl(t, stmt(t, sf, 0))
	if d.Names[0].Name != "grade" {
		t.Errorf("got %q, want grade", d.Names[0].Name)
	}
	m, ok := d.Values[0].(*ast.MatchExpr)
	if !ok {
		t.Fatalf("expected MatchExpr value, got %T", d.Values[0])
	}
	if len(m.Arms) != 3 {
		t.Fatalf("expected 3 arms, got %d", len(m.Arms))
	}
}
