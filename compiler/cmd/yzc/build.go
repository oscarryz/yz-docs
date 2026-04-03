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
	"yz/internal/ir"
	"yz/internal/parser"
	"yz/internal/sema"
)

// cmdBuild compiles the Yz project in dir to a binary at target/bin/app.
func cmdBuild(dir string) error {
	sources, err := compileProject(dir)
	if err != nil {
		return err
	}

	genDir := filepath.Join(dir, "target", "gen")
	binDir := filepath.Join(dir, "target", "bin")

	if err := writeGeneratedGo(genDir, sources, dir); err != nil {
		return err
	}

	binPath := filepath.Join(binDir, "app")
	if err := goBuild(genDir, binPath); err != nil {
		return err
	}

	fmt.Printf("yzc: built %s\n", binPath)
	return nil
}

// cmdRun compiles and immediately runs the Yz project in dir.
func cmdRun(dir string) error {
	if err := cmdBuild(dir); err != nil {
		return err
	}
	binPath := filepath.Join(dir, "target", "bin", "app")
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
	relDir  string // slash-separated path relative to source root, "" for root
	name    string // file name without .yz extension
}

// walkYzFiles recursively finds all .yz files under srcRoot, skipping
// target/ and hidden directories.
func walkYzFiles(srcRoot string) ([]fileEntry, error) {
	var entries []fileEntry
	err := filepath.WalkDir(srcRoot, func(path string, d fs.DirEntry, err error) error {
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
		rel, _ := filepath.Rel(srcRoot, path)
		relDir := filepath.ToSlash(filepath.Dir(rel))
		if relDir == "." {
			relDir = ""
		}
		name := strings.TrimSuffix(d.Name(), ".yz")
		entries = append(entries, fileEntry{absPath: path, relDir: relDir, name: name})
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

// compileProject walks all .yz files, compiles each directory as a separate
// Go package, and returns a map of relative output path → Go source.
// Sub-packages are compiled before the root; their exports are registered in
// the root's analyzer so FQN references (house.front.Host()) resolve correctly.
func compileProject(srcRoot string) (map[string]string, error) {
	entries, err := walkYzFiles(srcRoot)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("no .yz files found in %s", srcRoot)
	}

	// Group by relative directory.
	byDir := map[string][]fileEntry{}
	for _, e := range entries {
		byDir[e.relDir] = append(byDir[e.relDir], e)
	}

	// Sort dirs: deepest first so sub-packages are ready before root.
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
		goSrc, exp, err := compilePackageDir(byDir[dir], dir, nil)
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
		goSrc, _, err := compilePackageDir(rootFiles, "", rootAnalyzer)
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
func compilePackageDir(files []fileEntry, relDir string, a *sema.Analyzer) (string, *pkgExport, error) {
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
	}
	var pfiles []parsedFile
	for _, fe := range files {
		src, err := os.ReadFile(fe.absPath)
		if err != nil {
			return "", nil, fmt.Errorf("reading %s: %w", fe.absPath, err)
		}
		p := parser.New(src)
		sf, err := p.ParseFile()
		if err != nil {
			return "", nil, fmt.Errorf("parse %s: %w", fe.absPath, err)
		}
		pfiles = append(pfiles, parsedFile{sf: sf, path: fe.absPath})
	}

	if a == nil {
		a = sema.NewAnalyzer()
	}
	for _, pf := range pfiles {
		if err := a.AnalyzeFile(pf.sf); err != nil {
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

	// Find the yz compiler module root so generated code can reference yz/runtime/yzrt.
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
