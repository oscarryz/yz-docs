#solved in [Path Dependent Types](docs/Features/Path%20Dependent%20Types.md)

The questions raised here are resolved by the path-dependent types model:

- **Simple generics** — `T` is a field with implicit type `#()`. `List(Int)` passes `Int` as the `T` argument.
- **Associated types** — type fields inside a boc (`Node #()` in `Graph`) are accessed via path: `g.Node`.
- **Path-dependent link** — `process #(g Graph, n g.Node)` gives `n` the type stored in `g`'s `Node` field, resolved statically at compile time.
- **Structural vs nominal** — Yz uses structural typing throughout; no `implements` keyword.
- **Compilation strategy** — compile-time path resolution with specialized (monomorphized) Go output; no boxing.

The `#()` metatype (the type whose values are types) is the key: it is the implicit type of every bare type parameter, just as `Unit` is the implicit return type of bocs that return nothing.

See [Path Dependent Types](docs/Features/Path%20Dependent%20Types.md) for the full design.
