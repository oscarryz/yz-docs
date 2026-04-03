// Package conformance_test validates the Yz compiler end-to-end.
//
// Two test suites:
//
//   - TestGolden: .yz source → generated .go output compared to a golden file.
//   - TestErrors: .yz source → compile must fail; error message must contain the
//     substring stored in the matching .error file.
//
// To regenerate golden files after an intentional output change:
//
//	UPDATE_GOLDEN=1 go test ./test/conformance/...
package conformance_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"yz/internal/codegen"
	"yz/internal/ir"
	"yz/internal/parser"
	"yz/internal/sema"
)

// compile runs the full Yz pipeline and returns the generated Go source.
func compile(src []byte) (string, error) {
	p := parser.New(src)
	sf, err := p.ParseFile()
	if err != nil {
		return "", err
	}
	a := sema.NewAnalyzer()
	if err := a.AnalyzeFile(sf); err != nil {
		return "", err
	}
	f := ir.Lower(sf, a, "main")
	return codegen.Generate(f), nil
}

// TestGolden walks testdata/golden/, compiles each .yz file, and compares the
// output to the corresponding .go golden file.
//
// Set UPDATE_GOLDEN=1 to write/overwrite golden files instead of comparing.
func TestGolden(t *testing.T) {
	const dir = "testdata/golden"
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("reading %s: %v", dir, err)
	}
	update := os.Getenv("UPDATE_GOLDEN") != ""

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yz") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".yz")
		t.Run(name, func(t *testing.T) {
			src, err := os.ReadFile(filepath.Join(dir, e.Name()))
			if err != nil {
				t.Fatalf("reading source: %v", err)
			}

			got, err := compile(src)
			if err != nil {
				t.Fatalf("compile: %v", err)
			}

			goldenPath := filepath.Join(dir, name+".go")
			if update {
				if err := os.WriteFile(goldenPath, []byte(got), 0o644); err != nil {
					t.Fatalf("writing golden: %v", err)
				}
				t.Logf("updated %s", goldenPath)
				return
			}

			want, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("missing golden file %s\n\nRun with UPDATE_GOLDEN=1 to create it", goldenPath)
			}
			if got != string(want) {
				t.Errorf("output mismatch for %s\n\n--- want ---\n%s\n--- got ---\n%s", name, want, got)
			}
		})
	}
}

// TestErrors walks testdata/errors/, compiles each .yz file, and asserts that:
//  1. Compilation fails.
//  2. The error message contains the substring from the matching .error file.
func TestErrors(t *testing.T) {
	const dir = "testdata/errors"
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("reading %s: %v", dir, err)
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yz") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".yz")
		t.Run(name, func(t *testing.T) {
			src, err := os.ReadFile(filepath.Join(dir, e.Name()))
			if err != nil {
				t.Fatalf("reading source: %v", err)
			}

			_, cerr := compile(src)
			if cerr == nil {
				t.Fatal("expected a compile error, got none")
			}

			errPath := filepath.Join(dir, name+".error")
			want, err := os.ReadFile(errPath)
			if err != nil {
				t.Fatalf("missing .error file %s", errPath)
			}
			wantSubstr := strings.TrimSpace(string(want))
			if !strings.Contains(cerr.Error(), wantSubstr) {
				t.Errorf("error message does not contain expected substring\n\nwant substr: %q\ngot error:   %q", wantSubstr, cerr.Error())
			}
		})
	}
}
