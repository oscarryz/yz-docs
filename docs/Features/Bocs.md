#feature 
# Blocks of Code (Bocs)

In Yz, everything is a **block of code** (boc). A boc plays the role that packages, modules, functions, methods, closures, objects, and classes play in other languages — unified under one construct.

## The basics

A **boc literal** is a sequence of expressions between `{` and `}`. Assign it to a variable with `:` — a **short boc declaration** — and call it with `()`:

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

## Boc Interface `#()`

A `#(params)` declares the boc's inputs and outputs explicitly. Labeled params are input fields; unlabeled types at the end are outputs. This is the **boc declaration** form:

```yz
// named input, no output
greet #(name String) {
  print(name)
}

// two named inputs, one unlabeled output
add #(x Int, y Int, Int) {
  x + y
}
result: add(3, 4)  // 7
```

The **boc expanded form** uses `=` to separate interface from body:

```yz
greet #(name String) = {
  name String
  print(name)
}
```

See [Boc Interface](Boc%20Interface.md) for the full syntax, all entry forms, and encapsulation.

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

## Boc signatures

A boc interface without a body defines a structural type — a struct if it has only data fields, an interface if it has only boc fields:

```yz
Point #(x Int, y Int)            // struct: two Int fields
Greeter #(greet #())             // interface: any boc with greet qualifies
Named #(name String, greet #())  // mixed: data + method
```

## Initialization

Calling a boc with uninitialized variables is a compile error. What counts as "uninitialized" depends on the form.

In a **boc declaration** (`name #(params) { body }`): the interface is the complete contract. A variable in the body that is not in the interface and has no default has no way to receive a value — compile error:

```yz
n #(name String) {
    last_name String   // error: not in interface, no default
}
```

In a **short boc declaration** (`name : { ... }`): all uninitialized variables become required parameters in the inferred interface. Calling without providing them is a compile error:

```yz
n : {
    name String
    last_name String
}
// inferred: #(name String, last_name String)

n("yz")          // error: last_name not provided
n("yz", "lang")  // ok
```

A variable with a default value becomes an optional parameter.

### Declare and assign within a body

A variable can be declared uninitialized and assigned later within the same scope. The compiler verifies it is assigned on all control-flow paths before it is read:

```yz
classify: {
    n Int
    label String          // uninitialized
    (n > 0) ? {
        label = "positive"
    }, {
        label = "non-positive"
    }
    label                 // OK — assigned on both paths
}
```

### When to use `Option(T)` vs a default

- Use **`Option(T)`** when absence is a meaningful state — the value may legitimately never exist (e.g., `last_login Option(Date)`).
- Use a **default value** (`field String = ""`) when the field has a clear zero or fallback.
- Use **definite assignment** (declare then assign before use) when the value is always set within the scope before it is needed.

Wrapping every field in `Option` to work around initialization is an anti-pattern.

## Closures

Nested bocs capture variables from their enclosing scope:

```yz
Person: {
  name String
  greet: {
    print("I'm `name`")  // captures name from Person
  }
}

p: Person("Alice")
p.greet()  // I'm Alice
```

The compiler tracks captured variables and ensures they remain accessible.
