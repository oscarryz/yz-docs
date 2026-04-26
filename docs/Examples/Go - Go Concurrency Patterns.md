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
    messages [String]()
    next: { 
       messages.pop()
    }
    i: 1
    
    while({ true }, {
        time.delay(1)
        messages.push("`m` `i`")    // acquires {messages}, writes, releases
        i = i + 1
    })                      // ← natural yield point here!
}

main: {
    boring("sync")          // launches boring, doesn't block
    5.times().do({
        print(boring.next()) // acquires {messages}, reads, releases
    })
}
```

Because the BoC model, the boc in the `while` and the boc in the `do` loop have to gain access to the `messages` boc, and due to the _happens-before_ trait, the first message is written before the first read takes place.
