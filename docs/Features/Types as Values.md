#feature
#pending

# Yz Language Design: Types as First-Class Values

This document outlines the unification of types and values within the "everything is a boc" (block of code) philosophy. By treating types as values, we collapse the distinction between generics, interfaces, and variables into a single, consistent mental model.

---

## 1. The Core Grammar: The Space-Colon Rule

The language is built on three fundamental operations involving identifiers and expressions.

| Operator | Name | Purpose | Example |
| :--- | :--- | :--- | :--- |
| **Space (` `)** | **Constraint** | Defines what an identifier **must** be. | `age Int` |
| **Colon (`:`)** | **Binding** | Declares and initializes an identifier. | `x : 10` |
| **Equals (`=`)** | **Assignment** | Mutates an existing identifier. | `x = 11` |

### Unified Declaration
The full form of a declaration is: `identifier Constraint = Value`.
* **`a String = "hi"`**: Full form.
* **`a : "hi"`**: Abbreviated (Constraint is derived from the value).
* **`a String`**: Declaration only; `a` is a "hole" that must be filled.

---

## 2. Types are boc Values

In Yz, `Int`, `String`, and `Person` are not keywords or special metadata; they are **bocs** that reside in the environment.

* **Type Aliasing**: `MyInt : Int` (Binding the value `Int` to a new name).
* **Types as Arguments**: Since types are values, you can pass them to bocs just like strings or numbers.

### Generic bocs (boc Factories)
A "Generic" is simply a boc that takes a type-value and returns a new boc definition. No special `<T>` syntax is required.

```yz
Box : {
    T Type    // T is a parameter that must be a Type
    {         // This nested block is the returned "specialized" boc
        value T
    }
}

// Usage:
sbox : Box(String) 
sbox.value = "Hello"
```

---

## 3. Signatures: The "Hollow" boc

Signatures `#( ... )` are bocs that represent **Requirements**. They define the shape of a hole that needs to be filled. 

* **As Interfaces**: `HasName : #( name String )`
* **In Blocks**: A block `{ ... }` satisfies a signature if it contains the required members. There is no `implements` keyword; it is purely structural.

### Nested Requirements
Signatures can be nested to define complex structural requirements.

```yz
// greet takes 'obj'. 'obj' must have a 'name' block returning a String
greet : { 
    obj #( name #( String ) ) 
    print(obj.name())
}

Person : {
    n String
    name : { n }
}

greet(Person("Ann")) // Success: Person matches the required structure
```

---

## 4. Path-Dependent Types (Associated Types)

Because a boc can contain variables that hold Type-values, "Associated Types" are handled via standard member access.

```yz
Graph : {
    Node Type
    Edge Type
    neighbors #( Node, [Edge] )
}

// Instantiation: 3 variables in Graph, so 3 args in invocation
SocialGraph : Graph(Int, String, {
    n Node 
    ["The end"] 
})

// Function using Path-Dependent Types
print_neighbors : {
    g Graph
    n g.Node // 'n' must be the type defined by 'g's Node variable
    print(g.neighbors(n))
}

print_neighbors(SocialGraph, 1) // Valid because SocialGraph.Node is Int
```

---

## 5. Singleton and Literal Types

If `x Int` means `x` must satisfy the boc `Int`, then `x 200` means `x` must satisfy the value `200`. This allows for compile-time constant enforcement and powerful pattern matching.

```yz
// This signature only accepts the literal integer 200
SuccessHandler #( code 200, String )

handle : SuccessHandler {
    code // The parameter
    "OK" // The return
}

handle(200) // ✅ Works
handle(404) // ❌ Compiler error: 404 does not satisfy 200
```

---

## 6. Implementation Notes (Compiler Logic)

### The `#satisfies` Protocol
Internally, the compiler treats the space operator `A B` as a call to a protocol: `B.#satisfies(A)`.

1. **Structural Check**: If `B` is a signature `#( ... )`, the compiler verifies `A` has all members defined in `B`.
2. **Value Check**: If `B` is a literal (e.g., `200`), the compiler verifies `A` is exactly that literal.
3. **Type Check**: If `B` is a primitive boc (e.g., `Int`), the compiler performs machine-type validation.

### Parameter Mapping
When a boc is invoked `boc(arg1, arg2)`, the compiler maps these values to the **uninitialized declarations** inside the boc in the order they appear.

### Mutability
* **`:` (Binding)** is used for the initial declaration of a variable within a scope.
* **`=` (Assignment)** is used to update the value of an existing variable.