#feature

# String Interpolation

Embed any expression inside a string using `${}`:

```
name: "Alice"
age: 30
print("Hello, ${name}!")                    // Hello, Alice!
print("Age: ${age}, next: ${age + 1}")      // Age: 30, next: 31
```

Interpolation works inside both single- and double-quoted strings.

## Any expression is valid

```
x: 5
print("x squared: ${x * x}")               // x squared: 25

p: Person("Bob", 25)
print("Greeting: ${p.name}")               // Greeting: Bob
```

Complex expressions including method calls and closures are valid inside `${}`:

```
names: ["Alice", "Bob", "Carol"]
print("Members: ${names.join(", ")}")      // Members: Alice, Bob, Carol

print("fib(5) = ${fib(2 + 3)}")           // fib(5) = 5
```

## Prefer interpolation over concatenation

```
// preferred
print("User: ${user.name}, age: ${user.age}")

// avoid
print("User: " + user.name + ", age: " + user.age.to_str())
```

Interpolation handles type conversion automatically — no need to call `to_str()` on numeric values.

## See Also

- [Info strings](Info%20strings.md) — backtick-delimited compile-time metadata