What if the block type could use `::` and `;`

eg
```javascript
// current
walkie_talkie: {
  thing {talk {String} walk {String}} 
  print thing.walk()
  print thing.talk()
}
// this suggestion
walkie_talkie: {
  thing :: talk :: String; walk :: String
  print thing.walk()
  print thing.talk()
}
```

More
```js

a_two {a Int Int } = {
  a: 1
  2
}

//swap is a block with two variables of type int

swap {a Int b Int} = { a Int b Int; b a }
one two: swap 2 1 

// this proposal 

a_two :: a Int ; Int  = {
  a: 1
  2
}
// swap is a block with two variables of type int

swap :: a Int; b Int = { b a }
one two: swap 2 1 
```

```js
thing { {String} {String} }
thing ::  :: String  :: String
thing #( #(String), #(String))
```


```js
person {name String} = {
    name: 'Bob'
    age Int // not accessible from outside as the block only declared `name String`
}
//this

person :: name String = {
    name: 'Bob'
    age Int // 
}
```
```js
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
// this

// The producer takes a "channel" which is a block that receives an Int
// it produces a number and sends it to the channel 
producer :: channel :: Int = {
	1.to(100).each { e Int
		channel(e) // send e to the channel
	}
} 
// The squarer takes a "channel" (a block that takes an int)
// where to send the data and returns 
// a block that expects an int, that's how the producer and the printer 
// gets connected.
squarer: {
	channel :: Int
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
}```

```js
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
// this
counter: {
	out :: Int
	1.to(100).each{ e Int ; out(e) }
}
squarer: {
	in :: Int
	out :: Int
	v: in()
	out(v * v)
}
printer: {
  in :: Int
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
```


#answered  Use `#()` e.g. `max #(Int,Int,Int)` is a block that takes and/or returns 3 ints
