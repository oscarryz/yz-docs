package conformance_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestExamples builds every directory under examples/ (skipping those starting
// with '_') and fails if yzc build returns an error.  If a main.output sidecar
// exists the compiled binary is run and its stdout is compared against it.
//
// These tests are skipped under -short because building and running examples
// can take several seconds.
func TestExamples(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping example builds in short mode")
	}

	yzc := buildYzc(t)

	// Resolve the examples/ directory relative to this test file's module root.
	moduleRoot, err := findYzRoot()
	if err != nil {
		t.Fatalf("finding module root: %v", err)
	}
	examplesDir := filepath.Join(moduleRoot, "examples")

	entries, err := os.ReadDir(examplesDir)
	if err != nil {
		t.Fatalf("reading examples dir: %v", err)
	}

	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), "_") {
			continue
		}
		name := e.Name()
		exampleDir := filepath.Join(examplesDir, name)

		t.Run(name, func(t *testing.T) {
			// Build.
			cmd := exec.Command(yzc, "build", exampleDir)
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("yzc build failed:\n%s", out)
			}

			// Run + compare output if a sidecar exists.
			outputFile := filepath.Join(exampleDir, "main.output")
			wantBytes, err := os.ReadFile(outputFile)
			if os.IsNotExist(err) {
				return // compile-only check
			}
			if err != nil {
				t.Fatalf("reading main.output: %v", err)
			}
			want := strings.TrimRight(string(wantBytes), "\n")

			appBin := filepath.Join(exampleDir, "target", "bin", "app")
			out, err := exec.Command(appBin).Output()
			if err != nil {
				t.Fatalf("running app: %v", err)
			}
			got := strings.TrimRight(string(out), "\n")
			if got != want {
				t.Errorf("output mismatch\nwant:\n%s\ngot:\n%s", want, got)
			}
		})
	}
}

// buildYzc compiles the yzc binary into a temp directory and returns its path.
func buildYzc(t *testing.T) string {
	t.Helper()
	out := filepath.Join(t.TempDir(), "yzc")
	cmd := exec.Command("go", "build", "-o", out, "yz/cmd/yzc")
	if combined, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("building yzc: %v\n%s", err, combined)
	}
	return out
}

