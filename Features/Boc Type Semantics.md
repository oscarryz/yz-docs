#feature 
# Boc Type Semantics

This document summarizes the runtime semantics of boc types — how singletons differ from instantiable types, and when a boc returns a value vs. runs side effects.

## Singletons (lowercase names)

A lowercase-named boc is a singleton. Its state is shared — every call operates on the same instance:

```yz
counter: {
  count: 0
  increment: { count = count + 1 }
  value: { count }
}

counter.increment()
counter.increment()
print(counter.value())  // 2
```

Calling `counter.increment()` multiple times mutates the same `count` variable.

## Types (uppercase names)

An uppercase-named boc is a **type**. Each call constructs a new, independent instance:

```yz
Person: {
  name String
  age Int
  greet #() {
    print("Hi, I'm `name`")
  }
}

alice: Person("Alice", 30)
bob: Person("Bob", 25)

alice.greet()  // Hi, I'm Alice
bob.greet()    // Hi, I'm Bob
```

`alice` and `bob` are independent — modifying one does not affect the other.

## Return semantics

A boc returns the value of its **last expression**:

```yz
add #(x Int, y Int, Int) {
  x + y   // this is returned
}

result: add(3, 4)  // 7
```

When a boc has no explicit return type, calling it executes its body as a side effect. The boc itself (as a value) can be inspected via field access:

```yz
counter.count   // read the current count field
```

## Closures

Nested bocs capture variables from their enclosing scope:

```yz
Person: {
  name String
  greet: {
    print("I'm `name`")  // captures `name` from Person
  }
}

p: Person("Alice")
p.greet()  // I'm Alice
```

The compiler tracks which variables are captured and ensures they remain accessible.

## Type-only declarations

A `#(...)` declaration without a body defines only the **shape** of a boc — no runtime instance is created:

```yz
Greeter #(greet #())    // interface: requires greet #()
Point #(x Int, y Int)   // struct: two Int fields
```

These are used for structural typing. See [Structural typing](Structural%20typing.md) and [Block type](Block%20type.md).
