#doc
# The Yz Programming Language Documentation

> **Note**: Yz is currently in the design phase with no compiler implementation yet. All examples and features described here represent the intended design.

## Table of Contents

1. [Introduction](#introduction)
2. [Core Concepts](#core-concepts)
3. [Basic Syntax](#basic-syntax)
4. [Blocks of Code (Bocs)](#blocks-of-code-bocs)
5. [Types and Variables](#types-and-variables)
6. [Creating New Types](#creating-new-types)
7. [Generics](#generics)
8. [Type Variants](#type-variants)
9. [Structural Typing](#structural-typing)
10. [Concurrency](#concurrency)
11. [Error Handling](#error-handling)
12. [Control Flow](#control-flow)
13. [Standard Library](#standard-library)
14. [Examples](#examples)

## Introduction

Yz is a programming language that explores simplifying concurrency, data, objects, methods, functions, closures, and classes under a single artifact: **blocks of code**. The language aims to provide a unified abstraction that can serve multiple roles traditionally handled by separate language constructs.

### Quick Example

```javascript
// Factorial in Yz
factorial : {
   n Int 
   match {
	   n > 0 => n * factorial(n - 1)
   }, {
	   1
   }
}
print("`factorial(5)`")  // prints 120
```

### Design Philosophy

Yz operates on the principle that most programming constructs can be unified under the concept of blocks of code. A block can:
- Store data (like objects/structs)
- Execute code (like functions/methods)
- Run concurrently (like actors)
- Define types (like classes)

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
  match reponse {
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
  
  // The `process_data` will complete
  // until `create_profile` complets. 
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

## Standard Library

### Collections

```javascript
// Arrays
numbers: [1, 2, 3, 4, 5]
doubled: numbers.map({ n Int; n * 2 })
evens: numbers.filter({ n Int; n % 2 == 0 })
sum: numbers.reduce({ acc Int, n Int; acc + n })

// Dictionaries  
config: ["host": "localhost", "port": "8080"]
host: config["host"]  // Option type
```

### Option and Result Types

```javascript
// Option for nullable values
find_user: {
  id String
  // Returns Option(User)
  users.find({ user User; user.id == id })
}

user: find_user("123").or_else({ User("default") })

// Result for error handling
safe_divide: {
  a Int
  b Int
  match {
	  b == 0 => Err("Division by zero")
  }, {
	  Ok(a / b)
  }
}
```

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

This documentation provides a comprehensive overview of the Yz programming language design. The language aims to simplify concurrent programming while maintaining type safety through its innovative "blocks of code" abstraction that unifies many traditionally separate language constructs.

