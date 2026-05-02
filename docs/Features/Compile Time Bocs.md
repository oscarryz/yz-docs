#feature

# Yz Compile-Time Bocs

## Overview

Yz provides a compile-time execution system built from regular Yz code. Compile-time bocs run during compilation, have access to the full language, and receive the structural metadata of the boc being compiled as a regular value.

The system is built on one rule:

> **Any variable typed `Compile` or `[Compile]` is executed by the compiler during type inference. Its return value is merged into the parent boc.**

Everything else follows from existing Yz concepts.

See also: [Yz Language Overview](../../README.md) · [Generics](Generics%20Revisited.md) · [Structural Reflection](Structural%20Reflection.md)

---

## The `Compile` Interface

`Compile` is a structural interface in Yz:

```
Compile : {
    run #(Boc, Boc)   // receives the parent boc, returns a boc to merge
}
```

Any boc that satisfies this interface can be used in a `Compile` variable. The compiler recognises the type as the trigger. The variable can be named anything:

```
Person : {
    compile [Compile] = [Derive, JSON]  // named 'compile' — conventional
    name String
}

Person : {
    derives [Compile] = [Derive, JSON]  // named 'derives' — also valid
    name String
}
```

The type does the work. The name is irrelevant.

See also: [Structural Typing](Structural%20typing.md) · [Boc Type](Block%20type.md)

---

## Basic Usage

A single `Compile` implementation:

```
Logger : {
    run #(Boc, Boc) = {
        // generates a log() method on the parent boc
    }
}

Person : {
    logger Compile = Logger()
    name String
    age  Int
}
```

Multiple `Compile` implementations via array:

```
Person : {
    compile [Compile] = [Derive, JSON, Logger]
    name String
    age  Int
}
```

---

## Field Metadata — Infostrings

Infostrings provide passive metadata to `Compile` implementations. They attach to the boc or to individual fields and are readable via `self.infostring` and `self.fields[n].infostring` inside any `Compile` boc.

Infostrings are regular Yz strings — multiline by default — and their content can be a boc literal. A boc literal inside an infostring is compiled and available as structured data, but **never executed**. This gives `Compile` implementations a consistent, compiler-validated format to read from, rather than each implementation inventing its own string parsing convention.

```
"graphql: {
    schema: 'https://myapi.com/graphql'
    keep_foo: { 'bar' }
}"
Movies : {
    compile Compile = GraphQL

    "json: {
        field_name: 'movie_title'
    }"
    title String

    "json: {
        ignore: true
    }"
    internal_id String
}
```

A `Compile` implementation reads the infostring boc as a regular value — field access, iteration, pattern matching — with no string parsing required:

```
GraphQL : {
    run #(Boc, Boc) = {
        config    = self.infostring("graphql")  // returns a Boc
        schema_url = config.schema              // 'https://myapi.com/graphql'
        kept       = config.keep_foo            // { 'bar' } — a Boc, not executed
        ...
    }
}
```

The boc inside the infostring is data. Its fields are readable. Its body is never invoked. The contract is the same as a plain string infostring — describe, never act — but the structure is validated by the compiler and shared across any `Compile` implementation that reads it.

| Infostrings | Compile slots |
| --- | --- |
| Passive metadata | Active code generation |
| Attached to fields or bocs | Reads infostrings |
| Plain strings or structured boc literals | The engine |
| Never executed | Full language access |

See also: [Info strings](Info%20strings.md)

---

## What Compile Has Access To

A `Compile` boc has access to the full language. The same facilities available at runtime are available during compilation:

```
GraphQL : {
    run #(Boc, Boc) = {
        config     = self.infostring("graphql")  // returns the infostring boc
        schema_url = config.schema               // structured field access — no parsing

        // fetch a remote schema at compile time
        schema = http.get(schema_url)

        // read a local file at compile time
        query = file.read("queries/movies.gql")

        // use the full Yz standard library
        types = schema.parse().types.map({ t Type
            // generate a boc per type
        })

        { queries: types }
    }
}
```

The two special properties of a `Compile` boc are:

1. **When** it runs — during type inference, at boc definition time
2. **What happens to its return value** — merged into the parent boc's AST

Everything inside is regular Yz.

See also: [Structural Reflection](Structural%20Reflection.md) for the full `Boc` API available inside `run`

---

## Ordering of Multiple Compile Implementations

When `[Compile]` contains multiple implementations they run in **array order**. Each implementation receives the result of the previous one — the progressively merged boc, not the original.

```
compile [Compile] = [Derive, Logging, Metrics]
```

Execution is sequential:

```
original boc
    → Derive.run(original)       → merged boc₁
    → Logging.run(boc₁)          → merged boc₂  // can see Derive's output
    → Metrics.run(boc₂)          → merged boc₃  // can see both
    → boc₃ continues through type inference
```

**Array order is semantically significant.** If `Logging` generates a method that calls a method generated by `Derive`, `Derive` must appear first in the array. This is a first-class rule, not an implementation detail.

If two implementations generate the same slot name the later one wins. Check for slot name collisions when combining third-party `Compile` implementations.

---

## Constraint Declarations on Compile Implementations

When a `Compile` implementation generates code that calls methods on a type parameter, it introduces constraints on that parameter. These constraints are real — callers of the parent boc must satisfy them.

By strong convention, `Compile` implementations declare the constraints they require on their type parameter using the `#(...)` structural signature syntax:

```
// Inline structural constraint
Derive : {
    S #(serialize #(String))
    run #(Boc, Boc) = { ... }
}

// Via a named interface if one exists
Serializable : #(serialize #(String))

Derive : {
    S Serializable
    run #(Boc, Boc) = { ... }
}
```

Multiple constraints are listed inside `#(...)` with commas. Named interfaces and inline method signatures can be freely mixed:

```
Serializable : #(serialize #(String))
Debuggable   : #(toString  #(String))

Derive : {
    S #(Serializable, Debuggable, metricsId #(String))
    run #(Boc, Boc) = { ... }
}
```

This declaration is the difference between a library user seeing the constraint in tooling at authoring time versus discovering it as a compiler error after instantiation. For public `Compile` implementations this convention should be treated as mandatory.

The compiler surfaces the full flattened constraint set — including constraints from generated code — in error messages and tooling. When a call fails, attribution identifies the responsible `Compile` implementation:

```
error: Point does not satisfy constraint on T in Container
  toString #(String) required by Logging
  Logging is in Container's [Compile] array

// tooling surfaces:
Container #(value T)
  T requires:
    serialize  #(String)    ← declared by Derive
    toString   #(String)    ← declared by Logging
    metricsId  #(String)    ← declared by Metrics
```

See also: [Generics Revisited](Generics%20Revisited.md) — constraints are optional for regular generics but by convention required for `Compile` implementations

---

## Execution Timing

`Compile` slots run at **boc definition time** — when the boc is compiled, not when it is instantiated. For generic bocs this distinction matters:

```
Stack : {
    compile Compile = Derive()
    T
    items []T
    push  #(T)
    pop   #(T)
}

istack Stack(Int)    // Compile already ran — Stack is fully defined
sstack Stack(String) // same — no Compile re-execution
```

Generated code contains type parameters and is monomorphized later at instantiation, exactly as hand-written generic code would be.

---

## Compilation Lifecycle

`Compile` slots run during type inference, interleaved with it. This gives them access to the inferred type information of the boc they are processing.

```
Source
  │
  ▼
Tokenizer + Parser
  │
  ▼
Semantic Analysis
  │
  ▼
Inference ←──────────────────────────┐
  │                                  │
  │  encounters compile slot         │
  │  (or method on boc with          │
  │   Compile slots not yet run)     │
  ▼                                  │
Run Compile slot                     │
  │                                  │
  │  generated boc merged into AST   │
  └──────────────────────────────────┘  re-enter inference
  │
  │  all Compile slots exhausted
  ▼
IR Generation
  │
  ▼
Codegen
```

Inference and `Compile` execution are interleaved. When inference encounters a `Compile` slot, it triggers that slot, merges the returned `Boc` back into the AST, and continues inference on the result. The loop continues until no new `Compile` slots remain.

Generated code is fully analysed — its types are inferred, its constraints propagate, and it participates in the type system identically to hand-written code.

---

## Two-Phase Compilation

`Compile` implementations are full Yz programs that execute during compilation, so the compiler needs them as callable native code before it can compile the rest of the source. This is handled by a two-phase build.

**Phase 1 — Bootstrap:** The compiler scans source for bocs satisfying the `Compile` interface (identified by `run #(Boc, Boc)`). It compiles only those bocs to native executables ahead of everything else. The standard library and previously compiled packages are available at this phase.

**Phase 2 — Main compilation:** The compiler processes the full source normally. When inference encounters a `Compile` slot, it calls the pre-compiled executable via subprocess, passing the current partially-inferred `Boc` as serialised data, receives the generated `Boc` back, merges it into the AST, and resumes inference.

```
Phase 1:
  scan → find Compile implementations → compile to native executables

Phase 2:
  parse → build AST → begin inference
      ↓ encounter compile slot
      serialize current Boc → subprocess call → deserialize returned Boc
      merge into AST → resume inference
```

**Compile implementations must live in a separate package from the bocs they process.** Phase 1 compiles them before the rest of the package exists. A `Compile` implementation in the same package as its target creates a circular build dependency.

**Caching:** Compiled `Compile` executables are cached. The cache key is the hash of the implementation source combined with the hash of the input boc's structure. A subsequent compilation only rebuilds a `Compile` executable when its source changes, and only re-runs it when the input boc's structure changes.

---

## Circular Generation

Mutually triggering `Compile` slots produce an infinite loop:

```
A : {
    compile [Compile] = [GenB]   // generates something that triggers B's Compile
}

B : {
    compile [Compile] = [GenA]   // generates something that triggers A's Compile
}
```

The compiler detects this via cycle tracking on `Compile` slot execution — the same mechanism used for cycle detection in recursive type inference.

The rule is:

> **A `Compile` slot cannot trigger the `Compile` slots of a boc that is currently being compiled.**

A violation is a compile error. The error message names the full cycle.

---

## Compile Implementations Are Regular Bocs

`Compile` implementations can themselves have `Compile` slots, type parameters, and infostrings. They are bocs that happen to satisfy the `Compile` interface:

```
Derive : {
    compile [Compile] = [Validate]  // Derive has its own Compile slot
    S #(serialize #(String))
    run #(Boc, Boc) = {
        // implementation
    }
}
```

`Derive`'s own `Compile` slot runs when `Derive` is compiled as a library — long before any user boc includes it. By the time `Person` references `Derive` in its `[Compile]` array, `Derive` is fully compiled and cached.

---

## Standard Library Compile Implementations

The following `Compile` implementations are provided by the Yz standard library. Each documents its infostring format and the constraints it declares on its type parameter.

| Implementation | Infostring | Generates | Declared constraints on `S` |
| --- | --- | --- | --- |
| `Derive` | `"derive: { interfaces: [...] }"` | Interface implementations | Varies per interface |
| `JSON` | `"json: { field_name: '...' }"` / `"json: { ignore: true }"` | Serialization / deserialization | None |
| `GraphQL` | `"graphql: { schema: '...' }"` | Typed query bocs | None |
| `Debug` | None | `debug #(String)` method | None |
| `Validate` | `"validate: { rule: '...' }"` | Validation methods | None |

---

## Comparison With Other Languages

| Feature | Yz | Lisp | Rust | Zig | Haskell (`deriving`) |
| --- | --- | --- | --- | --- | --- |
| Same language at compile time | ✅ | ✅ | ❌ separate crate | ✅ | ❌ |
| Full type information available | ✅ | ❌ pre-type-check | ❌ pre-type-check | ✅ | ✅ |
| Attached to target | ✅ inside boc | ❌ separate | ❌ separate | ❌ separate | ✅ via `deriving` |
| Passive metadata | Infostrings | N/A | Attributes | N/A | N/A |
| Open / extensible | ✅ any boc | ✅ | ✅ proc-macros | ❌ | ❌ built-in only |
| Constraint declarations | Convention | N/A | Generated bounds | N/A | Implicit per class |
| Constraints visible in tooling | ✅ attributed | N/A | ✅ in generated code | ✅ | ✅ |

The key distinction from Haskell's `deriving` and Rust's `derive` is openness. Both are closed systems where the derivable set is fixed, making constraint side-effects predictable by definition. Yz `Compile` is open — any boc satisfying the interface qualifies — which is why constraint declarations by convention and clear tooling attribution are load-bearing parts of the design.

---

## See Also

- [Structural Reflection](Structural%20Reflection.md) — the full `Boc` API available inside `run`
- [Generics Revisited](Generics%20Revisited.md) — constraint inference and propagation
- [Structural Typing](Structural%20typing.md) — how `#(...)` signatures work
- [Boc Type](Block%20type.md) — boc type syntax including `#(...)`
- [Info strings](Info%20strings.md) — infostring format and conventions
- [Conditional Bocs](Conditional%20Bocs.md) — used in generated control flow
- [Type Variants](Type%20variants.md) — used in generated variant dispatch