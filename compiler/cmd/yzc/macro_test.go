package main

import (
	"strings"
	"testing"

	"yz/internal/ast"
	"yz/internal/parser"
	"yz/runtime/macrowire"
)

func parseStmts(t *testing.T, src string) []ast.Node {
	t.Helper()
	sf, err := parser.New([]byte(src)).ParseFile()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	return sf.Stmts
}

// ---------------------------------------------------------------------------
// scanMacroDefs
// ---------------------------------------------------------------------------

const debugMacroSrc = `
Debug : {
    Schema : NoConfig
    run #(subject Boc, config NoConfig, Boc) {
        generated("")
    }
}
`

func TestScanMacroDefsDebug(t *testing.T) {
	defs, err := scanMacroDefs(parseStmts(t, debugMacroSrc), "macros")
	if err != nil {
		t.Fatalf("scanMacroDefs: %v", err)
	}
	if len(defs) != 1 {
		t.Fatalf("expected 1 def, got %d", len(defs))
	}
	d := defs[0]
	if d.Name != "Debug" || d.RelDir != "macros" || d.SchemaTypeName != "NoConfig" || len(d.SchemaFields) != 0 {
		t.Errorf("unexpected def: %+v", d)
	}
}

func TestScanMacroDefsWithConfig(t *testing.T) {
	src := `
JSONConfig : {
    pretty Bool
    indent Int
}

JSON : {
    Schema : JSONConfig
    run #(subject Boc, config JSONConfig, Boc) {
        generated("")
    }
}
`
	defs, err := scanMacroDefs(parseStmts(t, src), "macros")
	if err != nil {
		t.Fatalf("scanMacroDefs: %v", err)
	}
	if len(defs) != 1 {
		t.Fatalf("expected 1 def, got %d (JSONConfig must not match: no run)", len(defs))
	}
	d := defs[0]
	if d.Name != "JSON" || d.SchemaTypeName != "JSONConfig" {
		t.Errorf("unexpected def: %+v", d)
	}
	want := []macrowire.FieldSpec{{Name: "pretty", Type: "Bool"}, {Name: "indent", Type: "Int"}}
	if len(d.SchemaFields) != 2 || d.SchemaFields[0] != want[0] || d.SchemaFields[1] != want[1] {
		t.Errorf("schema fields = %+v, want %+v", d.SchemaFields, want)
	}
}

func TestScanMacroDefsInlineSchema(t *testing.T) {
	src := `
Doc : {
    Schema #(documentation String)
    run #(subject Boc, config Schema, Boc) {
        generated("")
    }
}
`
	defs, err := scanMacroDefs(parseStmts(t, src), "macros")
	if err != nil {
		t.Fatalf("scanMacroDefs: %v", err)
	}
	if len(defs) != 1 {
		t.Fatalf("expected 1 def, got %d", len(defs))
	}
	d := defs[0]
	if d.SchemaTypeName != "" || len(d.SchemaFields) != 1 || d.SchemaFields[0].Name != "documentation" {
		t.Errorf("unexpected def: %+v", d)
	}
}

func TestScanMacroDefsNearMisses(t *testing.T) {
	cases := map[string]string{
		"no Schema": `
X : {
    run #(subject Boc, config NoConfig, Boc) { generated("") }
}
`,
		"no run": `
X : {
    Schema : NoConfig
}
`,
		"two-param run": `
X : {
    Schema : NoConfig
    run #(subject Boc, Boc) { generated("") }
}
`,
		"lowercase name": `
x : {
    Schema : NoConfig
    run #(subject Boc, config NoConfig, Boc) { generated("") }
}
`,
		"run first param not Boc": `
X : {
    Schema : NoConfig
    run #(subject Int, config NoConfig, Boc) { generated("") }
}
`,
		"run return not Boc": `
X : {
    Schema : NoConfig
    run #(subject Boc, config NoConfig, Int) { generated("") }
}
`,
	}
	for label, src := range cases {
		defs, err := scanMacroDefs(parseStmts(t, src), "macros")
		if err != nil {
			t.Errorf("%s: unexpected error: %v", label, err)
			continue
		}
		if len(defs) != 0 {
			t.Errorf("%s: expected no defs, got %+v", label, defs)
		}
	}
}

func TestScanMacroDefsNonScalarSchema(t *testing.T) {
	src := `
Other : {
    x Int
}

BadConfig : {
    inner Other
}

Bad : {
    Schema : BadConfig
    run #(subject Boc, config BadConfig, Boc) { generated("") }
}
`
	_, err := scanMacroDefs(parseStmts(t, src), "macros")
	if err == nil {
		t.Fatal("expected non-scalar schema error")
	}
	if !strings.Contains(err.Error(), "only scalar config fields") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestScanMacroDefsMissingSchemaType(t *testing.T) {
	src := `
Bad : {
    Schema : Missing
    run #(subject Boc, config Missing, Boc) { generated("") }
}
`
	_, err := scanMacroDefs(parseStmts(t, src), "macros")
	if err == nil || !strings.Contains(err.Error(), `schema type "Missing" not found`) {
		t.Errorf("unexpected error: %v", err)
	}
}

// ---------------------------------------------------------------------------
// annotationTriggers
// ---------------------------------------------------------------------------

// annotatedSubject parses src and returns the first annotated ShortDecl.
func annotatedSubject(t *testing.T, src string) *ast.ShortDecl {
	t.Helper()
	for _, stmt := range parseStmts(t, src) {
		if sd, ok := stmt.(*ast.ShortDecl); ok && sd.Annotation != nil {
			return sd
		}
	}
	t.Fatal("no annotated ShortDecl found")
	return nil
}

func TestAnnotationTriggers(t *testing.T) {
	src := "`\nauthor: \"oscar\"\nDebug: {}\nJSON: { pretty: true }\n`\nPerson : {\n    name String\n}\n"
	sd := annotatedSubject(t, src)
	triggers, err := annotationTriggers(sd.Annotation)
	if err != nil {
		t.Fatalf("annotationTriggers: %v", err)
	}
	if len(triggers) != 2 {
		t.Fatalf("expected 2 triggers, got %d: %+v", len(triggers), triggers)
	}
	if triggers[0].Name != "Debug" || triggers[1].Name != "JSON" {
		t.Errorf("trigger order/names wrong: %+v", triggers)
	}
	if len(triggers[0].Config.Elements) != 0 {
		t.Errorf("Debug config should be empty")
	}
	if len(triggers[1].Config.Elements) != 1 {
		t.Errorf("JSON config should have one entry")
	}
}

func TestAnnotationTriggersNonBocValue(t *testing.T) {
	src := "`\nDebug: 42\n`\nPerson : {\n    name String\n}\n"
	sd := annotatedSubject(t, src)
	if _, err := annotationTriggers(sd.Annotation); err == nil {
		t.Fatal("expected error for non-boc trigger value")
	}
}

func TestAnnotationTriggersNone(t *testing.T) {
	src := "`\nauthor: \"oscar\"\nversion: 2\n`\nPerson : {\n    name String\n}\n"
	sd := annotatedSubject(t, src)
	triggers, err := annotationTriggers(sd.Annotation)
	if err != nil || len(triggers) != 0 {
		t.Errorf("expected no triggers, got %+v (err=%v)", triggers, err)
	}
}

// ---------------------------------------------------------------------------
// validateConfig + encodeConfig
// ---------------------------------------------------------------------------

func triggerFor(t *testing.T, annSrc string) trigger {
	t.Helper()
	sd := annotatedSubject(t, annSrc+"\nPerson : {\n    name String\n}\n")
	triggers, err := annotationTriggers(sd.Annotation)
	if err != nil || len(triggers) != 1 {
		t.Fatalf("expected 1 trigger, got %+v (err=%v)", triggers, err)
	}
	return triggers[0]
}

func TestValidateConfigOK(t *testing.T) {
	def := &macroDef{Name: "JSON", SchemaFields: []macrowire.FieldSpec{
		{Name: "pretty", Type: "Bool"},
		{Name: "indent", Type: "Int"},
	}}
	tr := triggerFor(t, "`\nJSON: {\n    pretty: true\n    indent: 4\n}\n`")
	entries, err := encodeConfig(tr.Config)
	if err != nil {
		t.Fatalf("encodeConfig: %v", err)
	}
	if err := validateConfig(tr, def, entries); err != nil {
		t.Errorf("validateConfig: %v", err)
	}
}

func TestValidateConfigErrors(t *testing.T) {
	def := &macroDef{Name: "JSON", SchemaFields: []macrowire.FieldSpec{
		{Name: "pretty", Type: "Bool"},
	}}
	cases := map[string]struct {
		ann  string
		want string
	}{
		"unknown key":  {"`\nJSON: { prety: true }\n`", `config key "prety" does not match`},
		"wrong kind":   {"`\nJSON: { pretty: 1 }\n`", `expected Bool, got Int`},
		"missing key":  {"`\nJSON: {}\n`", `missing config field "pretty"`},
	}
	for label, c := range cases {
		tr := triggerFor(t, c.ann)
		entries, err := encodeConfig(tr.Config)
		if err != nil {
			t.Fatalf("%s: encodeConfig: %v", label, err)
		}
		err = validateConfig(tr, def, entries)
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: error = %v, want contains %q", label, err, c.want)
		}
	}
}

// ---------------------------------------------------------------------------
// buildSubjectPayload
// ---------------------------------------------------------------------------

func TestBuildSubjectPayload(t *testing.T) {
	src := `
Person : {
    name String
    age Int
    greet #(String) {
        "hi"
    }
}
`
	stmts := parseStmts(t, src)
	sd := stmts[0].(*ast.ShortDecl)
	bl := sd.Values[0].(*ast.BocLiteral)
	p := buildSubjectPayload("Person", bl, nil)
	if p.SubjectName != "Person" {
		t.Errorf("subject name = %q", p.SubjectName)
	}
	want := []macrowire.FieldSpec{{Name: "name", Type: "String"}, {Name: "age", Type: "Int"}}
	if len(p.Fields) != 2 || p.Fields[0] != want[0] || p.Fields[1] != want[1] {
		t.Errorf("fields = %+v, want %+v (methods must be skipped)", p.Fields, want)
	}
	// The payload must round-trip through the wire codec.
	decoded, err := macrowire.DecodePayload([]byte(p.Encode()))
	if err != nil {
		t.Fatalf("round trip: %v", err)
	}
	if decoded.SubjectName != "Person" || len(decoded.Fields) != 2 {
		t.Errorf("round trip mismatch: %+v", decoded)
	}
}

// ---------------------------------------------------------------------------
// registry
// ---------------------------------------------------------------------------

func TestRegistryDuplicate(t *testing.T) {
	reg := newMacroRegistry()
	if err := reg.register(&macroDef{Name: "Debug", RelDir: "macros"}); err != nil {
		t.Fatalf("first register: %v", err)
	}
	err := reg.register(&macroDef{Name: "Debug", RelDir: "other"})
	if err == nil || !strings.Contains(err.Error(), "defined in both") {
		t.Errorf("expected duplicate error, got %v", err)
	}
}
