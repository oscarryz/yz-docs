#impl
Ticket numbers are permanent. `[x]` = closed, `[ ]` = open. Next available: **YZC-0100**.

# Yz Compiler Implementation

## Status
- **103 golden + 25 error conformance tests passing** (+ macro driver suite: debug_merge + 4 error cases; multi_root + subdir_coexist + macro_debug examples) — `go test -race ./...` passes (test 51 has pre-existing timing flakiness)
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
YZC-0016 -- String `++` concatenation -- S -- needs YZC-0031
YZC-0013 -- Array `<<` append -- S -- needs YZC-0031  
YZC-0009 -- Range iteration -- S -- needs YZC-0031  
YZC-0019 -- `break`/`continue`/`return` in loops -- M -- needs YZC-0031  
YZC-0014 -- Option/Result method chaining -- M -- needs YZC-0031  
YZC-0039 -- Operators audit -- L -- needs YZC-0031  
YZC-0008 -- Same-cown reentrant scheduling deadlock -- M -- dormant  
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

