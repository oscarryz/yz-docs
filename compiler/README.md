#readme 
# yzc — The Yz Compiler

Compiles Yz source files (`.yz`) to Go source code, then invokes `go build` to produce a native binary.

## Quick Start

```bash
# Build the compiler
make build

# Create a new project
bin/yzc new hello

# Build a project
bin/yzc build .

# Build and run
bin/yzc run .
```

## Architecture

```
.yz → Lexer → Parser → AST → Sema → IR → Codegen → .go → go build → Binary
```

## Development

```bash
# Run all tests
make test

# Clean build artifacts
make clean
```
