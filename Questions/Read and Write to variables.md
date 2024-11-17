
Idea, make variable read only by default and require objects to provide a writer method.

Variables can be written only by objects that have direct scope visibility (loca and parent)

```js
parent: {
  x Int // 
  x = x + 1 // ok, it has local scope
  x_= : {
    y Int 
    x = x + y // writting to parent.x because it has visibility
  }
  x_== : { n Int ; x  == n }
}

other: {
  parent(1) // write to parent.x through its "execution" form `()`
  parent.x // can read it
  // parent.x = 2 // compilation error, can't write it, can see it through `parent` but is not in the "inner" scope
  parent.x_=(3)
  assert (parent.x == 4)
}
```

This ensures single writer multiple reader concurrency 