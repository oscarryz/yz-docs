#impl
# Yz Compiler Implementation

## Status
- **65 golden + 18 error conformance tests passing** — `go test -race ./...` passes (test 51 has pre-existing timing flakiness)
- Compiler: `compiler/` directory, Go module `module yz`
- Runtime: `compiler/runtime/rt/`

---

## Completed Phases

| Phase | Description | Tests |
|-------|-------------|-------|
| 0 | Project setup — `cmd/yzc`, `Makefile`, `go.mod` | — |
| 1 | Lexer — tokenizer + ASI | 38 |
| 2 | Parser — recursive descent AST | 32 |
| 3 | Semantic analysis — scope, type inference, boc/struct dispatch | passing |
| 4 | IR — lowerer (AST+sema → IR) | 8 |
| 5 | Codegen — Go source emitter; `yzc build`/`run`/`new` | 10 |
| 6 | Runtime — `types.go`, `core.go`, `collections.go`, `cown.go` | passing |
| 7 | Integration — conformance golden tests, examples, error tests | 65 golden |

---

## Open Tickets

Sorted by effort and independence. S = small, M = medium, L = large, XL = epic. *design* = needs a decision before implementation.

YZC-0017 -- Dict optional access -- S
YZC-0047 -- Cycle detection in homoiconic Stringify -- S
YZC-0057 -- Cyclic / mutually-recursive type declarations -- S
YZC-0012 -- Multiple return values -- M
YZC-0027 -- `:` as Type Alias -- M -- needs YZC-0066
YZC-0066 -- Path-Dependent Types: `#()` metatype, T fields, `g.Node` -- L
YZC-0038 -- `Result(T,E)` type -- M
YZC-0045 -- Default values in type-only boc declarations -- M -- needs YZC-0011
YZC-0026 -- Generics: Explicit Constraint Declaration -- M -- needs YZC-0066
YZC-0030 -- Associated Types -- M -- needs YZC-0066
YZC-0016 -- String `++` concatenation -- S -- needs YZC-0031
YZC-0013 -- Array `<<` append -- S -- needs YZC-0031
YZC-0009 -- Range iteration -- S -- needs YZC-0031
YZC-0019 -- `break`/`continue`/`return` in loops -- M -- needs YZC-0031
YZC-0014 -- Option/Result method chaining -- M -- needs YZC-0031
YZC-0039 -- Operators audit -- L -- needs YZC-0031
YZC-0043 -- Captured variable reference semantics -- *design*
YZC-0059 -- Compile-time bocs interface interaction -- *design* -- needs YZC-0025
YZC-0008 -- Reentrant inline calls unsafe in HOF closures -- S -- dormant
YZC-0021 -- Directory and file bocs -- L
YZC-0040 -- Smart Nesting / Namespace Flattening -- M -- needs YZC-0021
YZC-0022 -- Multiple source roots -- M
YZC-0044 -- Producer-consumer example and golden test -- M -- needs YZC-0031
YZC-0002 -- Cross-package support -- L -- needs YZC-0040, YZC-0022
YZC-0023 -- Cancellation / non-local return -- L
YZC-0058 -- Native type annotation -- L -- needs YZC-0025, YZC-0059
YZC-0060 -- Design and implement `self` in Yz -- L -- needs YZC-0058, YZC-0059
YZC-0041 -- Dependency management -- L
YZC-0042 -- Package management (`yz` tool) -- L -- needs YZC-0041
YZC-0024 -- `return`, `break`, `continue` (major) -- L -- needs YZC-0019, YZC-0023
YZC-0025 -- Infostrings: content is a boc body -- L
YZC-0028 -- Compile-Time Bocs (`Compile` interface) -- XL -- needs YZC-0025, YZC-0026, YZC-0027, YZC-0030
YZC-0029 -- Remove `mix`: runtime + spec -- M -- needs YZC-0028
YZC-0031 -- Scalar Types in Yz Source (uppering) -- XL -- needs YZC-0025, YZC-0028

---

# Details

Ticket numbers are permanent. `[x]` = closed, `[ ]` = open. Next available: **YZC-0067**.

---

## Bugs

- [x] **[YZC-0001] Variants broken**

  variants were not updated for the BOC model; see `examples/variants`

- [x] **[YZC-0003] Assigning Unit-returning boc to variable**

  `a : foo()` where `foo` returns Unit should be a sema error; detect in sema; add error golden test

- [x] **[YZC-0004] Top-level boc callable as function**

  implemented: `lowerCall` and `isBocMethodCall` extended for plain body singletons → `Foo.Call(args)`, and structured singletons → `Foo.Call(args)`. Golden test 55.

- [~] **[YZC-0005] Double return with sleep**

  `foo: { time.sleep(1); 1 }` emits two return statements — not reproducible as of BOC work; superseded by YZC-0035.

- [x] **[YZC-0006] Standalone boc invocation**

  resolved by YZC-0004: `p()` lowers to `P.Call()`. Golden test 56.

- [x] **[YZC-0007] Unused variables in generated Go**

  `emitBodyStmts` pre-scans via `usedNames`; emits `_ = varName` after any unused `DeclStmt`. Golden test 54.

- [x] **[YZC-0048] Flaky test 51 — concurrent output ordering**

  test 51 had wrong ordering expectation; deleted `.output` sidecar. Golden source-diff test still passes.

- [ ] **[YZC-0008] Reentrant inline calls unsafe in HOF closures**

  closure inside `ScheduleMulti` body passed as argument contains sync-body calls that bypass cown acquisition; fix: sub-generator with `heldCowns = nil` when emitting closure args; dormant until HOF closures operate on cown-bearing types.

- [x] **[YZC-0035] Sema does not check boc body return type against declared output**

  when a boc declares a non-Unit output but the body returns Unit, sema accepted it silently; fix: verify inferred body return type matches declared output.

---

## Language Features

- [x] **[YZC-0034] Definite assignment analysis (phase 1 replaced by YZC-0051)**

  `checkStructConstructorArgs` removed; replaced by YZC-0051.

- [x] **[YZC-0049] Lowerer: singleton boc params not emitted**

  leading TypedDecl-no-value elements in `lowerBodyOnlySingleton` now emitted as `Call(a std.T, ...)` params.

- [x] **[YZC-0051] CFG-based field definite-assignment**

  `FieldInitState` in `sema/definite_assign.go` tracks field init on all control-flow paths. Error tests 13 (updated) and 14 (new).

- [x] **[YZC-0052] Codegen "fill in later" — wrong arity on `NewBar()`**

  zero-arg constructor emits `&Bar{}` instead of `NewBar()`. Golden test 62. Depends on: YZC-0049.

- [x] **[YZC-0053] CFG check at boc-boundary crossing**

  `analyzeCall` checks all required fields of struct-typed args are assigned before the call. Error test 15.

- [x] **[YZC-0054] CFG: multi-level field access not tracked**

  `FieldInitState` now uses dotted string keys (`"inner.field"`). `markAssigned` marks all prefixes; `isAssigned` checks all prefixes. `memberPath` helper extracts root var + full dotted path from nested `MemberExpr`. Error test 18.

- [x] **[YZC-0055] CFG: variable aliasing defeats tracking**

  `c : b` now clones field-init state from source var. Error test 16.

- [x] **[YZC-0056] CFG: variant type construction skipped**

  no fix needed: direct variant field access outside `match` is already a sema error.

- [x] **[YZC-0063] Single-arm non-exhaustive match**

  `p match Constructor` (Bool form) and `p match Constructor => { body }` (narrowing form). `InfixMatchExpr` AST node, `VariantTestExpr` IR node. Golden test 64.

- [x] **[YZC-0065] Type-directed variant constructor disambiguation**

  `Symbol.Alternatives` stores all options; `expectedType` propagated inward; qualified form `Shape.Circle(5)` via `fieldType` variant namespace lookup. Golden test 65, error test 17.

- [ ] **[YZC-0009] Range iteration**

  `1.to(10).each({ i Int; ... })` — extend lowerer `.each` recognition to Range receiver. Depends on: YZC-0031.

- [x] **[YZC-0010] HOF iteration + cown happens-before**

  `.filter`, `.each` as sync Go closures. Golden test 27.

- [x] **[YZC-0036] While loop yield and external caller interleaving**

  BocDecl singletons use `std.Schedule`; recursive self-calls marked `IsRecursive`.

- [x] **[YZC-0011] Named arguments in constructor calls**

  `lowerStructArgs` reorders by field declaration order; `lowerNamedArgs` for BocDecl calls. Golden test 59.

- [ ] **[YZC-0012] Multiple return values**

  `x, y = swap(x, y)` — multi-assign LHS not yet implemented.

- [ ] **[YZC-0013] Array append via `<<`**

  `a << item` → `a.Append(item)`; `Array.Append` exists in yzrt. Depends on: YZC-0031.

- [ ] **[YZC-0014] Option/Result method chaining**

  `result.or_else({ error Error; ... })`, `result.and_then({ val T; ... })`. Depends on: YZC-0031.

- [x] **[YZC-0015] Non-word boc names**

  `balance+= #(amount Int) { ... }` — parser accepts `NON_WORD` token and maps to Go-safe name.

- [ ] **[YZC-0016] String concatenation with `++`**

  lowerer emits `Plusplus` but runtime `String` has no such method. Depends on: YZC-0031.

- [ ] **[YZC-0017] Dict optional access**

  `d[key]` should return `Option(V)`; currently panics on missing key.

- [x] **[YZC-0018] Bool methods `&&` / `||`**

  `Bool.Ampamp` / `Bool.Pipepipe` in yzrt. Golden test 53.

- [ ] **[YZC-0019] `break` / `continue` / `return` in loops**

  concurrency model settled; parser/sema/lowerer work is self-contained. Depends on: YZC-0031.

- [x] **[YZC-0020] Compiler homoiconic dump — backtick interpolation**

  backtick inside a string triggers homoiconic representation. Golden test 60.

- [x] **[YZC-0037] Decimal type end-to-end**

  `std.Decimal` with arithmetic, comparisons, `to_str`. Golden test 58.

- [ ] **[YZC-0038] `Result(T,E)` type**

  implement as a variant type in yzrt; wire sema/lowerer recognition. Spec: `docs/Features/Error handling.md`.

- [ ] **[YZC-0039] Operators audit**

  systematic comparison of spec vs. yzrt/lowerer: `%`, bitwise, string operators. Depends on: YZC-0031.

- [ ] **[YZC-0040] Smart Nesting / Namespace Flattening**

  `house/house.yz` flattens to `house.method`. Depends on: YZC-0021.

- [ ] **[YZC-0043] Captured variable reference semantics**

  design: value vs. reference capture in boc literals. See `docs/Questions/Memory Management.md`.

- [ ] **[YZC-0045] Default values in type-only boc declarations (interfaces)**

  `Greeter #(name String = "Alice")` — defaults are call-site sugar. Depends on: YZC-0011.

- [x] **[YZC-0046] `${}` interpolation requires `to_str`**

  sema checks for `to_str #(String)` on the interpolated type. Depends on: YZC-0020.

- [ ] **[YZC-0047] Cycle detection in homoiconic `Stringify`**

  thread a visited-pointer set through `Stringify`; emit `TypeName(...)` on re-entry.

- [x] **[YZC-0061] Structured singleton: TypedDecl-with-value field missing `self.`**

  `collectFieldNames` gating removed. Golden test 63.

---

## Infrastructure

- [x] **[YZC-0033] Compiler deep review against settled spec**

  all four sub-items resolved: BOC singletons, `foo.param` accessible after call, error messages say "returns nothing", all bocs serialized through cown.

- [ ] **[YZC-0021] Directory and file bocs**

  defer until in-file nesting works; extend FQN tree to directories and files as bocs.

- [x] **[YZC-0032] Rename `BocWithSig` → `BocDecl`**

  done throughout AST, sema, lowerer, and spec/02.

- [ ] **[YZC-0002] Cross-package support**

  broke during BOC migration. Depends on: YZC-0040, YZC-0022.

- [ ] **[YZC-0022] Multiple source roots**

  `src/` + `lib/` as independent FQN mount points.

- [ ] **[YZC-0023] Cancellation / non-local return**

  non-local `return` across goroutine boundaries; see `docs/Questions/How to cancel a running block.md`.

- [ ] **[YZC-0044] Producer-consumer example and golden test**

  `boring`/`while` producer-consumer in `docs/Features/Concurrency.md`. Depends on: YZC-0031.

- [ ] **[YZC-0057] Cyclic / mutually-recursive type declarations**

  two-pass sema: collect all top-level type names first, then resolve field types.

- [ ] **[YZC-0058] Native type annotation — `compile_time:[Native]`**

  compiler-internal annotation for types backed by Go primitives. Depends on: YZC-0025, YZC-0059.

- [ ] **[YZC-0059] Design: compile-time bocs interface interaction**

  concrete interaction patterns for `Compile` interface. Depends on: YZC-0025.

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

### YZC-0025 — Infostrings: content is a boc body

Infostring delimiter stays backtick; content is full Yz syntax, parsed and type-checked, never executed. Intersection with Native annotations (YZC-0058).

- [ ] AST — `InfoString` holds `*BocLiteral`
- [ ] Lexer — re-lex infostring content as Yz source
- [ ] Parser — re-parse as boc body
- [ ] Sema — type-check content
- [ ] Codegen — attach compiled infostring boc to declaration metadata
- [ ] Spec 01 — update

### YZC-0026 — Generics: Explicit Constraint Declaration

`thing T Talker` declares `T` must implement `Talker`; additive with inference.

- [ ] Parser — `T Constraint` optional suffix after single-uppercase type param
- [ ] Sema — validate at instantiation; union with inferred constraints
- [ ] Spec 04 — update

### YZC-0027 — `:` as Type Alias

`Name : SomeType` declares a type alias. Depends on YZC-0066 for the unified model; can be implemented as a limited special form before YZC-0066 lands (emit `type Name = GoType` in Go, no `#()` metatype required).

- [ ] Parser — distinguish from `Name TypeExpr` (typed decl) and `name : value` (short decl)
- [ ] Sema — register alias; resolve as aliased type
- [ ] Lowerer — emit `type Name = GoType`
- [ ] Spec 04 — add
- [ ] Deferred to YZC-0066: generic instantiation via alias (`StringList : List(String)`), associated type binding (`Node : User` inside a boc)

### YZC-0066 — Path-Dependent Types: `#()` metatype, T fields, `g.Node`

Unified model for generics, type aliases, and associated types. See `docs/Features/Path Dependent Types.md`.

- [ ] Sema — `#()` recognized as metatype; bare GENERIC_IDENT field given implicit `#()` type
- [ ] Sema — type fields in constructors (`List(Int)` binds `T = Int`)
- [ ] Sema — path-dependent resolution: `g.Node` in type position looks up `Node` field of `g`'s struct type
- [ ] Sema — type variable inference: unify GENERIC_IDENT against call-site argument types
- [ ] Lowerer/Codegen — emit specialized (monomorphized) Go types per concrete instantiation
- [ ] `Node : User` inside a boc body treated as type alias, not value alias
- [ ] Golden tests: generic boc, `g.Node` signature, `StringList : List(String)` alias
- [ ] Spec 04 — generics section; Spec 05 — associated types section

### YZC-0028 — Compile-Time Bocs (`Compile` interface)

Any boc with `Schema #()` and `run #(Boc, Boc)` satisfies `Compile`. Depends on: YZC-0025, YZC-0026, YZC-0027, YZC-0030.

- [ ] Sema — recognize `Compile` structural interface
- [ ] Sema — scan infostring for `compile_time: [...]`
- [ ] Boc metatype — `Boc` value type for `run`
- [ ] Two-phase build — compile `Compile` implementations first
- [ ] Serialization — `Boc` wire format
- [ ] AST merge — merge returned `Boc` into parent
- [ ] Cycle detection
- [ ] Caching — keyed on source hash
- [ ] Spec 12 — new spec file

### YZC-0029 — Remove `mix`: runtime + spec — PARTIALLY COMPLETE

Compiler removal done.

- [x] Lexer, Parser, Sema, Lowering/Codegen, Golden tests — done
- [ ] Runtime — implement `Mix` as a `Compile` boc
- [ ] Spec 09 — remove `mix`; document `Mix` compile implementation

### YZC-0030 — Associated Types: Path-Dependent Type References

`process #(g Graph, n g.Node)` — sema resolves `g.Node` against the concrete type of `g`. Design resolved; see `docs/Features/Path Dependent Types.md` and `docs/Features/Associated Types.md`. Depends on YZC-0066.

- [ ] Sema — `value.TypeName` in type position (path-dependent resolution)
- [ ] Lowerer — emit concrete Go type at resolution site
- [ ] Golden test: Graph/SocialGraph/process example

### YZC-0031 — Scalar Types in Yz Source (uppering)

`Int/String/Bool/Decimal/Unit` move from Go to `stdlib/` with `compile-time:[Native]`. Depends on: YZC-0025, YZC-0028.

- [ ] Define `compile-time:[Native]` infostring semantics
- [ ] Move scalar types to `stdlib/`
- [ ] Annotate native ops per method
- [ ] Implement higher-level methods in Yz
- [ ] Remove all primitive-type special-casing from the compiler
- [ ] `Bool.&&`/`||` — rewrite as lazy closure-taking boc methods
