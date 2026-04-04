package codegen

import (
	"strings"
	"testing"

	"yz/internal/ir"
	"yz/internal/lexer"
	"yz/internal/parser"
	"yz/internal/sema"
)

// genWithPackages compiles Yz source to Go, with pre-registered sub-packages
// (simulating what build.go does for cross-package FQN resolution).
func genWithPackages(t *testing.T, src string, pkgs map[string]map[string]*sema.Symbol) string {
	t.Helper()
	_ = lexer.Tokenize([]byte(src))
	p := parser.New([]byte(src))
	sf, err := p.ParseFile()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	a := sema.NewAnalyzer()
	for relDir, exports := range pkgs {
		parts := strings.Split(relDir, "/")
		pkgAlias := parts[len(parts)-1]
		importPath := "yzapp/" + relDir
		a.RegisterPackage(relDir, pkgAlias, importPath, exports)
	}
	if err := a.AnalyzeFile(sf); err != nil {
		t.Fatalf("sema: %v", err)
	}
	f := ir.Lower(sf, a, "main")
	return Generate(f)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// gen compiles Yz source all the way to Go source string.
func gen(t *testing.T, src string) string {
	t.Helper()
	_ = lexer.Tokenize([]byte(src))
	p := parser.New([]byte(src))
	sf, err := p.ParseFile()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	a := sema.NewAnalyzer()
	if err := a.AnalyzeFile(sf); err != nil {
		t.Fatalf("sema: %v", err)
	}
	f := ir.Lower(sf, a, "main")
	return Generate(f)
}

// contains asserts that got contains all of the listed substrings.
func contains(t *testing.T, got string, wants ...string) {
	t.Helper()
	for _, w := range wants {
		if !strings.Contains(got, w) {
			t.Errorf("generated code missing %q\n\n--- generated ---\n%s", w, got)
		}
	}
}

// notContains asserts that got does NOT contain substr.
func notContains(t *testing.T, got, substr string) {
	t.Helper()
	if strings.Contains(got, substr) {
		t.Errorf("generated code should NOT contain %q\n\n--- generated ---\n%s", substr, got)
	}
}

// ---------------------------------------------------------------------------
// 01 — Package and imports
// ---------------------------------------------------------------------------

func TestGeneratePackageDecl(t *testing.T) {
	got := gen(t, `x: 42`)
	contains(t, got, "package main")
}

func TestGenerateRuntimeImport(t *testing.T) {
	got := gen(t, `x: 42`)
	contains(t, got, `std "yz/runtime/yzrt"`)
}

// ---------------------------------------------------------------------------
// 02 — Literal boxing
// ---------------------------------------------------------------------------

func TestGenerateIntLitBoxed(t *testing.T) {
	got := gen(t, `x: 42`)
	contains(t, got, "std.NewInt(42)")
}

func TestGenerateStringLitBoxed(t *testing.T) {
	got := gen(t, `name: "Alice"`)
	contains(t, got, `std.NewString("Alice")`)
}

func TestGenerateBoolLit(t *testing.T) {
	got := gen(t, `flag: true`)
	contains(t, got, "std.NewBool(true)")
}

func TestGenerateDecimalLit(t *testing.T) {
	got := gen(t, `pi: 3.14`)
	contains(t, got, "std.NewDecimal(3.14)")
}

// ---------------------------------------------------------------------------
// 03 — Singleton struct type and var
// ---------------------------------------------------------------------------

func TestGenerateSingletonStructType(t *testing.T) {
	got := gen(t, `counter: {
    count: 0
}`)
	contains(t, got,
		"type _counterBoc struct",
		"count std.Int",
		"var Counter = &_counterBoc{",
		"count: std.NewInt(0)",
	)
}

// ---------------------------------------------------------------------------
// 04 — Method emission with goroutine thunk
// ---------------------------------------------------------------------------

func TestGenerateMethodThunk(t *testing.T) {
	got := gen(t, `counter: {
    count: 0
    value: { count }
}`)
	contains(t, got,
		"func (self *_counterBoc) Value()",
		"*std.Thunk[std.Int]",
		"return std.Go(func() std.Int {",
		"return self.count",
	)
}

func TestGenerateIncrementMethod(t *testing.T) {
	got := gen(t, `counter: {
    count: 0
    increment: { count = count + 1 }
}`)
	contains(t, got,
		"func (self *_counterBoc) Increment()",
		"std.Go(func() std.Unit {",
		"self.count = self.count.Plus(std.NewInt(1))",
	)
}

// ---------------------------------------------------------------------------
// 05 — Method call (binary operator → MethodCall)
// ---------------------------------------------------------------------------

func TestGenerateBinaryMethodCall(t *testing.T) {
	got := gen(t, `counter: {
    count: 0
    next: { count + 1 }
}`)
	contains(t, got, ".Plus(std.NewInt(1))")
}

// ---------------------------------------------------------------------------
// 06 — Struct type (uppercase boc)
// ---------------------------------------------------------------------------

func TestGenerateStructDecl(t *testing.T) {
	got := gen(t, `Person: {
    name String
    age Int
}`)
	contains(t, got,
		"type Person struct",
		"name std.String",
		"age std.Int",
	)
}

// ---------------------------------------------------------------------------
// 07 — main boc becomes func main
// ---------------------------------------------------------------------------

func TestGenerateMainFunc(t *testing.T) {
	got := gen(t, `main: {
    x: 42
}`)
	contains(t, got, "func main()")
	notContains(t, got, "type _mainBoc")
}

// ---------------------------------------------------------------------------
// 08 — Field access uses self.field
// ---------------------------------------------------------------------------

func TestGenerateFieldAccessInMethod(t *testing.T) {
	got := gen(t, `counter: {
    count: 0
    get: { count }
}`)
	// Method body must reference self.count, not bare count.
	contains(t, got, "self.count")
}

// ---------------------------------------------------------------------------
// 09 — Structured concurrency: BocGroup in main
// ---------------------------------------------------------------------------

func TestGenerateBocGroupInMain(t *testing.T) {
	got := gen(t, `counter: {
    count: 0
    increment: { count = count + 1 }
    value: { count }
}
main: {
    counter.increment()
    counter.increment()
    print(counter.value())
}`)
	contains(t, got,
		"_bg0 := &std.BocGroup{}",
		"_bg0.Go(func() any {",
		"Counter.Increment().Force()",
		"_bg0.Wait()",
		"std.Print(Counter.Value().Force())",
	)
}

// ---------------------------------------------------------------------------
// 10 — Full counter program compiles as valid Go
// ---------------------------------------------------------------------------

func TestGenerateCounterProgram(t *testing.T) {
	got := gen(t, `counter: {
    count: 0
    increment: { count = count + 1 }
    value: { count }
}`)
	contains(t, got,
		"package main",
		"type _counterBoc struct",
		"func (self *_counterBoc) Increment()",
		"func (self *_counterBoc) Value()",
		"var Counter = &_counterBoc{",
	)
}

// ---------------------------------------------------------------------------
// 11 — Cross-package FQN resolution
// ---------------------------------------------------------------------------

func TestFQNStructConstructor(t *testing.T) {
	// Simulate a "front" package (at relDir "front") exporting Host struct.
	hostStruct := &sema.StructType{Name: "Host", Fields: []sema.StructField{
		{Name: "name", Type: sema.TypString},
	}}
	exports := map[string]*sema.Symbol{
		"Host": {Name: "Host", Type: hostStruct},
	}
	got := genWithPackages(t, `main: {
    h: front.Host("Alice")
}`, map[string]map[string]*sema.Symbol{
		"front": exports,
	})
	contains(t, got,
		`"yzapp/front"`,
		`front.NewHost(std.NewString("Alice"))`,
	)
}

func TestFQNNestedNamespace(t *testing.T) {
	// Simulate "house/front" package exporting Host struct.
	hostStruct := &sema.StructType{Name: "Host", Fields: []sema.StructField{
		{Name: "name", Type: sema.TypString},
	}}
	exports := map[string]*sema.Symbol{
		"Host": {Name: "Host", Type: hostStruct},
	}
	got := genWithPackages(t, `main: {
    h: house.front.Host("Alice")
}`, map[string]map[string]*sema.Symbol{
		"house/front": exports,
	})
	contains(t, got,
		`"yzapp/house/front"`,
		`front.NewHost(std.NewString("Alice"))`,
	)
}
