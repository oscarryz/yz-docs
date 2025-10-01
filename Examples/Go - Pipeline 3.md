#example

https://github.com/adonovan/gopl.io/blob/master/ch8/pipeline3/main.go

In a loop:

- Produce numbers
- Square them
- Print them
-
The producer knows when to stop

```js

// WORKS
// To make it work as close as possible, the squarer needs to invoke the counter // until it closes. The counter needs to know when to stop
`
Creates a block that returns an Int and a Bool every times it is invoked.
It will continue to be invoked while n < 100
This counter is the "producer"
`
counter: {
	n: -1
	{
	    n < 100 ? {n = n + 1 }
	    n < 100
	    n
	}

}
squarer: {
	producer #(Int, Boolean)  //e.g. the counter
	consumer #(Int)

	value Int
	open  Bool
	is_open: {
	    value, open = producer()
	}
	// while calls a block `is_open`.. similar to
	// while { is_open() } but better because it
	// doesn't need to create an extra block
	while is_open {
		consumer (value * value)
	}
}
printer: {
	n Int
	print('`n`')
}
main: {
   squarer( counter() , printer )
}
```


```js

// DOESNT_WORK
// The counter takes a "channel" which is a block that receives an Int
// it produces a number and sends it to the channel
counter (out #(Int)) = {
	1.to(100).each({
		e Int
		out(e) // send e to the channel
	}
}
// The squarer takes a "channel" (a block that takes an int)
// where to send the data and returns a block that expects an int;
// that's how the producer and the printer
// gets connected.
squarer: {
	out #(Int)
	#(v Int) {
		out( v * v )
	}
}
// The printer is a block that takes an Int and prints it
printer #(n Int)  {
    print(`n`)
}
main: {
	producer(squarer(printer))
}

//
// Another approach that ressambles the actual Go code but doesn't work
//
counter: {
	out #(Int)
	1.to(100).each({ e Int ; out(e) }
}
squarer: {
	in #(Int)
	out #(Int)
	v: in()
	out(v * v)
}
printer: {
  in #(Int)
  v: in()
  print(`v`)
}
main: {
    naturals : {n Int; n}
    squares  : {n Int; n}
    counter(naturals) // ok
    squarer(naturals squares) // no function found
    printer(squares)
}

// The reson it doesn't work is because the block `(Int)` by itself cannot stop for the value to be overriden before is read (channels do), it can block until the value is avaiable which is righ away.
// It would work if a buffer is used to store the values so it is not overriden


// Another reason it doesn't work is because it tries to use a function as a buffer, but returning the last value wouldn't sync (f(1), f(2), f(), f() wouldnt necesarily 1 and 2 i the same order)

```
