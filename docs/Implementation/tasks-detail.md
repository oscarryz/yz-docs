#impl
Open ticket details. See tasks.md for the index.

---

## Bugs

- [ ] **[YZC-0008] Same-cown reentrant scheduling deadlock**

  Any code path that calls `Schedule(&self.Cown, ...)` while already executing inside a closure
  scheduled on `self.Cown` deadlocks — the outer task waits for its own completion.

  **Known manifestations:**

  1. **Local boc vars in main** (`37_local_boc_var` — confirmed deadlock with `TestRuntime`):
     Local boc variables (`foo #(String) = { ... }`) are lowered as methods on the enclosing
     singleton (`_mainBoc.Foo()`). When `Call()` — which holds `self.Cown` — calls
     `self.Foo().Force()`, `Foo()` schedules on the same `self.Cown` → deadlock.

  2. **HOF closures inside ScheduleMulti** (original case, still dormant):
     A closure passed as a callback argument and generated inside a `ScheduleMulti` body
     contains sync-body calls that assume the cown is held. If the closure escapes and is
     invoked outside the multi-cown body, those calls fire without holding the cown — data race.

  3. **Recursive local bocs** (was failing, now passing — see note):
     A local boc `f` calling itself via `self.F(n-1).Force()` inside `f()` would re-acquire
     `self.Cown` while held. This was the `39_local_boc_recursive` case; it currently passes,
     likely because the recursive call is handled inline rather than scheduled.

  **Root cause:** the lowerer emits all local boc vars as methods on the enclosing struct,
  sharing its cown. There is no mechanism to detect or prevent a task re-scheduling on a cown
  it already holds.

  **Fix direction:** Phase E.1 (implicit BocGroup per scope) removes statement-position `.Force()`
  calls, eliminating the blocking wait that causes the deadlock. Alternatively, local boc vars
  could be lowered to plain Go closures (not cown-scheduled methods) when they don't capture
  cown-bearing state — this would be a targeted fix without requiring the full Phase E rewrite.

---

## Language Features

- [ ] **[YZC-0009] Range iteration**

  `1.to(10).each({ i Int; ... })` — extend lowerer `.each` recognition to Range receiver. Depends on: YZC-0031.

- [ ] **[YZC-0013] Array append via `<<`**

  `a << item` → `a.Append(item)`; `Array.Append` exists in yzrt. Depends on: YZC-0031.

- [ ] **[YZC-0014] Option/Result method chaining**

  `result.or_else({ error Error; ... })`, `result.and_then({ val T; ... })`. Depends on: YZC-0031.

- [ ] **[YZC-0016] String concatenation with `++`**

  lowerer emits `Plusplus` but runtime `String` has no such method. Depends on: YZC-0031.

- [ ] **[YZC-0019] `break` / `continue` / `return` in loops**

  concurrency model settled; parser/sema/lowerer work is self-contained. Depends on: YZC-0031.

- [ ] **[YZC-0039] Operators audit**

  systematic comparison of spec vs. yzrt/lowerer: `%`, bitwise, string operators. Depends on: YZC-0031.

---

## Infrastructure

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

- [ ] **[YZC-0091] Nested singleton codegen: sub-singleton struct with own methods**

  `foo: { bar: { baz #() {} } }` — `bar` inside a singleton must lower to a
  sub-singleton struct with its own `Baz()` method, not a closure-returning
  `bar() Unit` method. Currently `foo.bar.baz()` fails:
  `Utils.extra.Help undefined (type func() rt.Unit has no field or method Help)`.

  Test: `examples/_wip/subdir_coexist` — promote when fixed.
  Depends on: YZC-0021. Will be superseded by YZC-0080 (uniform boc literal typing).

- [ ] **[YZC-0090] Multi-return for nested bocs (methods on singleton)**

  Multi-return (`wins, total : summary(3, 5)`) works for top-level singleton bocs
  but not for bocs that are methods on another singleton. `lowerMethod` only
  takes `Returns[0]`; `lowerBocBody` doesn't handle multi-return at all.

  Fix: detect `len(Returns) > 1` in `lowerMethod`, generate a result struct
  (same pattern as `lowerBodyOnlySingleton`), thread return count into
  `lowerBocBody` to collect and wrap the last N trailing expressions.

  Tests added here act as a regression guard when YZC-0080 supersedes this.

- [ ] **[YZC-0022] Multiple source roots**

  `src/` + `lib/` as independent FQN mount points. Depends on: YZC-0085.

- [ ] **[YZC-0023] Cancellation / non-local return**

  non-local `return` across goroutine boundaries; see `docs/Questions/How to cancel a running block.md`.

- [ ] **[YZC-0044] Producer-consumer example and golden test**

  `boring`/`while` producer-consumer in `docs/Features/Concurrency.md`. Depends on: YZC-0031.

- [ ] **[YZC-0058] Native type annotation — `macros: [Native]`**

  compiler-internal annotation for types backed by Go primitives. Depends on: YZC-0025, YZC-0059.

- [ ] **[YZC-0059] Design: macro interface interaction**

  concrete interaction patterns for `Macro` interface. Depends on: YZC-0025.

- [ ] **[YZC-0060] Design and implement `self` in Yz**

  `self` as compiler built-in or macro-generated binding. Depends on: YZC-0058, YZC-0059.

---

## Tooling

- [ ] **[YZC-0041] Dependency management**

  HTTPS-based import resolution; fetch and cache source. See `docs/Questions/Dependency Management.md`.

- [ ] **[YZC-0042] Package management (`yz` tool)**

  `yz init`, `yz add <url>`, lock file. Depends on: YZC-0041.

---

## Major Features

### YZC-0024 — `return`, `break`, `continue`

Blocked on concurrency model (YZC-0019, YZC-0023).

- [ ] Parser — `BreakStmt` / `ContinueStmt` AST nodes
- [ ] Sema — validate context
- [ ] Lowerer — emit compile error when encountered
- [ ] Spec 07 — update
- [ ] Golden tests — sema-level error tests

### YZC-0088 — Codegen: attach compiled annotation boc to declaration metadata

Deferred from YZC-0025. Once the macro system (YZC-0028) is defined, the compiler needs to store the parsed+type-checked annotation `*BocLiteral` alongside its target declaration so that macro passes can inspect and transform it.

- [ ] Define representation: how the annotation boc is stored on `ir.StructDecl` / `ir.SingletonDecl` / `ir.FuncDecl`
- [ ] Codegen — emit annotation metadata in generated Go (or as a side channel for the macro runner)
- [ ] Wire into macro invocation pipeline (YZC-0028)

### YZC-0028 — Macros (`Macro` interface)

Any boc with `Schema #()` and `run #(Boc, Boc)` satisfies `Macro`. Depends on: YZC-0025, YZC-0026, YZC-0027, YZC-0030, YZC-0066, YZC-0059.

- [ ] Sema — recognize `Macro` structural interface
- [ ] Sema — scan annotation for `macros: [...]`
- [ ] Boc metatype — `Boc` value type for `run`
- [ ] Two-phase build — compile `Compile` implementations first
- [ ] Serialization — `Boc` wire format
- [ ] AST merge — merge returned `Boc` into parent
- [ ] Cycle detection
- [ ] Caching — keyed on source hash
- [ ] Spec 12 — new spec file

### YZC-0031 — Scalar Types in Yz Source (uppering)

`Int/String/Bool/Decimal/Unit` move from Go to `stdlib/` with `compile-time:[Native]`. Depends on: YZC-0025, YZC-0028, YZC-0002.

- [ ] Define `macros: [Native]` annotation semantics
- [ ] Move scalar types to `stdlib/`
- [ ] Annotate native ops per method
- [ ] Implement higher-level methods in Yz
- [ ] Remove all primitive-type special-casing from the compiler
- [ ] `Bool.&&`/`||` — rewrite as lazy closure-taking boc methods

### YZC-0076 — Existential associated types: opaque-token / path-identity tracking

**Status note:** YZC-0075 was superseded by YZC-0079, which established that Yz uses structural typing rather than nominal path-identity for associated types. It is unclear whether this ticket is still needed — the opaque-token / cross-root rejection problem it describes may be moot in a fully structural system. Revisit after YZC-0079 has been used in real code; close if no concrete use case emerges.

Phase 2: the hard part. Deferred until YZC-0079 is settled and there is real usage demand.

- [ ] *design* — decide path-variable representation in the type system
- [ ] *design* — define scoping rules for opaque tokens (block-scoped vs field-storable)
- [ ] Sema — tag values with their existential path root at the point of production
- [ ] Sema — verify path roots match at call sites consuming opaque tokens
- [ ] Sema — reject cross-root usage with a clear error
- [ ] Conformance tests — opaque-token round-trip; cross-root rejection

### YZC-0080 — Uniform boc literal typing: one structural type derived from elements

#### Invariant

> Every boc literal, regardless of where it appears, receives one structural type derived mechanically from its elements. No code path branches on "is this a closure or a struct?" — that distinction is resolved at the use site by structural compatibility, not by classification during analysis.

#### Target design

Every boc literal gets one rich structural type:

```
BocLiteralType {
    Params    []BocParam      // TypedDecl nil-value entries → input signature
    Methods   []MethodField   // ShortDecl+BocLiteral or BocDecl-with-body entries
    Fields    []ValueField    // TypedDecl with value or ShortDecl with non-boc value
    Returns   []Type          // last-expression type(s)
}
```

#### Dependencies

Likely needs YZC-0025 (annotations / compile-time metadata). May also simplify YZC-0031.

- [ ] Design: define `BocLiteralType` in `sema/types.go`
- [ ] Sema: assign `BocLiteralType` to every `*ast.BocLiteral` in `analyzeExpr`; delete classification branches
- [ ] Sema: structural compatibility between `BocLiteralType` and `BocType` / `StructType` / interfaces
- [ ] Lowerer: dispatch on use-site expected type instead of sema classification flags
- [ ] Delete `hasInnerBocsOrMethods`, `bocLitHasParams`, `anonBocCache`, `anonDecls` from lowerer
- [ ] All existing tests pass

### YZC-0082 — Struct-outer nested type (concrete associated type)

`Foo: { Bar: {} }` — `Bar` is a type definition scoped to `Foo`; instances of `Foo` expose it as `f.Bar()`.

- [ ] *design* — decide whether inner type bodies can reference outer instance fields (path-dependent vs. self-contained)
- [ ] Sema: recognize uppercase struct-literal inside struct boc as concrete associated type definition
- [ ] Sema: `f.Bar()` resolves to the inner type; enforce no `Foo.Bar()` static access
- [ ] Lowerer: emit inner type as package-level Go struct; `f.Bar()` → constructor call
- [ ] Golden test: `Foo: { Bar: {} }` + `f.Bar()` compiles and runs

### YZC-0084 — Generic instantiation alias: `StringList : List(String)`

`StringList : List(String)` should declare a type alias for a concrete generic instantiation. Depends on: YZC-0027.

- [ ] *design* — decide emission: `type StringList = List[std.String]` (Go alias) or `type StringList struct { ... }` (copy)
- [ ] Sema: recognize `Name : GenericType(Args)` as instantiation alias
- [ ] Lowerer: emit appropriate Go type declaration
- [ ] `StringList(...)` constructor call works
- [ ] Golden test
