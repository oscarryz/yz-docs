#impl 
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

- [x] **Pass 1 — Sema: uniform type recording for nested bocs** — `analyzeBocDecl` produces `StructType{IsSingleton:true}` for lowercase bocs with inner structure; `analyzeStructBoc` returns return types; `lowerName` and `isSingletonBoc` accept singleton StructType; `TestSingletonBocFieldAccess` confirms field access from outside boc body resolves correctly.

- [x] **Pass 2 — Lowerer: lift nested boc structs to package level** — `liftLocalBoc` unified handler in `lower.go`: detects body-form bocs (`f: { n Int; ... }`) and BocWithSig bocs (`foo #(String) { ... }`) inside `main`; emits `_main_fBoc` struct + `Call(params)` method at package level; emits `_f := &_main_fBoc{}` local instance; recursive self-calls inside `Call()` emit `self.Call(args)`. Goldens 37 and 39 updated: struct + BocGroup for body-form, struct + forced call for BocWithSig.

- [x] **Pass 3 — Unify lowering paths** — all lowercase bocs now emit singleton structs uniformly; `main` is no longer special-cased; `lowerBodyOnlySingleton` handles pure-statement bocs (BocType from sema), `lowerStructuredSingleton` handles bocs with fields/methods (StructType{IsSingleton:true}) and now also emits a `Call()` method for any statement elements; `recvMethods` tracking enables unqualified method calls (e.g. `f(n-1)` → `self.F(n-1)`) inside method bodies and `Call()`; `func main()` is auto-emitted as a shim (`Main.Call().Force()`); all 39 goldens updated. `localBocVars` and `liftLocalBoc` are removed from the `main` path (local named bocs inside a singleton are now proper methods). The HOF callback / while predicate question (goldens 27, 34, 05) was also updated — these still use Go closures for anonymous boc literals in call arguments, which is the correct behavior since the design question only applies to NAMED local bocs.

- [x] **Pass 4 — Top-level BocWithSig as singleton** — `lowerBocWithSig` now emits `*SingletonDecl` (struct + `Call(params...)` method) instead of `*FuncDecl` for non-generic BocWithSig declarations. `lowerTopAssignment` (declare-then-assign pattern) also emits `SingletonDecl`. Call sites emit `Singleton.Call(args)` via updated `lowerCall`; `lowerName` no longer excludes BocWithSig from singleton capitalization. BocType params (e.g. `cond #(Bool)`) become `func() T` struct fields, called as `self.cond()` (not method calls). Generic BocWithSig keeps `FuncDecl` until generic singletons are designed. `go test -race ./...` passes; all 40 goldens regenerated.

- [ ] **Directory and file bocs** — defer until in-file nesting works; then extend the FQN tree to cover files and directories as bocs.

- [ ] **Multiple source roots** — a project may declare more than one source root (e.g. `src/` for app code, `lib/` for third-party). Each root is an independent FQN mount point; names inside `src/foo/bar.yz` resolve as `foo.bar.*`, names inside `lib/baz.yz` resolve as `baz.*`. No cross-root FQN collision is possible. Tooling convention (versioning, lock files) deferred; the compiler needs to accept a list of source roots and build one FQN forest per root.

## Language Design — Open Questions (tracked in Questions/)

- [ ] **Cancellation / non-local return across goroutine boundaries** — non-local `return` from a callback conflicts with structured concurrency. Three open sub-problems: goroutine leaks when a race-return fires, escaped non-local returns into completed bocs, and structured concurrency violation. See `Questions/How to cancel a running block.md`. No implementation work until the design question is resolved.

- [x] **SWMR write semantics in codegen** — field writes from outside a boc (`a.b = v` in a different boc) are wrapped in `std.Schedule(&target.Cown, func() std.Unit { target.field = value; return std.TheUnit }).Force()`. Implemented in Phase C; golden test 42.

- [x] **Phase D — Extend cowns to struct type instances** — struct bocs (`Foo: { ... }`) and their instances (`bar : Foo()`) must embed `std.Cown` just like singleton bocs do. Currently `emitStructDecl` omits `std.Cown`, so method bodies that reference `self.Cown` fail at `go build`. Also, multi-cown detection in `lowerBocWithSigAsSingleton` only recognizes singleton params, not regular struct params — so `transfer(src Account, dst Account, ...)` doesn't acquire both cowns atomically. Fix: (1) add `std.Cown` as first field in `emitStructDecl`; (2) extend multi-cown detection to include all non-interface struct types; (3) regenerate golden tests; (4) add a new golden test for struct-instance method concurrency (e.g. `account_balance`).

## Known Bugs
- [x] Dict literals — fixed: now emits `std.NewDict[K,V]().Set(k,v)...` chain; golden test 24
- [x] Array literals — already worked via variadic `std.NewArray(...)`; golden test 24

- [ ] **Assigning a Unit-returning boc to a variable** — `a : foo()` where `foo` returns nothing (Unit) should be a sema error, analogous to Go's `x := f()` where `f` returns nothing. Detect in sema: if a `ShortDecl` or `TypedDecl` RHS resolves to a boc call whose return type is `Unit`, report an error. Add an error golden test.

- [ ] **Top-level boc callable as function** — `foo: { time.sleep(1); "done" }` at the top level is currently lowered as a singleton struct (`*_fooBoc`), which is not callable. When a top-level boc has a body that returns a value and is invoked as `foo()`, it should be lowered as a Go function instead. Fixing this also fixes the fire-and-forget issue for that code path, since the body would go through `lowerBocBody` (which uses `BocGroup` + `Wait` for standalone thunk calls) rather than `lowerClosureBody`. Needs sema and lowerer changes; add a golden test.

- [x] **Standalone thunk calls inside closure bodies not forced** — fixed in `lowerClosureBody`: non-last ExprStmt where `isBocMethodCall` is true now wraps the lowered expr in `ForceExpr`; golden test 40.

- [ ] **Unused variables in generated Go** — Yz allows unused variables but Go does not. Fix: after lowering all statements in a scope (main boc, method body), scan the emitted IR for declared variable names (`DeclStmt.Name`) that never appear as `Ident` references in subsequent IR nodes. Append `_ = varName` (`AssignStmt` with blank target) for each unused name. Applies to `lowerMainBoc`, `lowerBocBody`, and `lowerClosureBody`. No change to sema or parser; pure IR post-processing. Add a golden test with a declared-but-unused variable.

## Documentation Gaps — Features Documented but Not Yet Implemented

These are documented in the language spec/features and need compiler implementation:

- [ ] **`break` / `continue` / `return` in loops** — spec defines these keywords; lexer/parser may tokenize them but lowerer/codegen don't emit them yet. *(Full semantics now in Breaking Change item 2)*

- [ ] **Range iteration** — `1.to(10).each({ i Int; ... })` — `Range.Each` exists in yzrt and `Int.To` returns a `Range`, but the lowerer doesn't recognize `.each(closure)` on a Range value (only on Array). Need HOF lowering for Range receivers.

- [ ] **Named arguments in constructor calls** — `Person(name: "Alice", age: 30)` — parser likely handles labeled args but lowerer may not reorder them to match struct field order. Add golden test.

- [ ] **Multiple return values** — `x, y = swap(x, y)` — multiple assignment on LHS is documented; not in any golden test. Requires parser and lowerer support for multi-assign statements.

- [ ] **Array append via `<<`** — `a << item` as sugar for `a.Append(item)` via non-word method invocation. `Array.Append` exists in yzrt. Need a golden test and lowerer to emit the `Append` call.

- [ ] **Option/Result method chaining** — `result.or_else({ error Error; ... })`, `result.and_then({ val T; ... })` — documented in error-handling features. Requires implementing `or_else`, `and_then`, `or` methods on the Option/Result types in yzrt, plus lowerer support for chained calls on variant types.

- [ ] **Non-word boc names** — bocs (methods on structs, standalone bocs) whose names are non-word identifiers: `balance+= #(amount Int) { ... }`, `balance-=`, `hola++`, etc. The lexer already produces `NON_WORD` tokens for these character sequences, but the parser currently only allows word identifiers as boc names in declarations. Fix: (1) parser: accept a `NON_WORD` token as a valid boc name in `ShortDecl` and `BocWithSig` forms; (2) lowerer: map the non-word name to a valid Go method name using the same symbol table used for operator methods (e.g. `+=` → `PlusEq`, `++` → `PlusPlus`); (3) call sites: `obj.balance+=(x)` must lower to `obj.BalancePlusEq(x)`; (4) add a golden test (e.g. the `account_balance` example). Needed to unblock `examples/account_balance/main.yz`.

- [ ] **`to_str()` method on user types** — examples use `n.to_string()` but yzrt uses `ToStr()` (mapped from `to_str()`). Ensure the compiler correctly maps `to_str()` calls on user-defined types; update examples to use `to_str()` not `to_string()`.

- [ ] **String concatenation with `++`** — `"hello" ++ " " ++ "world"` — `String.Plus` exists in yzrt; need a golden test to confirm the codegen path works end-to-end.

- [ ] **Dict Optional access** — `d[key]` should return `Option(V)` per spec; currently returns `V` directly (panics on missing key via `At()`). Needs yzrt change + codegen update.

- [ ] **Bool methods `&&` / `||`** — `Bool.Ampamp` and `Bool.Pipepipe` exist in yzrt; confirm they are wired through the operator lowering path (codegen for `&&`/`||` binary expressions). Add golden test.
 
- [ ] **Info strings** — `` `"doc string"` `` before a declaration; retrievable via `info(var).text` at runtime. The lexer captures info strings as AST nodes, and `yzrt.Info()` exists, but codegen doesn't attach info strings to declarations or emit `Info()` calls. See `Features/Info strings.md`. *(Superseded by Breaking Change item 3 — infostrings as boc bodies; that item covers the full redesign)*

- [x] **Explicit type on boc-call declarations** — `c String = http.get(url)`: fixed in `lowerMainStmt`, `lowerBocBody`, and `lowerClosureBody` — detect `isBocMethodCall` on the TypedDecl value and use inferred `:=` + `thunkVars`, same as `ShortDecl`; golden test 36.
- [x] **Local boc variable with explicit boc-type** — `foo #(String) = { "hello" }` and `foo #(String) { "hello" }` inside a boc body: sema now falls back to shorthand semantics (String = return type) when the body-only form has no TypedDecl params; lowerer emits a local function literal `func() *Thunk[String] { return std.Go(...) }` tracked in `localBocVars` so calls are emitted directly (not double-wrapped) and results are auto-forced; golden test 37.

---

## Yz 0.2.0

These items reflect decisions documented in `docs/Implementation/breaking-changes.md`. Each one requires changes across spec, compiler, and/or runtime.

### 1. String Interpolation: `${}` instead of backtick — COMPLETE

The old interpolation syntax `` `expr` `` inside strings is replaced by `${expr}`. Backtick is now reserved exclusively for infostrings. See `docs/Features/String interpolation.md`.

- [x] **Lexer** — already recognizes `${` / `}` as interpolation delimiters; backtick only starts an infostring (outside strings)
- [x] **AST** — `InterpolatedStringExpr` comment updated to reflect `${}` syntax
- [x] **Golden tests** — no conformance `.yz` files used backtick interpolation; examples updated to `${}`
- [x] **Spec 01** — updated: §1.10 and §1.12 ASI example now use `${}` syntax; all backtick-interpolation examples in spec files 01–08 replaced

### 2. `return`, `break`, `continue` — frontend up to sema

The spec (`docs/Features/return, break, continue.md`) now defines precise semantics. See also `docs/Features/Language Primitives.md`.

Codegen is deferred: `break`/`continue` only looked simple because `tryLowerWhile` was emitting Go `for` loops (see item 9). Once `while` goes through the normal recursive boc path, `break`/`continue` become non-local exits from a recursive call chain — the same class of hard problem as `return` through anonymous boc boundaries. None of the three have a clean codegen path until the concurrency model and recursive boc semantics are settled.

- [ ] **Parser** — parse `break` and `continue` as `BreakStmt` / `ContinueStmt` AST nodes (tokens already exist)
- [ ] **Sema** — validate context: `break`/`continue` only inside a loop body; `continue` in a match arm only inside a match; `return` tracks the nearest enclosing named boc
- [ ] **Lowerer** — emit `panic("not yet implemented")` or a compile error for all three when encountered, so they fail loudly rather than being silently dropped
- [ ] **Spec 07** — update control-flow spec to match semantics in `docs/Features/return, break, continue.md`
- [ ] **Golden tests** — add sema-level error tests for `break`/`continue` outside a loop, `return` type mismatch

### 3. Infostrings — content is a boc body

The infostring delimiter stays backtick, but its content is now a **boc body** (full Yz syntax, parsed and type-checked, never executed). Currently the AST stores the infostring as a plain `*StringLit`. See `docs/Features/Info strings.md`.

- [ ] **AST** — change `InfoString` to hold a `*BocLiteral` (or `[]ast.Node` elements) instead of `*StringLit`
- [ ] **Lexer** — infostring content needs to be re-lexed as Yz source, not treated as a string value; currently the lexer emits the raw string content
- [ ] **Parser** — after lexing the backtick span, re-parse its content as a boc body using the existing boc-body parser
- [ ] **Sema** — type-check infostring content; validate that referenced names (e.g. in `compile_time: [Derive]`) resolve to known bocs; report missing types at compile time
- [ ] **Codegen** — attach the compiled infostring boc to the declaration's metadata in generated code (prerequisite for compile-time bocs)
- [ ] **Spec 01** — update lexical structure to clarify backtick = infostring delimiter, content = boc body (not string value)

### 4. Generics — Explicit Constraint Declaration

Type parameters can now carry an explicit constraint: `thing T Talker` declares that `T` must implement `Talker`. Constraint inference from usage still works; explicit declaration is additive. See `docs/Features/Generics - Type Parameters.md`.

- [ ] **Parser** — in param lists and boc bodies, after a single-uppercase-letter type param, allow an optional uppercase-name constraint: `T Constraint`; produce a new `TypeParamDecl{Name: "T", Constraint: "Talker"}` AST node (or extend the existing generic param representation)
- [ ] **Sema** — when an explicit constraint is present, validate at instantiation that the concrete type satisfies it; combine with inferred constraints from usage (union, not replacement)
- [ ] **Error messages** — surface explicit constraint violations clearly, distinct from inferred violations
- [ ] **Spec 04** — update type-system spec to document `T Constraint` syntax and its interaction with inference

### 5. `:` as Type Alias

`Name : SomeType` (or `Name : #(...)`) declares a type alias — a new name for an existing type. This is a general-purpose feature usable anywhere, not exclusive to boc bodies or compile-time contexts. (Feature doc to be created.) The `Schema : #(...)` pattern in `Compile` implementations is one use of this.

- [ ] **Feature doc** — user to create `docs/Features/Type Alias.md` describing `: ` alias syntax and scoping rules
- [ ] **Parser** — distinguish `Name : TypeExpr` (alias) from `Name TypeExpr` (typed declaration) and from `name : value` (short decl); uppercase `Name` + `:` + type expression = alias
- [ ] **Sema** — register alias in current scope; resolve references to the alias name as the aliased type; aliases do not produce runtime fields or values
- [ ] **Lowering** — emit Go type alias (`type Name = GoType`) at the appropriate scope; do not emit struct fields for alias declarations
- [ ] **Spec 04** — add type alias syntax and scoping rules

### 6. Compile-Time Bocs (`Compile` interface)

Major new feature. Any boc with `Schema #()` and `run #(Boc, Boc)` satisfies the `Compile` interface. `compile_time: [Impl, ...]` in a boc's infostring schedules those implementations to run during type inference. Their return values are merged into the parent boc. See `docs/Features/Compile Time Bocs.md`.

Depends on: items 3 (infostring as boc body), 4 (explicit constraints), 5 (associated types).

- [ ] **Sema** — recognize the `Compile` structural interface (duck-typed: any boc with `Schema #()` and `run #(Boc, Boc)`)
- [ ] **Sema** — when resolving a boc definition: scan its infostring for `compile_time: [...]`; resolve each name; schedule them in array order during type inference
- [ ] **Boc metatype** — implement the `Boc` value type accessible inside `run`: `{name String, instantiable Bool, fields [Boc], methods [Boc], type_params [Boc], infostring Boc, source #()}`; the compiler must populate it for every boc definition
- [ ] **Two-phase build** — Phase 1: scan all source for `Compile` implementations; compile them to native executables first. Phase 2: main compilation; when inference hits a `compile_time` trigger, call the pre-compiled executable via subprocess, serialize the current `Boc` in, deserialize the returned `Boc`, merge into AST, re-enter inference
- [ ] **Serialization** — implement `Boc` → wire format (JSON or binary) for subprocess calls
- [ ] **AST merge** — implement merging a returned `Boc` into the parent boc's AST (add fields, methods, type params)
- [ ] **Cycle detection** — detect circular `compile_time` triggers and report as compile errors
- [ ] **Caching** — cache compiled `Compile` executables keyed on source hash + input boc structure hash
- [ ] **Spec** — add `spec/12-compile-time-bocs.md` covering the `Compile` interface, two-phase build, execution timing, and constraint propagation

### 7. Remove `mix` as a Keyword

`mix` is no longer a language keyword; it becomes a library `Compile` implementation. See `docs/Features/Replaced features/mix.md`.

Depends on: item 6 (compile-time bocs) landing first — `mix` functionality moves to a `Mix` compile implementation.

- [ ] **Lexer** — remove `mix` from the keyword list (currently `token.MIX`)
- [ ] **Parser** — remove `MixStmt` parsing; `mix` becomes a regular identifier
- [ ] **Sema** — remove `mix` analysis (embedding resolution, conflict detection)
- [ ] **Lowering/Codegen** — remove Go-embedding code path for `mix`
- [ ] **Runtime** — implement `Mix` as a `Compile` boc in yzrt or stdlib (flattens fields and promotes methods from named boc into host)
- [ ] **Golden tests** — update or remove conformance tests that use `mix` syntax; add new test using `compile_time: [Mix]` infostring form
- [ ] **Spec 09** — remove `mix` from modules-and-organization spec; document `Mix` compile implementation

### 8. Concurrency Model — Behaviour-Oriented Concurrency (BOC)

**Scope: complete runtime redesign.** The current model (goroutines + `*Thunk[T]` + `Force()`) is replaced by cown-based resource acquisition semantics: every value is a protected concurrent resource; every invocation acquires all needed resources atomically before running; ordering is determined by resource overlap and spawn order. See `docs/Features/Concurrency.md`.

This item is **large and architectural** — it touches the runtime, codegen, and the thunk transparency mechanism. It should be broken into sub-phases when implementation begins.

#### Phase A — Mutex cowns (data-race freedom) — COMPLETE
- [x] **Runtime** — add `Cown` struct (embeds `sync.Mutex`) and `Schedule[T]` to `runtime/rt/cown.go`
- [x] **Codegen** — emit `std.Cown` as embedded field in every singleton boc struct
- [x] **Codegen** — change singleton method thunk from `std.Go(...)` to `std.Schedule(&self.Cown, ...)`; use split-BocGroup pattern (BocGroup.Wait() after Schedule) to avoid re-entrancy deadlock
- [x] **BocWithSig singletons** — `Call(params...)` methods use `std.Go` (not `std.Schedule`) to avoid deadlock on recursive singletons (e.g. `countdown` calling itself)
- [x] All 40 conformance goldens updated; `go test -race ./...` passes

#### Phase B.1 — Queue-based cown scheduler (spawn-order guarantee) — COMPLETE
- [x] **Runtime** — replaced mutex-based `Cown` with atomic lock-free queue scheduler matching BOC paper (Cheeseman et al., OOPSLA 2023) section 3 algorithm. Each `Cown` holds an atomic tail pointer to a linked list of `request` nodes; a `behaviour` runs when its count (one per required cown) reaches zero. `Schedule[T]` interface unchanged — no codegen changes. `releaseCown` uses CAS-then-spin to hand token to successor. Added `TestScheduleSerializes`, `TestSchedulePreservesOrder`, `TestScheduleTwoIndependentCowns` to `yzrt_test.go`.

#### Phase B.2 — Multi-cown atomic acquisition — COMPLETE
- [x] **Runtime** — add `ScheduleMulti[T](cowns []*Cown, fn func() T)`: registers behaviour on all cowns simultaneously; behaviour runs when all grant tokens; no sorting needed (queue handles per-cown ordering)
- [x] **Lowerer** — detect cown-typed arguments in boc calls (singleton boc instances passed as params); emit `ScheduleMulti` instead of `Schedule` when multiple cowns needed
- [x] **Parser** — accept lowercase IDENT as type expressions for singleton boc type annotations in params
- [x] **IR** — add `ExtraCowns []string` to `ThunkExpr`; add `IsThunk bool` to `DeclStmt` for safe hoisting
- [x] **Codegen** — emit `std.ScheduleMulti([]*std.Cown{...}, ...)` when ExtraCowns present; hoist IsThunk DeclStmts outside Schedule closure (split-BocGroup pattern extended)
- [x] **Conformance** — test 41 `41_multi_cown_sync.yz`: sync boc acquires bank + ledger cowns atomically

#### Phase C — Ownership-based field writes (SWMR fix)

`Thunk[T]` / `Go()` / `Force()` are compatible with the BOC model and are not removed. Phase C applies BOC ownership to field mutation: a write to a field owned by a different boc must go through that boc's cown.

- [x] **Lowerer/Codegen** — `lowerAssignment` detects field assignment targets on top-level singletons (via `LookupInFile`); wraps in `std.Schedule(&Target.Cown, ...).Force()`.
- [x] **Happens-before** — already guaranteed by Phase B.1 queue scheduler.
- [ ] **Sema** — optional: record singleton ownership per field access to avoid re-deriving at codegen time.
- [x] **Conformance** — golden test 42 `42_cross_cown_write.yz`: main writes directly to bank's field; generated code uses `std.Schedule`.
- [x] **Spec 08** — rewritten: replaced actor/channel model with BOC/cown model; §8.4 now describes cowns and behaviours; §8.6 SWMR updated; §8.8 implementation table corrected.

### 9. Remove `while` built-in — let recursion be the implementation

`while(cond, body)` is currently intercepted by `tryLowerWhile` in the lowerer and emitted as a Go `for` loop. This is a premature optimization: the language design says `while` is user-land recursion, not a primitive. The Go `for` shortcut also masks the true difficulty of `break`/`continue` (see item 2) and is the same category of mistake as `mix` being a keyword (see item 7) — a built-in standing in for something that should be expressed in Yz itself.

- [x] **Lowerer** — remove `tryLowerWhile` from `lower.go`; `while(cond, body)` becomes a regular boc call that goes through the normal recursive boc path
- [x] **Runtime** — remove `yzrt.While` function (removed from `yzrt/core.go` and `yzrt_test.go`)
- [x] **IR** — remove `ForStmt` node from `ir.go` and its codegen case from `codegen.go`; remove `lowerBocAsExpr` (only used by `tryLowerWhile`)
- [x] **Golden tests** — updated `05_while.yz` to define `while` as a top-level recursive `BocWithSig`; regenerated `05_while.go` — output is now a recursive Go function instead of a `for` loop
