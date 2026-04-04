// Package diagnostic formats Rust-style error messages with source context.
package diagnostic

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

// ANSI color codes — only emitted when stderr is a color terminal.
const (
	red    = "\x1b[1;31m"
	cyan   = "\x1b[36m"
	bold   = "\x1b[1m"
	reset  = "\x1b[0m"
)

// colorEnabled reports whether color output should be used.
// Respects the NO_COLOR env var (https://no-color.org/) and checks if
// stderr is a real TTY.
func colorEnabled() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	fi, err := os.Stderr.Stat()
	if err != nil {
		return false
	}
	// ModeCharDevice is set for a real terminal.
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func col(color, s string) string {
	if !colorEnabled() {
		return s
	}
	return color + s + reset
}

// Format returns a Rust-style diagnostic string.
//
//	error: undefined variable: foo
//	 --> main.yz:3:5
//	  |
//	3 |     foo = 42
//	  |     ^^^
//
// src is the full source file content. line and col are 1-based.
// length is the number of source characters to underline (0 → 1 caret).
func Format(src []byte, filename string, line, col_, length int, msg string) string {
	var b strings.Builder

	// "error: msg"
	b.WriteString(col(red, "error"))
	b.WriteString(col(bold, ": "+msg))
	b.WriteString("\n")

	// " --> file:line:col"
	b.WriteString(col(cyan, fmt.Sprintf(" --> %s:%d:%d", filename, line, col_)))
	b.WriteString("\n")

	// Extract the source line.
	srcLine := extractLine(src, line)
	if srcLine == "" {
		return b.String()
	}

	lineStr := fmt.Sprintf("%d", line)
	pad := strings.Repeat(" ", len(lineStr))

	// "  |"
	b.WriteString(col(cyan, pad+" |"))
	b.WriteString("\n")

	// "3 | source line"
	b.WriteString(col(cyan, lineStr+" | "))
	b.WriteString(srcLine)
	b.WriteString("\n")

	// "  |     ^^^"
	if length <= 0 {
		length = 1
	}
	caretCol := col_ - 1 // zero-based
	if caretCol < 0 {
		caretCol = 0
	}
	carets := strings.Repeat("^", length)
	spaces := strings.Repeat(" ", caretCol)
	b.WriteString(col(cyan, pad+" | "))
	b.WriteString(col(red, spaces+carets))
	b.WriteString("\n")

	return b.String()
}

// extractLine returns the src line at lineNum (1-based), without trailing newline.
func extractLine(src []byte, lineNum int) string {
	cur := 1
	start := 0
	for i, c := range src {
		if cur == lineNum {
			// Find end of line.
			end := bytes.IndexByte(src[i:], '\n')
			if end < 0 {
				return string(src[i:])
			}
			return string(src[i : i+end])
		}
		if c == '\n' {
			cur++
			start = i + 1
		}
		_ = start
	}
	return ""
}
