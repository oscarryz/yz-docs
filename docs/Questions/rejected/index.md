# Rejected — discarded design proposals

Design ideas and alternatives for grammar, object literals, and type systems that were debated but ultimately declined for the Yz language.

## Files

- [Block type alternative.md](<./Block type alternative.md>) — Design documentation regarding Block type alternative (using `::` and `;`).
- [Casting.md](./Casting.md) — Discussion on explicit casting mechanisms that was declined.
- [Consider Immutability.md](<./Consider Immutability.md>) — Explorations into strict immutability constraints.
- [Consider each method invocation would create a boc copy.md](<./Consider each method invocation would create a boc copy.md>) — Hypothesis regarding method call overheads and BOC copying.
- [Do we need function overloading.md](<./Do we need function overloading.md>) — Argument against supporting overloaded functions in Yz.
- [Everything is an expression ?.md](<./Everything is an expression ?.md>) — Explorations into expression-only grammar.
- [Extension methods.md](<./Extension methods.md>) — Proposal for adding methods to external types without inheritance.
- [How to deep copy?.md](<./How to deep copy?.md>) — Resolving the mechanics of value vs reference copying for blocks.
- [How to make two identical blocks be different.md](<./How to make two identical blocks be different.md>) — Managing object identity and equality in a block-based system.
- [How to share async data between actors.md](<./How to share async data between actors.md>) — Resolved via channels/actors pattern.
- [Immutability.md](./Immutability.md) — General rules regarding state changes in Yz variables.
- [Immutable outside of block.md](<./Immutable outside of block.md>) — Constraints on external variable immutability.
- [Import - Embed.md](<./Import - Embed.md>) — File embedding vs standard import logic.
- [Instances for libraries.md](<./Instances for libraries.md>) — Creating instances from external library types.
- [No braces `{}`.md](<./No braces `{}`.md>) — Discarded proposal for significant whitespace syntax.
- [Object literal.md](<./Object literal.md>) — Alternative syntaxes rejected in favor of Bocs (Blocks of Code).
- [Read and Write to variables.md](<./Read and Write to variables.md>) — Grammar rules for variable state mutation.
- [References vs Copy.md](<./References vs Copy.md>) — Clarifications on pass-by-reference versus pass-by-value semantics.
- [Skip Commas?.md](<./Skip Commas?.md>) — Exploring optional comma separators in lists and arguments.
- [Stateless bocs and pure functions.md](<./Stateless bocs and pure functions.md>) — Identifying and optimizing side-effect-free blocks.
- [The block type.md](<./The block type.md>) — Defining the core structure of the `boc` data type.
- [Type definition with type keyword.md](<./Type definition with type keyword.md>) — Syntax variations involving a reserved 'type' keyword.
- [Types as Values.md](<./Types as Values.md>) — Treating type definitions themselves as passable variables.
- [Varargs?.md](<./Varargs?.md>) — Support for variable-length argument lists.
- [Variables as data types.md](<./Variables as data types.md>) — Meta-programming concepts regarding variable typing.
- [Wait for a value to be set.md](<./Wait for a value to be set.md>) — Awaiting asynchronous block completions.
- [import.md](./import.md) — Import mechanics and module dependencies.
