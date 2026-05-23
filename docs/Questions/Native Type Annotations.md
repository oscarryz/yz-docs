#open-question

## How should Yz types declare a native (Go) implementation?

Related tickets: YZC-0025 (infostrings / compile-time bocs), YZC-0058 (Native annotation implementation).

---

### The problem

Yz types like `Int`, `String`, `Bool`, and `Decimal` are ultimately backed by Go primitives. Once
these move to Yz source (YZC-0031), each method needs to tell the compiler what Go code to emit for that method — otherwise the compiler cannot produce working Go output.

The same mechanism would let users wrap arbitrary Go libraries (e.g. the `net/http` package) as Yz types, replacing the current hand-coded `http` singleton in the runtime.

---

### The two-fold nature

**Layer 1 — compile-time bocs (YZC-0025)**

The general mechanism: an infostring on a boc triggers a user-supplied `Compile` implementation
that transforms or extends the boc at build time. This is a general macro system.

See: [Compile Time Bocs](docs/Features/Compile%20Time%20Bocs.md) [Info strings](docs/Features/Info%20strings.md) [Compile Bocs Issues](docs/Questions/Compile%20Bocs%20Issues.md)

**Layer 2 — Native annotation (YZC-0058)**

A *specific* compile-time instance where the compiler itself provides the transformation, without
requiring a user-written `Compile` boc. This is needed because the `Compile` mechanism itself
depends on types like `Int` and `String` already existing — a bootstrapping problem.

The compiler must therefore detect a magic annotation key (e.g. `compile_time:[Native]` or
`!:[Native]`) and handle it internally, bypassing the `Compile` interface entirely.

---

### Proposed syntax (draft — subject to change)

```js
`compile_time: [Native]
 go_type: "int"
 go_zero: "0"`
Int : {
    `go: "self + other"`
    + #(other Int, Int)

    `go: "self - other"`
    - #(other Int, Int)

    `go: "self % other"`
    % #(other Int, Int)

    `go: "self == other"`
    == #(other Int, Bool)
}
```

The type-level infostring carries `compile_time:[Native]` plus any type-wide metadata (`go_type`,
`go_zero`). Each method has its own single infostring with the `go:` expression template; `self`
refers to the receiver, params are referenced by name.

For methods that require a Go import (both keys in one multiline infostring):

```js
`compile_time: [Native]`
Http : {
    `go_import: "net/http"
     go: "http.Get(url)"`
    get #(url String, String)
}
```

---

### Open questions

1. **Annotation key** — should it be `go:` (backend-explicit) or `native:` (abstract)?
   `go:` is simpler now; `native:` avoids binding source to a single backend but requires
   a mapping layer when a second backend is added.

2. **Template variables** — `self` and param names? Or positional (`{{0}}`, `{{1}}`)?
   Named is more readable; positional avoids needing to know param names for anonymous params.

3. **Go imports** — declared per-method (`go_import:`) or collected at the type level?
   Per-method is more self-contained; type-level avoids repetition for types with many methods.

4. **Return type mapping** — who maps the Go return value back to the Yz type?
   The compiler wraps the Go expression in the appropriate `std.New*()` constructor. The
   annotation only describes the raw Go computation, not the boxing.

5. **Error returns** — Go methods often return `(T, error)`. How does the native annotation
   express this? Does it map to `Result(T, Error)` automatically?

6. **Bootstrapping boundary** — which types are "primordially native" (handled with no
   annotation at all, purely compiler magic) vs. annotated in source? Candidates for primordial:
   `Int`, `Bool`, `String`, `Decimal`, `Unit`. Everything else could use the annotation.

7. **User-land Go interop** — should regular Yz code be able to use `compile_time:[Native]`
   to wrap any Go package, or is it restricted to stdlib? Unrestricted is powerful but opens
   the compiler to arbitrary Go; restricted keeps the language portable.
