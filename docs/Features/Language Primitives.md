#feature

# Yz Language Primitives

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

## Compiler Seam

### `compiler.read(b)` and `compiler.insert(b)`

**What it does:** `compiler.read` returns the `Boc` instance the compiler created for a
boc definition. `compiler.insert` splices a `Boc` into the current compilation context.
Both are available only inside `Compile` typed slots.

**Why it cannot be a regular boc:** The compiler builds structural descriptions of bocs
as it parses and infers them. Making that information available as a value requires a
controlled seam between the language and the compilation phase. `compiler` is that seam
— explicit, localized, and available only where compile-time execution is opted into via
the `Compile` type.

Everything inside a `Compile` slot is regular Yz. `compiler` is the only thing that is
not.

See also: [Structural Reflection](./yz-structural-reflection.md) ·
[Compile-Time Bocs](./yz-compile-time-bocs.md)

---

## Meta-Type

### `Boc`

**What it does:** Represents the structural description of any boc — its name,
fields, methods, type parameters, infostrings, and literal source — as a regular Yz
value. The compiler creates a `Boc` instance for every boc definition automatically.

**Why it cannot be a regular boc:** Every other boc is defined by the developer. `Boc`
instances are created by the compiler from source code. The compiler must populate them
— the developer cannot. This is the minimum meta-level concept required for structural
reflection, compile-time code generation, and constraint inference to work.

`Boc` is not a meta-type in the sense of living outside the type system. It is a regular
boc with a known structure. Its instances behave like any other value — they can be
passed, iterated, serialized, and sent across a wire. The only special thing about `Boc`
is who creates its instances.

```yz
Boc : {
    name         String
    instantiable Bool
    fields       [Boc]
    methods      [Boc]
    type_params  [Boc]
    infostrings  [String]
    source       #()
}
```

See also: [Structural Reflection](./yz-structural-reflection.md)

---

## Naming Convention

### Single Uppercase Letter — Type Parameter

**What it does:** A single uppercase letter identifier (`T`, `U`, `V`, etc.) in a boc
body declares a type parameter — a placeholder `Boc` that is filled with a concrete type
at instantiation time.

**Why it is a compiler primitive:** The compiler uses this rule to identify type
parameters and create placeholder `Boc` instances for them. It is not merely a
convention — a multi-letter identifier in the same position is a syntax error, not a
type parameter.

```yz
Box : { T }    // T is a type parameter — placeholder Boc, valid
Box : { Tom }  // syntax error — Tom is not a type parameter
```

This convention fits within the existing identifier system without ambiguity:

| Shape | Meaning |
|---|---|
| `lowercase` | singleton |
| `Uppercase` multi-letter | instantiable boc |
| Single uppercase letter `T` | type parameter |

See also:  [Generics - Type Parameters](Generics%20-%20Type%20Parameters.md)


---

## Everything Else Is A boc

The following features that are primitives or keywords in other languages are regular
bocs in Yz:

| Feature | How Yz expresses it |
|---|---|
| `if` / `else` | `if_true_if_false` method on `Bool` |
| `while` loop | Recursive boc with `Bool` condition |
| Operators `+`, `>`, etc. | Methods on boc — `Int` has `+ #(Int, Int)` |
| Imports | `mix` — composition via `Compile` and `Boc` |
| Generics constraints | Inferred from usage — no declaration needed |
| Annotations / macros | `Compile` typed slot — a boc like any other |
| Serialization / codegen | Library `Compile` implementations |
| Dependency management | `project` boc with a `Compile` slot |
| Reflection / introspection | `Boc` instances via `compiler.read()` |
| Associated types | Type slots — a `Boc` placeholder inside a boc |

---

## The Design Constraint Going Forward

Any proposed addition to Yz must answer:

1. **Can this be expressed as a boc?** If yes, it should be.
2. **If not, which of the four keywords covers it?** If one does, use it.
3. **If neither, is the new primitive genuinely unavoidable?** If yes, document it here
   with the same rigour as the existing primitives.

The list on this page should grow only under significant pressure.