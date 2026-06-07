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
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "build":
		dirs := os.Args[2:]
		if len(dirs) == 0 {
			fmt.Fprintln(os.Stderr, "yzc build: no source roots specified")
			fmt.Fprintln(os.Stderr, "usage: yzc build <project-dir> [extra-roots...]")
			os.Exit(1)
		}
		if err := cmdBuild(dirs[0], dirs[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "yzc build: %v\n", err)
			os.Exit(1)
		}

	case "run":
		dirs := os.Args[2:]
		if len(dirs) == 0 {
			fmt.Fprintln(os.Stderr, "yzc run: no source roots specified")
			fmt.Fprintln(os.Stderr, "usage: yzc run <project-dir> [extra-roots...]")
			os.Exit(1)
		}
		if err := cmdRun(dirs[0], dirs[1:]); err != nil {
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
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, `Usage: yzc <command> [arguments]

Commands:
  build <project-dir> [extra-roots...]   Compile and build the project
  run   <project-dir> [extra-roots...]   Build and run the project
  new   <name>                           Create a new Yz project
  version                                Print the compiler version`)
}

