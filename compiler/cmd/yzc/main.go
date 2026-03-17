// Command yzc is the Yz compiler. It compiles .yz source files to Go source
// code and then invokes `go build` to produce a binary.
//
// Usage:
//
//	yzc build [dir]     Compile and build the project in dir (default ".")
//	yzc run   [dir]     Build and immediately run the project
//	yzc new   <name>    Create a new Yz project skeleton
package main

import (
	"fmt"
	"os"
)

const version = "0.1.0-dev"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "build":
		dir := "."
		if len(os.Args) > 2 {
			dir = os.Args[2]
		}
		if err := cmdBuild(dir); err != nil {
			fmt.Fprintf(os.Stderr, "yzc build: %v\n", err)
			os.Exit(1)
		}

	case "run":
		dir := "."
		if len(os.Args) > 2 {
			dir = os.Args[2]
		}
		if err := cmdRun(dir); err != nil {
			fmt.Fprintf(os.Stderr, "yzc run: %v\n", err)
			os.Exit(1)
		}

	case "new":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "yzc new: missing project name")
			os.Exit(1)
		}
		if err := cmdNew(os.Args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "yzc new: %v\n", err)
			os.Exit(1)
		}

	case "version":
		fmt.Printf("yzc %s\n", version)

	default:
		fmt.Fprintf(os.Stderr, "yzc: unknown command %q\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, `Usage: yzc <command> [arguments]

Commands:
  build [dir]    Compile and build the project (default dir: ".")
  run   [dir]    Build and run the project
  new   <name>   Create a new Yz project
  version        Print the compiler version`)
}

// cmdBuild compiles the Yz project in dir to a binary.
func cmdBuild(dir string) error {
	// TODO: implement lexer → parser → sema → ir → codegen → go build pipeline
	fmt.Printf("yzc: building %s (not yet implemented)\n", dir)
	return nil
}

// cmdRun compiles and runs the Yz project in dir.
func cmdRun(dir string) error {
	if err := cmdBuild(dir); err != nil {
		return err
	}
	// TODO: execute the built binary
	fmt.Println("yzc: run not yet implemented")
	return nil
}

// cmdNew scaffolds a new Yz project directory.
func cmdNew(name string) error {
	// TODO: create project.yz, main.yz, etc.
	fmt.Printf("yzc: creating project %q (not yet implemented)\n", name)
	return nil
}
