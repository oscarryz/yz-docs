package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// cmdNew creates a new Yz project skeleton in a directory named <name>.
func cmdNew(name string) error {
	if err := os.MkdirAll(name, 0o755); err != nil {
		return fmt.Errorf("creating directory %s: %w", name, err)
	}

	mainYz := `main: {
    print("Hello, World!")
}
`
	mainPath := filepath.Join(name, "main.yz")
	if err := os.WriteFile(mainPath, []byte(mainYz), 0o644); err != nil {
		return fmt.Errorf("writing main.yz: %w", err)
	}

	gitignore := "target/\n"
	giPath := filepath.Join(name, ".gitignore")
	if err := os.WriteFile(giPath, []byte(gitignore), 0o644); err != nil {
		return fmt.Errorf("writing .gitignore: %w", err)
	}

	fmt.Printf("yzc: created project %q\n", name)
	fmt.Printf("  %s/main.yz\n", name)
	fmt.Printf("  %s/.gitignore\n", name)
	fmt.Printf("\nRun: yzc run %s\n", name)
	return nil
}
