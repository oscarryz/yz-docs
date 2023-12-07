https://github.com/adonovan/gopl.io/blob/master/ch8/pipeline3/main.go

In a loop: 

- Produce numbers
- Square them
- Print them
-
The producer knows when to stop 


```javascript

// The producer takes a "channel" which is a block that receives an Int
// it produces a number and sends it to the channel 
producer {channel {Int}} = {
	channel {Int}
	1.to(100).each { e Int
		channel(e) // send e to the channel
	}
} 
// The squarer takes a "channel" (a block that takes an in)
// where to send the data and returns 
// a block that expects an int, that's how the producer and the printer 
// gets connected.
squarer: {
	channel {Int}
	{   n Int 
		channel ( n * n )
	 }
}
// The printer is a block
printer: { n Int 
    print "$(n)"
}
main: {
	producer(squarer(printer)) 
}

// another approach that ressambles the actual Go code but doesn't work

counter: {
	out {Int}
	1.to(100).each{ e Int ; out(e) }
}
squarer: {
	in {Int}
	out {Int}
	v: in()
	out(v * v)
}
printer: {
  in {Int}
  v: in()
  print '$(v)'
}
main: {
    naturals : {n Int; n}
    squares  : {n Int; n}
    counter(naturals)
    squarer(naturals squares)
    printer(squares)
}

// The reson it doesn't work is because the block `{Int}` by itself cannot stop for the value to be overriden before is read (channels do), it can block until the value is avaiable which is righ away.

// To make it work as close as possible, the squarer needs to invoke the counter until it closes. The counter needs to know when to stop


`
Creates a block that returns an Int and a Bool every times it's invoked.
It will continue to be invoked while n < 100
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
	producer {Int Boolean} 
	consumer {Int}

	value Int
	open  Bool
	is_open: { 
	    value open = producer()
	}
	while is_open { 
		consumer (v * v)
	}
}
printer: {
	n Int 
	print '$(n)'
}
main: {
   squarer( counter() printer )
}
```