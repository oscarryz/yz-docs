#impl
Ticket numbers are permanent. `[x]` = closed, `[ ]` = open. Next available: **YZC-0090**.

# Yz Compiler Implementation

## Status
- **92 golden + 25 error conformance tests passing** — `go test -race ./...` passes (test 51 has pre-existing timing flakiness)
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

YZC-0076 -- Existential associated types: opaque-token / path-identity tracking -- L -- *design* -- needs YZC-0079 -- *may not be needed: see detail*  
YZC-0016 -- String `++` concatenation -- S -- needs YZC-0031
YZC-0013 -- Array `<<` append -- S -- needs YZC-0031  
YZC-0009 -- Range iteration -- S -- needs YZC-0031  
YZC-0019 -- `break`/`continue`/`return` in loops -- M -- needs YZC-0031  
YZC-0014 -- Option/Result method chaining -- M -- needs YZC-0031  
YZC-0039 -- Operators audit -- L -- needs YZC-0031  
YZC-0059 -- Macro interface interaction -- *design* -- needs YZC-0025  
YZC-0008 -- Same-cown reentrant scheduling deadlock -- M -- dormant  
YZC-0089 -- Invariant 5: foo.yz + foo/ coexistence — loader merge + nested singleton codegen -- M -- needs YZC-0021  
YZC-0022 -- Multiple source roots -- M -- needs YZC-0085  
YZC-0044 -- Producer-consumer example and golden test -- M -- needs YZC-0031  
YZC-0023 -- Cancellation / non-local return -- L  
YZC-0058 -- Native type annotation -- L -- needs YZC-0025, YZC-0059  
YZC-0060 -- Design and implement `self` in Yz -- L -- needs YZC-0058, YZC-0059  
YZC-0041 -- Dependency management -- L  
YZC-0042 -- Package management (`yz` tool ) -- L -- needs YZC-0041  
YZC-0024 -- `return`, `break`, `continue` (major) -- L -- needs YZC-0019, YZC-0023  
YZC-0088 -- Codegen: attach compiled annotation boc to declaration metadata -- M -- needs YZC-0028  
YZC-0028 -- Macros (`Macro` interface) -- XL -- needs YZC-0025, YZC-0026, YZC-0027, YZC-0030, YZC-0066, YZC-0059   
YZC-0031 -- Scalar Types in Yz Source (uppering) -- XL -- needs YZC-0025, YZC-0028, YZC-0002 
YZC-0080 -- Uniform boc literal typing: one structural type derived from elements -- XL -- *design* -- needs YZC-0025

---

Details: [open](tasks-detail.md) · [done](tasks-done.md)

