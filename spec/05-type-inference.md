#spec
# 5. Type Inference

This chapter defines how the Yz compiler infers types when they are not explicitly declared.

## 5.1 Overview

Yz supports **local type inference** — the compiler deduces the type of a variable from the expression assigned to it. Types flow **forward** from initializers, and **backward** from call sites to generic parameters.

Explicit type annotations are always optional but allowed for documentation or disambiguation.

## 5.2 Variable Type Inference

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

### Unit Return

If the last expression produces `Unit`, or if there are no expressions (only declarations/statements), the return type is `Unit`:

```yz
log: {
    message String
    print(message)    // print returns Unit
}
// log : #(message String)  — implicit Unit return
```

## 5.4 Generic Type Inference

Generic type parameters are resolved at the **use site** — the compiler examines the concrete types of arguments passed to determine generic bindings.

### Basic Inference

```yz
Box: {
    T
    value T
}

b: Box(42)         // T is bound to Int
s: Box("hello")    // T is bound to String
```

### Inference from Context

```yz
identity: {
    T
    x T
    x           // returns T
}

n: identity(42)        // T = Int, n : Int
s: identity("hello")   // T = String, s : String
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
    // a > b → a.>(b) → T must have method > #(T, Bool)
}

max(3, 7)              // OK: Int has >
max("a", "b")          // OK if String has >
max([1], [2])          // ERROR if [Int] has no >
```

## 5.5 Structural Compatibility Inference

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

## 5.6 Match Expression Type Inference

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

## 5.7 Array and Dictionary Inference

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

## 5.8 Inference Limitations

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

## 5.9 Summary of Inference Rules

| Context | Inferred From |
|---------|---------------|
| Short declaration `x: expr` | Type of `expr` |
| Boc return type | Last expression(s) in boc body |
| Generic parameter `T` | Concrete type of argument at call site |
| Array `[a, b, c]` | Common type of elements |
| Dict `[k1:v1, k2:v2]` | Common key type, common value type |
| Match expression | Common return type of all branches |
| Width subtyping | Structural check: target has subset of source fields |
