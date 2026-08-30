#spec
# 12. Macros

> Status: first slice (YZC-0028). This chapter is a stub covering what the
> compiler implements today; the full design is in
> [Macros](../docs/Features/Macros.md).

## 12.1 The Macro Interface

A macro is a boc satisfying the `Macro` structural interface:

```
Macro : {
    Schema #()
    run    #(subject Boc, config Schema, Boc)
}
```

`Schema` fixes the shape of the config block the macro accepts. In the
current slice it must be a type alias to a named config type declared in the
same package (`Schema : NoConfig`); config fields are limited to the scalar
types `String`, `Int`, `Decimal`, and `Bool`.

## 12.2 Triggering

Uppercase-keyed entries in a boc's annotation trigger macros. The value is
the config block, validated against the macro's `Schema`; `{}` means no
config:

```
`Debug: {}`
Person : {
    name String
    age  Int
}
```

Multiple triggers run sequentially in annotation order; each macro receives
the progressively merged subject. Lowercase entries are passive metadata and
never trigger. An uppercase key with no matching macro is a compile error.

## 12.3 Two-Phase Build

1. **Bootstrap** — directories defining macros are excluded from the app
   build and compiled to native executables (`target/macros/<pkg>/bin/`),
   with a compiler-provided prelude declaring `Boc`, `Field`, `NoConfig`,
   `generated`, and `while`. Executables are cached on a source hash.
2. **Expansion** — before semantic analysis, the compiler runs each
   triggered macro as a subprocess, passing the subject and config, and
   merges the returned slots into the subject boc. Runs are cached per
   subject structure.

## 12.4 Wire Format

Both directions use Yz source restricted to the non-executable data subset.

Compiler → macro (stdin):

```
subject: {
    name: "Person"
    fields: [
        {
            name: "name"
            type: "String"
        }
    ]
}
config: { pretty: true }
```

Macro → compiler (stdout): a raw Yz boc body containing the generated
slots. The `run` method returns it via `generated(src)`; the compiler parses
it and appends the elements to the subject boc literal.

## 12.5 Reflection Surface (this slice)

```
Boc : {
    name   String     // subject type name
    fields [Field]    // declared data fields (name Type form)
    source String     // generated source (outbound)
}

Field : {
    name String
    type String       // type as written in source, unresolved
}
```

## 12.6 Rules and Errors

- Macros must live in a separate package (directory) from the bocs they
  process; violations are a compile error.
- Macros cannot be defined in the project root package.
- Two macros with the same name in one project are a compile error.
- A macro cannot (transitively) trigger a macro of a package that is
  currently being compiled — the cycle is reported with its full chain.
- A macro package must not declare a `main` boc.

## 12.7 Deferred

Full structural reflection (methods, type params, field annotations),
inline `Schema #(...)` config typing, expansion interleaved with type
inference, bare-name triggers, `name.info` companion triggers, and
non-scalar config values. See the YZC-0028 follow-ups in the
implementation tracker.
