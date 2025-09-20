#feature


The whole language

## Variables

```js
// comment
/* 
multi line 
comment */
// variable declaration long format -> id type (=) expr
message String = "Hello"
// variable declaration short format ->  id : expr
to_whom: "World"
// String literal "", "", `` multiline and interpolated
print("`message`, `to_whom`!") // Hello, World!
n Int = 1
d : 3.14
```

## Blocks of code (boc)
```js
// put some variable sourroundded by `{` `}` and you have a block
{ 
	message String = "Hello"
	to_whom: "World"
  	print("`message`, `to_whom`")
}
// no output, because it"s not executed

// Blocks can be assigned to variabless
greet: { 
	message String = "Hello"
	to_who: "World"
	print("`message`, `to_who`")
}
// execute with `(` `)`
greet() // prints "Hello, World!"

// blocks variables can be accessed using `.` notation
greet.to_whom = "Everybody"
greet() // prints "Hello, Everybody!"
// even after execution
greet.message // returns "Hello"

// The invocation call can also use the variables as parameters
greet("Hola", "Mundo") // prints "Hola, Mundo!" 

// And the result of the invocation assigned to variables
what_was_printed: greet("Hallo", "Wereld") // prints "Hallo, Wereld!" and assings it to the variable `what_was_printed` because the last execution statement was the print which returns that value.

// Some values can be ommited if they have a default and they can also be named

greet(to_whom: "world", message: "Nice to meet you") // prints "Nice to meet you, world!"

greet(to_whom: "mundo") // prints "Nice to meet you, mundo!"

// So blocks are like objects + functions / closures / methods
```

### Block type / signature

```js
// Blocks have type too, this is specified between `#()` and the internal variables separated by `,`, if no variable is needed (e.g. returning a value) just add the type e.g #(s String, Bool).

// So a variable declaration long format -> id type (=) expr
greet #(message String, to_whom String, String) = {
	print("`message` `to_whom`") // no need to redleclare vars if in the same go
}

greet #(message String, to_whom String, String)
// something else 
print("blah")
greet = {
// vars need to be declared here to make sure the block is assignable to the type `greet`
	message String
	to_whom String
	message // just returns the variable message
}

// The type signature can also ommit the variable names
// Here `greet` is a varialble of type block that takes (or returns) three Strings
greet #(String, String, String) 
greet() // compilation error, needs to be assigned a value.
greet = {
	a String 
	b String
	c String
}
greet() // compialation error, a, b, c need a default value or one assigned
greet("uno", "dos", "tres" )


```

## Creating instances

```js
// A block can return another block
create_block: {
	{
		name String 
	}
}
x: create_block() 
x.name = "X"
x() // just returns `X`

```

## Creating new types

```js
// if name starting with uppercase then it declares a new type
Named: {
	name String
}
// when invoked it creates a copy 
x: Named("X")
x.name // x.name is "X"

// obviously it can have other blocks
Person: {
	name String
	introduce_yourself: {
		// and blocks can access the outer scope
		print("My name is `name`")
	}
}
a: Person("Alice") 
a.introduce_yourself() // My name is Alice

// acting like objects / clases with methods. 

// The difference between declaring a new type and 
// returning a block is, the type can be used when declaring variables...
// so

Person: {
	name String
}
// and 
create_named: {
	{
		name String
	}
}
//both create a copy of a block with a String
// x is a block that has a variable name of type String
x #(name String)
// can be assigned the value of `create_named()`
x = create_named() // Yes
// or the block created with the type Person
x = Person() // works too


// But the later is more natural to create custom types
// e.g Using a type `Person`
greet #(p Person) = {
	print("Hello `p.name`") 
}
// vs 
// using the block signature `(name String)`
greet #(p #(name String)) = {
	print("Hello `p.name`")
}

// To create an instance of the type invoke it as if it was a function and pass the parameters you need 
// Instead of retuning the last computer expresion it will return an instance of that tye 
Person : {
	name String
	last_name String
}
alice: Person("Alice" "Adams") 

// The same rules of named args and default values apply. 


```

## Generic types

When a type is a single uppercase letter, it represents a generic type

```js
Box: {
	data T
}
a: Box(1)
n Int = a.data // a.data is bound to an int
s String = a.data // compilation error

b: Box("Hola")
s String = b.data

// Then the generic type is unknow and a type needs to be declared, using () works as well as long is in the declaration part and the parameter is a type
ibox Box(Int) // declares a `box` of type Box of Integers
sbox Box(String) // declares a `box` fo type Box of String
gbox Box(T) // declares a box yet to be bound
```


## Concurrency 

When a block is invoked, it will run asynchronously if no value is assigned. 
If a value is assigned it will run synchronously. 

A block will finish its execution when all the inner blocks are completed, so all the blocks synchronize at the bottom of the enclosing block. 


```js
// Consumer / Producer example: 
// 
// Writes repeatedly to a buffer. 
produce: {
	buffer [String]
	count: 0
	while {true} {
		buffer.push("msg : `count`")
		count = count + 1
	}
}
// Reads continously from a buffer
// If the buffer is empty then just keeps trying
consume: {
	buffer []String
	while { true } {
		if buffer.len() == 0 { 
			continue
		} {
			println("You said: `buffer[0]`")
			buffer.shift()
		}
	}
}
buffer : []String
produce(buffer)
consume(buffer)
```

## Single owner and single writer

When a variable is defined inside a block, that block becomes the only owner and  only writer to that variable, even when some other blocks write outside of it or kept references to it;  internally the actual writer and owner is still  the creating block. 

```js
// Counter object
counter: {
	// Counter creates the variable `n`
	n Int = 0
}
// Even when other blocks try to write to it, internally the one that writes it is `counter`
one : {
	counter.n = counter.n + 1
}
two : {
	n: counter.n // two.n is not created here, is a reference to `counter.n`
	// so here `counter` is still the writer of `n` under the hood
	n = n + 1
}
// run blocks `one` and `two` concurrrently
one()
two()
// counter.n = counter.n + 1 is internally translated to something like 
// counter.set_n_value( some_value )
```

Is important to keep in mind this could make harder for a block to be garbage collected if other block have references to their values. To ensure this doesn"t happen make sure to copy the values instead of just referencing them 
```js
counter: {
	n Int = 0
}
two: {
	n std.copy( counter.n ) // makes a copy of `counter.n`
	n = n + 1 // two is the actual writer here
}
```


## Control flow

Most of the flow can be left to the types and their implementation. For instance the `Bool` class has methods that allows you to decide between two options or the `Option` boc would let you encapsulate a result

`Bool` example 
```js
// Returns a Bool instance 
f #(Bool) {
   1 < 2
}
r: f() 
// the method `?` decides between two bocs
r  ? {
   println("it was true")
}, {
   println("it was false")
}
```

`Option` example
```js

o #(Option(String)) {
    Some("hi")
}
f #(Option(String)) {
   None()
}

// the `or` method returns the value or an alternative
println(o.or("bye"))
```

But the main control structure will be most of the times a [Conditional Bocs](Conditional%20Bocs.md) which lets you evaluate a condition and execute the expressions next to the arrow `=>`

```js
match 
{ cond_1() => action_1() },
{ cond_2() => action_2() },
{ default_action() }
```

When passing a parameter to `match` the type varian will be matched

```js
x Option(String) = ...
match x 
	{ Option.Some() => print("We have `x.value`")},
	{ Option.None() => print("We have nothing")}

```
