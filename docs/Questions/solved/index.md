# Solved — resolved design inquiries

Conclusions regarding core semantic decisions including array initialization, pattern matching, generics constraints, variable lifetime, and structural typing rules.

## Subdirectories

- [Type signature/](./Type%20signature/index.md)— Details on how method and block signatures are defined and validated.
- [concurrency](./concurrency/index.md) — Conclusions on async evaluation, channels, and actor models.
- [generics](./generics/index.md) — Resolutions for generic type enforcement and polymorphic functions.

## Files

- [Array initialization.md](<./Array initialization.md>) — Rules for creating and populating arrays at runtime.
- [Associated Types.md](<./Associated Types.md>) — Design documentation regarding associated types in generics.
- [Associative array initialization.md](<./Associative array initialization.md>) — Syntax for creating and initializing key-value mappings.
- [Baby working.md](<./Baby working.md>) — Early proof-of-concept implementation notes.
- [Can we include constraints in the block definition to validate qm.md](<./Can we include constraints in the block definition to validate qm.md>) — Exploring inline type constraints for block definitions.
- [Consider assigning expressions.md](<./Consider assigning expressions.md>) — Design documentation regarding assignment operators and expression evaluation.
- [Consider type + definition special case.md](<./Consider type + definition special case.md>) — Handling cases where a type and its definition are merged.
- [Creating instances.md](<./Creating instances.md>) — Design documentation regarding object instantiation.
- [Default Args.md](<./Default Args.md>) — Syntax and behavior for optional block parameters.
- [Default values.md](<./Default values.md>) — Implementation of default argument binding in blocks.
- [Do block still have to be objects and functions at the same time qm.md](<./Do block still have to be objects and functions at the same time qm.md>) — Separating BOC capabilities for data structures vs callable behavior.
- [Enum.md](./Enum.md) — Design documentation regarding enumeration syntax.
- [Enums  and Disjointed Unions.md](<./Enums  and Disjointed Unions.md>) — Implementation of discrete type variations.
- [Error Handling.md](<./Error Handling.md>) — Standardizing `Result` types and error propagation mechanisms.
- [Existential Types and Associated Types.md](<./Existential Types and Associated Types.md>) — How Yz handles existential types over associated types.
- [Explicit return or last variable?.md](<./Explicit return or last variable?.md>) — Deciding between an explicit `return` keyword vs the "last expression" pattern.
- [Extensibility, Enums, Union Types, and Specialization.md](<./Extensibility, Enums, Union Types, and Specialization.md>) — Language features related to type specialization.
- [Flow typing and conditionals.md](<./Flow typing and conditionals.md>) — Dynamic type narrowing based on `if`/`?` operators.
- [Generics - Bound depending on use.md](<./Generics - Bound depending on use.md>) — Inferred generic boundaries based on block usage.
- [Generics - momorphized qm.md](<./Generics - momorphized qm.md>) — Polymorphic implementations across different types.
- [Generics without <>.md](Generics%20without%20<>.md) — Exploring syntax alternatives that omit angle brackets for generics.
- [Generics.md](./Generics.md) — Core design for generic type parameters in Yz.
- [How much magic should we allow qm.md](<./How much magic should we allow qm.md>) — Defining the limits of type inference and macro execution.
- [How to add sum types.md](<./How to add sum types.md>) — Implementation of SumTypes for pattern matching.
- [How to create asymetric matching types.md](<./How to create asymetric matching types.md>) — Handling mismatched data structures during pattern matching.
- [How to initialize literals like `Int`.md](<./How to initialize literals like `Int`.md>) — Design documentation regarding literal syntax and type binding.
- [How to support multiple interfaces in generics.md](<./How to support multiple interfaces in generics.md>) — Constraining generics to satisfy multiple structural requirements.
- [How to use utility functions.md](<./How to use utility functions.md>) — Standard library integration and extension points.
- [Instance creation, parameters.md](<./Instance creation, parameters.md>) — Design documentation regarding constructor arguments and object mapping.
- [Macro Interface Interaction Design.md](<./Macro Interface Interaction Design.md>) — How BOCs interact with the compile-time macro system.
- [Macros.md](./Macros.md) — Design documentation regarding the macro expansion process.
- [Match files and folders to blocks.md](<./Match files and folders to blocks.md>) — Organization strategies for module loading.
- [Memory management.md](<./Memory management.md>) — Rules for stack allocation, garbage collection points, and BOC lifetimes.
- [Native Type Annotations.md](<./Native Type Annotations.md>) — Design documentation regarding built-in type modifiers.
- [Nomenclature.md](./Nomenclature.md) — Design documentation regarding naming conventions and reserved keywords.
- [Optional syntax.md](<./Optional syntax.md>) — Exploring optional punctuation in block declarations.
- [Pattern Matching.md](<./Pattern Matching.md>) — Core implementation of the `match` control flow structure.
- [Pattern matching - More dicussion.md](<./Pattern matching - More dicussion.md>) — Extended discussion on pattern matching edge cases.
- [Re-declare variables if there's type signature?.md](<./Re-declare variables if there's type signature?.md>) — Handling variable scope and shadowing within BOCs.
- [Reference vs value.md](<./Reference vs value.md>) — Design documentation regarding pass-by-reference versus pass-by-value semantics.
- [Reuse context Or embedding.md](<./Reuse context Or embedding.md>) — Deciding between composing blocks via embedding vs passing context.
- [Semicolons.md](./Semicolons.md) — Design documentation regarding block separation and ASI (Automatic Semicolon Insertion).
- [Semicolons?.md](./Semicolons?.md) — Discussion on optional line termination in single-expression blocks.
- [Should structural typing include variable name.md](<./Should structural typing include variable name.md>) — Determining if record names matter for type compatibility.
- [Signatures + Literals duplication.md](<./Signatures + Literals duplication.md>) — Reducing boilerplate between block definitions and type signatures.
- [Special case for arrays and dictionaries?.md](<./Special case for arrays and dictionaries?.md>) — Syntax optimizations for collection types.
- [Tuples?.md](./Tuples?.md) — Exploring inline multi-value return structures.
- [Uniform Boc Literal Typing.md](<./Uniform Boc Literal Typing.md>) — Design documentation regarding consistent typing for code block literals.
- [Use `.` for method invocation always?.md](<./Use `.` for method invocation always?.md>) — Standardizing dot-notation for all property access.
- [Use alternative different syntax for generics?.md](<./Use alternative different syntax for generics?.md>) — Exploring non-bracketed approaches to generic parameters.
- [While loop yield and external caller interleaving.md](<./While loop yield and external caller interleaving.md>) — Implementation details of concurrent task yielding during iteration.
- [Yet another alternative to signature syntax.md](<./Yet another alternative to signature syntax.md>) — Exploring `#(...)` versus traditional function pointer signatures.
- [`self` and-or initialize literals.md](<./`self` and-or initialize literals.md>) — How `self` is referenced within literal block initialization blocks.
- [are annotations macros.md](<./are annotations macros.md>) — Implementation of the annotation expansion process.
- [non-local return.md](<./non-local return.md>) — Handling early exits that break out of nested BOC contexts.
- [type alias.md](<./type alias.md>) — Design documentation regarding type abbreviation and naming.
- [utility functions.md](<./utility functions.md>) — Design documentation regarding standard library helper methods.
