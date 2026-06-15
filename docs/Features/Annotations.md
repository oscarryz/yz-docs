#feature

# Yz Annotations

## Overview

An annotation is a boc body placed immediately before a boc or field definition. It is parsed and compiled as part of the program, but never executed. It carries structured metadata attached to a definition — readable by macros, build tools, and other compile-time consumers.

```yz
`
author: "oscar"
version: "1.0"
`
Movies : {
    ...
}
```

The backtick delimiters take the place of `{` and `}`. Everything inside is valid Yz.

See also: [Macros](Macros.md) · [Structural Reflection](Structural%20Reflection.md) · [Boc Type](Boc%20Interface.md)

---

## Syntax

There are two forms:

**Form 1 — Inline backtick literal** placed immediately before the definition it annotates:

```yz
`port: 8080`
Server : { ... }
```

**Form 2 — Companion file** (`name.info` in the same directory as the boc named `name`):

```
net/
  http.info   ← annotation for net.http
  http.yz     ← defines net.http
```

The content of `http.info` is the same annotation body syntax — a boc body without `{ }`. A `name.info` file is never a boc; the compiler attaches its content to the named boc's annotation slot.

Use Form 1 for inline annotations on individual declarations. Use Form 2 when the annotation targets the boc defined by the file itself, or when keeping metadata separate from code — for example, project dependency declarations.

---

## Content

The content of an annotation is a boc body — the same syntax as the inside of `{ ... }`, without the braces. Scalar values, nested bocs, and arrays are all valid:

```
`port: 8080`
Server : { ... }

`produces: "application/json"`
Handler : { ... }

`
author: "oscar"
tags:   ["http", "rest"]
config: {
    timeout: 30
    retries: 3
}
`
Client : { ... }
```

Type references are resolved at compile time — a missing type is a compile error. Multiple concerns live as separate entries inside a single annotation. There is one annotation per definition.

Annotations have two restrictions that distinguish them from regular bocs:

- **Never executed.** The boc body is data — no method calls, no side effects, no invocations.
- **No nested annotations.** An annotation cannot contain another annotation.

---

## Uppercase vs Lowercase Names

Inside an annotation, the case of a name carries meaning:

- **Uppercase names** refer to types resolved at compile time. They are the mechanism by which macros are triggered. A missing type is a compile error.
- **Lowercase names** are plain metadata fields — readable by any compile-time consumer, never dispatch triggers.

```
`
Debug
JSON: { ignore: false }
author: "oscar"
`
Movies : { ... }
```

`Debug` and `JSON` are type references; `author` is passive metadata. How uppercase names trigger macros is described in [Macros](Macros.md).

---

## Field-Level Annotations

Annotations can be placed on individual fields inside a boc using the same syntax:

```
Movies : {
    `json: { field_name: "movie_title" }`
    title String

    `json: { ignore: true }`
    internal_id String

    `
    json:     { field_name: "release" }
    validate: { min: 1888 }
    `
    year Int
}
```

A field annotation is attached to that field's `Boc` instance and is accessible to any compile-time consumer that walks the boc's fields.

---

## Annotations in the `Boc` Metatype

Every definition has exactly one annotation slot. The compiler parses the annotation and attaches the resulting compiled boc to the `Boc` instance of the definition:

```
Boc : {
    name        String
    instantiable Bool
    fields      [Boc]
    methods     [Boc]
    type_params [Boc]
    annotation  Boc      // compiled annotation boc, or empty boc if absent
    source      #()
}
```

`annotation` is a regular `Boc` value. It can be passed, iterated, and serialized the same as any other `Boc`.

See also: [Structural Reflection](Structural%20Reflection.md) for the full `Boc` API

---

## Summary of Rules

| Property             | Rule                                                                          |
| -------------------- | ----------------------------------------------------------------------------- |
| Syntax               | Boc body without `{` `}`, backtick delimited                                 |
| Companion form       | `name.info` file — same syntax, targets file-level boc                        |
| Placement            | Immediately before a boc or field definition                                  |
| Per definition       | One annotation per definition                                                 |
| Execution            | Never executed                                                                |
| String interpolation | Not available inside annotations                                              |
| Nested annotations   | Not allowed                                                                   |
| Uppercase names      | Type references resolved at compile time — missing types are errors           |
| Lowercase names      | Passive metadata — readable by any compile-time consumer                      |
| In `Boc` metatype   | `Boc.annotation` — a regular compiled `Boc` value                            |
