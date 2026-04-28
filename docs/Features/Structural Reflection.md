# Yz Structural Reflection — The `Boc` Type

## Overview

Every boc definition in Yz has a corresponding `Boc` instance created by the compiler.
This instance describes the boc's structure — its fields, methods, type parameters,
infostrings, and literal source — as a regular Yz value.

`Boc` is not a meta-type that appears by magic. It is a regular boc whose instances the
compiler populates. Once populated, a `Boc` instance behaves like any other value in the
language — it can be passed to functions, iterated, serialized, and sent across a wire.

This document describes the `Boc` type, how the compiler creates and exposes instances,
and how generics, compile-time execution, and `mix` are all built on top of it.

See also: [Generics](./yz-generics.md) · [Compile-Time boc](./yz-compile.md) ·
[mix](./yz-mix.md) · [boc Reference](./yz-boc.md)

---

## The `Boc` Type

```yz
Boc : {
    name         String
    instantiable Bool
    fields       [Boc]
    methods      [Boc]
    type_params  [Boc]
    infostrings  [String]
    source       #()
}
```

Every slot is a regular Yz value:

- `name` — the identifier the boc was defined with, or empty for anonymous bocs
- `instantiable` — whether the boc can be instantiated with uppercase convention
- `fields` — the boc's data slots, each described as a `Boc` instance
- `methods` — the boc's executable slots, each described as a `Boc` instance
- `type_params` — placeholder `Boc` instances representing unresolved type parameters
- `infostrings` — the raw infostring values attached to this boc or field
- `source` — the literal boc body as an executable value — homoiconicity lives here

Fields and methods are both `[Boc]` because in Yz everything is a boc. The distinction
between a field and a method is a matter of what `source` contains, not a difference in
kind.

---

## How The Compiler Populates `Boc` Instances

The compiler creates a `Boc` instance for every boc definition it encounters during
parsing and inference. This happens automatically — the developer does not request it.

The instance is made available during `Compile` execution via the `compiler` boc:

```yz
compiler.read(Named)    // returns the Boc instance for Named
compiler.insert(b)      // splices a Boc into the current compilation context
```

`compiler` is available only during `Compile` slot execution. Outside of compile-time
context, `compiler` does not exist. This is the explicit and localized seam between the
language and the compilation phase — there is no other.

See also: [Compile-Time boc — Compilation Lifecycle](./yz-compile.md#compilation-lifecycle)

---

## Type Parameters Are Placeholder Bocs

A type parameter such as `T` is not a special language primitive. It is a `Boc` instance
with no fields, no methods, and no concrete definition — a placeholder waiting to be
filled at instantiation time.

```yz
Box : {
    T        // a Boc with no definition — instantiable: true, everything else empty
    value T  // a field whose type slot points to T's Boc instance
}
```

When `Box(String)` is written, the compiler substitutes `String`'s `Boc` instance for
`T`'s placeholder `Boc` everywhere in `Box`. Generics are `Boc` substitution.

The single uppercase letter convention is a human readability signal — a way for
developers to identify placeholder bocs at a glance. The compiler identifies them by
their emptiness, not their name.

### Constraint Inference Through `Boc`

Constraint inference is the compiler querying `Boc` instances. When the compiler
encounters `thing.talk()` inside a generic boc, it asks:

```yz
compiler.read(T).methods.any({ m Boc
    m.name == "talk"
})
```

If the placeholder `Boc` for `T` does not yet have a `talk` method, the compiler records
the requirement. When `T` is eventually filled with a concrete `Boc`, the compiler
verifies the concrete `Boc`'s methods satisfy all recorded requirements.

There is no separate constraint machinery. Constraint inference is `Boc` queries.

See also: [Generics — Constraint Inference](./yz-generics.md#constraint-inference)

---

## Compile-Time Execution Through `Boc`

A `Compile` slot receives its parent boc as a `Boc` instance. Previously this appeared
as a magic `self` — through the `Boc` type it is simply the named argument passed to
`run`:

```yz
Compile : {
    run #(Boc, Boc)  // receives parent as Boc instance, returns Boc to merge
}
```

Inside any `Compile` implementation, the parent boc's full structure is available as
regular `Boc` slots:

```yz
Serialize : {
    run #(Boc, Boc) = {
        parent = compiler.read(self)

        generated_fields = parent.fields.map({ f Boc
            "{f.name}: {f.source().serialize()}"
        }).join(",")

        {
            serialize #(String) = { generated_fields }
        }
    }
}
```

`parent.fields` is `[Boc]`. Each field is a `Boc`. Everything is regular iteration over
regular values. The returned boc literal is merged into the parent by the compiler via
`compiler.insert`.

See also: [Compile-Time boc](./yz-compile.md)

---

## `mix` As An Example

`mix` demonstrates how `Boc` composes with `Compile` to implement a feature that could
have been a language primitive — without any new primitives.

```yz
mix : {
    embeds [#()]
    e Compile = Embed(embeds)
}

Embed : {
    targets [#()]
    run #(Boc, Boc) = {
        merged = targets.map({ t #()
            compiler.read(t)
        }).reduce({ acc Boc  b Boc
            // merge fields and methods from b into acc
            acc.fields  + b.fields
            acc.methods + b.methods
        })
        merged
    }
}
```

Usage:

```yz
Person : {
    mix([Named, Auditable, Logger])
    last_name String
}
```

`Named`, `Auditable`, and `Logger` are passed as values — the same way `String` is
passed to `Box(String)`. `compiler.read()` returns a `Boc` instance for each.
`Embed` merges their fields and methods. The result is returned as a `Boc` and spliced
into `Person` by the compiler.

`mix` is not a keyword. It is a boc that uses `Compile` which uses `Boc`. There is no
magic.

### Conflict Detection

When two mixed bocs contribute a slot with the same name, the compiler detects the
conflict after all merges are complete. The error message identifies the origin `Boc` of
each conflicting slot — whether it came from a `mix` target or a `Compile` generated
slot.

See also: [mix](./yz-mix.md)

---

## Infostrings Through `Boc`

Infostrings are not a separate metadata system. They are a slot on `Boc`:

```yz
"json:user_name"
name String
```

The compiler attaches `"json:user_name"` to the `Boc` instance for `name` via its
`infostrings` slot. A `Compile` implementation reads it like any other slot:

```yz
parent.fields.forEach({ f Boc
    f.infostrings.forEach({ s String
        // parse and act on field-level metadata
    })
})
```

The distinction between structure and metadata dissolves — both are slots on `Boc`.

---

## `Boc` As A Regular Value

Because `Boc` is a regular boc instance, it participates in everything Yz values can do:

```yz
// pass to a function
process(compiler.read(Person))

// serialize and send across a wire
network.send(compiler.read(Person).serialize())

// reconstruct on the other side
person_meta Boc = Boc.deserialize(network.receive())
person_meta.fields  // fully accessible remotely
```

This gives Yz structural reflection without a dedicated reflection API. The same `Boc`
type used for compile-time code generation is the type used for runtime introspection and
inter-process communication.

---

## The Complete Magic Surface

The entirety of what is special about `Boc` and the compiler seam is:

1. **The compiler creates a `Boc` instance for every boc definition** — automatically,
   without developer action
2. **`compiler.read(b)`** returns the `Boc` instance for any boc — available only during
   `Compile` execution
3. **`compiler.insert(b)`** splices a `Boc` into the current compilation context —
   available only during `Compile` execution

Everything else — generics, constraint inference, `mix`, `Compile`, infostrings,
inter-process reflection — is regular Yz code operating on regular `Boc` values.