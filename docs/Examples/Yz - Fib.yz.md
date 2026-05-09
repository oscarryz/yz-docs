#example
Fibonacci

```js
//
fibonacci: { n Int
  n <= 2 ? { 1 }, {
    fibonacci(n - 1) + fibonacci(n - 2)
  }
}
fibonacci(10)

//
fibonacci: { n Int; current Int; next Int; result Int
  n == 0 ? { current }, {
    fibonacci(n - 1, next, current + next)
  }
}

fibonacci(10, 0, 1)

//
fibonacci: { n Int
  first, second: 0, 1
  1.to(n).each({ _ Int
    tmp: second
    second = second + first
    first = tmp
  })
  first
}

//
main: {
  name: get_line()
  print("Hello ${name}")
}


main: {
  file1, file2: get_args()
  str: read_file(file1)
  write_file(file2, str)
}
```


Yz v1.0
```js
fib: { n Int
    f: 0
    s: 1
    0.to(n - 1).each({ _ Int
        tmp: s
        s = s + f
        f = tmp
    })
    f
}

```
