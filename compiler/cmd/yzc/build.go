package main

import (
	"fmt"
	"io/fs"
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
)

// cmdBuild compiles the Yz project. projectDir owns target/; extraRoots are
// additional source roots contributing FQNs to the same namespace.
func cmdBuild(projectDir string, extraRoots []string) error {
	srcRoots := append([]string{projectDir}, extraRoots...)
	sources, err := compileProject(projectDir, srcRoots)
	if err != nil {
		return err
	}

	genDir := filepath.Join(projectDir, "target", "gen")
	binDir := filepath.Join(projectDir, "target", "bin")

	if err := writeGeneratedGo(genDir, sources, projectDir); err != nil {
		return err
	}

	binPath := filepath.Join(binDir, "app")
	if err := goBuild(genDir, binPath); err != nil {
		return err
	}

	fmt.Printf("yzc(%s): built %s\n", version, binPath)
	return nil
}

// cmdRun compiles and immediately runs the Yz project.
func cmdRun(projectDir string, extraRoots []string) error {
	if err := cmdBuild(projectDir, extraRoots); err != nil {
		return err
	}
	binPath := filepath.Join(projectDir, "target", "bin", "app")
	absPath, err := filepath.Abs(binPath)
	if err != nil {
		return err
	}
	cmd := exec.Command(absPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// fileEntry is one .yz source file discovered during the project walk.
type fileEntry struct {
	absPath string // absolute path to the .yz file
	relDir  string // slash-separated path relative to its source root, "" for root
	name    string // file name without .yz extension
	srcRoot string // absolute path of the source root this file came from
}

// walkYzFiles recursively finds all .yz files under srcRoot, skipping
// target/ and hidden directories.
func walkYzFiles(srcRoot string) ([]fileEntry, error) {
	absSrcRoot, err := filepath.Abs(srcRoot)
	if err != nil {
		return nil, err
	}
	var entries []fileEntry
	err = filepath.WalkDir(absSrcRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			base := d.Name()
			if base == "target" || strings.HasPrefix(base, ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".yz") {
			return nil
		}
		rel, _ := filepath.Rel(absSrcRoot, path)
		relDir := filepath.ToSlash(filepath.Dir(rel))
		if relDir == "." {
			relDir = ""
		}
		name := strings.TrimSuffix(d.Name(), ".yz")
		entries = append(entries, fileEntry{absPath: path, relDir: relDir, name: name, srcRoot: absSrcRoot})
		return nil
	})
	return entries, err
}

// pkgNameFromDir returns the Go package name for a relative directory.
// The root ("") is the main package; subdirs use their last path segment.
func pkgNameFromDir(relDir string) string {
	if relDir == "" {
		return "main"
	}
	parts := strings.Split(relDir, "/")
	return parts[len(parts)-1]
}

// pkgExport holds the exported symbols of one compiled sub-package.
type pkgExport struct {
	relDir     string
	pkgAlias   string
	importPath string
	exports    map[string]*sema.Symbol
}

// compileProject walks all .yz files across all srcRoots, compiles each
// directory as a separate Go package, and returns a map of relative output
// path → Go source. projectDir owns target/; srcRoots includes projectDir
// plus any extra roots (stdlib, third-party, etc.).
// Sub-packages are compiled before the root; their exports are registered in
// the root's analyzer so FQN references (house.front.Host()) resolve correctly.
func compileProject(projectDir string, srcRoots []string) (map[string]string, error) {
	var entries []fileEntry
	for _, root := range srcRoots {
		es, err := walkYzFiles(root)
		if err != nil {
			return nil, err
		}
		entries = append(entries, es...)
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("no .yz files found in %s", strings.Join(srcRoots, ", "))
	}

	// Detect FQN collisions: same (relDir, name) pair from different roots.
	type fqnKey struct{ relDir, name string }
	seen := map[fqnKey]string{} // fqnKey → srcRoot of first occurrence
	for _, e := range entries {
		k := fqnKey{e.relDir, e.name}
		if prev, conflict := seen[k]; conflict {
			fqn := e.name
			if e.relDir != "" {
				fqn = strings.ReplaceAll(e.relDir, "/", ".") + "." + e.name
			}
			return nil, fmt.Errorf("FQN conflict: %q defined in both %s and %s", fqn, prev, e.srcRoot)
		}
		seen[k] = e.srcRoot
	}

	// Group by relative directory.
	byDir := map[string][]fileEntry{}
	for _, e := range entries {
		byDir[e.relDir] = append(byDir[e.relDir], e)
	}

	// Invariant 5 (spec §9): foo.yz + foo/ can coexist. Detect pairs where a
	// root file and a same-named directory share a stem, parse and wrap the
	// directory's files, and inject them into the root boc literal at analysis
	// time. The matched directory is removed from byDir so it is not compiled
	// as a separate Go package.
	// Only applies when foo/ has no sub-directories (deeper nesting deferred).
	pendingInjections := map[string][]ast.Node{} // root file name → nodes to inject
	for _, fe := range byDir[""] {
		subDir := fe.name
		if _, ok := byDir[subDir]; !ok {
			continue
		}
		// Skip if sub-dir has its own sub-directories (not supported yet).
		hasDeeper := false
		for d := range byDir {
			if strings.HasPrefix(d, subDir+"/") {
				hasDeeper = true
				break
			}
		}
		if hasDeeper {
			continue
		}
		for _, sf := range byDir[subDir] {
			src, err := os.ReadFile(sf.absPath)
			if err != nil {
				return nil, fmt.Errorf("reading %s: %w", sf.absPath, err)
			}
			p := parser.New(src)
			parsed, parseErr := p.ParseFile()
			if parseErr != nil {
				if pe, ok := parseErr.(*parser.ParseError); ok {
					fmt.Fprint(os.Stderr, diagnostic.Format(src, sf.absPath, pe.Line, pe.Col, pe.Len, pe.Msg))
					return nil, fmt.Errorf("parse error in %s", sf.absPath)
				}
				return nil, fmt.Errorf("parse %s: %w", sf.absPath, parseErr)
			}
			pendingInjections[fe.name] = append(pendingInjections[fe.name],
				&ast.ShortDecl{
					Names:  []*ast.Ident{{Name: sf.name}},
					Values: []ast.Expr{&ast.BocLiteral{Elements: parsed.Stmts}},
				},
			)
		}
		delete(byDir, subDir)
	}

	// Sort dirs: deepest first so sub-packages are ready before root.
	// Build after Invariant 5 deletion so merged dirs are excluded.
	var dirs []string
	for d := range byDir {
		dirs = append(dirs, d)
	}
	sort.Slice(dirs, func(i, j int) bool {
		di := strings.Count(dirs[i], "/")
		dj := strings.Count(dirs[j], "/")
		if di != dj {
			return di > dj
		}
		return dirs[i] < dirs[j]
	})

	// Compile all non-root directories first; collect their exports.
	var subExports []*pkgExport
	result := map[string]string{}
	for _, dir := range dirs {
		if dir == "" {
			continue // root compiled last
		}
		goSrc, exp, err := compilePackageDir(byDir[dir], dir, nil, nil)
		if err != nil {
			return nil, err
		}
		subExports = append(subExports, exp)
		pkgName := pkgNameFromDir(dir)
		outPath := filepath.Join(filepath.FromSlash(dir), pkgName+".go")
		result[outPath] = goSrc
	}

	// Build a root analyzer pre-seeded with sub-package exports.
	rootAnalyzer := sema.NewAnalyzer()
	for _, exp := range subExports {
		rootAnalyzer.RegisterPackage(exp.relDir, exp.pkgAlias, exp.importPath, exp.exports)
	}

	// Compile the root (main) package with the seeded analyzer.
	if rootFiles, ok := byDir[""]; ok {
		goSrc, _, err := compilePackageDir(rootFiles, "", rootAnalyzer, pendingInjections)
		if err != nil {
			return nil, err
		}
		result["main.go"] = goSrc
	}

	return result, nil
}

// compilePackageDir compiles all .yz files in one directory into a single Go
// source string and returns the exported symbols of the package.
// If a is non-nil it is used as the analyzer (for the root package which has
// pre-registered sub-package exports); otherwise a fresh analyzer is created.
func compilePackageDir(files []fileEntry, relDir string, a *sema.Analyzer, pendingInjections map[string][]ast.Node) (string, *pkgExport, error) {
	pkgName := pkgNameFromDir(relDir)

	// Sort: main.yz last within each dir.
	sort.Slice(files, func(i, j int) bool {
		if files[i].name == "main" {
			return false
		}
		if files[j].name == "main" {
			return true
		}
		return files[i].name < files[j].name
	})

	type parsedFile struct {
		sf   *ast.SourceFile
		path string
		src  []byte
	}
	var pfiles []parsedFile
	for _, fe := range files {
		src, err := os.ReadFile(fe.absPath)
		if err != nil {
			return "", nil, fmt.Errorf("reading %s: %w", fe.absPath, err)
		}
		p := parser.New(src)
		sf, parseErr := p.ParseFile()
		if parseErr != nil {
			if pe, ok := parseErr.(*parser.ParseError); ok {
				fmt.Fprint(os.Stderr, diagnostic.Format(src, fe.absPath, pe.Line, pe.Col, pe.Len, pe.Msg))
				return "", nil, fmt.Errorf("parse error in %s", fe.absPath)
			}
			return "", nil, fmt.Errorf("parse %s: %w", fe.absPath, parseErr)
		}
		// Invariant 1+2 (spec §9): every file's content is the body of a boc
		// named after the file. Always wrap — no special-casing for files that
		// already declare a same-named top-level boc.
		sf.Stmts = []ast.Node{
			&ast.ShortDecl{
				Names:         []*ast.Ident{{Name: fe.name, TokType: token.LookupIdent(fe.name)}},
				Values:        []ast.Expr{&ast.BocLiteral{Elements: sf.Stmts}},
				IsFileWrapper: true,
			},
		}
		// Invariant 5: inject sub-directory declarations into the matching
		// root boc literal (foo.yz + foo/ merge).
		if nodes, ok := pendingInjections[fe.name]; ok {
			injectIntoBocLiteral(sf, fe.name, nodes)
		}
		pfiles = append(pfiles, parsedFile{sf: sf, path: fe.absPath, src: src})
	}

	if a == nil {
		a = sema.NewAnalyzer()
	}
	for _, pf := range pfiles {
		if err := a.AnalyzeFile(pf.sf); err != nil {
			if ses, ok := err.(sema.SemaErrors); ok {
				for _, se := range ses {
					fmt.Fprint(os.Stderr, diagnostic.Format(pf.src, pf.path, se.Line, se.Col, se.Len, se.Msg))
				}
				return "", nil, fmt.Errorf("semantic errors in %s", pf.path)
			}
			return "", nil, fmt.Errorf("sema %s: %w", pf.path, err)
		}
	}

	combined := &ir.File{PkgName: pkgName}
	for _, pf := range pfiles {
		f := ir.Lower(pf.sf, a, pkgName)
		combined.Decls = append(combined.Decls, f.Decls...)
		for _, imp := range f.Imports {
			if !containsStr(combined.Imports, imp) {
				combined.Imports = append(combined.Imports, imp)
			}
		}
	}

	// Build the export record for this package (used by parent packages).
	importPath := "yzapp/" + strings.ReplaceAll(relDir, string(filepath.Separator), "/")
	exp := &pkgExport{
		relDir:     relDir,
		pkgAlias:   pkgName,
		importPath: importPath,
		exports:    a.ExportedSymbols(),
	}

	return codegen.Generate(combined), exp, nil
}

// injectIntoBocLiteral finds the top-level ShortDecl named bocName in sf and
// appends nodes to its BocLiteral body (spec §9 Invariant 5 loader merge).
// When the outer ShortDecl is a file wrapper that itself contains an inner
// ShortDecl with the same name (e.g. utils.yz wrapping `utils: {}`), the
// nodes are injected into the inner boc so they appear as fields/methods of
// the named singleton rather than as siblings of it.
// No-ops if no matching ShortDecl is found.
func injectIntoBocLiteral(sf *ast.SourceFile, bocName string, nodes []ast.Node) {
	for _, stmt := range sf.Stmts {
		sd, ok := stmt.(*ast.ShortDecl)
		if !ok || len(sd.Names) == 0 || len(sd.Values) == 0 {
			continue
		}
		if sd.Names[0].Name != bocName {
			continue
		}
		if bl, ok := sd.Values[0].(*ast.BocLiteral); ok {
			// If this is a file wrapper with an inner same-named boc, inject there.
			if sd.IsFileWrapper {
				for _, inner := range bl.Elements {
					if innerSD, ok2 := inner.(*ast.ShortDecl); ok2 &&
						len(innerSD.Names) > 0 && innerSD.Names[0].Name == bocName &&
						len(innerSD.Values) > 0 {
						if innerBL, ok3 := innerSD.Values[0].(*ast.BocLiteral); ok3 {
							innerBL.Elements = append(innerBL.Elements, nodes...)
							return
						}
					}
				}
			}
			bl.Elements = append(bl.Elements, nodes...)
			return
		}
	}
}


func containsStr(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

// writeGeneratedGo writes all generated Go files and go.mod into genDir.
// sources maps relative output path (e.g. "main.go", "house/front/front.go")
// to Go source content.
func writeGeneratedGo(genDir string, sources map[string]string, projectDir string) error {
	if err := os.MkdirAll(genDir, 0o755); err != nil {
		return fmt.Errorf("creating gen dir: %w", err)
	}

	for relPath, goSrc := range sources {
		fullPath := filepath.Join(genDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			return fmt.Errorf("creating dir for %s: %w", relPath, err)
		}
		if err := os.WriteFile(fullPath, []byte(goSrc), 0o644); err != nil {
			return fmt.Errorf("writing %s: %w", relPath, err)
		}
	}

	// Find the yz compiler module root so generated code can reference yz/runtime/rt.
	yzRoot, err := yzModuleDir()
	if err != nil {
		return fmt.Errorf("locating yz module: %w", err)
	}

	// Write go.mod with a replace directive pointing at the local yz module.
	goMod := fmt.Sprintf(`module yzapp

go 1.23

require yz v0.0.0

replace yz => %s
`, yzRoot)
	modPath := filepath.Join(genDir, "go.mod")
	if err := os.WriteFile(modPath, []byte(goMod), 0o644); err != nil {
		return fmt.Errorf("writing go.mod: %w", err)
	}

	return nil
}

// goBuild runs `go build -o binPath .` in genDir.
func goBuild(genDir, binPath string) error {
	absBin, err := filepath.Abs(binPath)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(absBin), 0o755); err != nil {
		return fmt.Errorf("creating bin dir: %w", err)
	}
	cmd := exec.Command("go", "build", "-o", absBin, ".")
	cmd.Dir = genDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go build: %w", err)
	}
	return nil
}

// yzModuleDir returns the absolute path to the directory containing the yz
// module (the one with `module yz` in its go.mod).
// Search order:
//  1. $YZ_ROOT env var
//  2. Walk up from the current working directory
//  3. Walk up from the executable's directory
func yzModuleDir() (string, error) {
	if root := os.Getenv("YZ_ROOT"); root != "" {
		return filepath.Abs(root)
	}

	// Walk up from cwd.
	cwd, err := os.Getwd()
	if err == nil {
		if found, ok := findGoModWithModule(cwd, "yz"); ok {
			return found, nil
		}
	}

	// Walk up from executable location.
	exe, err := os.Executable()
	if err == nil {
		if found, ok := findGoModWithModule(filepath.Dir(exe), "yz"); ok {
			return found, nil
		}
	}

	return "", fmt.Errorf("cannot locate yz module root (set YZ_ROOT env var)")
}

// findGoModWithModule walks up from dir looking for a go.mod that declares
// the given module name. Returns the directory and true if found.
func findGoModWithModule(dir, moduleName string) (string, bool) {
	for {
		modPath := filepath.Join(dir, "go.mod")
		data, err := os.ReadFile(modPath)
		if err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				line = strings.TrimSpace(line)
				if line == "module "+moduleName {
					return dir, true
				}
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", false
}
