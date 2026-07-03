package main

import (
	"testing"

	"yz/internal/ast"
	"yz/internal/parser"
	"yz/runtime/macrowire"
)

// parseStmts is a helper for tests that parse Yz source to AST nodes.
func parseStmts(t *testing.T, src string) []ast.Node {
	t.Helper()
	sf, err := parser.New([]byte(src)).ParseFile()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	return sf.Stmts
}

// ---------------------------------------------------------------------------
// scanMacroDefs
// ---------------------------------------------------------------------------

func TestScanMacroDefs_detects_debug_macro(t *testing.T) {
	stmts := parseStmts(t, `
Debug: {
    Schema: NoConfig
    run #(subject Boc, config Schema, Boc) {
        generated("debug")
    }
}
`)
	defs := scanMacroDefs(stmts, "macros")
	if len(defs) != 1 {
		t.Fatalf("want 1 def, got %d", len(defs))
	}
	if defs[0].Name != "Debug" {
		t.Errorf("name: got %q, want Debug", defs[0].Name)
	}
	if defs[0].RelDir != "macros" {
		t.Errorf("relDir: got %q", defs[0].RelDir)
	}
}

func TestScanMacroDefs_near_miss_lowercase(t *testing.T) {
	stmts := parseStmts(t, `
debug: {
    Schema: NoConfig
    run #(subject Boc, config Schema, Boc) {}
}
`)
	defs := scanMacroDefs(stmts, "macros")
	if len(defs) != 0 {
		t.Errorf("want 0 defs (lowercase), got %d", len(defs))
	}
}

func TestScanMacroDefs_near_miss_no_schema(t *testing.T) {
	stmts := parseStmts(t, `
Debug: {
    run #(subject Boc, config NoConfig, Boc) {}
}
`)
	defs := scanMacroDefs(stmts, "macros")
	if len(defs) != 0 {
		t.Errorf("want 0 defs (no Schema), got %d", len(defs))
	}
}

func TestScanMacroDefs_near_miss_no_run(t *testing.T) {
	stmts := parseStmts(t, `
Debug: {
    Schema: NoConfig
}
`)
	defs := scanMacroDefs(stmts, "macros")
	if len(defs) != 0 {
		t.Errorf("want 0 defs (no run), got %d", len(defs))
	}
}

func TestScanMacroDefs_inline_schema(t *testing.T) {
	stmts := parseStmts(t, `
JSON: {
    Schema #(pretty Bool, indent Int)
    run #(subject Boc, config Schema, Boc) {}
}
`)
	defs := scanMacroDefs(stmts, "macros")
	if len(defs) != 1 {
		t.Fatalf("want 1 def, got %d", len(defs))
	}
	if len(defs[0].SchemaFields) != 2 {
		t.Errorf("schema fields: got %d, want 2", len(defs[0].SchemaFields))
	}
	if defs[0].SchemaFields[0].Name != "pretty" {
		t.Errorf("field[0]: got %q", defs[0].SchemaFields[0].Name)
	}
	if defs[0].SchemaFields[1].Name != "indent" {
		t.Errorf("field[1]: got %q", defs[0].SchemaFields[1].Name)
	}
}

// ---------------------------------------------------------------------------
// annotationTriggers
// ---------------------------------------------------------------------------

func TestAnnotationTriggers_basic(t *testing.T) {
	stmts := parseStmts(t, "`Debug: {}`\nPerson: {}\n")
	// The annotation is on the ShortDecl for Person.
	var ann *ast.Annotation
	for _, n := range stmts {
		if sd, ok := n.(*ast.ShortDecl); ok && sd.Annotation != nil {
			ann = sd.Annotation
		}
	}
	if ann == nil {
		t.Fatal("no annotation found")
	}
	triggers := annotationTriggers(ann)
	if len(triggers) != 1 {
		t.Fatalf("want 1 trigger, got %d", len(triggers))
	}
	if triggers[0].Name != "Debug" {
		t.Errorf("name: got %q", triggers[0].Name)
	}
}

func TestAnnotationTriggers_with_config(t *testing.T) {
	stmts := parseStmts(t, "`JSON: { pretty: true }`\nData: {}\n")
	var ann *ast.Annotation
	for _, n := range stmts {
		if sd, ok := n.(*ast.ShortDecl); ok && sd.Annotation != nil {
			ann = sd.Annotation
		}
	}
	if ann == nil {
		t.Fatal("no annotation found")
	}
	triggers := annotationTriggers(ann)
	if len(triggers) != 1 {
		t.Fatalf("want 1 trigger, got %d", len(triggers))
	}
	if len(triggers[0].Config) != 1 {
		t.Fatalf("config len: got %d, want 1", len(triggers[0].Config))
	}
	if triggers[0].Config[0].Key != "pretty" {
		t.Errorf("config key: got %q", triggers[0].Config[0].Key)
	}
	if b, ok := triggers[0].Config[0].Value.(macrowire.ConfigBool); !ok || !b.V {
		t.Errorf("config value: want true bool")
	}
}

func TestAnnotationTriggers_lowercase_ignored(t *testing.T) {
	stmts := parseStmts(t, "`debug: {}`\nPerson: {}\n")
	var ann *ast.Annotation
	for _, n := range stmts {
		if sd, ok := n.(*ast.ShortDecl); ok && sd.Annotation != nil {
			ann = sd.Annotation
		}
	}
	if ann == nil {
		return
	}
	triggers := annotationTriggers(ann)
	if len(triggers) != 0 {
		t.Errorf("want 0 triggers (lowercase), got %d", len(triggers))
	}
}

// ---------------------------------------------------------------------------
// validateConfig
// ---------------------------------------------------------------------------

func TestValidateConfig_noconfig_empty(t *testing.T) {
	def := &macroDef{Name: "Debug", SchemaTypeName: "NoConfig"}
	tr := trigger{Name: "Debug"}
	if msg := validateConfig(tr, def); msg != "" {
		t.Errorf("want empty, got %q", msg)
	}
}

func TestValidateConfig_noconfig_with_keys(t *testing.T) {
	def := &macroDef{Name: "Debug", SchemaTypeName: "NoConfig"}
	tr := trigger{Name: "Debug", Config: []macrowire.ConfigEntry{
		{Key: "ignore", Value: macrowire.ConfigBool{V: true}},
	}}
	if msg := validateConfig(tr, def); msg == "" {
		t.Error("want error for unexpected config key")
	}
}

func TestValidateConfig_inline_schema_valid(t *testing.T) {
	def := &macroDef{Name: "JSON", SchemaFields: []schemaField{
		{Name: "pretty", Type: "Bool"},
		{Name: "indent", Type: "Int"},
	}}
	tr := trigger{Name: "JSON", Config: []macrowire.ConfigEntry{
		{Key: "pretty", Value: macrowire.ConfigBool{V: true}},
	}}
	if msg := validateConfig(tr, def); msg != "" {
		t.Errorf("want empty, got %q", msg)
	}
}

func TestValidateConfig_inline_schema_bad_key(t *testing.T) {
	def := &macroDef{Name: "JSON", SchemaFields: []schemaField{
		{Name: "pretty", Type: "Bool"},
	}}
	tr := trigger{Name: "JSON", Config: []macrowire.ConfigEntry{
		{Key: "ignor", Value: macrowire.ConfigBool{V: true}},
	}}
	if msg := validateConfig(tr, def); msg == "" {
		t.Error("want error for unknown config key")
	}
}

// ---------------------------------------------------------------------------
// buildSubjectPayload
// ---------------------------------------------------------------------------

func TestBuildSubjectPayload_basic(t *testing.T) {
	stmts := parseStmts(t, `
Person: {
    name String
    age Int
}
`)
	var bl *ast.BocLiteral
	for _, n := range stmts {
		if sd, ok := n.(*ast.ShortDecl); ok && sd.Names[0].Name == "Person" {
			if b, ok2 := sd.Values[0].(*ast.BocLiteral); ok2 {
				bl = b
			}
		}
	}
	if bl == nil {
		t.Fatal("no BocLiteral")
	}
	p := buildSubjectPayload("Person", bl, nil)
	if p.SubjectName != "Person" {
		t.Errorf("name: got %q", p.SubjectName)
	}
	if len(p.Fields) != 2 {
		t.Fatalf("fields: got %d, want 2", len(p.Fields))
	}
	if p.Fields[0].Name != "name" || p.Fields[0].Type != "String" {
		t.Errorf("field[0]: got {%q,%q}", p.Fields[0].Name, p.Fields[0].Type)
	}
	if p.Fields[1].Name != "age" || p.Fields[1].Type != "Int" {
		t.Errorf("field[1]: got {%q,%q}", p.Fields[1].Name, p.Fields[1].Type)
	}
}
