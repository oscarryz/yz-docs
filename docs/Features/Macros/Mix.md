#feature
# Mix — Compile-Time Code Reuse

## Purpose

`Mix` copies methods from one boc into another at compile time. It is not inheritance — it is a one-time structural copy. The result is as if you had written the methods by hand.

## Design

```yz
Animal: {
    talk #(String) { "hello" }
}

`
macros: [Mix]
mix: Animal
`
Cat: {
    // Cat now has talk #(String) as if it were written here
}
```

`Mix` reads the `mix:` field from the annotation, finds the named boc, and copies its methods into the annotated boc's body during compilation.

## Properties

- **Not inheritance**: `Cat` does not have an `is-a` relationship with `Animal`. `Cat` satisfies any interface `Animal` satisfies, but only because it has the same methods — not because of a runtime relationship.
- **Compile-time only**: No runtime overhead. The copied methods are indistinguishable from hand-written ones.
- **Conflicts**: If `Cat` already defines `talk`, it is a compile error (or last-wins — TBD).
- **Chaining**: `Mix` output participates in further type inference normally.

## Status

Design phase. Depends on YZC-0025, YZC-0028. See also `docs/Examples/Yz - Compile Time Mix.md`.
