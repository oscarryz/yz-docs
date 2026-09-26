# update log  

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

