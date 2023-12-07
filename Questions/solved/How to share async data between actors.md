[Wait for a value to be set](solved/Wait%20for%20a%20value%20to%20be%20set.md)

We have two blocks, one produces a value, the other consumes it. How to make this async while keeping the structured concurrency nature? (see [Concurrent by Default](../../Features/Concurrent%20by%20Default.md) )



// The example below doesn't show the problem I'm trying to establish. 
// but it's kind the idea 
```javascript
producer: { 
	s: 'Random number $(random.int(100))' 
}
consume: {
   thing {String} // a thing that retuns a string
   value: thing() // invoke and wait for a value
   print 'The value is: $(value)'
}
main: {
	consume(producer)
}

```


#answered  Just use the return value to sync and rely on syncing at the bottom, if something has to react to something else, do use a callback

As an example see: [Go - Pipeline 3](Go%20-%20Pipeline%203.md)
