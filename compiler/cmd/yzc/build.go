package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"yz/internal/ast"
	"yz/internal/codegen"
	"yz/internal/ir"
	"yz/internal/parser"
	"yz/internal/sema"
)

// cmdBuild compiles the Yz project in dir to a binary at target/bin/app.
func cmdBuild(dir string) error {
	paths, err := collectYzFiles(dir)
	if err != nil {
		return err
	}

	goSrc, err := compileFiles(paths)
	if err != nil {
		return err
	}

	genDir := filepath.Join(dir, "target", "gen")
	binDir := filepath.Join(dir, "target", "bin")

	if err := writeGeneratedGo(genDir, goSrc, dir); err != nil {
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

// collectYzFiles returns all .yz file paths in dir (non-recursive, flat only).
// main.yz is always last so it can reference types from other files.
func collectYzFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading directory %s: %w", dir, err)
	}
	var main, others []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yz") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		if e.Name() == "main.yz" {
			main = append(main, path)
		} else {
			others = append(others, path)
		}
	}
	if len(others)+len(main) == 0 {
		return nil, fmt.Errorf("no .yz files found in %s", dir)
	}
	// Non-main files first (alphabetical via ReadDir), main.yz last.
	return append(others, main...), nil
}

// compileFiles runs the full Yz pipeline over multiple source files and
// returns combined Go source. Files are analyzed in order (main.yz last)
// using a single shared sema scope so cross-file references resolve.
func compileFiles(paths []string) (string, error) {
	type parsedFile struct {
		sf   *ast.SourceFile
		path string
	}

	// Parse all files.
	var pfiles []parsedFile
	for _, path := range paths {
		src, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("reading %s: %w", path, err)
		}
		p := parser.New(src)
		sf, err := p.ParseFile()
		if err != nil {
			return "", fmt.Errorf("parse %s: %w", path, err)
		}
		pfiles = append(pfiles, parsedFile{sf: sf, path: path})
	}

	// Analyze all files with a single shared scope.
	a := sema.NewAnalyzer()
	for _, pf := range pfiles {
		if err := a.AnalyzeFile(pf.sf); err != nil {
			return "", fmt.Errorf("sema %s: %w", pf.path, err)
		}
	}

	// Lower all files and merge decls into one ir.File.
	combined := &ir.File{PkgName: "main"}
	for _, pf := range pfiles {
		f := ir.Lower(pf.sf, a, "main")
		combined.Decls = append(combined.Decls, f.Decls...)
		for _, imp := range f.Imports {
			if !containsStr(combined.Imports, imp) {
				combined.Imports = append(combined.Imports, imp)
			}
		}
	}

	return codegen.Generate(combined), nil
}

func containsStr(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

// writeGeneratedGo writes main.go and go.mod into genDir.
func writeGeneratedGo(genDir, goSrc, projectDir string) error {
	if err := os.MkdirAll(genDir, 0o755); err != nil {
		return fmt.Errorf("creating gen dir: %w", err)
	}

	// Write main.go.
	mainPath := filepath.Join(genDir, "main.go")
	if err := os.WriteFile(mainPath, []byte(goSrc), 0o644); err != nil {
		return fmt.Errorf("writing main.go: %w", err)
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
