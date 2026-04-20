#answered 
#rejected Follow the original design and let everything declared as in the signature to be public, but make the owner the solely writer for that variable [Single Writer](../../Features/Single%20Writer.md)

Idea, make variable read only by default and require objects to provide a writer method.

Variables can be written only by objects that have direct scope visibility (local and parent)

```js
parent: {
  x Int // 
  x = x + 1 // ok, it has local scope
  x_= : {
    y Int 
    x = x + y // writting to parent.x because it has visibility
  }
}

other: {
  parent(1) // write to parent.x through its "execution" form `()`
  parent.x // can read it
  // parent.x = 2 // compilation error, can't write to it
  parent.x_=(3) // setter method
  std.assert (parent.x == 4)
}
```

This ensures single writer multiple reader concurrency 