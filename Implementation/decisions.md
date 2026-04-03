# Yz Compiler — Implementation Decisions

Resolved progressively during planning sessions (2026-03-04 through 2026-04-03).

---

## Syntax & Lexer

| # | Decision | Resolution |
|---|----------|------------|
| 1 | ASI rules | Only insert `;`, never commas. Commas are explicit in `()` and `[]`. |
| 2 | Multi-line strings | Strings are multi-line by default. No raw/heredoc syntax needed. |
| 3 | `=` semantics | `=` is **not** an expression — it is a statement. |
| 4 | `:` semantics | `:` (short declaration) is a statement, not an expression. |
| 15 | `=>` as non-word | `=>` alone (standalone, not preceded by other non-word chars) is the FAT_ARROW delimiter and **cannot** be used as a non-word identifier. However, `=>` appearing mid-sequence (e.g. `!=>`, `<=>`) **is** part of a valid non-word identifier. Spaces determine token boundaries, just like word identifiers: `f oo` is two idents, `foo` is one. |
| 16 | Info string delimiters | Info strings use regular string delimiters (`'...'` or `"..."`) — **not** backtick. Backtick is only for string interpolation inside a string. Info strings are parsed and kept in the AST; no runtime code is generated for them (tooling support deferred). |

## Assignment & Destructuring

| # | Decision | Resolution |
|---|----------|------------|
| 5 | Multiple assignment scope | Support both simple multi-return (`a, b = foo()`) **and** nested destructuring with parentheses (`a, (b, c) = foo(), bar()`). |
| 17 | Grouping vs tuple | `(a, b)` as an expression is **grouping only** — a single-element group is valid, multi-element is not (no tuples). `lhs = (a, b)` is invalid syntax. |

## Boc Semantics

| # | Decision | Resolution |
|---|----------|------------|
| 6 | Lowercase boc state | State is a **singleton per fully-qualified name (FQN)**, not global. The boc body runs once; state persists across invocations of the same FQN. Example: `counter.increment(); counter()` → returns `1`. Each distinct FQN is its own singleton. |
| 18 | Boc FQN | A boc's FQN is determined by its **source path + nesting**. `src/a.yz` containing `b: {}` creates FQN `a.b`. A file at `src/a/b.yz` also creates FQN `a.b`. Both can coexist within the same source root (merging, deferred to phase 2). Across **different source roots**, the same FQN defined in two roots is a compilation error — the first source root wins. |
| 19 | BocWithSig special form | `greet #(n String) { body }` is a special form where params declared in the signature are **directly available inside the body without redeclaration**. This is syntactic sugar for declaring and initializing in one step. The body does NOT need to re-declare `n`. Contrast: `greet #(n String) = { n String /* must redeclare */ }`. The `=` form requires the body to declare the params itself. Resolving which params are in scope inside the body is a **sema concern**, not a parser concern. Same rule applies to uppercase type bocs. |
| 20 | Comma in boc bodies | Commas act as statement separators in boc bodies (alongside semicolons/ASI). This supports the `T, E,` generic param list syntax in type bocs and separators between match arms. |

## Concurrency Model

| # | Decision | Resolution |
|---|----------|------------|
| 21 | All boc invocations | **Every boc invocation runs in a goroutine.** There is no special async/thunk syntax — all calls are non-blocking by default. |
| 22 | Lazy thunk | The result of a boc invocation is a **lazy thunk** that materializes on first use. `result: add(1, 2)` is non-blocking; `result` materializes when it is first needed (e.g. passed to `print`, used in arithmetic). |
| 23 | Structured concurrency | A boc is not considered complete until **all child bocs it spawned have completed**. This provides implicit structured concurrency: the parent eventually waits for all descendants, even though the invocations are non-blocking. IO and network operations naturally trigger thunk materialization. |

## Entry Point

| # | Decision | Resolution |
|---|----------|------------|
| 24 | Entry point | The entry point is a **boc named `main`**. By convention this lives in `main.yz`, but only because `main.yz` defines the `main` boc — the filename is not magic. |

## Implementation Strategy

| # | Decision | Resolution |
|---|----------|------------|
| 7 | Compiler language | **Go** |
| 8 | Compiler architecture | Emit **Go source code** → then `go build` (not AST, not interpreter) |
| 9 | v0.1 milestone scope | **Full** — all spec features |
| 10 | Repository | `compiler/` directory in this repo (`yz-docs`). Go module: `module yz` (local-only, not go-getable). |
| 25 | Working style | Component-by-component with TDD: implement one complete logical unit + tests, review, then proceed. Tests are written first (red → green). Trust automated tests; verify generated code compiles and tests pass. |
| 26 | First milestone | A program using **concurrency**: fetching two resources concurrently and a counter boc. This exercises bocs, goroutines, thunks, structured concurrency, and arithmetic. |
| 27 | Conformance tests | Created **incrementally** as each phase lands (not upfront). Use `go test` framework. |

## Runtime Semantics

| # | Decision | Resolution |
|---|----------|------------|
| 11 | Unhandled errors | **Panic** (crash) |
| 12 | `nil` | **No `nil` concept.** `nil` is a valid identifier but has no special meaning. Use `Result`/`Option` instead. |
| 13 | `&&`/`||` short-circuit | They are regular methods, **not** compiler-special-cased. Short-circuit behavior is natural because the argument is a **boc** (lazy). e.g., `a \|\| { expensive() }` — the boc is only called if `a` is false. |
| 14 | Dependency config | No TOML. Use a **`.yz` file** (e.g., `project.yz`) implementing a project interface. Minimal format TBD. |
| 28 | Type representation | **All types are structs** in the generated Go runtime — there are no primitive Go types. `age Int` in Yz compiles to a field of type `std.Int` (a Go struct). Literals are also boxed: `1` in source code becomes `std.NewInt(1)` in generated Go. |
| 29 | Literal boxing | Boxing of literals (e.g. `1` → `std.Int`) is done by the **codegen phase**, not the runtime. |
| 30 | Standard library naming | The standard library types are named **`std.Int`**, **`std.Decimal`**, **`std.String`**, **`std.Bool`**, **`std.Unit`** — not `YzInt`, `YzDecimal`, etc. The runtime Go package is `yz/runtime/yzrt`; the types it exports use the `std.*` naming convention from the Yz perspective. |

## Code Generation

| # | Decision | Resolution |
|---|----------|------------|
| 31 | Non-word method naming | Non-word method names are mapped to Go-safe identifiers using the **symbol name** of each character — not semantically meaningful names. Examples: `+` → `plus`, `++` → `plusplus`, `?` → `qm`, `==` → `eqeq`, `!=` → `neq`, `&&` → `ampamp`, `\|\|` → `pipepipe`, `<=` → `lteq`. This is mechanical, not interpretive. |
| 32 | Build output directory | Generated `.go` files go to `target/gen/` inside the **project being compiled** (not the compiler directory). Binary goes to `target/bin/app`. The `target/` directory is added to `.gitignore` in `yzc new`-generated projects. |
