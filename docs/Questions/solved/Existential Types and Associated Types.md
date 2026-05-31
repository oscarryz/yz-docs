#open-question #resolved

## How should Yz handle existential types over associated types?

Related features: [Associated Types](docs/Features/Associated%20Types.md), YZC-0066, YZC-0074.

**Resolution (YZC-0079):** In Yz's structural type system, `g.Node` when `g` has an abstract type is simply equivalent to its bound. There are no existential types in the nominal sense. See the resolution section below.

---

### The problem

When a boc type has an associated type (e.g. `Node #()` inside `Graph`), and you put two different
concrete graphs into the same array, the compiler must generalise the element type to `Graph`. At
that point `Node` becomes *existential*: it definitely exists, but its concrete identity is hidden
(or unknown). The type checker must know what is safe to do with `g.Node` when `g : Graph` and
not `CityGraph`.

---

### Concrete scenario

```js
Graph: {
  Node #()
  name String
}

CityGraph: {
  Node: City
  name: "London Transit Graph"
}

SocialGraph: {
  Node: Person
  name: "Six Degrees of Alice"
}

main: {
    cg: CityGraph()
    sg: SocialGraph()

    // Compiler generalises element type to `Graph`; Node becomes existential.
    myGraphs: [cg, sg]

    // ALLOWED — name is a plain field on Graph, no Node involved.
    for g in myGraphs {
        print(g.name)
    }

    // ERROR — compiler cannot guarantee myGraphs[0].Node is City.
    london: City(name: "London", population: 9000000)
    visit(myGraphs[0], london)
}
```

---

### Resolution

**The existential framing was wrong for Yz.** The question assumed nominal type identity — that `City` and `Person` are distinct even if structurally identical, and that `g.Node` is an opaque token tied to a specific `g`. Yz uses full structural typing, which dissolves the problem:

> `g.Node` when `g` has an abstract type is structurally equivalent to its bound.

Concretely:

- `Node #(label #(String))` — `g.Node` ≡ "any value with a `label()` method returning `String`"
- `Node #()` — `g.Node` ≡ "any value" (unconstrained)

So `visit(g, london)` is **valid** as long as `london` satisfies the `Node` bound, regardless of whether `g` is concrete or abstract. The compiler checks the bound, not path identity.

```yz
Graph: { Node #(label #(String)) }

City: {
    name String
    label #(String) = { name }
}

describe #(g Graph, n g.Node) { print(n.label()) }

main: {
    g Graph = CityGraph()
    london: City(name: "London")
    describe(g, london)   // VALID — london satisfies Node #(label #(String))
}
```

Error only when the argument does not satisfy the bound:

```yz
City: { name String }  // no label() method

describe(g, london)    // YZC-0079 error: City does not satisfy g.Node bound
```

---

### The "opaque token" question (YZC-0076)

The opaque-token scenario from the original question:

```js
processElement #(g Graph) {
    token: g.firstNode()
    visit(g, token)        // allowed
    visit(otherGraph, token) // forbidden?
}
```

In structural typing, `token` has the type of `firstNode`'s return — which is the `Node` bound. Passing it to `visit(otherGraph, token)` is valid as long as `token` satisfies `otherGraph`'s `Node` bound too (which it does, since they share the same bound `#(label #(String))`). Path-identity tracking is not needed.

YZC-0076 is therefore lower priority and would only become relevant if Yz adds nominal type distinctions in the future.

---

### Open questions — status

1. **Implicit or explicit wildcard?** — Resolved: implicit. No syntax needed; bound check is sufficient.
2. **Syntax for the wildcard** — Resolved: no syntax needed.
3. **Opaque-token scope** — Resolved: not needed in a structural type system (see above).
4. **Error message quality** — Resolved: `YZC-0079: argument type X does not satisfy g.Node bound`.
5. **Interaction with constraints** — Resolved: the bound is the type; `g.firstNode().name()` is always valid if `Node` is bounded by `#(name #(String))`.
6. **Collections inference** — Still relevant for array literals; resolved at the literal site when elements unify to a common interface type.
