# Yz Compiler — Implementation Plan

All pre-implementation decisions resolved in [decisions.md](file:///Users/oscar/code/github/oscarryz/yz-docs-1/Implementation/decisions.md).

---

## Architecture

```
Source (.yz) → Lexer → Parser → AST → Sema → IR → Codegen → Go Source → go build → Binary
```

```mermaid
flowchart LR
    A[".yz files"] --> B["Lexer"]
    B --> C["Parser"]
    C --> D["AST"]
    D --> E["Sema"]
    E --> F["IR"]
    F --> G["Codegen"]
    G --> H[".go files"]
    H --> I["go build"]
    I --> J["Binary"]
```

## Project Structure

The compiler lives inside the existing `yz-docs-1` repository under the `compiler/` directory.

```
yz-docs-1/                    (existing repo)
├── README.md
├── spec/
├── Examples/
├── ...
└── compiler/                 ← NEW: all compiler code here
    ├── cmd/yzc/              CLI entry point
    │   └── main.go
    ├── internal/
    │   ├── lexer/            Tokenizer + ASI
    │   ├── token/            Token types
    │   ├── ast/              AST node definitions
    │   ├── parser/           Recursive descent parser
    │   ├── sema/             Semantic analysis (types, scopes)
    │   ├── ir/               Intermediate representation
    │   ├── codegen/          Go source code emitter
    │   └── build/            go build orchestration
    ├── runtime/yzrt/         Runtime library (imported by generated Go)
    │   ├── actor.go          Goroutine + channel actor infra
    │   ├── thunk.go          Lazy thunk + materialization
    │   ├── types.go          Built-in types (Int, Decimal, String, Bool, Unit)
    │   ├── collections.go    Array, Dictionary, Range
    │   ├── variants.go       Option, Result
    │   └── core.go           print, while, info
    ├── test/                 Conformance & integration tests
    ├── examples/             Example .yz programs
    ├── go.mod                module yz
    ├── Makefile
    └── README.md
```

---

## Phase Details

### Phase 0 — Project Setup ✅

| Item | Details |
|------|---------|
| Location | `compiler/` directory in this repository |
| Go module | `module yz` (local-only, not go-getable) |
| CLI | `cmd/yzc/main.go` with `build`, `run`, `new` subcommands |
| Makefile | `build`, `test`, `clean` targets |

> [!IMPORTANT]
> Internal imports use the `yz/internal/...` path (e.g. `import "yz/internal/lexer"`).
> The runtime package is imported by generated code as `yz/runtime/yzrt`.

---

### Phase 1 — Lexer ✅

#### `internal/token/token.go`
Token enum and `Token` struct (`Type`, `Literal`, `Line`, `Col`).

Token categories:
- Identifiers: `IDENT` (lowercase), `TYPE_IDENT` (uppercase multi-char), `GENERIC_IDENT` (single uppercase)
- Keywords: `BREAK`, `CONTINUE`, `RETURN`, `MATCH`, `MIX`
- Literals: `INT_LIT`, `DECIMAL_LIT`, `STRING_LIT`
- Non-word: `NON_WORD` (open char set)
- Delimiters: `LBRACE`, `RBRACE`, `LPAREN`, `RPAREN`, `LBRACKET`, `RBRACKET`, `COLON`, `ASSIGN`, `COMMA`, `SEMICOLON`, `DOT`, `HASH`, `FAT_ARROW`

#### `internal/lexer/lexer.go`
- UTF-8 rune scanning
- Multi-line strings (both `'...'` and `"..."`)
- String interpolation: track backtick nesting inside strings
- Escape sequences: `\n \t \r \\ \' \" \` \0`
- Comments: `//` to EOL, `/* ... */` nested
- ASI: insert `SEMICOLON` after newline when prev token is identifier, literal, `break`/`continue`/`return`, `)`, `]`, `}`

---

### Phase 2 — Parser ✅

#### `internal/ast/ast.go`
Node types for all constructs. Key nodes:

| Category | Nodes |
|----------|-------|
| Statements | `ShortDecl`, `TypedDecl`, `Assignment`, `ReturnStmt`, `BreakStmt`, `ContinueStmt`, `MixStmt`, `BocWithSig` |
| Expressions | `BinaryExpr` (non-word L-to-R), `UnaryExpr`, `CallExpr`, `MemberExpr`, `IndexExpr`, `Ident`, `IntLit`, `DecimalLit`, `StringLit`, `BocLiteral`, `ArrayLiteral`, `DictLiteral`, `MatchExpr`, `GroupExpr`, `InfoString` |
| Types | `BocTypeExpr`, `ArrayTypeExpr`, `DictTypeExpr`, `SimpleTypeExpr` |
| Top-level | `SourceFile`, `VariantDef`, `BocParam` |

#### `internal/parser/parser.go`
Recursive descent. Key grammar points:
- **Expression**: `UnaryExpr { NON_WORD UnaryExpr }` — flat, left-to-right
- **Multiple assignment**: `IdentifierList "=" ExpressionList` — only lowercase `IDENT` on LHS
- **Boc literal**: `"{" { BocElement sep } "}"` — sep is `;` (ASI) or `,`
- **BocWithSig**: `name #(params) { body }` — params available in body without redeclaration (sema concern). `name #(params) = { body }` — body must redeclare params.
- **Match**: two forms (with/without subject expression)
- **Commas as separators**: commas act as statement separators in boc bodies alongside semicolons (supports `T, E,` generic params and match arm lists)
- **`=>` non-word rule**: standalone `=>` is always FAT_ARROW; `!=>`, `<=>` etc. are single NON_WORD tokens

---

### Phase 3 — Semantic Analysis

#### [NEW] `internal/sema/`
- **Scope**: linked list of scope frames, lexical lookup with shadowing
- **Type checker**: structural compatibility, width subtyping
- **Inference**: variable types from RHS, return types from last expression, generic params from usage
- **BocWithSig param scoping**: when the `name #(params) { body }` form is used (no `=`), sema adds the signature params into the body's scope so they are available without redeclaration
- **Variants**: track discriminant tags per variant constructor
- **Mix**: flatten fields into host, error on conflicts
- **Access**: `#()` hides everything not in the signature
- **FQN resolution**: resolve each boc's fully-qualified name from source path + nesting; detect collisions across source roots

---

### Phase 4 — IR

#### [NEW] `internal/ir/`
Bridge between AST and Go codegen. Maps Yz concepts to Go-friendly constructs:

| Yz Concept | IR Representation |
|-----------|-------------------|
| Lowercase boc | Go struct (singleton state) + goroutine-wrapped method calls |
| Uppercase boc (type) | Go struct + constructor func; each invocation creates new instance |
| Boc with methods | Go struct with methods |
| Non-word method `a + b` | Method call using symbol name: `a.plus(b)`. `?` → `qm`, `==` → `eqeq`, `!=` → `neq`, `&&` → `ampamp`, `\|\|` → `pipepipe`, `<=` → `lteq`, etc. |
| Any boc invocation | Goroutine; result wrapped in `yzrt.Thunk[T]` (lazy, materializes on first use) |
| Structured concurrency | Parent boc has a `sync.WaitGroup`; each child goroutine registers; parent Done() waits |
| `match` variant | Go `switch` on discriminant tag |
| `match` condition | Go `if/else` chain |
| Literal `1`, `"x"` | Boxed: `std.NewInt(1)`, `std.NewString("x")` — boxing done in codegen |
| Closure captures | Go closure or explicit capture struct |

---

### Phase 5 — Code Generation

#### [NEW] `internal/codegen/`
Emit `.go` files:
- One Go package per Yz namespace (directory)
- Generated `main.go` for the `main` boc (entry point)
- Struct definitions for types; all fields typed as `std.*` (e.g. `std.Int`, `std.String`)
- Method definitions with non-word names mapped to symbol names (`plus`, `qm`, `eqeq`, etc.)
- All boc calls wrapped in goroutines; results are `yzrt.Thunk[T]`
- Structured concurrency via `sync.WaitGroup` in each boc
- Literal boxing at call sites: `1` → `std.NewInt(1)`
- `go.mod` with dependency on `yz/runtime/yzrt`

#### [NEW] `internal/build/`
- Write generated `.go` files to `target/gen/` inside the project being compiled
- Run `go mod init` + `go mod tidy`
- Run `go build -o target/bin/app`
- `target/` should be added to `.gitignore` in `yzc new`-generated projects

---

### Phase 6 — Runtime Library

#### [NEW] `runtime/yzrt/`
Go package imported by generated code:

| File | Contents |
|------|----------|
| `actor.go` | Actor struct, message queue (buffered chan), sequential processing loop, structured concurrency via `sync.WaitGroup` |
| `thunk.go` | `Thunk[T]` generic type, lazy eval, materialization on first use |
| `types.go` | `Int`, `Decimal`, `String`, `Bool`, `Unit` with all spec methods. Non-word method names use symbol naming: `plus`, `minus`, `qm`, `eqeq`, etc. Exported as `std.Int` etc. from Yz perspective. |
| `collections.go` | `Array[T]`, `Dict[K,V]`, `Range` |
| `variants.go` | `Option[T]`, `Result[T,E]`, discriminant tags |
| `core.go` | `Print()` (materializes thunks), `While()`, `Info()` |

---

### Phase 7 — Integration & Testing

- **Conformance tests**: Created **incrementally** as each phase lands (not upfront). `.yz` + `.expected` pairs in `compiler/test/conformance/`. Run via `go test`.
- **Golden tests**: `.yz` → expected `.go` output comparison
- **End-to-end**: compile + run + compare stdout to `.expected`
- **Error tests**: verify compile errors for invalid programs
- **CLI**: `yzc build`, `yzc run`, `yzc new <name>`
- **Project config**: minimal `project.yz` format

---

## Proposed Build Order

1. **Phase 0** → get the skeleton compiling ✅
2. **Phase 1** → lex a simple program, verify with tests ✅
3. **Phase 2** → parse to AST, pretty-print it
4. **Phase 6** (partial) → runtime types needed for codegen
5. **Phase 4 + 5** → IR + codegen for a minimal subset (variable decl, print, arithmetic)
6. **Phase 3** → add type checking incrementally
7. **Iterate 4-6** → add features one by one (bocs, types, match, actors, thunks...)
8. **Phase 7** → conformance tests as each feature lands

> [!TIP]
> First end-to-end milestone: a concurrent program — fetching two resources concurrently and a counter boc. Entry point is a boc named `main`.

## Verification Plan

### Automated Tests

All tests run from the `compiler/` directory:

```bash
cd /Users/oscar/code/github/oscarryz/yz-docs-1/compiler
go test ./...
```

Each phase adds tests in its own package (e.g. `internal/lexer/lexer_test.go`, `internal/parser/parser_test.go`).

### Manual Verification

After Phase 5, compile and run a concurrent `.yz` program end-to-end. The first milestone
program fetches two resources concurrently and implements a counter boc:

```bash
cd /Users/oscar/code/github/oscarryz/yz-docs-1/compiler
go run ./cmd/yzc run ../examples/concurrent.yz
```

Generated output goes to `target/gen/` inside the compiled project directory.
Binary is built to `target/bin/app`.
