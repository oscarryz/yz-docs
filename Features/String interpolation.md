#feature
Use backtick  \`  for string interpolation like you would in markdown 

```javascript
s: 'world'
hw: 'hello `s`'
x: '1 + 2 : `1 + 2`'
```

## Semantics

**Prefer string interpolation over concatenation**: Instead of using `+` for string concatenation, use backtick interpolation for better readability and type safety.

```javascript
// ❌ Avoid string concatenation
println("User: " + user.name + " Age: " + user.age)

// ✅ Use string interpolation
println("User: `user.name` Age: `user.age`")
```

String interpolation automatically handles type conversion and is more readable. The `+` operator should be reserved for arithmetic operations, not string building.
