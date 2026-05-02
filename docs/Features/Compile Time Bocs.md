#feature

# Yz Compile-Time Bocs

## Overview

Yz provides a compile-time execution system built from regular Yz code. Compile-time bocs run during compilation, have access to the full language, and receive the structural metadata of the boc being compiled as a regular value.

The system is built on one rule:

> **A `compile_time: [...]` variable inside a boc's infostring triggers `Compile` implementations during type inference. Their return values are merged into the parent boc.**

Everything else follows from existing Yz concepts.

See also: [Yz Language Overview](../../README.md) · [Generics](Generics%20Revisited.md) · [Structural Reflection](Structural%20Reflection.md)

---

## The `Compile` Interface

`Compile` is a structural interface in Yz:

```
Compile : {
    Schema #()
    run    #(Boc, Boc)   // receives the parent boc, returns a boc to merge
}
```

`Schema` is an associated type — each implementation fixes it to the concrete shape it expects from the infostring. `run` receives the parent boc and returns a boc whose slots are merged into it.

The `name` that maps an implementation to its infostring variable lives in the implementation's own infostring, not in the interface. When `name` is absent the compiler derives it from `Schema`: if `Schema` has exactly one field, that field's name is used. If `Schema` has multiple fields, `name` is required — a compile error if missing.

See also: [Structural Typing](Structural%20typing.md) · [Boc Type](Block%20type.md)

---

## Triggering Compile Implementations

`Compile` implementations are declared in the `compile_time` variable of the boc's infostring:

```
`compile_time: [Derive, JSON, Logger]`
Person : {
    name String
    age  Int
}
```

During parsing the compiler scans infostrings for `compile_time`. When found, the listed implementations are scheduled to run during type inference — sequentially, in array order. Referenced types are resolved at parse time — `compile_time: [Deribe]` is a compile error.

The boc body carries no compile-time triggering mechanism. All compile-time concerns live in the infostring.

---

## Infostrings as Boc Bodies

An infostring is a backtick-delimited literal placed immediately before a boc or field definition. Its content is a **boc body** — parsed and compiled, but never executed. Multiple concerns live as separate variables inside a single infostring.

```
`
compile_time: [Derive, JSON, GraphQL]
graphql: {
    schema: "https://myapi.com/graphql"
    keep_foo: { "bar" }
}
json: {
    ignore: false
}
`
Movies : {
    `json: { field_name: "movie_title" }`
    title String

    `json: { ignore: true }`
    internal_id String
}
```

`self.infostring` inside a `Compile` implementation is typed as `Boc` — the full shared infostring of the annotated type. Each implementation accesses only the variable it owns, typed via its `Schema` associated type:

```
GraphQL : {
    Schema : #(schema String)        // associated type — declares expected shape
    run #(Boc, Boc) = {
        // self.infostring is Boc
        // self.infostring.graphql is typed as Schema = #(schema String)
        schema_url = self.infostring.graphql.schema   // String — validated
        ...
    }
}

JSON : {
    Schema : #(field_name String, ignore Bool)
    run #(Boc, Boc) = {
        self.fields.forEach({ f Boc
            config      = f.infostring.json            // typed as Schema
            field_name  = config.field_name            // String
            should_skip = config.ignore                // Bool
        })
    }
}
```

`self.infostring.<name>` is not a regular `Boc` field lookup — the return type is determined by the implementation's `Schema`, not by `Boc`'s definition. Each implementation reads its own slice; others are ignored.

| Infostrings | Compile slots |
| --- | --- |
| Passive metadata, never executed | Active code generation |
| Attached to fields or bocs | Reads infostrings via `self.infostring.<name>` |
| Boc body — compiler validated | Full language access |

See also: [Info strings](Info%20strings.md)

---

## Declaring an Implementation

A complete `Compile` implementation declares its `Schema`, optionally its `name` in the infostring, and provides `run`:

```
// name derived from Schema's single field — "documentation"
Doc : {
    Schema : #(documentation String)
    run #(Boc, Boc) = {
        text = self.infostring.documentation   // String
        // generate documentation-related slots
    }
}

// name required — Schema has multiple fields
`name: "server"`
Server : {
    Schema : #(
        port     Int
        protocol String
        timeout  Int
    )
    run #(Boc, Boc) = {
        config   = self.infostring.server
        port     = config.port       // Int
        protocol = config.protocol   // String
        timeout  = config.timeout    // Int
        // generate server configuration slots
    }
}
```

Used as:

```
`
compile_time: [Doc, Server]
documentation: "Handles incoming HTTP requests."
server: {
    port:     8080
    protocol: "https"
    timeout:  30
}
`
Handler : { ... }
```

---

## What Compile Has Access To

A `Compile` boc has access to the full language. The same facilities available at runtime are available during compilation:

```
GraphQL : {
    Schema : #(schema String)
    run #(Boc, Boc) = {
        schema_url = self.infostring.graphql.schema

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

When `compile_time` lists multiple implementations they run in **array order**. Each receives the result of the previous one — the progressively merged boc, not the original.

```
`compile_time: [Derive, Logging, Metrics]`
Container : { T; value T }
```

Execution is sequential:

```
original boc
    → Derive.run(original)       → merged boc₁
    → Logging.run(boc₁)          → merged boc₂  // can see Derive's output
    → Metrics.run(boc₂)          → merged boc₃  // can see both
    → boc₃ continues through type inference
```

**Array order is semantically significant.** If `Logging` generates a method that calls a method generated by `Derive`, `Derive` must appear first. This is a first-class rule, not an implementation detail.

If two implementations generate the same slot name the later one wins. Check for slot name collisions when combining third-party `Compile` implementations.

---

## Constraint Declarations on Compile Implementations

When a `Compile` implementation generates code that calls methods on a type parameter, it introduces constraints on that parameter. Callers of the parent boc must satisfy them.

By strong convention, `Compile` implementations declare the constraints they require on their type parameter using the `#(...)` structural signature syntax. For public implementations this convention should be treated as mandatory.

```
// Inline structural constraint
`compile_time: [Validate]`
Derive : {
    Schema     : #(serialize #(String))
    S          #(serialize #(String))
    run #(Boc, Boc) = { ... }
}

// Via a named interface if one exists
Serializable : #(serialize #(String))
Debuggable   : #(toString  #(String))

Derive : {
    Schema : #(serialize #(String))
    S      #(Serializable, Debuggable, metricsId #(String))
    run #(Boc, Boc) = { ... }
}
```

The compiler surfaces the full flattened constraint set in error messages and tooling, with attribution:

```
error: Point does not satisfy constraint on T in Container
  toString #(String) required by Logging
  Logging is in Container's compile_time array

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

`Compile` implementations run at **boc definition time** — when the boc is compiled, not when it is instantiated. For generic bocs this distinction matters:

```
`compile_time: [Derive]`
Stack : {
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

`Compile` implementations run during type inference, interleaved with it. This gives them access to the inferred type information of the boc they are processing.

```
Source
  │
  ▼
Tokenizer + Parser
  │  (scans infostrings for compile_time)
  ▼
Semantic Analysis
  │
  ▼
Inference ←──────────────────────────┐
  │                                  │
  │  encounters compile_time         │
  │  in infostring of boc            │
  ▼                                  │
Run Compile implementations          │
  (sequential, in array order)       │
  │                                  │
  │  generated boc merged into AST   │
  └──────────────────────────────────┘  re-enter inference
  │
  │  all compile_time slots exhausted
  ▼
IR Generation
  │
  ▼
Codegen
```

When inference encounters a boc whose `compile_time` implementations have not yet run, it triggers them in order, merges each returned `Boc` back into the AST, and continues inference on the result. The loop continues until no new `compile_time` triggers remain.

Generated code is fully analysed — its types are inferred, its constraints propagate, and it participates in the type system identically to hand-written code.

---

## Two-Phase Compilation

`Compile` implementations are full Yz programs that execute during compilation. The compiler needs them as callable native code before it can compile the rest of the source. This is handled by a two-phase build.

**Phase 1 — Bootstrap:** The compiler scans source for bocs satisfying the `Compile` interface (identified by `Schema #()` and `run #(Boc, Boc)`). It compiles only those bocs to native executables ahead of everything else. The standard library and previously compiled packages are available at this phase.

**Phase 2 — Main compilation:** The compiler processes the full source normally. When inference encounters a `compile_time` infostring variable, it calls the pre-compiled executables via subprocess in array order — passing the current partially-inferred `Boc` as serialised data, receiving the generated `Boc` back, merging it into the AST, and resuming inference.

```
Phase 1:
  scan → find Compile implementations → compile to native executables

Phase 2:
  parse → scan infostrings for compile_time → build AST → begin inference
      ↓ encounter compile_time
      serialize current Boc → subprocess call → deserialize returned Boc
      merge into AST → resume inference
```

**Compile implementations must live in a separate package from the bocs they process.** Phase 1 compiles them before the rest of the package exists. A `Compile` implementation in the same package as its target creates a circular build dependency.

**Caching:** Compiled `Compile` executables are cached. The cache key is the hash of the implementation source combined with the hash of the input boc's structure. A subsequent compilation only rebuilds a `Compile` executable when its source changes, and only re-runs it when the input boc's structure changes.

---

## Circular Generation

Mutually triggering `compile_time` arrays produce an infinite loop:

```
`compile_time: [GenB]`
A : { ... }   // GenB generates something that triggers B's compile_time

`compile_time: [GenA]`
B : { ... }   // GenA generates something that triggers A's compile_time
```

The compiler detects this via cycle tracking on `compile_time` execution — the same mechanism used for cycle detection in recursive type inference.

The rule is:

> **A `Compile` implementation cannot trigger the `compile_time` of a boc that is currently being compiled.**

A violation is a compile error. The error message names the full cycle.

---

## Compile Implementations Are Regular Bocs

`Compile` implementations can themselves have infostrings with `compile_time`, type parameters, and associated types. They are bocs that satisfy the `Compile` interface:

```
`compile_time: [Validate]`
Derive : {
    Schema : #(serialize #(String))
    S      #(serialize #(String))
    run #(Boc, Boc) = {
        // implementation
    }
}
```

`Derive`'s own `compile_time` runs when `Derive` is compiled as a library — long before any user boc references it. By the time `Person` lists `Derive` in its `compile_time` array, `Derive` is fully compiled and cached.

---

## Standard Library Compile Implementations

The following `Compile` implementations are provided by the Yz standard library.

| Implementation | `name` | `Schema` | Generates | Declared constraints on `S` |
| --- | --- | --- | --- | --- |
| `Derive` | `"derive"` | `#(interfaces [String])` | Interface implementations | Varies per interface |
| `JSON` | `"json"` | `#(field_name String, ignore Bool)` | Serialization / deserialization | None |
| `Debug` | derived | `#(debug Bool)` | `debug #(String)` method | None |
| `Validate` | `"validate"` | `#(rule String, strict Bool)` | Validation methods | None |

---

## Comparison With Other Languages

| Feature | Yz | Lisp | Rust | Zig | Haskell (`deriving`) |
| --- | --- | --- | --- | --- | --- |
| Same language at compile time | ✅ | ✅ | ❌ separate crate | ✅ | ❌ |
| Full type information available | ✅ | ❌ pre-type-check | ❌ pre-type-check | ✅ | ✅ |
| Attached to target | ✅ infostring | ❌ separate | ❌ separate | ❌ separate | ✅ via `deriving` |
| Passive metadata | Infostrings | N/A | Attributes | N/A | N/A |
| Open / extensible | ✅ any boc | ✅ | ✅ proc-macros | ❌ | ❌ built-in only |
| Constraint declarations | Convention | N/A | Generated bounds | N/A | Implicit per class |
| Constraints visible in tooling | ✅ attributed | N/A | ✅ in generated code | ✅ | ✅ |
| Typed metadata schema | ✅ associated type | ❌ | ❌ | ❌ | ❌ |

The key distinction from Haskell's `deriving` and Rust's `derive` is openness. Both are closed systems where the derivable set is fixed, making constraint side-effects predictable by definition. Yz `Compile` is open — any boc satisfying the interface qualifies — which is why constraint declarations by convention, typed `Schema` declarations, and clear tooling attribution are load-bearing parts of the design.

---

## See Also

- [Structural Reflection](Structural%20Reflection.md) — the full `Boc` API available inside `run`
- [Generics Revisited](Generics%20Revisited.md) — constraint inference and propagation
- [Structural Typing](Structural%20typing.md) — how `#(...)` signatures work
- [Boc Type](Block%20type.md) — boc type syntax including `#(...)`
- [Info strings](Info%20strings.md) — infostring syntax and conventions
- [Conditional Bocs](Conditional%20Bocs.md) — used in generated control flow
- [Type Variants](Type%20variants.md) — used in generated variant dispatch