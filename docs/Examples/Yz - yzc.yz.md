#example

## yzc — compiler entry point (intent version)

> This is a speculative Yz transcription of the compiler's own `cmd/yzc/main.go`.
> It uses `os` and `cli` from the [stdlib proposal](../Questions/stdlib.md),
> which does not exist yet. It shows what the entry point _would_ look like once
> the standard library lands.

The Go original dispatches on `os.Args` with a `switch` statement. In Yz,
`os.args()` returns only the user arguments (no program name at index 0),
and `match` handles the command dispatch with condition-form branches.

```js
// yzc — The Yz compiler.
// Compiles .yz source files to Go, then invokes go build to produce a binary.
//
// Usage:
//   yzc [dir]          Build and run the project in dir (default ".")
//   yzc build [dir]    Compile and build the project
//   yzc run   [dir]    Build and immediately run the project
//   yzc new   <name>   Create a new Yz project skeleton

version: "dev"

// cmd_build, cmd_run, and cmd_new are defined in build.yz / run.yz / new.yz
// in the same directory. In a multi-file Yz project every file in a directory
// shares one namespace, so no import or forward declaration is needed here.
// These signature-only lines document the expected shape.
cmd_build #(dir String, Result(Unit, String))
cmd_run   #(dir String, Result(Unit, String))
cmd_new   #(name String, Result(Unit, String))

print_usage #() {
    cli.print('
Usage: yzc [dir]
       yzc <command> [arguments]

Commands:
  [dir]          Build and run the project (default dir: ".")
  build [dir]    Compile and build the project (default dir: ".")
  run   [dir]    Build and run the project
  new   <name>   Create a new Yz project
  version        Print the compiler version
')
}

main: {
  args: os.args()

  // No arguments: run the current directory.
  args.len() == 0 ? {
    cmd_run(".").?(err_on(""))
   }, {
    cmd: args[0]
    match
    { cmd == "build" =>
        dir: args.len() > 1 ? { args[1] }, { "." }
        cmd_build(dir).?(err_on("build"))
    },
    { cmd == "run" =>
        dir: args.len() > 1 ? { args[1] }, { "." }
        cmd_run(dir).?(err_on("run"))
    },
    { cmd == "new" =>
        args.len() < 2 ? {
	      cli.print("yzc new: missing project name")
	      os.exit(1)
        }
        cmd_new(args[1]).?(err_on("new"))
    },
    { cmd == "version" =>
        cli.print("yzc ${version}")
    },
    {
        // Unrecognised first argument — treat as a directory to run.
        cmd_run(cmd).?(err_on("new"))
    }
  }

  err_on #(message String) {
	  { 
		e String
	    cli.print("yzc ${message} ${e}")
	    os.exit(1)
	  }
  }
}
```

### What this shows

**Condition-form `match` as a dispatch table.**
Each branch is a boc `{ condition => body }`. The final branch `{ body }` with no
condition is the default — equivalent to Go's `default` case.

**`.?()` for inline error handling.**
`cmd_run(".").?({ e String; cli.print(...); os.exit(1) })` binds the error value to
`e` when the result is `Err`, runs the boc, then continues. Because `os.exit` never
returns, there is no "continue" here — but the pattern composes cleanly when the
handler just logs and propagates.

**Signature-only boc declarations.**
`cmd_build #(dir String, Result(Unit, String))` without a body declares the expected
shape for structural typing. In practice the body lives in another file; this line
documents the contract without creating a duplicate definition.

**`version` as a plain declaration.**
`version: "dev"` is a top-level singleton field. It can be overridden at link time by
the build system (same role as `var version = "dev"` in Go with `-ldflags`).

