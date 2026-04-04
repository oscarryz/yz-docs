# Yz Compiler Implementation

## Phase 0 — Project Setup
- [x] Create `compiler/` directory skeleton
- [x] Initialize `go.mod` (`module yz`)
- [x] Create `cmd/yzc/main.go` (CLI with `build`, `run`, `new`)
- [x] Create `Makefile` (`build`, `test`, `clean`)
- [x] Create `compiler/README.md`
- [x] Verify: `go build ./...` passes

## Phase 1 — Lexer
- [x] `internal/token/token.go` — token types
- [x] `internal/lexer/lexer.go` — tokenizer + ASI
- [x] `internal/lexer/lexer_test.go` — 38 tests, all passing

## Phase 2 — Parser
- [x] `internal/ast/ast.go` — AST node types
- [x] `internal/parser/parser.go` — recursive descent
- [x] `internal/parser/parser_test.go` — 32 tests, all passing

## Phase 3 — Semantic Analysis
- [x] `internal/sema/analyzer.go` — scope, type inference, boc/struct dispatch
- [x] `internal/sema/analyzer_test.go` — tests passing

## Phase 4 — IR
- [x] `internal/ir/ir.go` — IR node type definitions
- [x] `internal/ir/lower.go` — AST+sema → IR lowerer
- [x] `internal/ir/ir_test.go` — 8 tests, all passing

## Phase 5 — Code Generation
- [x] `internal/codegen/codegen.go` — Go source emitter
- [x] `internal/codegen/codegen_test.go` — 10 tests, all passing
- [x] `cmd/yzc/build.go` — full pipeline: parse→sema→IR→codegen→go build
- [x] `cmd/yzc/new.go` — project scaffolding

## Phase 6 — Runtime Library
- [x] `runtime/yzrt/types.go` — Int, Decimal, String, Bool, Unit with symbol-named methods
- [x] `runtime/yzrt/thunk.go` — Thunk[T], Go[T] (goroutine spawn), Force()
- [x] `runtime/yzrt/collections.go` — Array[T], Dict[K,V], Range
- [x] `runtime/yzrt/core.go` — Print, While, Info, BocGroup (structured concurrency)
- [x] `runtime/yzrt/yzrt_test.go` — tests passing

## Phase 7 — Integration & Testing
- [x] `compiler/test/conformance/` — golden tests, 18 passing
- [x] `compiler/examples/` — counter, milestone (concurrent fetch + counter)
- [ ] Error tests — programs that should fail with specific errors

## Language Features — Implemented
- [x] `while` loop
- [x] `BocWithSig` — top-level functions and methods inside singleton/struct bocs
- [x] `match` expression (condition form)
- [x] `mix` statement — Go embedding
- [x] Multi-file projects — flat and subdirectory (cross-package FQN)
- [x] Type-only BocWithSig — `Name #(params)`: data params → struct (no constructor); all-boc params → Go interface (structural typing)
- [x] `http` built-in singleton — `http.get(uri)`, `http.post(uri, body)`
- [x] First milestone — concurrent HTTP fetch + counter (`examples/milestone/`)

## Language Features — Not Yet Implemented
- [ ] Mixed type-only decl — `Name #(name String, greet #())` both data + BocType params (deferred)
- [ ] `BocWithSig` body-only form — `name #(params) = { body }` (body re-declares params)
- [ ] Variant/discriminant match — `match expr { Variant.Case => body }`
- [ ] Cross-package singleton method calls — `pkg.singletonVar.method()`
- [ ] Error reporting — useful sema errors with line/col
- [ ] `yzc run` — compile + execute in one step
