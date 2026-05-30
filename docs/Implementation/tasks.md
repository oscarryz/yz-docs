#impl
# Yz Compiler Implementation

## Status
- **78 golden + 22 error conformance tests passing** — `go test -race ./...` passes (test 51 has pre-existing timing flakiness)
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
YZC-0012 -- Multiple return values -- M  
YZC-0027 -- `:` as Type Alias -- M  
YZC-0038 -- `Result(T,E)` type -- M  
YZC-0045 -- Default values in type-only boc declarations -- M -- needs YZC-0011  
YZC-0071 -- Implicit constraint synthesis for type params used in method params -- M  
YZC-0070 -- Anonymous boc literal as structural interface value -- M  
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
YZC-0028 -- Compile-Time Bocs (`Compile` interface) -- XL -- needs YZC-0025, YZC-0026, YZC-0027, YZC-0030, YZC-0059   
YZC-0029 -- Remove `mix`: runtime + spec -- M -- needs YZC-0028  
YZC-0031 -- Scalar Types in Yz Source (uppering) -- XL -- needs YZC-0025, YZC-0028 

---

# Details

Ticket numbers are permanent. `[x]` = closed, `[ ]` = open. Next available: **YZC-0073**.

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

- [x] **[YZC-0057] Cyclic / mutually-recursive type declarations**

  two-pass sema: collect all top-level type names first, then resolve field types.
  Implemented: `AnalyzeFile` first pass pre-registers stubs; `analyzeStructBoc` reuses
  stub pointer so forward/mutual refs stay valid. Golden test: `66_forward_type_ref.yz`.

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

### [x] YZC-0026 — Generics: Explicit Constraint Declaration ✓

`thing T Talker` declares `T` must implement `Talker`; additive with inference.
Multiple constraints supported: `T Talker Serializable`.

- [x] Parser — `T Constraint` optional suffix after single-uppercase type param; `parseConstraintList` collects trailing TYPE_IDENTs; new `TypeParamDecl` AST node for body-context form (`V Talker` as a statement)
- [x] Sema — `StructType.ExplicitConstraints map[string][]string`; constraints stored from both `TypeParamDecl` (body) and `BocParam.Constraints` (signature); pre-scan updated for `TypeParamDecl`; abstract method return types now correctly propagated from signature when body is nil
- [x] IR — `StructDecl.ExplicitConstraints`; lowerer propagates from sema; `isVariantBoc`/`lowerVariantBoc` accept `TypeParamDecl` elements
- [x] Codegen — `buildTypeParamConstraints` emits `[V Talker]` (single), `[V interface{A;B}]` (multiple), or `[V any]` (none); replaces inline loop in both struct and variant paths
- [x] Golden test 76 — `Box[V Describable]` + `Animal` satisfying `Describable`
- [ ] Spec 04 — update

### YZC-0070 — Anonymous boc literal as structural interface value

A boc literal with inner boc-valued fields (`{ describe: { "a boc" } }`) should satisfy
a structural interface constraint at the call site:

```yz
Describable: { describe #(String) }

Box: {
    V Describable
    value V
}

main: {
    c : Box(value: { describe: { "a boc" } })  // should work
    print(c.value.describe())
}
```

Currently fails with a constraint-violation sema error because boc literals are typed as
`BocType` (a plain function type) rather than as anonymous structs.

#### What needs to change

**Sema** — In `analyzeExpr` for `*ast.BocLiteral`, detect when the literal has named
boc-valued fields (the existing `hasInnerBocsOrMethods` predicate already covers this).
When true, type the literal as an anonymous `StructType{IsSingleton:true}` whose fields
are the inner boc-field names and types, rather than as a `BocType`. This makes
`typeHasMethod` find the methods during constraint checking, and makes structural
compatibility with interfaces work.

**Lowerer** — `lowerBocLitExpr` currently always emits a `ClosureExpr` (Go func
literal). When the sema type for the `BocLiteral` is a `StructType`, generate an
anonymous Go struct type instead:
1. Assign a unique name `_anonBoc<N>` (counter on the lowerer).
2. Build a `StructDecl{Name:"_anonBoc0", NoConstructor:true}` with one `MethodDecl`
   per inner boc-valued field (using the existing `lowerMethod` helper).
3. Collect these into `l.anonDecls []*StructDecl` on the lowerer; prepend to
   `f.Decls` after the file is fully lowered.
4. Return `&_anonBoc0{}` as the call-site expression.

**Codegen** — No changes needed: `StructDecl` with `NoConstructor=true` already
emits the struct type + methods without a constructor function.

#### Edge cases to defer

- Anonymous boc literals that capture outer variables (closures over `self.*` fields).
  These require storing captured values in the anonymous struct as fields.
- Nested anonymous bocs inside an anonymous boc method body.
- Multiple uses of structurally identical anonymous boc patterns (dedup opportunity).

#### Acceptance criteria

- `Box(value: { describe: { "a boc" } })` compiles and runs.
- The anonymous struct type satisfies the `Describable` Go interface.
- All existing tests continue to pass.
- New golden test (e.g. test 77) covers the pattern.

- [ ] Sema — type boc literals with inner boc fields as anonymous `StructType`
- [ ] Lowerer — emit anonymous Go struct type + methods; collect as `anonDecls`
- [ ] Golden test 77 — anonymous boc literal satisfying interface constraint

### [x] YZC-0072 — Inline anonymous interface constraint in type params: `V #(method #(T))` ✓

Allow a generic type parameter to be constrained by an inline anonymous interface signature instead of requiring a named interface:

```yz
// Desired — inline constraint, no separate Describable declaration needed
Box: {
    V #( describe #(String) )
    value V
}

// Equivalent to (already works):
Describable #( describe #(String) )
Box: {
    V Describable
    value V
}
```

#### Analysis of related syntax

**`Foo #(describe #(String))`** (named, no body) — already supported via `analyzeBocDeclNode`. Any uppercase-name boc declaration with no body creates a structural interface. You can then use `V Foo` as an explicit constraint (YZC-0026). So the named form is fully working.

**`foo #(describe #(String)) = { describe: { "hola" } }`** — this is a concrete boc implementation in expanded form; `describe` is a callback parameter. Not a type constraint — this is a call signature.

**`V #( describe #(String) )`** — NOT currently supported. The parser sees `GENERIC_IDENT` followed by `#` (not `TYPE_IDENT`), falls through to `parseTypedDecl`, and treats `V` as a field name with boc type `#(describe #(String))`. This creates a regular boc-typed field called `V`, not a constrained type parameter.

#### Required changes

**Parser** — in `parseStatement`, before the `GENERIC_IDENT TYPE_IDENT` → `parseTypeParamDecl` check, add:

```go
// V #(...) — type param with inline anonymous constraint
if tok.Type == token.GENERIC_IDENT && p.peekAt(token.HASH) {
    return p.parseInlineConstraintTypeParam()
}
```

`parseInlineConstraintTypeParam` should:
1. Consume the GENERIC_IDENT (name)
2. Parse the `#(...)` as a `BocTypeExpr` (reuse `parseBocTypeExpr`)
3. Return a `TypeParamDecl{Name: name, InlineConstraint: bocTypeExpr}`

**AST** — add `InlineConstraint *BocTypeExpr` to `TypeParamDecl` (alongside the existing `Constraints []TypeExpr` for named constraints).

**Sema** — in `storeExplicitConstraints` (or new helper): when `InlineConstraint` is present, synthesise an anonymous interface name (e.g. `_V_constraint`) and register it as a `StructType{IsInterface:true}` in the current scope, then store it as the explicit constraint for the type param.

**IR/Codegen** — no changes needed; the explicit constraint path already handles named interfaces; the anonymous one just gets a generated name.

#### Acceptance criteria

`Box: { V #(describe #(String)); value V; desc #(String) { value.describe() } }` compiles and runs without a separate named interface declaration.

- [x] AST — `TypeParamDecl.InlineConstraint *ast.BocTypeExpr`
- [x] Parser — detect `GENERIC_IDENT HASH` before `isBocDeclStart`; route to `parseInlineConstraintTypeParam`
- [x] Sema — `storeInlineConstraint` synthesises `_StructParamConstraint` at file scope; stored in `ExplicitConstraints`
- [x] Lowerer — `emitSyntheticInterface` emits `InterfaceDecl` for synthetic names before the struct
- [x] Golden test 78 — `V #(method #(T))` inline constraint used and satisfied
- [ ] Spec 04 — document inline constraint syntax

### YZC-0071 — Implicit constraint synthesis for type params used in method params

When a struct has a bare type param `V` (no explicit constraint) and a method uses `V` as a **method parameter type** and calls methods on it, the compiler must infer the required constraint from those calls and emit it in the Go type parameter.

#### Problem

The existing inferred-constraint mechanism (`activeConstraints`) works when `V` is a struct **field** type and methods are called via `self.field.method()`. It does NOT currently work when `V` appears as a method **parameter** type:

```yz
Foo: {
    V
    do #(value V) {
        value.hola()   // V must have hola() — but constraint is NOT inferred
    }
}
```

Generated Go fails: `V any` doesn't have `hola()`.

#### Desired behaviour

The compiler synthesises the constraint from usage. Two equivalent spellings the user may use:

**Option A — named interface** (compiler generates):
```yz
Foo: {
    V Holer    // compiler synthesises Holer
    ...
}
Holer: {
    hola #()
}
```

**Option B — inline constraint syntax** (user writes):
```yz
Foo: {
    V #( hola #() )
    ...
}
```

For now, Option A (infer from method-param usage, emit as `interface{ ... }` inline in Go type param, same as the existing field-usage path) is sufficient.

#### Root cause

`analyzeCall` records constraints when `a.activeConstraints != nil` and the receiver is `*GenericType`. This fires for field accesses because `self.value.hola()` resolves `self.value` to `*GenericType{V}`. But for a method parameter `value V`, the resolution of `value` inside the method body also returns `*GenericType{V}` — so the same path *should* fire.

Investigate whether `activeConstraints` is nil when the method body is entered (e.g. if `analyzeBocDeclNode` clears it), or if there is a scope/resolution issue preventing `value` from being typed as `*GenericType` inside the method body.

#### Acceptance criteria

```yz
Foo: {
    V
    do #(value V) {
        value.hola()
    }
}
Holer: {
    hola #()
}
main: {
    h : Holer()
    Foo(do: h)
}
```
Compiles and runs without explicit `V Holer` in `Foo`.

- [ ] Investigate why `activeConstraints` doesn't fire for method-param receivers
- [ ] Fix constraint recording for method-param usage of generic type params
- [ ] Golden test 78 — bare `V` inferred from method-param usage
- [ ] Spec 04 — document implicit constraint inference

### YZC-0027 — `:` as Type Alias

`Name : SomeType` declares a type alias. Depends on YZC-0066 for the unified model; can be implemented as a limited special form before YZC-0066 lands (emit `type Name = GoType` in Go, no `#()` metatype required).

- [ ] Parser — distinguish from `Name TypeExpr` (typed decl) and `name : value` (short decl)
- [ ] Sema — register alias; resolve as aliased type
- [ ] Lowerer — emit `type Name = GoType`
- [ ] Spec 04 — add
- [ ] Deferred to YZC-0066: generic instantiation via alias (`StringList : List(String)`), associated type binding (`Node : User` inside a boc)

### [x] YZC-0066 — Associated Types: `#()` metatype, T fields, type aliases, call-site unification ✓

Unified model for generics, type aliases, and associated types. See `docs/Features/Path Dependent Types.md`.

Full implementation plan: [`docs/Implementation/yzc-0066-plan.md`](yzc-0066-plan.md)

Note: was originally named "Path-Dependent Types" — name corrected; YZC-0030 covers the remaining path-dependent resolution for abstract types.

- [x] Sema — `#()` recognized as metatype; bare GENERIC_IDENT field given implicit `#()` type
- [x] Sema — type fields in constructors (`List(Int)` binds `T = Int`) — Go inference handles monomorphization
- [x] Sema — `g.Node` in type position resolves when `g`'s concrete type is statically known
- [x] Sema — type variable inference: unify GENERIC_IDENT against call-site argument types (`GenericInstType`)
- [x] `Node : User` inside a boc body treated as type alias (IsTypeField), not value alias
- [x] Golden tests: 68 (type alias), 69 (implicit TypeParams), 70 (path-dependent), 71 (type var unification)
- [ ] Spec 04 — generics section; Spec 05 — associated types section

### YZC-0028 — Compile-Time Bocs (`Compile` interface)

Any boc with `Schema #()` and `run #(Boc, Boc)` satisfies `Compile`. Depends on: YZC-0025, YZC-0026, YZC-0027, YZC-0030, YZC-0066, YZC-0059.

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

### [x] YZC-0030 — Path-Dependent Types: abstract `g.Node` resolution ✓

`process #(g Graph, n g.Node)` — sema resolves `g.Node` against the **abstract** type of `g` (interface parameter), not just the concrete static type. Design resolved; see `docs/Features/Path Dependent Types.md` and `docs/Features/Associated Types.md`.

Note: was originally named "Associated Types" — name corrected; the associated-type machinery (YZC-0066) is now complete. Depends on YZC-0067: until Graph is emitted as a Go interface, passing a concrete subtype (SocialGraph) as an abstract parameter (Graph) fails Go's type checker.

When `g` is a concrete local variable, `g.Node` already resolves correctly (done in YZC-0066). This ticket covers the abstract case: two different `g1: Graph` and `g2: Graph` values have distinct, incompatible `g1.Node` vs `g2.Node` types at the type-checker level.

- [x] Sema — `g.Node` in type position when `g` has an abstract/interface type — PathDependentType returned by resolveTypeExpr; call-site check in analyzeCall
- [x] Sema — enforce `g1.Node` and `g2.Node` are distinct types even when both satisfy `Graph` — error test 20
- [x] Lowerer — sema substitutes concrete return type at call site; goTypeForVar uses resolved *StructType, var gets concrete Go type (e.g. `*User`) when called from concrete context
- [x] Golden test: Graph/SocialGraph/accept — test 72 passes; *SocialGraph satisfies Graph interface

### [x] YZC-0067 — Emit Go interfaces for structural Yz types ✓

In Yz, any struct that has the required fields/methods satisfies a type structurally. In Go, this only works when the target type is a Go `interface`, not a Go `struct`. Currently all Yz boc types (including those with only method fields) are emitted as Go structs, so passing `*SocialGraph` where `*Graph` is expected fails Go's type checker.

The fix: boc types that have `IsInterface=true` (all fields are BocType methods) should be emitted as Go interfaces. Any Yz struct that satisfies the interface structurally will then automatically satisfy the Go interface, no casting required.

YZC-0030 depends on this: path-dependent type params (`g Graph, n g.Node`) resolve correctly in sema but the generated Go doesn't compile when passing `*SocialGraph` as `*Graph` until Graph is a Go interface.

- [x] Codegen — emit `type Name interface { ... }` for `IsInterface=true` structs instead of `type Name struct { ... }`
- [x] Codegen — emit Go interface methods (no receiver, no `std.Cown` embed)
- [x] Lowerer — when a param type is an interface, pass the arg directly (no pointer wrapping)
- [x] Sema — extend `IsInterface` detection: a boc type with a mix of abstract type fields (`Node #()`) and method fields should also be treated as an interface
- [x] Golden test: Graph/SocialGraph/process — `process(sg, u)` compiles in Go with `sg *SocialGraph` satisfying `Graph` interface
- [x] Verify existing `IsInterface` golden tests (structural typing tests) still pass

### [x] YZC-0068 — GoStore type mismatch for path-dependent return types ✓

Functions with path-dependent return types (e.g. `makeNode #(g Graph, g.Node)`) are emitted as singleton boc methods that return `*std.Thunk[any]` — Go does not support generic methods, so the return type cannot be parameterized. At the call site, sema correctly resolves the return type to a concrete type (e.g. `*User`) and the generated variable is `var node *User`, but `std.GoStore(_bg0, MakeNode.Call(sg), &node)` fails the Go type checker because the thunk is `*std.Thunk[any]` while the destination pointer is `*User`.

Conformance test 73 (`73_path_dependent_return.yz`) has the correct golden `.go` file but the generated Go does not compile end-to-end due to this mismatch.

Options:
- **Option A** — Add `std.GoStoreAny[T any](bg *BocGroup, thunk *Thunk[any], dest *T)` runtime helper that does a type-assertion `thunk.Force().(T)`. Simple, no codegen change.
- **Option B** — Emit path-dependent-return functions as Go free functions with a type parameter instead of boc methods (`func MakeNode[N any](g Graph) *std.Thunk[N]`). Requires codegen changes; more correct but complex.

Recommended: Option A (runtime helper) as the pragmatic fix; Option B as a follow-on if the generic method restriction is lifted.

- [x] Add `GoStoreAny[T any]` to `compiler/runtime/rt/core.go`
- [x] Codegen — emit `GoStoreAny` when `GoStore` has a `*Thunk[any]` source and concrete `*T` dest
- [x] Golden test 73 updated; end-to-end compilation verified

### [x] YZC-0069 — Call-site type variable unification (Phase C generics) ✓

This is Phase C of YZC-0066. Phases A and B (direct type variable unification) already work — see golden test 71 (`71_type_var_unification.yz`). Phase C covers the two harder cases where the existing unifier falls short.

#### What already works (test 71)

When a type variable appears *directly* as a parameter type, it is unified with the argument type at the call site and substituted into the return type:

```yz
identity #(val A, A)    // A = typeof(val) — direct match → var n std.Int  ✓
wrap #(val A, Box(A))   // A = typeof(val), return Box(A) → Box[Int]       ✓
```

#### Gap 1 — Structural (nested) unification

A type variable buried inside a generic wrapper in the parameter type:

```yz
map #(collection List(A), fn #(A, B), List(B))
//               ^^^^^^^
// Matching List(Int) against List(A) requires recursing into the generic's type arguments.
// Current behavior: A stays GenericType{A} → emitted as `any` in Go.
```

The current unifier only matches `argType == GenericType{A}` directly. It does not recurse into `GenericInstType` wrappers.

#### Gap 2 — Boc-argument inference

A type variable that is only knowable from the *return type* of a boc (closure) argument:

```yz
map #(collection List(A), fn #(A, B), List(B))
//                            ^^^
// B comes from what fn returns, not from fn's type directly.
// Need to analyze the closure argument, observe it returns Int, then bind B = Int.
```

This requires a two-pass strategy: first unify all non-boc arguments (binding as many variables as possible), then analyze boc-literal arguments with the partial substitution already in scope, then unify the boc's inferred return type to bind remaining variables.

#### Why this is different from YZC-0030

YZC-0030 (path-dependent types) is field lookup: `g.Node` means "find the `Node` field on `g`'s concrete type." It is a single-step table lookup on a named parameter.

Phase C is a proper unification algorithm: match a *pattern* type (containing free variables) against a *concrete* type (no free variables), collecting a substitution map. The two problems are orthogonal.

#### Implementation

Add a `unify` function in `compiler/internal/sema/`:

```go
// unify matches pattern (which may contain GenericType free variables) against
// concrete, adding bindings to subst. One-directional: pattern has the free vars.
func unify(pattern, concrete Type, subst map[string]Type) {
    switch p := pattern.(type) {
    case *GenericType:
        if existing, ok := subst[p.Name]; ok {
            // consistency check: existing binding must match concrete
        } else {
            subst[p.Name] = concrete
        }
    case *GenericInstType:
        // e.g. List(A) vs List(Int): recurse into type arguments
        if c, ok := concrete.(*GenericInstType); ok && p.Name == c.Name {
            for i := range p.Args {
                unify(p.Args[i], c.Args[i], subst)
            }
        }
    case *BocType:
        // e.g. #(A, B) vs #(Int, String): unify params and returns pairwise
        if c, ok := concrete.(*BocType); ok {
            for i := range p.Params { unify(p.Params[i].Type, c.Params[i].Type, subst) }
            for i := range p.Returns { unify(p.Returns[i], c.Returns[i], subst) }
        }
    }
}
```

Wire it into `analyzeCall` when the callee has generic type variables:

1. **Pass 1** — non-boc arguments: for each `(param, arg)` pair where the arg is not a boc literal, call `unify(param.Type, argType, subst)`.
2. **Pass 2** — boc-literal arguments: analyze each closure argument with the partial substitution applied to its expected parameter types; infer the closure's return type; call `unify(param.ReturnType, closureReturnType, subst)` to bind any remaining variables.
3. **Apply substitution** to the callee's declared return type to produce the concrete return type for this call expression.

The substitution application (`applySubst(t Type, subst map[string]Type) Type`) mirrors `unify` structurally: replace `GenericType` leaves, recurse into `GenericInstType` args and `BocType` params/returns.

#### Implementation notes (actual vs. planned)

`unifyTypes` and `substituteType` already existed and already handled `GenericType`, `ArrayType`, `BocType`, and `GenericInstType` structurally — the two-pass unification was already wired into `analyzeCall`. The real gap was that generic struct **constructor calls** returned bare `StructType` (no concrete type args), so passing the result to a generic HOF couldn't unify the type variables. Also, the lowerer used `goType` (which emits the raw variable name e.g. `"A"`) instead of `goTypeForVar` in GoStore paths, producing invalid Go like `var v A`.

#### Checklist

- [x] `unifyTypes(formal, actual, bindings)` already in sema — handles GenericType/ArrayType/BocType/GenericInstType
- [x] `substituteType(t, bindings)` already in sema — mirrors unify structurally
- [x] Two-pass unification already wired in `analyzeCall` (boc-literal args produce a BocType, unifyTypes handles BocType returns)
- [x] Constructor calls: `analyzeCall` for `*StructType` with TypeParams now infers concrete type args and returns `GenericInstType{Name,[concreteArgs]}` — `Box(value:42)` → `GenericInstType{Box,[Int]}`
- [x] `fieldType` extended with `*GenericInstType` case: looks up base struct, builds subst TypeParams→TypeArgs, returns substituted field type
- [x] `isBocMethodCall` in lowerer extended to recognise `GenericInstType` as struct-like (method calls on generic struct instances still treated as boc calls)
- [x] GoStore and method-body paths use `goTypeForVar` (with `"any"` fallback) instead of `goType` — prevents invalid `var v A` when type var is unresolved
- [x] Golden test 74 (`74_phase_c_generic_hof.yz`): `transform` (boc-arg inference) + `unwrap(Box(...))` (generic struct → HOF) — both result vars typed concretely
- [x] All 52 golden + 20 error tests pass

### YZC-0031 — Scalar Types in Yz Source (uppering)

`Int/String/Bool/Decimal/Unit` move from Go to `stdlib/` with `compile-time:[Native]`. Depends on: YZC-0025, YZC-0028.

- [ ] Define `compile-time:[Native]` infostring semantics
- [ ] Move scalar types to `stdlib/`
- [ ] Annotate native ops per method
- [ ] Implement higher-level methods in Yz
- [ ] Remove all primitive-type special-casing from the compiler
- [ ] `Bool.&&`/`||` — rewrite as lazy closure-taking boc methods
