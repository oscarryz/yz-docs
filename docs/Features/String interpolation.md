#feature

# String Interpolation

Yz has two interpolation forms inside string literals, each with a distinct contract.

## `${}` — user-facing interpolation

The value's type must implement `to_str #(String)`. Compile error if it does not.

Built-in types (`Int`, `Decimal`, `String`, `Bool`) have `to_str` built in:

```
name: "Alice"
age:  30
print("Hello, ${name}!")                   // Hello, Alice!
print("Age: ${age}, next: ${age + 1}")     // Age: 30, next: 31
```

For user-defined types, `to_str` must be explicitly provided:

```
Person : {
    name String
    age  Int
    to_str : { "${name} (age ${age})" }
}

p: Person("Alice", 30)
print("User: ${p}")    // User: Alice (age 30)
```

Omitting `to_str` on a type used in `${}` is a **compile error**. This is intentional — it forces the programmer to decide what the user-facing representation should be.

## `` ` `` — compiler homoiconic interpolation

Embeds an expression using a compiler-generated representation that mirrors the Yz source that would recreate the value. Always works — no method required.

```
p: Person("Alice", 30)
print("Debug: `p`")
// Debug: Person(name: "Alice", age: 30)
```

Arrays are pretty-printed one element per line:

```
people: [Person("Alice", 30), Person("Bob", 42)]
print("List: `people`")
// List: [
//    Person(name: "Alice", age: 30),
//    Person(name: "Bob", age: 42)
// ]
```

Nested types render recursively with cycle detection. The type itself (not an instance) renders as its signature:

```
print("Type: `Person`")
// Type: Person #(name String, age Int)
```

### Representation rule

The backtick form uses the most specific named type the compiler has resolved for the value:

- If the compiler resolved the value to a named type (through an explicit declaration or a structural match at a call site), it renders using that type's constructor form: `Person(name: "Alice")`.
- If the value was never matched to a named type (a truly anonymous boc), it renders using the structural form: `{ name: "Alice" }`.

```
Person : { name String }

e : Person("Alice")      // type resolved to Person
f : { name: "Alice" }    // anonymous — no named type

print("`e`")   // Person(name: "Alice")
print("`f`")   // { name: "Alice" }
```

The named-type form is preferred because it preserves semantic identity: two bocs with identical fields but different names (`Person` vs `Address`) would otherwise be indistinguishable.

### Backtick as text serialization

For pure data bocs (fields only, no closures or external resources), the backtick output is valid Yz source that recreates the value:

```
file.write("`e`")   // writes: Person(name: "Alice")
// reading back and evaluating gives an equivalent Person instance
```

This makes backtick a natural text serialization format for data bocs — no separate library needed. Binary serialization is out of scope for this mechanism.

Backtick interpolation inside strings is distinct from info-string backticks at statement level (see [Info strings](Info%20strings.md)).

## Prefer interpolation over concatenation

```
// preferred
print("User: ${user.name}, age: ${user.age}")

// avoid
print("User: " + user.name + ", age: " + user.age.to_str())
```

## See Also

- [Info strings](Info%20strings.md) — backtick-delimited compile-time metadata
- [Define new types](Define%20new%20types.md) — adding `to_str` to a user type
