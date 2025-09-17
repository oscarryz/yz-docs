#answered [Async + Lazy eval + Structured Concurrency](Async%20+%20Lazy%20eval%20+%20Structured%20Concurrency.md)


Update (sometime in 2023) Yz will be [Concurrency](../../../Features/Concurrency.md)

Single thread model with channels / process 
// status: I don't know but the important part is to avoid threads 
// as that implies sharing state
// Is better to share data to another thread/process/channel
// See Elixir and Go



@22 May 2022

Concurrency model to follow Go's  with channels (pipes)
See [go's example](https://github.com/adonovan/gopl.io/blob/master/ch8/pipeline3/main.go)

```go
NumberPipe <-> { // <-> says is't a pipe
    n Int // it handles Ints
}


/*
    Creates natural numbers ad-infinitum and writes them 
    in the `out` Pipe (or channel)
    by writing / setting the value in the Pipe `n` variable
*/
generate_numbers: { out Pipe
    n : 0
    while { true } {
        out.n = n 
        n = n + 1
    }
}
/*
    Square numbers takes two pipes
    `in` where to read numbers
    `out` where to write numbers
    
*/
square_numbers: { in Pipe; out Pipe
    n: in.n // blocks until there's something in `n` when run in a goroutine
    out.n = n * n // writes into the out
}
/*
    Print numbers coming from the `in` Pipe
*/
print_numbers: { in Pipe 
    print `Squared number: {in.n}` // blocks until `in.n` has a value
}

// Create two pipes
naturals: Pipe{}
squares: Pipe{}

// `go` keyword runs the block in an asynchronous way
go generate_numbers(naturals)
go square_numbers(naturals, squares)
go print_numbers(squares)


```

A possible approach without needing to use the `go` keyword would be to execute methods inside a pipe block

The same as before would be
```go

//
Squares <-> {
  natural Int
  square Int
  generate: {
    i Int
    while { continue } {
      natural = i
      i = i + 1
    }
  }
  square: {
    n: natural // blocks until natural has something
    // natural has been read, now the next line wil execute
    square = n * n
  }
  print: {
    print "{square}"
  }
}
sqrt: Squares{}
sqrt.generate() // Will run concurrently (as in go generate_numbers())
sqrt.square()   // Will run concurrently (as in go square_numbers())
sqrt.print()    // Will run concurrently (as in go print()())


```

A similar example with ping pong
```js
PingPong <-> {
  data String
  out  String
  ping: {
    while {true} {
      data = 'ping' // writes to data and blocks until data is read
    }
  }
  pong: {
    m: data // blocks until data has something
    out = m.reverse()
  }
  print: {
    print '{out}' // blocks until out has something
  }
}

pp: PingPong{}
pp.ping() // go pp.ping
pp.pong()
pp.print()

print '{pp.out}' // will block until pp.out has a value

```


// 01/15/2023

```javascript
ping <-> {
    s String
}
x: ping 'Pong'
do <-> {
    p <->{String}
    while{true}{p 'ping'}
}
do ping
5.times {_ Int
    print ping()
}


// with `go` keyword
ping:{
    while {true} {
        print 'Hello'
    }
}
go ping() 
print 'Continue'

do: {
    p <=> {String}
    while{true}{ p 'ping' }
}
p: <=> {s String}
go do p 
5.times{
    print p()
}



```

## Actors 
I just learned Elixir actually uses Actors

```javascript
'Prints Hello in the background'
hello <-> {
    time.sleep('1s')
    print 'Hello'
    hello()
}
print 'Starting'
hello() // starts printing 'Hello' in a separate process/thread
print 'Continue'
/*
Starting
Continue
Hello
Hello
...
*/
```

Counter, Squerer, Printer

1. Create numbers ad infinitum
2. Square them
3. Print them 

```javascript

Counter: {
    n Int
    x:0
    next : {
       n = n + 1 
       more ? {
           sleep 100ms
            x = x + 1
           next()
       }
    }
    more: {
        x < 100
    }
}
Squarer: {
    nn Int
    square: {
        n Int
        nn = n * n
        sleep 100ms
        counter.more()?{
            square()
        }
    }
}
Printer: {
    print: {n Int; print '{n}'}
}

```


If we use Actors, the code blocks wouldn't be too different from 
regular blocks. To send message the block would need to be constructed with the actors they interact. 

```javascript
// Producer -> Consumer/Producer -> Consumer
NumberGenerator: {
    s Squarer
    n: 1
    next: {
        while {true} {
            s.square s
            n = n + 1
            sleep 100.ms()
        }
    }
}
Squarer: {
    p Printer
    square: {n Int
        p.print n * n 
    }
}
Printer: {
    print: {n Int
        print '{n}'
    }
}
p <-> Printer{}
s <-> Squarer{p}
g <-> NumberGenerator{s}
p.next()
print 'Continues'
```

29 January 2023

Structured Concurrency. Have a scope marked by the `<->` simbol and any function call in there would be async but will wait (or abort) if with the rest

Example
```javascript
print 'Starting'
// Box is a wrapper to be able to have a reference to an object
Box: <T>{
    data T
    empty: true
    set: {t T
        data = t
        empty = false
    }
    pop: {
        empty = true
        data
    }
    
}
example <-> {
    channel Box<Int> = Box{0} // not really a channel but will be shared 
    produce(channel)
    consume(channel)
}
print 'End'
produce: { i Box<Int>
    while {true} {
     i.set random.int(100)
     time.sleep(100.ms())
    }
}
consume: { i Box<Int>
    while {i.data != int.nan} {
        n i.pop()
        print '{n}'
        time.sleep(50.ms())
        if n = 41 { 
            break// or cancel 
        }
    }
}
```
Problem:
- How to coordinate read/write on the Box object.
    - Maybe becuase it's in a coroutine they would be auto synced like channels do? It woudl be confusing though to beheve differently.
- What if I want to run sync code to prepare things?
    - Run it in a separate block?
        ```javascript
        nursery <-> {
            a b c = setup()
            do_with(a)
            do_with(b)
            do_with(c)
        }
        ```
        Would run all of them at the same time, create a block?
            ```
            nursery <-> {
               a b c = Int Int Int
               {
                   a b c = setup()
               }()
               do_with(a)
               do_with(b)
               do_with(c)
            }```
    Other obvious answer: Do it outside the nursery
        ```
        a b c = setup()
        nursery <-> {
               do_with(a)
               do_with(b)
               do_with(c)
        }```

### 30 January 2023
Nurseries and Channels: 

```javascript
// Nursery
concurrency <-> {
    out Channel<Int> = Channel{0}
    in  Channel<Int> = Channel{0}
    generate out
    square in out
    print in
}

generate: { 
    n Channel<Int>
    while {true} {
        n -> random.int 100
    }
}
square: {
    in Channel<Int>
    out Channel<Int>
    while {true} {
        //x: in.receive()
        x: in <-()
        out -> in <-()
    }
}
print: {
    in Channel<Int>
    print '{in <-()}'
}
```

```javascript

main: {
    generate: {
        random.int(1 100)
    }
    square: { n Int 
        n * n
    }
    print: {
        n Int
        print 'The number is {n}'
    }
}
```


[Structured Concurrency in Java](https://openjdk.org/jeps/428#:~:text=For%20example%2C%20in%20this%20single%2Dthreaded%20version%20of%20handle()%20the%20task%2Dsubtask%20relationship%20is%20apparent%20from%20the%20syntactic%20structure%3A)

```java
Response handle() throws IOException {
    String theUser  = findUser();
    int    theOrder = fetchOrder();
    return new Response(theUser, theOrder);
}
```
```javascript
//Yz sequential
handle: {
    the_user: find_user()
    the_order: find_order()
    Response { the_user the_order}
}
// Yz concurrent
handle <-> {
    the_user: find_user()
    the_order: find_order()
    if some_logic() {
        other_thing() // is other t
    }
    Response { the_user the_order }
}

```
Todo: structural concurrency

Cancellation
ErrorHandling
Fan in
Fan out?
Other logic inside nursery
Channels (most likely with generics Channel<Int>, but they are special)



