# 3. Expressions and Statements

This chapter defines the semantics of Yz expressions and statements, complementing the syntactic grammar in Chapter 2.

## 3.1 Statements vs. Expressions

Yz distinguishes between statements and expressions:

- **Expression** — produces a value. Expressions can be used as statements (their value is the boc's return candidate).
- **Statement** — performs an action but does not produce a value that can be used in a larger expression. Declarations and keyword statements are statements.

### Statements

| Form | Kind | Example |
|------|------|---------|
| Short declaration | `id : expr` | `name: "Alice"` |
| Typed declaration | `id Type` | `age Int` |
| Typed declaration + init | `id Type = expr` | `age Int = 30` |
| Boc signature + body | `id #(...) { ... }` | `greet #(String) { "Hi" }` |
| Assignment | `target = expr` | `name = "Bob"` |
| `return` | | `return value` |
| `break` | | `break` |
| `continue` | | `continue` |
| `mix` | | `mix Named` |

### Expressions

Everything else is an expression. Expressions include:

- Literals: `42`, `3.14`, `"hello"`, `{ ... }`
- Identifiers: `name`, `Person`
- Method calls: `greet("Alice")`, `a + b`, `list << item`
- Member access: `person.name`
- Index access: `array[0]`
- Match expressions: `match { ... }`
- Array/Dict literals: `[1, 2, 3]`, `["a": 1]`
- Grouped: `(expr)`

## 3.2 Declaration Semantics

### Short Declaration (`:`)

```yz
name: "Alice"
```

- Creates a new variable `name`
- **Infers the type** from the right-hand side (`String`)
- The variable is initialized and available from this point forward
- Equivalent to: `name String = "Alice"`

### Typed Declaration

```yz
age Int
```

- Creates a new variable `age` of type `Int`
- The variable is **uninitialized** — it must be assigned before use, or passed as an argument during invocation
- Uninitialized typed declarations serve as **parameters** when the boc is invoked

### Typed Declaration with Initialization

```yz
age Int = 30
```

- Creates a new variable with explicit type **and** a default value
- The default can be overridden during invocation

### Summary: Parameters vs. Fields

Inside a boc, variables serve dual roles:

| Has default value? | Role when boc is invoked |
|--------------------|--------------------------|
| No (`age Int`) | **Required parameter** — must be provided |
| Yes (`age: 30`) | **Optional parameter / field** — uses default if not provided |

## 3.3 Assignment Semantics

```yz
name = "Bob"
```

- `=` is the **only operator** in Yz — it assigns a value to an existing variable
- `=` is **not** an expression — it cannot be used inside other expressions
- The target must be a previously declared variable or a member access (`obj.field = value`)

### Multiple Assignment

```yz
a, b = swap("hello", "world")
```

- The right-hand side must produce enough values for the left-hand side
- Values are matched right-to-left: the last value goes to the last variable, etc.
- Parentheses on the LHS group assignments for destructuring from multi-return bocs

## 3.4 Expression Evaluation

### Non-Word Method Invocation

All arithmetic, comparison, and logical operations are method calls using non-word invocation syntax:

```yz
a + b        // → a.+(b)
x == y       // → x.==(y)
-n           // → n.-()  (unary negation)
```

All non-word methods have **equal precedence** and are evaluated **left-to-right** (Smalltalk-style). Use parentheses to control grouping:

```yz
1 + 2 * 3        // → (1.+(2)).*(3) = 9
1 + (2 * 3)      // → 1.+(2.*(3)) = 7
(a > 0) && (a < 10)  // Parentheses needed for intended grouping
```

### Method Invocation

```yz
greet("Alice")           // Positional argument
greet(name: "Alice")     // Named argument
person.greet()           // Member method call
1.to(10)                 // Method on literal
```

### Conditional Expression

The `?` method on `Bool` selects between two bocs:

```yz
condition ? { expr_true }, { expr_false }
```

This is equivalent to `condition.?(true_boc, false_boc)`. The `?` method executes one of the two bocs and returns its result.

### Member Access

```yz
person.name        // Access field 'name'
person.greet()     // Invoke method 'greet'
person.name = "X"  // Assign to field
```

### Index Access

```yz
array[0]            // Access element — returns the value
dict["key"]         // Access element — returns Option(V)
array[0] = 42       // Assign to element
```

## 3.5 Block Return Values

The **last expression(s)** in a block body are the block's return value(s):

```yz
add: {
    a Int
    b Int
    a + b   // ← this is the return value
}
```

### Multiple Return Values

Multiple trailing expressions produce multiple return values:

```yz
swap: {
    a String
    b String
    b       // second-to-last: first return value
    a       // last: second return value
}

x, y = swap("hello", "world")  // x = "world", y = "hello"
```

### Explicit Return Type

When a boc has an explicit return type in its signature, the last expression must match:

```yz
greet #(name String, String) {
    "Hello, `name`!"   // Must be String
}
```

### Implicit Return Type

When no explicit return type is given:
- If the last expression returns `Unit` → the boc returns `Unit`
- If assigned to a variable → the boc returns its **instance** (for field access)
- If not assigned → the boc's logic executes (side effects)

## 3.6 Boc Invocation Semantics

### Positional Arguments

```yz
multiply(5, 3)
```

Arguments are assigned to the boc's variables in **declaration order**, mapping to variables that lack a default value first (required parameters).

### Named Arguments

```yz
divide(numerator: 10, denominator: 2)
```

Arguments are assigned to specifically named variables, regardless of order.

### Type Instantiation

When invoking a user-defined type (uppercase name):

```yz
p: Person("Alice", 30)      // New instance (positional)
p: Person(name: "Alice")    // New instance (named)
```

Each invocation creates a **new, independent instance**.

### Boc Invocation (lowercase name)

When invoking a regular boc (lowercase name):

```yz
result: greet("Alice")      // Executes the boc
```

The boc is executed with the provided arguments. If the boc has state, it is **shared** (singleton semantics).

## 3.7 Commas vs. Semicolons

| Context | Separator | Example |
|---------|-----------|---------|
| Inside `( )` | `,` (comma) | `greet("Alice", 30)` |
| Inside `[ ]` | `,` (comma) | `[1, 2, 3]` |
| Inside `{ }` | `;` (semicolon, usually via ASI) | `{ a: 1; b: 2 }` |
| Between conditional bocs | `,` (comma) | `{ cond1 => a }, { cond2 => b }` |

## 3.8 Expression as Statement

Any expression can appear as a statement. Its value becomes a return-value candidate for the enclosing boc:

```yz
process: {
    data: fetch()
    transform(data)     // Expression-statement: return value candidate
    validate(data)      // Expression-statement: this becomes the return value (last expression)
}
```

## 3.9 Operator Summary

Yz has exactly **one** operator: `=` (assignment).

Everything else that looks like an operator is a **method call** using non-word invocation:

| Looks like | Actually is |
|-----------|-------------|
| `a + b` | `a.+(b)` |
| `a == b` | `a.==(b)` |
| `-n` | `n.-()` (unary) |
| `x ? { a }, { b }` | `x.?(boc_a, boc_b)` |
| `list << item` | `list.<<(item)` |
| `a = b` | Assignment (the only real operator) |
