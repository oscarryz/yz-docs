#feature

A single uppercase letter declares a **generic type parameter**. Yz follows the same convention as Go, Rust, Java, and Scala: type parameters are always declared explicitly in the type body.

## Declaring a generic type

Type parameters are declared as bare uppercase single-letter identifiers inside the type's body, before the fields that use them:

```yz
Box: {
    T          // T is a type parameter — not a field
    value T    // field whose type is T
}
```

The bare `T` line declares `T` as a type parameter. Any field after it can use `T` as its type.

> **Why explicit declaration?** Requiring the bare `T` line (rather than inferring T from field types) is consistent with every mainstream language and avoids surprising behaviour. The compiler needs to know which single-letter identifiers are type parameters vs concrete types.

## Construction — type inferred from value

When constructing a generic type, the type argument is **inferred from the constructor arguments**:

```yz
b: Box(42)         // T inferred as Int
s: Box("hello")    // T inferred as String
```

This is equivalent to Go's `NewBox(42)` with type inference.

## Typed variable declaration

When you want to name the type explicitly — for documentation, for an uninitialized variable, or to constrain later assignment — use a typed declaration:

```yz
s Box(String) = Box("hello")    // T is explicitly String
```

`Box(String)` in the type-annotation position means "a Box parameterized with String", the same way `Box<String>` means in Java/Rust or `Box[String]` in Go/Scala. The `()` syntax is Yz's notation for type arguments.

## Multiple type parameters

```yz
Pair: {
    K, V       // two type parameters
    key K
    value V
}

p: Pair("name", 42)   // K = String, V = Int
```

## Generic variant types

Variant (sum) types follow the same rules:

```yz
Option: {
    V
    Some(value V)
    None()
}

x: Some("hello")    // V = String
match x
    { Some => print(x.value) },
    { None => print("nothing") }
```

## Generic methods

Methods defined inside a generic type receive `self` as a pointer to the parameterized type:

```yz
Container: {
    T
    value T
    get #(T) {
        value
    }
}

c: Container(42)
print(c.get())    // prints 42
```

The compiler emits `func (self *Container[T]) Get() *std.Thunk[T]` with the correct type parameter on the receiver.

## Constraint inference

The compiler **automatically infers** what operations a type parameter T must support by scanning how T-typed values are used inside method bodies. This happens without any syntax — you never write `T: Comparable` or `where T:`.

When you call a method or use an operator on a T-typed value:

```yz
Ordered: {
    T
    value T
    is_less #(other T, Bool) {
        value < other    // ← compiler records: T must support < (lt method)
    }
}
```

The compiler records: **T requires `lt`**.

At every constructor call site the compiler checks that the concrete type satisfies **all** inferred requirements and reports every missing method at once:

```
error: type constraint violation for Ordered:
Item is missing methods required by T:
  lt [used in Ordered.is_less]
  to_string [used in Ordered.describe]
```

This is **Option 4** behaviour — all violations reported together, not one at a time after repeated fix-compile cycles.

### Success: constraint satisfied

```yz
o: Ordered(42)      // T=Int; Int has lt → OK, no error
```

### Failure: constraint violated

```yz
Item: {
    name String     // no lt method
}

o: Ordered(Item("x"))   // compile error: Item is missing lt
```

## Deferred / not yet implemented

- **`Box(String)` as a type-only constructor** — `word: Box(String)` to create a Box[String] without providing a value yet, then `word.value = "hello"` later. This requires passing a type as a constructor argument, which has no equivalent in mainstream languages and is complex to implement.
- **Implicit type parameters** — `Box: { value T }` where `T` is used in a field but not declared. This would cause an "undefined type: T" error currently.
- **Named constraints** — `T Comparable` or `T Printable`; constraint inference is automatic (no syntax needed).
- **Multiple type parameters in BocWithSig** — `#(key K, value V)`.
