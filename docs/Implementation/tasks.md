---
type: impl
generated: { by: "oscarryz", at: 2026-05-27T22:35:53+02:00 }
---
#impl
Ticket numbers are permanent. `[x]` = closed, `[ ]` = open. Next available: **YZC-0130**.

# Yz Compiler Implementation

## Status
- **116 golden + 26 error conformance tests passing** (+ macro driver suite: debug_merge + 4 error cases; multi_root + subdir_coexist + macro_debug examples) — `go test -race ./...` passes (test 51 has pre-existing timing flakiness)
- Plus `examples/sibling_calls` (YZC-0102), `examples/bare_variant_root` (YZC-0103), `examples/recursive_type_forward_ref` (YZC-0109), and `examples/file_wrapper_sig_form` (YZC-0122): all have zero golden-test coverage since golden fixtures bypass the real CLI's mandatory root-file-wrap path entirely; these are `examples/`-level regression tests instead
- Compiler: `compiler/` directory, Go module `module yz`
- Runtime: `compiler/runtime/rt/`, macro wire codec: `compiler/runtime/macrowire/`

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

~~YZC-0076 -- Existential associated types -- closed: not needed under current macro dispatch model~~  
~~YZC-0102 -- Sibling calls inside any auto-wrapped root file lower to a nonexistent global `Name.Call()` instead of `self.name()`~~
~~YZC-0103 -- Uppercase root file containing only bare variant constructors (no plain fields) is not recognized as a type decl~~
~~YZC-0104 -- Dict dot-call methods (`.at`, `.has`, ...) had no sema type-inference case at all and widened to `any`~~
~~YZC-0105 -- `?:` conditional as a boc's final return value was discarded; codegen emitted an if/else statement and fell through to `return Unit`~~
~~YZC-0106 -- Self-recursive `#(...)`-declared boc with non-Unit return type and a value-returning conditional/expression as its last statement produces a Schedule closure hardcoded to `func() std.Unit`, breaking the build~~
~~YZC-0107 -- Sibling call to a non-leaf held-cown callee deadlocks when its result is forced in the same method~~
~~YZC-0108 -- Variant match with a trailing default arm (no `=>`) misdetected as a boolean-condition match, undefined variant-name globals at build~~
~~YZC-0109 -- Recursive/mutually-recursive type declarations (YZC-0057/0077) resolve in golden tests but fail with "undefined type" through the real CLI, which always wraps a file's top-level statements one level deeper than AnalyzeFile's forward-reference pre-registration pass scans~~  
~~YZC-0110 -- Bare identifier inside any closure/anonymous-boc body (HOF closures, match arms) is misdetected as a generic type param declaration, silently skipping undefined-symbol checking~~  
~~YZC-0111 -- Cond-match arm with a multi-statement body block (`cond => { s1; s2 }`) parses as one opaque boc-literal element instead of flattening, emitting an uncalled closure literal~~  
YZC-0082 -- `Foo.Bar()` static access is not rejected by yzc -- S -- [ticket](tickets/YZC-0082.md)  
YZC-0112 -- `continue` in a cond-match arm does not fall through -- M -- *design* -- [ticket](tickets/YZC-0112.md)  
YZC-0116 -- A non-leaf held-cown call inside a match arm or closure body deadlocks -- M -- *design* -- [ticket](tickets/YZC-0116.md)  
YZC-0118 -- Macro run-output cache serves stale output after a macro edit -- S -- [ticket](tickets/YZC-0118.md)  
YZC-0121 -- A call to an unresolvable identifier fails only at `go build` -- S -- [ticket](tickets/YZC-0121.md)  
~~YZC-0122 -- Same-named file-wrapper boc in signature form fails at `go build`~~ -- [ticket](tickets/YZC-0122.md)  
YZC-0123 -- A local declared without an initializer is dropped -- *design* -- [ticket](tickets/YZC-0123.md)  
YZC-0124 -- Boc body is not atomic across inner calls -- L -- *design* -- [ticket](tickets/YZC-0124.md)  
YZC-0125 -- Sibling boc calls run sequentially -- S -- [ticket](tickets/YZC-0125.md)  
YZC-0126 -- Constructor call at the start of a statement fails to parse in a type's methods -- S -- [ticket](tickets/YZC-0126.md)  
YZC-0127 -- Boc-typed fields of a type are dropped -- S -- [ticket](tickets/YZC-0127.md)  
YZC-0128 -- A type's methods are not visible on its instances inside its own methods -- S -- [ticket](tickets/YZC-0128.md)  
YZC-0129 -- A singleton and a type whose names differ only in case fail to build -- M -- [ticket](tickets/YZC-0129.md)  
~~YZC-0119 -- `unquoteString` didn't resolve `\`` (or `\0`), and chained ReplaceAll passes let `\\n` be reinterpreted as `\n` by a later pass~~
~~YZC-0120 -- `yzc run` never wired the compiled binary's stdin, so the `read` builtin hit instant EOF and returned "" instead of waiting for input~~
~~YZC-0117 -- Array.Each / Dict.Each are Go-void but sema types them as Unit-returning, so `.each(...)` as a boc method's tail statement generates `return <void call>`, failing at `go build`~~
~~YZC-0113 -- `Dict.each` documented in spec §10.8 but had no sema case and no runtime method at all, failing only at `go build`~~  
~~YZC-0114 -- Generic struct's homoiconic (backtick) repr leaks the full Go import path for a nested generic type argument (e.g. `Box(Array[yz/runtime/rt.Int], ...)`)~~  
~~YZC-0115 -- Zero-argument generic variant constructor (e.g. `Option.None()`) called from another generic context fails to compile with "cannot infer T"~~
YZC-0016 -- String `++` concatenation -- S -- needs YZC-0031 -- [ticket](tickets/YZC-0016.md)  
YZC-0013 -- Array `<<` append -- S -- needs YZC-0031 -- [ticket](tickets/YZC-0013.md)  
YZC-0009 -- Range iteration -- S -- needs YZC-0031 -- [ticket](tickets/YZC-0009.md)  
YZC-0019 -- `break`/`continue`/`return` in loops -- M -- needs YZC-0031 -- [ticket](tickets/YZC-0019.md)  
YZC-0014 -- Option/Result method chaining -- M -- needs YZC-0031 -- [ticket](tickets/YZC-0014.md)  
YZC-0039 -- Operators audit -- L -- needs YZC-0031 -- [ticket](tickets/YZC-0039.md)  
~~YZC-0101 -- Sibling method call fails sema resolution when callee is declared after caller~~
~~YZC-0100 -- Boc-typed field in a body-only singleton silently dropped as a param~~
~~YZC-0008 -- Same-cown reentrant scheduling deadlock~~
~~YZC-0091 -- Nested singleton codegen: sub-singleton struct with own methods~~
YZC-0044 -- Producer-consumer example and golden test -- M -- needs YZC-0031 -- [ticket](tickets/YZC-0044.md)  
YZC-0023 -- Cancellation / non-local return -- L -- [ticket](tickets/YZC-0023.md)  
YZC-0058 -- GoSource: Go-backed type implementations -- L -- needs ~~YZC-0025~~, ~~YZC-0059~~ -- [ticket](tickets/YZC-0058.md)  
YZC-0060 -- Implement `self` as a lexically-resolved compiler built-in -- L -- needs ~~YZC-0059~~ (YZC-0058 only for the Go-backed-method sub-case) -- [ticket](tickets/YZC-0060.md)  
YZC-0099 -- `self` / native-builtins: confirm go_source calling convention covers built-in operators -- S -- needs YZC-0058, YZC-0060, YZC-0031 -- [ticket](tickets/YZC-0099.md)  
~~YZC-0041 -- `Deps` macro: compile-time dependency validation -- cancelled, superseded by YZC-0097~~
YZC-0096 -- `yz fetch`: dependency fetcher -- M -- needs ~~YZC-0097~~, ~~YZC-0022~~ -- [ticket](tickets/YZC-0096.md)  
YZC-0042 -- `yz` tool: run, new, add, init (wraps yzc + yz fetch) -- L -- needs ~~YZC-0041~~, YZC-0096, ~~YZC-0097~~ -- [ticket](tickets/YZC-0042.md)  
YZC-0024 -- `return`, `break`, `continue` (major) -- L -- needs YZC-0019, YZC-0023 -- [ticket](tickets/YZC-0024.md)  
YZC-0088 -- Codegen: attach compiled annotation boc to declaration metadata -- M -- needs ~~YZC-0028~~ -- [ticket](tickets/YZC-0088.md)  
~~YZC-0098 -- Self-scope associated type resolution + structural bound codegen~~
~~YZC-0028 -- Macros (`Macro` interface) -- first slice complete; follow-ups listed in tasks-done~~
YZC-0031 -- Scalar Types in Yz Source (uppering) -- XL -- needs ~~YZC-0025~~, ~~YZC-0028~~, ~~YZC-0002~~, ~~YZC-0022~~ -- [ticket](tickets/YZC-0031.md)  

---

Details: [open](tickets/index.md) · [done](tasks-done.md) · [legacy](tasks-detail.md)

