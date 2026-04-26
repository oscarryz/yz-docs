#example

#concurrency 

[A boring function](https://go.dev/talks/2012/concurrency.slide#12)
```js
boring: {
  msg String
  i: 0
  while({true}: {
    print("`msg`, `i`")
    time.delay(random.int(3))
    i = i + 1
  }
}
main: {
  boring("boring!") // launches
  // and waits at the end of the main block because of structural concurrency
}
```


[Channels](https://go.dev/talks/2012/concurrency.slide#19)

```js
boring: {
    m String
    next: { value String }
    i: 1
    
    while({ true }, {
        time.delay(1)
        next("`m` `i`")    // acquires {next}, writes, releases
        i = i + 1
    })                      // ← natural yield point here!
}

main: {
    boring("sync")          // launches boring, doesn't block
    5.times().do({
        print(boring.next()) // acquires {next}, reads, releases
    })
}
```

Because the BoC model, the boc in the `while` and the boc in the `do` loop have to gain access to the `next` boc, and due to the _happens-before_ trait, they interleave taking turns to write and read to it

The queue on next is the magic:

```
boring: acquires {next}, writes "sync 1", releases
main:   acquires {next}, reads  "sync 1", releases
boring: acquires {next}, writes "sync 2", releases
main:   acquires {next}, reads  "sync 2", releases
```

The while loop's recursive nature creates natural yield points because each iteration has to re-acquire next — and if main is waiting on it, boring yields implicitly without any new keyword.