# Non-Word Method Invocation

When a method name is a non-word symbol (e.g. `<<`, `+`, `==`), it can be invoked without `.` or parentheses, as long as it has at least one parameter. The receiver comes first, then the method name, then the arguments.

```yz
Example: {
  <<: {
    n Int
    print(n)
  }
}

e: Example()
e << 1    // same as e.<<(1) — prints 1
```

## Defining non-word methods

Non-word methods are declared like any other boc variable, using the symbol as the name. They must appear inside a boc (type or singleton):

```yz
Vec: {
  x Int
  y Int
  +: {
    other Vec
    Vec(x + other.x, y + other.y)
  }
}

a: Vec(1, 2)
b: Vec(3, 4)
c: a + b    // Vec(4, 6)
```

## Operators as methods

All arithmetic and comparison operators are methods on the built-in types. For example, `a + b` is sugar for `a.+(b)`, which the compiler translates to `a.plus(b)` (using the symbol name convention).

Built-in symbol → method name mapping:
- `+`  → `plus`
- `-`  → `minus`
- `*`  → `star`
- `/`  → `slash`
- `%`  → `percent`
- `==` → `eqeq`
- `!=` → `neq`
- `<`  → `lt`
- `>`  → `gt`
- `<=` → `lteq`
- `>=` → `gteq`
- `&&` → `ampamp`
- `||` → `pipepipe`
- `?`  → `qm`

## Relation to trailing-block syntax

Trailing-block syntax (omitting `()` for a single boc argument) works for word-named methods. Non-word invocation works for symbol-named methods. They are complementary. See [Trailing block syntax](Trailing%20block%20syntax.md).
