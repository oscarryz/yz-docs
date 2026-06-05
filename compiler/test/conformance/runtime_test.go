package conformance_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestRuntime compiles every golden .yz that has a matching .output sidecar,
// runs the generated binary in a temp directory, and asserts that stdout
// matches the .output file. This catches runtime bugs (deadlocks, wrong
// ordering, incorrect values) that the source-diff TestGolden cannot see.
//
// Skipped under `go test -short` — use that flag during development for a
// faster feedback loop, and run without -short at the end of a ticket or in CI.
func TestRuntime(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping runtime tests in short mode; run without -short for full validation")
	}
	const dir = "testdata/golden"
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("reading %s: %v", dir, err)
	}

	yzRoot, err := findYzRoot()
	if err != nil {
		t.Fatalf("locating yz module: %v", err)
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".output") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".output")
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			wantBytes, err := os.ReadFile(filepath.Join(dir, e.Name()))
			if err != nil {
				t.Fatalf("reading .output: %v", err)
			}
			want := strings.TrimRight(string(wantBytes), "\n")

			src, err := os.ReadFile(filepath.Join(dir, name+".yz"))
			if err != nil {
				t.Fatalf("reading source: %v", err)
			}
			goSrc, err := compile(src)
			if err != nil {
				t.Fatalf("compile: %v", err)
			}

			tmp := t.TempDir()
			if err := os.WriteFile(filepath.Join(tmp, "main.go"), []byte(goSrc), 0o644); err != nil {
				t.Fatalf("writing main.go: %v", err)
			}
			goMod := fmt.Sprintf("module yzapp\n\ngo 1.23\n\nrequire yz v0.0.0\n\nreplace yz => %s\n", yzRoot)
			if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte(goMod), 0o644); err != nil {
				t.Fatalf("writing go.mod: %v", err)
			}

			cmd := exec.Command("go", "run", ".")
			cmd.Dir = tmp
			out, err := cmd.Output()
			if err != nil {
				t.Fatalf("go run failed: %v\noutput: %s", err, out)
			}
			got := strings.TrimRight(string(out), "\n")
			if got != want {
				t.Errorf("output mismatch\nwant: %q\ngot:  %q", want, got)
			}
		})
	}
}

// findYzRoot walks up from the current working directory to find the directory
// containing a go.mod that declares "module yz".
func findYzRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		data, err := os.ReadFile(filepath.Join(dir, "go.mod"))
		if err == nil && strings.Contains(string(data), "module yz\n") {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod with 'module yz' not found walking up from %s", dir)
		}
		dir = parent
	}
}
