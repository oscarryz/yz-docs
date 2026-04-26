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

// no channels in Yz
// just use a boc to
// send and read the message
boring: {

  msg String
  nc #(m String)

  i: 1

  while {true}, {
    time.delay(1)
    nc("`msg` `i` back")
    i = i + 1
  }
}

main:{
  // Can't launch and forget
  // if we want to hear back
  // nc : { m String }
  // boring("b", nc )

  // It needs to be used as callback
  5.times().do({
    boring("fn", { m String; print("You said: `m`") })
  })
}
```

Because the BoC model, the boc in the `while` and the boc in the `do` loop have to gain access to the `nc` boc, due to the _happens-before_ trait, they interleave taking turns to write and read the `m` variable.
