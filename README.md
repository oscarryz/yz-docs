# The Yz Programming Language

> **Note**: Yz is currently in the design phase with no compiler implementation yet. All examples and features described here represent the intended design.

## Table of Contents

1. [Quick Example](#quick-example)
2. [Core Concepts](#core-concepts)
3. [Basic Syntax](#basic-syntax)
4. [Blocks of Code (Bocs)](#blocks-of-code-bocs)
5. [Types and Variables](#types-and-variables)
6. [Creating New Types](#creating-new-types)
7. [Generics](#generics)
8. [Type Variants](#type-variants)
9. [Structural Typing](#structural-typing)
10. [Arrays and Dictionaries](#arrays-and-dictionaries)
11. [Concurrency](#concurrency)
12. [Error Handling](#error-handling)
13. [Control Flow](#control-flow)
14. [Non-Word Method Invocation](#non-word-method-invocation)
15. [Info Strings](#info-strings)
16. [Code Organization](#code-organization)
17. [Examples](#examples)
18. [Reserved Words and Symbols](#reserved-words-and-symbols)
19. [Design Philosophy](#design-philosophy)

## Quick Example

```javascript
// Factorial in Yz
factorial: { n Int
  n > 0 ? { n * factorial(n - 1) }
          { 1 }
}
print("`factorial(5)`")  // prints 120
```

Yz is a programming language that explores simplifying concurrency, data, objects, methods, functions, closures, and classes under a single artifact: **blocks of code**. The language aims to provide a unified abstraction that can serve multiple roles traditionally handled by separate language constructs.

## Core Concepts

### Blocks of Code (Bocs)

The fundamental unit in Yz is a **block of code** (boc). Everything is a block:
- Variables are blocks
- Functions are blocks  
- Objects are blocks
- Classes are blocks
- Modules are blocks
- Actors are blocks

A block is a series of expressions between `{` and `}`:

```javascript
{
  message: "Hello"
  recipient: "World"
  print("`message`, `recipient`!")
}
```

### Unified Abstraction

Blocks serve multiple roles:

```javascript
// As data (object-like)
person: {
  name: "Alice"
  age: 30
}

// As behavior (function-like)  
greet: {
  name String
  print("Hello, `name`!")
}

// As both (object with methods)
counter: {
  count: 0
  increment: {
    count = count + 1
  }
}
```

## Basic Syntax

### Comments

```javascript
// Single line comment

/* 
   Multiline comment
*/
```

### Variables

```javascript
// Long form declaration
message String = "Hello"

// Short form with type inference
name: "World"

// Type declaration without initialization
age Int
```

### String Interpolation

Use backticks for string interpolation:

```javascript
name: "Alice"
greeting: "Hello, `name`!"  // "Hello, Alice!"
```

## Blocks of Code (Bocs)

### Basic Block Structure

```javascript
// A simple block
{
  a: 1
  b: 2
  a + b  // Last expression(s) are the "return value"
}
```

### Assigning Blocks to Variables

```javascript
calculator: {
  a: 0
  b: 0
  add: {
    a + b
  }
}
```

### Executing Blocks

Use `()` to execute a block:

```javascript
result: calculator()  // Executes the block
calculator.a = 5      // Access variables
calculator.b = 3      // Access variables  
sum: calculator.add() // Call methods
```

### Block Variable Access

Block variables can be accessed using `.` notation and modified before execution:

```javascript
greet: { 
  message String = "Hello"
  to_whom: "World"
  print("`message`, `to_whom`")
}

// Change variables before execution
greet.to_whom = "Everybody"
greet() // prints "Hello, Everybody!"

// Variables can be accessed even after execution
greet.message // returns "Hello"
```

### Block Parameters

Variables in blocks serve as both parameters and return values:

```javascript
multiply: {
  x Int
  y Int
  x * y
}

result: multiply(5, 3)  // result = 15
```

### Named Parameters

```javascript
divide: {
  numerator Int
  denominator Int
  numerator / denominator
}

result: divide(numerator: 10, denominator: 2)  // result = 5
```

### Default Values

Parameters can have default values:

```javascript
f: {
  a Int
  b: 0  // b has default value of 0
}

f(1)     // b defaults to 0
f(1, 2)  // assigning 2 to b
```

### Multiple Return Values

The last n expressions can be used as return values:

```javascript
swap: {
  a String
  b String
  b  // Second to last
  a  // Last
}

x: "hello"
y: "world"
x, y = swap(x, y)  // x = "world", y = "hello"
```

## Types and Variables

### Built-in Types

```javascript
// Numbers
n Int = 42
m : -1
pi Decimal = 3.14

// Strings
message String = "Hello"
name: "World"  // Type inferred

// Booleans
flag Bool = true

// Arrays
numbers [Int] = [1, 2, 3]
words: ["hello", "world"]

// Dictionaries
ages [String:Int] = ["Alice": 30, "Bob": 25]

// Bocs
greet #(msg String,String) {
    "Hello `msg`"
}

hi: {
   42
}
```

### Block Signatures

Block types are defined by their signature:

```javascript
// Block that takes two Ints and returns an Int
add #(x Int, y Int, Int) {
  x + y
}

// Block with just return type
get_answer #(Int) { 42 }

// Block with no parameters or return
do_something #() {
  print("Done")
}
```

The signature is inferred when using the short declaration + assignment operator `:`

```js
add : {
   x Int
   y Int
   x + y
}
get_answer : { 42 }
do_something: {
   print("Done")
}
```

### Block Signature Declaration and Assignment

Blocks can be declared with signatures and assigned later:

```javascript
// Declare signature
greet #(message String, to_whom String, String)

// Assign implementation later
greet = {
  message String
  to_whom String
  message // just returns the variable message
}

// Type signature can omit variable names
greet #(String, String, String) 
greet() // compilation error, needs to be assigned a value

greet = {
  a String 
  b String
  c String
}
greet() // compilation error, a, b, c need a default value or one assigned
greet("uno", "dos", "tres")
```

## Creating New Types

### Type Declaration

Uppercase names define new types:

```javascript
Person : {
  name String
  age Int
  greet: {
    print("Hello, I'm `name`")
  }
}
```

### Creating Instances

```javascript
alice: Person(name: "Alice", age: 30)
// or
bob: Person("Bob", 25)

alice.greet()  // "Hello, I'm Alice"
```

### Multi-field Types

```javascript
Person : {
  name String
  last_name String
}
alice: Person("Alice", "Adams") 

// The same rules of named args and default values apply
```

### Blocks Returning Blocks

A block can return another block:

```javascript
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

### Type vs Block Factory

The difference between declaring a new type and returning a block:

```javascript
// Type declaration
Person: {
  name String
}

// Block factory
create_named: {
  {
    name String
  }
}

// Both create a copy of a block with a String
// x is a block that has a variable name of type String
x #(name String)
// can be assigned the value of `create_named()`
x = create_named() // Yes
// or the block created with the type Person
x = Person() // works too

// But the latter is more natural to create custom types
// e.g., using a type `Person`
greet #(p Person) = {
  print("Hello `p.name`") 
}
// vs 
// using the block signature `(name String)`
greet #(p #(name String)) = {
  print("Hello `p.name`")
}
```

### Type Signatures for Custom Types

```javascript
// Explicit signature
Point #(x Int, y Int) {
  distance_to_origin: {
    sqrt(x * x + y * y)
  }
}
```

## Generics

Single uppercase letters represent generic types:

```javascript
Box: {
  data T  // T is generic
}

int_box: Box(42)        // T becomes Int
string_box: Box("Hi")   // T becomes String
```

### Generic Functions

```javascript
identity: {
  value T
  value  // Returns whatever type was passed in
}

number: identity(42)    // number: Int
text: identity("hi")    // text: String
```

### Constrained Generics

The constraints are inferred from the usage:
```javascript
printable: {
  value T  // T must have a print method
  value.print()
}

Person: {
   name String
   print : {
     println("My name is `name`")
   }
}
printable(Person("Yz"))
printable("oh oh") // "oh oh" doesn't have a `print` block
```

## Type Variants

Type variants provide sum type functionality:

```javascript
Option: {
  T
  Some(value T),
  None()
}

maybe_number: Option.Some(42)
nothing: Option.None()

// Pattern matching with match
result: match maybe_number {
  Some => "Got value: `maybe_number.value`"
}, {
  None => "No value"
}
```

### More Complex Variants

```javascript
Result: {
  T, E
  Ok(value T),
  Err(error E)
}

NetworkResponse: {
  Success(data String),
  Failure(error String),
  Timeout()
}

handle_response: {
  response NetworkResponse
  match response {
    Success => print("Data: `response.data`")
  }, {
    Failure => print("Error: `response.error`")  
  }, {
    Timeout => print("Request timed out")
  }
}
```

## Structural Typing

Yz uses structural typing - types match based on structure, not names:

```javascript
Point: {
  x Int
  y Int
}

Vector: {
  x Int
  y Int  
}

process_coordinates: {
  coords #(x Int, y Int)  // Any type with x, y Int fields
  coords.x + coords.y
}

p: Point(3, 4)
v: Vector(1, 2)

process_coordinates(p)  // Works - Point has x, y Int
process_coordinates(v)  // Works - Vector has x, y Int
```

## Arrays and Dictionaries

### Arrays

```javascript
// Type declaration
a [Int]
// initialization
a = [1, 2, 3]

// decl + init
a [Int] = [1, 2, 3]

// short declr + init
a : [1, 2, 3]

// empty decl + init
a [Int] = []Int // Is an empty array
// short declr + init
a : []Int // empty array of ints 

// Generic
a [T] = [1, 2, 3]
a : []T

// Array operations
a << 'Hello' // or a.add 'Hello'
print(a[0]) // access element 0 of the array
a[0] = "Hola"
```

### Dictionaries (Associative Arrays)

```javascript
// Type
[Key_Type : Value_Type] 

// declaration
d [String:Int] 
// initialization 
d = [ "one": 1, "two": 2]

// decl + init 
e [String:Int] = ["one":1, "two":2]

// short decl + init
f : ["one":1, "two":2 ]

// empty 
g2 [String:Int] = [String]Int
// short decl + init empty
g1 : [String]Int

// generic + initialization
g3 [K:V] = [String]Int 
g4 [K:V]
g4["hello":1]

// Dictionary access returns Optional(V)
d : [ 1 : 2, 3: 4] // [Int: Int]
d[1] // Some(2)
d[5] // None()
```

## Concurrency

### Async by Default

Every block call is asynchronous, the value will be resolved by the time it is used:

```javascript
// These run concurrently
fetch_user("alice")
fetch_orders("alice")

user: fetch_user("alice")  
print(user) // might be resolved by then, it will block if it hasn't completed.
```

### Structured Concurrency

Blocks synchronize at the end of their enclosing scope:

```javascript
process_data: {
  // Both operations start concurrently
  img: fetch_image("123")
  usr: fetch_user("alice")
  
  // The `process_data` will not complete
  // until `create_profile` completes. 
  // It will block execution until it does. 
  create_profile(img, usr)
}
```

### Single Writer Principle

Each variable has only one writer (its declaring block):

```javascript
counter: {
  count: 0  // counter owns count
}

increment: {
  counter.count = counter.count + 1  // Modified through counter
}

decrement: {
  counter.count = counter.count - 1  // Also modified through counter  
}
```

## Error Handling

Yz uses `Result` and `Option` types for error handling:

```javascript
divide: {
  a Int
  b Int
  b == 0 ? {
    Err("Division by zero")
  } {
    Ok(a / b)
  }
}

result: divide(10, 2).or_else({
  error Error
  print("Error: `error`")
  0  // Default value
})
```

### Chaining Operations

```javascript
process_file: {
  filename String
  read_file(filename)
    .and_then { content String
      parse_content(content)
    }
    .and_then { data Data
      validate_data(data)  
    }
    .or_else { error Error
      print("Processing failed: `error`")
    }
}
```

## Control Flow

### Conditional Expressions

```javascript
max: {
  a Int
  b Int
  a > b ? { a },{ b }
}

status: user.age >= 18 ? { "adult" },{ "minor" }
```

### Bool Type Control Flow

The `Bool` type has methods that allow you to decide between two options:

```javascript
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

### Option Type Control Flow

```javascript
o #(Option(String)) {
    Some("hi")
}
f #(Option(String)) {
   None()
}

// the `or` method returns the value or an alternative
println(o.or("bye"))
```

### Match expressions

```javascript
describe_number: {
  n Int
  match  { 
	  n < 0  => "negative"
  },{ 
	  n == 0 => "zero"
  },{ 
	  n > 0  => "positive" 
  }
}
```

### Match with Type Variants

When passing a parameter to `match` the type variant will be matched:

```javascript
x Option(String) = ...
match x 
  { Option.Some() => print("We have `x.value`")},
  { Option.None() => print("We have nothing")}
```

### Generic Match Syntax

```javascript
match 
{ cond_1() => action_1() },
{ cond_2() => action_2() },
{ default_action() }
```

### Loops

```javascript
// Range iteration
1.to(10).each({ i Int
  print("`i`")
})

// Array iteration
names: ["Alice", "Bob", "Charlie"]
names.each({ name String
  print("Hello, `name`!")
})

// "While" loops
factorial: {
  n Int
  result: 1
  current: n
  while({ current > 1 }, {
    result = result * current
    current = current - 1
  })
  result
}
```

### Control Flow Keywords

```javascript
// Early return
check: {
   age Int
   age < 21 ? {
       message: 'You have to be over 21 to access this site'
	   return
   }
   message: 'Welcome to the site'
}

// Break from loops
max_from_list: {
  list []
  m Int
  list.for_each {
    item Int
    item < 0 ? {
      break // will exit the loop 
    }
    m = max(max, item)
  }
  m
}

// Continue in loops
n Int
match {
   n % 3 == 0 => print("Fizz")
   continue // evaluated the following condition
}, {
	n % 5 == 0 => print("Buzz")
}, {
	print("`n`")
}
```

## Non-Word Method Invocation

When boc name is non-word, we can invoke it without `.` and `parenthesis` as long as it has at least one parameter.

```js
Example: { 
   // the "<<" variable
   << : { 
	   n Int
	   printnln(n)
   }
}
e : Example()
e << 1 // same as e.<<(1)
```

## Info Strings

Information Strings. You can add a string before any element and will be available by calling `std.info(element)` 

```javascript
`A message`
message: 'Hello'
info(message).text // A message
```

You can add blocks there too, these block don't need to be valid yz code, but different tools might require them to be.

```javascript
`
Prints the classics "Hello, World!" program to the screen

variables: {
  what String = 'Hello' // what message to display
  who  String = 'World' // what to say
}

tests: {
  assert say_hello()            == "Hello, World!"                   // Uses defaults      
  assert say_hello 'Hola'       == "Hola, World!"                    // Overrides first variable 'what'
  assert say_hello who: 'there' == "Hello, there!"                   // Explicitly overrides variable 'who'
  assert say_hello who: 'home' what: 'Welcome'  == "Welcome, home!"  // "Named parameters"
}
version: 1.0
author: 'Yz developers'
`
say_hello: {
   // Any element can have an info string
  'What message to display'
  what: 'Hello' 

   // Can have validation info
   // Or serialization info
  `
   validation: "\w*"
   json_field: 'xyz'
  `
  who:  'World' 
  // Could also be used as running examples
  // that will be validated with yzc test  
  `
  Example: print 'Hello' 'world'
  Output: Hello, world!
  `
  print '{what}, {who}!'
}
```

To retrieve it use the `std.info` block and pass the element

```javascript
info: std.info say_hello
print info.text  // Prints the classics "Hello, World!"... etc.etc
info.tests()     // Runs the tests
info.version     // 1.0
info.examples()  // run the examples 
```

## Code Organization

### Simple Projects

For simple projects, the `yz` build tool will compile each individual file and will create an executable if they are named `main.yz`, have a `main` method, or have free floating code. You can also pass the filename to process a single file. If no entry point is found, they are considered libraries and no executable will be created.

### Larger Projects

If a `yz` file contains a `configuration` structure, then it will be used to create the executable. You can also create this configuration structure by invoking the build tool `init` _`project_name`_, which can create additional folders: 

```
yz init project_name
project_name/
   project_name.yz
   doc/
   src/
   lib/
   test/
   vendor/
```

### Configuration File

A configuration file is a `.yz` file that contains information like version, entry point and dependencies:

```javascript
version: '0.1.0'
entry: 'main.yz'
src_path: ['./src/' './vendor/' './lib/']
vendor: ['./vendor']
dependencies: []
```

### Dependencies

To add a dependency use `yz install printer`, the dependency will be added to the dependency section

```js
dependencies: [
     printer: {version: "1.0.0"   url: 'https://example.org/print.git'} 
]
```

### Filesystem Block Name Resolution

The compiler resolves block names to filesystem files using the `src_path` variable using the following strategy:

1. File name including subdirectories will create a block, even if it is empty, excluding directories defined in the project's `src_path` variable
   `./src/house/front/Host.yz` will create the `house.front.Host` block

2. Files and subdirectories can be used to create blocks and nested blocks for better code organization.

## Examples

### Counter with State

```javascript
Counter: {
  count Int = 0
  
  increment: {
    count = count + 1
  }
  
  decrement: {
    count = count - 1
  }
  
  get: { count }
}

counter: Counter()
counter.increment()
counter.increment()
print(counter.get())  // prints 2
```

### Binary Tree

```javascript
Tree: {
  T
  Empty(),
  Node(value T, left Tree(T), right Tree(T))
  
  insert: {
    value T
    match {
      Empty() => Node(value, Empty(), Empty())
    }, {
      Node() => value < self.value ? {
        Node(value, left.insert(value), right)
      } {
        Node(value, left, right.insert(value))
      }
    }
  }
}

tree: Tree.Empty()
tree = tree.insert(5).insert(3).insert(7)
```

### HTTP Server Concept

```javascript
Server: {
  port Int
  routes [String:#(Request, Response)]
  
  listen: {
    // Server implementation would go here
    print("Server listening on port `port`")
  }
  
  route: {
    path String
    handler #(Request, Response)
    routes[path] = handler
  }
}

server: Server(8080)
server.route("/hello", {
  request Request
  Response(body: "Hello, World!")
})
server.listen()
```

## Reserved Words and Symbols

The following cannot be identifiers or part of an identifier:

```
break
continue
return
match
=>
:
`
'
"
[]
{}
()
, ; . #
```

`=` might be part of an identifier, but there are also `=` and `==` operators.

## Design Philosophy

Yz operates on the principle that most programming constructs can be unified under the concept of blocks of code. A block can:
- Store data (like objects/structs)
- Execute code (like functions/methods)
- Run concurrently (like actors)
- Define types (like classes)

This documentation provides a comprehensive overview of the Yz programming language design. The language aims to simplify concurrent programming while maintaining type safety through its innovative "blocks of code" abstraction that unifies many traditionally separate language constructs.