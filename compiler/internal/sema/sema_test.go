package sema

import (
	"strings"
	"testing"

	"yz/internal/lexer"
	"yz/internal/parser"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func analyzeSource(t *testing.T, src string) (*Analyzer, error) {
	t.Helper()
	toks := lexer.Tokenize([]byte(src))
	_ = toks
	p := parser.New([]byte(src))
	sf, err := p.ParseFile()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	a := NewAnalyzer()
	err = a.AnalyzeFile(sf)
	return a, err
}

func mustAnalyze(t *testing.T, src string) *Analyzer {
	t.Helper()
	a, err := analyzeSource(t, src)
	if err != nil {
		t.Fatalf("unexpected analysis error: %v", err)
	}
	return a
}

func expectError(t *testing.T, src string, contains string) {
	t.Helper()
	_, err := analyzeSource(t, src)
	if err == nil {
		t.Fatalf("expected error containing %q, got nil", contains)
	}
	if !strings.Contains(err.Error(), contains) {
		t.Errorf("error %q does not contain %q", err.Error(), contains)
	}
}

// ---------------------------------------------------------------------------
// 01 — Literal types
// ---------------------------------------------------------------------------

func TestIntLiteralType(t *testing.T) {
	a := mustAnalyze(t, "42")
	// The expression '42' in the source should have type Int.
	typ := a.ExprType(a.LastExpr())
	if typ != TypInt {
		t.Errorf("got %v, want Int", typ)
	}
}

func TestDecimalLiteralType(t *testing.T) {
	a := mustAnalyze(t, "3.14")
	typ := a.ExprType(a.LastExpr())
	if typ != TypDecimal {
		t.Errorf("got %v, want Decimal", typ)
	}
}

func TestStringLiteralType(t *testing.T) {
	a := mustAnalyze(t, `"hello"`)
	typ := a.ExprType(a.LastExpr())
	if typ != TypString {
		t.Errorf("got %v, want String", typ)
	}
}

// ---------------------------------------------------------------------------
// 02 — Declarations and identifier resolution
// ---------------------------------------------------------------------------

func TestShortDeclInference(t *testing.T) {
	a := mustAnalyze(t, `x: 42`)
	sym := a.LookupInFile("x")
	if sym == nil {
		t.Fatal("symbol 'x' not found")
	}
	if sym.Type != TypInt {
		t.Errorf("x: got type %v, want Int", sym.Type)
	}
}

func TestShortDeclString(t *testing.T) {
	a := mustAnalyze(t, `name: "Alice"`)
	sym := a.LookupInFile("name")
	if sym == nil {
		t.Fatal("symbol 'name' not found")
	}
	if sym.Type != TypString {
		t.Errorf("name: got type %v, want String", sym.Type)
	}
}

func TestTypedDecl(t *testing.T) {
	a := mustAnalyze(t, `age Int`)
	sym := a.LookupInFile("age")
	if sym == nil {
		t.Fatal("symbol 'age' not found")
	}
	if sym.Type != TypInt {
		t.Errorf("age: got type %v, want Int", sym.Type)
	}
}

func TestTypedDeclWithInit(t *testing.T) {
	a := mustAnalyze(t, `age Int = 30`)
	sym := a.LookupInFile("age")
	if sym == nil {
		t.Fatal("symbol 'age' not found")
	}
	if sym.Type != TypInt {
		t.Errorf("age: got type %v, want Int", sym.Type)
	}
}

func TestIdentResolves(t *testing.T) {
	a := mustAnalyze(t, `x: 42
y: x`)
	sym := a.LookupInFile("y")
	if sym == nil {
		t.Fatal("symbol 'y' not found")
	}
	if sym.Type != TypInt {
		t.Errorf("y: got type %v, want Int", sym.Type)
	}
}

func TestUnknownIdentError(t *testing.T) {
	expectError(t, `y`, "undefined")
}

func TestMultipleDecls(t *testing.T) {
	a := mustAnalyze(t, `name: "Alice"
age: 30
active: true`)
	for _, tc := range []struct {
		name string
		want Type
	}{
		{"name", TypString},
		{"age", TypInt},
		{"active", TypBool},
	} {
		sym := a.LookupInFile(tc.name)
		if sym == nil {
			t.Fatalf("symbol %q not found", tc.name)
		}
		if sym.Type != tc.want {
			t.Errorf("%s: got %v, want %v", tc.name, sym.Type, tc.want)
		}
	}
}

// ---------------------------------------------------------------------------
// 03 — Arithmetic and binary expressions
// ---------------------------------------------------------------------------

func TestArithmeticType(t *testing.T) {
	a := mustAnalyze(t, `1 + 2`)
	typ := a.ExprType(a.LastExpr())
	if typ != TypInt {
		t.Errorf("1 + 2: got %v, want Int", typ)
	}
}

func TestComparisonType(t *testing.T) {
	a := mustAnalyze(t, `1 == 2`)
	typ := a.ExprType(a.LastExpr())
	if typ != TypBool {
		t.Errorf("1 == 2: got %v, want Bool", typ)
	}
}

func TestUnaryNegType(t *testing.T) {
	a := mustAnalyze(t, `-42`)
	typ := a.ExprType(a.LastExpr())
	if typ != TypInt {
		t.Errorf("-42: got %v, want Int", typ)
	}
}

// ---------------------------------------------------------------------------
// 04 — Boc types and return inference
// ---------------------------------------------------------------------------

func TestBocReturnInference(t *testing.T) {
	// counter: { count: 0; count }
	// The boc's return type should be Int (last expr is 'count' which is Int).
	a := mustAnalyze(t, `counter: {
    count: 0
    count
}`)
	sym := a.LookupInFile("counter")
	if sym == nil {
		t.Fatal("symbol 'counter' not found")
	}
	bt, ok := sym.Type.(*BocType)
	if !ok {
		t.Fatalf("counter: got type %T, want *BocType", sym.Type)
	}
	if len(bt.Returns) != 1 || bt.Returns[0] != TypInt {
		t.Errorf("counter return: got %v, want [Int]", bt.Returns)
	}
}

func TestBocParams(t *testing.T) {
	// add: { a Int; b Int; a + b }
	a := mustAnalyze(t, `add: {
    a Int
    b Int
    a + b
}`)
	sym := a.LookupInFile("add")
	if sym == nil {
		t.Fatal("symbol 'add' not found")
	}
	bt, ok := sym.Type.(*BocType)
	if !ok {
		t.Fatalf("add: got type %T, want *BocType", sym.Type)
	}
	if len(bt.Params) != 2 {
		t.Fatalf("add: got %d params, want 2", len(bt.Params))
	}
	if bt.Params[0].Label != "a" || bt.Params[0].Type != TypInt {
		t.Errorf("param[0]: got %v, want a Int", bt.Params[0])
	}
	if bt.Params[1].Label != "b" || bt.Params[1].Type != TypInt {
		t.Errorf("param[1]: got %v, want b Int", bt.Params[1])
	}
	if len(bt.Returns) != 1 || bt.Returns[0] != TypInt {
		t.Errorf("add return: got %v, want [Int]", bt.Returns)
	}
}

func TestNestedBocScoping(t *testing.T) {
	// Inner boc can read outer variable.
	a := mustAnalyze(t, `counter: {
    count: 0
    increment: { count = count + 1 }
}`)
	sym := a.LookupInFile("counter")
	if sym == nil {
		t.Fatal("'counter' not found")
	}
	if _, ok := sym.Type.(*BocType); !ok {
		t.Fatalf("counter: got %T, want *BocType", sym.Type)
	}
}

// ---------------------------------------------------------------------------
// 05 — BocWithSig special form
// ---------------------------------------------------------------------------

func TestBocWithSigParamsInScope(t *testing.T) {
	// greet #(name String) { name }
	// 'name' should be in scope in the body without re-declaration.
	a := mustAnalyze(t, `greet #(name String) { name }`)
	sym := a.LookupInFile("greet")
	if sym == nil {
		t.Fatal("'greet' not found")
	}
	bt, ok := sym.Type.(*BocType)
	if !ok {
		t.Fatalf("greet: got %T, want *BocType", sym.Type)
	}
	if len(bt.Params) != 1 || bt.Params[0].Label != "name" {
		t.Errorf("greet params: got %v", bt.Params)
	}
	// Return type should be String (last expr is 'name' which is String)
	if len(bt.Returns) != 1 || bt.Returns[0] != TypString {
		t.Errorf("greet return: got %v, want [String]", bt.Returns)
	}
}

func TestBocWithSigReturnType(t *testing.T) {
	// greet #(name String, String) — return type is explicit
	// (the anonymous String at end of param list is the return type)
	a := mustAnalyze(t, `greet #(name String, String) { name }`)
	sym := a.LookupInFile("greet")
	if sym == nil {
		t.Fatal("'greet' not found")
	}
	bt, ok := sym.Type.(*BocType)
	if !ok {
		t.Fatalf("greet: got %T, want *BocType", sym.Type)
	}
	if len(bt.Returns) != 1 || bt.Returns[0] != TypString {
		t.Errorf("greet return: got %v, want [String]", bt.Returns)
	}
}

// ---------------------------------------------------------------------------
// 06 — Structural types (user-defined)
// ---------------------------------------------------------------------------

func TestUserDefinedType(t *testing.T) {
	a := mustAnalyze(t, `Person: {
    name String
    age Int
}`)
	sym := a.LookupInFile("Person")
	if sym == nil {
		t.Fatal("'Person' not found")
	}
	st, ok := sym.Type.(*StructType)
	if !ok {
		t.Fatalf("Person: got %T, want *StructType", sym.Type)
	}
	if st.Name != "Person" {
		t.Errorf("struct name: got %q, want Person", st.Name)
	}
	if len(st.Fields) != 2 {
		t.Fatalf("Person fields: got %d, want 2", len(st.Fields))
	}
}

func TestStructuralCompatibility(t *testing.T) {
	// Employee is compatible with Person (has all of Person's fields + more).
	person := &StructType{Name: "Person", Fields: []StructField{
		{Name: "name", Type: TypString},
		{Name: "age", Type: TypInt},
	}}
	employee := &StructType{Name: "Employee", Fields: []StructField{
		{Name: "name", Type: TypString},
		{Name: "age", Type: TypInt},
		{Name: "id", Type: TypInt},
	}}
	if !employee.IsCompatibleWith(person) {
		t.Errorf("Employee should be compatible with Person")
	}
	if person.IsCompatibleWith(employee) {
		t.Errorf("Person should NOT be compatible with Employee (missing 'id')")
	}
}

func TestIncompatibleTypes(t *testing.T) {
	point := &StructType{Name: "Point", Fields: []StructField{
		{Name: "x", Type: TypInt},
		{Name: "y", Type: TypInt},
	}}
	person := &StructType{Name: "Person", Fields: []StructField{
		{Name: "name", Type: TypString},
		{Name: "age", Type: TypInt},
	}}
	if point.IsCompatibleWith(person) {
		t.Errorf("Point should not be compatible with Person")
	}
}

// ---------------------------------------------------------------------------
// 07 — Mix
// ---------------------------------------------------------------------------

func TestMixConflictError(t *testing.T) {
	expectError(t, `A: { name String }
B: { name String; age Int }
C: {
    mix A
    mix B
}`, "conflict")
}

func TestMixSuccess(t *testing.T) {
	a := mustAnalyze(t, `Timestamps: {
    created_at Int
    updated_at Int
}
Post: {
    mix Timestamps
    title String
}`)
	sym := a.LookupInFile("Post")
	if sym == nil {
		t.Fatal("'Post' not found")
	}
	st, ok := sym.Type.(*StructType)
	if !ok {
		t.Fatalf("Post: got %T, want *StructType", sym.Type)
	}
	// Post should have all 3 fields after mix.
	if len(st.Fields) != 3 {
		t.Fatalf("Post fields: got %d, want 3 (created_at, updated_at, title)", len(st.Fields))
	}
}

// ---------------------------------------------------------------------------
// 08 — FQN
// ---------------------------------------------------------------------------

func TestTopLevelFQN(t *testing.T) {
	a := mustAnalyze(t, `counter: { count: 0 }`)
	sym := a.LookupInFile("counter")
	if sym == nil {
		t.Fatal("'counter' not found")
	}
	// Top-level FQN in a file with no path prefix is just the name.
	if sym.FQN != "counter" {
		t.Errorf("FQN: got %q, want %q", sym.FQN, "counter")
	}
}

// ---------------------------------------------------------------------------
// 09 — Non-word method name mapping
// ---------------------------------------------------------------------------

func TestNonWordMethodName(t *testing.T) {
	cases := []struct{ op, want string }{
		{"+", "plus"},
		{"-", "minus"},
		{"*", "star"},
		{"/", "slash"},
		{"==", "eqeq"},
		{"!=", "neq"},
		{"&&", "ampamp"},
		{"||", "pipepipe"},
		{"?", "qm"},
		{"<=", "lteq"},
		{">=", "gteq"},
		{"++", "plusplus"},
	}
	for _, tc := range cases {
		got := NonWordMethodName(tc.op)
		if got != tc.want {
			t.Errorf("NonWordMethodName(%q): got %q, want %q", tc.op, got, tc.want)
		}
	}
}

// ---------------------------------------------------------------------------
// 10 — Array and dict types
// ---------------------------------------------------------------------------

func TestArrayLiteralType(t *testing.T) {
	a := mustAnalyze(t, `nums: [1, 2, 3]`)
	sym := a.LookupInFile("nums")
	if sym == nil {
		t.Fatal("'nums' not found")
	}
	at, ok := sym.Type.(*ArrayType)
	if !ok {
		t.Fatalf("nums: got %T, want *ArrayType", sym.Type)
	}
	if at.Elem != TypInt {
		t.Errorf("nums elem: got %v, want Int", at.Elem)
	}
}

func TestDictLiteralType(t *testing.T) {
	a := mustAnalyze(t, `ages: ["Alice": 30, "Bob": 25]`)
	sym := a.LookupInFile("ages")
	if sym == nil {
		t.Fatal("'ages' not found")
	}
	dt, ok := sym.Type.(*DictType)
	if !ok {
		t.Fatalf("ages: got %T, want *DictType", sym.Type)
	}
	if dt.Key != TypString || dt.Val != TypInt {
		t.Errorf("ages: got [%v:%v], want [String:Int]", dt.Key, dt.Val)
	}
}

// ---------------------------------------------------------------------------
// 11 — Realistic: counter program
// ---------------------------------------------------------------------------

func TestCounterProgram(t *testing.T) {
	src := `counter: {
    count: 0
    increment: { count = count + 1 }
    value: { count }
}`
	a := mustAnalyze(t, src)
	sym := a.LookupInFile("counter")
	if sym == nil {
		t.Fatal("'counter' not found")
	}
	if _, ok := sym.Type.(*BocType); !ok {
		t.Fatalf("counter: got %T, want *BocType", sym.Type)
	}
}

// ---------------------------------------------------------------------------
// 15 — Variant (sum) types
// ---------------------------------------------------------------------------

func TestVariantTypeDecl(t *testing.T) {
	a := mustAnalyze(t, `Pet: {
    Cat(name String, lives Int),
    Dog(name String, years Int),
}`)
	sym := a.LookupInFile("Pet")
	if sym == nil {
		t.Fatal("'Pet' not found")
	}
	st, ok := sym.Type.(*StructType)
	if !ok {
		t.Fatalf("Pet: got %T, want *StructType", sym.Type)
	}
	if !st.IsVariant {
		t.Error("Pet.IsVariant: want true")
	}
	if len(st.Variants) != 2 {
		t.Fatalf("variants: want 2, got %d", len(st.Variants))
	}
	if st.Variants[0].Name != "Cat" {
		t.Errorf("variants[0].Name: got %q, want Cat", st.Variants[0].Name)
	}
	if st.Variants[1].Name != "Dog" {
		t.Errorf("variants[1].Name: got %q, want Dog", st.Variants[1].Name)
	}
	// Merged flat fields: name (shared), lives, years
	if len(st.Fields) != 3 {
		t.Errorf("merged fields: want 3, got %d: %+v", len(st.Fields), st.Fields)
	}
}

func TestVariantConstructorRegistered(t *testing.T) {
	a := mustAnalyze(t, `Pet: {
    Cat(name String, lives Int),
    Dog(name String, years Int),
}`)
	catSym := a.LookupInFile("Cat")
	if catSym == nil {
		t.Fatal("'Cat' constructor not registered")
	}
	bt, ok := catSym.Type.(*BocType)
	if !ok {
		t.Fatalf("Cat: got %T, want *BocType", catSym.Type)
	}
	if len(bt.Params) != 2 {
		t.Fatalf("Cat params: want 2, got %d", len(bt.Params))
	}
	if bt.Params[0].Label != "name" || bt.Params[1].Label != "lives" {
		t.Errorf("Cat param names: got %q %q", bt.Params[0].Label, bt.Params[1].Label)
	}
	// Return type is the parent StructType (Pet).
	if len(bt.Returns) != 1 {
		t.Fatalf("Cat returns: want 1, got %d", len(bt.Returns))
	}
	retSt, ok := bt.Returns[0].(*StructType)
	if !ok {
		t.Fatalf("Cat return: got %T, want *StructType", bt.Returns[0])
	}
	if retSt.Name != "Pet" {
		t.Errorf("Cat return name: got %q, want Pet", retSt.Name)
	}
}
