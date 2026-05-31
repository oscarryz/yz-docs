#feature
# Path-Dependent Types

Path-dependent types (PDT) are how Yz handles generics, type aliases, and associated types — all through a single, uniform mechanism. There is no separate generics syntax, no special `<T>` brackets, no `type` keyword for aliases, and no `associated type` declaration form. Everything collapses into the same "types are values, fields hold them" model that drives the rest of the language.

See also: [Generics — Type Parameters](Generics%20-%20Type%20Parameters.md), [Type Alias](Type%20Alias.md), [Structural Typing](Structural%20typing.md)

---

## The metatype: `#()`

Every UpperCase boc — every type — satisfies the empty interface `#()`. This makes `#()` the natural metatype: the type whose values are types.

```yz
Int    // satisfies #()
String // satisfies #()
Person // satisfies #()
```

Just as `Unit` (the return type of bocs that produce no useful value) is invisible in syntax, `#()` as the type of a type parameter is invisible. You never write it; the compiler infers it.

---

## Type parameters as fields

A bare single-uppercase-letter line inside a boc body declares a type parameter. Semantically it is a field with implicit type `#()`:

```yz
List : {
    T             // field of type #() — holds the element type
    add    #( T )
    remove #( T )
    size   #( Int )
}
```

Constructing a `List` passes a type as the first argument, exactly like any other field:

```yz
intList : List(Int)       // T = Int
stringList : List(String) // T = String
```

The stored type is accessible as a field:

```yz
intList.T   // evaluates to the type Int
```

This is the path: `intList.T`. It is path-*dependent* because `intList.T` and `stringList.T` resolve to different types depending on which instance you use.

---

## Three levels, one mechanism

### Level 1 — Type alias

```yz
Bar : Foo
```

`Bar` is a new name for `Foo`. Expanded form:

```yz
Bar #( name String ) = Foo   // signature copied from Foo; body is Foo
```

The compiler copies `Foo`'s signature into `Bar`'s, then emits `type Bar = Foo` in Go. No `T` field involved; the simplest case.

### Level 2 — Generic instantiation

```yz
StringList : List(String)
```

`StringList` is `List` with `T` bound to `String`. `StringList.T` is `String`. Expanded form:

```yz
StringList #( add #(String), remove #(String), size #(Int) ) = List(String)
```

`T` is gone in the expanded form because it has been substituted. The signature spells out the fully concrete interface.

### Level 3 — Associated types (via path-dependent access)

A boc can declare type fields that its callers use as associated types:

```yz
Graph : {
    Node #()
    Edge #()
    add_node  #( Node )
    add_edge  #( Node, Node, Edge )
    neighbors #( Node, List(Node) )
}
```

A concrete graph binds those types:

```yz
SocialGraph : {
    Node : User
    Edge : Relationship
    add_node  #( Node ) = { ... }
    add_edge  #( Node, Node, Edge ) = { ... }
    neighbors #( Node, List(Node) ) = { ... }
}
```

`Node : User` inside `SocialGraph` is a type alias (level 1) used as an associated type. A function that works with any graph accesses the associated type through the instance path:

```yz
process #( g Graph, n g.Node )
```

`n`'s type is `g.Node`. At the call site:

```yz
sg : SocialGraph()
u  : User("Alice")
process(sg, u)   // g.Node = SocialGraph.Node = User — compiler verifies u is User
```

The compiler resolves `g.Node` statically from the known type of `sg`; no runtime lookup is needed.

#### Constrained associated types

A type field can be constrained by an interface, giving the compiler enough information to type-check method calls even when the concrete graph is not known:

```yz
Graph : {
    Node #( label #(String) )
}

describe #( g Graph, n g.Node ) {
    print(n.label())   // valid — the bound guarantees label() exists
}
```

#### When g has an abstract type

Because Yz uses structural typing, `g.Node` when `g`'s static type is the interface (not a concrete implementation) is structurally equivalent to its bound. There is no nominal existential identity — the bound is the whole story:

```yz
g Graph = SocialGraph()
u : User("Alice")        // User satisfies Node #(label #(String))
describe(g, u)           // valid — u satisfies the bound
```

The compiler errors only if the argument does not satisfy the bound:

```yz
describe(g, 42)   // error: Int does not satisfy g.Node bound
```

If `Node` has no bound (`Node #()`), any value is accepted.

---

## Type variables in signatures

When a function works with multiple related generic types, single-uppercase letters (`A`, `B`, `T`, etc.) serve as type variables in the signature. Their `#()` type is implicit:

```yz
map #( collection List(A), fn #(A, B), List(B) )
```

`A` and `B` are unbound type variables. At the call site the compiler unifies `List(A)` with the type of the first argument to infer `A`, then infers `B` from `fn`'s return type.

---

## Constraints

A constraint is written by replacing the implicit `#()` with a named interface:

```yz
SortedList : {
    T Comparable      // T must satisfy Comparable; #() is implicit behind Comparable
    add #( T )
    ...
}

serialize_all #( collection List(A Serializable), String )
```

`T Comparable` means T's type is `Comparable` rather than the unconstrained `#()`. No extra syntax is needed.

---

## Compile-time resolution

`g.Node` in a signature is resolved at compile time from the static type of `g`:

- **Concrete `g`** (`g : SocialGraph`): `g.Node` resolves to the concrete bound type (`User`). The compiler verifies the argument is exactly `User`.
- **Abstract `g`** (`g : Graph`): `g.Node` resolves to its bound interface (e.g. `#(label #(String))`). The compiler verifies the argument satisfies that bound structurally.

In both cases the check is purely static. The generated Go output is fully specialized — no `interface{}`, no boxing, no runtime type descriptors.

Unlike Scala's path-dependent types, Yz does not track nominal path identity. Two values from different abstract graphs are considered compatible as long as they satisfy the same bound, consistent with Yz's structural type system throughout.

---

## Relationship to other features

| Concept | In other languages | In Yz |
|---|---|---|
| Generics | `List<T>`, `List[T]` | `List : { T; ... }` — T is a field |
| Type alias | `type Bar = Foo` | `Bar : Foo` |
| Generic instantiation | `List<String>`, `List[String]` | `StringList : List(String)` |
| Associated types | Rust `trait Graph { type Node; }`, Scala `g.Node` | `Graph : { Node #(); ... }` + `g.Node` |
| Existential associated types | Rust `dyn Trait`, Scala existential types | not needed — bound is the type (structural) |
| Metatype | Rust `PhantomData<T>`, Scala `Type`, Kotlin `KClass<T>` | `#()` — implicit, never written |
