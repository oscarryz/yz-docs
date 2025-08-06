(Invalid, see below) How concurrency is going to work is still TBD
```javascript
spawn: yz.concurrent.spawn
context: yz.concurrent.context

main: spawn {
    c: context.channel({n:0})
    a: spawn { 
        x: 10
        // sending 10 to channel 
        c(x)
    }
    b: spawn {
        // Receiving 10 from channel
        x: c()
    }
}

```

With the latest content from:  [[../Features/Concurrency]]
```javascript
// Attempt to run the above
main: {
    c: {n Int}
    a: {
       x: 10
       c(x) // calls `c` with arg value 10
    }
    b: {
       x: c.n // receives the last arg from c
    }
    // Problem, b has to run until a has set the value.
    // Solution: call b from c
    a()
    b()
}
```
```javascript
//Corrected
main: {
    b {Int}
    c: {
        n Int
        b(n)
    }
    a: {
        x:10
        c(x)
    }
    b = {
       x Int
    }
    a()
    // b() no need to call `b` will be called from `c`
    // slightly problem, this whole thing is synchronized
    // we might want to run it async, but that's another example
}
```

```javascript
// Behaving like channels
main: {
    b #(Int)
    c: #(n Int) // c: {n Int = ints.NAN}
    a: {
        // do stuff if some cond is true, send x
        loop {
          x: int.random()
          if x == 10 {
              c(x) //
          }
        }
    }
    b: {
        // do stuff... when c has a value act on it.
        loop {
            print "I'm doing stuff"
            if c.n != int.nan {
               print 'Hey I can see c has something' // at the cost of always be checking...
            }
        }
    }
    // who calls b though? And how to make b wait 
}
```



[Alternate wait in ballerina]()
```ballerina
// ballerina
function fetch(string url) returns string|error { ….
}
function altFetch(string urlA, string urlB) returns string|error {
    worker A returns string|error { 
        return fetch(urlA); 
    } 
    worker B returns string|error { 
       return fetch(urlB);
	} 
	return wait A | B; 
}
```

Alternate wait in Yz
```js
fetch: {
   string Url
   // ... return string or error 	   
}
alt_fetch: {
   url_a String
   url_b String
   a: {
      fetch(url_a)
      return
   }
   b: {
      fetch(url_b)
      return
   }
   // Launches both, but when one of them finishes it will exit the enclosing `alt_fetch`
   a()
   b()
}
```
(see: [return, break, continue](../../Features/return,%20break,%20continue.md))

Inter worker message passing: 

```js

main: {
   b: {
	   x1 Int
   }
   c: {
	   x2 Int
   }
   a: {
	   b(1)
	   c(2)
   }
   wait_for_a: a() // waits until a finishes 
   y1: b.x1 
   y2: c.x2
   z:  y1 + y2
}
```