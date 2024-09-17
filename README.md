This repository contains design notes for the Yz programming language, like examples, features and design questions. 


> [!WARNING]
> There's no [compiler for the Yz programming language yet](https://github.com/oscarryz/yzc), making all the content here vaporware. 

# Quick overview

```javascript
// Factorial in Yz
factorial: { n Int
  n > 0 ? { n * factorial n - 1 }
          { 1 }
}
print "`factorial 5`"  // prints 120
```
Yz is a programming language that explores the possibility to simplify concurrency, data, objects, methods, functions, closures, classes under a single artifact: a code of block.

## Blocks
A code of block is a series of expressions between `{` `}`, their variables can be accessed from outside the block (behaving like objects) and they can be executed concurrently (behaving like methods/functions/closures).


```javascript
language: {
  name: 'Yz'         //  a variable `name` with the value `'Yz'` and infers the type String
  description String //  variable `description` of type `String`
}
main: {
  print(language.name) // prints: Yz
  language.description = 'A programming language' // sets the value for `description`
  print(language())    // Executes the block and prints `A programming language`
}
```

## Executing blocks

A block of code is executed by placing `()` after the block name, they are executed concurrently and synchronized at the end of the enclosing block of code.

The parameters are any of the variables defined in the block (ordered top down) and the return values are any of the variables or expressions (starting from down to top reversed). If the last value has a name, that can be used, if it was an expression an index could be used instead.

Here's an example: 
```javascript
// Fetches customer and product using customer_id and product_id parameters
// This is concurrent by default
retrieve_order: {
  // "parameters"
  cid Int
  pid Int
  // The following three blocks invocation run concurrently
  print('Retrieving information')
  fetch_customer(cid)
  fetch_product(pid)
  // The blocks finish excecution when all the blocks above had finished.
}

// Executing the block
// Defines two variables: `customer` and `product` using the last 2 computed values in the block
customer product: retrieve_order(89 012)
```

## Creating instances of a block

A new type can be defined if the identifier stars with uppercase. Use `()` to create an instance.

Example: 

```javascript
// Defines a new type `Point` with integer variables `x` and `y` 
Point : {
    x Int
    y Int
}
p1 Point =  Point(x:10, y:20) // p1 variable declared and initialized
p2: Point(40, 40) // p2 type `Point` inferred.

// An anonymous block can still structurally match the new type
p3 = {x: 10, y: 100} 
p1 = p3 // Can be assigned because it is a block with a `x` and a `y` of type Int
```

Variables can be accessed from outside the block and blocks can be variables too, thus a nested block can be used as "methods". 

E.g. The following defines the `to_string` block (method) that access the variables `x`  and `y` from the outer scope

```javascript
Point : {
    x Int
    y Int
    to_string: {
      "`x`,`y`" // $(expr) for string interpolation
   }
}
p: Point(0 0)
print(p.to_string()) //prints `0,0`
```

If a "method" is a non-word name ( e.g. `+`, `>`, `<`  etc., it can be executed without the `.` notation, this is a convenience to make the code look more traditional.

```javascript
Point : {
   x Int
   y Int

   // `+` is a non-word name block/method
   +: { 
       other Point
       Point( x + other.x 
              y + other.y)
    }
}
p1: Point(1 2) + Point(3 4) // invoking `+` without `.` results in a new Point{x: 4 y:6} 
// same as 
Point(1 2).+(Point(3 4))
```

### What about `this` or `self`?

There's no `this` or `self` variable but if needed it can be defined just like any other variable but it has to be assigned elsewhere (most likely in a factory method)

```javascript
Person : {
    name String
    self Person
    introduce_yourself: {
        'My name is `self.name`'
    }
}
alice:{
    p: Person('Alice')
    p.self = p
    p
}
print alice().introduce_yourself() // My name is Alice
```


## Expression blocks

It is also possible to have blocks with no variables and only expressions.

```javascript
one_two: { 1 ; 2 }
```
They behave the same, except they don't take parameters when executed, and naturally their values cannot be accessed by name, only by index.

```javascript
one_two: { 1; 2 }
a b: one_two() // a:1 b:2

// Desugared version
one_two #(Int, Int) = { 1; 2 } // `one_two` is a block of product types `Int` `Int` initialized withbthe block `{1 2}`
a Int 
b Int
one_two()
a = one_two.0 // 1 is the first computed value
b = one_two.1 // 2 is the second computed value
```

## The block type

If we want to use a block as a parameter we have to declare its type. The syntax for a block type is `#()`  and can contain other types which can be optionally be named as variables. If not named they are expected to be expressions.

The following declares a variable `a_two` of type block, that contains a variable `a` of type `Int` and a expression of type `Int` 

```javascript
a_two #(a Int, Int ) = {
  a: 1
  2
}
```

`swap` is a block with two variables of type int
```javascript
swap #(a Int, b Int) = {
    a Int
    b Int 
    b
    a
}
one two: swap(2, 1)
```

Most of the times you don't need to specify the block type as it can be inferred and/or a uppercase type can be used instead. It could be beneficial for instance to keep variables "private"(ish)

For instance if the block type defines a string `{name String}`  only that variable can be seen, all the variables inside are not accessible: 

```javascript
person #(name String) = {
    name: 'Bob'
    age Int // not accessible from outside as the block only declared `name String`
}
person.name // Bob
person.age  // compilation error: no variable named age
```


## Structural typing

Yz uses structural typing to know if a block can be used as a parameter and/or where other block is defined. The match of the structure includes the variable names.

```js
print_it: {
    thing #(Int, Int ) // a `thing` is a block that returns two integers
    a b :thing() // execute it
    print("`a` and `b`")
}
// All of the following will structurally match because they're things that return two ints

time: {
   hour: 12
   minute: 24
}  
Point : {
   x Int
   y Int
}
print_it(time) // a regular named block
print_it(Point(1 2)) // an instance of a `Point` type
print_it({ 4 2 }) // and expression block
```

### Special characters

The following cannot be used as identifiers or  part of identifier
- `{`,`}` used to create blocks.  
- `#` Used for block signature
- `(`,`)` used to execute blocks.  
- `[`,`]` used to declare/access arrays/dictionaries.   
-  ``"``, ``'`` and \` used to declare strings.  
- `:` used for type inference.  
- `;` used to separate expressions in the same line. 
- `,` not used, but reserved for clarity.  


# Other things
There are many other things not covered in this overview which is already too big, these topics are memory management, control structures (or the lack thereof), reusing code, generics, error handling, among others.
