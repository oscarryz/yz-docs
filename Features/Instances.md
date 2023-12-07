Replaced with [Creating instances](../Questions/solved/Creating%20instances.md)

A new type can be created by creating an alias of a type signature. 

First let's see what a type signature is and how to use it
```javascript
// Type signature: Is a block with a variable count of type Int
{count Int} 

// Can declare a var with that type signature
b1 {count Int}
// Can initialize a block matching that signature
b1 = {
    x: 1
    y: 2
    count: x + y // count of type Int declared and signature satisfied
}
b1() // returns count, 3
```

Usually the type signature can be inferred from the content
```javascript
b1: {
    count: 10 // b1 type signature is {count Int}
}
b2: {
    10 // b2 type signature is {Int}
}

```

To create copies from a block create a ~~type signature alias using the `::` operator and then assign a block to that alias~~  (there's no alias)

```javascript
b1: {count: 10}
Counter { count Int }
Counter = b1 // assign the body of b1 `{count:10}` to the type alias `Counter`
// To create a copy use the instantiation construct
c1: Counter() // creates an instance
b1.count // 10
c1.count // 10
c1.count = 12 // 12
b1.count // still 10 
```
Usually you assign the block directly to the alias

```javascript
Counter { 
    count: 10
}
c1: Counter()
b1: Counter()
b1.count // 10
c1.count //10
c1.count = 12 // 12
b1.count // still 10
```

// Old info follows, might be outdated


Creating instances of a block (objects)

You can create a copy of a block, by using the block literal `{}` 

```javascript
thing: {
    name String
}
// Regular without copy 
a: thing('Ping') // execute with param 'Ping' and value is returned so `a` is 'Ping'
a // Ping
thing.name // Ping
b: thing('Pong')
a // Still Ping
thing.name // Pong
b // Pon

// This is fine because we're executing the block and then taking the value
// `thing.name` changes reference but `a` still refers to the origina.

// With copy
a: thing{'Ping'} // creates a copy of `thing` so `a` now has a `name` variable with the value "Ping"
a.name // Ping
b: thing{'Pong'} // creates a copy of 'thing' so `b` now has a `name` variable with the value "Pong"
b.name // Pong
thing.name // compiler error, is not defined, only the copies were affected

```

However, we cannot use a block name as type: 
```javascript
a thing // a type thing is not valid
// we have to use the full block type
a {name String} // a is a block that has a name variable of type String
// Or we can define an uppercase named block `Thing` and then use that
Thing: {
    name String
}
a Thing // now we can declare a of type Thing

// What's the type of `Thing` though? 
// The following is not valid because that's a Declaration and Instantiation
Thing { name String } = { 
    name String
}
thing {name String String} = {
    name String
    name
}

```


```javascript

    // Definition 
	Person {
		name String
		email String
		dob Date
		say_hi: {
			print('Hi')
		}
	}

	 p: Person (
		 name: 'Oscar'
		 email: 'email@oscarryz.com'
		 dob: Date { 01  01  1980}
		 say_hi: { println 'yes' }
	 )
    // also 
    p: Person (
        name: 'Foo'
        email: 'bar'
        dob: {
            day: 01
            month: 01
            year: 1980
        }
        say_hi:{ print 'yes' }
    )
     p: Person (
        name: 'Foo'
        email: 'bar'
        dob: {
            01
            01
            1980
        }
        say_hi:{ print 'yes' }
    )
```

April 20th 2023

~~The upper case thing is just a convention, objects can be created with closures~~
This is another way, by creating anonymous blocks that capture the environment of the enclosing block, but this style might be too cumbersome, yet, valid.

```javascript
point: {
   x Int
   y Int
   {x:x;y:y}
}
p: point(10 20) // returns the last executed
p.x  // 10
p.y  // 20

thing: {
   a_name String
   {
      name: a_name
   }
}
duck: thing("Duck")
duck.name // 

```

#open-question How using closures fits with the whole language? 
Answer, it doesn't 


```javascript
// Not valid nor wanted
person: {
    Person {
        _name String
        name {String}
        greet {Person String}
        builder {String {String}}
        call_me {{}}
        ok      {{{}}}
    }
}
// Nope
p: person()(
    _name: 'Oscar'
    name: {s String}
    greet: {p Person; "Hola {p.name()}"}
    builder: {s String; {s}}
    call_me: {{print 'OK'}}
    ok     : { do {} { print 'Ok'; do()}}
)
```
Correct
```javascript
// Alias
Person {
        _name String
        name {String}
        greet {Person String}
        builder {String {String}}
        call_me {{}}
        ok      {{{}}}
} 
p: Person {
    _name: 'Oscar'
    name: {s String}
    greet: {p Person; "Hola {p.name()}"}
    builder: {s String; {s}}
    call_me: {{print 'OK'}}
    ok     : { do {} { print 'Ok'; do()}}
}
{name String}{'Oscar'}
```

So, if creating block instances is conveyed by placing {} after the block, are types still needed? 

```javascript
person: {
    name String
}
alice: person{}('Alice' )
bob: person{}('Bob')
// What if it returns a block
person: {
    name String
    {name:{name}}
}
alice: person('Alice')
alice.name() // returns Alice
bob: person('Bob')
alice.name() // returns Bob because person.name was set to Bob
// To avoid this we can create a copy of the block:
alice: person{'Alice'} // {} creates a copy... Alice is the first variable
//alice: person{name:'Alice'}  same as above
alice.name() // err name is a String
bob: person{'Bob'}
// definition should've been: 
person: {
    name String
    get_name: {name}
}
alice: person{name:'Alice'}
alice.get_name() // Alice
bob: person{name:'Bob'}
bob.get_name() // Bob
person.name // compiler error ? 
// Which is exactly what this is: 
Person: {
    name String
    get_name: {name}
}

```

### Good

```javascript
a:thing{}
a.name 
// a type 
a thing // not possible no differentiation from calling a(thing)
a Thing // Yes, this is what I've been using
```

### Problem

I cannot differentiate instantiation from type
```javascript
a {name String} // a is a block of name and String or a is instntiating?
```

Possible solution: 

1. Note, the block is only defining types `name String` and not code `{name: 'Alice'}`
2. Use `::` everywhere 
    1. `a::{name String}`
    but it might be problematic to use as parameters, lets see

```js
do: {
    action::{}
    action()
}
// vs
do: {
    action {}
    action()
}
// but if is not a block?
do: {
    action::Int
    action = action + 1
}
```
it's kind of terrible. Might work if used only to define new types and not parameters, but then there's the problem there's no difference between parameters and local variables

More on option 1

```js
do0: {
    action {}
    action()
}
do1: {
    action {String}
    action {} // coping the action? is action a block?
    action() // action returns a String
}
```

I thing the uppercase is the thing along with the `::`

Creating instances:

```javascript
counter: {
    count Int
    increment: {count = count + 1}
}
// types
counter {count Int increment {}}
// creating an instance
counter{} // not possible? 
counter::{count Int increment{}} // now it is instantiable
counter{} // yes, 
counter::{count Int increment {} } = {
    count Int
    increment: {count = count + 1}
}
c: counter{10}


```

#open-question  How to define a new type and instantiate without using `{}` everywhere 
Still nothing decided but so far I think:


`Upper{}` will create an instance
`Upper::{}` will define an instantiable type
`lower {}`  will define a variable of type block
`lower{}`  maybe will create an instance
`lower::{}` maybe will define an instianble type

Things to avoid: Using `::` for every type declarations e.g. 

```
block: {
    variable :: {Int Int Int} 
    int_variable :: Int
}
```


```javascript
// Type
i Int
b Bool
d Decimal
f {} // block takes and receives nothing
c::{} // instantiable block that takes and receives nothing

i = 1
b = true
d = 0.0
f = {} // block used as block literal 
c = {} // same, but this is instantiable 
z {} // defines a block 
z = c{} // z is a copy of c
y c // y is of type c which is the same as {}
y = z 
y = c{} 
```


```javascript
get_age :: { get_age {Int} }
x: get_age()
Person::{ 
    age:0 
    get_age: {age}
}
print_age: {
    ager get_age
    ager.get_age()
}
```

The reason why `Person:{}` is not enough to create a new instantiable  is because its type is the same as just a regular block `b:{}` both have type `{}` making it indistiguisable? 
Also using `:` alone would create a variable (although uppercase would desamabiguate).


[Type Alias (wip)](Type%20Alias%20(wip).md)