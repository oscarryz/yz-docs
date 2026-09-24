---
type: spec
generated: { by: "oscarryz", at: 2026-09-24T00:00:00+00:00 }
---
#spec
# Index — Yz Language Specification

The normative specification for the Yz programming language, organized as numbered chapters covering syntax, semantics, and the standard library.

## Chapters

- [1. Lexical Structure](./01-lexical-structure.md) — Source text encoding, character categories, and how source is tokenized.
- [2. Grammar](./02-grammar.ebnf.md) — The syntactic grammar of Yz in EBNF, built on the lexer's token stream.
- [3. Expressions and Statements](./03-expressions-and-statements.md) — Semantics of expressions vs. statements, complementing the grammar.
- [4. Type System](./04-type-system.md) — Yz's structural type system: type kinds, compatibility rules, variants, and generics.
- [5. Type Inference](./05-type-inference.md) — How the compiler infers types from initializers and call sites when not explicitly declared.
- [6. Blocks and Scoping](./06-blocks-and-scoping.md) — The boc (block of code) as Yz's unified abstraction, and variable scoping rules.
- [7. Control Flow](./07-control-flow.md) — Conditional and iteration constructs expressed via methods and `match`, without dedicated `if`/`for`/`while` keywords.
- [8. Concurrency](./08-concurrency.md) — Behaviour-oriented concurrency (BOC): async-by-default boc invocation, transparent thunks, and structured concurrency.
- [9. Modules and Code Organization](./09-modules-and-organization.md) — The file system as module system: how directory structure defines the namespace hierarchy.
- [10. Standard Library](./10-standard-library.md) — Built-in types (`Int`, `Decimal`, `String`, `Bool`, `Unit`), collection types, and core functions.
- [11. Conformance Tests](./11-conformance-tests.md) — The canonical test suite of Yz programs with expected behavior, serving as living specification and compiler validation.
- [12. Macros](./12-macros.md) — The `Macro` interface and compile-time macro execution (first-slice stub; full design in `docs/Features/Macros.md`).
