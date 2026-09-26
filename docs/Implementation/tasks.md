---
type: impl
generated: { by: "oscarryz", at: 2026-05-27T22:35:53+02:00 }
---
#impl
Ticket numbers are permanent. `[x]` = closed, `[ ]` = open. Next available: **YZC-0110**.

# Yz Compiler Implementation

## Status
- **108 golden + 25 error conformance tests passing** (+ macro driver suite: debug_merge + 4 error cases; multi_root + subdir_coexist + macro_debug examples) — `go test -race ./...` passes (test 51 has pre-existing timing flakiness)
- Plus `examples/sibling_calls` (YZC-0102), `examples/bare_variant_root` (YZC-0103), and `examples/recursive_type_forward_ref` (YZC-0109): all have zero golden-test coverage since golden fixtures bypass the real CLI's mandatory root-file-wrap path entirely; these are `examples/`-level regression tests instead
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
YZC-0107 -- Sibling call to a non-leaf held-cown callee deadlocks when its result is forced in the same method -- L -- *design*
~~YZC-0108 -- Variant match with a trailing default arm (no `=>`) misdetected as a boolean-condition match, undefined variant-name globals at build~~
~~YZC-0109 -- Recursive/mutually-recursive type declarations (YZC-0057/0077) resolve in golden tests but fail with "undefined type" through the real CLI, which always wraps a file's top-level statements one level deeper than AnalyzeFile's forward-reference pre-registration pass scans~~
YZC-0016 -- String `++` concatenation -- S -- needs YZC-0031  
YZC-0013 -- Array `<<` append -- S -- needs YZC-0031  
YZC-0009 -- Range iteration -- S -- needs YZC-0031  
YZC-0019 -- `break`/`continue`/`return` in loops -- M -- needs YZC-0031  
YZC-0014 -- Option/Result method chaining -- M -- needs YZC-0031  
YZC-0039 -- Operators audit -- L -- needs YZC-0031  
~~YZC-0101 -- Sibling method call fails sema resolution when callee is declared after caller~~
~~YZC-0100 -- Boc-typed field in a body-only singleton silently dropped as a param~~
~~YZC-0008 -- Same-cown reentrant scheduling deadlock~~
~~YZC-0091 -- Nested singleton codegen: sub-singleton struct with own methods~~
YZC-0044 -- Producer-consumer example and golden test -- M -- needs YZC-0031  
YZC-0023 -- Cancellation / non-local return -- L  
YZC-0058 -- GoSource: Go-backed type implementations -- L -- needs ~~YZC-0025~~, ~~YZC-0059~~  
YZC-0060 -- Implement `self` as a lexically-resolved compiler built-in -- L -- needs ~~YZC-0059~~ (YZC-0058 only for the Go-backed-method sub-case)  
YZC-0099 -- `self` / native-builtins: confirm go_source calling convention covers built-in operators -- S -- needs YZC-0058, YZC-0060, YZC-0031  
~~YZC-0041 -- `Deps` macro: compile-time dependency validation -- cancelled, superseded by YZC-0097~~
YZC-0096 -- `yz fetch`: dependency fetcher -- M -- needs ~~YZC-0097~~, ~~YZC-0022~~
YZC-0042 -- `yz` tool: run, new, add, init (wraps yzc + yz fetch) -- L -- needs ~~YZC-0041~~, YZC-0096, ~~YZC-0097~~  
YZC-0024 -- `return`, `break`, `continue` (major) -- L -- needs YZC-0019, YZC-0023  
YZC-0088 -- Codegen: attach compiled annotation boc to declaration metadata -- M -- needs ~~YZC-0028~~  
~~YZC-0098 -- Self-scope associated type resolution + structural bound codegen~~
~~YZC-0028 -- Macros (`Macro` interface) -- first slice complete; follow-ups listed in tasks-done~~
YZC-0031 -- Scalar Types in Yz Source (uppering) -- XL -- needs ~~YZC-0025~~, ~~YZC-0028~~, ~~YZC-0002~~, ~~YZC-0022~~ 

---

Details: [open](tasks-detail.md) · [done](tasks-done.md)

