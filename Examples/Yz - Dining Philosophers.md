[Dining Philosophers](https://en.wikipedia.org/wiki/Dining_philosophers_problem)

The following is a basic setup. Will still need some time to do a test run and see what happens with the resources.

```javascript
Fork : {
    taken   : false
    is_free : { taken == false } 
    take    : { taken = true   } 
    drop    : { taken = false  } 
} 
Philosopher : { 
    name String
    left  Fork
    right Fork

    think () = {
           print("$(name) is eating...")
           time.sleep(random(1 5), time.SECONDS)
           eat()
    }
    eat   () = {
        left.is_free() && { right.is_free() } { 
           left.take()
           take.take()
           print("$(name) is eating...")
           wait: time.sleep(random(1 5), time.SECONDS)
           left.drop()
           right.drop()
        } {
            think()
        }
    }
}

dine: { 
    philosophers : [ 
       Philosopher("Plato")
       Philosopher("Socrates")
       Philosopher("Kant")
    }
    philosophers = setup(philosophers)
    while { true } { 
        philosophers.each { 
            p Philosopher
            eat()
        }
    }

}
setup: {
    philosophers [] Philosopher
    w: philosophers.for_each { 
        i Int
        p Philosopher
        p.right = Fork()
        
    
        next : i == philosophers.len() - 1  {
            0 
        } { 
            i + 1
        } 
        philosophers[next].left = p.right 
    }
    philosophers
}

