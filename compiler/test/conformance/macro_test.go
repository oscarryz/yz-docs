package conformance_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestMacros runs driver-level macro compilation cases under testdata/macros/.
// Each sub-directory is one case:
//   - expected.error: yzc build must fail; error output must contain the substring.
//   - expected.output: yzc build must succeed; the compiled binary must print this.
//   - If neither file is present, only a successful build is required.
func TestMacros(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping macro build tests in short mode")
	}

	yzc := buildYzc(t)

	moduleRoot, err := findYzRoot()
	if err != nil {
		t.Fatalf("finding module root: %v", err)
	}
	macrosDir := filepath.Join(moduleRoot, "test", "conformance", "testdata", "macros")

	entries, err := os.ReadDir(macrosDir)
	if err != nil {
		t.Fatalf("reading %s: %v", macrosDir, err)
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		caseDir := filepath.Join(macrosDir, name)

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cmd := exec.Command(yzc, "build", caseDir)
			out, buildErr := cmd.CombinedOutput()

			errFile := filepath.Join(caseDir, "expected.error")
			wantErrBytes, readErr := os.ReadFile(errFile)
			if readErr == nil {
				// Error case: build must fail and output must contain the substring.
				wantErr := strings.TrimSpace(string(wantErrBytes))
				if buildErr == nil {
					t.Fatalf("expected build error containing %q, got success", wantErr)
				}
				if !strings.Contains(string(out), wantErr) {
					t.Errorf("error mismatch\nwant substring: %q\ngot output:     %q", wantErr, string(out))
				}
				return
			}

			// Success case: build must pass.
			if buildErr != nil {
				t.Fatalf("yzc build failed:\n%s", out)
			}

			outputFile := filepath.Join(caseDir, "expected.output")
			wantBytes, err := os.ReadFile(outputFile)
			if os.IsNotExist(err) {
				return // compile-only check
			}
			if err != nil {
				t.Fatalf("reading expected.output: %v", err)
			}
			want := strings.TrimRight(string(wantBytes), "\n")

			appBin := filepath.Join(caseDir, "target", "bin", "app")
			appOut, err := exec.Command(appBin).Output()
			if err != nil {
				t.Fatalf("running app: %v", err)
			}
			got := strings.TrimRight(string(appOut), "\n")
			if got != want {
				t.Errorf("output mismatch\nwant: %q\ngot:  %q", want, got)
			}
		})
	}
}
