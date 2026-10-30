---
type: impl
generated: { by: "oscarryz", at: 2026-06-05T23:56:26+02:00 }
---
#impl
Legacy detail entries, none of them in the open index. Open tickets live in
[tickets/](tickets/index.md).

---

## Thunk / Concurrency

- [ ] **[YZC-0094] Fully lazy thunk model: propagate thunks through all expressions, force only at BocGroup boundary** *(design)*

  ### History

  - **Phase E.2** (commit `64ea0ea`): tried putting laziness *inside* scalar
    types — `Int`, `String`, etc. got optional lazy fields; boc calls returned
    the scalar type directly. `GoWait` replaced `Go+Force`.
  - **Phase E.3** (commit `fc0989a`): reverted. Scalars are plain again; all
    boc methods return `*Thunk[T]`. Added `GoStore[T]` for typed-decl boc
    calls. `thunkVars` / auto-forcing re-introduced.

  Neither phase reached the correct model because bigger tickets were in the
  way. The backlog is now clear enough to revisit.

  ### Target model

  Every boc call returns a `*Thunk[T]`. **Thunks are never forced eagerly.**
  Instead they propagate through all expressions and are forced only when the
  boc body's BocGroup resolves them at the end:

  ```yz
  g : greet()           // Thunk[String]  — goroutine launched
  w : world.say()       // Thunk[String]  — goroutine launched, both concurrent

  g == "hi from greet"    // Thunk[Bool]  — no goroutine, cold chain
  ? { print("ok") }, {}  // Thunk[Unit]  — still lazy

  w == "hello from world"
  ? { print("ok") }, {}  // Thunk[Unit]
  ```

  Main's BocGroup forces every statement-level `Thunk[Unit]` at the bottom.
  Forcing cascades inward (pull-based): `Thunk[Unit]` → forces `Thunk[Bool]`
  → forces `Thunk[String]` → goroutine completes.

  BocGroup is **necessary**: Go kills un-waited goroutines when `main()`
  returns, so structured concurrency must always wait at the end of each boc.

  ### Why `==`, `?`, etc. produce thunks

  In Yz, `==` and `?` are ordinary method calls, not operators with special
  compiler rules. Consistency requires they behave the same regardless of
  whether their arguments are concrete values or thunks. Therefore:

  - `Thunk[String].Eqeq(String)` → `Thunk[Bool]`
  - `Thunk[Bool].Qm(Thunk[Unit], Thunk[Unit])` → `Thunk[Unit]`

  This is implemented by making `Thunk[T]` forward all of `T`'s methods,
  returning `Thunk[ReturnType]` rather than `ReturnType`.

  Variables and inline calls are identical — `g == "hi"` and `greet() == "hi"`
  both produce `Thunk[Bool]` with no forcing difference.

  ### Waiting IS forcing — BocGroup simplification

  Today `GoWait` spawns a *meta-goroutine* that calls `.Force()` on a thunk
  and signals `sync.WaitGroup`. This is a two-layer design: boc goroutine +
  meta-goroutine just to do the force.

  In the fully lazy model this collapses. The boc goroutines launched by
  `Schedule` are already running concurrently. Forcing their thunks
  *sequentially* loses nothing — by the time we force thunk 2, goroutine 2 is
  already running (or done). Total wait time is still the slowest goroutine,
  same as a WaitGroup.

  Therefore `BocGroup` becomes:

  ```go
  type BocGroup struct{ thunks []Forceable }

  func (g *BocGroup) Add(th Forceable)  { g.thunks = append(g.thunks, th) }
  func (g *BocGroup) Wait() {
      for _, th := range g.thunks { th.Force() }
  }
  ```

  - No `sync.WaitGroup`
  - No `GoWait` (meta-goroutine eliminated)
  - No `GoStore` (`g : greet()` holds `*Thunk[String]`, not a concrete value)
  - `Wait()` IS the forcing mechanism — waiting and forcing are the same act

  ### Note on "truly lazy"

  `Thunk[T]` already supports both forms (see `thunk.go`):

  - `Go[T](fn)` — hot thunk: goroutine launched immediately, result cached
    when the goroutine completes.
  - `NewThunk[T](fn)` — cold thunk: no goroutine, `fn` runs in the caller's
    goroutine the first time `Force()` is called, result cached via `sync.Once`.

  Forwarding methods on `Thunk[T]` are therefore trivial cold thunks:

  ```go
  func (th *Thunk[String]) Eqeq(other String) *Thunk[Bool] {
      return NewThunk(func() Bool { return th.Force().Eqeq(other) })
  }
  ```

  The goroutines are genuinely concurrent — `Force()` on a hot thunk blocks
  until the goroutine completes. The current overhead is the meta-goroutines
  spawned by `GoStore` to force hot thunks — those are eliminated here.

  ### Implementation plan

  **Phase 1 — Runtime (`runtime/rt/`)** ✓ DONE

  - [x] `thunk.go`: no changes needed; `NewThunk` and `Go` are already correct.
  - [x] `core.go`:
    - [x] Replace `BocGroup.sync.WaitGroup` with `[]func()` (slice of force closures).
    - [x] `BocGroup.Add(fn func())` — append closure to pending list.
    - [x] `BocGroup.Wait()` — iterate and call each closure. No WaitGroup.
    - [x] `GoWait`, `GoStore`, `GoStoreAny` kept temporarily (marked with
      TODO comments); **must be deleted as the final step of Phase 3** once
      codegen no longer emits them.
  - [x] `thunk_scalars.go` (new file — per scalar: `String`, `Int`, `Bool`, `Decimal`, `Unit`):
    - [x] Go generics do not allow type-specific methods on `*Thunk[T]`, so
      each scalar gets a concrete wrapper: `ThunkString`, `ThunkInt`,
      `ThunkDecimal`, `ThunkBool`, `ThunkUnit`, `ThunkRange`.
    - [x] Each wrapper exposes forwarding methods that return new cold thunks.
    - [x] Example: `ThunkString.Eqeq(String) ThunkBool`, `ThunkBool.Qm(...)`,
      `ThunkInt.Plus(Int) ThunkInt`.
    - [x] Constructors: `GoStringThunk(fn)` (hot), `NewStringThunk(fn)` (cold).

  **Phase 2+3 — IR + Codegen** ✓ DONE

  - [x] `ir.go`: SpawnExpr simplified — `Body []Stmt` removed, `Thunk Expr`
    added. `StoreAnyType` retained for one path-dependent test case.
  - [x] `lower.go`: all 12 SpawnExpr construction sites updated to use
    `Thunk:` field. Variable type stays concrete (`var n T`) for now —
    full type change to `*Thunk[T]` deferred (no existing tests trigger it).
    Note: `thunkVars` and `lowerExprForced` removal deferred because
    concrete-variable use sites (argument passing, string interpolation) still
    need forcing; removing those requires the ThunkX codegen path (future work).
  - [x] `codegen.go`:
    - `spawnForceInner` removed.
    - `emitSpawnStmt` added: always hoists thunk to `_thN` var (goroutine
      starts at registration time, not deferred to Wait), then emits Add closure.
    - `emitImmediateBody` updated to use `Thunk` field directly, preserving
      `StoreVar` in closures.
    - `thunkCount *int` added to generator (shared across sub-generators).
    - `collectUsedExpr` updated for new SpawnExpr shape.
  - [x] `GoStore`/`GoWait`/`GoStoreAny` removed from runtime (no longer emitted).

  **Phase 4 — Golden files** ✓ DONE

  - [x] 43 golden `.go` files regenerated. All 92 golden + error + runtime
    tests pass.
  - [ ] Add new golden tests for: `greet() == "hi" ? { print("ok") }, {}`,
    inline arithmetic on boc results, string interpolation of thunks.
    (Blocked on ThunkX codegen integration — separate follow-up.)

  **Out of scope for this ticket**

  - Full propagation through ALL method calls (e.g. `Thunk[String].Length`
    on user-defined types) — covered when scalar types move to Yz source
    (YZC-0031), which can generate the forwardings from the method definitions.
  - Changes to sema type inference for thunk-returning expressions.

- [ ] **[YZC-0095] Phase 7: dethunkification — scalar-intrinsic lazy types, eliminate ThunkX wrappers** *(M)*

  ### Context

  YZC-0094 Phases 1–6 introduced `ThunkX` wrapper types (`ThunkString`, `ThunkBool`, etc.)
  and `WrapXThunk` constructors so boc call results could chain through method calls lazily.
  Phase 6 added 33 T-variant methods (`EqeqT`, `PlusT`, …) for both-sides-lazy binary ops.

  Commit `64ea0ea` (Phase E.2, pre-YZC-0094) tried an alternative — embedding lazy state
  directly in scalar types — but was reverted due to an unrelated design confusion, not
  because the approach was wrong.

  ### Goal

  Public boc methods return the scalar type directly (`std.String`, `std.Int`, etc.) with
  lazy state carried internally. The `*Thunk[T]` is hidden inside the scalar; callers never
  see or wrap it. `ThunkX` wrapper types, `WrapXThunk` constructors, and T-variant methods
  are eliminated from the runtime and the lowerer.

  ### Target generated code

  ```go
  // was: func (self *_greeterBoc) Greet() *std.Thunk[std.String]
  func (self *_greeterBoc) Greet() std.String {
      return std.LazyString(std.Schedule(&self.Cown, func() std.String {
          return self.greet()
      }))
  }

  // binary op: Greet().Eqeq(NewString("hello")) — no wrapping, no T-variant
  // result is std.Bool (lazy), passed to Bool.Qm
  ```

  ### What is eliminated

  - `ThunkString`, `ThunkInt`, `ThunkBool`, `ThunkDecimal`, `ThunkUnit`, `ThunkRange`
  - `WrapStringThunk`, `WrapIntThunk`, `WrapBoolThunk`, `WrapDecimalThunk`, `WrapUnitThunk`, `WrapRangeThunk`
  - All 33 T-variant methods (`EqeqT`, `PlusT`, `AmpampT`, …)
  - `isBocMethodCall`, `isThunkXExpr`, `thunkWrapFuncFor`, `lowerThunkConditional`,
    `bocBranchAsClosure` from `lower.go`
  - Special-case WrapXThunk emission in `lowerExpr/BinaryExpr`

  ### What stays / changes

  - **`LazyX` constructors** added to each scalar type in `types.go` (as in `64ea0ea`).
  - **`Bool.Qm(trueCase, falseCase func() any) any`** — the `?` operator as a method on
    `Bool` directly (replaces `ThunkBool.Qm`). Forces `self`, picks branch, returns result.
    Return type is `any`; the lowerer knows the concrete type statically and casts if needed.
  - **BocGroup interaction**: unchanged — `_bg0.Add(func() { result.Force() })` where
    `result` is the `*Thunk[Unit]` backing the lazy Unit returned by `Qm`. Or the BocGroup
    switches to a `Waitable` interface and calls `Await()` on the scalar.
  - **Conditional lowering**: `isThunkXConditional` check disappears — `Bool` is always
    `Bool` (lazy or not), so the lowerer uses `if/else` for statement position and
    `Bool.Qm(closures)` for expression position, unconditionally.
  - **User-defined types**: not affected. Struct boc types (`*_personBoc`) still return
    `*Thunk[Person]` from public methods; forcing is explicit. This ticket covers the 5
    built-in scalar types only.

  ### Implementation steps

  1. `runtime/rt/types.go`: add `*lazyX` structs and `LazyX(th *Thunk[X]) X` constructors
     for String, Int, Bool, Decimal, Unit. Add `Bool.Qm`.
  2. `runtime/rt/thunk_scalars.go`: delete entire file (ThunkX types).
  3. `internal/codegen/codegen.go`: `emitMethodDecl` wraps scalar return with `LazyX`.
  4. `internal/ir/lower.go`: remove ThunkX detection paths; `lowerExpr/BinaryExpr` no
     longer wraps or T-variants; `lowerConditional` unconditionally uses `if/else` or Qm.
  5. Regenerate all golden files.

  ### Scope note

  This covers only the 5 built-in scalar types. Full propagation through user-defined type
  methods remains deferred to YZC-0031 (scalar types in Yz source).

---

## Infrastructure

- [x] **[YZC-0093] Uppercase root file (`Foo.yz`) always-wrap: example + spec §9 clarification**

  YZC-0092 implemented always-wrap for lowercase root files only. Uppercase root
  files (e.g. `Foo.yz`) follow the same invariant but have no conformance test.

  Two sub-cases:
  - `Foo.yz` with free-floating fields (`name String; age Int`) → wraps to
    `Foo: { name String; age Int }` — the wrapper IS the struct type.
  - `Foo.yz` with an inner `Foo: {}` → wraps to `Foo: { Foo: {} }` — inner
    becomes an associated type (struct-outer, YZC-0082); `fileWrapperHasInnerBoc`
    should trigger unwrap so the inner `Foo` becomes top-level.

  Work:
  - Add a small root-level uppercase example (e.g. `examples/root_type/`) with
    a `Foo.yz` struct file and a `main.yz` that constructs it
  - Add spec §9 Invariant 1 clarifying note covering both sub-cases
  - Add or promote conformance test

- [x] **[YZC-0092] Always-wrap root files; `main()` as explicit entry invocation**

  Remove the `hasTopLevelBocNamed` guard in `build.go` so all root files are
  always wrapped in a boc named after the file — consistent with spec §9
  Invariant 1 ("file content = boc body named after the file").

  Consequences:
  - `main.yz` with `main: {}` wraps to `main: { main: {} }`. The outer `main`
    executes its body; to run the inner boc the file must call `main()` explicitly.
  - `main.yz` with free-floating statements (no inner `main: {}`) works as-is —
    the outer `main` body just executes them directly.
  - `Foo.yz` with `Foo: {}` wraps to `Foo: { Foo: {} }` — `Foo` inner becomes
    an associated type (YZC-0082). `Foo.yz` with `name String; age Int` wraps
    to `Foo: { name String; age Int }` — struct type, constructor works.

  Work:
  - Remove `hasTopLevelBocNamed` guard and helper from `build.go`
  - Update all existing examples: drop explicit same-name wrapper, or add
    `main()` call at end of `main.yz`
  - Update conformance tests accordingly
  - Add clarifying example to spec §9 Invariant 1 covering `Foo.yz` and
    `main.yz` with inner `main: {}`

- [ ] **[YZC-0022] Multiple source roots**

  Enable `yzc build` to accept multiple source root directories. Each root
  contributes `.yz` files to the same FQN namespace; the first argument owns
  `target/gen/` and `target/bin/`. Source roots themselves are never written to.

  **Motivation:** required for stdlib and third-party libraries to live outside
  the user project without copying source in. Foundation for YZC-0041 (dependency
  management), which will derive the root list from `project.info` instead of CLI
  args, and YZC-0031 (scalar types in Yz source), which needs a stdlib root.

  **CLI:**
  ```
  yzc build myproject/ stdlib/ somelib/
  ```
  No implicit default — require at least one argument. `yzc build` with no args
  is an error with a usage hint. (Smart defaults belong in YZC-0041.)

  **Semantics (spec §9.2 Invariant 3):**
  - FQN is computed relative to each root: `stdlib/net/http.yz` → `net.http`
  - Same FQN path + different declaration names across roots → merged into one boc
  - Same FQN path + same declaration name across roots → compilation error
  - File/dir coexistence (Invariant 5) applies per root independently

  **Implementation delta:**
  1. `cmd/yzc/main.go` — parse positional args; first arg = project dir, rest =
     extra source roots; error if none given
  2. `compileProject(dir string)` → `compileProject(projectDir string, srcRoots []string)`
  3. `walkYzFiles` called per root; `fileEntry` gets `srcRoot` field so FQN is
     computed relative to the correct root
  4. `byDir` grouping extended to group by `(relDir, name)` across all roots;
     detect and report same-FQN same-name collisions
  5. `target/` path hardcoded to `projectDir` (first arg) throughout

  **Out of scope:** stdlib auto-injection, `project.info` resolution, dependency
  fetching — those are YZC-0041/0042.

  Depends on: YZC-0085 (done).

- [x] **[YZC-0059] Design: macro interface interaction** — [resolved](../Questions/solved/Macro%20Interface%20Interaction%20Design.md)

  Settle the taxonomy and interaction model for macros before implementing
  YZC-0028 or any macro-dependent feature (YZC-0041, YZC-0058, YZC-0060).
  Depends on: ~~YZC-0025~~.

  **Questions to resolve:**

  **1. Macro taxonomy — two distinct concepts, not one**

  - *Code-generating macros*: receive an annotated Boc, return a transformed
    Boc merged back into the AST. Closest to Rust proc-macros. Declare
    `schema #()` so the compiler can validate the annotation body at
    compile time. Run during compilation.
  - *Build-hook programs*: `yz fetch`, `Native`, Go generator wrappers. Act
    on the environment (network, filesystem), not on code. No schema, no AST
    output. Run before or alongside compilation. `yz fetch` is the canonical
    example — a standalone program that reads annotations and populates a cache.

  The key differentiator: **schema**. Macros declare `schema #()`; the
  compiler validates the annotation body against it at compile time. Programs
  parse annotations ad hoc with no compile-time guarantees.

  **2. Macro identification — three candidate models**

  - *Explicit list only*: `!:[JSON]` / `macros: [JSON]` is the authority.
    Unambiguous; slightly verbose when the field name already implies the macro.
  - *Attribute name implies macro*: `json: {...}` in an annotation body
    implicitly triggers the `JSON` macro. Concise; risks silent failure on typo
    or when macro is not in scope.
  - *Hybrid (recommended starting point)*: explicit list (`!:[...]`) is the
    authority and always required. Macros find their config by convention —
    a macro named `JSON` looks for a `json:` field in the annotation body.
    No implicit triggering; no redundant naming. Explicit list = what runs;
    field name = where config lives.

  **3. Annotation format for tooling (gates YZC-0041)**

  `yz fetch` needs to locate dependency declarations in annotations. The
  agreed format (to be finalised here):

  ```
  `!:[Deps]
  dependencies: [
      my_lib: { version: "1.0.0", url: "https://..." }
  ]
  `
  ```

  - `!:[Deps]` (or `macros: [Deps]`) = explicit trigger
  - `dependencies:` = config field, routed to `Deps` by convention
  - Deps can appear inline in source or in a `name.info` companion file
  - Multiple declarations across files are merged; one `yz.lock` per project root
  - Lock file pins SHA (git commit SHA preferred — immutable by definition)
  - `Deps` macro = thin schema-only macro for compile-time validation;
    actual fetching is done by `yz fetch` (a standalone program, not a macro)

  **4. Programs vs macros — clear boundary**

  Programs (`yz fetch`) are not macros. They scan annotations for markers,
  act on the environment, and pass results to the compiler (e.g. extra source
  roots via `yzc build proj/ ~/.yz/cache/lib/`). They do not satisfy the
  `Macro` interface and do not transform AST. The `Deps` macro and `yz fetch`
  are complementary: the macro validates, the tool acts.

---

## Tooling

- ~~[ ] **[YZC-0041] `Deps` macro — compile-time dependency validation**~~

  ~~Cancelled.~~ Dependencies are passive annotation metadata, not macros.
  Superseded by YZC-0097.

- [x] **[YZC-0097] Annotation metadata contract for project and dependency configuration** *(replaces YZC-0041)*

  Define the annotation format for project-level and dependency metadata so
  that external tools (`yz fetch`, build tool) have a stable, documented
  contract to read from. This is passive metadata — lowercase keys, no macro
  dispatch, no AST transformation. The compiler validates the annotation shape
  at compile time (typed metadata) but never fetches or resolves anything.
  Compilation remains predictable: no external processes are triggered by the
  compiler as a side effect of annotation processing.

  - Define `dependencies:` annotation format (fields: `version`, `url`, `sha`)
  - Define `project:` annotation format (name, source paths, description)
  - Specify where these annotations may appear (project root `.yz` file or
    `project.info` companion)
  - Document the contract between annotation metadata and external tools

  ```yz
  `dependencies: [
      my_lib: { version: "1.0.0", url: "https://..." }
  ]
  project: { name: "my_project", source_paths: ["src/"] }`
  ```

  Depends on: ~~YZC-0059~~.

---

## Major Features

### ~~YZC-0076 — Existential associated types: opaque-token / path-identity tracking~~ — CLOSED

**Closed.** The original motivation was to support an array of heterogeneous macros, each carrying a different `Schema` shape, validated at compile time:

```yz
macros: [Debug, Serialize, Ord]   // each has a different Schema
```

Iterating over that list and validating each schema would require knowing the specific `Schema` type per element — which requires either path-dependent types (`m.Schema` where `m` is a variable) or first-class existential types.

The dispatch mechanism changed (YZC-0059): macros are now triggered by uppercase type name resolution in the annotation body. The compiler always dispatches to a specific concrete type (`Debug`, `Serialize`, etc.) — it never holds a generic `Macro` variable and dispatches through it. Therefore `m.Schema` on a runtime variable never arises, and the cross-root rejection problem does not occur.

**What Yz has:**
- Associated types (YZC-0066) — compile-time, concrete shape known when the specific type is in hand
- Structural compatibility — checks annotation body against the concrete `MacroImpl.Schema`

**What Yz does not have:**
- Runtime path-dependent types — `m.Schema` where `m` is a variable of interface type `Macro`
- First-class existential types — holding a heterogeneous collection and unpacking each element's specific associated type shape

These gaps would matter if macros were ever runtime-dispatched objects or if a user wanted a heterogeneous list of macros with per-element schema introspection. Neither is needed under the current dispatch model.


### YZC-0082 — Struct-outer nested type (concrete associated type)

`Foo: { Bar: {} }` — `Bar` is a type definition scoped to `Foo`; instances of `Foo` expose it as `f.Bar()`.

- [ ] *design* — decide whether inner type bodies can reference outer instance fields (path-dependent vs. self-contained)
- [ ] Sema: recognize uppercase struct-literal inside struct boc as concrete associated type definition
- [ ] Sema: `f.Bar()` resolves to the inner type; enforce no `Foo.Bar()` static access
- [ ] Lowerer: emit inner type as package-level Go struct; `f.Bar()` → constructor call
- [ ] Golden test: `Foo: { Bar: {} }` + `f.Bar()` compiles and runs

### YZC-0084 — Generic instantiation alias: `StringList : List(String)`

`StringList : List(String)` should declare a type alias for a concrete generic instantiation. Depends on: ~~YZC-0027~~.

- [ ] *design* — decide emission: `type StringList = List[std.String]` (Go alias) or `type StringList struct { ... }` (copy)
- [ ] Sema: recognize `Name : GenericType(Args)` as instantiation alias
- [ ] Lowerer: emit appropriate Go type declaration
- [ ] `StringList(...)` constructor call works
- [ ] Golden test
