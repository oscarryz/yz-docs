#feature 
# Int

Integer literals:

```yz
n: 42
m: -1
```

## Operations

```yz
a: 10
b: 3

a + b   // 13
a - b   // 7
a * b   // 30
a / b   // 3   (integer division)
a % b   // 1   (remainder)
-a      // -10 (negation)
```

## Comparisons

```yz
a == b   // false
a != b   // true
a < b    // false
a > b    // true
a <= b   // false
a >= b   // true
```

## Conversion and range

```yz
n.to_str()    // "10" — convert to String
n.abs()       // absolute value
n.to(m)       // Range from n to m (exclusive), for iteration
```

## Range iteration

```yz
1.to(10).each({ i Int
  print(i)
})
// prints 1 through 9
```
