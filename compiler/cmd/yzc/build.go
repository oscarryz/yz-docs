package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"yz/internal/ir"
	"yz/internal/parser"
	"yz/internal/sema"
	"yz/internal/codegen"
)

// cmdBuild compiles the Yz project in dir to a binary at target/bin/app.
func cmdBuild(dir string) error {
	src, err := readProjectDir(dir)
	if err != nil {
		return err
	}

	goSrc, err := compileSource(src)
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

// readProjectDir finds main.yz (or any single .yz file) in dir and returns its source.
func readProjectDir(dir string) ([]byte, error) {
	mainPath := filepath.Join(dir, "main.yz")
	if src, err := os.ReadFile(mainPath); err == nil {
		return src, nil
	}

	// Fall back: any single .yz file in dir.
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading directory %s: %w", dir, err)
	}
	var yzFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".yz") {
			yzFiles = append(yzFiles, filepath.Join(dir, e.Name()))
		}
	}
	switch len(yzFiles) {
	case 0:
		return nil, fmt.Errorf("no .yz files found in %s", dir)
	case 1:
		return os.ReadFile(yzFiles[0])
	default:
		return nil, fmt.Errorf("multiple .yz files in %s: use main.yz as entry point", dir)
	}
}

// compileSource runs the full Yz pipeline and returns Go source bytes.
func compileSource(src []byte) (string, error) {
	p := parser.New(src)
	sf, err := p.ParseFile()
	if err != nil {
		return "", fmt.Errorf("parse: %w", err)
	}

	a := sema.NewAnalyzer()
	if err := a.AnalyzeFile(sf); err != nil {
		return "", fmt.Errorf("sema: %w", err)
	}

	f := ir.Lower(sf, a, "main")
	return codegen.Generate(f), nil
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
