---
type: impl
generated: { by: "oscarryz", at: 2026-09-23T00:00:00+00:00 }
---
#impl
# Conformance Golden Tests — the `.output` Sidecar

## The trap

A golden conformance test (`compiler/test/conformance/testdata/golden/NNN_name.{yz,go}`) only proves the compiler produces the *expected Go text*. It does not prove that text compiles, or that running it does the right thing.

`TestGolden` (`compiler/test/conformance/conformance_test.go`) runs every `.yz` file through `compile()` — parse → sema → lower → codegen — and diffs the result against the checked-in `.go` file. `compile()` stops at codegen. It never calls `go build`, so a golden pair can lock in generated Go that doesn't even compile, and `TestGolden` will pass forever: it's diffing text against itself, not against reality.

`TestRuntime` (`compiler/test/conformance/runtime_test.go`) is the only suite that actually `go build`s and `go run`s the generated code and checks stdout. But it only considers golden entries that have a matching `.output` sidecar file:

```go
if e.IsDir() || !strings.HasSuffix(e.Name(), ".output") {
    continue
}
```

A `.yz`/`.go` pair with no `.output` is **invisible** to `TestRuntime`. Add a new golden test without its `.output` file, and both `make test` (which skips `TestRuntime` entirely, via `-short`) and `make test-full` will report green — even if the generated code fails to compile or silently drops a side effect.

## Concrete case: YZC-0100

The regression test for YZC-0100 (`107_boc_typed_param.yz`) was originally added with only `.yz`/`.go`, no `.output`. It exercised a boc-typed field but never actually *called* the callback it declared, so it couldn't have caught the bug even if it had run. Worse, the sema fix in that ticket had a second, independent bug: a boc-typed parameter's invocation was being wrapped in `std.Go(...)` (async, fire-and-forget) instead of emitted as a direct call, so the callback's side effect (a `print`) was silently dropped. `TestGolden` was green throughout, because the generated Go — wrong as it was — matched the checked-in golden file byte for byte.

The bug only surfaced by writing `107_boc_typed_param.output`, letting `TestRuntime` build and run the binary, and noticing the expected second line of output never printed.

## Rule

When adding a golden test for any construct that has runtime-observable behavior (anything that prints, or whose evaluation order matters), create all three files:

1. `NNN_name.yz` — the source.
2. `NNN_name.go` — via `UPDATE_GOLDEN=1 go test ./test/conformance/...` (or by hand, then verified with `-run`).
3. `NNN_name.output` — the exact expected stdout, captured by actually running the compiled binary (`yzc run <dir>`, or building the generated Go directly). Don't hand-write the expected output from what you *think* the code should print — run it.

A golden test with no `.output` is a text-diff fixture only. Treat it as a **weaker** test than one with an `.output` sidecar, not an equivalent one — it does not confirm the generated code compiles, let alone behaves correctly.

## The opposite failure mode

`.output` sidecars aren't free either: they only work for deterministic output. YZC-0048 (see `tasks-done.md`) had to delete a `.output` sidecar for a test whose concurrent output ordering was legitimately non-deterministic — the source-diff `TestGolden` test was kept, but `TestRuntime` coverage for that case was dropped rather than made flaky. Before adding a `.output` file, make sure the construct under test has a single correct output ordering; if it doesn't, that's a real coverage gap to accept, not a reason to write a flaky expectation.
