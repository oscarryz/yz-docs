#spec
# 5. Type Inference

This chapter defines how the Yz compiler infers types when they are not explicitly declared.

## 5.1 Overview

Yz supports **local type inference** — the compiler deduces the type of a variable from the expression assigned to it. Types flow **forward** from initializers, and **backward** from call sites to generic parameters.

Explicit type annotations are always optional but allowed for documentation or disambiguation.

## 5.2 Variable Type Inference

### Typed Declaration (no initializer)

When a variable is declared with an explicit type and no value, no type inference occurs — the type is given directly. The variable is **uninitialized** and subject to definite assignment (§3.2): it must be assigned on all control-flow paths before it is read.

```yz
age Int          // type is Int; uninitialized — must be assigned before use
age = 30         // now initialized
```

### Short Declaration

When a variable is declared with `:`, its type is inferred from the initializing expression:

```yz
name: "Alice"       // name : String
age: 30             // age  : Int
pi: 3.14            // pi   : Decimal
flag: true          // flag : Bool
items: [1, 2, 3]    // items: [Int]
```

### Boc Type Inference

When a boc literal is assigned, the type is the structural signature of the boc:

```yz
greet: {
    name String
    "Hello, ${name}!"
}
// greet : #(name String, String)
```

### Type Identity from Uppercase Names

When applying a user-defined type constructor, the variable has that structural type:

```yz
p: Person("Alice", 30)
// p : #(name String, age Int)  — structurally equivalent to Person
```

## 5.3 Return Type Inference

A boc's return type is inferred from its **last expression(s)**:

```yz
add: {
    a Int
    b Int
    a + b       // a.+(b) returns Int → boc return type is Int
}
// add : #(a Int, b Int, Int)
```

### Multiple Return Values

```yz
swap: {
    a String
    b String
    b           // String
    a           // String
}
// swap : #(a String, b String, String, String)
```

### No Return (Unit internally)

If the last expression produces no meaningful value, or if there are no expressions (only declarations/statements), the boc returns nothing. Internally this is represented as `Unit`, but it does not surface to users.

```yz
log: {
    message String
    print(message)    // print returns nothing
}
// log : #(message String)  — returns nothing
```

## 5.4 Generic Type Inference

Generic type parameters are resolved at the **use site** — the compiler examines the concrete types of arguments passed to determine generic bindings.

### Basic Inference

With **explicit T** (`T` declared as a `#()` field), T is inferred when a named parameter is used for the value argument:

```yz
Box: {
    T          // #() field — holds the element type
    value T
}

b: Box(value: 42)    // T inferred as Int from named argument
b: Box(Int, 42)      // T explicit — also valid
b: Box(42)           // ERROR: 42 assigned to T field which expects #()
```

With **implicit T** (no bare declaration line), T is inferred directly from the positional argument:

```yz
Box: {
    value T    // anonymous type variable — inferred from constructor argument
}

b: Box(42)        // T inferred as Int
s: Box("hello")   // T inferred as String
```

### Inference from Context

```yz
identity: {
    T
    x T
    x           // returns T
}

n: identity(x: 42)        // T inferred as Int from named x argument; n : Int
s: identity(x: "hello")   // T inferred as String; s : String
```

### Multiple Type Parameters

```yz
Pair: {
    K, V
    key K
    value V
}

p: Pair("name", 42)   // K = String, V = Int
```

### Inference with Collections

```yz
first: {
    T
    list [T]
    list[0]
}

names: ["Alice", "Bob"]
n: first(names)        // T = String, n : String
```

### Constraint Checking

After inference, the compiler verifies that the resolved type satisfies all structural requirements at the use site:

```yz
max: {
    T
    a T
    b T
    (a > b) ? { a }, { b }
    // a > b → a.>(b) → T must have method > #(other T, Bool)
}

max(3, 7)              // OK: Int has >
max("a", "b")          // OK if String has >
max([1], [2])          // ERROR if [Int] has no >
```

## 5.5 Path-Dependent Type Inference

When a function signature contains a path-dependent type (`g.Node`), the compiler resolves the type by:

1. Inferring or checking the type of the leading variable (`g`)
2. Looking up the named field on that type (`Node` on whatever struct `g` is)
3. Using the stored type value as the expected type for the parameter

```yz
Graph : {
    Node #()
    ...
}

SocialGraph : {
    Node : User
    ...
}

process #( g Graph, n g.Node )

sg : SocialGraph()
u  : User("Alice")
process(sg, u)   // step 1: g = SocialGraph; step 2: SocialGraph.Node = User; step 3: n must be User
```

Parameters to the left of a path-dependent parameter are resolved first, in declaration order. If the leading variable's type cannot be determined statically, the compiler reports an error.

When a function uses unbound single-uppercase type variables (`A`, `B`) across multiple parameters, the compiler unifies them against argument types in order:

```yz
map #( collection List(A), fn #(A, B), List(B) )
```

Step 1: bind `A` from `collection`'s element type. Step 2: verify `fn`'s first parameter type matches `A`; bind `B` from `fn`'s return type. Step 3: the return type `List(B)` is now concrete.

> **Implementation note:** path-dependent inference is tracked in YZC-0066.

---

## 5.6 Structural Compatibility Inference

When a value is assigned to a variable with an explicit type, the compiler checks structural compatibility:

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

e: Employee("Alice", 30, 1001)

// The compiler checks: does Employee have name:String, age:Int? Yes.
p Person = e     // OK — width subtyping
```

## 5.7 Match Expression Type Inference

### Condition Match

The type of a condition match is the **common structural type** of all branch return values:

```yz
label: match {
    score >= 90 => "A"
}, {
    score >= 80 => "B"
}, {
    "C"
}
// All branches return String → label : String
```

### Variant Match

In variant match, the matched variable's type is narrowed to the variant's fields:

```yz
Result: {
    T, E
    Ok(value T)
    Err(error E)
}

r: Result.Ok(42)   // r : Result(Int, E)

match r {
    Result.Ok => r.value    // r narrowed: .value : Int is accessible
}, {
    Result.Err => r.error   // r narrowed: .error : E is accessible
}
```

## 5.8 Array and Dictionary Inference

### Array Literal

The element type is inferred from the elements:

```yz
nums: [1, 2, 3]          // [Int]
names: ["Alice", "Bob"]  // [String]
```

All elements must have a common structural type. If they don't, it's a compile error.

### Dictionary Literal

Key and value types are inferred from entries:

```yz
ages: ["Alice": 30, "Bob": 25]    // [String:Int]
```

### Empty Collections

Empty collections require explicit type annotations:

```yz
items: [Int]()             // Empty [Int]
lookup: [String:Bool]()    // Empty [String:Bool]
```

## 5.9 Inference Limitations

Type inference is **local** — it does not propagate across boc boundaries in complex ways. When inference is ambiguous, the compiler requires an explicit annotation:

```yz
// Ambiguous — compiler cannot infer T without a usage context
identity: {
    T
    x T
    x
}

// Must provide concrete type at call site
n: identity(42)  // OK — T = Int from argument
```

## 5.10 Summary of Inference Rules

| Context | Inferred From |
|---------|---------------|
| Short declaration `x: expr` | Type of `expr` |
| Boc return type | Last expression(s) in boc body |
| Generic parameter `T` | Concrete type of argument at call site |
| Path-dependent `g.Node` | Static type of `g`, field lookup on its struct type |
| Type variable `A` in signature | Unified against concrete argument types in order |
| Array `[a, b, c]` | Common type of elements |
| Dict `[k1:v1, k2:v2]` | Common key type, common value type |
| Match expression | Common return type of all branches |
| Width subtyping | Structural check: target has subset of source fields |
