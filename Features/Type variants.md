#feature

## What are type variants?

A type variant is a type that can hold exactly one of several named forms at a time — also called a **sum type** or **discriminated union**. Each form (constructor) carries its own data, and you always know which form a value holds.

## Syntax

Declare a variant type as a regular boc, using zero-argument or value-carrying constructors as its members:

```yz
Option: {
    V
    Some(value V)
    None()
}
```

- `V` declares a generic type parameter (see [Generics](Generics.md)).
- `Some(value V)` is a constructor that carries one field named `value` of type `V`.
- `None()` is a constructor that carries no data.

## Constructing a variant value

Call one of the constructors directly:

```yz
x: Some("hello")    // x is an Option holding a String
y: None()           // y is an Option carrying no data
```

The type argument is inferred from the constructor argument.

## Matching on a variant

Use `match` with one branch per constructor. Each branch is a boc literal with the constructor name as the pattern:

```yz
match x
    { Some => print(x.value) },
    { None => print("nothing") }
```

- The pattern `Some` (or `None`) is the constructor name — no explicit deconstruction syntax.
- `x.value` accesses the field carried by the matched constructor.
- Every constructor must be handled; the compiler enforces exhaustiveness.

## More examples

### Result — two type parameters

```yz
Result: {
    T, E
    Ok(value T)
    Err(error E)
}

r: Ok(42)
match r
    { Ok  => print(r.value) },
    { Err => print(r.error) }
```

### Shape — constructors with different fields

```yz
Shape: {
    Empty()
    Circle(radius Decimal)
    Rectangle(width Decimal, height Decimal)
}

s: Circle(3.0)
match s
    { Empty     => print("no shape") },
    { Circle    => print(s.radius) },
    { Rectangle => print(s.width) }
```


## See also

- [SumTypes](SumTypes.md) — brief conceptual overview
- [Generics](Generics.md) — how generic type parameters work