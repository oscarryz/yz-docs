#example

Fibonacci sequence using range iteration:

```yz
main: {
  a: 0
  b: 1
  c: 1
  print("`a + b + c`")
  0.to(50).each({ _ Int
    a = b
    b = c
    c = a + b
    print("`c`")
  })
}
```
