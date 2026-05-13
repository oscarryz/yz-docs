#readme 
# The Yz Programming Language

> **Note**: [Yz compiler](./compiler) is work in progress. All examples and features described here represent the intended design.


## Quick Example

```javascript
// Factorial in Yz
factorial: { n Int
  n > 0 ? { n * factorial(n - 1) },
          { 1 }
}
print("${factorial(5)}")  // prints 120
```

Yz is a programming language built around a single construct: the **block of code** (boc). Variables, functions, objects, types, modules, concurrent behaviours, and protected resources are all blocks. Rather than separate constructs for each role, you compose everything from one idea.

A block is a series of expressions between `{` and `}`, and the same block can act as data, behaviour, or both:

```javascript
// As data
person: {
  name: "Alice"
  age: 30
}

// As behaviour
greet: {
  name String
  print("Hello, ${name}!")
}

// As both
counter: {
  count: 0
  increment: { count = count + 1 }
}
```

## Basic Syntax

### Comments

*→ [Details](docs/Features/Comments.md)*

```javascript
// Single line comment

/* 
   Multiline comment
*/
```

### Variables

*→ [Details](docs/Features/Variables.md)*

```javascript
// Long form declaration
message String = "Hello"

// Short form with type inference
name: "World"

// Type declaration without initialization
age Int
```

### Strings

*→ [Details](docs/Features/Strings.md)*

Both `"double"` and `'single'` quotes create strings; they are interchangeable:

```javascript
a: "Hello"
b: 'Hello'  // identical
```

### String Interpolation

*→ [Details](docs/Features/String interpolation.md)*

Use `${...}` inside a string literal for interpolation:

```javascript
name: "Alice"
greeting: "Hello, ${name}!"   // "Hello, Alice!"
greeting: 'Hello, ${name}!'   // same
```

## Blocks of Code (Bocs)

*→ [Details](docs/Features/Bocs.md)*

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
  print("${message}, ${to_whom}")
}

// Change variables before execution
greet.to_whom = "Everybody"
greet() // prints "Hello, Everybody!"

// Variables can be accessed even after execution
greet.message // returns "Hello"
```

### Block Parameters and Return Values

In Yz there is no separate concept of "parameter" or "return value" — they are just variables. A variable declared without a value is a required input; one declared with a value is optional (defaults apply). The last expression(s) in the body are the output.

```javascript
greet: {
  message String       // required — caller must provide
  to_whom: "World"     // optional — defaults to "World"
  "${message}, ${to_whom}!"  // return value
}

greet("Hello")            // "Hello, World!"
greet("Hi", "Alice")      // "Hi, Alice!"
greet(to_whom: "Bob", message: "Hey")  // named args, any order
```

Because parameters are fields, they are accessible before and after the call:

```javascript
greet.to_whom = "Everyone"
greet("Hello")    // "Hello, Everyone!"
greet.message     // "Hello" — readable after call
```

The last N expressions are the return values — no `return` keyword needed:

```javascript
swap: {
  a String
  b String
  b    // second-to-last — first return value
  a    // last — second return value
}

x, y = swap("hello", "world")  // x = "world", y = "hello"
```

## Concurrency

*→ [Details](docs/Features/Concurrency.md)*

### Async by Default

Every block call is asynchronous. The value is resolved by the time it is used:

```javascript
// These run concurrently
fetch_user("alice")
fetch_orders("alice")

user: fetch_user("alice")
print(user) // blocks here only if fetch_user hasn't completed yet
```

### Structured Concurrency

A boc does not complete until all bocs it spawned have completed:

```javascript
process_data: {
  // Both operations start concurrently
  img: fetch_image("123")
  usr: fetch_user("alice")

  // process_data will not complete until create_profile completes
  create_profile(img, usr)
}
```

### Exclusive Access (BOC Model)

Every value in Yz is a protected concurrent resource — a **cown** (concurrent owner). Only one running boc can hold a cown at a time; all others queue behind it. Cowns are acquired atomically: a boc that needs multiple resources gets all of them at once or waits until it can.

```javascript
Account: {
  balance Int
  balance+= #(amount Int) { balance = balance + amount }
  balance-= #(amount Int) { balance = balance - amount }
}

// transfer acquires src and dst atomically before running
transfer #(src Account, dst Account, amount Int) {
  src.balance >= amount ? {
    src.balance-=(amount)
    dst.balance+=(amount)
  }, {
    print("insufficient funds")
  }
}

main: {
  alice: Account(100)
  bob:   Account(0)

  transfer(alice, bob, 30)  // acquires alice + bob
  transfer(bob, alice, 10)  // waits — bob is taken by the first transfer
}
```

Two bocs that need **different** resources run in parallel automatically. Two bocs that share a resource serialize in the order they were spawned. No locks, no `synchronized`, no `async/await`.

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
greet: { msg String
    "Hello ${msg}"
}
hi: { 42 }
```

### Block Signatures

*→ [Details](docs/Features/Block type.md)*

There are four ways to declare a block, from most implicit to most explicit.

**1. Identifier + inferred boc** — signature is inferred from the body:

```javascript
add: {
  x Int
  y Int
  x + y       // return type Int inferred from last expression
}
get_answer: { 42 }
do_something: { print("Done") }
```

**2. Identifier + explicit boc type** — declares the signature without a body (assigned later):

```javascript
add #(x Int, y Int, Int)   // takes two Ints, returns Int
get_answer #(Int)           // no inputs, returns Int
do_something #()            // no inputs, no return
```

**3. Identifier + explicit boc type + body (assigned separately):**

```javascript
add #(x Int, y Int, Int)
add = {
  x Int
  y Int
  x + y
}
```

**4. Typed block — identifier + explicit signature + inline body:**

```javascript
add #(x Int, y Int, Int) {
  x + y
}
get_answer #(Int) { 42 }
do_something #() { print("Done") }
```

This last form is the most common for named blocks that need an explicit signature.

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

*→ [Details](docs/Features/Define new types.md)*

### Type Declaration

Uppercase names define new types:

```javascript
Person : {
  name String
  age Int
  greet: {
    print("Hello, I'm ${name}")
  }
}
```

### Creating Instances

*→ [Details](docs/Features/Create instances.md)*

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

Constraints are inferred from usage by default:

```javascript
printable: {
  value T  // T must have a print method — inferred from usage below
  value.print()
}

Person: {
   name String
   print: {
     print("My name is ${name}")
   }
}
printable(Person("Yz"))
printable("oh oh") // error: String doesn't have a `print` block
```

Constraints can also be declared explicitly as an optional annotation:

```javascript
serialize: {
  value T Serializable  // T must satisfy the Serializable interface
  value.to_json()
}
```

An explicit constraint is checked at the call site; an inferred constraint is checked at usage inside the body. Both forms are valid.

## Type Variants

*→ [Details](docs/Features/Type variants.md)*

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
  Some => "Got value: ${maybe_number.value}"
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
    Success => print("Data: ${response.data}")
  }, {
    Failure => print("Error: ${response.error}")  
  }, {
    Timeout => print("Request timed out")
  }
}
```

## Structural Typing

*→ [Details](docs/Features/Structural typing.md)*

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

*→ [Details](docs/Features/Array.md)*

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
a [Int] = [Int]() // Is an empty array
// short declr + init
a : [Int]() // empty array of ints 

// Generic
a [T] = [1, 2, 3]
a : [T]()

// Array operations
a << 'Hello' // or a.add 'Hello'
print(a[0]) // access element 0 of the array
a[0] = "Hola"
```


### Dictionaries (Associative Arrays)

*→ [Details](docs/Features/Associative arrays.md)*

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
g2 [String:Int] = [String:Int]()
// short decl + init empty
g1 : [String:Int]()

// generic + initialization
g3 [K:V] = [String:Int]() 
g4 [K:V]
g4["hello":1]

// Dictionary access returns Optional(V)
d : [ 1 : 2, 3: 4] // [Int: Int]
d[1] // Some(2)
d[5] // None()
```

## Error Handling

*→ [Details](docs/Features/Error handling.md)*

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
  print("Error: ${error}")
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
      print("Processing failed: ${error}")
    }
}
```

## Control Flow

*→ [Details](docs/Features/Conditional Bocs.md)*

*→ [Details](docs/Features/return, break, continue.md)*

```javascript
// ? is a method on Bool — true-branch, false-branch
max: {
  a Int
  b Int
  a > b ? { a }, { b }
}

// match — first branch whose condition is true runs
describe: {
  n Int
  match {
    n < 0  => "negative"
  }, {
    n == 0 => "zero"
  }, {
    n > 0  => "positive"
  }
}

// match on a type variant
x Option(String) = ...
match x {
  Option.Some => print("Got ${x.value}")
}, {
  Option.None => print("Nothing")
}

// iteration
1.to(10).each { i Int; print("${i}") }

names: ["Alice", "Bob", "Charlie"]
names.each { name String; print("Hello, ${name}!") }

while({ current > 0 }, { current = current - 1 })

// return, break, continue work as in most languages
check: {
  age Int
  age < 21 ? { return }
  print("Welcome")
}
```

## Trailing Block Syntax

*→ [Details](docs/Features/Trailing block syntax.md)*

When the only argument to a method is a block literal, the parentheses can be omitted. Write the block directly after the method name on the same line:

```javascript
// Both are identical
list.filter({ item Int; item > 10 })
list.filter { item Int; item > 10 }

// Chaining
[1, 2, 3, 10, 20]
    .filter { n Int; n > 5 }
    .each   { n Int; print(n) }
```

The `{` must appear on the same line as the method name (a newline causes ASI to insert a semicolon, and the block becomes a separate statement).

## Non-Word Method Invocation

*→ [Details](docs/Features/Non-Word invocation.md)*

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

*→ [Details](docs/Features/Info strings.md)*

An infostring is a boc body delimited by backticks placed immediately before a definition. Its content is valid Yz — compiled but never executed — and can be used at compile time to augment or extend the language:

```javascript
`
compile_time: [JSON, Embed]
`
Movie : {
    title  String

    `json: "release_date"`
    year   Int

    `json: "ignore"`
    internal_id String

    `embed: "icon.png"`
    image Image
}
```

`compile_time` lists the extensions to run on the annotated boc. Each extension reads only the variable it owns (`json`, `embed`, …). Referenced names are resolved at compile time — a typo is a compile error.

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

### Bank Account Transfers

Five concurrent transfers. Some share accounts (serialized), others don't (parallel). No locks written anywhere.

```javascript
Account: {
  balance Int
  balance+= #(amount Int) { balance = balance + amount }
  balance-= #(amount Int) { balance = balance - amount }
}

transfer #(src Account, dst Account, amount Int) {
  src.balance >= amount ? {
    src.balance-=(amount)
    dst.balance+=(amount)
  }, {
    print("insufficient funds: need ${amount}, have ${src.balance}")
  }
}

main: {
  alice: Account(100)
  bob:   Account(0)
  carol: Account(50)

  transfer(alice, bob,   30)  // alice + bob
  transfer(bob,   alice, 10)  // serialized after above
  transfer(bob,   carol, 20)  // serialized after above
}
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
`
compile_time: [http.HttpServer]
port: 8080
`
Server: {

  `route: "/hello"`
  hello #(r Request, w Response) {
    Response(body: "Hello, World!")
  }

  `route: "/users/${id}"`
  get_user #(r Request, w Response) {
    id: r.params.id
    Response(body: "User: ${id}")
  }

  `route: "/users"; method: http.Post`
  create_user #(r Request, w Response) {
    Response(body: "Created")
  }
}

server: Server()
server.listen()
```

## Reserved Words and Symbols

*→ [Details](docs/Features/Reserved words and characters and symbols.md)*

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

## Repository Structure

- **`compiler/`** — Go implementation of the Yz compiler. Includes the lexer, parser, AST, lowerer, and code generator. Emits Go source and invokes `go build` to produce binaries.
- **`docs/`** — Additional documentation, design notes, and implementation decisions.
  - **[`docs/Features/`](docs/Features/README.md)** — Full feature reference, one page per language feature.
- **`spec/`** — Language specification split across numbered sections (01–11), describing syntax, semantics, and type system.
