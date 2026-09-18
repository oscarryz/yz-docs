# Features — core language capabilities

Documentation outlining the syntax, semantics, and design principles of the Yz programming language's core components: Bocs (Blocks of Code), types, generics, arrays, and control flow.

## Directories

- [Macros](./Macros/index.md) — Compile-time code generation documentation and macro-driven build automation details.
- [Pending features](<./Pending features/index.md>) — Documented language extensions that are proposed but not yet implemented.
- [Replaced features](<./Replaced features/index.md>) — Archived descriptions of older capabilities superseded by newer versions.

## Files

- [Annotations.md](./Annotations.md) — Design documentation regarding custom source code markers and metadata.
- [Array.md](./Array.md) — Rules for array initialization, declaration, and bounds handling in Yz.
- [Associated Types.md](<./Associated Types.md>) — Design documentation regarding path-dependent type definitions in generics.
- [Associative arrays.md](<./Associative arrays.md>) — Dictionary-like syntax for storing and retrieving key-value pairs.
- [Boc Interface.md](<./Boc Interface.md>) — Defining signatures and behavioral contracts for code blocks using `#(...)` .
- [Bocs.md](./Bocs.md) — Core documentation regarding the Block of Code (BOC) data structure and execution model.
- [Code organization.md](<./Code organization.md>) — Design documentation regarding how BOCs are split, named, and grouped across files.
- [Comments.md](./Comments.md) — Syntax for single-line `//` and block `/* */` comments in Yz source code.
- [Concurrency.md](./Concurrency.md) — Overview of asynchronous execution models, channel mechanics, and actor isolation.
- [Conditional Bocs.md](<./Conditional Bocs.md>) — Usage patterns for conditional logic within code blocks.
- [Create instances.md](<./Create instances.md>) — Design documentation regarding object instantiation and constructor syntax.
- [Decimal.md](./Decimal.md) — Documentation regarding precision handling and decimal arithmetic rules in Yz.
- [Define new types.md](./Define new types.md) — Syntax for declaring custom object structures and data types.
- [Dependencies.md](./Dependencies.md) — Design documentation regarding project dependency management and linking.
- [Error handling.md](<./Error handling.md>) — Design documentation regarding `Result` types and error propagation mechanisms.
- [Generics - Type Parameters.md](<./Generics - Type Parameters.md>) — Constraints and definitions for generic type parameters in Yz.
- [GoExtensions.md](./GoExtensions.md) — Documentation regarding Go-backed native type implementations and bindings.
- [Int.md](./Int.md) — Syntax rules and behavior for integer numbers and arithmetic in Yz.
- [Macros.md](./Macros.md) — Documentation regarding the compile-time macro expansion system.
- [Non-Word invocation.md](<./Non-Word invocation.md>) — Design documentation regarding method calls on non-alphanumeric identifiers.
- [Path Dependent Types.md](<./Path Dependent Types.md>) — Handling types that rely on internal block states and context values.
- [README.md](./README.md) — Yz Language Features overview and summary.
- [Reserved words and characters and symbols.md](<./Reserved words and characters and symbols.md>) — Symbols and keywords strictly prohibited or reserved in the Yz grammar.
- [String interpolation.md](<./String interpolation.md>) — Design documentation regarding embedding variables within string literals using `${}` .
- [Strings.md](./Strings.md) — Rules for string literal quoting, escaping, and text manipulation in Yz.
- [Structural Reflection.md](./Structural Reflection.md) — Inspecting types at compile-time based on their internal structure fields.
- [Structural typing.md](<./Structural typing.md>) — Defining type compatibility based on field layout rather than inheritance names.
- [SumTypes.md](./SumTypes.md) — Implementation of union/sum types and pattern matching in Yz.
- [Trailing block syntax.md](<./Trailing block syntax.md>) — Rules for omitting parentheses when passing block literals as arguments.
- [Type Alias.md](<./Type Alias.md>) — Syntax rules for creating shorthand names for complex type structures.
- [Type variants.md](<./Type variants.md>) — Design documentation regarding distinct data structures within a sum type.
- [Variables.md](./Variables.md) — Rules for declaring, typing, and scoping variables within Yz blocks.
- [Yz Language Design.md](<./Yz Language Design.md>) — Core design philosophy and architecture principles of the Yz language.
- [return, break, continue.md](<./return, break, continue.md>) — Control flow operators for exiting loops and functions in Yz.
