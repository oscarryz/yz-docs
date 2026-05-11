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
| `Unit` | No meaningful value (like `void`) | (implicit) |

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

### Stateful vs. Stateless Bocs

The presence of `#(...)` in a boc declaration determines its execution model:

- **Body form** `foo: { field T; ... }` — stateful actor. Fields persist between calls. Concurrent calls are serialized through the boc's actor queue. `foo.field` is accessible from outside.
- **BocWithSig form** `foo #(param T, ...) { ... }` — stateless function. Parameters are local to each call. Concurrent calls run in parallel. `foo.param` does not exist between calls.

This distinction also applies to boc types used as HOF parameters (see §4.3.1).

### 4.3.1 Named vs. Anonymous Params in Boc Types

When a boc type is used as a parameter type, the named/anonymous distinction declares what the callee will do with it:

| Signature type | What callee expects | Who satisfies it |
|---|---|---|
| `#(String, Int)` | Only callability — `func("x")` | Both stateful and stateless bocs |
| `#(name String, Int)` | Field access + callability — `person.name` | Only stateful bocs |

```yz
// Anonymous: only calls func — stateless or stateful both work
map #(func #(String, Int)) { ... }
map(#(item String, Int) { item.length() })   // stateless ✓
map({ item String; item.length() })          // stateful ✓

// Named: accesses person.name — must be stateful
greet #(person #(name String, Int)) {
    println(person.name)   // requires persistent field
}
greet({ name: "Alice"; name.length() })      // stateful ✓
greet(#(name String, Int) { name.length() }) // stateless — ERROR: no .name field
```

A stateful boc satisfies both forms (wider capability). A stateless boc satisfies only the anonymous form.

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

### Unit Type

A boc that doesn't return a meaningful value has return type `Unit`:

```yz
say_hi #(name String) {
    print("Hi, ${name}!")
    // Returns Unit implicitly
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

### Declaration

Type parameters are always **declared explicitly** in the type body as bare uppercase identifiers, before the fields that reference them:

```yz
Box: {
    T          // T is a type parameter — not a field
    value T    // field whose type is T
}
```

A bare single-letter uppercase identifier on its own line declares a type parameter. This is required — using a single-letter type name in a field without declaring it first is an error:

```yz
// WRONG — T not declared, causes "undefined type: T"
Bad: {
    value T
}
```

This explicit-declaration rule is consistent with every mainstream statically-typed language and avoids ambiguity about which identifiers are type parameters vs concrete type names.

### Instantiation by Inference

Generic types are resolved at the **use site** by inferring the type argument from the constructor arguments:

```yz
b: Box(42)          // T inferred as Int
s: Box("hello")     // T inferred as String
```

### Typed Variable Declaration

To declare a variable with an explicit type annotation, use `TypeName(TypeArg)` in type position:

```yz
s Box(String) = Box("hello")    // explicit: s is a Box[String]
b Box(Int) = Box(42)            // explicit: b is a Box[Int]
```

`Box(String)` in a type-annotation position means "Box parameterized with String" — analogous to `Box<String>` in Java/Rust or `Box[String]` in Go/Scala. This is distinct from `Box("hello")` in expression position, which is a constructor call.

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

## 4.8 Equality Semantics

The `==` method is defined on **every type** and performs **structural equality**:

- **Built-in types**: Value comparison (`42 == 42`, `"a" == "a"`)
- **Bocs**: Recursive field-by-field comparison — all fields must be `==`
- **Variants**: Same discriminant tag AND all fields `==`
- **Arrays**: Same length AND element-wise `==`
- **Dictionaries**: Same keys AND value-wise `==`

## 4.9 Type Summary

```
Types:
  Built-in     : Int, Decimal, String, Bool, Unit
  Boc type     : #(params..., return_types...)
  Array type   : [ElementType]
  Dict type    : [KeyType:ValueType]
  User-defined : Uppercase-named boc with fields/methods
  Variant      : Constructor-based subtypes within a user-defined type
  Generic      : Single uppercase letter (T, E, K, V, etc.)

Compatibility:
  Structural — based on field names + types, not type names
  Width subtyping — wider types assignable to narrower types
  No nominal typing — name is a label, not identity
```
