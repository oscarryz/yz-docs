#example

#concurrency 

[A boring function](https://go.dev/talks/2012/concurrency.slide#12)
```js
boring: {
  msg String
  i: 0
  loop: {
    print('`msg`, `i`')
    time.delay(1)
    // less boring
    time.delay(random.int(3))
    i = i + 1
    loop()
  }
  loop()
}
main: {
  boring("boring!") // launches and continues
  // but then waits at the end of the main block because of structural concurrency
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
  nc : { m String }
  boring("b", c )
  5.times().do({
    print("You said: `nc.m`)
  })
}
```

Because the BoC model, the boc in the `while` and the boc in the `do` loop have to gain access to the `nc` boc, due to the _happens-before_ trait, they interleave taking turns to write and read the `m` variable.
