# Yz Compiler Implementation

## Phase 0 — Project Setup
- [x] Create `compiler/` directory skeleton
- [x] Initialize `go.mod` (`module yz`)
- [x] Create `cmd/yzc/main.go` (CLI with `build`, `run`, `new`)
- [x] Create `Makefile` (`build`, `test`, `clean`)
- [x] Create `compiler/README.md`
- [x] Verify: `go build ./...` passes

## Phase 1 — Lexer
- [x] `internal/token/token.go` — token types
- [x] `internal/lexer/lexer.go` — tokenizer + ASI
- [x] `internal/lexer/lexer_test.go` — 38 tests, all passing

## Phase 2 — Parser
- [x] `internal/ast/ast.go` — AST node types
- [x] `internal/parser/parser.go` — recursive descent
- [x] `internal/parser/parser_test.go` — 32 tests, all passing

## Phase 3 — Semantic Analysis
- [x] `internal/sema/analyzer.go` — scope, type inference, boc/struct dispatch
- [x] `internal/sema/analyzer_test.go` — tests passing

## Phase 4 — IR
- [x] `internal/ir/ir.go` — IR node type definitions
- [x] `internal/ir/lower.go` — AST+sema → IR lowerer
- [x] `internal/ir/ir_test.go` — 8 tests, all passing

## Phase 5 — Code Generation
- [x] `internal/codegen/codegen.go` — Go source emitter
- [x] `internal/codegen/codegen_test.go` — 10 tests, all passing
- [x] `cmd/yzc/build.go` — full pipeline: parse→sema→IR→codegen→go build
- [x] `cmd/yzc/new.go` — project scaffolding

## Phase 6 — Runtime Library
- [x] `runtime/yzrt/types.go` — Int, Decimal, String, Bool, Unit with symbol-named methods
- [x] `runtime/yzrt/thunk.go` — Thunk[T], Go[T] (goroutine spawn), Force()
- [x] `runtime/yzrt/collections.go` — Array[T], Dict[K,V], Range
- [x] `runtime/yzrt/core.go` — Print, While, Info, BocGroup (structured concurrency)
- [x] `runtime/yzrt/yzrt_test.go` — tests passing

## Phase 7 — Integration & Testing
- [x] `compiler/test/conformance/` — golden tests, 18 passing
- [x] `compiler/examples/` — counter, milestone (concurrent fetch + counter)
- [x] Error tests — 7 cases: parse errors, undefined variable/type, mix undefined/conflict/not-struct

## Language Features — Implemented
- [x] `while` loop
- [x] `BocWithSig` — top-level functions and methods inside singleton/struct bocs
- [x] `match` expression (condition form)
- [x] `mix` statement — Go embedding
- [x] Multi-file projects — flat and subdirectory (cross-package FQN)
- [x] Type-only BocWithSig — `Name #(params)`: data params → struct (no constructor); all-boc params → Go interface (structural typing)
- [x] `http` built-in singleton — `http.get(uri)`, `http.post(uri, body)`
- [x] First milestone — concurrent HTTP fetch + counter (`examples/milestone/`)

## Language Features — Implemented (continued)
- [x] Variant/discriminant sum types — `Pet: { Cat(...), Dog(...) }` with per-variant constructors
- [x] Discriminant match — `match expr { Cat => body }, { Dog => body }` → Go switch
- [x] Cross-package singleton method calls — `pkg.singleton.method()`
- [x] `yzc run` — compile + execute in one step
- [x] `http` built-in singleton — `http.get(uri)`, `http.post(uri, body)`
- [x] thunk transparency — `a: boc.call()` auto-forced on use

## Language Features — Not Yet Implemented
- [x] Mixed type-only decl — `Name #(name String, greet #())` → struct with data fields + function-typed fields + method wrappers
- [x] `BocWithSig` body-only form — `name #(params) = { body }` — named and anonymous param matching
- [x] Error reporting — Rust-style diagnostics with source context and caret underlines

## BocWithSig Body-Only — Deferred
- [x] Default values in params — `#(name String = "hello")` — injected at call sites; golden test 21
- [x] `ShortDecl` as param — `name : "default"` in sig — type inferred from default; golden test 22
- [x] Generic variant types — `Option: { V; Some(value V); None() }` with `[V any]` on struct and constructors; discriminant match works; golden test 25
- [x] Generic type vars in sig — `identity #(value V, V)` → `func identity[V any](value V) *Thunk[V]`; golden test 26
- [x] Uninstantiated generics — `Option(String)` → `*Option[std.String]` in type positions
- [x] Declare-only then assign-later — `greet #(name String)` then `greet = { name String; … }` → FuncDecl; golden test 23

## Language Features — Already Implemented (discovered)
- [x] Multiline strings — strings span lines naturally; `"` or `'` closes on any line (lexer handles `\n` inside string literals)

## Language Features — Implemented (continued)
- [x] HOF / closures as arguments — `list.filter({ item Int; item > 10 })` — sync closures with typed params; `Array.Filter`, `Array.Each`, `ArrayMap`; golden test 27

## Generics — Future Work
- [x] HOF: `list.map({ item Int; item * 2 })` — `lowerCall` detects `.map(boc)` on ArrayType → emits `std.ArrayMap(recv, closure)`; result type inferred via `:=`; golden test 28
- [x] Generic constraint inference (Option 4) — sema scans method bodies for T-method calls; records constraints; checks all at instantiation; reports all missing methods at once; error test 09; golden test 31 (generic method receiver)
- [x] Go constraint generation — emit `[T interface{ ToStr() std.String }]` from inferred constraints; lowerMethodName fixes to_string→ToStr; golden test 32
- [x] Multiple type params — `#(key K, value V)` → `[K any, V any]`; Pair[K,V] struct + makePair[K,V] function; parser fix: TYPE_IDENT'(' only parsed as VariantDef in type boc bodies (inTypeBoc flag); golden test 33
- [x] Generic structs (non-variant) — `Box: { T; value T }` → `type Box[T any] struct { value T }`; golden test 29
- [x] Typed generic declaration — `b Box(String) = Box("hello")` → `var b *Box[std.String] = NewBox(...)`; golden test 30; TypedDecl in lowerMainStmt
- [x] Optional parens for trailing-block calls — `list.filter { block }` without `()`; in `parsePostfix`, LBRACE after MemberExpr → CallExpr with BocLiteral arg; golden test 34
- [x] Unary minus on variables — `-x` → `x.Neg()`; `a - -b` → `a.Minus(b.Neg())`; pipeline was already wired (parser+sema+lowerer+codegen); golden test 35

## Boc Uniformity — Major Architectural Work

See full analysis in [`boc_uniformity.md`](boc_uniformity.md).

The compiler currently has three separate lowering paths for bocs depending on where they appear (file-scope, BocWithSig, local/nested). The intended design is one uniform construct: a boc is a boc regardless of nesting depth. Local bocs should produce the same struct+method+concurrency output as file-scope bocs.

- [ ] **Pass 1 — Sema: uniform type recording for nested bocs** — `analyzeBocDecl` should produce a `StructType` (not just `BocType`) for lowercase local bocs that contain inner bocs or BocWithSig methods; FQN registration should mirror file-scope behavior.

- [ ] **Pass 2 — Lowerer: lift nested boc structs to package level** — introduce a pre-pass that collects all boc declarations at any nesting depth; emit `_fqnBoc` struct + methods at package level; emit instance creation (`&_fqnBoc{}`) at the point of declaration in the enclosing function body.

- [ ] **Pass 3 — Unify lowering paths** — after Passes 1+2, merge `lowerTopLevel` / `lowerBocBody` / `lowerClosureBody` / `lowerBocAsStmts` into a single path; remove `var f any` hacks and `localBocVars` tracking.

- [ ] **Directory and file bocs** — defer until in-file nesting works; then extend the FQN tree to cover files and directories as bocs.

- [ ] **Multiple source roots** — a project may declare more than one source root (e.g. `src/` for app code, `lib/` for third-party). Each root is an independent FQN mount point; names inside `src/foo/bar.yz` resolve as `foo.bar.*`, names inside `lib/baz.yz` resolve as `baz.*`. No cross-root FQN collision is possible. Tooling convention (versioning, lock files) deferred; the compiler needs to accept a list of source roots and build one FQN forest per root.

## Language Design — Open Questions (tracked in Questions/)

- [ ] **Cancellation / non-local return across goroutine boundaries** — non-local `return` from a callback conflicts with structured concurrency. Three open sub-problems: goroutine leaks when a race-return fires, escaped non-local returns into completed bocs, and structured concurrency violation. See `Questions/How to cancel a running block.md`. No implementation work until the design question is resolved.

- [ ] **Stateless bocs and pure functions** — BocWithSig form (`foo #(params) { ... }`) is the stateless function form; body form (`foo: { ... }`) is the stateful actor form. The compiler needs to enforce: (1) `foo.field` is a type error on a BocWithSig boc; (2) passing a stateless boc where a named-param signature type is expected (e.g., `#(name String, Int)`) is a type error; (3) BocWithSig calls emit free goroutines (no actor queue). See `Questions/Stateless bocs and pure functions.md` and `Features/Bocs.md`.

- [ ] **SWMR write semantics in codegen** — field writes from outside a boc (`a.b = v` in a different boc) should be emitted as queued actor messages, not direct struct field assignments. Currently the codegen emits direct field writes which is a data race. Requires runtime support for a write-message channel per boc instance. Depends on the cancellation/actor-queue design.

## Known Bugs
- [x] Dict literals — fixed: now emits `std.NewDict[K,V]().Set(k,v)...` chain; golden test 24
- [x] Array literals — already worked via variadic `std.NewArray(...)`; golden test 24

- [ ] **Assigning a Unit-returning boc to a variable** — `a : foo()` where `foo` returns nothing (Unit) should be a sema error, analogous to Go's `x := f()` where `f` returns nothing. Detect in sema: if a `ShortDecl` or `TypedDecl` RHS resolves to a boc call whose return type is `Unit`, report an error. Add an error golden test.

- [ ] **Top-level boc callable as function** — `foo: { time.sleep(1); "done" }` at the top level is currently lowered as a singleton struct (`*_fooBoc`), which is not callable. When a top-level boc has a body that returns a value and is invoked as `foo()`, it should be lowered as a Go function instead. Fixing this also fixes the fire-and-forget issue for that code path, since the body would go through `lowerBocBody` (which uses `BocGroup` + `Wait` for standalone thunk calls) rather than `lowerClosureBody`. Needs sema and lowerer changes; add a golden test.

- [ ] **Standalone thunk calls inside closure bodies not forced** — `time.sleep(1)` as an expression statement inside a local boc (closure body) is emitted as a raw `std.Time.Sleep(...)` call with no `.Force()` or `BocGroup` wrapping. The goroutine fires but nothing waits for it, so the sleep is effectively skipped. Fix: in `lowerClosureBody`, detect standalone ExprStmt calls that return a thunk (via `isBocMethodCall`) and emit `.Force()` on the result (or wrap in a `BocGroup` + `Wait()` like `lowerMainBoc` does). Add a golden test.

- [ ] **Unused variables in generated Go** — Yz allows unused variables but Go does not. Fix: after lowering all statements in a scope (main boc, method body), scan the emitted IR for declared variable names (`DeclStmt.Name`) that never appear as `Ident` references in subsequent IR nodes. Append `_ = varName` (`AssignStmt` with blank target) for each unused name. Applies to `lowerMainBoc`, `lowerBocBody`, and `lowerClosureBody`. No change to sema or parser; pure IR post-processing. Add a golden test with a declared-but-unused variable.

## Documentation Gaps — Features Documented but Not Yet Implemented

These are documented in the language spec/features and need compiler implementation:

- [ ] **`break` / `continue` / `return` in loops** — spec defines these keywords; lexer/parser may tokenize them but lowerer/codegen don't emit them yet. Needed for `while` loops, `each` callbacks, and early exit from bocs.

- [ ] **Range iteration** — `1.to(10).each({ i Int; ... })` — `Range.Each` exists in yzrt and `Int.To` returns a `Range`, but the lowerer doesn't recognize `.each(closure)` on a Range value (only on Array). Need HOF lowering for Range receivers.

- [ ] **Named arguments in constructor calls** — `Person(name: "Alice", age: 30)` — parser likely handles labeled args but lowerer may not reorder them to match struct field order. Add golden test.

- [ ] **Multiple return values** — `x, y = swap(x, y)` — multiple assignment on LHS is documented; not in any golden test. Requires parser and lowerer support for multi-assign statements.

- [ ] **Array append via `<<`** — `a << item` as sugar for `a.Append(item)` via non-word method invocation. `Array.Append` exists in yzrt. Need a golden test and lowerer to emit the `Append` call.

- [ ] **Option/Result method chaining** — `result.or_else({ error Error; ... })`, `result.and_then({ val T; ... })` — documented in error-handling features. Requires implementing `or_else`, `and_then`, `or` methods on the Option/Result types in yzrt, plus lowerer support for chained calls on variant types.

- [ ] **`to_str()` method on user types** — examples use `n.to_string()` but yzrt uses `ToStr()` (mapped from `to_str()`). Ensure the compiler correctly maps `to_str()` calls on user-defined types; update examples to use `to_str()` not `to_string()`.

- [ ] **String concatenation with `++`** — `"hello" ++ " " ++ "world"` — `String.Plus` exists in yzrt; need a golden test to confirm the codegen path works end-to-end.

- [ ] **Dict Optional access** — `d[key]` should return `Option(V)` per spec; currently returns `V` directly (panics on missing key via `At()`). Needs yzrt change + codegen update.

- [ ] **Bool methods `&&` / `||`** — `Bool.Ampamp` and `Bool.Pipepipe` exist in yzrt; confirm they are wired through the operator lowering path (codegen for `&&`/`||` binary expressions). Add golden test.
 
- [ ] **Info strings** — `` `"doc string"` `` before a declaration; retrievable via `info(var).text` at runtime. The lexer captures info strings as AST nodes, and `yzrt.Info()` exists, but codegen doesn't attach info strings to declarations or emit `Info()` calls. See `Features/Info strings.md`.

- [x] **Explicit type on boc-call declarations** — `c String = http.get(url)`: fixed in `lowerMainStmt`, `lowerBocBody`, and `lowerClosureBody` — detect `isBocMethodCall` on the TypedDecl value and use inferred `:=` + `thunkVars`, same as `ShortDecl`; golden test 36.
- [x] **Local boc variable with explicit boc-type** — `foo #(String) = { "hello" }` and `foo #(String) { "hello" }` inside a boc body: sema now falls back to shorthand semantics (String = return type) when the body-only form has no TypedDecl params; lowerer emits a local function literal `func() *Thunk[String] { return std.Go(...) }` tracked in `localBocVars` so calls are emitted directly (not double-wrapped) and results are auto-forced; golden test 37.
