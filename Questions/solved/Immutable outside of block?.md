#idea 

Make local variables mutable but outside the block/object they are immutable, they can only be modified through "messages"

```javascript
Counter: {
    count: 0
    increment: {
        count = count + 1
    }
    reset: { count = 0 }
}
c1: Counter{}
stuff: {
    a Counter
    if a.count > 10 {
        do_something()
        // compilation error: cannot modify a.count
        // a.count = 0 
        a.reset()
        a = Counter{} // maybe reassing it just for an example 
    } {
        a.increment()
    }
}
11.times { 
    stuff c1
    c1 = Counter{c.count} // just because 
}
```

#rejected If we achieve "1 writer" with "everything is an actor" this is not needed. 