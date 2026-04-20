#feature 
# Blocks of Code (Bocs)

In Yz, everything is a **block of code** (boc). A boc plays the role that packages, modules, functions, methods, closures, objects, and classes play in other languages — unified under one construct.

## The basics

A block is a sequence of expressions between `{` and `}`. Assign it to a variable with `:` and call it with `()`:

```yz
hello: {
  print("Hello, World!")
}
hello()  // Hello, World!
```

## Variables inside bocs

Variables declared inside a boc belong to that boc. They can be other bocs:

```yz
hi: {
  text: "Hello"
  recipient: "World"
  action: {
    print("`text`, `recipient`!")
  }
  action()
}
hi()  // Hello, World!
```

Access variables via `.` notation and call inner bocs directly:

```yz
hi.action()  // Hello, World!
```

## Modifying variables before execution

Use dot notation to modify a boc's variables, then execute it:

```yz
hi.text = "Goodbye"
hi.recipient = "everybody"
hi()  // Goodbye, everybody!

hi.text  // "Goodbye"  (accessible after execution too)
```

Or pass values in the call — positional order matches declaration order:

```yz
hi("Goodbye", "everybody")  // same result
```

Use named arguments to target specific variables:

```yz
hi(recipient: "Yz world")  // Hello, Yz world!
hi(action: { print("Nothing") })  // (prints "Nothing")
```

## Multiple return values

The last N expressions in a boc are its return values. Assign them to multiple variables:

```yz
swap: {
  a String
  b String
  b  // second-to-last
  a  // last
}

x: "hello"
y: "world"
x, y = swap(x, y)  // x = "world", y = "hello"
```

The assignment order is right-to-left on the left side, bottom-to-top on the return values.

## Singletons vs types

**Lowercase** names define **singleton** bocs — there is exactly one, shared across all callers:

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

**Uppercase** names define **instantiable** bocs (types). Each call produces an independent instance:

```yz
Person: {
  name String
  age Int
  greet #() {
    print(name)
  }
}

alice: Person("Alice", 30)
bob: Person("Bob", 25)

alice.greet()  // Alice
bob.greet()    // Bob
```

Instances can be created with positional or named arguments:

```yz
alice: Person("Alice", 30)            // positional
alice: Person(name: "Alice", age: 30) // named
```

## Boc signatures

When a boc has explicit input and output types, declare them with `#(params)`. This is the **shorthand form** — signature and body together:

```yz
// Takes a String, returns nothing
greet #(name String) {
  print(name)
}
greet("Alice")

// Takes two Ints, returns an Int (last unnamed type = return type)
add #(x Int, y Int, Int) {
  x + y
}
result: add(3, 4)  // 7
```

The **body-only form** separates the signature from the body with `=`. The body re-declares the parameters:

```yz
greet #(name String) = {
  name String
  print(name)
}
```

See [Block type](Block%20type.md) for the full boc type syntax.

## Default parameters

```yz
greet #(name String = "Alice") {
  print(name)
}
greet()       // Alice
greet("Bob")  // Bob
```

Shorthand (type inferred from default):

```yz
greet #(name: "Alice") {
  print(name)
}
```

## Declare first, assign later

Declare a signature without a body, then assign it:

```yz
greet #(name String)

greet = { name String
  print(name)
}

greet("Alice")
greet("Bob")
```

## Type-only declarations

A signature without a body defines a **type** — struct or interface depending on contents:

```yz
// Struct: data fields only
Point #(x Int, y Int)

// Interface: method signatures only
Greeter #(greet #())

// Mixed: data + methods
Named #(name String, greet #())
```

These participate in structural typing — any boc with the right shape satisfies the type. See [Structural typing](Structural%20typing.md).

## Stateful actors vs. stateless functions

Every boc call is async and routed through an actor queue — but the two syntactic forms have meaningfully different semantics:

### The body form — stateful actor

```yz
add: { a Int; b Int; a + b }
```

`add` is a singleton actor. Its fields `a` and `b` persist between calls. After `add(3, 4)`, `add.a == 3`. Concurrent calls are serialized through `add`'s queue — only one runs at a time. This is correct for objects and singletons with meaningful shared state, but wrong for pure utility functions (concurrent callers queue unnecessarily).

### The BocWithSig form — stateless function

```yz
add #(a Int, b Int, Int) { a + b }
```

The `#(...)` declaration changes everything: `a` and `b` are **parameters local to each call**, not persistent fields. Each invocation is an independent goroutine. `add.a` does not exist. Concurrent calls run in **parallel**. Intermediate variables in the body are also always local:

```yz
add #(a Int, b Int, Int) {
    r Int = a + b   // local — add.r does not exist
    r
}
```

### The rule

> If a boc's declaration includes `#(...)`, it is a **stateless function** — each call is independent and parallel-safe.
> If it has no `#(...)`, it is a **stateful actor** — fields persist between calls, concurrent calls serialize.

This is why the `#(...)` syntax emerged: not just for type annotation, but as the declaration of intent that "this boc is a function, not an object."

### Higher-order bocs: named vs. anonymous params

When a boc is passed as an argument, the parameter type in the signature declares what the callee expects of it:

**Anonymous types** (`#(String, Int)`) — callability only:

```yz
map #(func #(String, Int)) { ... }
```

`map` only calls `func(item)` — it never accesses `func`'s fields. Both stateful and stateless bocs satisfy `#(String, Int)`.

**Named params** (`#(name String, Int)`) — field access required:

```yz
greet #(person #(name String, Int)) {
    println(person.name)   // accessing .name — person must have this field
}
```

`greet` accesses `person.name`, so `person` must be a stateful boc with a persistent `name` field. A stateless boc does not satisfy this — it has no `.name`. Passing a stateless boc here is a type error.

This gives a natural subtyping: a stateful boc satisfies **both** `#(name String, Int)` (field-accessible + callable) and `#(String, Int)` (callable only). A stateless boc satisfies **only** the anonymous form.

### Form summary

| Form | Fields persist | Concurrent calls | `.field` access |
|---|---|---|---|
| `foo: { field T; ... }` | Yes — singleton | Serialized | `foo.field` valid |
| `Foo: { field T; ... }` | Yes — per instance | Parallel (each fresh) | `instance.field` valid |
| `foo #(param T, ...) { ... }` | No | Parallel | `foo.param` invalid |

For utility functions (`add`, `max`, `filter`, etc.), always use the BocWithSig form. For objects and actors with meaningful shared state, use the body form. See also [Questions/Stateless bocs and pure functions](../Questions/Stateless%20bocs%20and%20pure%20functions.md).
