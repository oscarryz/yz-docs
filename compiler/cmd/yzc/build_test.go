package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCmdRunConnectsStdin guards against a regression where cmdRun's
// exec.Command never set cmd.Stdin: the compiled binary's stdin then fell
// back to Go's default (the null device), so any `read(...)` call hit
// instant EOF and returned "" instead of waiting for input. See YZC-0120.
func TestCmdRunConnectsStdin(t *testing.T) {
	if testing.Short() {
		t.Skip("compiles and runs a real binary; skip in short mode")
	}

	dir := t.TempDir()
	src := "main: {\n    name: read(\"\")\n    print(\"got:${name}\")\n}\nmain()\n"
	if err := os.WriteFile(filepath.Join(dir, "main.yz"), []byte(src), 0o644); err != nil {
		t.Fatalf("writing main.yz: %v", err)
	}

	stdinR, stdinW, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe (stdin): %v", err)
	}
	stdoutR, stdoutW, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe (stdout): %v", err)
	}

	oldStdin, oldStdout := os.Stdin, os.Stdout
	os.Stdin, os.Stdout = stdinR, stdoutW
	defer func() { os.Stdin, os.Stdout = oldStdin, oldStdout }()

	if _, err := stdinW.WriteString("Alice\n"); err != nil {
		t.Fatalf("writing to stdin pipe: %v", err)
	}
	stdinW.Close()

	outCh := make(chan string, 1)
	go func() {
		out, _ := io.ReadAll(stdoutR)
		outCh <- string(out)
	}()

	runErr := cmdRun(dir, nil)
	stdoutW.Close()
	out := <-outCh

	if runErr != nil {
		t.Fatalf("cmdRun: %v\noutput: %s", runErr, out)
	}
	if !strings.Contains(out, "got:Alice") {
		t.Errorf("expected the child process's stdin to be connected to ours (via os.Stdin); got output %q", out)
	}
}
