# update log  

## 2026-09-28 (1)
- **Update**: [tasks.md](./tasks.md) -- Closed YZC-0107; bumped golden-test count to 113.
- **Update**: [tasks-detail.md](./tasks-detail.md) -- Removed the YZC-0107 Bugs entry (moved to tasks-done.md).
- **Update**: [tasks-done.md](./tasks-done.md) -- Added completion write-up for YZC-0107: `lowerExpr`'s `*ast.BinaryExpr` case unconditionally force-wrapped a non-scalar boc-call operand, deadlocking when the operand was a non-leaf, self-recursive sibling call whose cown the enclosing method already held (e.g. `check() { countdown(3) == 0 }`). Diagnosed via a user question challenging whether `==` itself was the cause -- disproved by rewriting the same call through a hand-written `is` method (deadlocked identically, confirming the trigger is positional -- inline vs. bound -- not operator-specific) and by hand-verifying the fix shape (`r : countdown(3); r == 0`) already worked with zero compiler changes before automating it.
- **Update**: `internal/ir/lower.go` -- Added `hoistHeldCownCalls`/`isNonLeafHeldCownCall`, run at the top of `lowerBocBody`: splices any non-leaf held-cown call found nested inside a larger expression (binary operand, call argument, `?` condition, match subject, string interpolation part, array/dict literal element) out into its own synthetic `name : call(...)` element beforehand, reusing the already-correct bound-call Schedule/BocGroup/Wait lowering instead of adding new runtime machinery.
- **Update**: [concurrency-design.md](./concurrency-design.md) -- Replaced the YZC-0107 "Known gap" note with a "Resolved gap" note describing the hoist mechanism.
- **Creation**: `119_reentrant_expr_position_force` golden test + `.output` sidecar -- covers both the binary-operand and call-argument shapes; verified deterministic across 5 repeated `TestRuntime` runs. Also hand-verified against the ticket's original `Waitlist`/`find_index`/`has` repro and against two independently-hoisted calls combined in one expression (`countdown(3) == countdown(5)`).

## 2026-09-27 (7)
- **Update**: [tasks.md](./tasks.md) -- Closed YZC-0115; bumped golden-test count to 112, next available to YZC-0116.
- **Update**: [tasks-done.md](./tasks-done.md) -- Added completion write-up for YZC-0115: `variantTypeArgs` only recognized a resolved `*sema.GenericInstType` for computing a qualified variant constructor's explicit Go type args; a zero-arg constructor (`Option.None()`) called from a boc generic over the same type param instead gets a bare, uninstantiated `*sema.StructType` from sema, which fell through to "no explicit args needed" -- and Go's own inference can't recover T for a zero-argument generic call from return-type context alone.
- **Update**: `internal/ir/lower.go` -- `variantTypeArgs` now also accepts a bare `*sema.StructType` (falling back to its `TypeParams` names) and, in both shapes, emits an unresolved `*sema.GenericType` argument by name instead of bailing out.
- **Creation**: `118_generic_variant_zero_arg_ctor` golden test + `.output` sidecar.
- **Note**: a second failure surfaced while stress-testing against a generic `Stack` struct (`r.value.ToStr undefined`) turned out to be the pre-existing empty-array-literal type-inference limitation, not a new bug -- no ticket filed.

## 2026-09-27 (6)
- **Update**: [tasks.md](./tasks.md) -- Closed YZC-0114; bumped golden-test count to 111, next available to YZC-0115.
- **Update**: [tasks-done.md](./tasks-done.md) -- Added completion write-up for YZC-0114: a generic struct's homoiconic `String()` calls `std.YzTypeName` on its type-param field, which fell back to `reflect.TypeOf(v).Name()` -- correct for an ordinary struct, but Go's reflect synthesizes a generic instantiation's `Name()` from its type arguments' fully package-qualified `String()` form, leaking e.g. `yz/runtime/rt.Int` for a std collection type argument.
- **Update**: `runtime/rt/types.go` -- `YzTypeName` now strips any `path/to/pkg.`-shaped prefix from `t.Name()` via a regex (`importPathPrefix`).
- **Creation**: `117_generic_homoiconic_type_arg` golden test + `.output` sidecar.

## 2026-09-27 (5)
- **Update**: [tasks.md](./tasks.md) -- Closed YZC-0113; bumped golden-test count to 110, next available to YZC-0114.
- **Update**: [tasks-done.md](./tasks-done.md) -- Added completion write-up for YZC-0113: `Dict.each` is documented in spec §10.8 but had neither a sema `fieldType` case (fell through to the same "extensible" `Unknown` YZC-0104 already fixed for `at`/`has`/`set`/`length`) nor a runtime `Each` method at all, failing only at `go build`.
- **Update**: `internal/sema/analyzer.go` -- Added `each` to `DictType`'s `fieldType` switch.
- **Update**: `runtime/rt/collections.go` -- Added `Dict[K, V].Each`, iterating in ascending `StringifyRepr` key order (matching `Dict.String`'s existing sort) instead of Go's randomized map order, so output is reproducible run to run.
- **Creation**: `116_dict_each` golden test + `.output` sidecar; verified deterministic across 5 repeated `TestRuntime` runs.
- **Note**: `remove`/`keys`/`values` (spec §10.8's other missing Dict methods) deliberately left unimplemented -- `remove`'s mutate-vs-copy-on-write semantics need a decision this fix didn't make (see write-up).

## 2026-09-27 (4)
- **Update**: [tasks.md](./tasks.md) -- Closed YZC-0111; filed YZC-0112 (found immediately after, via spec's own "Match with continue" example); bumped golden-test count to 109, next available to YZC-0113.
- **Update**: [tasks-done.md](./tasks-done.md) -- Added completion write-up for YZC-0111: `parseConditionalBoc` parsed a match arm's `{ multi; statement }` body block as one opaque `*ast.BocLiteral` element instead of flattening it, so it got lowered as an anonymous closure value and emitted as a dead, uncalled `func(){}` literal.
- **Update**: [tasks-detail.md](./tasks-detail.md) -- Added a Bugs entry for YZC-0112: `continue` inside a match arm has zero lowering/codegen support (only parsed + no-op'd in sema) and is silently dropped instead of falling through to the next branch, per spec §7.3. Marked *design* -- Go's if/else-if chain has no native construct for conditional fallthrough, so the fix needs a different generated shape for match, not just wiring up an IR node. Distinct from open YZC-0019 (loop break/continue).
- **Update**: `internal/parser/parser.go` -- `parseConditionalBoc` now detects a `{` immediately after the optional `cond =>`, parses it as a nested boc literal via `parseBocLiteral`, and flattens its `.Elements` directly into `arm.Body` instead of leaving it as a single wrapping node.
- **Creation**: `115_match_arm_body_block` golden test + `.output` sidecar.

## 2026-09-27 (3)
- **Update**: [tasks.md](./tasks.md) -- Closed YZC-0110; bumped error-conformance-test count to 26, next available to YZC-0111.
- **Update**: [tasks-done.md](./tasks-done.md) -- Added completion write-up for YZC-0110: `analyzeStructBoc`'s `case *ast.Ident` (meant only for bare single-letter generic type params) matched on Go AST type instead of `TokType`, so any bare-Ident statement inside a closure/anonymous-boc body -- HOF closures, match arms -- was silently registered as a fabricated generic type instead of being resolved and validated.
- **Update**: `internal/sema/analyzer.go` -- Gated the generic-type-param registration on `e.TokType == token.GENERIC_IDENT`; every other bare Ident element now falls through to ordinary expression analysis so undefined/malformed identifiers are flagged with a proper `undefined: %s` error.
- **Creation**: `test/conformance/testdata/errors/29_undefined_ident_in_closure.yz` + `.error` -- reproduces identically through the golden-test driver's unwrapped `compile()` path (no file-wrap needed), so a plain error-conformance fixture suffices.

## 2026-09-27 (2)
- **Update**: [tasks.md](./tasks.md) -- Closed YZC-0109; next available bumped to YZC-0110.
- **Update**: [tasks-done.md](./tasks-done.md) -- Added completion write-up for YZC-0109: golden tests for YZC-0057/YZC-0077 (recursive/mutually-recursive types) call `sema.AnalyzeFile` directly and never exercise the real CLI's mandatory per-file boc-wrap, which buries top-level type declarations one level below where `AnalyzeFile`'s forward-reference stub pre-registration pass looks.
- **Update**: `internal/sema/analyzer.go` -- Added `preRegisterNestedTypes` (mirrors `AnalyzeFile`'s first pass for types nested inside any boc body); `analyzeStructBoc`'s stub-reuse lookup changed from `fileScope.LookupLocal` to a `currentScope.Lookup` chain walk so it finds stubs at either level.
- **Creation**: `examples/recursive_type_forward_ref/` -- regression example (plain self-reference, mutual recursion, variant-constructor self-reference) proven to fail pre-fix and pass post-fix through the actual CLI; golden fixtures cannot exercise this path (same reasoning as `examples/sibling_calls`/`examples/bare_variant_root`).

## 2026-09-27
- **Update**: [tasks.md](./tasks.md) -- Closed YZC-0108; bumped golden-test count to 108, next available to YZC-0109.
- **Update**: [tasks-done.md](./tasks-done.md) -- Added completion write-up for YZC-0108: a variant match's default arm (no `=>`) aborted discriminant-switch lowering for the *entire* match, falling back to a boolean-condition path that treated variant names as undefined globals.
- **Update**: `internal/ir/ir.go` -- Added `SwitchCase.IsDefault`.
- **Update**: `internal/ir/lower.go` -- `tryLowerDiscriminantMatch`/`tryLowerDiscriminantMatchExpr` now handle a default arm (`arm.Condition == nil`) as an `IsDefault` case instead of aborting discriminant-match detection.
- **Update**: `internal/codegen/codegen.go` -- `emitSwitchStmt`/`emitSwitchIIFE` emit Go's `default:` for an `IsDefault` case.
- **Creation**: `114_variant_match_default` golden test + `.output` sidecar (covers both expression- and statement-position match).

## 2026-09-26 (7)
- **Update**: [tasks.md](./tasks.md) -- Filed YZC-0107 (found while continuing to dogfood struct-instance bocs beyond `examples/_wip/library`); next available bumped to YZC-0108.
- **Update**: [tasks-detail.md](./tasks-detail.md) -- Added a Bugs section entry for YZC-0107: a struct method that calls a self-recursive (non-leaf) sibling and forces its result in the same expression deadlocks, since neither YZC-0008's sync-rewrite nor its deferred-`BocGroup.Wait()` mitigation cover an expression-position force. Marked *design* -- needs a caller-side split-before/after-Schedule mechanism, not yet root-caused to a fix.
- **Verification**: minimal repro (`Counter` with `countdown`/`check`) and a realistic repro (`Waitlist` with `find_index`/`has`) both reproduce the exact same `fatal error: all goroutines are asleep - deadlock!`; kept in the session scratchpad, not committed, since there's no fix yet to regression-test.

## 2026-09-26 (6)
- **Update**: [tasks.md](./tasks.md) -- Closed YZC-0106; bumped golden-test count to 107.
- **Update**: [tasks-detail.md](./tasks-detail.md) -- Removed the closed YZC-0106 writeup (Bugs section now empty).
- **Update**: [tasks-done.md](./tasks-done.md) -- Added completion write-up for YZC-0106: the real root cause was `bodyHasBocCallsInStmtPos` false-flagging a last-element value-returning conditional's branch calls as needing a `BocGroup`, not the `lowerBocDeclAsSingleton` special case the ticket originally suspected.
- **Update**: `internal/ir/lower.go` -- `bodyHasBocCallsInStmtPos` now takes `resultType` and skips a last-element conditional's branches when the result isn't `std.Unit`; the trailing `bgVar` Wait append in `lowerBocBody` now also checks the body hasn't already terminated in a `ReturnStmt`.
- **Creation**: `113_self_recursive_return` golden test + `.output` sidecar.
- **Verification**: `examples/_wip/library` now builds and runs end-to-end (previously blocked by this bug).

## 2026-09-26 (5)
- **Update**: [tasks.md](./tasks.md) -- Closed YZC-0103; noted `examples/bare_variant_root` alongside `examples/sibling_calls` as example-only regression coverage.
- **Update**: [tasks-detail.md](./tasks-detail.md) -- Removed the closed YZC-0103 writeup.
- **Update**: [tasks-done.md](./tasks-done.md) -- Added completion write-up for YZC-0103 (parser `inTypeBoc` never set for bare-root-file variant constructors; confirmed a trailing comma was not the missing ingredient before diagnosing the real cause).
- **Creation**: `compiler/examples/bare_variant_root/` -- Regression example for YZC-0103 (golden tests bypass the root-file auto-wrap path, same reasoning as YZC-0102/`sibling_calls`).
- **Update**: `compiler/examples/_wip/library/BorrowResult.yz` -- Reverted to its natural bare-constructor form now that YZC-0103 is fixed.

## 2026-09-26 (4)
- **Update**: [tasks.md](./tasks.md) -- Closed YZC-0105; filed YZC-0106 (surfaced while fixing it); bumped golden-test count to 106.
- **Update**: [tasks-detail.md](./tasks-detail.md) -- Removed the closed YZC-0105 writeup; added YZC-0106.
- **Update**: [tasks-done.md](./tasks-done.md) -- Added completion write-up for YZC-0105 (conditional-as-return-value in `lowerBocBody`/`lowerConditionalExpr`), including the regression caught and fixed before landing.
- **Creation**: `112_conditional_return_value` golden test + `.output` sidecar.

## 2026-09-26 (3)
- **Update**: [tasks.md](./tasks.md) -- Closed YZC-0104; bumped golden-test count to 105.
- **Update**: [tasks-detail.md](./tasks-detail.md) -- Removed the closed YZC-0104 writeup.
- **Update**: [tasks-done.md](./tasks-done.md) -- Added completion write-up for YZC-0104 (missing `DictType` case in `Analyzer.fieldType`).
- **Creation**: `111_dict_methods` golden test + `.output` sidecar.

## 2026-09-26 (2)
- **Creation**: [tasks.md](./tasks.md) -- Filed YZC-0102 through YZC-0105, found while dogfooding a larger example (`examples/_wip/library/`); closed YZC-0102.
- **Update**: [tasks-detail.md](./tasks-detail.md) -- Added a Bugs section with details for YZC-0103/0104/0105; removed the closed YZC-0102 writeup.
- **Update**: [tasks-done.md](./tasks-done.md) -- Added completion write-up for YZC-0102 (sibling-call resolution priority in `lowerCall`).
- **Creation**: `compiler/examples/sibling_calls/` -- Regression example for YZC-0102 (golden tests bypass the root-file auto-wrap path, so this bug needed an `examples/`-level test instead).

## 2026-09-26
- **Update**: [tasks.md](./tasks.md) -- Closed YZC-0101; bumped golden-test count to 104.
- **Update**: [tasks-detail.md](./tasks-detail.md) -- Removed the closed YZC-0101 writeup (now-empty Bugs section removed).
- **Update**: [tasks-done.md](./tasks-done.md) -- Added completion write-up for YZC-0101 (sibling-method pre-scan in `analyzeStructBoc`).

## 2026-09-24
- **Update**: [tasks.md](./tasks.md) -- Closed YZC-0100 and YZC-0008; added and indexed YZC-0101.
- **Update**: [tasks-detail.md](./tasks-detail.md) -- Removed the closed YZC-0100/YZC-0008 writeups; added YZC-0101.
- **Update**: [tasks-done.md](./tasks-done.md) -- Added completion write-ups for YZC-0100 and YZC-0008.
- **Update**: [concurrency-design.md](./concurrency-design.md) -- Added a re-entrancy addendum: sync-rewriting a held-cown call is only safe for leaf/non-recursive callees.
- **Creation**: [doc-lifecycle.md](./doc-lifecycle.md) -- Rules for plan-doc lifecycle, `index.md`, and `log.md` in this OKF bundle.
- **Update**: [index.md](./index.md) -- Indexed `doc-lifecycle.md`.

## 2026-09-23
- **Creation**: [./conformance-golden-tests.md](./conformance-golden-tests.md) -- Document the `.output` sidecar requirement for `TestRuntime` to catch broken generated Go, using YZC-0100 as the example.

## 2026-09-16
- **Creation**: Added OKF `index.md` and `log.md` for bundle structure.

