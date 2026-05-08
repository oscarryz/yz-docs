#feature 
# return, break, and continue

These are three of the four reserved keywords in Yz. Each one performs a non-local jump —
something that cannot be expressed as a boc calling another boc. The fourth keyword,
`match`, is documented in [Conditional Bocs](Conditional%20Bocs.md).


See also: [Language Primitives](Language%20Primitives.md) · [Concurrency](Concurrency.md)

---

## return

By default, a boc returns the value of its last expression. Use `return` to exit early
with a value, or without one.

```yz
check : {
    age Int
    age < 21 ? {
        return "too young"
    }, {}
    "welcome"
}

print(check(20))   // too young
print(check(25))   // welcome
```

### Named Bocs vs Anonymous Bocs

`return` exits the nearest **named** boc. Anonymous bocs — those passed as arguments to
other bocs — are transparent to `return`. This is what makes guard clauses work:

```yz
foo : {                          // named — return targets this
    u User
    u.age < 18 ? { return }, {}  // exits foo, not the anonymous boc
    // ... rest of code only runs if u.age >= 18
}
```

Without this rule, `return` would exit only the anonymous boc passed to `?`, which is
not useful — the anonymous boc's result is already its last expression. Guard clauses
require `return` to reach the enclosing named boc.

The same rule applies inside `forEach`, `map`, and any other boc that takes callbacks:

```yz
find_first : {                        // named — return targets this
    items [Item]
    target Item
    items.forEach({ item Item         // anonymous — transparent to return
        item == target ? { return item }, {}
    })
    None()
}
```

### You Never Need To Exit An Anonymous Boc Early

Anonymous bocs are always used in one of three ways:

- **Producing a value** — the last expression is the value. No `return` needed.
- **Loop body** — use `continue` to skip an iteration or `break` to exit the loop.
- **Conditional branch** — the branch evaluates to its last expression naturally.

```yz
// value-producing — last expression is the result
result : items.map({ item Item
    item.invalid() ? { Item.Empty() }, { transform(item) }
})

// loop body — break and continue handle early exit
items.forEach({ item Item
    item.done() ? { break }, {}
    item.skip() ? { continue }, {}
    process(item)
})
```

There is no case where `return` from an anonymous boc is meaningful that is not already
covered by `break`, `continue`, or the last-expression rule.

### return And Concurrency

When a named boc exits via `return`, all values it acquired are released atomically.
Pending invocations waiting on those values are free to proceed.

---

## break

`break` exits the innermost enclosing loop immediately:

```yz
max_from_list : {
    list [Int]
    m : 0
    list.forEach({ item Int
        item < 0 ? { break }, {}   // stops iteration immediately
        m = max(m, item)
    })
    m
}
```

`break` targets the innermost loop, not the enclosing named boc. To exit the enclosing
named boc use `return`.

### break And Concurrency

`break` exits a loop inside a running invocation. The invocation itself continues
after the loop — no values are released. Only `return` releases acquired values.

---

## continue

`continue` has two distinct uses depending on context.

### In Loops

`continue` skips the rest of the current iteration and moves to the next:

```yz
list.forEach({ n Int
    n % 2 == 0 ? { continue }, {}   // skip even numbers
    print(n)                          // prints only odd numbers
})
```

### In match expressions

`continue` inside a `match` arm falls through to the next arm:

```yz
n Int
match {
    n % 3 == 0 => {
        print("Fizz")
        continue        // fall through to check next condition
    }
}, {
    n % 5 == 0 => print("Buzz")
}, {
    print("`n`")
}
```

Both uses share one keyword because the underlying operation is the same — advance past
the current case and evaluate the next one.

---

## Summary

| Keyword | Targets | Releases acquired values |
|---|---|---|
| `return` | Nearest named boc | ✅ yes — atomically |
| `break` | Nearest enclosing loop | ❌ invocation continues |
| `continue` (loop) | Current loop iteration | ❌ invocation continues |
| `continue` (match) | Current match arm | ❌ invocation continues |