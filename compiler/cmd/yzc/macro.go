package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"yz/internal/ast"
	"yz/internal/codegen"
	"yz/internal/ir"
	"yz/internal/parser"
	"yz/internal/sema"
	"yz/runtime/macrowire"
)

// ---------------------------------------------------------------------------
// Macro definition
// ---------------------------------------------------------------------------

// schemaField is one field declared in a macro's Schema type.
type schemaField struct {
	Name string
	Type string // Yz type expression as string
}

// macroDef records a macro boc found during the scan phase.
type macroDef struct {
	Name            string       // e.g. "Debug"
	RelDir          string       // relative dir of the package containing this macro
	SchemaTypeName  string       // e.g. "NoConfig" for alias form, "" for inline form
	SchemaFields    []schemaField // fields from the inline Schema type (empty for alias)
	Decl            *ast.ShortDecl
}

// ---------------------------------------------------------------------------
// Macro registry
// ---------------------------------------------------------------------------

type macroState int

const (
	macroUnbuilt macroState = iota
	macroBuilding
	macroBuilt
)

// macroRegistry tracks macro definitions and their compiled binaries.
type macroRegistry struct {
	byName map[string]*macroDef  // macro name → def
	byDir  map[string]*macroDef  // relDir → def (one macro pkg per dir)
	binPath map[string]string    // relDir → path to compiled binary
	state  map[string]macroState // relDir → build state
}

func newMacroRegistry() *macroRegistry {
	return &macroRegistry{
		byName:  make(map[string]*macroDef),
		byDir:   make(map[string]*macroDef),
		binPath: make(map[string]string),
		state:   make(map[string]macroState),
	}
}

// register adds a macro definition to the registry.
// Returns an error if the name is already registered.
func (r *macroRegistry) register(def *macroDef) error {
	if existing, ok := r.byName[def.Name]; ok {
		return fmt.Errorf("macro %q already defined in %q, cannot redefine in %q",
			def.Name, existing.RelDir, def.RelDir)
	}
	r.byName[def.Name] = def
	r.byDir[def.RelDir] = def
	r.state[def.RelDir] = macroUnbuilt
	return nil
}

// ---------------------------------------------------------------------------
// Macro scan (structural AST, pre-sema)
// ---------------------------------------------------------------------------

// scanMacroDefs scans top-level stmts for macro definitions.
// A macro is an uppercase-named ShortDecl whose BocLiteral has a Schema
// entry and a run BocDecl with 3 parameters.
func scanMacroDefs(stmts []ast.Node, relDir string) []*macroDef {
	var defs []*macroDef
	for _, node := range stmts {
		sd, ok := node.(*ast.ShortDecl)
		if !ok || len(sd.Names) != 1 || !isUpperName(sd.Names[0].Name) {
			continue
		}
		if len(sd.Values) != 1 {
			continue
		}
		bl, ok := sd.Values[0].(*ast.BocLiteral)
		if !ok {
			continue
		}
		def := extractMacroDef(sd.Names[0].Name, relDir, bl, sd)
		if def != nil {
			defs = append(defs, def)
		}
	}
	return defs
}

// extractMacroDef inspects a BocLiteral to see if it satisfies the Macro
// interface: has a Schema entry and a run boc decl with 3 params.
// Returns nil if the boc does not satisfy the interface.
func extractMacroDef(name, relDir string, bl *ast.BocLiteral, decl *ast.ShortDecl) *macroDef {
	hasSchema := false
	hasRun := false
	var schemaTypeName string
	var schemaFields []schemaField

	for _, elem := range bl.Elements {
		switch e := elem.(type) {
		case *ast.ShortDecl:
			if len(e.Names) == 1 && e.Names[0].Name == "Schema" {
				hasSchema = true
				// Schema : AliasType — alias form
				if len(e.Values) == 1 {
					if id, ok := e.Values[0].(*ast.Ident); ok {
						schemaTypeName = id.Name
					}
				}
			}
			if len(e.Names) == 1 && e.Names[0].Name == "run" {
				if len(e.Values) == 1 {
					if innerBl, ok := e.Values[0].(*ast.BocLiteral); ok {
						if innerBl.Annotation != nil {
							// run with annotation — skip
							_ = innerBl
						}
						hasRun = true // simplified detection
					}
				}
			}
		case *ast.TypedDecl:
			if e.Name.Name == "Schema" {
				hasSchema = true
			}
		case *ast.BocDecl:
			if e.Name.Name == "Schema" {
				hasSchema = true
				// Inline Schema #(fields...) — extract field specs
				if e.Sig != nil {
					for _, p := range e.Sig.Params {
						if p.Label != "" && p.Type != nil {
							schemaFields = append(schemaFields, schemaField{
								Name: p.Label,
								Type: ast.TypeExprString(p.Type),
							})
						}
					}
				}
			}
			if e.Name.Name == "run" && e.Sig != nil {
				// Must have exactly 3 params: subject Boc, config SchemaType, return Boc
				if len(e.Sig.Params) >= 3 {
					hasRun = true
				}
			}
		}
	}

	if !hasSchema || !hasRun {
		return nil
	}
	return &macroDef{
		Name:           name,
		RelDir:         relDir,
		SchemaTypeName: schemaTypeName,
		SchemaFields:   schemaFields,
		Decl:           decl,
	}
}

func isUpperName(name string) bool {
	if name == "" {
		return false
	}
	ch := rune(name[0])
	return ch >= 'A' && ch <= 'Z'
}

// ---------------------------------------------------------------------------
// Trigger detection
// ---------------------------------------------------------------------------

// trigger is one macro invocation extracted from an annotation boc.
type trigger struct {
	Name   string            // macro name (uppercase key in annotation)
	Config []macrowire.ConfigEntry
	Pos    ast.Pos
}

// annotationTriggers extracts macro triggers from an annotation's BocLiteral.
// Triggers are ShortDecls with uppercase names in the annotation body.
func annotationTriggers(ann *ast.Annotation) []trigger {
	if ann == nil || ann.Body == nil {
		return nil
	}
	var result []trigger
	for _, elem := range ann.Body.Elements {
		sd, ok := elem.(*ast.ShortDecl)
		if !ok || len(sd.Names) != 1 {
			continue
		}
		name := sd.Names[0].Name
		if !isUpperName(name) {
			continue
		}
		// Collect config from the trigger value (BocLiteral).
		var config []macrowire.ConfigEntry
		if len(sd.Values) == 1 {
			if bl, ok := sd.Values[0].(*ast.BocLiteral); ok {
				config = extractConfig(bl)
			}
		}
		result = append(result, trigger{Name: name, Config: config, Pos: sd.Pos})
	}
	return result
}

// extractConfig reads ShortDecl scalar entries from a config BocLiteral.
func extractConfig(bl *ast.BocLiteral) []macrowire.ConfigEntry {
	var entries []macrowire.ConfigEntry
	for _, elem := range bl.Elements {
		sd, ok := elem.(*ast.ShortDecl)
		if !ok || len(sd.Names) != 1 || len(sd.Values) != 1 {
			continue
		}
		key := sd.Names[0].Name
		val := configValueFromExpr(sd.Values[0])
		if val == nil {
			continue
		}
		entries = append(entries, macrowire.ConfigEntry{Key: key, Value: val})
	}
	return entries
}

func configValueFromExpr(e ast.Expr) macrowire.ConfigValue {
	switch v := e.(type) {
	case *ast.StringLit:
		// Strip quotes
		s := v.Value
		if len(s) >= 2 {
			s = s[1 : len(s)-1]
		}
		return macrowire.ConfigString{V: s}
	case *ast.IntLit:
		var n int64
		fmt.Sscanf(v.Value, "%d", &n)
		return macrowire.ConfigInt{V: n}
	case *ast.Ident:
		switch v.Name {
		case "true":
			return macrowire.ConfigBool{V: true}
		case "false":
			return macrowire.ConfigBool{V: false}
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Config validation
// ---------------------------------------------------------------------------

// validateConfig checks that all config keys in the trigger match the macro's
// schema fields (for non-NoConfig macros). Returns an error string or "".
func validateConfig(t trigger, def *macroDef) string {
	if def.SchemaTypeName != "" {
		// Alias form (e.g. Schema : NoConfig) — no fields expected
		if len(def.SchemaFields) == 0 && len(t.Config) > 0 {
			// Any keys in config are unexpected
			for _, e := range t.Config {
				return fmt.Sprintf("macro %q: config key %q does not match any Schema field (expected: no fields)", t.Name, e.Key)
			}
		}
		return ""
	}
	// Inline Schema form: validate each key
	allowed := make(map[string]bool)
	for _, f := range def.SchemaFields {
		allowed[f.Name] = true
	}
	for _, e := range t.Config {
		if !allowed[e.Key] {
			var fieldDesc []string
			for _, f := range def.SchemaFields {
				fieldDesc = append(fieldDesc, f.Name+" "+f.Type)
			}
			return fmt.Sprintf("macro %q: config key %q does not match any Schema field (expected: %s)",
				t.Name, e.Key, strings.Join(fieldDesc, ", "))
		}
	}
	return ""
}

// ---------------------------------------------------------------------------
// Subject payload builder
// ---------------------------------------------------------------------------

// buildSubjectPayload constructs a macrowire.Payload from a subject ShortDecl
// whose value is a BocLiteral. Only TypedDecl fields are serialized.
func buildSubjectPayload(name string, bl *ast.BocLiteral, config []macrowire.ConfigEntry) *macrowire.Payload {
	p := &macrowire.Payload{
		SubjectName: name,
		Config:      config,
	}
	for _, elem := range bl.Elements {
		switch e := elem.(type) {
		case *ast.TypedDecl:
			p.Fields = append(p.Fields, macrowire.FieldSpec{
				Name: e.Name.Name,
				Type: ast.TypeExprString(e.Type),
			})
		}
	}
	return p
}

// ---------------------------------------------------------------------------
// Source hash for caching
// ---------------------------------------------------------------------------

// sourceHash returns a hex sha256 over the sorted source contents + the
// macro prelude + yzc version. Used to skip rebuilding unchanged macro pkgs.
func sourceHash(sources [][]byte) string {
	h := sha256.New()
	h.Write([]byte(version))
	h.Write([]byte(macroPrelude))
	for _, s := range sources {
		h.Write(s)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

// ---------------------------------------------------------------------------
// Macro prelude (compiler-injected into macro packages)
// ---------------------------------------------------------------------------

// macroPrelude is prepended to every macro package compilation so that macro
// authors can reference Boc, Field, NoConfig, and the generated() helper.
const macroPrelude = `
Boc: {
    name String
    fields [Field]
    source String
}
Field: {
    name String
    type_ String
}
NoConfig: {}
generated #(src String, Boc) {
    b: Boc(name: "", fields: [Field](), source: src)
    b
}
`

// preludeStmts returns the parsed AST nodes of the macro prelude (memoized).
var _preludeStmts []ast.Node
var _preludeErr error
var _preludeOnce bool

func preludeStmts() ([]ast.Node, error) {
	if _preludeOnce {
		return _preludeStmts, _preludeErr
	}
	_preludeOnce = true
	sf, err := parseSource([]byte(macroPrelude))
	if err != nil {
		_preludeErr = fmt.Errorf("macro prelude parse: %w", err)
		return nil, _preludeErr
	}
	_preludeStmts = sf.Stmts
	return _preludeStmts, nil
}

// ---------------------------------------------------------------------------
// Main synthesizer (macro executable entry point)
// ---------------------------------------------------------------------------

// genMacroMain generates the Go `main` function for a macro executable.
// It dispatches on os.Args[1] (the macro name), decodes stdin via macrowire,
// constructs the Boc and config values, calls the macro's Run method,
// and prints the Source of the returned Boc.
func genMacroMain(defs []*macroDef) string {
	var b strings.Builder
	b.WriteString(`package main

import (
	"fmt"
	"io"
	"os"

	std "yz/runtime/rt"
	"yz/runtime/macrowire"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: macros <MacroName>")
		os.Exit(1)
	}
	macroName := os.Args[1]
	stdin, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "read stdin:", err)
		os.Exit(1)
	}
	payload, err := macrowire.DecodePayload(stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "decode payload:", err)
		os.Exit(1)
	}
	_ = payload
	switch macroName {
`)
	for _, def := range defs {
		b.WriteString(fmt.Sprintf("\tcase %q:\n", def.Name))
		// The macro struct constructor takes the Schema field as its only arg.
		// Schema is always an interface type, so nil is a valid zero value.
		b.WriteString(fmt.Sprintf("\t\tm := New%s(nil)\n", def.Name))
		b.WriteString(fmt.Sprintf("\t\tsubject := buildBoc(payload)\n"))
		// config is the Schema interface type — nil is valid for NoConfig.
		b.WriteString("\t\tout := m.Run(subject, nil)\n")
		// Run returns *std.Thunk[*Boc]; force it then read the source field.
		b.WriteString("\t\tfmt.Print(out.Force().Source().GoString())\n")
	}
	b.WriteString(`	default:
		fmt.Fprintf(os.Stderr, "unknown macro %q\n", macroName)
		os.Exit(1)
	}
}

func buildBoc(p *macrowire.Payload) *Boc {
	var fields []*Field
	for _, f := range p.Fields {
		fields = append(fields, NewField(
			std.NewString(f.Name),
			std.NewString(f.Type),
		))
	}
	return NewBoc(std.NewString(p.SubjectName), std.NewArray[*Field](fields...), std.NewString(""))
}
`)
	return b.String()
}

// ---------------------------------------------------------------------------
// Cache helpers
// ---------------------------------------------------------------------------

// macroCacheDir returns target/macros/<pkgKey> inside projectDir.
func macroCacheDir(projectDir, relDir string) string {
	key := strings.ReplaceAll(relDir, "/", "_")
	if key == "" {
		key = "_root"
	}
	return filepath.Join(projectDir, "target", "macros", key)
}

// macroBinPath returns the path to the compiled macro binary.
func macroBinPath(projectDir, relDir string) string {
	return filepath.Join(macroCacheDir(projectDir, relDir), "bin", "macros")
}

// runCachePath returns the path to a cached macro output file.
func runCachePath(projectDir, relDir, macroName string, payload *macrowire.Payload) string {
	h := sha256.New()
	h.Write([]byte(macroName))
	h.Write([]byte(payload.Encode()))
	key := fmt.Sprintf("%x", h.Sum(nil))
	return filepath.Join(macroCacheDir(projectDir, relDir), "runs", key+".out")
}

// readCache reads a cache file if it exists.
func readCache(path string) (string, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	return string(data), true
}

// writeCache writes data to a cache file, creating parent dirs as needed.
func writeCache(path, data string) {
	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	_ = os.WriteFile(path, []byte(data), 0o644)
}

// parseSource parses Yz source bytes using the standard parser.
func parseSource(src []byte) (*ast.SourceFile, error) {
	return parser.New(src).ParseFile()
}

// ---------------------------------------------------------------------------
// Phase 4: Macro expansion
// ---------------------------------------------------------------------------

// expandMacros walks top-level stmts looking for annotated uppercase ShortDecls
// that carry macro triggers. For each trigger it invokes the macro binary and
// appends the returned boc body to the subject BocLiteral.
// reg may be nil (no macros registered) — in that case expandMacros is a no-op.
func expandMacros(stmts []ast.Node, relDir string, reg *macroRegistry, stack []string, projectDir string) error {
	if reg == nil {
		return nil
	}
	for _, node := range stmts {
		sd, ok := node.(*ast.ShortDecl)
		if !ok || len(sd.Names) != 1 {
			continue
		}
		// File wrappers contain the actual declarations inside their BocLiteral.
		// Walk into them transparently so annotations on top-level bocs are found.
		if sd.IsFileWrapper {
			if len(sd.Values) == 1 {
				if bl, ok2 := sd.Values[0].(*ast.BocLiteral); ok2 {
					if err := expandMacros(bl.Elements, relDir, reg, stack, projectDir); err != nil {
						return err
					}
				}
			}
			continue
		}
		if sd.Annotation == nil {
			continue
		}
		name := sd.Names[0].Name
		if !isUpperName(name) {
			continue
		}
		if len(sd.Values) != 1 {
			continue
		}
		bl, ok := sd.Values[0].(*ast.BocLiteral)
		if !ok {
			continue
		}

		triggers := annotationTriggers(sd.Annotation)
		for _, tr := range triggers {
			def, exists := reg.byName[tr.Name]
			if !exists {
				return fmt.Errorf("unknown macro %q: no Macro implementation with that name found", tr.Name)
			}
			// Same-package check
			if def.RelDir == relDir {
				return fmt.Errorf("macro %q is defined in package %q and cannot be applied to a boc in the same package; macros must live in a separate package from the bocs they process", tr.Name, relDir)
			}
			// Config validation
			if msg := validateConfig(tr, def); msg != "" {
				return fmt.Errorf("%s", msg)
			}
			// Ensure macro binary is built
			if err := ensureMacroPackageBuilt(projectDir, def.RelDir, reg, stack); err != nil {
				return err
			}
			// Build payload from current subject elements (progressive merge)
			payload := buildSubjectPayload(name, bl, tr.Config)
			// Invoke macro (with run cache)
			generated, err := invokeMacro(projectDir, def.RelDir, tr.Name, payload, reg)
			if err != nil {
				return fmt.Errorf("macro %q failed: %w", tr.Name, err)
			}
			// Parse returned Yz source as boc body
			sf, parseErr := parseSource([]byte(generated))
			if parseErr != nil {
				return fmt.Errorf("macro %q returned invalid Yz source: %v", tr.Name, parseErr)
			}
			// Append generated stmts to subject BocLiteral
			bl.Elements = append(bl.Elements, sf.Stmts...)
		}
	}
	return nil
}

// invokeMacro runs the macro binary with the given payload on stdin.
// Output is cached keyed by (macroName + payload hash).
func invokeMacro(projectDir, relDir, macroName string, payload *macrowire.Payload, reg *macroRegistry) (string, error) {
	cachePath := runCachePath(projectDir, relDir, macroName, payload)
	if cached, ok := readCache(cachePath); ok {
		return cached, nil
	}

	binPath, ok := reg.binPath[relDir]
	if !ok {
		return "", fmt.Errorf("macro binary not found for %q", relDir)
	}

	cmd := exec.Command(binPath, macroName)
	cmd.Stdin = strings.NewReader(payload.Encode())
	out, err := cmd.Output()
	if err != nil {
		var ee *exec.ExitError
		if ok2 := isExitError(err, &ee); ok2 {
			return "", fmt.Errorf("exit 1:\n%s", string(ee.Stderr))
		}
		return "", err
	}
	result := string(out)
	writeCache(cachePath, result)
	return result, nil
}

func isExitError(err error, target **exec.ExitError) bool {
	if ee, ok := err.(*exec.ExitError); ok {
		*target = ee
		return true
	}
	return false
}

// ---------------------------------------------------------------------------
// Phase 3: Bootstrap build
// ---------------------------------------------------------------------------

// ensureMacroPackageBuilt builds the macro binary for relDir if not already
// built (cached). Detects cycles via the visiting stack.
func ensureMacroPackageBuilt(projectDir, relDir string, reg *macroRegistry, stack []string) error {
	switch reg.state[relDir] {
	case macroBuilt:
		return nil
	case macroBuilding:
		// Cycle detected
		return fmt.Errorf("macro cycle detected: %s", strings.Join(append(stack, relDir), " -> "))
	}

	reg.state[relDir] = macroBuilding
	stack = append(stack, relDir)

	absRelDir := filepath.Join(projectDir, filepath.FromSlash(relDir))
	files, err := walkYzFiles(absRelDir)
	if err != nil {
		return fmt.Errorf("walking macro dir %s: %w", relDir, err)
	}
	// Sort files for deterministic hash
	sort.Slice(files, func(i, j int) bool { return files[i].absPath < files[j].absPath })

	var sources [][]byte
	for _, fe := range files {
		src, err := os.ReadFile(fe.absPath)
		if err != nil {
			return fmt.Errorf("reading %s: %w", fe.absPath, err)
		}
		sources = append(sources, src)
	}

	// Check cache
	cacheDir := macroCacheDir(projectDir, relDir)
	hashFile := filepath.Join(cacheDir, "source.hash")
	binPath := macroBinPath(projectDir, relDir)
	hash := sourceHash(sources)

	if existingHash, err := os.ReadFile(hashFile); err == nil {
		if string(existingHash) == hash {
			if _, err := os.Stat(binPath); err == nil {
				reg.binPath[relDir] = binPath
				reg.state[relDir] = macroBuilt
				return nil
			}
		}
	}

	if err := bootstrapMacroPackage(projectDir, relDir, files, sources, hash, reg, stack); err != nil {
		reg.state[relDir] = macroUnbuilt
		return err
	}

	reg.binPath[relDir] = binPath
	reg.state[relDir] = macroBuilt
	return nil
}

// bootstrapMacroPackage compiles a macro package to a native executable.
// The macro package source files are compiled without file-wrapping (stmts
// concatenated directly) and with the macro prelude prepended. The generated
// Go main is appended, then compiled with go build.
func bootstrapMacroPackage(projectDir, relDir string, files []fileEntry, sources [][]byte, hash string, reg *macroRegistry, stack []string) error {
	prelStmts, err := preludeStmts()
	if err != nil {
		return err
	}

	// Concatenate all parsed stmts (unwrapped — no file wrapper).
	combined := &ast.SourceFile{}
	combined.Stmts = append(combined.Stmts, prelStmts...)

	var defs []*macroDef
	for i, fe := range files {
		sf, parseErr := parseSource(sources[i])
		if parseErr != nil {
			return fmt.Errorf("parse %s: %w", fe.absPath, parseErr)
		}
		// Reject top-level `main` boc in macro packages
		for _, node := range sf.Stmts {
			if sd, ok := node.(*ast.ShortDecl); ok && len(sd.Names) == 1 && sd.Names[0].Name == "main" {
				return fmt.Errorf("macro package %q must not define a 'main' boc", relDir)
			}
		}
		combined.Stmts = append(combined.Stmts, sf.Stmts...)

		// Scan for macro defs (recursive macro expansion on macro code is a
		// documented follow-up; for now, just collect defs for genMacroMain).
		defs = append(defs, scanMacroDefs(sf.Stmts, relDir)...)
	}

	if len(defs) == 0 {
		return fmt.Errorf("macro package %q defines no macros", relDir)
	}

	// Sema + lower + codegen for the macro package.
	a := sema.NewAnalyzer()
	if err := a.AnalyzeFile(combined); err != nil {
		return fmt.Errorf("sema macro package %q: %w", relDir, err)
	}
	irFile := ir.Lower(combined, a, "main")
	goSrc := codegen.Generate(irFile)

	// Generate the Go main dispatcher.
	mainSrc := genMacroMain(defs)

	// Write generated Go files.
	genDir := filepath.Join(macroCacheDir(projectDir, relDir), "gen")
	// Remove stale .go files before writing new ones so old codegen outputs
	// from previous builds don't conflict with current declarations.
	if entries, _ := os.ReadDir(genDir); entries != nil {
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".go") {
				_ = os.Remove(filepath.Join(genDir, e.Name()))
			}
		}
	}
	if err := os.MkdirAll(genDir, 0o755); err != nil {
		return fmt.Errorf("creating macro gen dir: %w", err)
	}
	if err := os.WriteFile(filepath.Join(genDir, "macro.go"), []byte(goSrc), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(genDir, "main.go"), []byte(mainSrc), 0o644); err != nil {
		return err
	}

	// Write go.mod for the macro executable.
	yzRoot, err := yzModuleDir()
	if err != nil {
		return err
	}
	goMod := fmt.Sprintf("module yzapp\n\ngo 1.23\n\nrequire yz v0.0.0\n\nreplace yz => %s\n", yzRoot)
	if err := os.WriteFile(filepath.Join(genDir, "go.mod"), []byte(goMod), 0o644); err != nil {
		return err
	}

	// Build the macro binary.
	binPath := macroBinPath(projectDir, relDir)
	if err := os.MkdirAll(filepath.Dir(binPath), 0o755); err != nil {
		return err
	}
	absBin, _ := filepath.Abs(binPath)
	cmd := exec.Command("go", "build", "-o", absBin, ".")
	cmd.Dir = genDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("building macro %q: %w", relDir, err)
	}

	// Write hash file.
	cacheDir := macroCacheDir(projectDir, relDir)
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(cacheDir, "source.hash"), []byte(hash), 0o644)
}
