#feature
# Go Extensions — Go-backed type implementations

## Overview

`go_source:` is a passive annotation key that links a Yz type declaration to a Go source file. Methods declared without a body in the Yz type are implemented in the linked Go file. The compiler includes the Go file directly in the build output.

This is the mechanism for:
- Stdlib scalar types (`Int`, `String`, `Bool`, `Decimal`) — backed by Go primitives
- Standard library extensions — Go-backed methods alongside pure Yz methods
- Third-party library wrappers — adapting Go libraries to Yz's type and concurrency model

`go_source:` is **not a macro** — it does not dispatch through the macro system and does not transform the AST. It is processed in the compiler's first pass, before type resolution.

---

## Annotation syntax

### Type-level

All body-less methods in the type delegate to the named Go file:

```yz
`go_source: "stdlib/int.go"`
Int: {
    + #(other Int, Int)           // body-less — implemented in int.go
    parse #(s String, Int)        // body-less — implemented in int.go
    times #(n Int, Range) {       // has a body — pure Yz, not delegated
        Range(0, n)
    }
}
```

### Method-level

Only this method delegates; the file may differ from a type-level annotation:

```yz
Int: {
    `go_source: "stdlib/int.go"`
    parse #(s String, Int)
}
```

Method-level is useful when methods come from different Go files, or when only a subset of methods are Go-backed.

### Third-party wrapper

A type with only a signature (no body) can be fully backed by Go:

```yz
`go_source: "clients/http_client.go"`
HttpClient #(
    get  #(url String, String)
    post #(url String, body String, String)
)
```

---

## Go binding comment

Each Go function that implements a Yz method must have a `//yz:bind` comment on the line **immediately above** the `func` declaration. No blank lines between the comment and the `func`.

```go
// stdlib/int.go
package yzstd

import "github.com/oscarryz/yz/runtime/std"

//yz:bind Int parse #(s String, Int)
func IntParse(s std.String) std.Int { ... }

//yz:bind Int + #(other Int, Int)
func IntPlus(a, b std.Int) std.Int { ... }

//yz:bind Int == #(other Int, Bool)
func IntEqeq(a, b std.Int) std.Bool { ... }
```

### Format

```
//yz:bind TypeName methodSignature
```

- `TypeName` — the Yz type that owns the method
- `methodSignature` — a standard Yz method declaration; parsed by the existing Yz parser
- Non-word method names (`+`, `==`, `?`) work naturally — they are valid in Yz signatures

---

## Compiler first-pass

`go_source:` is processed before type resolution in a dedicated first pass:

1. Scan all annotations across source roots for `go_source:` keys
2. Collect the listed `.go` files
3. Line-scan each `.go` file for `//yz:bind` lines; the binding comment must be
   immediately above a `func` declaration
4. Parse the method signature using the existing Yz parser
5. Build a map: `TypeName.methodName → GoFuncName`
6. Validate: every body-less method on a `go_source:`-annotated type must have
   a `//yz:bind` entry — missing binding is a compile error
7. Include the Go files in the build output alongside the generated Go code

---

## Go API contract

Go files linked via `go_source:` must follow these rules:

**Types** — all parameters and return values use `std.*` types:

| Yz type   | Go type      |
|-----------|--------------|
| `Int`     | `std.Int`    |
| `String`  | `std.String` |
| `Bool`    | `std.Bool`   |
| `Decimal` | `std.Decimal`|
| `Unit`    | `std.Unit`   |

**Errors** — returned as `std.Result[T]`, never as Go `error`:
```go
//yz:bind Int parse #(s String, Int)
func IntParse(s std.String) std.Result[std.Int] { ... }
```

**Concurrency** — the Go function is synchronous from Yz's perspective. Do not spawn goroutines directly. Concurrency is managed by the Yz runtime (cowns, BocGroup). The Go function receives values, computes, returns — nothing else.

---

## File location

Go source files live alongside their `.yz` counterparts in the source tree:

```
stdlib/
  int.yz      ← Yz declaration with go_source: "stdlib/int.go"
  int.go      ← Go implementation
  string.yz
  string.go
```

The path in `go_source:` is relative to the project source root.

---

## What this does not cover

- `self` inside Go-backed methods — see YZC-0060
- Generic Go-backed methods — deferred
- Goroutine-spawning Go code — not allowed in `go_source:` files; use the Yz runtime API instead

---

## See also

- [Annotations](Annotations.md) — annotation syntax
- [Dependencies](Dependencies.md) — declaring external Yz packages
- [Macros](Macros.md) — compile-time code generation
