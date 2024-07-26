[Dining Philosophers](https://en.wikipedia.org/wiki/Dining_philosophers_problem)

The following is a basic setup. Will still need some time to do a test run and see what happens with the resources.


```javascript
// In lieu of imports, declare three typs to be the same as those returned 
// by executing the anonymous block with those three types full qualified name
Option Some None : ({
    std.option.Option
    std.option.Some
    std.option.None
}())

/*
If that's not possible the following should be

Option : std.option.Option
Some : std.option.Some
None : std.option.None
*/

// A fork has two operations: 
// `try_take` by a Philosopher and returns true or false if it was possible
// `try_drop` if the fork is currenly holded by the given Philosopher
Fork (
    try_take (by Philosopher, Bool),
    try_drop (by Phisolopher)
)  = {

    current_user Option(Philosopher)
    taken   : false
    is_free : { taken == false } 

    try_take : {
        by Philosopher
        is_free() ? { 
            taken = true
            current_user = Some(by)
            true 
        } { 
            false
        }
    }
    try_drop : {
        by Philosopher
        current_user.value_is(by) ? {
            taken = false
            current_user = None()
        }
    }
}

Philosopher : { 
    name String
    left  Fork
    right Fork
    self Philosopher

    think () = {
           print("$(name) is eating...")
           time.sleep(random(1 5), time.SECONDS)
           eat()
    }
    eat   () = {
        left.try_take(self) && { right.try_take(self) } {
           print("$(name) is eating...")
           wait: time.sleep(random(1 5), time.SECONDS)
        }
        left.try_drop(self)
        right.try_drop(self)
        think()
    }
}

main: { 
    philosophers = init(["Plato" "Socrates" "Kant"])
    while { true } { 
        philosophers.each { 
            p Philosopher
            eat()
        }
    }

}
init: {
    names [String]
    philosophers : []Philosophers
    w: names.for_each { 
        i Int
        name String

        p : Philosopher(name)
        p.self = p
        p.right = Fork()
        is_last :  i == names.len() - 1
        next : is_last ? { 0 } { i + 1 } 
        philosophers[next].left = p.right 
    }
    philosophers
}

