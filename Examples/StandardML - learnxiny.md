```javascript
first_element: <T> {
  list []T
  list[0]
}
second_element: <T> {
  list []T
  list[1]
}

// http://go/learnxinyminutes
/*
fun evenly_positioned_elems (odd::even::xs) = even::evenly_positioned_elems xs
  | evenly_positioned_elems [odd] = []  (* Base case: throw away *)
  | evenly_positioned_elems []    = []  (* Base case *)
*/
evenly_positioned_element: <T> {
  list []T

  l: len(list)
  l == 0 || { l == 1 } ? {
    []
  } {
    list[1] ++ sublist(list,1)
  }
}

```

Fibonacci

```javascript

fibonacci: { 
    n Int
    n == 0 ? 
    {0} 
    { n == 1 ? 
         {1} 
         { fibonacci n - 1 + fibonacci n - 2 }
    }
}

fibonacci: {
    n Int
    when_eq n [
     {0}: {0}
     {1}: {1}
     {true}: {fibonacci n - 1 + fibonacci n-2}
    ]
}
```

