This repository contains design notes for the Yz programming language, like examples, features and design questions. 


> [!WARNING]
> There's no [compiler for the Yz programming language yet](https://github.com/oscarryz/yzc), making all the content here vaporware. 

# Quick overview

```javascript
// Factorial in Yz
factorial: { n Int
  n > 0 ? { n * factorial(n - 1) }
          { 1 }
}
print("`factorial(5)`")  // prints 120
```
Yz is a programming language that explores the possibility to simplify concurrency, data, objects, methods, functions, closures, classes under a single artifact: a code of block.

## Blocks
A code of block is a series of expressions between `{` `}`, their variables can be accessed from outside the block (behaving like objects) and they can be executed concurrently (behaving like methods/functions/closures).

To execute a block, use `()` as you would with a function

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

When invoking a block parameters will be assigned to the block variables starting from the top.
The return are the last expressions evaluated starting from the bottom - n e.g. 

```js
sum: {
   a : 0
   b : 0
   a + b
}
c: sum(1,2) // assing 1 to `a` and 2 to `b`.
// c gets the last expression `a + b`

x, y, z : sum(1, 2)
// x, y, z are the  last n - 3 expressions thus:
// x : 1 (the variable a)
// y : 2 ( the variable b)
// z : 3 (the expresion a + b) 

// Also, the `()` invocation can name the variables, behaving like named parameters
sum(b:10, a:20) // 20 + 10 
```

## Creating instances of a block

If the identifier starts with uppercase then it defines a new type. When you execute an instance is created. 
Example: 

```javascript
// Defines a new type `Point` with integer variables `x` and `y` 
Point : {
    x Int
    y Int
}
p1 Point =  Point(x:10, y:20) // p1 variable declared and initialized an a new instace created
p2: Point(40, 40) // p2 type `Point` inferred.
```

Variables can be accessed from outside the block and blocks can be variables too, thus a nested block can be used as "methods". 

E.g. The following defines the `to_string` block (method) that access the variables `x`  and `y` from the outer scope

```javascript
Point : {
    ...
    to_string: {
      "`x`,`y`" // `expr` for string interpolation
   }
}
p: Point(0, 0)
print(p.to_string()) //prints `0,0`
```

### Non-word names
If a "method" is a non-word name ( e.g. `+`, `>`, `<`  etc.), it can be executed without the `.` notation, this is a convenience to make the code look more traditional.

```javascript
Point : {
   ...
   // `+` is a non-word name block/method
   +: { 
       other Point
       Point( x + other.x, 
              y + other.y)
    }
}
p1: Point(1, 2) + Point(3, 4) // invoking `+` without `.`
// same as 
Point(1, 2).+(Point(3, 4))
```

## Expression blocks

It is also possible to have blocks with no variables and only expressions.

```javascript
one_two: { 1, 2 }
```
They behave the same as the other blocks, except they don't take parameters when executed, and naturally their values cannot be accessed by name, only by index.

```javascript
one_two: { 1, 2 }
a, b: one_two() // a:1 b:2
// or 
one_two()
a : one_two.0 // 1 is the first computed value
b : one_two.1 // 2 is the second computed value
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


Most of the times you don't need to specify the block type as it can be inferred and/or a uppercase type can be used instead. It could be beneficial for instance to keep variables "private"(ish)

For instance if the block type defines a string `{name String}`  only that variable can be seen, all the variables inside are not accessible: 

```javascript
person #(name String) = {
    name: 'Bob'
    age Int // not accessible from outside as the block only declared `name String`
}
person.name // Bob
person.age  // compilation error: no variable named `age`
```

If the type uses a single uppercase letter, then is a generic:

```js
Box #(T) = {
   data T
}
int_box : Box(1)
string_box : Box("Hi")
```

A generic without variable can be used to "instantiate" the generic type:

```js
Node: {
   T 
   data T
   left Node(T) // left hast to be of the same type T
   right Node(T)
}
root : Node(String) // assigns the type String to the parameter type T
root.left = Node(data:"a") // asserted vs the type T
root.right = Node(data: -1 ) // compilation error
``` 


## Structural typing

Yz uses structural typing to know if a block can be assigned to a variable. The match of the structure includes the variable names.
```js
thing #(Int, Int) // `thing` is a variable of type block with two integers
// Matches because `time` has two integers
thing = time: {
   hour: 12
   minute: 24
}
Point: {
  x Int
  y Int
}
// Can be assigned because `Point` has two Int
thing = Point( 1 , 2 )

// Can be assigned becaue the expresion block has two Ints
thing = { 3, 4 }
```

## Concurrency (or in Yz the functions color are ... purpple)<sup>1</sup>
 
Every block executes concurrently and synchornizes at the end of the parent block. 
To wait for a block to finish its execution, assign the return value to a variable (or put it inside another block and wait until that ends)
```js
main: {
   fetch( "Order123ABC" )
   print("Fetching order...")
   fetch_order.order // This is the end of the `main` block so it waits until both `fetch_order` and `print` finish executing. 
   // fetch.order has a value at this point
}

```
<sup>1: https://journal.stuffwithstuff.com/2015/02/01/what-color-is-your-function/</sup>
