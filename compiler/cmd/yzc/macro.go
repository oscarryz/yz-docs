package main

// Macro system (YZC-0028): scan, registry, trigger extraction, config
// validation, and subject payload building.
//
// A macro is a boc satisfying the Macro structural interface:
//
//	Name : {
//	    Schema : SomeConfigType        // or: Schema #(field Type, ...)
//	    run #(subject Boc, config SomeConfigType, Boc) { ... }
//	}
//
// Detection is a structural AST scan — no sema. Macros are triggered by
// uppercase-keyed entries in an annotation (`Debug: {}`); the entry value is
// the config block, validated against the macro's schema fields.

import (
	"fmt"

	"yz/internal/ast"
	"yz/internal/token"
	"yz/runtime/macrowire"
)

// macroDef describes one detected macro implementation.
type macroDef struct {
	Name           string // macro type name — the dispatch key
	RelDir         string // package dir relative to the source root ("" = root)
	SchemaTypeName string // concrete config type name; "" when schema is the inline #(...) form
	SchemaFields   []macrowire.FieldSpec
	Decl           *ast.ShortDecl
}

type buildState int

const (
	macroUnbuilt buildState = iota
	macroBuilding
	macroBuilt
)

// macroRegistry holds all macro definitions discovered in a project.
type macroRegistry struct {
	byName  map[string]*macroDef
	byDir   map[string][]*macroDef
	binPath map[string]string     // relDir → built executable path
	state   map[string]buildState // relDir → bootstrap state
}

func newMacroRegistry() *macroRegistry {
	return &macroRegistry{
		byName:  map[string]*macroDef{},
		byDir:   map[string][]*macroDef{},
		binPath: map[string]string{},
		state:   map[string]buildState{},
	}
}

func (r *macroRegistry) register(def *macroDef) error {
	if prev, ok := r.byName[def.Name]; ok {
		return fmt.Errorf("macro %q defined in both %q and %q", def.Name, prev.RelDir, def.RelDir)
	}
	r.byName[def.Name] = def
	r.byDir[def.RelDir] = append(r.byDir[def.RelDir], def)
	return nil
}

// ---------------------------------------------------------------------------
// Macro definition scan
// ---------------------------------------------------------------------------

// scanMacroDefs walks top-level stmts of a package and returns the macro
// definitions found. A definition is a ShortDecl with a single uppercase name
// whose boc literal contains a Schema entry and a matching run declaration.
// Schema fields are resolved against the same stmt list (plus builtin prelude
// types); an unresolvable or non-scalar schema is an error.
func scanMacroDefs(stmts []ast.Node, relDir string) ([]*macroDef, error) {
	var defs []*macroDef
	for _, stmt := range stmts {
		sd, ok := stmt.(*ast.ShortDecl)
		if !ok || len(sd.Names) != 1 || len(sd.Values) != 1 {
			continue
		}
		if sd.Names[0].TokType != token.TYPE_IDENT {
			continue
		}
		bl, ok := sd.Values[0].(*ast.BocLiteral)
		if !ok {
			continue
		}
		schemaEntry, runDecl := findMacroShape(bl)
		if schemaEntry == nil || runDecl == nil {
			continue
		}
		def := &macroDef{Name: sd.Names[0].Name, RelDir: relDir, Decl: sd}
		if err := resolveSchema(def, schemaEntry, runDecl, stmts); err != nil {
			return nil, err
		}
		defs = append(defs, def)
	}
	return defs, nil
}

// findMacroShape returns the Schema entry (ShortDecl alias or bodiless
// BocDecl) and the run BocDecl when the boc literal has the Macro shape.
func findMacroShape(bl *ast.BocLiteral) (schemaEntry ast.Node, runDecl *ast.BocDecl) {
	for _, el := range bl.Elements {
		switch e := el.(type) {
		case *ast.ShortDecl:
			if len(e.Names) == 1 && e.Names[0].Name == "Schema" && len(e.Values) == 1 {
				if _, ok := e.Values[0].(*ast.Ident); ok {
					schemaEntry = e
				}
			}
		case *ast.BocDecl:
			switch {
			case e.Name.Name == "Schema" && e.Body == nil:
				schemaEntry = e
			case e.Name.Name == "run" && e.Body != nil && isRunSig(e.Sig):
				runDecl = e
			}
		}
	}
	return schemaEntry, runDecl
}

// isRunSig checks for `#(subject Boc, config X, Boc)`: two labeled inputs —
// the first typed Boc — and an unlabeled Boc return.
func isRunSig(sig *ast.BocTypeExpr) bool {
	if sig == nil || len(sig.Params) != 3 {
		return false
	}
	p0, p1, p2 := sig.Params[0], sig.Params[1], sig.Params[2]
	return p0.Label != "" && simpleTypeName(p0.Type) == "Boc" &&
		p1.Label != "" && simpleTypeName(p1.Type) != "" &&
		p2.Label == "" && simpleTypeName(p2.Type) == "Boc"
}

func simpleTypeName(te ast.TypeExpr) string {
	if st, ok := te.(*ast.SimpleTypeExpr); ok && len(st.TypeArgs) == 0 {
		return st.Name
	}
	return ""
}

// resolveSchema fills def.SchemaTypeName and def.SchemaFields.
//
// Supported schema declarations:
//   - alias:  `Schema : NoConfig` — fields come from the named type's boc
//     literal in the same stmt list (prelude NoConfig = no fields)
//   - inline: `Schema #(pretty Bool, ...)` — fields come from the sig params
//
// run's second param names the concrete config type; `Schema` itself is
// accepted and resolved through the alias/inline declaration.
func resolveSchema(def *macroDef, schemaEntry ast.Node, runDecl *ast.BocDecl, stmts []ast.Node) error {
	configTypeName := simpleTypeName(runDecl.Sig.Params[1].Type)

	switch se := schemaEntry.(type) {
	case *ast.ShortDecl: // alias form
		target := se.Values[0].(*ast.Ident).Name
		if configTypeName != "Schema" && configTypeName != target {
			return fmt.Errorf("macro %q: run config type %q does not match Schema alias %q",
				def.Name, configTypeName, target)
		}
		def.SchemaTypeName = target
		fields, err := schemaFieldsFromType(def.Name, target, stmts)
		if err != nil {
			return err
		}
		def.SchemaFields = fields
	case *ast.BocDecl: // inline form: Schema #(fields...)
		if configTypeName != "Schema" {
			return fmt.Errorf("macro %q: run config type %q does not match inline Schema declaration",
				def.Name, configTypeName)
		}
		fields, err := schemaFieldsFromSig(def.Name, se.Sig)
		if err != nil {
			return err
		}
		def.SchemaFields = fields
	}
	return nil
}

// scalarConfigTypes are the only config field types supported in this slice.
var scalarConfigTypes = map[string]bool{
	"String": true, "Int": true, "Decimal": true, "Bool": true,
}

// schemaFieldsFromType extracts config fields from a named struct boc
// declared in the same stmt list. "NoConfig" is the builtin empty schema.
func schemaFieldsFromType(macroName, typeName string, stmts []ast.Node) ([]macrowire.FieldSpec, error) {
	if typeName == "NoConfig" {
		return nil, nil
	}
	for _, stmt := range stmts {
		sd, ok := stmt.(*ast.ShortDecl)
		if !ok || len(sd.Names) != 1 || sd.Names[0].Name != typeName || len(sd.Values) != 1 {
			continue
		}
		bl, ok := sd.Values[0].(*ast.BocLiteral)
		if !ok {
			continue
		}
		var fields []macrowire.FieldSpec
		for _, el := range bl.Elements {
			td, ok := el.(*ast.TypedDecl)
			if !ok {
				continue
			}
			ft := ast.TypeExprString(td.Type)
			if !scalarConfigTypes[ft] {
				return nil, fmt.Errorf("macro %q: config field %q: only scalar config fields are supported, got %s",
					macroName, td.Name.Name, ft)
			}
			fields = append(fields, macrowire.FieldSpec{Name: td.Name.Name, Type: ft})
		}
		return fields, nil
	}
	return nil, fmt.Errorf("macro %q: schema type %q not found in package", macroName, typeName)
}

// schemaFieldsFromSig extracts config fields from an inline `Schema #(...)`.
func schemaFieldsFromSig(macroName string, sig *ast.BocTypeExpr) ([]macrowire.FieldSpec, error) {
	if sig == nil {
		return nil, nil
	}
	var fields []macrowire.FieldSpec
	for _, p := range sig.Params {
		if p.Label == "" {
			return nil, fmt.Errorf("macro %q: inline Schema params must be labeled", macroName)
		}
		ft := ast.TypeExprString(p.Type)
		if !scalarConfigTypes[ft] {
			return nil, fmt.Errorf("macro %q: config field %q: only scalar config fields are supported, got %s",
				macroName, p.Label, ft)
		}
		fields = append(fields, macrowire.FieldSpec{Name: p.Label, Type: ft})
	}
	return fields, nil
}

// ---------------------------------------------------------------------------
// Trigger extraction
// ---------------------------------------------------------------------------

// trigger is one macro invocation requested by an annotation entry.
type trigger struct {
	Name   string
	Config *ast.BocLiteral
	Pos    ast.Pos
}

// annotationTriggers extracts macro triggers from an annotation: entries with
// a single uppercase key. The value must be a boc literal (the config block).
func annotationTriggers(ann *ast.Annotation) ([]trigger, error) {
	if ann == nil || ann.Body == nil {
		return nil, nil
	}
	var triggers []trigger
	for _, el := range ann.Body.Elements {
		sd, ok := el.(*ast.ShortDecl)
		if !ok || len(sd.Names) != 1 || sd.Names[0].TokType != token.TYPE_IDENT {
			continue
		}
		if len(sd.Values) != 1 {
			continue
		}
		bl, ok := sd.Values[0].(*ast.BocLiteral)
		if !ok {
			return nil, fmt.Errorf("macro trigger %q: value must be a boc literal config block", sd.Names[0].Name)
		}
		triggers = append(triggers, trigger{Name: sd.Names[0].Name, Config: bl, Pos: sd.Position()})
	}
	return triggers, nil
}

// ---------------------------------------------------------------------------
// Config validation and encoding
// ---------------------------------------------------------------------------

// kindName maps a decoded config value kind to its Yz scalar type name.
func kindName(k macrowire.ConfigKind) string {
	switch k {
	case macrowire.ConfigString:
		return "String"
	case macrowire.ConfigInt:
		return "Int"
	case macrowire.ConfigDecimal:
		return "Decimal"
	case macrowire.ConfigBool:
		return "Bool"
	}
	return "?"
}

func schemaFieldList(fields []macrowire.FieldSpec) string {
	if len(fields) == 0 {
		return "no config fields"
	}
	s := "expected: "
	for i, f := range fields {
		if i > 0 {
			s += ", "
		}
		s += f.Name + " " + f.Type
	}
	return s
}

// encodeConfig converts a trigger's config block into ordered wire entries.
func encodeConfig(cfg *ast.BocLiteral) ([]macrowire.ConfigEntry, error) {
	var entries []macrowire.ConfigEntry
	if cfg == nil {
		return nil, nil
	}
	for _, el := range cfg.Elements {
		sd, ok := el.(*ast.ShortDecl)
		if !ok || len(sd.Names) != 1 || len(sd.Values) != 1 {
			return nil, fmt.Errorf("config entries must be key: value pairs")
		}
		v, err := macrowire.ValueFromExpr(sd.Values[0])
		if err != nil {
			return nil, fmt.Errorf("config key %q: %w", sd.Names[0].Name, err)
		}
		entries = append(entries, macrowire.ConfigEntry{Key: sd.Names[0].Name, Value: v})
	}
	return entries, nil
}

// validateConfig checks a trigger's config entries against the macro's
// schema fields: every key must name a schema field with a matching scalar
// type, and every schema field must be present.
func validateConfig(tr trigger, def *macroDef, entries []macrowire.ConfigEntry) error {
	byName := map[string]macrowire.FieldSpec{}
	for _, f := range def.SchemaFields {
		byName[f.Name] = f
	}
	seen := map[string]bool{}
	for _, e := range entries {
		f, ok := byName[e.Key]
		if !ok {
			return fmt.Errorf("macro %q: config key %q does not match any Schema field (%s)",
				def.Name, e.Key, schemaFieldList(def.SchemaFields))
		}
		if got := kindName(e.Value.Kind); got != f.Type {
			return fmt.Errorf("macro %q: config key %q: expected %s, got %s",
				def.Name, e.Key, f.Type, got)
		}
		seen[e.Key] = true
	}
	for _, f := range def.SchemaFields {
		if !seen[f.Name] {
			return fmt.Errorf("macro %q: missing config field %q (%s)",
				def.Name, f.Name, schemaFieldList(def.SchemaFields))
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Subject payload
// ---------------------------------------------------------------------------

// buildSubjectPayload serializes the subject boc's current data fields
// (TypedDecl elements) plus the validated config into a wire payload.
func buildSubjectPayload(name string, bl *ast.BocLiteral, cfg []macrowire.ConfigEntry) *macrowire.Payload {
	p := &macrowire.Payload{SubjectName: name, Config: cfg}
	for _, el := range bl.Elements {
		td, ok := el.(*ast.TypedDecl)
		if !ok {
			continue
		}
		p.Fields = append(p.Fields, macrowire.FieldSpec{
			Name: td.Name.Name,
			Type: ast.TypeExprString(td.Type),
		})
	}
	return p
}
