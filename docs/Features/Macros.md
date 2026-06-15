#feature

# Yz Macros

## Overview

Yz provides a macro system built from regular Yz code. Macros run during compilation, have access to the full language, and receive the structural metadata of the boc being compiled as a regular value.

The system is built on one rule:

> **Uppercase type names in a boc's annotation trigger `Macro` implementations during type inference. Their return values are merged into the subject boc.**

Everything else follows from existing Yz concepts.

See also: [Yz Language Overview](../../README.md) · [Generics](Generics%20-%20Type%20Parameters.md) ·  [Structural Reflection](Structural%20Reflection.md) · [Annotations](Annotations.md)

---

## The `Macro` Interface

`Macro` is a structural interface in Yz:

```
Macro : {
    Schema #()
    run    #(subject Boc, config Schema, Boc)
}
```

`Schema` is an associated type — each implementation fixes it to the concrete shape it expects from the annotation config block. The compiler extracts the config block from the annotation, validates it against `Schema`, and passes it as `config`. `run` returns a boc whose slots are merged into the subject.
	
Since `Schema` is an associated type, `config` is already typed as the concrete schema of each implementation — no casting or special lookup needed.

The macro's own type name is the dispatch key. No `name` annotation is needed.

See also: [Structural Typing](Structural%20typing.md) · [Boc Type](Boc%20Interface.md)

---

## Triggering Macros

Macros are triggered by uppercase type names in a boc's annotation. A bare name triggers with no config; a named boc provides config validated against the macro's `Schema`:

```
`
Derive
JSON: { ignore: false }
Logger
`
Person : {
    name String
    age  Int
}
```

During parsing the compiler scans annotations for uppercase names. When found, the referenced macros are scheduled to run during type inference — sequentially, in top-to-bottom declaration order. Referenced types are resolved at parse time — `Deribe` in an annotation is a compile error if no such macro exists.

The boc body carries no macro-triggering mechanism. All macro concerns live in the annotation.

---

## Annotation Config

The compiler extracts the uppercase block matching the macro's type name, validates it against `Schema`, and passes it as `config` to `run`. Each macro receives only its own config — already typed. Lowercase entries in the annotation are passive metadata for field-walking:

```
`
Derive
JSON: { ignore: false }
Logger: { level: "debug", format: "json" }
`
Movies : {
    `json: { field_name: "movie_title" }`
    title String

    `json: { ignore: true }`
    internal_id String
}
```

```
Logger : {
    Schema : #(level String, format String)
    run #(subject Boc, config Schema, Boc) = {
        level  = config.level    // String
        format = config.format   // String
        // generate logging methods
    }
}

JSON : {
    Schema : #(field_name String, ignore Bool)
    run #(subject Boc, config Schema, Boc) = {
        subject.fields.forEach({ f Boc
            field_config = f.annotation.json         // field-level metadata, lowercase
            field_name   = field_config.field_name   // String
            should_skip  = field_config.ignore       // Bool
        })
    }
}
```

See also: [Annotations](Annotations.md) — annotation syntax, uppercase vs lowercase rules, companion files

---

## Declaring a Macro

A complete macro declares its `Schema` and provides `run`. The macro's type name is the dispatch key — no annotation needed on the macro itself:

```
Doc : {
    Schema : #(documentation String)
    run #(subject Boc, config Schema, Boc) = {
        text = config.documentation   // String
        // generate documentation-related slots
    }
}

Server : {
    Schema : #(
        port     Int
        protocol String
        timeout  Int
    )
    run #(subject Boc, config Schema, Boc) = {
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
Doc: { documentation: "Handles incoming HTTP requests." }
Server: {
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
Routes : {
    Schema : #(spec_file String)
    run #(subject Boc, config Schema, Boc) = {
        // fetch a remote spec at compile time
        spec = http.get(config.spec_file)

        // read a local file at compile time
        query = file.read("routes/base.json")

        // use the full Yz standard library
        routes = spec.parse().routes.map({ r Route
            // generate a boc per route
        })

        { handlers: routes }
    }
}
```

The two special properties of a macro are:

1. **When** it runs — during type inference, at boc definition time
2. **What happens to its return value** — merged into the subject boc's AST

Everything inside is regular Yz.

See also: [Structural Reflection](Structural%20Reflection.md) for the full `Boc` API available inside `run`

---

## Ordering of Multiple Macros

When multiple macros appear in an annotation they run in **top-to-bottom declaration order**. Each receives the result of the previous one — the progressively merged boc, not the original.

```
`
Derive
Logging
Metrics
`
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

**Declaration order is semantically significant.** If `Logging` generates a method that calls a method generated by `Derive`, `Derive` must appear first. This is a first-class rule, not an implementation detail.

If two macros generate the same slot name the later one wins. Check for slot name collisions when combining third-party macros.

---

## Constraints Introduced by Generated Code

Macros add methods and slots to the annotated boc — not to its type parameters. `Debug` adds a `debug #(String)` method to `Container` itself:

```
`Debug`
Container : {
    value T
}
// Container gains: debug #(String)
// T is unaffected — unless the generated implementation calls something on it
```

If the generated `debug` implementation calls a method on `T` (e.g. `value.to_string()`), type inference introduces a constraint on `T`. That constraint is invisible in `Container`'s source — it comes from generated code. The compiler attributes it to the macro:

```
error: Point does not satisfy constraint on T in Container
  to_string #(String)    required by Debug

Container #(value T)
  T requires:
    to_string #(String)    ← Debug
```

There is no declaration mechanism for this in the macro definition. The constraint emerges from type inference on the merged generated code, and attribution is automatic. Whether a macro introduces constraints on type parameters depends entirely on its generated implementation — macro authors should document this when it applies.

See also: [Generics - Type Parameters](Generics%20-%20Type%20Parameters.md)

---

## Execution Timing

Macros run at **boc definition time** — when the boc is compiled, not when it is instantiated. For generic bocs this distinction matters:

```
`Derive`
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
  │  (scans annotations for uppercase names)
  ▼
Semantic Analysis
  │
  ▼
Inference ←──────────────────────────┐
  │                                  │
  │  encounters macro names          │
  │  in annotation of boc            │
  ▼                                  │
Run macros                           │
  (sequential, declaration order)    │
  │                                  │
  │  generated boc merged into AST   │
  └──────────────────────────────────┘  re-enter inference
  │
  │  all macro triggers exhausted
  ▼
IR Generation
  │
  ▼
Codegen
```

When inference encounters a boc whose macros have not yet run, it triggers them in declaration order, merges each returned `Boc` back into the AST, and continues inference on the result. The loop continues until no new macro triggers remain.

Generated code is fully analysed — its types are inferred, its constraints propagate, and it participates in the type system identically to hand-written code.

---

## Two-Phase Compilation

Macros are full Yz programs that execute during compilation. The compiler needs them as callable native code before it can compile the rest of the source. This is handled by a two-phase build.

**Phase 1 — Bootstrap:** The compiler scans source for bocs satisfying the `Macro` interface (identified by `Schema #()` and `run #(subject Boc, config Schema, Boc)`). It compiles only those bocs to native executables ahead of everything else. The standard library and previously compiled packages are available at this phase.

**Phase 2 — Main compilation:** The compiler processes the full source normally. When inference encounters uppercase macro names in an annotation, it calls the pre-compiled executables via subprocess in declaration order — passing the current partially-inferred `Boc` as serialised data, receiving the generated `Boc` back, merging it into the AST, and resuming inference.

```
Phase 1:
  scan → find Macro implementations → compile to native executables

Phase 2:
  parse → scan annotations for uppercase names → build AST → begin inference
      ↓ encounter macro names
      serialize current Boc → subprocess call → deserialize returned Boc
      merge into AST → resume inference
```

**Macros must live in a separate package from the bocs they process.** Phase 1 compiles them before the rest of the package exists. A macro in the same package as its target creates a circular build dependency.

**Caching:** Compiled macro executables are cached. The cache key is the hash of the macro source combined with the hash of the input boc's structure. A subsequent compilation only rebuilds a macro executable when its source changes, and only re-runs it when the input boc's structure changes.

---

## Circular Generation

Mutually triggering macro annotations produce an infinite loop:

```
`GenB`
A : { ... }   // GenB generates something that triggers B's macro annotations

`GenA`
B : { ... }   // GenA generates something that triggers A's macro annotations
```

The compiler detects this via cycle tracking on macro execution — the same mechanism used for cycle detection in recursive type inference.

The rule is:

> **A macro cannot trigger a macro annotation of a boc that is currently being compiled.**

A violation is a compile error. The error message names the full cycle.

---

## Macros Are Regular Bocs

Macros can themselves have macro annotations, type parameters, and associated types. They are bocs that satisfy the `Macro` interface:

```
`Validate`
Derive : {
    Schema : #(interfaces [String])
    run #(subject Boc, config Schema, Boc) = {
        // implementation
    }
}
```

`Derive`'s own macro annotations run when `Derive` is compiled as a library — long before any user boc references it. By the time `Person`'s annotation references `Derive`, `Derive` is fully compiled and cached.

---

## Macro Catalogue

All known and planned macros. Status: **std** = in standard library, **planned** = designed but not yet implemented. Individual design notes for planned ones: [`Macros/`](Macros/)

| Macro      | Status  | `Schema`                                                       | Generates / Purpose                                              |
| ---------- | ------- | -------------------------------------------------------------- | ---------------------------------------------------------------- |
| `Derive`   | std     | `#(interfaces [String])`                                       | Interface implementations                                        |
| `JSON`     | std     | `#(field_name String, ignore Bool)`                            | Serialization / deserialization methods                          |
| `Debug`    | std     | #()                                                            | `debug #(String)` method                                         |
| `Validate` | std     | `#(rule String, strict Bool)`                                  | Validation methods                                               |
| `Build`    | planned | `#(boc String, platform String)`                               | Conditional file inclusion — picks variant body per target       |
| `Mix`      | planned | `#(mix String)`                                                | Copies methods from named boc into annotated boc at compile time |
| `Test`     | planned | `#(boc String, test_engine String)`                            | Marks file as test fragment; grants access to target internals   |
| `Deps`     | planned | `#(dependencies [#(name String, version String, uri String)])` | Fetches dependencies; wires into source root                     |
| `Doc`      | planned | `#(documentation String)`                                      | Generates API docs / IDE hover text from annotation field        |

---

## See Also

- [Structural Reflection](Structural%20Reflection.md) — the full `Boc` API available inside `run`
- [Generics - Type Parameters](Generics%20-%20Type%20Parameters.md) — constraint inference and propagation
- [Structural Typing](Structural%20typing.md) — how `#(...)` signatures work
- [Boc Type](Boc%20Interface.md) — boc type syntax including `#(...)`
- [Annotations](Annotations.md) — annotation syntax and conventions
- [Conditional Bocs](Conditional%20Bocs.md) — used in generated control flow
- [Type Variants](Type%20variants.md) — used in generated variant dispatch
