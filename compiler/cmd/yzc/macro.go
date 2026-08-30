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
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"yz/internal/ast"
	"yz/internal/codegen"
	"yz/internal/diagnostic"
	"yz/internal/ir"
	"yz/internal/parser"
	"yz/internal/sema"
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
	files   map[string][]fileEntry // relDir → macro package source files
	binPath map[string]string      // relDir → built executable path
	state   map[string]buildState  // relDir → bootstrap state
}

func newMacroRegistry() *macroRegistry {
	return &macroRegistry{
		byName:  map[string]*macroDef{},
		byDir:   map[string][]*macroDef{},
		files:   map[string][]fileEntry{},
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
		// The abstract associated-type form has no concrete constructor for
		// the synthesized macro main. Deferred — use a named config type.
		return fmt.Errorf("macro %q: inline `Schema #(...)` is not yet supported; declare a named config type and alias it (Schema : MyConfig)",
			def.Name)
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

// ---------------------------------------------------------------------------
// Macro prelude
// ---------------------------------------------------------------------------

// macroPrelude is Yz source prepended to every macro package during
// bootstrap compilation. It declares the reflection types a macro's run
// method works with, in the annotated-field form documented in spec 12.5.
const macroPrelude = `
Field : {
    name String
    type String
}

Boc : {
    name String
    fields [Field]
    source String
}

NoConfig : {
    none #(Bool) { true }
}

generated #(src String, Boc) {
    Boc(name: "generated", fields: [Field](), source: src)
}

while #(cond #(Bool), body #()) {
    cond() ? { body(), while(cond, body) }, {}
}
`

// preludeStmts parses the macro prelude. Parsed fresh per call — callers
// splice the returned nodes into their own SourceFile, so sharing AST nodes
// across analyses would leak sema state.
func preludeStmts() ([]ast.Node, error) {
	sf, err := parser.New([]byte(macroPrelude)).ParseFile()
	if err != nil {
		return nil, fmt.Errorf("internal: macro prelude does not parse: %w", err)
	}
	return sf.Stmts, nil
}

// ---------------------------------------------------------------------------
// Synthesized macro main
// ---------------------------------------------------------------------------

// genMacroMain generates the Go main for a macro package executable. The
// executable reads a wire payload on stdin, dispatches on os.Args[1] to the
// requested macro, and writes the returned Boc's source to stdout.
func genMacroMain(defs []*macroDef) string {
	var sb strings.Builder
	sb.WriteString(`package main

import (
	"fmt"
	"io"
	"os"

	mw "yz/runtime/macrowire"
	std "yz/runtime/rt"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: macros <MacroName>")
		os.Exit(2)
	}
	in, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "macro: reading stdin: %v\n", err)
		os.Exit(1)
	}
	p, err := mw.DecodePayload(in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "macro: bad payload: %v\n", err)
		os.Exit(1)
	}
	var fs []*Field
	for _, f := range p.Fields {
		fs = append(fs, NewField(std.NewString(f.Name), std.NewString(f.Type)))
	}
	subject := NewBoc(std.NewString(p.SubjectName), std.NewArray(fs...), std.NewString(""))
	cfg := map[string]mw.ConfigValue{}
	for _, e := range p.Config {
		cfg[e.Key] = e.Value
	}
	_ = cfg

	var out *Boc
	switch os.Args[1] {
`)
	for _, def := range defs {
		sb.WriteString(fmt.Sprintf("\tcase %q:\n", def.Name))
		sb.WriteString(fmt.Sprintf("\t\tout = New%s().Run(subject, %s).Force()\n",
			def.Name, configCtorExpr(def)))
	}
	sb.WriteString(`	default:
		fmt.Fprintf(os.Stderr, "macro: unknown macro %q\n", os.Args[1])
		os.Exit(2)
	}
	os.Stdout.WriteString(out.Source().GoString())
}
`)
	return sb.String()
}

// configCtorExpr returns the Go expression constructing a macro's config
// value from the decoded payload entries.
func configCtorExpr(def *macroDef) string {
	if len(def.SchemaFields) == 0 {
		return "New" + def.SchemaTypeName + "()"
	}
	args := make([]string, len(def.SchemaFields))
	for i, f := range def.SchemaFields {
		key := fmt.Sprintf("%q", f.Name)
		switch f.Type {
		case "String":
			args[i] = "std.NewString(cfg[" + key + "].Str)"
		case "Int":
			args[i] = "std.NewInt(cfg[" + key + "].Int)"
		case "Decimal":
			args[i] = "std.NewDecimal(cfg[" + key + "].Dec)"
		case "Bool":
			args[i] = "std.NewBool(cfg[" + key + "].Bool)"
		}
	}
	return "New" + def.SchemaTypeName + "(" + strings.Join(args, ", ") + ")"
}

// ---------------------------------------------------------------------------
// Bootstrap build (Phase 1 of the two-phase build)
// ---------------------------------------------------------------------------

// macroPkgKey converts a package relDir into a directory-name-safe key.
func macroPkgKey(relDir string) string {
	return strings.ReplaceAll(relDir, "/", "_")
}

// macroSourceHash fingerprints a macro package: its sources, the prelude,
// and the compiler version.
func macroSourceHash(files []fileEntry) (string, error) {
	h := sha256.New()
	sorted := append([]fileEntry(nil), files...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].absPath < sorted[j].absPath })
	for _, fe := range sorted {
		src, err := os.ReadFile(fe.absPath)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(h, "%s\x00", fe.name)
		h.Write(src)
		h.Write([]byte{0})
	}
	h.Write([]byte(macroPrelude))
	h.Write([]byte(version))
	return hex.EncodeToString(h.Sum(nil)), nil
}

// bootstrapMacroPackage compiles a macro package to a native executable at
// target/macros/<pkgKey>/bin/macros, reusing a cached binary when the source
// hash matches. The package compiles unwrapped (statements concatenated, no
// file-wrapper bocs) with the macro prelude prepended, so Boc/Field/NoConfig
// resolve and macro types stay top-level Go types. stack carries the macro
// expansion chain for cycle detection (macros can themselves be annotated
// with macros from other packages).
func bootstrapMacroPackage(projectDir string, files []fileEntry, relDir string, reg *macroRegistry, stack []string) error {
	pkgKey := macroPkgKey(relDir)
	macroDir := filepath.Join(projectDir, "target", "macros", pkgKey)
	binPath := filepath.Join(macroDir, "bin", "macros")
	hashPath := filepath.Join(macroDir, "source.hash")

	hash, err := macroSourceHash(files)
	if err != nil {
		return err
	}
	if prev, err := os.ReadFile(hashPath); err == nil && string(prev) == hash {
		if _, err := os.Stat(binPath); err == nil {
			reg.binPath[relDir] = binPath
			return nil
		}
	}

	// Parse all package files (sorted for determinism) and concatenate.
	sorted := append([]fileEntry(nil), files...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].name < sorted[j].name })
	prelude, err := preludeStmts()
	if err != nil {
		return err
	}
	combinedSF := &ast.SourceFile{Stmts: prelude}
	type parsedSrc struct {
		path string
		src  []byte
	}
	var srcs []parsedSrc
	for _, fe := range sorted {
		src, err := os.ReadFile(fe.absPath)
		if err != nil {
			return fmt.Errorf("reading %s: %w", fe.absPath, err)
		}
		sf, parseErr := parser.New(src).ParseFile()
		if parseErr != nil {
			if pe, ok := parseErr.(*parser.ParseError); ok {
				fmt.Fprint(os.Stderr, diagnostic.Format(src, fe.absPath, pe.Line, pe.Col, pe.Len, pe.Msg))
				return fmt.Errorf("parse error in %s", fe.absPath)
			}
			return fmt.Errorf("parse %s: %w", fe.absPath, parseErr)
		}
		for _, stmt := range sf.Stmts {
			if sd, ok := stmt.(*ast.ShortDecl); ok && len(sd.Names) == 1 && sd.Names[0].Name == "main" {
				return fmt.Errorf("macro package %q must not declare a main boc", relDir)
			}
		}
		combinedSF.Stmts = append(combinedSF.Stmts, sf.Stmts...)
		srcs = append(srcs, parsedSrc{path: fe.absPath, src: src})
	}

	// Macros are regular bocs: they may themselves carry macro annotations
	// (from other packages). Expand them before analysis; the stack detects
	// mutually-triggering cycles.
	if err := expandMacros(combinedSF.Stmts, relDir, reg, stack, projectDir); err != nil {
		return err
	}

	a := sema.NewAnalyzer()
	if err := a.AnalyzeFile(combinedSF); err != nil {
		if ses, ok := err.(sema.SemaErrors); ok && len(srcs) > 0 {
			// Positions may point into any file (or the prelude); report
			// against the first file for context.
			for _, se := range ses {
				fmt.Fprint(os.Stderr, diagnostic.Format(srcs[0].src, srcs[0].path, se.Line, se.Col, se.Len, se.Msg))
			}
			return fmt.Errorf("semantic errors in macro package %q", relDir)
		}
		return fmt.Errorf("macro package %q: %w", relDir, err)
	}

	f := ir.Lower(combinedSF, a, "main")
	goSrc := codegen.Generate(f)

	defs := reg.byDir[relDir]
	sources := map[string]string{
		"macros_gen.go": goSrc,
		"macro_main.go": genMacroMain(defs),
	}
	genDir := filepath.Join(macroDir, "gen")
	if err := writeGeneratedGo(genDir, sources, projectDir); err != nil {
		return err
	}
	if err := goBuild(genDir, binPath); err != nil {
		return fmt.Errorf("building macro package %q: %w", relDir, err)
	}
	if err := os.WriteFile(hashPath, []byte(hash), 0o644); err != nil {
		return err
	}
	reg.binPath[relDir] = binPath
	return nil
}

// ---------------------------------------------------------------------------
// Expansion (Phase 2 of the two-phase build)
// ---------------------------------------------------------------------------

// ensureMacroPackageBuilt lazily bootstraps the macro package that defines a
// triggered macro. stack carries the chain of "<dir>.<Macro>" frames being
// expanded; re-entering a package that is already building is a cycle.
func ensureMacroPackageBuilt(projectDir string, def *macroDef, reg *macroRegistry, stack []string) error {
	switch reg.state[def.RelDir] {
	case macroBuilt:
		return nil
	case macroBuilding:
		chain := append(append([]string(nil), stack...), def.RelDir+"."+def.Name)
		return fmt.Errorf("macro cycle detected: %s", strings.Join(chain, " -> "))
	}
	reg.state[def.RelDir] = macroBuilding
	err := bootstrapMacroPackage(projectDir, reg.files[def.RelDir], def.RelDir, reg,
		append(stack, def.RelDir+"."+def.Name))
	if err != nil {
		return err
	}
	reg.state[def.RelDir] = macroBuilt
	return nil
}

// invokeMacro runs a macro executable with the payload on stdin and returns
// its stdout (the generated Yz source). Results are cached per payload under
// target/macros/<pkgKey>/runs/ — the payload is the input boc's structure,
// so an unchanged subject reuses the previous output.
func invokeMacro(projectDir string, def *macroDef, reg *macroRegistry, payload string) (string, error) {
	runKey := sha256.Sum256([]byte(def.Name + "\x00" + payload))
	runsDir := filepath.Join(projectDir, "target", "macros", macroPkgKey(def.RelDir), "runs")
	cachePath := filepath.Join(runsDir, hex.EncodeToString(runKey[:])+".out")
	if out, err := os.ReadFile(cachePath); err == nil {
		return string(out), nil
	}

	cmd := exec.Command(reg.binPath[def.RelDir], def.Name)
	cmd.Stdin = strings.NewReader(payload)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("macro %q failed (%v):\n%s", def.Name, err, stderr.String())
	}
	if err := os.MkdirAll(runsDir, 0o755); err == nil {
		_ = os.WriteFile(cachePath, stdout.Bytes(), 0o644)
	}
	return stdout.String(), nil
}

// expandMacros walks top-level statements (descending into file-wrapper
// bocs) looking for annotated uppercase declarations, and runs each
// triggered macro in annotation order — merging the generated slots into
// the subject boc literal before sema sees it. Each macro receives the
// progressively merged subject.
func expandMacros(stmts []ast.Node, relDir string, reg *macroRegistry, stack []string, projectDir string) error {
	if reg == nil || len(reg.byName) == 0 {
		return nil
	}
	for _, stmt := range stmts {
		sd, ok := stmt.(*ast.ShortDecl)
		if !ok {
			continue
		}
		if sd.IsFileWrapper {
			if bl, ok := sd.Values[0].(*ast.BocLiteral); ok {
				if err := expandMacros(bl.Elements, relDir, reg, stack, projectDir); err != nil {
					return err
				}
			}
			continue
		}
		if err := expandSubject(sd, relDir, reg, stack, projectDir); err != nil {
			return err
		}
	}
	return nil
}

// expandSubject runs all macros triggered by one annotated declaration.
func expandSubject(sd *ast.ShortDecl, relDir string, reg *macroRegistry, stack []string, projectDir string) error {
	if sd.Annotation == nil || len(sd.Names) != 1 || sd.Names[0].TokType != token.TYPE_IDENT || len(sd.Values) != 1 {
		return nil
	}
	subjectBl, ok := sd.Values[0].(*ast.BocLiteral)
	if !ok {
		return nil
	}
	triggers, err := annotationTriggers(sd.Annotation)
	if err != nil {
		return err
	}
	for _, tr := range triggers {
		def, ok := reg.byName[tr.Name]
		if !ok {
			return fmt.Errorf("unknown macro %q: no Macro implementation with that name found", tr.Name)
		}
		if def.RelDir == relDir {
			return fmt.Errorf("macro %q is defined in package %q and cannot be applied to a boc in the same package; macros must live in a separate package from the bocs they process",
				def.Name, def.RelDir)
		}
		if err := ensureMacroPackageBuilt(projectDir, def, reg, stack); err != nil {
			return err
		}
		entries, err := encodeConfig(tr.Config)
		if err != nil {
			return fmt.Errorf("macro %q: %w", def.Name, err)
		}
		if err := validateConfig(tr, def, entries); err != nil {
			return err
		}
		payload := buildSubjectPayload(sd.Names[0].Name, subjectBl, entries).Encode()
		out, err := invokeMacro(projectDir, def, reg, payload)
		if err != nil {
			return err
		}
		genSF, parseErr := parser.New([]byte(out)).ParseFile()
		if parseErr != nil {
			return fmt.Errorf("macro %q returned invalid Yz source: %v\noutput:\n%s", def.Name, parseErr, out)
		}
		subjectBl.Elements = append(subjectBl.Elements, genSF.Stmts...)
	}
	return nil
}
