#open-question

## How should Yz handle existential types over associated types?

Related features: [Associated Types](docs/Features/Associated%20Types.md), YZC-0066, YZC-0074.

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

### The "opaque token" exception (advanced existentialism)

Even when `g.Node` is unknown, a value *produced by `g`* is guaranteed to have type `g.Node`.
If the value never leaves the path, the compiler can allow it to be passed back to operations on
the same `g`:

```js
Graph: {
  Node #()
  name      String
  firstNode #(Graph.Node)   // returns this graph's own Node type
}

processElement #(g Graph, Void) {
    // token : g.Node — opaque, but a real value
    token: g.firstNode()

    // ALLOWED — token and g share the same existential path; types align.
    sameToken: visit(g, token)

    // FORBIDDEN — token is g.Node, not some other Graph's Node.
    // visit(otherGraph, token)
}
```

The key invariant: operations on an existential value are safe only when all path-dependent
arguments were produced by the *same* root binding.

---

### Design crossroads

**Implicit vs. explicit wildcard syntax**

When a user writes `myGraphs Array Graph`, should the compiler silently make `Node` existential,
or should the user be required to annotate it explicitly (e.g. `Array Graph{Node: #}` or
`Array Graph{Node: ?}`) so the loss of `visit()` is visible at the declaration site?

Implicit is ergonomic but may surprise users when `visit()` stops working.
Explicit is verbose but documents the intent.

**Path identity tracking**

The opaque-token rule requires the compiler to track that `token` was produced by `g` and not
some other binding. This may need a notion of *path variables* (named existential witnesses) in
the type system, similar to Scala's path-dependent types or Haskell's `ST s` trick.

---

### Open questions

1. **Implicit or explicit wildcard?** Should `Array Graph` silently erase `Node`, or should the
   user write something like `Array Graph{Node: #}` to opt in to existential behaviour?

2. **Syntax for the wildcard.** If explicit: `#`, `?`, `_`, or something else? Should it be the
   same token used in other contexts (e.g. anonymous param `#()`)?

3. **Opaque-token scope.** Can an opaque `g.Node` value be stored in a field and used later, or
   only within the same block where `g` is in scope?

4. **Error message quality.** When `visit(myGraphs[0], london)` is rejected, what message should
   the compiler produce? Something like _"Node is existential here; cannot match against City"_?

5. **Interaction with constraints (YZC-0074).** If `Node #(name #(String))`, the erased type
   still satisfies the constraint. Should the compiler allow `g.firstNode().name()` even on an
   existential `g.Node`?

6. **Collections inference.** When building `[cg, sg]`, at what point does the compiler decide
   to generalise? Is it at the array literal, or deferred until the binding `myGraphs`?
