#example


https://web.archive.org/web/20230329041616/https://www.golangprograms.com/go-language/concurrency.html

```yz
Philosopher : {
  name String
  left Fork
  right Fork
  bites : 0
  // To eat, a philosopher needs to gain
  // access to the `left` and `right` forks
  // -- also to `name` and `bites`, but those are not shared)
  eat : {
    left.pick_up()
    right.pick_up()
    println("${name} is eating...")
    bites = bites + 1
    left.put_down()
    right.put_down()
  }
  // done eating needs to gain access to `bites`
  done_eating? {
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
  philosophers: ["a", "b", "c", "d", "e"].map({ name String; Philosopher(name, Fork())})
  // their right fork is the left of the next philosopher
  philosophers.for_each({
    idx Int
    ri : idx == philosophers.len() -1  ? {0} : {idx + 1}
    philosophers[idx].right = philosophers[ri].left
  })

  // Run
  //
  // Each philosopher has a turn to each
  // if done eating, gets removed from the list
  while({philosopher.len() > 0 }, {
    philosopher.each({
      p Philosopher
      p.eat()
      p.done_eating ? {philosophers.remove(p)}
    })
  })
}
```