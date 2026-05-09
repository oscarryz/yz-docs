#spec 
# 7. Control Flow

This chapter defines Yz's control flow mechanisms. Yz has no built-in control flow keywords for conditionals or loops — they are expressed through methods and boc composition.

## 7.1 Design Principle

Yz minimizes special syntax. Instead of dedicated `if`/`else`/`for`/`while` constructs, Yz uses:

- **`?` method on `Bool`** — conditional branching
- **`match` expression** — pattern matching (the only control-flow keyword)
- **Methods on ranges and collections** — iteration
- **`break`, `continue`, `return`** — early exit keywords

## 7.2 Conditional — The `?` Method

`?` is a method on `Bool` that takes two boc arguments and executes one based on the boolean value:

```yz
condition ? { true_branch }, { false_branch }
```

### Semantics

- If the receiver is `true`, the first boc is executed and its value is returned
- If the receiver is `false`, the second boc is executed and its value is returned
- Both bocs must have compatible return types

### Examples

```yz
// Simple conditional
(x > 0) ? { "positive" }, { "non-positive" }

// With side effects
(logged_in) ? {
    show_dashboard()
}, {
    redirect_to_login()
}

// Nested
(x > 0) ? {
    (x > 100) ? { "large" }, { "small" }
}, {
    "non-positive"
}

// As expression
label: (score >= 90) ? { "A" }, { "B" }
```

### Single-Branch Conditional

To execute something only when true, use an empty false branch:

```yz
(debug) ? { print("Debug mode") }, { }
```

## 7.3 Match Expressions

`match` is the only control-flow keyword in Yz. It comes in two forms.

### Condition Match (Cond-Style)

Evaluates boolean conditions in order and executes the first matching branch:

```yz
grade: match {
    score >= 90 => "A"
}, {
    score >= 80 => "B"
}, {
    score >= 70 => "C"
}, {
    "F"    // Default — no condition
}
```

**Rules:**
- Conditions are tested top-to-bottom
- The first `true` condition's boc is executed
- A boc without `=>` is the **default** (always matches)
- If no branch matches and there's no default, the result is `Unit`

### Variant Match

Discriminates between type variants of a specific value:

```yz
match result {
    Result.Ok => print("Value: ${result.value}")
}, {
    Result.Err => print("Error: ${result.error}")
}
```

**Rules:**
- The subject expression is evaluated once
- Each branch names a variant constructor
- The **runtime discriminant tag** (see §4.5) determines which branch executes
- Inside a matched branch, the variant's fields are accessible on the subject
- If no branch matches, the result is `Unit`

### Match with `continue`

The `continue` keyword inside a match branch causes evaluation to fall through to the **next** branch:

```yz
match {
    x > 0 => {
        print("positive")
        continue    // Also check next condition
    }
}, {
    x > 10 => print("and greater than 10")
}, {
    print("done")
}
```

### Match as Expression

Match returns a value — the value of the executed branch:

```yz
description: match response {
    Success => "OK"
}, {
    NotFound => "Not found"
}, {
    Timeout => "Timed out"
}, {
    "Unknown"
}
```

## 7.4 Iteration

Yz has no loop syntax. Iteration is achieved through methods on ranges and collections.

### Range Iteration

```yz
1.to(10).each({ i Int
    print("${i}")
})
```

- `1.to(10)` creates a `Range` value
- `.each(boc)` invokes the boc for each element

### Collection Iteration

```yz
names: ["Alice", "Bob", "Charlie"]
names.each({ name String
    print("Hello, ${name}!")
})
```

### While Loop

`while` is a function (not syntax) that takes two bocs — a condition and a body:

```yz
count: 0
while({ count < 10 }, {
    print("${count}")
    count = count + 1
})
```

- The condition boc is evaluated before each iteration
- If it returns `true`, the body boc is executed
- If `false`, iteration stops

### Indexed Iteration

```yz
names.each_with_index({ name String, i Int
    print("${i}: ${name}")
})
```

### Transformation

```yz
doubled: numbers.map({ n Int; n * 2 })
evens: numbers.filter({ n Int; n % 2 == 0 })
sum: numbers.reduce(0, { acc Int, n Int; acc + n })
```

## 7.5 Early Exit

### `return`

Exits the current boc and optionally returns a value:

```yz
find_first: {
    items [String]
    target String
    items.each({ item String
        (item == target) ? { return item }, { }
    })
    ""   // Not found
}
```

### `break`

Exits the current iteration (used inside `.each`, `while`, etc.):

```yz
1.to(100).each({ i Int
    (i > 10) ? { break }, { }
    print("${i}")
})
```

### `continue`

In iteration context: skips to the next iteration. In match context: falls through to the next branch (see §7.3).

```yz
1.to(20).each({ i Int
    (i % 2 == 0) ? { continue }, { }
    print("${i}")    // Only odd numbers
})
```

## 7.6 Control Flow Summary

| Mechanism | How |
|-----------|-----|
| If/else | `condition ? { then }, { else }` |
| Multi-way conditional | `match { cond1 => a }, { cond2 => b }, { default }` |
| Pattern matching | `match value { Variant1 => a }, { Variant2 => b }` |
| For loop | `range.each({ i; ... })` or `collection.each({ item; ... })` |
| While loop | `while({ condition }, { body })` |
| Early return | `return [value]` |
| Break from loop | `break` |
| Skip iteration | `continue` |
| Fallthrough in match | `continue` |
