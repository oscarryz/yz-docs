#example

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
// no channels in Yz, just send the message
```js
not_a_channel: {
	c Int
}
not_a_channel(1) // sending
value: not_a_channel.c // "receiving", value has type `Int` but is a thunk, it will materialize when used.
```

```js
main: {
  nac #(s String) {}    
  boring("boring", nac)
  0.to(5).each({ _ Int
    wait_for: nac.value_set()
    print('You say: `nac.s`') // right now this would just print nac.s 5 times
  })
  print("You're boring; I'm leaving")
}
boring: {
	msg String
	nac #(String)
	loop: {
		i Int
		nac('`msg` `i`') // callback, and continues
		time.sleep(time.duration(rand.int(1) * time.millisecond))
	}
}
```

Ideas:
#open-question  How can we wait for a value to be set? [../../Questions/solved/Wait for a value to be set](../../Questions/solved/Wait%20for%20a%20value%20to%20be%20set.md)

```js
loop: {
   s: nac.value // will wait until nac.value has something
   clear(nac.value)
   print(s)
}
```
