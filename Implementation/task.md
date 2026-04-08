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
- [x] Error tests — 7 cases: parse errors, undefined variable/type, mix undefined/conflict/not-struct

## Language Features — Implemented
- [x] `while` loop
- [x] `BocWithSig` — top-level functions and methods inside singleton/struct bocs
- [x] `match` expression (condition form)
- [x] `mix` statement — Go embedding
- [x] Multi-file projects — flat and subdirectory (cross-package FQN)
- [x] Type-only BocWithSig — `Name #(params)`: data params → struct (no constructor); all-boc params → Go interface (structural typing)
- [x] `http` built-in singleton — `http.get(uri)`, `http.post(uri, body)`
- [x] First milestone — concurrent HTTP fetch + counter (`examples/milestone/`)

## Language Features — Implemented (continued)
- [x] Variant/discriminant sum types — `Pet: { Cat(...), Dog(...) }` with per-variant constructors
- [x] Discriminant match — `match expr { Cat => body }, { Dog => body }` → Go switch
- [x] Cross-package singleton method calls — `pkg.singleton.method()`
- [x] `yzc run` — compile + execute in one step
- [x] `http` built-in singleton — `http.get(uri)`, `http.post(uri, body)`
- [x] thunk transparency — `a: boc.call()` auto-forced on use

## Language Features — Not Yet Implemented
- [x] Mixed type-only decl — `Name #(name String, greet #())` → struct with data fields + function-typed fields + method wrappers
- [x] `BocWithSig` body-only form — `name #(params) = { body }` — named and anonymous param matching
- [x] Error reporting — Rust-style diagnostics with source context and caret underlines

## BocWithSig Body-Only — Deferred
- [x] Default values in params — `#(name String = "hello")` — injected at call sites; golden test 21
- [x] `ShortDecl` as param — `name : "default"` in sig — type inferred from default; golden test 22
- [x] Generic variant types — `Option: { V; Some(value V); None() }` with `[V any]` on struct and constructors; discriminant match works; golden test 25
- [x] Generic type vars in sig — `identity #(value V, V)` → `func identity[V any](value V) *Thunk[V]`; golden test 26
- [x] Uninstantiated generics — `Option(String)` → `*Option[std.String]` in type positions
- [x] Declare-only then assign-later — `greet #(name String)` then `greet = { name String; … }` → FuncDecl; golden test 23

## Language Features — Already Implemented (discovered)
- [x] Multiline strings — strings span lines naturally; `"` or `'` closes on any line (lexer handles `\n` inside string literals)

## Language Features — Implemented (continued)
- [x] HOF / closures as arguments — `list.filter({ item Int; item > 10 })` — sync closures with typed params; `Array.Filter`, `Array.Each`, `ArrayMap`; golden test 27

## Generics — Future Work
- [x] HOF: `list.map({ item Int; item * 2 })` — `lowerCall` detects `.map(boc)` on ArrayType → emits `std.ArrayMap(recv, closure)`; result type inferred via `:=`; golden test 28
- [x] Generic constraint inference (Option 4) — sema scans method bodies for T-method calls; records constraints; checks all at instantiation; reports all missing methods at once; error test 09; golden test 31 (generic method receiver)
- [x] Go constraint generation — emit `[T interface{ ToStr() std.String }]` from inferred constraints; lowerMethodName fixes to_string→ToStr; golden test 32
- [x] Multiple type params — `#(key K, value V)` → `[K any, V any]`; Pair[K,V] struct + makePair[K,V] function; parser fix: TYPE_IDENT'(' only parsed as VariantDef in type boc bodies (inTypeBoc flag); golden test 33
- [x] Generic structs (non-variant) — `Box: { T; value T }` → `type Box[T any] struct { value T }`; golden test 29
- [x] Typed generic declaration — `b Box(String) = Box("hello")` → `var b *Box[std.String] = NewBox(...)`; golden test 30; TypedDecl in lowerMainStmt
- [x] Optional parens for trailing-block calls — `list.filter { block }` without `()`; in `parsePostfix`, LBRACE after MemberExpr → CallExpr with BocLiteral arg; golden test 34
- [x] Unary minus on variables — `-x` → `x.Neg()`; `a - -b` → `a.Minus(b.Neg())`; pipeline was already wired (parser+sema+lowerer+codegen); golden test 35

## Known Bugs
- [x] Dict literals — fixed: now emits `std.NewDict[K,V]().Set(k,v)...` chain; golden test 24
- [x] Array literals — already worked via variadic `std.NewArray(...)`; golden test 24
