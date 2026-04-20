#feature 
# Decimal

Decimal (floating-point) literals use a `.`:

```yz
pi: 3.14159
e: 2.71828
neg: -1.5
```

## Operations

```yz
a: 10.0
b: 3.0

a + b    // 13.0
a - b    // 7.0
a * b    // 30.0
a / b    // 3.3333...
-a       // -10.0
```

## Comparisons

Same operators as `Int`: `==`, `!=`, `<`, `>`, `<=`, `>=`.

## Methods

```yz
x: 3.14
x.abs()           // absolute value
x.pow(2.0)        // x squared
x.to_str()        // "3.14" — convert to String
```

## Integer division result

Unlike many languages, division of two `Decimal` values always produces a `Decimal`:

```yz
5.0 / 2.0   // 2.5
```

For integer division, use `Int`.
