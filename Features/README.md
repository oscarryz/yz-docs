# Yz Language Features

## Feature Reference

- [Bocs.md](Bocs.md) — The core abstraction: blocks of code that hold state and behavior
- [Variables.md](Variables.md) — Variable declaration, initialization, and their role as parameters
- [Block type.md](Block%20type.md) — The `#(...)` type syntax for describing boc signatures
- [Define new types.md](Define%20new%20types.md) — Uppercase bocs as named types (structs)
- [Create instances.md](Create%20instances.md) — Constructing instances of types
- [Structural typing.md](Structural%20typing.md) — Structural (duck-typing) compatibility between bocs
- [Generics.md](Generics.md) — Generic type parameters
- [Type variants.md](Type%20variants.md) — Sum types and variant constructors
- [mix.md](mix.md) — Code composition: merge fields and methods from another boc
- [Concurrency.md](Concurrency.md) — Async by default, lazy evaluation, and structured concurrency
- [Conditional Bocs.md](Conditional%20Bocs.md) — Boolean conditionals, `match`, and pattern matching
- [Array.md](Array.md) — Arrays: literals, indexed access, filter/each/map, and append
- [Associative arrays.md](Associative%20arrays.md) — Dictionaries (key-value maps)
- [Strings.md](Strings.md) — The String type and its methods
- [String interpolation.md](String%20interpolation.md) — Backtick interpolation inside string literals
- [Int.md](Int.md) — Integer type
- [Decimal.md](Decimal.md) — Decimal (floating-point) type
- [Trailing block syntax.md](Trailing%20block%20syntax.md) — Omitting parentheses when a boc is the sole argument
- [Non-Word invocation.md](Non-Word%20invocation.md) — Operator-style method calls (`+`, `<<`, `==`, etc.)
- [return, break, continue.md](return%2C%20break%2C%20continue.md) — Control flow keywords
- [Error handling.md](Error%20handling.md) — Result and Option types for error handling
- [Code organization.md](Code%20organization.md) — Files, packages, and modules
- [Info strings.md](Info%20strings.md) — Documentation strings attached to declarations
- [Comments.md](Comments.md) — Code comments
- [Single Writer.md](Single%20Writer.md) — Concurrency safety: each object has one logical writer
