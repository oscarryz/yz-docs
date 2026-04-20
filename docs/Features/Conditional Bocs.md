#feature
#pattern-matching

Yz provides two ways to branch on conditions: a simple boolean conditional and `match`, which evaluates a list of conditional bocs in order.

## Conditional boc

A conditional boc has a condition and an action, separated by `=>`:

```yz
{ 1 < 2 => "One is lower than two" }

{ is_monday() =>
    print("wake up")
    print("have breakfast")
    print("go to work")
}
```

The simpler boolean form uses the `?` method on `Bool`:

```yz
x: 5
x > 3 ? {
    print("big")
}, {
    print("small")
}
```

## match

`match` takes a comma-separated list of conditional bocs and evaluates them in order, returning the result of the first matching branch. The last boc in the list acts as the default (no condition required):

```yz
score: 85
grade: match {
    score >= 90 => "A"
}, {
    score >= 80 => "B"
}, {
    score >= 70 => "C"
}, {
    "F"
}
print(grade)
```

`match` can be used anywhere an expression is expected.

### Factorial example

```yz
factorial #(n Int, Int) {
    match {
        n == 0 => 1
    }, {
        n > 0 => n * factorial(n - 1)
    }
}
```

### Cascade evaluation with `continue`

If more than one branch should execute, use `continue` to fall through to the next branch:

```yz
n Int
match {
    n % 3 == 0 => print("Fizz")
    continue
}, {
    n % 5 == 0 => print("Buzz")
}, {
    print("`n`")
}
```

## Match on type variants

When `match` is given a subject, each branch can pattern-match on a type variant constructor:

```yz
x Option(String)   // could be Some or None

// Execute statements based on variant
match x {
    Some => print("The value is `x.value`")
}, {
    None => print("There was no value")
}

// Return a value based on variant
value: match x {
    Some => x.value
}, {
    "No value"
}
```

The last branch can omit its condition when the variants are exhaustive.