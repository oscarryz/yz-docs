#feature

## Design Philosophy

These principles guide every decision in the language, from syntax to semantics to error messages.

**Homoiconicity** — the output of a program looks like its source. A value rendered as a string resembles the expression that created it. What you write is what you see.

**Zero to no visible magic** — the programmer should never encounter behaviour with no apparent source. Internal compiler machinery is acceptable; invisible methods, implicit coercions, or auto-generated behaviour the user cannot read in the source are not. If something happens, there should be a place to look for why.

**Make things obvious** — when there are two ways to read a piece of code, the language should make the correct reading the natural one. Ambiguity is a cost.

**Simple exterior, complex interior** — simplicity is a user-facing property, not an implementation constraint. A simple syntax can be backed by a sophisticated runtime. The programmer pays the conceptual cost of the language's surface; the compiler pays the cost of making that surface work.

**Ergonomics matter** — a correct language that is unpleasant to write is a failed design. Friction accumulates. The common case should feel effortless; the uncommon case should be possible without ceremony.

**The syntax is the interface** — the way code looks is not decoration. Syntax communicates intent, scope, ownership, and structure. Every syntactic choice is a decision about what the programmer should notice and what they should be able to ignore.

## One Concept

These principles converge on a single structural decision: Yz is built on one composable concept — the **boc** (block of code). Every feature of the language is expressed as a boc, a composition of bocs, or a call between bocs.

No special forms. No hidden dispatch. No context-sensitive syntax. A boc is a boc whether it holds data, behaviour, a type, a concurrent actor, or a module. The concurrency model is powerful precisely because it does not look special — it is the same boc model with a runtime that takes concurrency seriously underneath.

The goal was zero keywords. The result is four.

## Why These Four Keywords

A keyword is a place where the boc model reaches its limit — where something categorically different from data transformation is required and no boc composition can fill the gap. Each keyword on this page represents that moment.

A new keyword is a cost, not a feature. The list below is a design invariant: it should grow only under significant pressure, and any candidate must answer the same question each of these four had to answer — *why can this not be a boc?*

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

```yz
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

See also: [Conditional Bocs](docs/Features/Conditional%20Bocs.md) ·  [Type variants](docs/Features/Type%20variants.md)


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

## Everything Else Is A Boc

The following are primitives or keywords in other languages. In Yz they are regular bocs:

| Feature | Other languages | Yz |
| ------- | --------------- | -- |
| Conditionals | `if` / `else` | `?` method on `Bool` — `cond ? { then }, { else }` |
| Loops | `while`, `for` | `while` boc — `while { cond } { body }` |
| Operators | `+`, `>`, `==`, etc. | Methods on types — `Int` has `+ #(other Int, Int)` |
| User-defined types | `class`, `struct`, `record` | Uppercase-starting boc — `Person : { name String }` |
| Generics | `<T>`, `[T]`, template params | Single uppercase identifier — `T`, `E`, `K` |
| Imports | `import`, `use`, `require` | FQN reference — `net : some.util.net` |
| Annotations / macros | `@annotation`, `#[attr]` | Info strings — `` `annotation` `` before any declaration |
| Dependency management | `package.json`, `go.mod` | Info strings in source |
| Associated types | `type Item = ...` | A type-of-type boc |

## Adding to the Language

Any proposed addition to Yz must answer:

1. **Can this be expressed as a boc?** If yes, it should be.
2. **If not, which of the four keywords covers it?** If one does, use it.
3. **If neither, is the new primitive genuinely unavoidable?** If yes, document it here
   with the same rigour as the existing four.

Then check it against the design philosophy: does it introduce visible magic? Does it add
surface complexity? Does the syntax make the intent obvious? If any answer is uncomfortable,
the design needs more work.
