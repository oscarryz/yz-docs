# 12 — Macros

## Overview

A **macro** is a compile-time boc that receives an annotated boc as input and returns a new boc body that is merged into the source before semantic analysis of the target package.

Macros live in a dedicated `macros/` package (a subdirectory named `macros` within the project). The compiler builds macro packages first, producing a native binary, then invokes it once per trigger before compiling the main project.

## Trigger form

An annotation backtick-block immediately before an uppercase boc declaration triggers one or more macros:

```yz
`Debug: {}`
Person: {
    name String
    age  Int
}
```

Each uppercase key in the annotation body names a macro. The value (a boc literal) carries the config for that macro.

## Macro interface

A boc satisfies the `Macro` contract when it declares:

```yz
Schema: NoConfig          // or a custom config schema
run #(subject Boc, config NoConfig, Boc) { ... }
```

`NoConfig` — the config schema type when the macro takes no configuration.

`Boc` — the compiler-provided metatype:

```yz
Boc: {
    name   String
    fields [Field]
    source String
}
Field: {
    name  String
    type_ String
}
```

`generated #(src String, Boc)` — helper that wraps a raw Yz source string as the return value.

## Wire format

The compiler serializes the subject boc as a newline-delimited text payload on the macro binary's stdin:

```
<subjectName>
<fieldName> <fieldType>
...
---config---
<key>=<value>
...
```

The macro binary writes the raw Yz boc body (statements only, no surrounding `{}`) to stdout.

## Two-phase build

1. **Macro compilation** — each `macros/` directory is compiled with the macro prelude prepended. The result is a native Go binary cached under `target/macros/<pkgKey>/bin/macros`.
2. **Expansion** — before sema, annotated bocs invoke the macro binary. The returned source is parsed and appended to the subject's boc literal elements.
3. **Main build** — the expanded AST is processed normally.

## Macro prelude

The compiler injects the following Yz source into every macro package before compilation so macro authors can reference `Boc`, `Field`, `NoConfig`, and `generated()` without imports:

```yz
Boc: {
    name   String
    fields [Field]
    source String
}
Field: {
    name  String
    type_ String
}
NoConfig: {}
generated #(src String, Boc) {
    b: Boc(name: "", fields: [Field](), source: src)
    b
}
```

## Restrictions (first slice)

- Macros must live in a `macros/` subdirectory; macros in the project root package are rejected.
- Only `NoConfig` schema is supported; custom config schemas are deferred.
- Macro-on-macro (a macro triggering another macro) is not supported in this slice.
- Same-package macros (a macro applied to a boc in its own package) are rejected.

## Example

```yz
// macros/debug.yz
Debug: {
    Schema: NoConfig
    run #(subject Boc, config NoConfig, Boc) {
        body: subject.name + "("
        sep:  ""
        subject.fields.each({ f Field
            body = body + sep + f.name + ": $" + "{" + f.name + "}"
            sep  = ", "
        })
        generated("debug #(String) {\n    \"" + body + ")\"\n}")
    }
}

// main.yz
`Debug: {}`
Person: {
    name String
    age  Int
}

main: {
    p: Person(name: "Ada", age: 36)
    print(p.debug())
}
// Output: Person(name: Ada, age: 36)
```
