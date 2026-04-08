# String Interpolation

Embed any expression inside a string using backticks:

```yz
name: "Alice"
age: 30
print("Hello, `name`!")          // Hello, Alice!
print("Age: `age`, next: `age + 1`")  // Age: 30, next: 31
```

Interpolation works inside both single- and double-quoted strings.

## Any expression is valid

```yz
x: 5
print("x squared: `x * x`")     // x squared: 25

p: Person("Bob", 25)
print("Greeting: `p.name`")     // Greeting: Bob
```

## Prefer interpolation over concatenation

```yz
// preferred
print("User: `user.name`, age: `user.age`")

// avoid
print("User: " + user.name + ", age: " + user.age.to_str())
```

Interpolation handles type conversion automatically — no need to call `to_str()` on numeric values.
