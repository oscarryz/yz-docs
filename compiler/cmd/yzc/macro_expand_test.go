package main

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"yz/internal/ast"
	"yz/internal/sema"
)

// findShortDecl returns the top-level ShortDecl with the given name.
func findShortDecl(t *testing.T, stmts []ast.Node, name string) *ast.ShortDecl {
	t.Helper()
	for _, stmt := range stmts {
		if sd, ok := stmt.(*ast.ShortDecl); ok && len(sd.Names) == 1 && sd.Names[0].Name == name {
			return sd
		}
	}
	t.Fatalf("ShortDecl %q not found", name)
	return nil
}

// primeRunCache writes canned macro output into the run cache so invokeMacro
// never execs a binary.
func primeRunCache(t *testing.T, projectDir string, def *macroDef, payload, output string) {
	t.Helper()
	runKey := sha256.Sum256([]byte(def.Name + "\x00" + payload))
	runsDir := filepath.Join(projectDir, "target", "macros", macroPkgKey(def.RelDir), "runs")
	if err := os.MkdirAll(runsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	cachePath := filepath.Join(runsDir, hex.EncodeToString(runKey[:])+".out")
	if err := os.WriteFile(cachePath, []byte(output), 0o644); err != nil {
		t.Fatal(err)
	}
}

// builtRegistry returns a registry with one already-built Debug macro in
// package "macros".
func builtRegistry(t *testing.T) (*macroRegistry, *macroDef) {
	t.Helper()
	reg := newMacroRegistry()
	def := &macroDef{Name: "Debug", RelDir: "macros", SchemaTypeName: "NoConfig"}
	if err := reg.register(def); err != nil {
		t.Fatal(err)
	}
	reg.state["macros"] = macroBuilt
	return reg, def
}

const annotatedPersonSrc = "`\nDebug: {}\n`\nPerson : {\n    name String\n}\n\nmain: {\n    p: Person(\"Ada\")\n    print(p.debug())\n}\n"

// TestExpandMergesGeneratedSlots runs expansion with a cached macro result
// and verifies the generated method lands in the subject and resolves in
// sema (the call site p.debug() type-checks).
func TestExpandMergesGeneratedSlots(t *testing.T) {
	projectDir := t.TempDir()
	reg, def := builtRegistry(t)

	stmts := parseStmts(t, annotatedPersonSrc)
	person := findShortDecl(t, stmts, "Person")
	personBl := person.Values[0].(*ast.BocLiteral)
	before := len(personBl.Elements)

	payload := buildSubjectPayload("Person", personBl, nil).Encode()
	primeRunCache(t, projectDir, def, payload, "debug #(String) {\n    \"Person(...)\"\n}")

	if err := expandMacros(stmts, "", reg, nil, projectDir); err != nil {
		t.Fatalf("expandMacros: %v", err)
	}
	if len(personBl.Elements) != before+1 {
		t.Fatalf("expected 1 merged element, got %d new", len(personBl.Elements)-before)
	}
	bd, ok := personBl.Elements[before].(*ast.BocDecl)
	if !ok || bd.Name.Name != "debug" {
		t.Fatalf("merged element is %T, want BocDecl debug", personBl.Elements[before])
	}

	// The merged AST must type-check as if hand-written.
	sf := &ast.SourceFile{Stmts: stmts}
	if err := sema.NewAnalyzer().AnalyzeFile(sf); err != nil {
		t.Errorf("sema after merge: %v", err)
	}
}

// TestExpandInsideFileWrapper checks expansion descends into file-wrapper
// bocs (the shape compilePackageDir produces).
func TestExpandInsideFileWrapper(t *testing.T) {
	projectDir := t.TempDir()
	reg, def := builtRegistry(t)

	stmts := parseStmts(t, annotatedPersonSrc)
	person := findShortDecl(t, stmts, "Person")
	personBl := person.Values[0].(*ast.BocLiteral)
	before := len(personBl.Elements)
	wrapped := []ast.Node{
		&ast.ShortDecl{
			Names:         []*ast.Ident{{Name: "main"}},
			Values:        []ast.Expr{&ast.BocLiteral{Elements: stmts}},
			IsFileWrapper: true,
		},
	}

	payload := buildSubjectPayload("Person", personBl, nil).Encode()
	primeRunCache(t, projectDir, def, payload, "debug #(String) {\n    \"Person(...)\"\n}")

	if err := expandMacros(wrapped, "", reg, nil, projectDir); err != nil {
		t.Fatalf("expandMacros: %v", err)
	}
	if len(personBl.Elements) != before+1 {
		t.Errorf("expansion did not descend into file wrapper")
	}
}

func TestExpandUnknownMacro(t *testing.T) {
	reg, _ := builtRegistry(t)
	stmts := parseStmts(t, "`\nDebgu: {}\n`\nPerson : {\n    name String\n}\n")
	err := expandMacros(stmts, "", reg, nil, t.TempDir())
	if err == nil || !strings.Contains(err.Error(), `unknown macro "Debgu"`) {
		t.Errorf("error = %v", err)
	}
}

func TestExpandSamePackage(t *testing.T) {
	reg, _ := builtRegistry(t)
	stmts := parseStmts(t, "`\nDebug: {}\n`\nPerson : {\n    name String\n}\n")
	err := expandMacros(stmts, "macros", reg, nil, t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "cannot be applied to a boc in the same package") {
		t.Errorf("error = %v", err)
	}
}

func TestEnsureMacroPackageBuiltCycle(t *testing.T) {
	reg := newMacroRegistry()
	def := &macroDef{Name: "GenA", RelDir: "macros_b"}
	if err := reg.register(def); err != nil {
		t.Fatal(err)
	}
	reg.state["macros_b"] = macroBuilding
	err := ensureMacroPackageBuilt(t.TempDir(), def, reg, []string{"macros_a.GenB"})
	if err == nil || !strings.Contains(err.Error(), "macro cycle detected: macros_a.GenB -> macros_b.GenA") {
		t.Errorf("error = %v", err)
	}
}

func TestExpandNoRegistryNoop(t *testing.T) {
	stmts := parseStmts(t, annotatedPersonSrc)
	if err := expandMacros(stmts, "", newMacroRegistry(), nil, t.TempDir()); err != nil {
		t.Errorf("empty registry must be a no-op, got %v", err)
	}
	if err := expandMacros(stmts, "", nil, nil, t.TempDir()); err != nil {
		t.Errorf("nil registry must be a no-op, got %v", err)
	}
}
