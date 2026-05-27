#feature
# Generics — Type Parameters

## The Rule

A **single uppercase letter** (any Unicode uppercase character) is a type parameter. This is a language-level rule.

---

## Declaring a Type Parameter

### Explicit declaration

A bare single-letter line inside a boc body declares `T` as a type parameter. Under the PDT model, `T` is a **field of type `#()`** — it stores the type value. Construction must supply it:

```yz
Box: {
    T          // field of type #() — holds the element type
    value T
}
```

Because `T` is the first declared field, positional args map to it first. To pass only the value and let the compiler infer T, use a named parameter:

```yz
b: Box(Int, 42)      // positional: T = Int, value = 42
b: Box(value: 42)    // named: T inferred as Int from value argument
b.T                  // valid — returns the type Int
```

Passing a value positionally where T expects a `#()` type is an error:

```yz
b: Box(42)   // ERROR — 42 is an Int value, not a type; T expects #()
```

### Implicit (inferred from use)

If `T` is used in a field type without a bare declaration line, `T` is an anonymous inferred type variable — not a field. Construction infers T from the value argument directly:

```yz
Box: {
    value T    // T not declared — anonymous type variable, inferred from use
}

b: Box(42)       // T inferred as Int — no T field passed
b: Box("hola")   // T inferred as String
b.T              // INVALID — no T field exists
```

Both forms are valid with different constructor syntax and different capabilities:

| Form | Constructor | `b.T` |
|---|---|---|
| `{ T; value T }` | `Box(Int, 42)` or `Box(value: 42)` | valid |
| `{ value T }` | `Box(42)` | invalid — no T field |

---

## Typed Variable Declaration

To name the type argument explicitly in a variable declaration, use `TypeName(TypeArg)` in type-annotation position:

```yz
// Using implicit-T Box: { value T }
s Box(String) = Box("hello")    // T is explicitly String

// Using explicit-T Box: { T; value T }
s Box(String) = Box(value: "hello")   // named param — T inferred and matches annotation
s Box(String) = Box(String, "hello")  // fully explicit
```

`Box(String)` in type-annotation position means "a Box parameterized with String" — analogous to `Box<String>` in Java/Rust or `Box[String]` in Go/Scala. The `()` syntax is Yz's notation for type arguments.

---

## Multiple Type Parameters

When type parameters are declared explicitly (bare letter lines), the caller must supply the concrete types before the values:

```yz
Pair: {
    K, V       // explicit declarations — caller must provide K and V
    key K
    value V
}

p: Pair(String, Int, "name", 42)   // K = String, V = Int, then the values
```

When type parameters are implicit (used in fields only, not declared bare), they are inferred from the constructor arguments:

```yz
Pair: {
    key K      // K inferred from first argument
    value V    // V inferred from second argument
}

p: Pair("name", 42)   // K = String, V = Int — inferred
```

**Variant exception:** variant constructors don't carry the enclosing type's type argument — the typed declaration does:

```yz
Option: {
    V
    Some(value V)
    None()
}

o Option(String) = Some("hi")   // V = String comes from the declaration, not Some(...)
```

---

## Type Parameters in Boc Declarations

Type parameters work in boc declarations (`name #(params) { body }`) too:

```yz
// T is a type parameter; x is a value of type T; return type is T
f: #(T, x T, T) { x }

// S must satisfy the Serializable constraint (see below)
g: #(S Serializable) { ... }

// Multiple type parameters
identity: #(T, value T, T) { value }
```

---

## Generic Variant Types

Variant (sum) types follow the same rules:

```yz
Option: {
    V
    Some(value V)
    None()
}

x: Some("hello")   // V = String

match x
    { Some => print(x.value) },
    { None => print("nothing") }
```

---

## Constraints

A constraint specifies what operations a type parameter must support. Constraints are **optional**; the compiler infers them automatically from usage.

### Inferred constraints

The compiler scans how a type-parameter-typed value is used inside the boc body and records every method call or operation as a requirement:

```yz
printable: {
    value T
    value.print()   // compiler infers: T must have print #()
}

Person: {
    name String
    print: { print("My name is ${name}") }
}

printable(Person("Yz"))   // ok — Person has print
printable("oh oh")        // error: String doesn't have print
```

All constraint violations across the entire body are reported at once (not one at a time).

### Explicit constraints

Optionally, a constraint can be declared directly next to the type parameter:

```yz
// Standalone type with constrained parameter
serialize: {
    value T Serializable   // T must implement Serializable
    value.to_json()
}

// In a signature
g: #(T Serializable) { ... }

// Explicit constraint on a variable
a T Serializable
```

`Serializable` here is any boc type or structural interface already in scope. An explicit constraint is checked at the call site; an inferred constraint is checked against usage inside the body. Both forms can coexist on the same type parameter — the union of all requirements applies.

### Named constraints (structural interfaces)

A constraint is just any boc type. A structural interface captures a named shape:

```yz
Talker: #( talk #(String) )

greet: {
    thing T Talker
    thing.talk()
}
```

If the compiler can find an existing boc type in scope that exactly matches the inferred constraint shape, it names it. If multiple types match, the inferred constraint is expressed as an anonymous signature instead.

---

## Constraint Propagation

When a generic boc calls another generic boc, constraints propagate upward automatically:

```yz
greet: {
    thing T
    thing.talk()      // inferred: T has talk #()
}

greet_all: {
    things []T
    things.each({
        thing T
        greet(thing)  // greet's constraint propagates to greet_all's T
    })
}
```

`greet_all` ends up requiring `T has talk #()` even though `talk` is never called directly in its body. Chains of arbitrary depth are flattened — the full set of requirements is always surfaced at the outermost call site.

---

## Compile-Time Bocs and Constraints

> **Caution:** `Compile` implementations can silently add constraints to type parameters.

Code generated by a `Compile` implementation participates in constraint inference equally with hand-written code. If a `Compile` boc generates a method that calls `value.serialize()`, the compiler infers `serialize #(String)` as a constraint on `T` — even though the developer never wrote that call:

```yz
Box: {
    `compile_time: [Serialize]`
    T
    value T
    // Serialize generates: value.serialize() — T now requires serialize
}
```

The compiler always surfaces the full flattened constraint set in errors and tooling, including the origin of each requirement.

See also: [Compile Time Bocs](Compile%20Time%20Bocs.md)

---

## Path-Dependent Types and Associated Types

Type parameters are fields that hold type values. This enables path-dependent access (`list.T`) and associated types (`g.Node`) as natural extensions:

```yz
Graph : {
    Node #()                    // Node is a type field, implicitly #()
    neighbors #( Node, List(Node) )
}

process #( g Graph, n g.Node )  // n's type depends on g
```

The `#()` type on `Node` is implicit and never written. Constraints narrow it: `T Comparable` means T's type is `Comparable` rather than the unconstrained `#()`.

See [Path Dependent Types](Path%20Dependent%20Types.md) for the full model.

---

## Not Yet Implemented

- **`Box(String)` as type-only constructor** — declare `word Box(String)` without providing a value, then assign later. Requires passing a type as a runtime constructor argument.
- **Generic instantiation via alias** — `StringList : List(String)` — tracked in YZC-0066.
- **Path-dependent types in signatures** — `g.Node` in parameter position — tracked in YZC-0066.
- **Explicit constraint propagation across module boundaries** — tooling surfaces inferred constraints; explicit cross-module constraint annotations are not yet designed.
