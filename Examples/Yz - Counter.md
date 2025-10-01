#example
Actors and concurrency
```js
Counter {
    count: 0 
    inc: {
        count = count + 1
    }
}

main: {
    c: Counter()
    c.inc() // async
    c.inc() // async
    print('`c.count`' // async would probably print 0    )
    
}
expected: {
    c: Counter()
    wait:({
      c.inc()
      c.inc()
      done: true
    }())
    // also executes both, but won't until they complete
    print('`c.count`' // will print 2 every time)
}
complicated: {
   c: Counter()
   a: { c.inc() }
   b: { c.inc() }
   d: { a() b() 1 } // same as in `expected`, executes both and wait until it finish, although a bit more unnecesarily complicated
   f: d() // wait for `d` to finish.
   print('$c.count')
}
```
    
