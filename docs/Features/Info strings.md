#feature

# Yz Infostrings

## Overview

An infostring is a string literal that appears immediately before a boc definition or a field declaration. Its content is a **boc body** — parsed and compiled by the language, but never executed. It carries structured metadata that `Compile` implementations and tooling can read at compile time.

```
"
compile_time: [Derive, Debug, Logging]
graphql: {
    schema: "https://myapi.com/graphql"
}
"
Movies : {
    ...
}
```

The string delimiters take the place of `{` and `}`. Everything inside is valid Yz — field declarations, nested bocs, scalar values — subject to the restrictions described below.

See also: [Compile-Time Bocs](Compile%20Time%20Bocs.md) · [Structural Reflection](Structural%20Reflection.md) · [Boc Type](Block%20type.md)

---

## Syntax

An infostring is placed on the line(s) immediately before the definition it annotates. Yz strings are multiline by default, so a single string spans as many lines as needed.

A minimal infostring with a single scalar:

```
"port: 8080"
Server : { ... }
```

A more structured infostring using nested bocs:

```
"
graphql: {
    schema: "https://myapi.com/graphql"
    keep_foo: { "bar" }
}
json: {
    ignore: false
}
produces: "application/json"
"
Movies : { ... }
```

The content is a boc body — the same syntax as the inside of `{ ... }` — without the braces. Multiple concerns are expressed as separate variables inside the same infostring. There is one infostring per definition.

---

## Infostrings Are Boc Bodies

The content of an infostring is parsed and compiled as a boc. This means:

- Field declarations, nested boc literals, scalar values — all valid
- Referenced types are resolved at compile time — a typo like `Deribe` is a compile error
- Infostrings can reference bocs defined elsewhere in the program

```
"
compile_time: [Derive, JSON, some.package.GraphQL]
"
Movies : { ... }
// error if Derive, JSON, or some.package.GraphQL do not exist
```

Infostrings have restrictions that distinguish them from regular bocs:

- **Never executed.** The boc body is data. No method calls, no side effects, no invocations.
- **No string interpolation.** An infostring lives inside a string literal — template expressions are unavailable.
- **No nested infostrings.** An infostring cannot contain another infostring.

Everything else a boc can contain is available: nested boc literals, arrays, scalar values, references to named types.

---

## `compile_time` — Triggering Compile Implementations

The reserved variable `compile_time` inside an infostring declares the `Compile` implementations the compiler should run on the annotated boc. It is always an array:

```
"compile_time: [Derive, Debug, Logging]"
Person : {
    name String
    age  Int
}
```

During parsing, the compiler scans infostrings for `compile_time`. When found, the listed implementations are scheduled to run during type inference — sequentially, in array order — exactly as if they had been declared in a `compile [Compile] = [...]` field on the boc.

A `Compile` implementation field on the boc itself (`compile [Compile] = [...]`) and `compile_time` in the infostring are distinct mechanisms and can coexist. `Movies.compile` is a field on the `Movies` boc. The infostring's `compile_time` is a property of `Movies`'s metadata — accessible as `Movies_boc.infostring.compile_time`. They do not conflict.

---

## How Compile Implementations Read Infostrings

Inside a `Compile` implementation's `run` method, the parent boc's infostring is accessible as a regular boc value via `self.infostring`. Each variable in the infostring is a field on that value:

```
GraphQL : {
    run #(Boc, Boc) = {
        config     = self.infostring.graphql    // the graphql: { ... } boc
        schema_url = config.schema              // "https://myapi.com/graphql"
        ...
    }
}

JSON : {
    run #(Boc, Boc) = {
        config      = self.infostring.json      // the json: { ... } boc
        should_skip = config.ignore             // false
        ...
    }
}
```

Each `Compile` implementation reads the variable it cares about. Others are ignored. No string parsing is required — field access is direct and type-checked.

Dynamic access by string key (`infostring["graphql"]`) is under consideration for cases where the key is not known at author time. For now, direct field access is the supported form.

---

## Field-Level Infostrings

Infostrings follow the same rules on individual fields:

```
Movies : {
    compile_time: [GraphQL, JSON]

    "graphql: { rename: 'movieTitle' }"
    title String

    "json: { ignore: true }"
    internal_id String

    "
    json: { field_name: 'release' }
    validate: { min: 1888 }
    "
    year Int
}
```

A `Compile` implementation reads field infostrings via `self.fields[n].infostring`:

```
JSON : {
    run #(Boc, Boc) = {
        self.fields.forEach({ f Boc
            config = f.infostring.json
            config.ignore ? { /* skip this field */ }
        })
    }
}
```

`compile_time` is syntactically valid in a field-level infostring, but triggering per-field `Compile` implementations is a topic for a separate design.

---

## Simple Values

Not all infostrings need nested boc structure. Scalar values are valid:

```
"port: 8080"
Server : { ... }

"produces: 'application/json'"
Handler : { ... }

"json: 'ignore'"
internal_id String
```

A `Compile` implementation reading `self.infostring.json` receives the string `"ignore"` directly. Whether to use a scalar or a nested boc is a choice for the `Compile` implementation's author to document.

---

## Documentation in Infostrings

Documentation can live in an infostring as a regular variable:

```
"
documentation: "
    The Movies boc represents a film record as returned by the catalogue API.
    Fields map directly to the GraphQL schema unless annotated otherwise.
"
compile_time: [GraphQL, JSON]
graphql: { schema: 'https://myapi.com/graphql' }
"
Movies : { ... }
```

A `Documentation` compile implementation could read `self.infostring.documentation` and generate API docs, IDE hover text, or reference pages.

Regular code comments (`/* ... */` and `//`) remain available for inline code explanation and are not accessible to `Compile` implementations. Documentation intended for tooling or generated output belongs in the infostring.

---

## Infostrings in the `Boc` Metatype

Infostrings are part of the `Boc` metatype. The compiler parses the infostring and attaches the resulting compiled boc to the `Boc` instance it creates for every definition:

```
Boc : {
    name        String
    instantiable Bool
    fields      [Boc]
    methods     [Boc]
    type_params [Boc]
    infostring  Boc      // the compiled infostring boc, or an empty boc if absent
    source      #()
}
```

`infostring` is a regular `Boc` value. It can be passed, iterated, serialized, and sent across a wire — the same as any other `Boc`. The `compile_time` field within it is what the compiler acts on; all other fields are data for `Compile` implementations and tooling to use freely.

See also: [Structural Reflection](Structural%20Reflection.md) for the full `Boc` API

---

## Summary of Rules

| Property | Rule |
| --- | --- |
| Syntax | Boc body without `{` `}`, inside a string literal |
| Placement | Immediately before a boc or field definition |
| Per definition | One infostring per definition |
| Multiple concerns | Separate variables inside the same infostring |
| Execution | Never executed |
| String interpolation | Not available |
| Nested infostrings | Not allowed |
| Type references | Resolved at compile time — missing types are errors |
| `compile_time` | Always an array; triggers listed `Compile` implementations |
| Access in `Compile` | `self.infostring.<variable>` |
| In `Boc` metatype | `Boc.infostring` — a regular compiled `Boc` value |