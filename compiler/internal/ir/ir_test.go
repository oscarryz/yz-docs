package ir

import (
	"testing"

	"yz/internal/lexer"
	"yz/internal/parser"
	"yz/internal/sema"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func lower(t *testing.T, src string) *File {
	t.Helper()
	_ = lexer.Tokenize([]byte(src)) // validate lexer is happy
	p := parser.New([]byte(src))
	sf, err := p.ParseFile()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	a := sema.NewAnalyzer()
	if err := a.AnalyzeFile(sf); err != nil {
		t.Fatalf("sema error: %v", err)
	}
	return Lower(sf, a, "main")
}

// ---------------------------------------------------------------------------
// 01 — Literals produce the right IR leaf nodes
// ---------------------------------------------------------------------------

func TestLowerIntLit(t *testing.T) {
	f := lower(t, `x: 42`)
	if len(f.Decls) != 1 {
		t.Fatalf("want 1 decl, got %d", len(f.Decls))
	}
	sd, ok := f.Decls[0].(*SingletonDecl)
	if !ok {
		t.Fatalf("want *SingletonDecl, got %T", f.Decls[0])
	}
	if len(sd.Fields) != 1 {
		t.Fatalf("want 1 field, got %d", len(sd.Fields))
	}
	field := sd.Fields[0]
	if field.Name != "x" {
		t.Errorf("field name: got %q, want x", field.Name)
	}
	if field.Type != "std.Int" {
		t.Errorf("field type: got %q, want std.Int", field.Type)
	}
	if _, ok := field.Init.(*IntLit); !ok {
		t.Errorf("field init: got %T, want *IntLit", field.Init)
	}
}

func TestLowerStringLit(t *testing.T) {
	f := lower(t, `name: "Alice"`)
	sd := f.Decls[0].(*SingletonDecl)
	if sd.Fields[0].Type != "std.String" {
		t.Errorf("got %q, want std.String", sd.Fields[0].Type)
	}
	if _, ok := sd.Fields[0].Init.(*StringLit); !ok {
		t.Errorf("init: got %T, want *StringLit", sd.Fields[0].Init)
	}
}

func TestLowerBoolLit(t *testing.T) {
	f := lower(t, `active: true`)
	sd := f.Decls[0].(*SingletonDecl)
	if sd.Fields[0].Type != "std.Bool" {
		t.Errorf("got %q, want std.Bool", sd.Fields[0].Type)
	}
}

// ---------------------------------------------------------------------------
// 02 — Singleton boc lowering
// ---------------------------------------------------------------------------

func TestLowerSingletonBoc(t *testing.T) {
	f := lower(t, `counter: {
    count: 0
}`)
	if len(f.Decls) != 1 {
		t.Fatalf("want 1 decl, got %d", len(f.Decls))
	}
	outer, ok := f.Decls[0].(*SingletonDecl)
	if !ok {
		t.Fatalf("want *SingletonDecl, got %T", f.Decls[0])
	}
	if outer.VarName != "counter" {
		t.Errorf("VarName: got %q, want counter", outer.VarName)
	}
	if outer.TypeName != "_counterBoc" {
		t.Errorf("TypeName: got %q, want _counterBoc", outer.TypeName)
	}
	if len(outer.Fields) != 1 || outer.Fields[0].Name != "count" {
		t.Errorf("fields: %v", outer.Fields)
	}
}

func TestLowerSingletonWithMethods(t *testing.T) {
	f := lower(t, `counter: {
    count: 0
    increment: { count = count + 1 }
    value: { count }
}`)
	outer := f.Decls[0].(*SingletonDecl)
	// Fields: only "count" (increment and value are methods)
	if len(outer.Fields) != 1 {
		t.Errorf("fields: want 1, got %d", len(outer.Fields))
	}
	if len(outer.Methods) != 2 {
		t.Errorf("methods: want 2, got %d", len(outer.Methods))
	}
	names := map[string]bool{}
	for _, m := range outer.Methods {
		names[m.Name] = true
	}
	if !names["increment"] || !names["value"] {
		t.Errorf("method names: %v", names)
	}
}

// ---------------------------------------------------------------------------
// 03 — Struct type (uppercase boc)
// ---------------------------------------------------------------------------

func TestLowerStructDecl(t *testing.T) {
	f := lower(t, `Person: {
    name String
    age Int
}`)
	if len(f.Decls) != 1 {
		t.Fatalf("want 1 decl, got %d", len(f.Decls))
	}
	sd, ok := f.Decls[0].(*StructDecl)
	if !ok {
		t.Fatalf("want *StructDecl, got %T", f.Decls[0])
	}
	if sd.Name != "Person" {
		t.Errorf("Name: got %q, want Person", sd.Name)
	}
	if len(sd.Fields) != 2 {
		t.Fatalf("fields: want 2, got %d", len(sd.Fields))
	}
	if sd.Fields[0].Name != "name" || sd.Fields[0].Type != "std.String" {
		t.Errorf("field[0]: %+v", sd.Fields[0])
	}
	if sd.Fields[1].Name != "age" || sd.Fields[1].Type != "std.Int" {
		t.Errorf("field[1]: %+v", sd.Fields[1])
	}
}

// ---------------------------------------------------------------------------
// 04 — Binary expression → MethodCall
// ---------------------------------------------------------------------------

func TestLowerBinaryExprToMethodCall(t *testing.T) {
	f := lower(t, `counter: {
    count: 0
    next: { count + 1 }
}`)
	outer := f.Decls[0].(*SingletonDecl)
	// Find the "next" method
	var nextMethod *MethodDecl
	for _, m := range outer.Methods {
		if m.Name == "next" {
			nextMethod = m
		}
	}
	if nextMethod == nil {
		t.Fatal("method 'next' not found")
	}
	// Body should be a ThunkExpr wrapping a ReturnStmt with a MethodCall
	if len(nextMethod.Body) != 1 {
		t.Fatalf("method body stmts: want 1, got %d", len(nextMethod.Body))
	}
	th, ok := nextMethod.Body[0].(*ExprStmt)
	if !ok {
		t.Fatalf("body[0]: got %T, want *ExprStmt", nextMethod.Body[0])
	}
	thunk, ok := th.Expr.(*ThunkExpr)
	if !ok {
		t.Fatalf("expr: got %T, want *ThunkExpr", th.Expr)
	}
	// The thunk body should have a ReturnStmt with the MethodCall
	if len(thunk.Body) < 1 {
		t.Fatal("thunk body is empty")
	}
	ret, ok := thunk.Body[len(thunk.Body)-1].(*ReturnStmt)
	if !ok {
		t.Fatalf("thunk body last: got %T, want *ReturnStmt", thunk.Body[len(thunk.Body)-1])
	}
	mc, ok := ret.Value.(*MethodCall)
	if !ok {
		t.Fatalf("return value: got %T, want *MethodCall", ret.Value)
	}
	if mc.Method != "Plus" {
		t.Errorf("method name: got %q, want Plus", mc.Method)
	}
}

// ---------------------------------------------------------------------------
// 05 — Method return type wraps in Thunk
// ---------------------------------------------------------------------------

func TestLowerMethodReturnsThunk(t *testing.T) {
	f := lower(t, `counter: {
    count: 0
    value: { count }
}`)
	outer := f.Decls[0].(*SingletonDecl)
	var valueMethod *MethodDecl
	for _, m := range outer.Methods {
		if m.Name == "value" {
			valueMethod = m
		}
	}
	if valueMethod == nil {
		t.Fatal("method 'value' not found")
	}
	// Results should be "*std.Thunk[std.Int]"
	if len(valueMethod.Results) != 1 {
		t.Fatalf("results: want 1, got %d", len(valueMethod.Results))
	}
	if valueMethod.Results[0] != "*std.Thunk[std.Int]" {
		t.Errorf("result type: got %q, want *std.Thunk[std.Int]", valueMethod.Results[0])
	}
}

// ---------------------------------------------------------------------------
// 06 — main boc becomes FuncDecl
// ---------------------------------------------------------------------------

func TestLowerMainBecomesFuncDecl(t *testing.T) {
	f := lower(t, `main: {
    x: 42
}`)
	if len(f.Decls) != 1 {
		t.Fatalf("want 1 decl, got %d", len(f.Decls))
	}
	fn, ok := f.Decls[0].(*FuncDecl)
	if !ok {
		t.Fatalf("want *FuncDecl, got %T", f.Decls[0])
	}
	if fn.Name != "main" {
		t.Errorf("func name: got %q, want main", fn.Name)
	}
}

// ---------------------------------------------------------------------------
// 07 — Typed declaration becomes FieldSpec with correct type
// ---------------------------------------------------------------------------

func TestLowerTypedDeclField(t *testing.T) {
	f := lower(t, `Point: {
    x Int
    y Int
}`)
	sd := f.Decls[0].(*StructDecl)
	if sd.Fields[0].Type != "std.Int" {
		t.Errorf("x type: got %q, want std.Int", sd.Fields[0].Type)
	}
}

// ---------------------------------------------------------------------------
// 08 — Receiver name in methods
// ---------------------------------------------------------------------------

func TestLowerMethodReceiverName(t *testing.T) {
	f := lower(t, `counter: {
    count: 0
    get: { count }
}`)
	outer := f.Decls[0].(*SingletonDecl)
	m := outer.Methods[0]
	if m.RecvName != "self" {
		t.Errorf("RecvName: got %q, want self", m.RecvName)
	}
	if m.RecvType != "*_counterBoc" {
		t.Errorf("RecvType: got %q, want *_counterBoc", m.RecvType)
	}
}
