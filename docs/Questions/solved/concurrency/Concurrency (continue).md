#answered  [Async + Lazy eval + Structured Concurrency](Async%20+%20Lazy%20eval%20+%20Structured%20Concurrency.md)

Feb 27 2023

Similar to Structural Concurrency + Actors  

- Launch actors in a launch block
- Only the actors created in the block will work concurrently (actors are just instantiable blocks)
- Share actors among themselves for cooperation
- Only methods (or blocks) will be used to communicate and not primitives

```javascript
CreateNumbers: {
    n Int
    create: {
        while { true } {
            set_n( x + 1 )
        }
    }
    set_n: { a Int 
        n = a
    }
    get_n: { n }
}
PrintNumber: {
    print: { cn CreateNumber
        while {true} {
            n = cn.get_n()
            print '{n}'
        }
    }
}
main: {
    f: Foo{}
    x: concurrency.launch { 
        cn: CreateNumbers{}
        pn: PrintNubers{}
        f.bar() // Foo was created outside, no couroutine created
        cn.create() // starts in a new coroutine
        pn.print(cn) // separate couroutine
        print 'Hello' // also created outside, no couroutine created
    }
    time.sleep 100.ms()
    x.cancel() // stops all the coroutines
}

```