# Task: Update All Examples to Current Yz Spec

## Context

The `Examples/` directory contains files collected from various languages and early Yz drafts.
Many of these predate the current spec (`spec/01` through `spec/11`) and use outdated syntax.

## Goal

Review every `.md` file in `Examples/` that contains Yz code and update it to use correct,
current Yz syntax as defined in the spec.

## Reference

The authoritative spec files are:

- `spec/01-lexical-structure.md` — identifiers, keywords, ASI, non-word identifiers, strings
- `spec/02-grammar.ebnf` — full grammar
- `spec/03-expressions-and-statements.md` — semantics
- `spec/04-type-system.md`
- `spec/05-type-inference.md`
- `spec/06-blocks-and-scoping.md`
- `spec/07-control-flow.md`
- `spec/08-concurrency.md`
- `spec/09-modules-and-organization.md`
- `spec/10-standard-library.md`

## Key Syntax Rules to Check

- **Declarations**: use `:` for short decl (`name: "Alice"`), not `=` or `let`/`var`
- **Assignment**: `=` is the only operator and is a statement, not an expression
- **Bocs**: blocks of code use `{ }`. Lowercase = singleton boc. Uppercase = type.
- **Parameters**: uninitialized typed declarations inside a boc are parameters: `age Int`
- **Boc with signature**: `greet #(name String, String) { "Hello, `name`!" }`
- **Non-word methods**: `+`, `-`, `*`, `/`, `==`, `?`, etc. are methods, not operators.
  Equal precedence, left-to-right: `1 + 2 * 3 = 9`. Use parens for grouping.
- **Conditional**: `condition ? { true_branch }, { false_branch }` (method on Bool, bocs 
  separated by comma)
- **Match**: `match expr { Variant => ... }, { ... }` or `match { cond => ... }, { ... }` (not when)
- **No nil**: use `Option` and `Result` instead
- **Keywords**: only `break`, `continue`, `return`, `match`, `mix`
- **`true`/`false`**: constants, not keywords
- **Strings**: `"..."` or `'...'`, multi-line by default. Interpolation: `` `expr` `` inside string.
- **Info strings**: string literal immediately before a declaration attaches metadata
- **Mix**: `mix Name` flattens another boc's fields into current boc
- **No loop syntax**: use `range.each({...})`, `while({ cond }, { body })`
- **Concurrency**: all boc calls are non-blocking (goroutine + lazy thunk). Structured concurrency
  ensures a boc waits for all child bocs before it is considered complete.
- **Generics** : They are single letter identifier, not `<` `>` enclosed
- **Function calls**": Always use parenthesis. `foo bar` => `foo(bar)` 

## What to Do Per File

1. Read the file
2. If it contains Yz code (look for `.yz` code blocks, most of the time they are in markdown as 
   `js` or `javascript` but they would contain Yz-specific syntax), update it
3. If it is purely another language's example with no Yz translation, either:
   - Add a correct Yz translation section, OR
   - Mark it as "Reference only — not Yz" at the top
4. Fix any outdated Yz syntax per the rules above
5. Keep the original non-Yz code intact (it serves as reference inspiration)

## Files to Review

Run: `ls Examples/*.md` to get the full list (~80 files).

## Done When

Every `.md` file in `Examples/` either:
- Contains correct, current Yz syntax, OR
- Is clearly marked as a reference-only file from another language
