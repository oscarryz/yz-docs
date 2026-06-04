#feature

# Yz Macros

## Overview

Yz provides a macro system built from regular Yz code. Macros run during compilation, have access to the full language, and receive the structural metadata of the boc being compiled as a regular value.

The system is built on one rule:

> **A `macros: [...]` variable inside a boc's annotation triggers `Macro` implementations during type inference. Their return values are merged into the parent boc.
> The variable name `!` is an alias for `macros`.**

Everything else follows from existing Yz concepts.

See also: [Yz Language Overview](../../README.md) · [Generics](Generics%20-%20Type%20Parameters.md) ·  [Structural Reflection](Structural%20Reflection.md) · [Annotations](Annotations.md)

---

## The `Macro` Interface

`Macro` is a structural interface in Yz:

```
Macro : {
    Schema #()
    run    #(parent Boc, Boc)   // parent = input boc; output = boc to merge
}
```

`Schema` is an associated type — each implementation fixes it to the concrete shape it expects from the annotation. `run` receives the parent boc (`parent`) and returns a boc whose slots are merged into it.

The `name` that maps a macro to its annotation variable lives in the macro's own annotation, not in the interface. When `name` is absent the compiler derives it from `Schema`: if `Schema` has exactly one field, that field's name is used. If `Schema` has multiple fields, `name` is required — a compile error if missing.

See also: [Structural Typing](Structural%20typing.md) · [Boc Type](Boc%20Interface.md)

---

## Triggering Macros

Macros are declared in the `macros` variable of the boc's annotation (`!` is an alias):

```
`macros: [Derive, JSON, Logger]`
Person : {
    name String
    age  Int
}
```

During parsing the compiler scans annotations for `macros`. When found, the listed macros are scheduled to run during type inference — sequentially, in array order. Referenced types are resolved at parse time — `macros: [Deribe]` is a compile error.

The boc body carries no macro-triggering mechanism. All macro concerns live in the annotation.

---

## Annotations as Boc Bodies

An annotation is a backtick-delimited literal placed immediately before a boc or field definition. Its content is a **boc body** — parsed and compiled, but never executed. Multiple concerns live as separate variables inside a single annotation.

```
`
macros: [Derive, JSON, GraphQL]
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

`self.annotation` inside a macro is typed as `Boc` — the full shared annotation of the annotated type. Each macro accesses only the variable it owns, typed via its `Schema` associated type:

```
GraphQL : {
    Schema : #(schema String)        // associated type — declares expected shape
    run #(parent Boc, Boc) = {
        // self.annotation is Boc
        // self.annotation.graphql is typed as Schema = #(schema String)
        schema_url = self.annotation.graphql.schema   // String — validated
        ...
    }
}

JSON : {
    Schema : #(field_name String, ignore Bool)
    run #(parent Boc, Boc) = {
        self.fields.forEach({ f Boc
            config      = f.annotation.json            // typed as Schema
            field_name  = config.field_name            // String
            should_skip = config.ignore                // Bool
        })
    }
}
```

`self.annotation.<name>` is not a regular `Boc` field lookup — the return type is determined by the macro's `Schema`, not by `Boc`'s definition. Each macro reads its own slice; others are ignored.

| Annotations | Macro slots |
| --- | --- |
| Passive metadata, never executed | Active code generation |
| Attached to fields or bocs | Reads annotations via `self.annotation.<name>` |
| Boc body — compiler validated | Full language access |

See also: [Annotations](Annotations.md)

---

## Declaring a Macro

A complete macro declares its `Schema`, optionally its `name` in the annotation, and provides `run`:

```
// name derived from Schema's single field — "documentation"
Doc : {
    Schema : #(documentation String)
    run #(parent Boc, Boc) = {
        text = self.annotation.documentation   // String
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
    run #(parent Boc, Boc) = {
        config   = self.annotation.server
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
macros: [Doc, Server]
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

## What a Macro Has Access To

A macro has access to the full language. The same facilities available at runtime are available during compilation:

```
GraphQL : {
    Schema : #(schema String)
    run #(parent Boc, Boc) = {
        schema_url = self.annotation.graphql.schema

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

The two special properties of a macro are:

1. **When** it runs — during type inference, at boc definition time
2. **What happens to its return value** — merged into the parent boc's AST

Everything inside is regular Yz.

See also: [Structural Reflection](Structural%20Reflection.md) for the full `Boc` API available inside `run`

---

## Ordering of Multiple Macros

When `macros` lists multiple macros they run in **array order**. Each receives the result of the previous one — the progressively merged boc, not the original.

```
`macros: [Derive, Logging, Metrics]`
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

If two macros generate the same slot name the later one wins. Check for slot name collisions when combining third-party macros.

---

## Constraint Declarations on Macros

When a macro generates code that calls methods on a type parameter, it introduces constraints on that parameter. Callers of the parent boc must satisfy them.

By strong convention, macros declare the constraints they require on their type parameter using the `#(...)` structural signature syntax. For public macros this convention should be treated as mandatory.

```
// Inline structural constraint
`macros: [Validate]`
Derive : {
    Schema     #(serialize #(String))
    S          #(serialize #(String))
    run #(parent Boc, Boc) = { ... }
}

// Via a named interface if one exists
Serializable : #(serialize #(String))
Debuggable   : #(toString  #(String))

Derive : {
    Schema : #(serialize #(String))
    S      #(Serializable, Debuggable, metricsId #(String))
    run #(parent Boc, Boc) = { ... }
}
```

The compiler surfaces the full flattened constraint set in error messages and tooling, with attribution:

```
error: Point does not satisfy constraint on T in Container
  toString #(String) required by Logging
  Logging is in Container's macros array

// tooling surfaces:
Container #(value T)
  T requires:
    serialize  #(String)    ← declared by Derive
    toString   #(String)    ← declared by Logging
    metricsId  #(String)    ← declared by Metrics
```

See also: [Generics - Type Parameters](Generics%20-%20Type%20Parameters.md) — constraints are optional for regular generics but by convention required for macros

---

## Execution Timing

Macros run at **boc definition time** — when the boc is compiled, not when it is instantiated. For generic bocs this distinction matters:

```
`macros: [Derive]`
Stack : {
    T
    items []T
    push  #(T)
    pop   #(T)
}

istack Stack(Int)    // Macro already ran — Stack is fully defined
sstack Stack(String) // same — no macro re-execution
```

Generated code contains type parameters and is monomorphized later at instantiation, exactly as hand-written generic code would be.

---

## Compilation Lifecycle

Macros run during type inference, interleaved with it. This gives them access to the inferred type information of the boc they are processing.

```
Source
  │
  ▼
Tokenizer + Parser
  │  (scans annotations for macros)
  ▼
Semantic Analysis
  │
  ▼
Inference ←──────────────────────────┐
  │                                  │
  │  encounters macros               │
  │  in annotation of boc            │
  ▼                                  │
Run macros                           │
  (sequential, in array order)       │
  │                                  │
  │  generated boc merged into AST   │
  └──────────────────────────────────┘  re-enter inference
  │
  │  all macros slots exhausted
  ▼
IR Generation
  │
  ▼
Codegen
```

When inference encounters a boc whose macros have not yet run, it triggers them in order, merges each returned `Boc` back into the AST, and continues inference on the result. The loop continues until no new macro triggers remain.

Generated code is fully analysed — its types are inferred, its constraints propagate, and it participates in the type system identically to hand-written code.

---

## Two-Phase Compilation

Macros are full Yz programs that execute during compilation. The compiler needs them as callable native code before it can compile the rest of the source. This is handled by a two-phase build.

**Phase 1 — Bootstrap:** The compiler scans source for bocs satisfying the `Macro` interface (identified by `Schema #()` and `run #(parent Boc, Boc)`). It compiles only those bocs to native executables ahead of everything else. The standard library and previously compiled packages are available at this phase.

**Phase 2 — Main compilation:** The compiler processes the full source normally. When inference encounters a `macros` annotation variable, it calls the pre-compiled executables via subprocess in array order — passing the current partially-inferred `Boc` as serialised data, receiving the generated `Boc` back, merging it into the AST, and resuming inference.

```
Phase 1:
  scan → find Macro implementations → compile to native executables

Phase 2:
  parse → scan annotations for macros → build AST → begin inference
      ↓ encounter macros
      serialize current Boc → subprocess call → deserialize returned Boc
      merge into AST → resume inference
```

**Macros must live in a separate package from the bocs they process.** Phase 1 compiles them before the rest of the package exists. A macro in the same package as its target creates a circular build dependency.

**Caching:** Compiled macro executables are cached. The cache key is the hash of the macro source combined with the hash of the input boc's structure. A subsequent compilation only rebuilds a macro executable when its source changes, and only re-runs it when the input boc's structure changes.

---

## Circular Generation

Mutually triggering `macros` arrays produce an infinite loop:

```
`macros: [GenB]`
A : { ... }   // GenB generates something that triggers B's macros

`macros: [GenA]`
B : { ... }   // GenA generates something that triggers A's macros
```

The compiler detects this via cycle tracking on macro execution — the same mechanism used for cycle detection in recursive type inference.

The rule is:

> **A macro cannot trigger the `macros` of a boc that is currently being compiled.**

A violation is a compile error. The error message names the full cycle.

---

## Macros Are Regular Bocs

Macros can themselves have annotations with `macros`, type parameters, and associated types. They are bocs that satisfy the `Macro` interface:

```
`macros: [Validate]`
Derive : {
    Schema : #(serialize #(String))
    S      #(serialize #(String))
    run #(parent Boc, Boc) = {
        // implementation
    }
}
```

`Derive`'s own `macros` run when `Derive` is compiled as a library — long before any user boc references it. By the time `Person` lists `Derive` in its `macros` array, `Derive` is fully compiled and cached.

---

## Macro Catalogue

All known and planned macros. Status: **std** = in standard library, **planned** = designed but not yet implemented. Individual design notes for planned ones: [`Macros/`](Macros/)

| Macro      | Status  | `name`            | `Schema`                                                       | Generates / Purpose                                            | Constraints on `S` |
| ---------- | ------- | ----------------- | -------------------------------------------------------------- | -------------------------------------------------------------- | ------------------ |
| `Derive`   | std     | `"derive"`        | `#(interfaces [String])`                                       | Interface implementations                                      | Varies per iface   |
| `JSON`     | std     | `"json"`          | `#(field_name String, ignore Bool)`                            | Serialization / deserialization methods                        | None               |
| `Debug`    | std     | derived           | `#(debug Bool)`                                                | `debug #(String)` method                                       | None               |
| `Validate` | std     | `"validate"`      | `#(rule String, strict Bool)`                                  | Validation methods                                             | None               |
| `Build`    | planned | `"build"`         | `#(boc String, platform String)`                               | Conditional file inclusion — picks variant body per target     | None               |
| `Mix`      | planned | `"mix"`           | `#(mix String)`                                                | Copies methods from named boc into annotated boc at compile time | None             |
| `Test`     | planned | `"test"`          | `#(boc String, test_engine String)`                            | Marks file as test fragment; grants access to target internals | None               |
| `Deps`     | planned | `"deps"`          | `#(dependencies [#(name String, version String, uri String)])` | Fetches dependencies; wires into source root                   | None               |
| `Doc`      | planned | `"documentation"` | `#(documentation String)`                                      | Generates API docs / IDE hover text from annotation field      | None               |

---

## See Also

- [Structural Reflection](Structural%20Reflection.md) — the full `Boc` API available inside `run`
- [Generics - Type Parameters](Generics%20-%20Type%20Parameters.md) — constraint inference and propagation
- [Structural Typing](Structural%20typing.md) — how `#(...)` signatures work
- [Boc Type](Boc%20Interface.md) — boc type syntax including `#(...)`
- [Annotations](Annotations.md) — annotation syntax and conventions
- [Conditional Bocs](Conditional%20Bocs.md) — used in generated control flow
- [Type Variants](Type%20variants.md) — used in generated variant dispatch
