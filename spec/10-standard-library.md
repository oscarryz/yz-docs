#spec 
# 10. Standard Library

This chapter defines the built-in types, their methods, and core functions provided by the Yz standard library.

## 10.1 Overview

The standard library provides:

1. **Built-in types** — `Int`, `Decimal`, `String`, `Bool`, `Unit`
2. **Collection types** — `Array` (`[T]`), `Dictionary` (`[K:V]`)
3. **Common variant types** — `Option(T)`, `Result(T, E)`
4. **Core functions** — `print`, `while`, `info`
5. **Utility modules** — `time`, `io`, `math`

## 10.2 Int

Integer type with arbitrary precision.

### Methods

| Method | Signature | Description |
|--------|-----------|-------------|
| `+` | `#(other Int, Int)` | Addition |
| `-` | `#(other Int, Int)` | Subtraction (binary) |
| `-` | `#(Int)` | Negation (unary) |
| `*` | `#(other Int, Int)` | Multiplication |
| `/` | `#(other Int, Int)` | Integer division |
| `%` | `#(other Int, Int)` | Modulo |
| `<` | `#(other Int, Bool)` | Less than |
| `>` | `#(other Int, Bool)` | Greater than |
| `<=` | `#(other Int, Bool)` | Less or equal |
| `>=` | `#(other Int, Bool)` | Greater or equal |
| `==` | `#(other Int, Bool)` | Equality |
| `!=` | `#(other Int, Bool)` | Inequality |
| `to` | `#(end Int, Range)` | Create range `[self, end)` |
| `to_string` | `#(String)` | String representation |
| `to_decimal` | `#(Decimal)` | Convert to Decimal |

## 10.3 Decimal

Decimal floating-point type.

### Methods

| Method | Signature | Description |
|--------|-----------|-------------|
| `+` | `#(other Decimal, Decimal)` | Addition |
| `-` | `#(other Decimal, Decimal)` | Subtraction |
| `-` | `#(Decimal)` | Negation (unary) |
| `*` | `#(other Decimal, Decimal)` | Multiplication |
| `/` | `#(other Decimal, Decimal)` | Division |
| `<` | `#(other Decimal, Bool)` | Less than |
| `>` | `#(other Decimal, Bool)` | Greater than |
| `<=` | `#(other Decimal, Bool)` | Less or equal |
| `>=` | `#(other Decimal, Bool)` | Greater or equal |
| `==` | `#(other Decimal, Bool)` | Equality |
| `!=` | `#(other Decimal, Bool)` | Inequality |
| `to_string` | `#(String)` | String representation |
| `to_int` | `#(Int)` | Truncate to Int |

## 10.4 String

UTF-8 encoded text. Strings support interpolation with backticks (see §1.10).

### Methods

| Method | Signature | Description |
|--------|-----------|-------------|
| `+` | `#(other String, String)` | Concatenation |
| `==` | `#(other String, Bool)` | Equality |
| `!=` | `#(other String, Bool)` | Inequality |
| `length` | `#(Int)` | Number of characters |
| `at` | `#(index Int, String)` | Character at index |
| `slice` | `#(start Int, end Int, String)` | Substring |
| `contains` | `#(sub String, Bool)` | Contains substring |
| `starts_with` | `#(prefix String, Bool)` | Starts with prefix |
| `ends_with` | `#(suffix String, Bool)` | Ends with suffix |
| `split` | `#(sep String, [String])` | Split by separator |
| `trim` | `#(String)` | Remove leading/trailing whitespace |
| `to_int` | `#(Option(Int))` | Parse as integer |
| `to_decimal` | `#(Option(Decimal))` | Parse as decimal |
| `to_upper` | `#(String)` | Uppercase |
| `to_lower` | `#(String)` | Lowercase |

## 10.5 Bool

Boolean type. `true` and `false` are standard library constants, not keywords.

### Methods

| Method | Signature | Description |
|--------|-----------|-------------|
| `&&` | `#(other Bool, Bool)` | Logical AND |
| `\|\|` | `#(other Bool, Bool)` | Logical OR |
| `?` | `#(then #(T), else #(T), T)` | Conditional — execute one of two bocs |
| `==` | `#(other Bool, Bool)` | Equality |
| `!=` | `#(other Bool, Bool)` | Inequality |
| `to_string` | `#(String)` | `"true"` or `"false"` |

## 10.6 Unit

The `Unit` type represents the absence of a meaningful value. Bocs that perform side effects without returning a value have return type `Unit`.

`Unit` has no methods.

## 10.7 Array — `[T]`

Ordered, indexed collection of elements.

### Construction

```yz
nums: [1, 2, 3]         // [Int]
empty: [String]()        // Empty [String]
```

### Methods

| Method | Signature | Description |
|--------|-----------|-------------|
| `length` | `#(Int)` | Number of elements |
| `at` | `#(index Int, T)` | Element at index |
| `<<` | `#(item T)` | Append element |
| `each` | `#(f #(T))` | Iterate over elements |
| `each_with_index` | `#(f #(T, Int))` | Iterate with index |
| `map` | `#(f #(T, U), [U])` | Transform elements |
| `filter` | `#(pred #(T, Bool), [T])` | Filter elements |
| `reduce` | `#(init U, f #(U, T, U), U)` | Fold/reduce |
| `contains` | `#(item T, Bool)` | Contains element (uses `==`) |
| `first` | `#(Option(T))` | First element |
| `last` | `#(Option(T))` | Last element |
| `slice` | `#(start Int, end Int, [T])` | Sub-array |
| `==` | `#(other [T], Bool)` | Element-wise equality |

## 10.8 Dictionary — `[K:V]`

Key-value collection.

### Construction

```yz
ages: ["Alice": 30, "Bob": 25]   // [String:Int]
empty: [String:Int]()              // Empty
```

### Methods

| Method | Signature | Description |
|--------|-----------|-------------|
| `at` | `#(key K, Option(V))` | Get value for key |
| `set` | `#(key K, value V)` | Set key-value pair |
| `has` | `#(key K, Bool)` | Key exists |
| `remove` | `#(key K, Option(V))` | Remove and return |
| `keys` | `#([K])` | All keys |
| `values` | `#([V])` | All values |
| `each` | `#(f #(K, V))` | Iterate over entries |
| `length` | `#(Int)` | Number of entries |
| `==` | `#(other [K:V], Bool)` | Key-value equality |

## 10.9 Range

Created by `Int.to()`. Represents a half-open interval `[start, end)`.

### Methods

| Method | Signature | Description |
|--------|-----------|-------------|
| `each` | `#(f #(Int))` | Iterate over range |
| `contains` | `#(n Int, Bool)` | Check membership |
| `to_array` | `#([Int])` | Convert to array |

## 10.10 Option(T)

Variant type for optional values.

```yz
Option: {
    T
    Some(value T)
    None()
}
```

### Usage

```yz
find: {
    items [String]
    target String
    result Option(String) = Option.None()
    items.each({ item String
        (item == target) ? {
            result = Option.Some(item)
        }, { }
    })
    result
}
```

## 10.11 Result(T, E)

Variant type for operations that may fail.

```yz
Result: {
    T, E
    Ok(value T)
    Err(error E)
}
```

### Usage

```yz
parse_int: {
    s String
    // ...
    match {
        valid => Result.Ok(parsed_value)
    }, {
        Result.Err("invalid integer: ${s}")
    }
}
```

## 10.12 Core Functions

### `print`

```yz
print #(value String)
```

Prints a string to standard output followed by a newline. This is an **IO operation** that triggers thunk materialization (see §8.3).

### `while`

```yz
while #(condition #(Bool), body #())
```

Repeatedly evaluates the condition boc; if `true`, executes the body boc. Stops when condition returns `false`.

### `info`

```yz
info #(target, InfoData)
```

Retrieves the info string metadata attached to any element (see §1.14).

## 10.13 Utility Modules

### `time`

| Function | Description |
|----------|-------------|
| `time.now()` | Current timestamp (milliseconds) |
| `time.sleep(ms Int)` | Sleep for duration |

### `io`

| Function | Description |
|----------|-------------|
| `io.read_line()` | Read line from stdin |
| `io.read_file(path String)` | Read file contents → `Result(String, Error)` |
| `io.write_file(path String, content String)` | Write file → `Result(Unit, Error)` |

### `math`

| Function | Description |
|----------|-------------|
| `math.abs(n Int)` | Absolute value |
| `math.max(a Int, b Int)` | Maximum |
| `math.min(a Int, b Int)` | Minimum |
| `math.pow(base Int, exp Int)` | Exponentiation |

> **Note:** The standard library will evolve. This chapter documents the **minimum viable** set for v0.1.
