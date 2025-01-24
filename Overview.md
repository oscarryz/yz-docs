This is an overview of the language and its features. A more comprehensive and detailed explanation can be found in the [Yz language](Yz%20language.md)  document.

```js
yz: {
  type_discipline: ["static", "strong", "structural", "inferred"]
  memory_management: ["garbage collected", "single writer"]
  concurrency_model: ["actor oriented"]
  programming_style: ["Object Oriented"]
}
```
## What is so different 
*Features that programming language try to avoid or are not common*

- Everything is a block of code.
  - objects, methods, functions, threads, control structures.
- Everything is public, mutable and concurrent
- Single writer
## What is not that different
*Features other programming languages include*

- Strong, static typing with structural type and type inference.
- Garbage Collected
- Result based error handling
- Minimal syntax

## What looks similar but is not
*Features that programming languages have... but different*

- Object Oriented; but without classes, nor inheritance and limited reusability
- Structural Concurrency; but without Scopes/Nurseries
- Actor Oriented; without whatever actor oriented actually is.
- (Barely) Functional; but obviously with side effects

# Everything is a Block of Code (boc)

```js
block_name : { 
  // block body
}
```
