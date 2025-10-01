#example

https://arturo-lang.io/playground/?example=factorial%20-%20recursive

```js
factorial: { n Int
    n > 0 ? { n * factorial(n - 1)}
            { 1 }
}

1.to(19 ).each({ x Int
    print('Factorial of `x` = `factorial x`')
}
`Output:
Factorial of 1 = 1
Factorial of 2 = 2
Factorial of 3 = 6
Factorial of 4 = 24
Factorial of 5 = 120
Factorial of 6 = 720
Factorial of 7 = 5040
Factorial of 8 = 40320
Factorial of 9 = 362880
Factorial of 10 = 3628800
Factorial of 11 = 39916800
Factorial of 12 = 479001600
Factorial of 13 = 6227020800
Factorial of 14 = 87178291200
Factorial of 15 = 1307674368000
Factorial of 16 = 20922789888000
Factorial of 17 = 355687428096000
Factorial of 18 = 6402373705728000
Factorial of 19 = 121645100408832000`

```

```js
arr: [1,2,3,4,5,6,7,8,9,10]
print(arr.filter { x Int; x % 2 == 0 })
```

Other formatting

```js
factorial: {n Int
  n > 0 ? {
    n * factorial(n - 10)
  } {
    1
  }
}

// 1 .to(19).each({ x Int
//   print('Factorial of `x` = `factorial x`')
// }
// alternate formatting
factorial: {n Int
   n > 0 ? { n * factorial n - 1} { 1 }
}
// using if
factorial: {n Int
    if n > 0 { n * factorial n - 1 } { 1 }
}
// using `when` (silly)
factorial: {n Int
  when [
     {n > 0}: {n * factorial n - 1}
  ] /*else*/ {1}
}

// using match (real)
factorial: {
	n Int
	match {
	   n > 0 => n * factorial(n - 1)
	}, {
	   1
	}
}

```
