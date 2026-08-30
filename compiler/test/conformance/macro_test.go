package conformance_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestMacros exercises the macro system (YZC-0028) through the real build
// driver. Each directory under testdata/macros/ is a full project:
//
//   - with an expected.error file: yzc build must fail and its combined
//     output must contain the file's substring.
//   - without one: the build must succeed; if main.output exists the binary
//     is run and compared; a second build must reuse the cached macro binary.
//
// Skipped under -short (each case builds macro executables with go build).
func TestMacros(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping macro driver tests in short mode")
	}

	yzc := buildYzc(t)

	const dir = "testdata/macros"
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("reading %s: %v", dir, err)
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		projectDir := filepath.Join(dir, name)

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Clean target/ so cached state from previous runs cannot mask
			// regressions.
			os.RemoveAll(filepath.Join(projectDir, "target"))

			wantErr, expectFailure := readFileIfExists(t, filepath.Join(projectDir, "expected.error"))

			out, err := exec.Command(yzc, "build", projectDir).CombinedOutput()
			if expectFailure {
				if err == nil {
					t.Fatalf("expected build failure, got success:\n%s", out)
				}
				if !strings.Contains(string(out), strings.TrimSpace(wantErr)) {
					t.Fatalf("build output missing %q:\n%s", strings.TrimSpace(wantErr), out)
				}
				return
			}
			if err != nil {
				t.Fatalf("yzc build failed:\n%s", out)
			}

			if want, ok := readFileIfExists(t, filepath.Join(projectDir, "main.output")); ok {
				appBin := filepath.Join(projectDir, "target", "bin", "app")
				got, err := exec.Command(appBin).Output()
				if err != nil {
					t.Fatalf("running app: %v", err)
				}
				if strings.TrimRight(string(got), "\n") != strings.TrimRight(want, "\n") {
					t.Errorf("output mismatch\nwant:\n%s\ngot:\n%s", want, got)
				}
			}

			// Rebuild: the macro executable must come from the cache.
			macroBins, _ := filepath.Glob(filepath.Join(projectDir, "target", "macros", "*", "bin", "macros"))
			if len(macroBins) == 0 {
				t.Fatal("no macro binary produced")
			}
			before, err := os.Stat(macroBins[0])
			if err != nil {
				t.Fatal(err)
			}
			if out, err := exec.Command(yzc, "build", projectDir).CombinedOutput(); err != nil {
				t.Fatalf("rebuild failed:\n%s", out)
			}
			after, err := os.Stat(macroBins[0])
			if err != nil {
				t.Fatal(err)
			}
			if !before.ModTime().Equal(after.ModTime()) {
				t.Error("macro binary rebuilt despite unchanged source")
			}
		})
	}
}

func readFileIfExists(t *testing.T, path string) (string, bool) {
	t.Helper()
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return "", false
	}
	if err != nil {
		t.Fatal(err)
	}
	return string(data), true
}
