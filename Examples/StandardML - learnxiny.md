#example

https://learnxinyminutes.com/docs/standard-ml/

```js
first_element: {
  list [T]
  list[0]
}
second_element: {
  list [T]
  list[1]
}

// http://go/learnxinyminutes
/*
fun evenly_positioned_elems (odd::even::xs) = even::evenly_positioned_elems xs
  | evenly_positioned_elems [odd] = []  (* Base case: throw away *)
  | evenly_positioned_elems []    = []  (* Base case *)
*/
evenly_positioned_elements: {
  list [T]
  (list.length() == 0 || list.length() == 1) ? {
    [T]()
  }, {
    list[1] ++ evenly_positioned_elements(list.sublist(1))
  }
}

```

Fibonacci

```js

// With nested conditionals
fibonacci: {
  n Int
  n == 0 ? {
    0
  }, {
    n == 1 ? {
      1
    }, {
      fibonacci(n - 1) + fibonacci(n - 2)
    }
  }
}
// Same with different format
fibonacci: {
    n Int
    n == 0 ?
    { 0 },
    { n == 1 ? { 1 },
    { fibonacci(n - 1) + fibonacci(n - 2) } }
}

// With `match`
fibonacci: {
    n Int
    match {
     n == 0 => 0
    }, {
     n == 1 => 1
    }, {
     fibonacci(n - 1) + fibonacci(n - 2)
    }
}
```