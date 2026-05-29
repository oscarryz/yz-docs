#open-question

How can `.?()` — a no-argument call to the `?` method on `Result` — propagate errors non-locally to the enclosing named boc?

See: [return, break, and continue](../Features/return,%20break,%20continue.md)

---

The `?` method on `Result` is defined as:

```yz
? #( handler #(E, V) = { e; return e }, V )
```

When called with an explicit boc at the call site, non-local return already works under the existing rule — anonymous bocs passed as arguments are transparent to `return`, so `return` exits the nearest named boc at the *definition site* (the caller):

```yz
save_report : {
    data Data
    fs.save(data).?({ e FsError; return e })  // return exits save_report ✓
    log.info("saved")
}
```

The problem is the no-argument form `.?()`. The default boc `{ e; return e }` would be defined in the stdlib, so its `return` would target a named boc there — not the call site.

```yz
save_report : {
    data Data
    fs.save(data).?()  // where does return exit? ✗
    log.info("saved")
}
```

The convention would be to always name it `.?()` — mirroring how Rust's `?` reads, but as a regular method call.

---

## Options

**A — Default boc arguments inherit the call site context**

Modify the semantics so that a default boc argument behaves as if it were written at the call site. Its `return` would exit the caller's nearest named boc, not the definition site's.

This is a targeted relaxation of the current rule and would make `.?()` work with no other changes.

**B — A "transparent" modifier on the default boc**

Allow a boc parameter to be marked so that its default is evaluated in the caller's context:

```yz
? #( handler^ #(E, V) = { e; return e }, V )
```

More explicit but adds syntax.

**C — Defer — require the explicit form for now**

`.?({ e; return e })` works today. Defer the shorthand until the semantics are clearer.

**D — Make `?` a postfix operator eventually**

Accept that the true zero-noise form `result?` requires a dedicated operator and design it properly when the time comes.

**E — Infostring macro (`Nlr`) with synthesized call-site rewrite**

Use the existing infostring/compile-time mechanism. The `Nlr` compile implementation scans the annotated boc body, finds every `.?` or `.?()` occurrence, and rewrites each one to `.?({ e E; return e })` with `E` inferred from the `Result(T, E)` type at that specific call site. The injected boc is anonymous and lands at the call site, so `return` exits correctly under the existing non-local return rules — no new language semantics required.

Roughly equivalent to how Rust's `?` started as the `try!()` macro before becoming a language operator.

*Per-boc* — annotate once, use `.?` freely inside:

```yz
`compile_time: [Nlr]`
save_report: {
    config: toml.parse(...).?    // rewritten to .?({ e EncodingError; return e })
    conn:   db.connect(...).?    // rewritten to .?({ e DbError; return e })
    log.info("done")
}
```

*Per-declaration* — surgical, without annotating the whole boc:

```yz
`nlr`
config: toml.parse(...)
`nlr`
conn: db.connect(...)
```

The explicit form remains valid anywhere and composes freely — for cases that need more than propagation:

```yz
conn: db.connect(...).?({ e DbError; log.warn(e.message); return e })
```

The key mechanism: `Nlr` does not change the semantics of `return` or `.?`. It is purely a source rewrite — generating code the developer could have written by hand. A reader who knows `Nlr` is active understands exactly what `.?` expands to.

Tradeoff: the transformation is invisible without knowing `Nlr` is in scope. The per-boc form in particular is opt-in implicit propagation, which in practice moves closer to exceptions — the happy path reads cleanly but error paths are hidden from the call site.
