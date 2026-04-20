#feature

When the last (or only) argument to a method is a block literal, the enclosing parentheses can be omitted. The block is written directly after the method name, separated by a space.

```yz
// With parentheses (always valid)
list.filter({ item Int; item > 10 })

// Without parentheses — trailing block syntax
list.filter { item Int; item > 10 }
```

Both forms are identical. The parser recognises a `{` immediately following a member access (`.method`) on the same line and treats it as a single-argument call.

## Rules

- The `{` must be on the **same line** as the method name. If a newline separates them, ASI inserts a semicolon and the block becomes a new statement, not an argument.
- Trailing-block syntax only applies after a **member access** (`receiver.method`). A free function call still requires parentheses: `filter({ block })`.
- The block is always the **sole** argument in this syntax. If you need additional arguments, use the parenthesised form.

## Examples

```yz
// Filter elements from an array
numbers: [1, 2, 3, 10, 20]
big: numbers.filter { n Int; n > 5 }

// Apply a side effect to each element
big.each { n Int; print(n) }

// Chaining works naturally
[1, 2, 3, 4, 5]
    .filter { n Int; n > 2 }
    .each   { n Int; print(n) }
```

Note that chaining across lines requires the `.` to start the continuation line (or use explicit parentheses), because ASI would otherwise insert a semicolon after the closing `}`.

## Relation to non-word invocation

Non-word method invocation (`e << 1`) is a separate but complementary feature: it eliminates both the `.` and the parentheses for non-word method names. Trailing-block syntax eliminates the parentheses for word method names when passing a single block argument.
