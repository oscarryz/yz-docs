# update log  

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

