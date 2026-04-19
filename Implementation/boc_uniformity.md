# Boc Uniformity — Design Gap Analysis

## The Intended Design

In Yz, **everything is a boc**. There is one construct. The same rules apply regardless of where a boc appears — in a directory, in a file, inside another boc, inside that boc. Nesting depth is irrelevant to semantics.

### Source roots

A project has one or more **source root directories**. Everything inside a source root is part of the boc tree; the source root directory itself is not a boc — it is just the mount point. FQNs are relative to the source root, not to the project directory.

A typical project might have:

```
project/
  src/          ← source root 1
    foo/
      bar.yz
  lib/          ← source root 2  (3rd-party / stdlib)
    baz.yz
```

This produces two FQN trees:

```
src root:  foo → bar → (declarations inside bar.yz)
lib root:  baz → (declarations inside baz.yz)
```

FQNs from `src`: `foo.bar`, `foo.bar.something`
FQNs from `lib`: `baz`, `baz.something`

A single source root (the default: `.`) means the project root is the mount point, giving the same behavior as today. Multiple source roots allow third-party libraries and a future stdlib to live in separate trees without polluting the app's namespace.

> The `src/` vs `lib/` tooling convention — versions, lock files, download caching — is a future concern. What matters now is the principle: source roots are mount points; bocs are everything inside them.

### Bocs form a forest

Given a single source root:

```
project/
  xyz/
    foo.yz:
      bar: {}
      baz: {
        qux: {}
      }
```

The FQN graph is:

```
xyz → foo → bar
          → baz → qux
```

- `xyz` is a boc (the directory)
- `foo` is a boc (the file)
- `bar` is a boc (a declaration inside `foo`)
- `baz` is a boc (another declaration inside `foo`)
- `qux` is a boc (a declaration inside `baz`)

`bar` being "file-scope" and `qux` being "nested" is not a semantic distinction — it is only a position in the tree. The same rules apply to both:

- Lowercase name → singleton (shared by all callers at that scope)
- Uppercase name → type (each call creates a fresh instance)
- `#(params)` → stateless function (no persistent fields, parallel calls)
- Inner declarations → children in the FQN tree
- Calls are async, return a thunk
- Access via FQN from outside, local name from inside

## Current Implementation: Three Separate Code Paths

The compiler currently has three distinct lowering paths that diverged as implementation shortcuts:

### Path 1: File-scope boc (`counter: { ... }`)

```yz
counter: {
    count: 0
    increment: { count = count + 1 }
    value: { count }
}
```

Generates a Go struct + package-level singleton var + goroutine-wrapped methods. **Works correctly.**

```go
type _counterBoc struct { count std.Int }
func (self *_counterBoc) Increment() *std.Thunk[std.Unit] { return std.Go(...) }
var Counter = &_counterBoc{count: std.NewInt(0)}
```

### Path 2: File-scope BocWithSig (`countdown #(n Int) { ... }`)

```yz
countdown #(n Int) {
    n == 0 ? { print("done") }, { print(n); countdown(n - 1) }
}
```

Generates a Go top-level function returning `*Thunk[T]`. **Works correctly.**

```go
func countdown(n std.Int) *std.Thunk[std.Unit] { return std.Go(...) }
```

### Path 3: Local boc inside a boc body (`f: { n Int; ... }`)

```yz
main: {
    f: { n Int; n == 0 ? { print("fin") }, { print(n); f(n-1) } }
    f(3)
}
```

Currently generates a plain Go function literal stored in an `any` variable. **Broken:**

```go
func main() {
    var f any = func(n std.Int) std.Unit { ... }   // typed as `any`
    std.Go(func() std.Unit { return f(std.NewInt(3)) })  // no BocGroup, exits immediately
}
```

Inner bocs and BocWithSig methods on local bocs are silently dropped. The `baz #() { ... }` in:

```yz
main: {
    boc: {
        bar: {}
        baz #() { print("baz") }
    }
    boc.baz()
}
```

...generates `boc.Baz()` on an `any`-typed function literal — this doesn't compile.

## Gap Summary

| Capability | File-scope | Local (nested) |
|---|---|---|
| Singleton struct emitted | Yes | No — plain func literal |
| Methods (BocWithSig) on boc | Yes | Dropped |
| Inner bocs as fields | Yes | Dropped |
| Goroutine-wrapped method calls | Yes | Sometimes — wrong type |
| Structured concurrency (BocGroup) | Yes for `main:` | No |
| Recursive self-calls | Yes (BocWithSig fixed) | Fire-and-forget |
| `.field` access | Yes | No |
| FQN registration in sema | Yes | Partial (pre-reg fix) |

## Root Cause

The lowerer (`ir/lower.go`) dispatches on the **position** of a node in the source file:

- `lowerTopLevel` for file-scope declarations
- `lowerMainBoc` for `main:` body statements
- `lowerBocBody` for method bodies
- `lowerClosureBody` / `lowerBocAsStmts` for inline boc literals

Each path has different code for handling nested `BocLiteral` nodes, and none of the inner paths produce proper struct types with methods.

The sema analyzer (`sema/analyzer.go`) has a corresponding split:

- `analyzeBocDecl` for `ShortDecl` with `BocLiteral` RHS (both file-scope and local)
- `analyzeStructBoc` for Uppercase boc bodies
- `analyzeBocWithSig` for `BocWithSig` nodes

`analyzeStructBoc` (the struct-emitting path) is only called for Uppercase names — so a lowercase local boc never gets the struct treatment even though it semantically should.

## Desired State

A locally-defined boc like:

```yz
main: {
    boc: {
        bar: {}
        baz #() { print("baz") }
    }
    boc.baz()
}
```

Should compile to roughly:

```go
// Lifted to package level (Go doesn't support methods on local types)
type _main_bocBoc struct{}

func (self *_main_bocBoc) Bar() *std.Thunk[std.Unit] {
    return std.Go(func() std.Unit { return std.TheUnit })
}
func (self *_main_bocBoc) Baz() *std.Thunk[std.Unit] {
    return std.Go(func() std.Unit {
        std.Print(std.NewString("baz"))
        return std.TheUnit
    })
}

func main() {
    _boc := &_main_bocBoc{}   // scoped to this invocation
    _bg0 := &std.BocGroup{}
    _bg0.Go(func() any { return _boc.Baz().Force() })
    _bg0.Wait()
}
```

Key principles:
- The struct type is **lifted to package level** (Go constraint), but the **instance** is created inside the enclosing function body (so Uppercase outer bocs get fresh inner bocs per instance)
- The struct name is derived from the **FQN** to avoid collisions: `_main_bocBoc`, `_main_baz_quxBoc`, etc.
- Method calls on a local boc instance behave identically to method calls on a file-scope singleton
- `BocGroup` coordination applies at the call site

## Implementation Approach

This is a non-trivial change. It likely requires two or three passes:

### Pass 1 — Sema: uniform FQN + type recording (lower risk)

- Ensure `analyzeBocDecl` records a proper `StructType` (not just `BocType`) for lowercase local bocs that contain inner bocs or BocWithSig methods — the same struct-shape analysis already done for Uppercase bocs
- FQN registration for nested bocs should mirror the file-scope case
- Sema changes are safer since they don't affect generated output until the lowerer is updated

### Pass 2 — Lowerer: lift nested boc structs to package level

- Introduce a **pre-pass** (or recursive collection phase) that walks the AST and collects all boc declarations at any nesting depth
- For each one, generate a `_fqnBoc` struct + methods at package level (same as today's `lowerTopLevel`)
- In the enclosing function body, emit a local variable holding a `new(_fqnBoc)` instance instead of a function literal

This requires restructuring the lowerer's top-level collection phase to handle nested bocs before the function bodies that reference them. The current single-pass approach emits declarations in encounter order; we'd need a two-phase approach: collect-all-types first, then emit function bodies.

### Pass 3 — Unify the lowering paths

Once Passes 1 and 2 work, several of the separate code paths in `lower.go` can be merged:
- `lowerTopLevel` and the new nested-boc lowering become one path
- `lowerBocBody`, `lowerClosureBody`, `lowerBocAsStmts` can be reviewed for consolidation
- The `localBocVars` tracking and ad-hoc `var f any` hacks can be removed

## Risks and Unknowns

- **Go local type restriction**: Go does not allow methods on locally-declared types. Lifting to package level means the struct name must be globally unique — FQN-based naming solves this but adds naming complexity.
- **Uppercase outer boc with lowercase inner boc**: `Foo: { bar: { count: 0 } }` — `bar` should be a per-`Foo`-instance singleton, so `_foo_barBoc` is instantiated inside `NewFoo()`, not as a package-level var. This is the main complexity difference from file-scope singletons.
- **Recursive local bocs**: `f: { n Int; f(n-1) }` — the pre-registration fix in sema already handles the name resolution; the lowerer needs to handle the function-var self-reference correctly.
- **Cross-scope method calls**: `boc.baz()` from an outer scope where `boc` is a local var — sema needs to track that `boc`'s type is `_main_bocBoc`, not `any`, so member access resolves correctly.
- **File/directory bocs**: The directory-as-boc and file-as-boc levels are not yet implemented at all. This design unification is a prerequisite for them making semantic sense.

## Relation to Existing Open Items

This work touches or resolves several items in `task.md`:

- **"Top-level boc callable as function"** — subcase of the uniformity problem
- **"SWMR write semantics"** — still deferred; this work doesn't address actor queues, only struct emission
- **"Standalone thunk calls inside closure bodies not forced"** — the lowerClosureBody path cleanup in Pass 3 would be the right time to fix this

## What NOT to Change Yet

- **Directory and file bocs** — defer until the in-file nesting case works correctly first
- **Actor queue / SWMR** — separate concern; struct emission can land without it
- **FQN cross-file references to nested bocs** — defer until single-file case is solid
