#impl
Completed tickets. Ticket numbers are permanent.

---

- [x] **[YZC-0076] Existential associated types: opaque-token / path-identity tracking** — closed, not implemented

  Original motivation was heterogeneous macro arrays with per-element schema
  validation, requiring runtime path-dependent types or first-class existentials.
  Macro dispatch changed to compile-time type name resolution (YZC-0059) —
  the compiler always has the concrete type in hand, so the problem never arises.
  Yz has associated types (YZC-0066) and structural compatibility but not runtime
  path-dependent types or first-class existential types.

- [x] **[YZC-0097] Annotation metadata contract for project and dependency configuration** *(replaces YZC-0041)*

  Defined the format for `project:` and `dependencies:` as passive annotation
  metadata in `project.info`. Lock file is Yz array syntax. `yz fetch` owns
  resolution and pinning; compiler never fetches. See
  [Dependencies](../Features/Dependencies.md).

- [x] **[YZC-0059] Design: macro interface interaction** — [resolved](../Questions/solved/Macro%20Interface%20Interaction%20Design.md)

  Settled annotation taxonomy: runtime metadata (typed boc, accessed via
  `value.annotation.field`), macro dispatch (uppercase type name in annotation
  body, `run #(Boc, Boc)`), and native extension (`GoSource`, hardcoded
  first-pass before type resolution). Conditional compilation, compiler
  directives, and ABI layout deferred to future work.

- [x] **[YZC-0022] Multiple source roots**

  `yzc build <project-dir> [extra-roots...]` — multiple positional source root
  directories, each contributing FQNs to the same namespace with its prefix
  stripped. First arg owns `target/`. Same-FQN same-name collision across roots
  is a compilation error. No implicit default root. Added `examples/multi_root/`.

- [~] **[YZC-0040] Smart Nesting / Namespace Flattening**

  `house/house.yz` flattens to `house.method`. Superseded by YZC-0085.

- [x] **[YZC-0093] Uppercase root file (`Foo.yz`) always-wrap: example + spec §9 clarification**

  Fixed `build.go` to set the correct `TokType` (`token.LookupIdent`) on the
  synthesized wrapper ident, so uppercase file names (e.g. `Pet.yz`) produce a
  `StructType` in sema rather than a `BocType`. Added `examples/root_type/`
  with `Pet.yz` + `main.yz` as a conformance test. Updated spec §9 Invariant 1
  with notes on uppercase files and explicit same-named inner boc files.

- [x] **[YZC-0092] Always-wrap root files; `main()` as explicit entry invocation**

  Every `.yz` file's content is wrapped in a `ShortDecl` named after the file
  (`IsFileWrapper=true` in `build.go`). When the wrapper contains an explicit
  same-named inner boc (e.g. `main: {}` inside `main.yz`), the lowerer unwraps
  it so inner declarations become package-level items. Files without an inner boc
  (root_wrap pattern) keep the wrapper as the singleton body. Cross-file type
  resolution re-registers the inner same-named symbol via sema's
  `fileWrapperScopes` map. All examples consolidated into single files.

- [x] **[YZC-0021] Directory and file bocs**

  Files in sub-directories are auto-wrapped in a boc named after the file
  before analysis (`build.go` `compilePackageDir`). `ledger/ledger.yz` content
  becomes `ledger: { ... }`, making FQN `ledger.ledger` accessible from the
  parent package. Root files without an explicit same-name top-level boc are
  also auto-wrapped (`world.yz` with free-floating code → callable as `world()`).
  Implements spec §9 Invariants 1+2.

- [x] **[YZC-0089] Invariant 5: `foo.yz` + `foo/` coexistence — loader merge**

  Loader merge implemented in `build.go` `compileProject`: detects root-level
  `foo.yz` + `foo/` pairs, injects wrapped sub-directory declarations into the
  root boc literal, removes `foo/` from separate-package compilation.
  Sema now resolves `foo.bar` correctly after the merge.
  Nested singleton codegen blocker extracted to YZC-0091.
  Test: `examples/_wip/subdir_coexist` — promoted once YZC-0091 lands.

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

- [x] **[YZC-0010] HOF iteration + cown happens-before**

  `.filter`, `.each` as sync Go closures. Golden test 27.

- [x] **[YZC-0036] While loop yield and external caller interleaving**

  BocDecl singletons use `std.Schedule`; recursive self-calls marked `IsRecursive`.

- [x] **[YZC-0011] Named arguments in constructor calls**

  `lowerStructArgs` reorders by field declaration order; `lowerNamedArgs` for BocDecl calls. Golden test 59.

- [x] **[YZC-0012] Multiple return values**

  `x, y = swap(x, y)` — multi-assign LHS. Multiple trailing non-Unit expressions from a boc body produce a `_<name>BocResult` plain struct; call sites Force() the thunk and destructure into individual variables. Golden 91.

- [x] **[YZC-0015] Non-word boc names**

  `balance+= #(amount Int) { ... }` — parser accepts `NON_WORD` token and maps to Go-safe name.

- [x] **[YZC-0018] Bool methods `&&` / `||`**

  `Bool.Ampamp` / `Bool.Pipepipe` in yzrt. Golden test 53.

- [x] **[YZC-0020] Compiler homoiconic dump — backtick interpolation**

  backtick inside a string triggers homoiconic representation. Golden test 60.

- [x] **[YZC-0037] Decimal type end-to-end**

  `std.Decimal` with arithmetic, comparisons, `to_str`. Golden test 58.

- [x] **[YZC-0038] `Result(T,E)` type**

  Implemented as user-level Yz code (no compiler built-in needed). Fixed the general sum-type
  issue: when a generic variant constructor doesn't constrain all parent type params (e.g.
  `Err(error E)` in `Result[T,E]` — `T` is unconstrained), the lowerer now emits explicit Go
  type args (`NewResultErr[std.Int, std.String](...)`). Sema fills in unbound type params from
  the call site's `expectedType` (TypedDecl annotation). Golden test 86.

- [x] **[YZC-0043] Captured variable reference semantics**

  Decision: always reference capture. Bocs are reference types; Go closures capture by reference. No copy semantics, no implementation work needed.

- [x] **[YZC-0045] Default values in type-only boc declarations (interfaces)**

  Struct field defaults (`next: Option.None()`) implemented: `DefaultExpr ast.Expr` stored in
  `StructField`; lowerer emits the default expression when field is omitted from a constructor
  call. Interface-level defaults (`Greeter #(name String = "Alice")`) deferred; depends on YZC-0011.

- [x] **[YZC-0046] `${}` interpolation requires `to_str`**

  sema checks for `to_str #(String)` on the interpolated type. Depends on: YZC-0020.

- [x] **[YZC-0078] `print` should require `String`; use `` "`x`" `` for debug output**

  Currently `print(a)` accepts any value and calls `Stringify` (homoiconic `String()` method).
  This conflates two distinct intents:

  - **Display**: `print("${a}")` — user-facing output; requires `to_str #(String)` on the type
  - **Debug**: `print("`a`")` — homoiconic structural dump; uses `String()`, no `to_str` needed

  `print(a)` silently falls through to the debug path, making it easy to accidentally ship
  debug output. The fix: restrict `print` to `String` only; `print(a)` where `a` is not a
  `String` becomes a sema error with message _"print requires String; use \"`a`\" for debug
  output or \"${a}\" for display output"_.

  Design decision:  `print` should be a regular Yz boc with the signature:
  `print #(String, nl:true)` // new line defaults to true
   that enforces the constraint naturally. 
  

  Current behaviour to preserve:
  - `print("hello")` — valid ✓
  - `print("${a}")` — valid when `a` has `to_str` ✓
  - `` print("`a`") `` — valid; always works ✓
  - `print(a)` — currently valid; should become a sema error after this ticket: a is not a string

- [x] **[YZC-0047] Cycle detection in homoiconic `Stringify`** ✓

  - [x] Runtime — per-goroutine visited set in `Stringify`/`StringifyRepr` via `sync.Map`
        keyed on `(goroutineID, ptr)`; cyclic references print as `TypeName(...)`
  - [x] Runtime — nil pointer guard in both functions (interface-wrapped nil no longer panics)
  - [x] Unit tests — self-cycle, indirect cycle, linear chain (no false positive), concurrent
        same-pointer (four tests in `runtime/rt/rt_test.go`)
  - [x] Golden test 84 — cyclic linked list via locally-declared `Option` variant; `b.next =
        Option.Some(a)` creates a cycle; `print(a)` emits `Node(..., Node(..., Node(...)))` ✓

- [x] **[YZC-0077] Recursive struct types: cycle guard in `IsCompatibleWith` + sema support** ✓

  - [x] Sema — pointer equality check `if t == u { return true }` at top of `*StructType`
        case in `IsCompatibleWith`; breaks infinite recursion without changing the interface
  - [x] (No lowerer/codegen change needed — struct fields of struct type already emit as `*Node`)
  - [x] Golden test 83 — `Node: { value Int; next Node }` + function over it compiles and runs

- [x] **[YZC-0061] Structured singleton: TypedDecl-with-value field missing `self.`**

  `collectFieldNames` gating removed. Golden test 63.

- [x] **[YZC-0017] Dict optional access**

  **Invariant:** For `d [K:V]`, `d[k]` returns `Option(V)` and `d[k] = v` takes `V`. The `V` inside `Option(V)` is the same type parameter — all constraints on `V` in the declaration carry through unchanged.

- [x] **[YZC-0087] Dict assignment syntax: `d["key"] = value`**

  The compiler currently accepts `d["key":value]` (key:value pair notation). Replace with standard assignment syntax `d[key] = value` to match user expectations and free up the oddness budget. Feature doc updated.

---

## Infrastructure

- [x] **[YZC-0033] Compiler deep review against settled spec**

  all four sub-items resolved: BOC singletons, `foo.param` accessible after call, error messages say "returns nothing", all bocs serialized through cown.

- [x] **[YZC-0032] Rename `BocWithSig` → `BocDecl`**

  done throughout AST, sema, lowerer, and spec/02.

- [x] **[YZC-0002] Cross-package support**

  Fixed: `isSingletonExport` in lowerer now handles `StructType{IsSingleton:true}` exports,
  resolving `pkg.singleton.method()` calls correctly. Example promoted to `examples/cross_pkg/`.

- [x] **[YZC-0057] Cyclic / mutually-recursive type declarations**

  two-pass sema: collect all top-level type names first, then resolve field types.
  Implemented: `AnalyzeFile` first pass pre-registers stubs; `analyzeStructBoc` reuses
  stub pointer so forward/mutual refs stay valid. Golden test: `66_forward_type_ref.yz`.

---

## Major Features

### [x] YZC-0025 — Annotations: content is a boc body ✓

Annotation delimiter stays backtick; content is full Yz syntax, parsed and type-checked, never executed. Intersection with Native annotations (YZC-0058).

- [x] AST — `Annotation` holds `*BocLiteral`; `BocLiteral.Annotation *Annotation`
- [x] Lexer — `ANNOTATION` token type; `scanAnnotation()` scans backtick-delimited content
- [x] Parser — sub-parser re-lexes and re-parses annotation content as boc body
- [x] Sema — traverses annotation body elements for type checking
- [~] Codegen — deferred to YZC-0088; depends on macro system (YZC-0028) to define what "declaration metadata" means
- [x] Spec 01 — interpolation restrictions documented (§1.14)

### [x] YZC-0026 — Generics: Explicit Constraint Declaration ✓

`thing T Talker` declares `T` must implement `Talker`; additive with inference.
Multiple constraints supported: `T Talker Serializable`.

- [x] Parser — `T Constraint` optional suffix after single-uppercase type param; `parseConstraintList` collects trailing TYPE_IDENTs; new `TypeParamDecl` AST node for body-context form (`V Talker` as a statement)
- [x] Sema — `StructType.ExplicitConstraints map[string][]string`; constraints stored from both `TypeParamDecl` (body) and `BocParam.Constraints` (signature); pre-scan updated for `TypeParamDecl`; abstract method return types now correctly propagated from signature when body is nil
- [x] IR — `StructDecl.ExplicitConstraints`; lowerer propagates from sema; `isVariantBoc`/`lowerVariantBoc` accept `TypeParamDecl` elements
- [x] Codegen — `buildTypeParamConstraints` emits `[V Talker]` (single), `[V interface{A;B}]` (multiple), or `[V any]` (none); replaces inline loop in both struct and variant paths
- [x] Golden test 76 — `Box[V Describable]` + `Animal` satisfying `Describable`
- [ ] Spec 04 — update

### [x] YZC-0070 — Anonymous boc literal as structural interface value ✓

A boc literal with inner boc-valued fields (`{ describe: { "a boc" } }`) should satisfy
a structural interface constraint at the call site.

- [x] Sema — type boc literals with inner boc fields as anonymous `StructType`
- [x] Lowerer — emit anonymous Go struct type + methods; collect as `anonDecls`
- [x] Golden test 88 — anonymous boc literal satisfying interface constraint

### [x] YZC-0072 — Inline anonymous interface constraint in type params: `V #(method #(T))` ✓

- [x] AST — `TypeParamDecl.InlineConstraint *ast.BocTypeExpr`
- [x] Parser — detect `GENERIC_IDENT HASH` before `isBocDeclStart`; route to `parseInlineConstraintTypeParam`
- [x] Sema — `storeInlineConstraint` synthesises `_StructParamConstraint` at file scope; stored in `ExplicitConstraints`
- [x] Lowerer — `emitSyntheticInterface` emits `InterfaceDecl` for synthetic names before the struct
- [x] Golden test 78 — `V #(method #(T))` inline constraint used and satisfied
- [ ] Spec 04 — document inline constraint syntax

### [x] YZC-0071 — Implicit constraint synthesis for type params used in method params ✓

- [x] Fix interface classification: structs with method bodies are not interfaces
- [x] Fix constraint recording: `ExplicitConstraints` updated in both sema and IR
- [x] Fix return type: `findInterfaceMethodReturnType` infers concrete return type from matching interface
- [x] Golden test 79 — `Wrapper[V Holer]` inferred from `value.hola()` in method param

### [x] YZC-0073 — Synthesize anonymous interface constraint when no named interface in scope ✓

- [x] Sema — `synthesizeConstraints`: generates `_StructVConstraint` `StructType` for non-builtin methods not matched by named interface
- [x] Sema — registers synthesized interface in file scope; adds to `ExplicitConstraints`
- [x] Lowerer — `emitSyntheticInterface` already handles `_`-prefixed names; no lowerer change needed
- [x] Golden test 80 — synthesized constraint (no named interface in scope)

### [x] YZC-0075 — Existential associated types: implicit erasure (superseded by YZC-0079)

All checklist items completed; approach replaced by YZC-0079 which uses structural bound-checking instead of existential rejection.

- [x] Sema — detect when a path-dependent type's root is an abstract (interface) binding rather than a concrete struct; mark as existential
- [x] Sema — allow method calls on existential `g.Node` when `Node` has a YZC-0074 bound
- [x] Sema — error at the use site when an existential `g.Node` is used in a position that requires a concrete type
- [x] Conformance tests — golden 87 (constrained method call allowed), error 22 (existential violation)

### [x] YZC-0079 — Associated type call-site check: bound-check instead of existential rejection ✓

Replaces YZC-0075. Structural bound-checking: any value satisfying the bound is a valid `g.Node`.

- [x] Revert YZC-0075 existential rejection in `analyzeCall`
- [x] Add bound-compatibility check: when `g` is abstract and `Node` has a bound, verify arg satisfies the bound
- [x] Updated conformance tests: error 22 now tests bound mismatch; golden 87 shows both concrete and abstract `g` working

### [x] YZC-0074 — Constrained associated types ✓

- [x] Sema — `StructField.Bound Type` stores the constraint; `buildAssocTypeBound` creates synthetic `_GraphNodeBound` interface from inline params; `TypeParamDecl` with TYPE_IDENT name handles `Node Sizer` form
- [x] Sema — bind site check in `analyzeCall` (YZC-0074 error) + `IsCompatibleWith` check in `StructType`
- [x] Sema — `fieldType` PathDependentType case returns bound interface so method calls type-check in bodies
- [x] Parser — `TYPE_IDENT + TYPE_IDENT` case routes to `parseTypeParamDecl`; preserves actual token type
- [x] Lowerer — `resolvePDTGoType` emits bound interface as Go type; `isBocMethodCall` extended for PDT-typed values; `lowerStructBoc` emits synthetic interfaces before the containing interface
- [x] Error test 21 — bind site violation: concrete type missing required method
- [x] Golden test 82 — function body calls `node.label()` via constrained `g.Node`; output `1,2`

### [x] YZC-0027 — `:` as Type Alias ✓

- [x] Sema — `analyzeShortDecl` already registers `Bar` with `*StructType{Name:"Foo"}`; no parser change needed
- [x] IR — `TypeAliasDecl{Name, Target}` added
- [x] Lowerer — `lowerTopShortDecl` detects type alias; constructor calls use `st.Name` (not callee id)
- [x] Codegen — emits `type Name = Target`
- [x] Golden test 81 — `Bar: Foo`; both `Foo(...)` and `Bar(...)` constructors work
- [ ] Generic instantiation (`StringList : List(String)`) — deferred to YZC-0084

### [x] YZC-0066 — Associated Types: `#()` metatype, T fields, type aliases, call-site unification ✓

- [x] Sema — `#()` recognized as metatype; bare GENERIC_IDENT field given implicit `#()` type
- [x] Sema — type fields in constructors (`List(Int)` binds `T = Int`) — Go inference handles monomorphization
- [x] Sema — `g.Node` in type position resolves when `g`'s concrete type is statically known
- [x] Sema — type variable inference: unify GENERIC_IDENT against call-site argument types (`GenericInstType`)
- [x] `Node : User` inside a boc body treated as type alias (IsTypeField), not value alias
- [x] Golden tests: 68 (type alias), 69 (implicit TypeParams), 70 (path-dependent), 71 (type var unification)
- [ ] Spec 04 — generics section; Spec 05 — associated types section

### [x] YZC-0067 — Emit Go interfaces for structural Yz types ✓

- [x] Codegen — emit `type Name interface { ... }` for `IsInterface=true` structs instead of `type Name struct { ... }`
- [x] Codegen — emit Go interface methods (no receiver, no `std.Cown` embed)
- [x] Lowerer — when a param type is an interface, pass the arg directly (no pointer wrapping)
- [x] Sema — extend `IsInterface` detection: a boc type with a mix of abstract type fields (`Node #()`) and method fields should also be treated as an interface
- [x] Golden test: Graph/SocialGraph/process — `process(sg, u)` compiles in Go with `sg *SocialGraph` satisfying `Graph` interface
- [x] Verify existing `IsInterface` golden tests (structural typing tests) still pass

### [x] YZC-0068 — GoStore type mismatch for path-dependent return types ✓

- [x] Add `GoStoreAny[T any]` to `compiler/runtime/rt/core.go`
- [x] Codegen — emit `GoStoreAny` when `GoStore` has a `*Thunk[any]` source and concrete `*T` dest
- [x] Golden test 73 updated; end-to-end compilation verified

### [x] YZC-0069 — Call-site type variable unification (Phase C generics) ✓

- [x] `unifyTypes(formal, actual, bindings)` already in sema — handles GenericType/ArrayType/BocType/GenericInstType
- [x] `substituteType(t, bindings)` already in sema — mirrors unify structurally
- [x] Two-pass unification already wired in `analyzeCall` (boc-literal args produce a BocType, unifyTypes handles BocType returns)
- [x] Constructor calls: `analyzeCall` for `*StructType` with TypeParams now infers concrete type args and returns `GenericInstType{Name,[concreteArgs]}`
- [x] `fieldType` extended with `*GenericInstType` case: looks up base struct, builds subst TypeParams→TypeArgs, returns substituted field type
- [x] `isBocMethodCall` in lowerer extended to recognise `GenericInstType` as struct-like
- [x] GoStore and method-body paths use `goTypeForVar` (with `"any"` fallback) instead of `goType`
- [x] Golden test 74 (`74_phase_c_generic_hof.yz`): `transform` + `unwrap(Box(...))` — both result vars typed concretely

### [x] YZC-0029 — Remove `mix`: runtime + spec — PARTIALLY COMPLETE

Compiler removal done.

- [x] Lexer, Parser, Sema, Lowering/Codegen, Golden tests — done
- [x] Spec 09 — remove `mix`; document `Mix` compile implementation

### [x] YZC-0030 — Path-Dependent Types: abstract `g.Node` resolution ✓

- [x] Sema — `g.Node` in type position when `g` has an abstract/interface type — PathDependentType returned by resolveTypeExpr; call-site check in analyzeCall
- [x] Sema — enforce `g1.Node` and `g2.Node` are distinct types even when both satisfy `Graph` — error test 20
- [x] Lowerer — sema substitutes concrete return type at call site; goTypeForVar uses resolved *StructType
- [x] Golden test: Graph/SocialGraph/accept — test 72 passes; *SocialGraph satisfies Graph interface

### [x] YZC-0081 — Singleton-outer nested type factory ✓

- [x] Sema: recognize uppercase struct-literal inside singleton as nested type definition
- [x] Sema: give it a scoped name (`singleton.Type`), register as a `StructType` in scope
- [x] Lowerer: emit nested type as package-level `StructDecl`; singleton field becomes the constructor
- [x] Field access `room.Window(...)` → constructor call on the nested Go struct
- [x] Golden test 90: `room: { Window: { size Int } }` + `room.Window(size: 3)` compiles and runs

### [x] YZC-0083 — Spec consolidation ✓

- [x] Spec 04 — generics: document explicit constraints (YZC-0026), type params, generic instantiation
- [x] Spec 04 — inline anonymous constraint syntax `V #(method #(T))` (YZC-0072)
- [x] Spec 04/05 — associated types: `#()` metatype, type aliases, call-site unification (YZC-0066)
- [x] Spec 04 — type alias `Name : SomeType` (YZC-0027)
- [x] Spec 03 — multi-variable short declaration `x, y : swap(...)`
- [x] Spec 04 — nested type declarations (singleton-outer, struct-outer) (YZC-0081/0082)
- [x] Spec 04 — struct field defaults (YZC-0045), recursive struct types (YZC-0077)
- [x] Spec 04/05 — constrained associated types + abstract g.Node resolution (YZC-0074/0079)

### [x] YZC-0085 — Module system design: file/dir invariants, `name.info` companion ✓

- [x] `docs/Features/Smart Nesting and Namespace Flattening.md` → move to `Replaced features/`
- [x] `docs/Features/Annotations.md` (was `Info strings.md`) — two declaration forms
- [x] `docs/Features/Code organization.md` → rewrite around these invariants
- [x] `spec/09-modules-and-organization.md` → rewrite
- [x] `docs/Features/Macros.md` (was `Compile Time Bocs.md`) — macro catalogue
- [x] New: `docs/Features/Macros/` subdir with individual macro docs

### [x] YZC-0086 — Rename: infostring → annotation, compile-time boc → macro, `_name.yz` → `name.info` ✓

- [x] `docs/Features/Info strings.md` → `docs/Features/Annotations.md`
- [x] `docs/Features/Compile Time Bocs.md` → `docs/Features/Macros.md`
- [x] `docs/Features/Compile Time Bocs/` → `docs/Features/Macros/`
- [x] `spec/09-modules-and-organization.md` — invariant 5
- [x] `spec/01-lexical-structure.md` — backtick comment
- [x] `docs/Features/Code organization.md` — invariant table
- [x] `docs/Questions/` — terminology updates
- [x] `compiler/internal/ast/ast.go` — `InfoString` → `Annotation`
- [x] `compiler/internal/parser/parser.go` — function/var renames
- [x] `compiler/internal/sema/analyzer.go` — case arm renames
