#spec
# 4. Type System

This chapter defines Yz's type system: the kinds of types, structural compatibility rules, type variants, and generics.

## 4.1 Overview

Yz uses a **structural** type system. Two types are compatible if they have the same structure (fields and methods), regardless of their names. There is no nominal typing — type names are convenient labels, not identity markers.

Every value in Yz has a type. Types are either:

- **Built-in types** — provided by the standard library
- **User-defined types** — defined via boc declarations with uppercase names

## 4.2 Built-In Types

| Type | Description | Literal Examples |
|------|-------------|-----------------|
| `Int` | Integer (arbitrary precision) | `0`, `42`, `1000` |
| `Decimal` | Decimal floating-point | `3.14`, `0.5` |
| `String` | UTF-8 text | `"hello"`, `'world'` |
| `Bool` | Boolean | `true`, `false` |
| `Unit` | No meaningful value (like `void`) — **internal term; user-facing code says "returns nothing"** | (implicit) |

### Methods on Built-In Types

Built-in types have methods defined in the standard library. These include non-word methods like `+`, `-`, `*`, `==`, etc. (see §1.9) and word methods like `to_string()`, `to(n)`, etc.

All types support the `==` method (structural equality) and `!=` (structural inequality).

The `Bool` type supports `&&`, `||`, and `?`.

## 4.3 Boc Types

A boc's type is defined by the **signature** of its public interface — the set of fields and methods it exposes.

### Explicit Signature

A boc type is explicitly declared using the `#(...)` syntax:

```yz
greet #(name String, String) {
    "Hello, ${name}!"
}
```

The signature `#(name String, String)` means: "accepts a `String` parameter named `name`, returns a `String`."

### Signature Components

```
#( [parameters] [, return_types] )
```

| Part | Description | Example |
|------|-------------|---------|
| Parameters | Named typed values the boc accepts | `name String`, `age Int` |
| Default values | Parameters with defaults are optional | `age Int = 0` |
| Return type(s) | Type(s) of the last expression(s) | `String`, `Int, Error` |

### Examples

```yz
// No parameters, returns Int
counter #(Int) { 0 }

// One required parameter, returns String
greet #(name String, String) { "Hello, ${name}!" }

// Two parameters (one optional), returns Unit
log #(message String, level String = "INFO") {
    print("${level}: ${message}")
}

// Multiple return values
divide #(a Int, b Int, Int, Error) {
    b == 0 ? {
        0
        Error("division by zero")
    }, {
        a / b
        nil
    }
}
```

### All Bocs Have Persistent Fields

All boc forms share the same model: fields persist between calls and are accessible from outside.

- **Short boc declaration** `foo: { field T; ... }` — fields persist. `foo.field` is accessible from outside.
- **Boc declaration** `foo #(param T, ...) { ... }` — syntactic sugar for the same model. `foo.param` persists between calls and is accessible from outside.

A field declared without a default (`field T`) is a **required field**. It must be assigned on all control-flow paths before it is read, and must be provided by the caller before the boc value is passed across a boc boundary. See §3.2 (Definite Assignment).

### 4.3.1 Named vs. Unlabeled Params in Boc Types

When a boc interface is used as a parameter type, the labeled/unlabeled distinction declares what the callee will do with the passed boc:

| Signature type | What callee expects | Who satisfies it |
|---|---|---|
| `#(String, Int)` | Reads two return values (String, Int) | Any boc whose last two expressions are String and Int |
| `#(name String, Int)` | Field access by name — `person.name` — and reads Int return | Any boc with a field named `name` of type String that returns Int |

```yz
// Callee reads two return values — any boc whose last two expressions are String, Int qualifies
describe #(source #(String, Int)) {
    label, count = source()
    print("${label}: ${count}")
}
describe({ "items"; 42 })                    // ✓ — returns String and Int

// Callee accesses person.name — boc must have a field named name
greet #(person #(name String, Int)) {
    println(person.name)   // requires field named name
}
greet({ name: "Alice"; name.length() })      // ✓ — boc literal with name field, returns Int
greet(#(name String, Int) { name.length() }) // ✓ — boc declaration with name field, returns Int
```

A boc satisfies `#(name String, Int)` if it has a field named `name` of type String and returns Int. Both boc literals and boc declarations can satisfy it.

### Synthetic Signature

When no explicit signature is given, the compiler creates a **synthetic signature** from all the boc's internal variables:

```yz
// No explicit signature
person: {
    name: "Alice"
    age: 30
}
// Synthetic signature: #(name String = "Alice", age Int = 30)
// Everything is public
```

### Unit Type (internal)

`Unit` is the compiler's internal representation for "returns nothing". It does not appear in user-facing code or error messages — bocs that produce no output are described as "returns nothing".

```yz
say_hi #(name String) {
    print("Hi, ${name}!")
    // returns nothing (Unit internally)
}
```

## 4.4 User-Defined Types

A user-defined type is a boc whose identifier starts with an uppercase letter:

```yz
Person: {
    name String
    age Int
}
```

### Instantiation

User-defined types are instantiated by invocation:

```yz
p: Person("Alice", 30)        // Positional
p: Person(name: "Alice", age: 30)  // Named
```

Each invocation creates a **new, independent instance**.

### Type as Structure

A user-defined type defines a **structural shape**. The name `Person` is a label — any boc with the fields `name String` and `age Int` is structurally compatible with `Person` (see §4.6).

## 4.5 Type Variants

A type can have multiple **variant constructors**, similar to algebraic data types / sum types. Each variant is an alternative way to construct a value of the type.

### Declaration

```yz
Result: {
    T, E
    Ok(value T)
    Err(error E)
}
```

Each variant constructor (e.g., `Ok`, `Err`) is declared with parenthesized parameter(s) inside the type's body.

### Construction

```yz
success: Result.Ok(42)
failure: Result.Err("not found")
```

### Variant Matching

The `match` expression discriminates between variants:

```yz
match result {
    Result.Ok => print("Value: ${result.value}")
}, {
    Result.Err => print("Error: ${result.error}")
}
```

### Runtime Discriminant Tag

Under structural typing, variants that have the same field structure would be indistinguishable without a hidden marker. Therefore:

> Each variant carries a **hidden runtime discriminant tag** that records which constructor was used. This tag is used exclusively by `match` expressions.

The tag is an implementation detail — it is not accessible to user code.

### Common Variant Types

The standard library provides common variant types:

```yz
Option: {
    T
    Some(value T)
    None()
}

Result: {
    T, E
    Ok(value T)
    Err(error E)
}
```

## 4.6 Structural Compatibility

### Definition

Type `A` is **compatible with** type `B` if `A` has at least all the fields and methods declared in `B` with compatible types. This is also known as **width subtyping**.

### Rules

1. **Exact match**: If `A` and `B` have identical fields/methods → compatible in both directions
2. **Width subtyping**: If `A` has all of `B`'s fields plus extras → `A` is compatible with `B`, but not the reverse
3. **Method compatibility**: Method signatures must match structurally (parameter types and return types)
4. **Generic compatibility**: Generic type parameters are matched by their structure at the use site (see §4.7)

### Examples

```yz
Person: {
    name String
    age Int
}

Employee: {
    name String
    age Int
    id Int
}

// Employee is compatible with Person (has all of Person's fields)
greet #(p Person) {
    print("Hello, ${p.name}!")
}

e: Employee("Alice", 30, 1001)
greet(e)         // OK — Employee has name, age
// e.id is accessible on e, but greet only sees name, age
```

### Non-Compatible Cases

```yz
Point: {
    x Int
    y Int
}

Person: {
    name String
    age Int
}

// Point and Person are NOT compatible even though both have 2 Int-like fields
// because the field names differ
```

### Assignability

- A value of type `A` can be assigned to a variable of type `B` if `A` is compatible with `B`
- Through the assigned variable, only `B`'s fields/methods are accessible
- The original value retains its full type

```yz
e: Employee("Alice", 30, 1001)
p Person = e     // OK — width subtyping
p.name           // OK — "Alice"
p.id             // ERROR — Person has no field 'id'
e.id             // OK — e is still Employee
```

## 4.7 Generic Types

Yz uses **single uppercase letters** as generic type parameters, following the same convention as Go, Rust, Java, and Scala.

### The `#()` metatype

`#()` — the empty interface — is the **metatype**: the type whose values are types. Every UpperCase boc (`Int`, `String`, `Person`, …) satisfies `#()` structurally. It is the implicit type of every bare type parameter, analogous to how `Unit` is the implicit return type of bocs that return nothing: you never write it, the compiler infers it.

### Declaration

Type parameters are declared as bare uppercase identifiers before the fields that reference them. Semantically they are **fields of type `#()`** — they hold a type value:

```yz
Box: {
    T          // field of type #() — holds the element type
    value T    // field whose type is T
}
```

There are two forms:

**Explicit declaration** — `T` on its own line makes it a named `#()` field. The constructor must supply it, either positionally or via a named parameter for inference:

```yz
Box: {
    T          // field of type #() — holds the element type
    value T
}

b: Box(Int, 42)      // positional: T = Int, value = 42
b: Box(value: 42)    // named: T inferred as Int — preferred ergonomic form
b.T                  // valid — resolves to the type Int

b: Box(42)           // ERROR — 42 is an Int value, not a type; T expects #()
```

**Implicit declaration** — `T` used in field types without a bare declaration line is an anonymous inferred type variable. It is not a field; `b.T` is invalid.

```yz
Box: {
    value T    // T not declared — anonymous, inferred from constructor argument
}

b: Box(42)       // T inferred as Int
b: Box("hello")  // T inferred as String
b.T              // ERROR — no T field
```

The two forms have different constructor syntax and different capabilities. Use explicit T when callers need to inspect the type via path-dependent access (`b.T`) or when T is part of the public interface.

### Typed Variable Declaration

To declare a variable with an explicit type annotation, use `TypeName(TypeArg)` in type position:

```yz
// Implicit-T Box: { value T }
s Box(String) = Box("hello")    // T inferred from value, annotation confirms String

// Explicit-T Box: { T; value T }
s Box(String) = Box(value: "hello")   // named param; T inferred and matches annotation
s Box(String) = Box(String, "hello")  // fully positional
```

`Box(String)` in type-annotation position means "Box parameterized with String" — analogous to `Box<String>` in Java/Rust or `Box[String]` in Go/Scala. This is distinct from the constructor call in expression position.

### Generic Constraints — Inferred Automatically

Yz does **not** have explicit constraint syntax (e.g., `T: Comparable`). Instead, the compiler **infers constraints** by scanning how T-typed values are used inside the generic type's method bodies.

When a method calls a method or applies an operator to a T-typed value, the compiler records that T must support that operation:

```yz
Ordered: {
    T
    value T
    is_less #(other T, Bool) {
        value < other    // compiler infers: T must support < (lt)
    }
}
```

At every constructor call site, the compiler checks that the concrete type satisfies **all** inferred requirements. If any are missing, **all** violations are reported at once:

```
error: type constraint violation for Ordered:
Item is missing methods required by T:
  lt [used in Ordered.is_less]
```

A type that satisfies the constraint compiles without error:

```yz
o: Ordered(42)           // OK: Int has lt
o2: Ordered("hello")     // OK: String has lt (lexicographic)
```

A type that is missing the required method is rejected at the constructor call site:

```yz
Item: { name String }
o: Ordered(Item("x"))    // compile error: Item missing lt
```

No annotation is required on the type definition. Constraints emerge naturally from usage.

### Multiple Type Parameters

```yz
Pair: {
    K, V
    key K
    value V
}

p: Pair("name", 42)  // K = String, V = Int
```

## 4.8 Type Aliases

`Name : ExistingType` declares a **type alias** — `Name` is another name for `ExistingType`. They are structurally identical; the compiler treats them as the same type.

```yz
Bar : Foo     // Bar and Foo are the same type
```

This uses the same `:` short-declaration syntax as value declarations (`x : 42`). The compiler copies `Foo`'s signature into `Bar`.

### Generic instantiation via alias

A type alias where the right-hand side is a parameterized type creates a concrete specialization:

```yz
StringList : List(String)   // StringList.T = String
IntPair    : Pair(Int, Int)
```

### Associated type binding

Inside a concrete boc, type aliases bind the abstract type fields of an interface (see §4.9):

```yz
SocialGraph : {
    Node : User         // Node is bound to User in this boc
    Edge : Relationship
    ...
}
```

> **Implementation note:** simple structural aliases (`Bar : Foo`) are tracked in YZC-0027. Generic instantiation and associated type binding depend on YZC-0066.

---

## 4.9 Path-Dependent Types and Associated Types

A boc can declare **type fields** — fields whose values are types. Other bocs that satisfy the interface bind those fields to concrete types. Functions access them via the instance path.

### Declaring type fields (associated types)

```yz
Graph : {
    Node #()                       // Node is a type field (type = #())
    Edge #()
    add_node  #( Node )
    add_edge  #( Node, Node, Edge )
    neighbors #( Node, List(Node) )
}
```

### Binding in a concrete boc

```yz
SocialGraph : {
    Node : User          // type alias — Node = User in SocialGraph
    Edge : Relationship
    add_node  #( Node ) = { ... }
    neighbors #( Node, List(Node) ) = { ... }
}
```

`SocialGraph.Node` is `User`. `SocialGraph.Edge` is `Relationship`.

### Path-dependent type in a signature

A function that works on any `Graph` accesses the associated type through the instance:

```yz
process #( g Graph, n g.Node )
```

`g.Node` is a **path-dependent type**: its concrete value depends on which graph `g` holds. The compiler resolves it statically at the call site:

```yz
sg : SocialGraph()
u  : User("Alice")
process(sg, u)   // g.Node = SocialGraph.Node = User — verified at compile time
process(sg, 42)  // ERROR: Int does not satisfy SocialGraph.Node (User)
```

### Type variable form

When a function relates two type-parameterized arguments, single-uppercase letters serve as unbound type variables in the signature (implicit `#()`):

```yz
map #( collection List(A), fn #(A, B), List(B) )
```

`A` and `B` are inferred at the call site from the concrete types of the arguments.

> **Implementation note:** path-dependent types and associated types are tracked in YZC-0066 and YZC-0030.

---

## 4.10 Equality Semantics

The `==` method is defined on **every type** and performs **structural equality**:

- **Built-in types**: Value comparison (`42 == 42`, `"a" == "a"`)
- **Bocs**: Recursive field-by-field comparison — all fields must be `==`
- **Variants**: Same discriminant tag AND all fields `==`
- **Arrays**: Same length AND element-wise `==`
- **Dictionaries**: Same keys AND value-wise `==`

## 4.11 Type Summary

```
Types:
  Built-in         : Int, Decimal, String, Bool, Unit
  Boc type         : #(params..., return_types...)
  Array type       : [ElementType]
  Dict type        : [KeyType:ValueType]
  User-defined     : Uppercase-named boc with fields/methods
  Variant          : Constructor-based subtypes within a user-defined type
  Generic          : Single uppercase letter (T, E, K, V, etc.) — field of type #()
  Type alias       : Name : ExistingType
  Metatype         : #() — the type of types; implicit, never written

Compatibility:
  Structural — based on field names + types, not type names
  Width subtyping — wider types assignable to narrower types
  No nominal typing — name is a label, not identity

Path-dependent:
  g.Node — type stored in field Node of instance g, resolved at compile time
  T in signatures — unbound type variable, inferred at call site
```
