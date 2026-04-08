# Structural Typing

Yz uses **structural typing**: a boc satisfies a type if it has the required fields and methods, regardless of its declared name. No explicit "implements" declaration is needed.

## Interface types

A type-only boc with method signatures acts as an interface:

```yz
Greeter #(greet #())

Person: {
  name String
  secret String
  greet #() {
    print(name)
  }
}

greet_all #(g Greeter) {
  g.greet()
}

main: {
  p: Person("Alice", "my secret")
  greet_all(p)  // works: Person has greet #()
}
```

`Person` satisfies `Greeter` because it has a `greet #()` method. The extra `secret String` field is ignored — structural typing only checks what the interface requires.

## Struct-like types

A type-only boc with data fields defines a named struct shape:

```yz
Point #(x Int, y Int)
```

Any boc with `x Int` and `y Int` fields is structurally compatible with `Point`.

## Mixed types

A type can combine data fields and method signatures:

```yz
Named #(name String, greet #())
```

## No declaration needed

You never annotate a type with "implements X". Compatibility is checked at the call site:

```yz
Robot: {
  model String
  greet #() {
    print("Unit `model` reporting")
  }
}

main: {
  // Both Person and Robot have greet #() — both satisfy Greeter
  greet_all(Person("Alice", ""))  // works
  greet_all(Robot("R2-D2"))       // works
}
```

If a boc is missing a required method, the compiler reports the error at the call site.
