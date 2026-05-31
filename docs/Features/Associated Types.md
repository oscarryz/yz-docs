#feature
# Associated Types

Associated types are type fields declared inside a boc interface and bound to concrete types by each implementation. They let a function work with an interface without knowing the concrete types in advance — the types are "associated" with the specific instance.

This is a direct application of [Path Dependent Types](Path%20Dependent%20Types.md). No separate mechanism is needed.

See also: [Generics — Type Parameters](Generics%20-%20Type%20Parameters.md), [Structural Typing](Structural%20typing.md)

---

## Declaring associated types

A boc interface declares type fields with implicit `#()` type:

```yz
Graph : {
    Node #()
    Edge #()
    add_node  #( Node )
    add_edge  #( Node, Node, Edge )
    neighbors #( Node, List(Node) )
}
```

`Node` and `Edge` are type fields — they hold type values. `add_node`, `add_edge`, and `neighbors` use them as parameter/return types.

---

## Binding associated types in a concrete boc

A concrete boc satisfies the interface by providing type aliases for each type field:

```yz
SocialGraph : {
    Node : User         // Node is bound to User
    Edge : Relationship // Edge is bound to Relationship

    add_node #( Node ) = {
        n Node
        // ...
    }
    add_edge #( Node, Node, Edge ) = {
        // ...
    }
    neighbors #( Node, List(Node) ) = {
        node Node
        [User("The end")]
    }
}
```

`Node : User` inside `SocialGraph` is a type alias: `SocialGraph.Node` is `User`.

---

## Using associated types in function signatures

A function that works with any `Graph` accesses the associated type through the instance:

```yz
process #( g Graph, n g.Node )
```

`g.Node` is a path-dependent type: its concrete value depends on which graph `g` holds. The compiler resolves it statically at the call site.

```yz
sg : SocialGraph()
u  : User("Alice")
process(sg, u)   // g.Node = SocialGraph.Node = User — verified at compile time
```

If you pass the wrong type for `n`:

```yz
process(sg, 42)  // error: Int does not satisfy SocialGraph.Node (User)
```

---

## Constrained associated types

A type field can carry a constraint — an interface the bound type must satisfy:

```yz
Graph : {
    Node #( label #(String) )   // Node must have a label() method
    first_node #( Node )
}
```

The constraint gives the compiler enough information to type-check method calls on `g.Node` values, even when `g`'s concrete type is not known:

```yz
describe #( g Graph, n g.Node ) {
    print(n.label())   // valid — Node is bounded by #(label #(String))
}
```

A concrete implementation must bind `Node` to a type that satisfies the constraint:

```yz
City : {
    name String
    label #(String) = { name }
}

CityGraph : { Node : City }   // City has label() — satisfies the bound
```

If the bound type does not satisfy the constraint, the compiler errors at the binding site.

---

## Abstract g: the bound is the type

Because Yz uses structural typing, `g.Node` when `g` has an abstract type (e.g. a local variable of type `Graph`) is structurally equivalent to its bound. Any value that satisfies the bound is a valid `g.Node`.

```yz
g Graph = CityGraph()
london : City(name: "London")   // City has label() — satisfies Node #(label #(String))
describe(g, london)             // VALID
```

There is no hidden nominal identity check. The concrete type of `g` at runtime does not matter — only whether the argument satisfies the declared bound.

If the argument does not satisfy the bound, the compiler errors:

```yz
Dot : { x Int }   // no label() method

dot : Dot(x: 3)
describe(g, dot)  // YZC-0079: Dot does not satisfy g.Node bound #(label #(String))
```

If `Node` has no bound (`Node #()`), any value is accepted — the type is fully unconstrained.

---

## Accessing the associated type as a value

Because `Node` is a stored field, it is readable:

```yz
sg : SocialGraph()
sg.Node   // evaluates to the type User
```

This allows functions to inspect or pass the associated type as a `#()` value.

---

## Multiple associated types

Interfaces can declare as many type fields as needed:

```yz
Container : {
    Item #()
    Cursor #()
    next   #( Cursor, Item )
    rewind #( Cursor )
}
```

Each implementation binds all of them:

```yz
IntStack : {
    Item   : Int
    Cursor : Int   // index into the stack
    next   #( Cursor, Item ) = { ... }
    rewind #( Cursor )       = { ... }
}
```

---

## Comparison with other languages

| Language | Associated type syntax | Yz equivalent |
|---|---|---|
| Rust | `trait Graph { type Node; }` | `Graph : { Node #(); ... }` |
| Scala | `trait Graph { type Node }` | same |
| Swift | `associatedtype Node` | same |
| Haskell | `type family Node g` | same |

Yz uses the same field-access syntax for associated types as for value fields. There is no special keyword. Unlike Rust or Scala, there is no notion of existential or opaque type identity — `g.Node` when `g` is abstract is simply its bound, consistent with structural typing throughout.
