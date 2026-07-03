package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"yz/internal/ast"
	"yz/internal/codegen"
	"yz/internal/ir"
	"yz/internal/sema"
	"yz/runtime/macrowire"
)

// debugMacroYz is the Debug macro used across bootstrap tests. It generates a
// `debug #(String)` method that interpolates every field of the subject.
// `${` never appears literally in macro source (it would start interpolation
// in the macro itself) — it is assembled via "$" + "{".
const debugMacroYz = `
Debug : {
    Schema : NoConfig
    run #(subject Boc, config NoConfig, Boc) {
        body: subject.name + "("
        sep: ""
        i: 0
        while({ i < subject.fields.length() }, {
            f: subject.fields.at(i)
            body = body + sep + f.name + ": $" + "{" + f.name + "}"
            sep = ", "
            i = i + 1
        })
        generated("debug #(String) {\n    \"" + body + ")\"\n}")
    }
}
`

// TestPreludeCompiles pushes the prelude plus the Debug macro through the
// full unwrapped pipeline (parse → sema → lower → codegen) and checks the
// key generated shapes the synthesized main depends on.
func TestPreludeCompiles(t *testing.T) {
	prelude, err := preludeStmts()
	if err != nil {
		t.Fatal(err)
	}
	stmts := parseStmts(t, debugMacroYz)
	sf := &ast.SourceFile{Stmts: append(append([]ast.Node{}, prelude...), stmts...)}
	a := sema.NewAnalyzer()
	if err := a.AnalyzeFile(sf); err != nil {
		t.Fatalf("sema: %v", err)
	}
	goSrc := codegen.Generate(ir.Lower(sf, a, "main"))
	for _, want := range []string{
		"func NewBoc(name std.String, fields std.Array[*Field], source std.String) *Boc",
		"func NewField(name std.String, type_ std.String) *Field",
		"func (self *Boc) Source() std.String",
		"func (self *Debug) Run(subject *Boc, config *NoConfig) *std.Thunk[*Boc]",
		"func NewDebug() *Debug",
		"func NewNoConfig() *NoConfig",
	} {
		if !strings.Contains(goSrc, want) {
			t.Errorf("generated Go missing %q\n\n%s", want, goSrc)
		}
	}
}

// TestGenMacroMain checks the synthesized dispatch main.
func TestGenMacroMain(t *testing.T) {
	defs := []*macroDef{
		{Name: "Debug", SchemaTypeName: "NoConfig"},
		{Name: "JSON", SchemaTypeName: "JSONConfig", SchemaFields: []macrowire.FieldSpec{
			{Name: "pretty", Type: "Bool"},
			{Name: "indent", Type: "Int"},
		}},
	}
	src := genMacroMain(defs)
	for _, want := range []string{
		`case "Debug":`,
		"NewDebug().Run(subject, NewNoConfig()).Force()",
		`case "JSON":`,
		`NewJSON().Run(subject, NewJSONConfig(std.NewBool(cfg["pretty"].Bool), std.NewInt(cfg["indent"].Int))).Force()`,
		"out.Source().GoString()",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("macro main missing %q\n\n%s", want, src)
		}
	}
}

// TestMacroSourceHash checks hash stability and change detection.
func TestMacroSourceHash(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "debug.yz")
	if err := os.WriteFile(path, []byte(debugMacroYz), 0o644); err != nil {
		t.Fatal(err)
	}
	files := []fileEntry{{absPath: path, relDir: "macros", name: "debug"}}
	h1, err := macroSourceHash(files)
	if err != nil {
		t.Fatal(err)
	}
	h2, err := macroSourceHash(files)
	if err != nil {
		t.Fatal(err)
	}
	if h1 != h2 {
		t.Error("hash not stable")
	}
	if err := os.WriteFile(path, []byte(debugMacroYz+"\n// changed"), 0o644); err != nil {
		t.Fatal(err)
	}
	h3, err := macroSourceHash(files)
	if err != nil {
		t.Fatal(err)
	}
	if h3 == h1 {
		t.Error("hash did not change with source")
	}
}

// TestBootstrapDebugEndToEnd builds the Debug macro package to a native
// executable and runs it with a Person payload. Slow (go build) — skipped
// under -short.
func TestBootstrapDebugEndToEnd(t *testing.T) {
	if testing.Short() {
		t.Skip("bootstrap build is slow; run without -short")
	}
	projectDir := t.TempDir()
	macroDir := filepath.Join(projectDir, "macros")
	if err := os.MkdirAll(macroDir, 0o755); err != nil {
		t.Fatal(err)
	}
	macroPath := filepath.Join(macroDir, "debug.yz")
	if err := os.WriteFile(macroPath, []byte(debugMacroYz), 0o644); err != nil {
		t.Fatal(err)
	}

	files := []fileEntry{{absPath: macroPath, relDir: "macros", name: "debug"}}
	defs, err := scanMacroDefs(parseStmts(t, debugMacroYz), "macros")
	if err != nil || len(defs) != 1 {
		t.Fatalf("scan: defs=%v err=%v", defs, err)
	}
	reg := newMacroRegistry()
	if err := reg.register(defs[0]); err != nil {
		t.Fatal(err)
	}

	if err := bootstrapMacroPackage(projectDir, files, "macros", reg); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}
	binPath := reg.binPath["macros"]
	if binPath == "" {
		t.Fatal("no binary path registered")
	}

	payload := (&macrowire.Payload{
		SubjectName: "Person",
		Fields: []macrowire.FieldSpec{
			{Name: "name", Type: "String"},
			{Name: "age", Type: "Int"},
		},
	}).Encode()

	cmd := exec.Command(binPath, "Debug")
	cmd.Stdin = strings.NewReader(payload)
	out, err := cmd.Output()
	if err != nil {
		stderr := ""
		if ee, ok := err.(*exec.ExitError); ok {
			stderr = string(ee.Stderr)
		}
		t.Fatalf("running macro: %v\n%s", err, stderr)
	}
	want := "debug #(String) {\n    \"Person(name: ${name}, age: ${age})\"\n}"
	if string(out) != want {
		t.Errorf("macro output:\n%s\nwant:\n%s", out, want)
	}

	// Second bootstrap must reuse the cached binary.
	info1, _ := os.Stat(binPath)
	if err := bootstrapMacroPackage(projectDir, files, "macros", reg); err != nil {
		t.Fatalf("re-bootstrap: %v", err)
	}
	info2, _ := os.Stat(binPath)
	if !info1.ModTime().Equal(info2.ModTime()) {
		t.Error("cached binary was rebuilt despite unchanged source")
	}
}
