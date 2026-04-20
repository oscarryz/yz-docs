#feature 
# Strings

Both single and double quotes delimit string literals:

```yz
greeting: "Hello, World!"
name: 'Alice'
```

Multi-line strings span lines naturally:

```yz
message: "
  This is a
  multi-line string
"
```

## String interpolation

Use backticks inside a string to embed expressions:

```yz
name: "Alice"
age: 30
print("Name: `name`, Age: `age`")   // Name: Alice, Age: 30
print("Next year: `age + 1`")       // Next year: 31
```

See [String interpolation](String%20interpolation.md) for details.

## String methods

```yz
s: "Hello, World!"
s.length()              // 13
s.to_upper()            // "HELLO, WORLD!"
s.to_lower()            // "hello, world!"
s.contains("World")     // true
s.has_prefix("Hello")   // true
s.has_suffix("!")       // true
s.trim()                // removes leading/trailing whitespace
s.to_str()              // returns the string itself
```

## String concatenation

Use `+` to concatenate strings, though interpolation is preferred:

```yz
full_name: first + " " + last
```
