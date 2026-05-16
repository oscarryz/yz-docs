#spec 
# 6. Blocks and Scoping

This chapter defines the semantics of blocks of code (bocs) — Yz's unified abstraction — and variable scoping rules.

## 6.1 Blocks Are Everything

In Yz, a boc is the universal building block. A single `{ }` construct serves as:

| Traditional Concept | Yz Equivalent |
|---------------------|---------------|
| Variable | `name: "Alice"` — a boc returning a value |
| Function | `greet: { name String; "Hello, `name`!" }` |
| Class / Struct | `Person: { name String; age Int }` — uppercase name |
| Module | A `.yz` file is a boc |
| Closure | A boc that captures variables from its enclosing scope |
| Actor | Every boc instance (see Chapter 8) |

## 6.2 Boc Lifecycle

### 1. Declaration

A boc is declared and given a name:

```yz
counter: {
    count: 0
    increment: { count = count + 1 }
    count
}
```

### 2. Instantiation / Invocation

There are three distinct boc forms with different execution semantics:

| Form | Name | Semantics |
|---|---|---|
| `Foo: { field T; ... }` | Uppercase, body | Type — each call creates a new independent instance |
| `foo: { field T; ... }` | Lowercase, body | Singleton actor — shared state, calls serialize |
| `foo #(param T, ...) { ... }` | Any case, with `#(...)` | Stateless function — each call is an independent goroutine, calls parallel |

```yz
Person: { name String; age Int }
p1: Person("Alice", 30)   // New instance — independent of p2
p2: Person("Bob", 25)     // Another new instance

greet: { name String; print("Hi, ${name}!") }
greet("Alice")             // Executes greet singleton (serialized if concurrent)
greet("Bob")               // Also executes greet singleton

hi #(name String) { print("Hi, ${name}!") }
hi("Alice")                // Independent goroutine — parallel-safe
hi("Bob")                  // Another independent goroutine
```

The key distinction: **body-form bocs** (with or without uppercase) are actors whose fields persist. **Boc declaration form** bocs (with `#(...)`) are stateless — parameters are local to each call, no persistent state, `hi.name` does not exist. See §4.3 for how this affects type compatibility when bocs are passed as arguments.

### 3. Completion

A boc completes when:
- All expressions in its body have been evaluated
- All inner bocs have completed (structured concurrency — see Chapter 8)

## 6.3 Scope Rules

### Lexical Scoping

Yz uses **lexical (static) scoping**. A variable is visible from its declaration point to the end of the enclosing boc:

```yz
outer: {
    x: 10
    inner: {
        y: 20
        z: x + y    // x is visible (from outer scope)
    }
    // y and z are NOT visible here
}
```

### Scope Hierarchy

```
Source File (top-level boc)
  └── Declaration
  └── Boc
        └── Declaration
        └── Nested Boc
              └── ...
```

Each `{ }` creates a new scope. Inner scopes can read and write variables from outer scopes.

### Variable Shadowing

A variable in an inner scope can shadow a variable with the same name from an outer scope:

```yz
x: 10
inner: {
    x: 20           // Shadows outer x
    print("${x}")    // Prints 20
}
print("${x}")        // Prints 10
```

## 6.4 Variable Capture (Closures)

Bocs capture variables from enclosing scopes by **reference** — modifications to captured variables are visible in the original scope:

```yz
count: 0
increment: {
    count = count + 1   // Captures 'count' by reference
}
increment()
print("${count}")   // Prints 1
```

### Capture in Type Instances

When a boc is instantiated as a type, captured variables become part of the instance's state:

```yz
make_counter: {
    start Int
    count: start
    increment: { count = count + 1; count }
    get: { count }
}

c: make_counter(0)
c.increment()      // count = 1
c.get()            // 1
```

## 6.5 Nested Bocs

Bocs can be nested to any depth. Each nested boc has access to its own scope and all enclosing scopes:

```yz
outer: {
    a: 1
    middle: {
        b: 2
        inner: {
            c: 3
            result: a + b + c   // All accessible
        }
    }
}
```

### Nested Types

A type defined inside another boc is scoped to that boc:

```yz
app: {
    Config: {
        host String
        port Int
    }
    config: Config("localhost", 8080)
}

// Config is accessible as app.Config from outside
```

## 6.6 Boc as Value

Bocs are first-class values. They can be:

### Passed as Arguments

```yz
apply: {
    f #(Int, Int)
    x Int
    f(x)
}

double: { n Int; n * 2 }
result: apply(double, 5)   // 10
```

### Returned from Bocs

```yz
make_adder: {
    n Int
    { x Int; x + n }    // Returns a boc that captures n
}

add5: make_adder(5)
add5(3)                  // 8
```

### Stored in Collections

```yz
operations: [
    { x Int; x + 1 },
    { x Int; x * 2 },
    { x Int; x - 3 }
]
```

## 6.7 Boc Identity

Two boc values are **not** `==` based on their code — they are `==` only if they are the **same instance** (reference equality for boc values used as functions) or if they have the same structural fields with equal values (for data-carrying bocs / types).

```yz
a: { x Int; x + 1 }
b: { x Int; x + 1 }
a == b    // false — different instances, even if same structure

p1: Person("Alice", 30)
p2: Person("Alice", 30)
p1 == p2  // true — same field values (structural equality)
```

## 6.8 Top-Level Scope

Each `.yz` file has an implicit top-level boc. Declarations at the top level are the file's exports:

```yz
// file: math.yz
pi: 3.14159
e:  2.71828

circle_area: {
    radius Decimal
    pi * radius * radius
}
```

All top-level declarations are accessible from other files via the module system (see Chapter 9).
