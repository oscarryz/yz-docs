#feature

# Yz Keyword Constraint

## The Design Invariant

Yz is built on a single composable concept — the **boc** (block of code). Every feature
of the language is expressed as a boc, a composition of bocs, or a call between bocs.

The goal was zero keywords. The result is four. Every primitive on this page represents
a place where the boc model reaches its limit — where something categorically different
from data transformation is required. Each one is documented here with the reason it
exists and why it could not be expressed as a boc.

This list is a design invariant. Future features must justify themselves against it. A
new keyword or primitive is a cost, not a feature.

---

## Keywords

### `match`

**What it does:** Evaluates a list of conditional bocs in order, returning the result of
the first matching branch. When given a subject, dispatches on type variant constructors
with exhaustiveness checking.

**Why it cannot be a boc:** Two reasons. First, exhaustiveness checking requires the
compiler to know the complete set of constructors for a variant type — a closed-world
assumption that a boc calling other bocs cannot enforce. Second, variant dispatch
requires the compiler to narrow the type of the subject inside each branch (flow typing)
— something only the compiler can do at the point of match.

The boolean form of `match` comes close to being expressible as a boc — a list of
conditional bocs evaluated in order. The variant form is what requires the keyword.

```js
// boolean form — evaluates conditions in order
grade : match {
    score >= 90 => "A"
}, {
    score >= 80 => "B"
}, {
    "F"
}

// variant form — exhaustive dispatch on constructors
match result {
    Ok  => print(result.value)
}, {
    Err => print(result.error)
}
```

See also: [Conditional Bocs](./yz-conditional-bocs.md) ·
[Type Variants](./yz-type-variants.md)

---

### `return`

**What it does:** Exits the current boc immediately, optionally producing a value.

**Why it cannot be a boc:** A boc calling another boc cannot escape its own execution
context. `return` is a jump — it unwinds the call stack to the boundary of the enclosing
boc. No amount of boc composition can replicate a non-local jump without the compiler's
cooperation.

In practice, `return` is often unnecessary — the last expression in a boc is its value.
It exists for early exit when a condition is met before the end of the body.

```yz
find : {
    items []T
    target T
    items.forEach({ item T
        item == target ? { return item }
    })
    None()
}
```

---

### `break`

**What it does:** Exits the nearest enclosing loop immediately.

**Why it cannot be a boc:** Same reason as `return` — it is a non-local jump that
escapes a loop context. A boc cannot reach outside itself to terminate a loop it did not
start.

```yz
items.forEach({ item T
    item.is_invalid() ? { break }
    process(item)
})
```

---

### `continue`

**What it does:** In a loop, advances to the next iteration. In a `match`, falls through
to the next branch.

**Why it cannot be a boc:** The loop form is a non-local jump — same argument as
`break`. The match fallthrough form could theoretically be expressed differently but
sharing one keyword for both uses keeps the language smaller.

```yz
// loop form — skip to next iteration
items.forEach({ item T
    item.is_processed() ? { continue }
    process(item)
})

// match fallthrough — continue to next branch
match {
    n % 3 == 0 => print("Fizz")
    continue
}, {
    n % 5 == 0 => print("Buzz")
}, {
    print("`n`")
}
```

---

## Everything Else Is A boc

The following features that are primitives or keywords in other languages are regular
bocs in Yz:

| Feature                  | How Yz expresses it                                        |
| ------------------------ | ---------------------------------------------------------- |
| `if` / `else`            | `if_true_if_false` method on `Bool`                        |
| `while` loop             | Recursive boc with `Bool` condition                        |
| Operators `+`, `>`, etc. | Methods on boc — `Int` has `+ #(Int, Int)`                 |
| New user-defined types   | Uppercase-starting boc                                     |
| Generics                 | Single uppercase identifier — `T`, `U`, `V`                |
| Imports                  | FQN reference or local alias — `net : some.util.place.net` |
| Annotations / macros     | Info strings — bocs themselves                             |
| Dependency management    | Info strings                                               |
| Associated types         | A type-of-type boc                                         |

---

## The Design Constraint Going Forward

Any proposed addition to Yz must answer:

1. **Can this be expressed as a boc?** If yes, it should be.
2. **If not, which of the four keywords covers it?** If one does, use it.
3. **If neither, is the new primitive genuinely unavoidable?** If yes, document it here
   with the same rigour as the existing primitives.

The list on this page should grow only under significant pressure.