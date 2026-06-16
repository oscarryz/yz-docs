#example
[Dining Philosophers](https://en.wikipedia.org/wiki/Dining_philosophers_problem)

The following is a basic setup. Will still need some time to do a test run and see what happens with the resources.


```js
// In lieu of imports, declare three typs to be the same as those returned
// by executing the anonymous block with those three types full qualified name
Option, Some, None: std.option.Option,  std.option.Some,  std.option.None

// A fork has two operations:
// `try_take` by a Philosopher and returns true or false if it was possible
// `try_drop` if the fork is currenly holded by the given Philosopher
Fork #(
    try_take #(by Philosopher, Bool),
    try_drop #(by Philosopher)
) = {

    current_user Option(Philosopher)
    taken:   false
    is_free: { taken == false }

    try_take: {
        by Philosopher
        is_free() ? {
            taken = true
            current_user = Some(by)
            true
        }, {
            false
        }
    }
    try_drop: {
        by Philosopher
        current_user.value_is(by) ? {
            taken = false
            current_user = None()
        }
    }
}

Philosopher: {
    name String
    left  Fork
    right Fork
    self Philosopher

    think #() {
           print("${name} is thinking...")
           time.sleep(random(1, 5), time.SECONDS)
           eat()
    }
    eat  #() {
        left.try_take(self) && { right.try_take(self) } ? {
           print("${name} is eating...")
           wait: time.sleep(random(1, 5), time.SECONDS)
        }
        left.try_drop(self)
        right.try_drop(self)
        think()
    }
}

main: {
    philosophers = init(["Plato", "Socrates", "Kant"])
    while({ true }, {
        philosophers.each({
            p Philosopher
            eat()
        })
    })

}
init: {
    names [String]
    philosophers: [Philosophers]()
    w: names.each({
        i Int
        name String

        p: Philosopher(name)
        p.self = p
        p.right = Fork()
        is_last: i == names.length() - 1
        next: is_last ? { 0 }, { i + 1 }
        philosophers[next].left = p.right
    })
    philosophers
}


```


A simplified version where they keep eating until they are done 


```js
// Using #(...) to hide the `bites` variable
Philosopher #( 
  name String,
  left Fork,
  right Fork,
  eat #(),
  is_done? #(Bool) 
) {
  bites : 0
  // To eat, a philosopher needs to gain
  // access to the `left` and `right` forks
  // -- also to `name` and `bites`, but those are not shared)
  eat = {
    left.pick_up()
    right.pick_up()
    println("`${name} is eating...")
    bites = bites + 1
    left.put_down()
    right.put_down()
  }
  // done eating needs to gain access to `bites`
  is_done? = {
    bites >= 3
  }
}
// Simple structure with two methods, pick_up and put_down
Fork : {
  pick_up : {}
  put_down : {}
}

main: {
  // Setup
  //
  // Each philosopher starts with a left fork.
  table: ["a", "b", "c", "d", "e"].map({ name String; Philosopher(name, Fork())})

 
  // each philosopher right fork, is the next one's left
  table.for_each({
    i Int
    current Philosopher = table[i]
    next Philosopher = table[ i == table.len() -1 ? { 0 } : { i  + 1} ] 
    current.right = next.left
  })

  // Run
  //
  // Each philosopher has a turn to each
  // if done eating, gets removed from the list
  while({ philosopher.is_empty() == false }, {
    philosopher.each({
      p Philosopher
      p.eat()
      p.is_done ? { philosophers.remove(p) }
    })
  })
}
/*
  lf a rf
  queue [(lf, b, rf)]
  lf c rf
  queue [(), [lf d rf]]
  lf e
*/
```


