---
type: feature
generated: { by: "oscarryz", at: 2026-10-06T10:27:00+02:00 }
---
#feature 
# Non-Word Method Invocation

When a method name is a non-word symbol (e.g. `<<`, `+`, `==`), it can be invoked without `.` or parentheses, as long as it has at least one parameter. The receiver comes first, then the method name, then the arguments.

```yz
Example: {
  // Declares the `<<` method 
  << : {
    n Int
    print(n)
  }
}

e: Example()
// invokes the `<<` method on `e`
e << 1    // same as `e.<<(1)` .  prints 1
```

## Defining non-word methods

Non-word methods are declared like any other boc variable, using the symbol as the name. They must appear inside a boc (type or singleton):

```yz
Vec: {
  x Int
  y Int
  + : {
    other Vec
    Vec(x + other.x, y + other.y)
  }
}

a: Vec(1, 2)
b: Vec(3, 4)
c: a + b    // Vec(4, 6)
```

## Precedence and associativity

All non-word methods have **equal precedence** and chains are **left-associative**. Since `a op b` is sugar for `a.op(b)`, a chain is just a sequence of method calls where the result of each call is the receiver of the next:

```yz
a ++ b ++ c        // (a ++ b) ++ c  ==  a.++(b).++(c)
a ++ (b ++ c)      // parentheses force the other grouping
1 + 2 * 3          // (1 + 2) * 3 = 9, there is no precedence table
```

This holds for every symbol, including mixed ones, so `x == 0 ? { ... }, { ... }` reads as `(x == 0) ? { ... }, { ... }`. The only exception is unary `-`, which binds to the expression right after it: `-n + 1` is `(-n) + 1`.

Associativity is a property of how the parser groups a chain, not of the method. A type's `++` need not be associative as an operation (`(a ++ b) ++ c` and `a ++ (b ++ c)` may differ); the language always picks the left grouping.

## Relation to trailing-block syntax

Trailing-block syntax (omitting `()` for a single boc argument) works for word-named methods. Non-word invocation works for symbol-named methods. They are complementary. See [Trailing block syntax](Trailing%20block%20syntax.md).
