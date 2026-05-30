// Command yzc is the Yz compiler. It compiles .yz source files to Go source
// code and then invokes `go build` to produce a binary.
//
// Usage:
//
//	yzc [dir]           Build and run the project in dir (default ".")
//	yzc build [dir]     Compile and build the project in dir (default ".")
//	yzc run   [dir]     Build and immediately run the project
//	yzc new   <name>    Create a new Yz project skeleton
package main

import (
	"fmt"
	"os"
)

var version = "dev"

func main() {
	// No arguments: run the current directory.
	if len(os.Args) < 2 {
		if err := cmdRun("."); err != nil {
			fmt.Fprintf(os.Stderr, "yzc: %v\n", err)
			os.Exit(1)
		}
		return
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
		// Treat an unrecognised first argument as a directory to run.
		if err := cmdRun(os.Args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "yzc: %v\n", err)
			os.Exit(1)
		}
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, `Usage: yzc [dir]
       yzc <command> [arguments]

Commands:
  [dir]          Build and run the project (default dir: ".")
  build [dir]    Compile and build the project (default dir: ".")
  run   [dir]    Build and run the project
  new   <name>   Create a new Yz project
  version        Print the compiler version`)
}

