#answered No channels 
#rejected 
# Actors and Channels

Every instantiable code is an actor (and/or can be used as actor)
Every block can be used as channel using the `<-` operator

```javascript
Producer: {
   produce: {
       channel { Int } // a channel of ints
       x: 0
       while { true } {
           channel <- x
           // or just channel(x)
           x = x + 1
        }
    }
}
Consumer: {
   consume: { 
       channel {Int}
       while {true} {
           n: <- channel // receives a number from channel
           // or just n: channel()
          // print '{n*n}'
       }
   }
}
main: {
    go : yz.concurrency.spaw
    // Option A
    work: go { 
       p : Producer{}
       c : Consumer{}
       x : { n Int } // the channel
       p.produce x  // because it's invoked in a spawn block, this method is async
       c.consume x  // so is this
       print 'Hola' // what about this?
    }
    time.wait 1,000.ms()
    work.stop()
    // Option B
    w2: go {
        actors: yz.concurrency.actors
        p: actors Producer{} 
        c: actors Consumer{}
        x: {n Int}
        p.produce x
        c.consume x
    }
    // Option C
    Employee {
        name String
        _id Int
        id: { 
            n Int
            _id = n
            _id 
        }

    }
    e Employee = Employee { 'Someone' }
    concurrent: yz.concurrency
    work: spawn {

        p: concurrent(Producer)
        p: concurrent(Consumer)
        p.produce(concurrent(e.id)
        c.consume(e.id)
    }
}
```

# Using a library
`spawn` is a function that can execute asynchronous methods on objects
obtained through `context.get`. These objects act as actors and each 
method call would behave async. 
```javascript
spawn: yz.concurrent.spawn
context: yz.concurrent.context
main: {
    e: Person{}
    w: spawn {
        p: context.get<Producer>()
        c: context.get<Consumer>()
        x: context.channel(e.id)
        p.produce(x)
        c.consume(x)
    }
    time.wait 1000.ms()
    w.cancel()
}
```

### Download image and metadata

```javascript
S3: aws.something.S3
Redis: my.imagination.library.Redis
spawn: yz.concurrent.spawn
context: yz.concurrent.context
display: yz.ui.display

main: { 
    s: {
        s3: context.get<S3>()
        redis: context.get<Redis>()
        image_blob: s3.fetch('s3://something/something')
        metadata: redis.query('abcd')
        print ('Fetching stuff')
        result: Image{image_blob medata}
    }
    load_image: spawn s
    load_image.wait()
    display(load_image.result)

}
```
### With regular actors

```javascript
Producer: {
    produce: { c Consumer 
        while {true} {
            n = n + 1
            c.consume(n)
        }
    }
}
Consumer: {
    consume: {n Int
        print '{n}'
    }
}
main: spawn {
    p = Producer{}
    c = Consumer{}
    p.produce(c)
}
```

_March 12 2023_

Elements created in a `spawn` block are "actors" and "channels" 


```javascript 
external: {
    message String
    print message
}
work: spawn => {
    f: Foo{} // Foo methods will be async. If return value is expected
    // the spawn will wait 
    f.something()
    external 'Will print right away'
    x: f.something_else() // 
    external 'Will print right away too'
    // before the block finishes
    // `x` will be "awaited" 
}
// inside foo
Foo :{ 
    something: {
        // long computation
    }
    something_else: {
        // another long conputation,this one
        // will run after the first one sequentially 
        42
    }
}
```