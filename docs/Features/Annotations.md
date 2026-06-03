#feature

# Yz Annotations

## Overview

An annotation is a boc body that appears immediately before a boc definition or a field declaration. Its content is parsed and compiled by the language, but never executed. It carries structured metadata that macros and tooling can read at compile time.

```yz
`
macros: [Debug, Logging]
graphql: {
    schema: "https://myapi.com/graphql"
}
`
Movies : {
    ...
}
```

The backtick delimiters take the place of `{` and `}`. Everything inside is valid Yz — field declarations, nested bocs, scalar values — subject to the restrictions described below.

See also: [Macros](Macros.md) · [Structural Reflection](Structural%20Reflection.md) · [Boc Type](Boc%20Interface.md)

---

## Syntax

There are two declaration forms:

**Form 1 — Inline backtick literal** (placed immediately before the definition it annotates):

```yz
`port: 8080`
Server : { ... }
```

**Form 2 — Companion file** (`name.info` in the same directory as the boc named `name`):

```
net/
  http.info   ← annotation body for net.http
  http.yz     ← defines net.http
```

The content of `http.info` is the same annotation body syntax — a boc body without `{ }`. A `name.info` file is never a boc; the compiler attaches its content to the named boc's annotation slot and triggers macros from there.

Use Form 1 for inline annotations on individual declarations. Use Form 2 when the annotation targets the boc defined by the file itself (which has no "line before it"), or when you want to keep metadata separate from code — for example, project dependency declarations.

A minimal annotation with a single scalar:

```
`port: 8080`
Server : { ... }
```

A more structured annotation using nested bocs:

```
`
graphql: {
    schema: "https://myapi.com/graphql"
    keep_foo: { "bar" }
}
json: {
    ignore: false
}
produces: "application/json"
`
Movies : { ... }
```

The content is a boc body — the same syntax as the inside of `{ ... }` — without the braces. Multiple concerns are expressed as separate variables inside the same annotation. There is one annotation per definition.

---

## Annotations Are Boc Bodies

The content of an annotation is parsed and compiled as a boc. This means:

- Field declarations, nested boc literals, scalar values — all valid
- Referenced types are resolved at compile time — a typo like `Deribe` is a compile error
- Annotations can reference bocs defined elsewhere in the program

```
`macros: [Derive, JSON, some.package.GraphQL]`
Movies : { ... }
// error if Derive, JSON, or some.package.GraphQL do not exist
```

Annotations have restrictions that distinguish them from regular bocs:

- **Never executed.** The boc body is data. No method calls, no side effects, no invocations.
- **No nested annotations.** An annotation cannot contain another annotation.

Everything else a boc can contain is available: nested boc literals, arrays, scalar values, references to named types, string interpolation etc.

---

## `macros` — Triggering Macros

The reserved variable `macros` inside an annotation declares the macros the compiler should run on the annotated boc. It is always an array:

```
`macros: [Derive, Debug, Logging]`
Person : {
    name String
    age  Int
}
```

`!` is an alias for `macros` (defined as `! : macros`), so the shorthand form is:

```
`!: [Derive, Debug, Logging]`
Person : {
    name String
    age  Int
}
```

During parsing, the compiler scans annotations for `macros`. When found, the listed macros are scheduled to run during type inference — sequentially, in array order. The boc body carries no macro-triggering mechanism.

---

## How Macros Read Annotations

Inside a macro's `run` method, the parent boc's annotation is accessible via `self.annotation`. It is typed as `Boc` — the full shared annotation. Each macro accesses only the variable it owns, and the return type of that access is determined by the macro's `Schema` associated type:

```
GraphQL : {
    Schema : #(schema String)
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
        config      = self.annotation.json    // typed as Schema
        should_skip = config.ignore           // Bool
        ...
    }
}
```

Each macro reads the variable it cares about. Others are ignored. No string parsing is required — field access is direct and validated against `Schema` at compile time.

---

## Field-Level Annotations

Annotations follow the same rules on individual fields:

```
`macros: [GraphQL, JSON]`
Movies : {
    `graphql: { rename: "movieTitle" }`
    title String

    `json: { ignore: true }`
    internal_id String

    `
    json: { field_name: "release" }
    validate: { min: 1888 }
    `
    year Int
}
```

A macro reads field annotations via `self.fields[n].annotation`:

```
JSON : {
    Schema : #(field_name String, ignore Bool)
    run #(parent Boc, Boc) = {
        self.fields.forEach({ f Boc
            config = f.annotation.json         // typed as Schema
            config.ignore ? { /* skip this field */ }
        })
    }
}
```

`macros` is syntactically valid in a field-level annotation, but triggering per-field macros is a topic for a separate design.

---

## Simple Values

Not all annotations need nested boc structure. Scalar values are valid:

```
`port: 8080`
Server : { ... }

`produces: "application/json"`
Handler : { ... }

`json: "ignore"`
internal_id String
```

A macro reading `self.annotation.json` receives the string `"ignore"` directly. Whether to use a scalar or a nested boc is a choice for the macro's author to document.

---

## Documentation in Annotations

Documentation can live in an annotation as a regular variable:

```
`
documentation: "
    The Movies boc represents a film record as returned by the catalogue API.
    Fields map directly to the GraphQL schema unless annotated otherwise.
"
macros: [GraphQL, JSON]
graphql: { schema: "https://myapi.com/graphql" }
`
Movies : { ... }
```

A `Doc` macro could read `self.annotation.documentation` and generate API docs, IDE hover text, or reference pages.

Regular code comments (`/* ... */` and `//`) remain available for inline code explanation and are not accessible to macros. Documentation intended for tooling or generated output belongs in the annotation.

---

## Annotations in the `Boc` Metatype

Annotations are part of the `Boc` metatype. The compiler parses the annotation and attaches the resulting compiled boc to the `Boc` instance it creates for every definition:

```
Boc : {
    name        String
    instantiable Bool
    fields      [Boc]
    methods     [Boc]
    type_params [Boc]
    annotation  Boc      // the compiled annotation boc, or an empty boc if absent
    source      #()
}
```

`annotation` is a regular `Boc` value. It can be passed, iterated, serialized, and sent across a wire — the same as any other `Boc`. The `macros` field within it is what the compiler acts on; all other fields are data for macros and tooling to use freely.

See also: [Structural Reflection](Structural%20Reflection.md) for the full `Boc` API

---

## Summary of Rules

| Property              | Rule                                                                     |
| --------------------- | ------------------------------------------------------------------------ |
| Syntax                | Boc body without `{` `}`, backtick delimited                             |
| Companion form        | `name.info` file — same syntax, targets file-level boc                   |
| Placement             | Immediately before a boc or field definition                             |
| Per definition        | One annotation per definition                                            |
| Multiple concerns     | Separate variables inside the same annotation                            |
| Execution             | Never executed                                                           |
| String interpolation  | Not available — `${}` expressions unavailable inside annotations         |
| Nested annotations    | Not allowed                                                              |
| Type references       | Resolved at compile time — missing types are errors                      |
| `macros`              | Always an array; the only trigger for macros (`!` is an alias)           |
| Access in macros      | `self.annotation.<variable>` — return type determined by `Schema`        |
| In `Boc` metatype     | `Boc.annotation` — a regular compiled `Boc` value                       |
