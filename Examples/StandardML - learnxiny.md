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
  (list.len() == 0 || {list.len() == 1}) ? {
    []T
  }, {
    list[1] ++ evenly_positioned_elements(list.sublist(1))
  }
}

```

Fibonacci

```js

// With non-local returns (or just with `return`)
fibonacci: {n Int
  n == 0 ? {0}
  n == 1 ? {1}
  fibonacci(n - 1) + fibonacci( n - 2)
}
// Bool and nested bool
fibonacci: {n Int
  n == 0 ? {
    0
  }, {
    n == 1 ? {
      1
    }, {
      fibonacci(n - 1) + fibonacci( n - 2)
    }
  }
}
// Same with different format
fibonacci: {
    n Int
    n == 0 ?
    {0},
    { n == 1 ?  {1},
    { fibonacci(n - 1) + flibonacci(n - 2) } }
}

// With `when_eq`
fibonacci: {
    n Int
    when_eq n [
     {0}: {0}
     {1}: {1}
     {true}: {fibonacci(n - 1) + fibonacci(n - 2)}
    ]
}
// With ifs
fibonacci: {n Int
  if(n == 0 , {0},
    else_if(n == 1, {1},
      else({ fibonacci(n - 1) + flibonacci(n - 2) })
    )
  )
}
// Same ifs without parenthesis
fibonacci: {
  n Int
  if n == 0 , {0},
  else_if n == 1, {1},
  else { fibonacci(n - 1) + flibonacci(n - 2) }
}
// Where the ifs signatures are
if #(Bool, then #(V), else #(V))
else_if #(Bool, then #(V), else #(V))
else #(#(V))
```
