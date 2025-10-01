#example

[A boring function](https://go.dev/talks/2012/concurrency.slide#12)
```js
boring: {
	i:0
	msg String
	loop: {
		print('`msg`, `i`')
		time.sleep(time.second)
		// less boring
		time.sleep(time.duration(random.int(1) * time.millisecond))
		i = i + 1
	}
}
main: {
	boring("boring!") // launches and continues
	// but then waits at the end of the main block
}
```


[Channels](https://go.dev/talks/2012/concurrency.slide#19)
// no channels in Yz, just send the message
```js
not_a_channel: {
	c Int
}
not_a_channel(1) // sending
value: not_a_channel.c // "receiving"
```

```js
not_a_channel
main:{
	nac (String) = {
		s String
	}
	boring("boring" nac)
	5.times({
		wait_for: nac.value_set()
		print('You say: `nac.s`' // right now this would just print nac.s 5 times)
	}
	print("You're boring; I'm leaving")
}
boring: {
	msg String
	nac (String)
	loop: {
		i Int
		nac('`msg` `i`') // callback, and continues
		time.sleep(time.duration(rand.int(1)* time.millisecond))
	}
}
```

Ideas:
#open-question  How can we wait for a value to be set? [../../Questions/solved/Wait for a value to be set](../../Questions/solved/Wait%20for%20a%20value%20to%20be%20set.md)


- We could potentially block until a variable has a value, just like `c <- chan`  in Go blocks until the channel has a value.
- We could check for a special value and clear it after use

```js
loop: {
   s: nac.value // will wait until nac.value has something
   clear(nac.value)
   print(s)
}
```
