#impl 
# Yz Compiler Implementation

## Status
- **53 golden + 2 error conformance tests passing** — `go test -race ./...` passes
- Compiler: `compiler/` directory, Go module `module yz`
- Runtime: `compiler/runtime/rt/`

---

## Completed Phases

All foundational phases are done. Details are in git history.

| Phase | Description | Tests |
|-------|-------------|-------|
| 0 | Project setup — `cmd/yzc`, `Makefile`, `go.mod` | — |
| 1 | Lexer — tokenizer + ASI | 38 |
| 2 | Parser — recursive descent AST | 32 |
| 3 | Semantic analysis — scope, type inference, boc/struct dispatch | passing |
| 4 | IR — lowerer (AST+sema → IR) | 8 |
| 5 | Codegen — Go source emitter; `yzc build`/`run`/`new` | 10 |
| 6 | Runtime — `types.go`, `core.go`, `collections.go`, `cown.go` | passing |
| 7 | Integration — conformance golden tests, examples, error tests | 51 golden |

---

## Implemented Features

### Language
- Singleton bocs, struct bocs, main boc — all uniform (boc uniformity passes 1–4)
- Boc declarations as methods; boc expanded form with named/anonymous param matching
- Type-only boc declarations: data params → struct; all-boc params → Go interface
- Mixed type-only decl: `Name #(name String, greet #())` → struct + method wrappers
- Variant/discriminant sum types: `Pet: { Cat(...), Dog(...) }` with per-variant constructors
- Discriminant match: `match expr { Cat => body }` → Go switch
- Condition match in statement position (if/else) and expression position (IIFE)
- `while` — user-land recursion via boc declaration; `tryLowerWhile` and `yzrt.While` removed
- HOF / closures as arguments: `.filter`, `.each`, `.map` on Array
- Default values in params: `#(name String = "hello")`
- `ShortDecl` as param: `name : "default"` — type inferred from default
- Declare-only then assign-later: `greet #(name String)` then `greet = { ... }`
- Optional parens for trailing-block calls: `list.filter { block }`
- Unary minus: `-x` → `x.Neg()`
- Multiline strings
- String interpolation: `${}` (backtick reserved for infostrings)
- Error reporting: Rust-style diagnostics with source context and caret underlines

### Types & Generics
- All types as `std.*` structs; literal boxing in codegen
- Generic structs: `Box: { T; value T }` → `Box[T any]`; generic variant types: `Option: { V; Some(value V); None() }`
- Generic type vars in boc declarations: `identity #(value V, V)` → `func identity[V any]`
- Generic constraint inference: sema infers from usage; reports all violations at once
- Go constraint generation: emits `[T interface{ Method() }]` from inferred constraints
- Multiple type params: `#(key K, value V)` → `[K any, V any]`
- Typed generic declaration: `b Box(String) = Box("hello")`; uninstantiated generics in type positions

### Concurrency (BOC — all phases complete)
- A: mutex cowns — data-race freedom
- B.1: queue-based cown scheduler — lock-free, spawn-order guarantee
- B.2: `ScheduleMulti` — atomic multi-cown acquisition
- C: ownership-based field writes (SWMR); cross-cown writes via `Schedule`
- D: struct boc instances embed `std.Cown`; fresh instance per call site for multi-cown boc declarations
- E.1: implicit BocGroup per scope; split-BocGroup pattern; `ScheduleAsSuccessor`
- E.3: plain scalar types (no lazy fields); `GoStore[T]`/`GoWait`; `*Thunk[T]` internal to runtime

### Runtime / Built-ins
- `http` singleton: `http.get(uri)`, `http.post(uri, body)`
- `print`, `Info`, `BocGroup` structured concurrency
- `Array[T]`, `Dict[K,V]`, `Range` with HOF: `.filter`, `.each`, `.map`
- `yzc run` — compile + execute in one step
- Cross-package singleton method calls
- `examples/milestone/` — concurrent HTTP fetch + counter boc (first milestone)

---

## Open Work

Ticket numbers: `YZC-NNNN`. Numbers are permanent — closed tickets keep their number.

### Bugs

- [x] **[YZC-0001] Variants broken** — variants were not updated for the BOC model; see `examples/variants`
- [ ] **[YZC-0002] Cross-package broken** — broke during BOC migration
- [x] **[YZC-0003] Assigning Unit-returning boc to variable** — `a : foo()` where `foo` returns Unit should be a sema error (analogue to Go's `x := f()` where `f` returns nothing); detect in sema; add error golden test
- [ ] **[YZC-0004] Top-level boc callable as function** — `foo: { time.sleep(1); "done" }` lowers as singleton struct, not callable as `foo()`; needs sema + lowerer fix
- [~] **[YZC-0005] Double return with sleep** — `foo: { time.sleep(1); 1 }` emits two return statements in generated Go — *not reproducible as of BOC work; superseded by YZC-0035*
- [ ] **[YZC-0006] Standalone boc invocation** — `p : { print("hello") }; p()` requires `p.call()` workaround; blocked on YZC-0004
- [ ] **[YZC-0007] Unused variables in generated Go** — Yz allows unused vars; Go rejects them; fix: after lowering a scope, append `_ = varName` for any declared name never referenced in subsequent IR nodes
- [ ] **[YZC-0008] Reentrant inline calls unsafe in HOF closures** — closure emitted inside a `ScheduleMulti` body and passed as argument to another boc contains sync-body calls that bypass cown acquisition; fix: sub-generator with `heldCowns = nil` when emitting closure args; dormant until HOF closures operate on cown-bearing types
- [x] **[YZC-0035] Sema does not check boc body return type against declared output** — when a boc declares a non-Unit output type (e.g. `foo #(Int)`) but the body's last expression returns Unit (e.g. only `time.sleep` or `print` calls), sema accepts it silently; the lowerer then emits `return std.TheUnit` which fails at `go build` with a type error; affects any void-returning call in that position, not just sleep; fix: after inferring the body's return type, verify it matches the declared output type and report a sema error

### Language Features

- [ ] **[YZC-0034] Definite assignment analysis** — `name Type` (uninitialized typed declaration) must be assigned before first use; `Bar()` with unassigned required fields is a compile error unless all paths assign before read; crossing a boc boundary requires fully initialized values at the call site. Sema pass: build per-scope control-flow graph; track assigned set; report "field f used before initialization" on unassigned reads. Spec: §3.2. Depends on: YZC-0033.

- [ ] **[YZC-0009] Range iteration** — `1.to(10).each({ i Int; ... })` — lowerer recognizes `.each` on Array only; extend to Range receiver
- [ ] **[YZC-0010] HOF iteration + cown happens-before** — resolved by YZC-0036: HOF methods use `a→b→a` indirect recursion → ScheduleAsSuccessor at each step → sequential processing behaviour. No further design needed; implement as sequential closure calls.
- [ ] **[YZC-0036] While loop yield and external caller interleaving** — **design resolved**: recursive call (callee FQN = caller FQN) → tail enqueue; non-recursive → ScheduleAsSuccessor. Gives `while` (direct recursion) producer semantics; HOF/indirect recursion gets sequential processor semantics. Indirect cycles (`a→b→a`) remain a known limitation. See `docs/Questions/While loop yield and external caller interleaving.md`. Implementation: lowerer must compare callee FQN to enclosing boc FQN before emitting Schedule vs ScheduleAsSuccessor.
- [ ] **[YZC-0011] Named arguments in constructor calls** — `Person(name: "Alice", age: 30)`
- [ ] **[YZC-0012] Multiple return values** — `x, y = swap(x, y)` — multi-assign LHS not in any golden test
- [ ] **[YZC-0013] Array append via `<<`** — `a << item` → `a.Append(item)`; `Array.Append` exists in yzrt
- [ ] **[YZC-0014] Option/Result method chaining** — `result.or_else({ error Error; ... })`, `result.and_then({ val T; ... })`
- [ ] **[YZC-0015] Non-word boc names** — `balance+= #(amount Int) { ... }` — parser only allows word identifiers in boc declarations; fix: accept `NON_WORD` token; map to Go-safe name via symbol table; add golden test
- [ ] **[YZC-0016] String concatenation with `++`** — lowerer emits `Plusplus` but runtime `String` has no such method; fix: add `Plusplus` to `String` in yzrt
- [ ] **[YZC-0017] Dict optional access** — `d[key]` should return `Option(V)`; currently panics on missing key via `At()`
- [x] **[YZC-0018] Bool methods `&&` / `||`** — `Bool.Ampamp` / `Bool.Pipepipe` exist in yzrt; golden test 53 confirms end-to-end. *Note: current operators are eager sync calls, special-cased on built-in Bool; when Bool moves to Yz source (YZC-0031), `&&`/`||` become lazy closure-taking boc methods that go through the normal BOC cycle — see YZC-0031 sub-item.*
- [ ] **[YZC-0019] `break` / `continue` / `return` in loops** — blocked on concurrency model settling; lowerer should emit compile error when encountered rather than silently dropping
- [ ] **[YZC-0020] `to_str()` mapping on user types** — sema rejects `p.to_str()` on user structs ("no field to_str"); needs sema to expose `to_str` on all struct types and lowerer to emit `ToStr()` or a default Go `String()` fallback

### Infrastructure

- [ ] **[YZC-0033] Compiler deep review against settled spec** — the compiler was generated using the old spec which had stateless/stateful boc distinction, `Unit` as user-facing type, and boc declarations treated as Go functions (no cown, no persistent fields). Review and align: (1) boc declarations should lower to singleton structs with cowns, not plain Go functions; (2) `foo.param` should be accessible after call; (3) sema errors should say "returns nothing" not mention Unit; (4) all bocs serialized through cown — no "fully parallel" fast path for boc declarations. Depends on: YZC-0004, YZC-0006.
  - [x] spec/02 grammar updated: labeled=input/unlabeled=output rule, BocDecl three forms, MixStmt removed
  - [x] sema errors say "returns nothing" instead of "Unit" (`displayType` helper, YZC-0003 check)
- [ ] **[YZC-0021] Directory and file bocs** — defer until in-file nesting works; extend FQN tree to directories and files as bocs
- [x] **[YZC-0032] Rename `BocWithSig` in compiler code** — AST node `BocWithSig`, sema path `analyzeBocWithSig`, lowerer path `lowerBocWithSig`, and all related identifiers should be renamed to `BocDecl` / `analyzeBocDecl` / `lowerBocDecl` to match the settled terminology; also rename the `BocWithSig` → `BocDecl` grammar production in spec/02
- [ ] **[YZC-0022] Multiple source roots** — `src/` + `lib/` as independent FQN mount points; compiler accepts list of source roots; builds one FQN forest per root
- [ ] **[YZC-0023] Cancellation / non-local return** — non-local `return` across goroutine boundaries conflicts with structured concurrency; see `docs/Questions/How to cancel a running block.md`

---

## Major Features

### YZC-0024 — `return`, `break`, `continue`

Blocked on concurrency model settling (see YZC-0019 and YZC-0023).

- [ ] Parser — `BreakStmt` / `ContinueStmt` AST nodes (tokens already exist)
- [ ] Sema — validate context: `break`/`continue` only inside loop; `return` tracks nearest named boc
- [ ] Lowerer — emit compile error when encountered (fail loudly)
- [ ] Spec 07 — update control-flow spec
- [ ] Golden tests — sema-level error tests

### YZC-0025 — Infostrings: content is a boc body

Infostring delimiter stays backtick; content is full Yz syntax, parsed and type-checked, never executed.

- [ ] AST — `InfoString` holds `*BocLiteral` instead of `*StringLit`
- [ ] Lexer — re-lex infostring content as Yz source
- [ ] Parser — re-parse as boc body using existing boc-body parser
- [ ] Sema — type-check content; validate referenced names
- [ ] Codegen — attach compiled infostring boc to declaration metadata
- [ ] Spec 01 — update

### YZC-0026 — Generics: Explicit Constraint Declaration

`thing T Talker` declares `T` must implement `Talker`; additive with inference.

- [ ] Parser — `T Constraint` optional suffix after single-uppercase-letter type param
- [ ] Sema — validate at instantiation; union with inferred constraints
- [ ] Error messages — explicit vs inferred violations distinct
- [ ] Spec 04 — update

### YZC-0027 — `:` as Type Alias

`Name : SomeType` declares a type alias usable anywhere.

- [ ] Feature doc — `docs/Features/Type Alias.md`
- [ ] Parser — distinguish `Name : TypeExpr` (alias) from `Name TypeExpr` (typed decl) and `name : value` (short decl)
- [ ] Sema — register alias; resolve as aliased type; no runtime fields
- [ ] Lowerer — emit `type Name = GoType`
- [ ] Spec 04 — add

### YZC-0028 — Compile-Time Bocs (`Compile` interface)

Any boc with `Schema #()` and `run #(Boc, Boc)` satisfies `Compile`. Depends on: YZC-0025, YZC-0026, YZC-0027, YZC-0030.

- [ ] Sema — recognize `Compile` structural interface (duck-typed)
- [ ] Sema — scan infostring for `compile_time: [...]`; schedule during type inference
- [ ] Boc metatype — `Boc` value type for `run`: `{name String, fields [Boc], methods [Boc], ...}`
- [ ] Two-phase build — compile `Compile` implementations first; call via subprocess during main compilation
- [ ] Serialization — `Boc` wire format (JSON or binary) for subprocess calls
- [ ] AST merge — merge returned `Boc` into parent boc's AST
- [ ] Cycle detection — circular `compile_time` triggers → compile error
- [ ] Caching — keyed on source hash + input boc structure hash
- [ ] Spec 12 — new spec file

### YZC-0029 — Remove `mix`: runtime + spec — PARTIALLY COMPLETE

Compiler removal done. Remaining work depends on YZC-0028.

- [x] Lexer — removed `token.MIX`
- [x] Parser — removed `MixStmt`; `mix` is now a regular identifier
- [x] Sema — removed mix analysis (embedding resolution, conflict detection)
- [x] Lowering/Codegen — removed Go-embedding path
- [x] Golden tests — updated / removed mix-using conformance tests
- [ ] Runtime — implement `Mix` as a `Compile` boc in yzrt or stdlib
- [ ] Spec 09 — remove `mix`; document `Mix` compile implementation

### YZC-0030 — Associated Types: Path-Dependent Type References

`process(g Graph, n g.Node)` — no new syntax; sema resolves `g.Node` at the call site by looking up `Node` on the concrete type bound to `g`. See decisions 50–51 in `decisions.md`.

- [ ] Sema — `value.TypeName` in type position; resolve against concrete type of `value`
- [ ] Lowerer — emit concrete Go type at resolution site
- [ ] Golden test — `associated_types.yz`

### YZC-0031 — Scalar Types in Yz Source (uppering)

Prerequisite: E.3 complete (done). `Int/String/Bool/Decimal/Unit` move from Go to `stdlib/` with `compile-time:[Native]` annotation. Native ops annotated per method; higher-level methods (`times`, `to`, `clamp`, `>=`, `Ord`) in plain Yz. Depends on: YZC-0025, YZC-0028.

- [ ] Define `compile-time:[Native]` infostring semantics (depends on YZC-0025)
- [ ] Move scalar types to `stdlib/`
- [ ] Annotate native ops per method
- [ ] Implement higher-level methods in Yz
- [ ] Remove all primitive-type special-casing from the compiler
- [ ] `Bool.&&` / `Bool.||` — rewrite as lazy closure methods `#(other #(Bool), Bool)`; calls go through the normal BOC cycle (return `*Thunk[Bool]`, participate in BocGroup/GoWait) instead of the current eager sync `Ampamp`/`Pipepipe`; lowerer wraps bare expression operands in a closure: `a && b` → `a.Ampamp({ b })`

---

## Ticket Rules

- `YZC-NNNN` numbers are permanent and never reused; closed items keep their number
- Numbers are assigned in creation order; next available: **YZC-0037**
- `depends-on` is a flat reference to ticket numbers — no nested phase hierarchy
- Reference tickets in commit messages and code comments for easy grep: `// YZC-0008`
- When the open list in any section exceeds ~10 items, split into a `tickets/` directory with one file per ticket
