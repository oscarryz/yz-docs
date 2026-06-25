#readme #feature 

# Yz Language Features

## Core

- [Bocs.md](Bocs.md) — The core abstraction: blocks of code that hold state and behavior
- [Boc Interface.md](Boc%20Interface.md) — The `#(...)` notation — boc interface, inputs/outputs, encapsulation
- [Variables.md](Variables.md) — Variable declaration, initialization, and their role as parameters
- [Define new types.md](Define%20new%20types.md) — Uppercase bocs as named types (structs)
- [Create instances.md](Create%20instances.md) — Constructing instances of types
- [Structural typing.md](Structural%20typing.md) — Structural (duck-typing) compatibility between bocs
- [Type Alias.md](Type%20Alias.md) — Naming an existing type with a new name
- [Type variants.md](Type%20variants.md) — Sum types and variant constructors
- [Associated Types.md](Associated%20Types.md) — Associated types and the `#()` metatype
- [Generics - Type Parameters.md](Generics%20-%20Type%20Parameters.md) — Generic type parameters and constraints
- [Path Dependent Types.md](Path%20Dependent%20Types.md) — Types resolved through a value's path

## Control flow

- [Conditional Bocs.md](Conditional%20Bocs.md) — Boolean conditionals, `match`, and pattern matching
- [return, break, continue.md](return%2C%20break%2C%20continue.md) — Control flow keywords
- [Error handling.md](Error%20handling.md) — Result and Option types for error handling

## Concurrency

- [Concurrency.md](Concurrency.md) — Async by default, lazy evaluation, and structured concurrency

## Built-in types

- [Int.md](Int.md) — Integer type
- [Decimal.md](Decimal.md) — Decimal (floating-point) type
- [Strings.md](Strings.md) — The String type and its methods
- [String interpolation.md](String%20interpolation.md) — Interpolation inside string literals
- [Array.md](Array.md) — Arrays: literals, indexed access, filter/each/map, and append
- [Associative arrays.md](Associative%20arrays.md) — Dictionaries (key-value maps)

## Syntax

- [Non-Word invocation.md](Non-Word%20invocation.md) — Operator-style method calls (`+`, `<<`, `==`, etc.)
- [Trailing block syntax.md](Trailing%20block%20syntax.md) — Omitting parentheses when a boc is the sole argument
- [Comments.md](Comments.md) — Code comments
- [Reserved words and characters and symbols.md](Reserved%20words%20and%20characters%20and%20symbols.md) — Reserved identifiers and symbols

## Annotations and metaprogramming

- [Annotations.md](Annotations.md) — Structured metadata attached to declarations
- [Macros.md](Macros.md) — Compile-time code generation via the `Macro` interface
- [Structural Reflection.md](Structural%20Reflection.md) — The `Boc` metatype API for introspecting declarations

## Interop and extensions

- [Go Extensions.md](GoExtensions.md) — Go-backed type implementations via `go_source:`

## Project and modules

- [Code organization.md](Code%20organization.md) — Files, packages, source roots, and the module system
- [Dependencies.md](Dependencies.md) — Declaring external Yz packages in `project.info`

## Design overview

- [Yz Language Design.md](Yz%20Language%20Design.md) — High-level design goals and principles
