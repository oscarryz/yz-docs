https://ucs.mlscript.dev/  > MLscript Guide

```js
0 // Integers
3.14159 // Decimal
true // a variable that happens to be Boolean
"this is a string"
{ true,  false} // "tuples"

```

```js
// Expressions can be bound to variables
x : 1
y : 2
z : x + y // 3
```

```js
// Functions can be defined with `{}`, parameters are variables
id : { x T ; x }
id("hello") // hello 
sum : { x Int; y Int ; x + y}
sum(40, 2) // 42
```

```js
// The types of functions are automatically inferred. But you can also specify the types explicitly. Type parameters are single upper case letter
id #(x A; A) = { 
  x A // #to-do we need to redefine here or not?
  x
}
```

```js
// "Classes" start with uppercase. For example, the following is a class definition for a simple `Point` class. It has two parameters `x` and `y` of type `Int`
Point: { x Int; y Int }
// variables are accesible for reading
origin : Point(0, 0)
origin.x + origin.y // 0 

// We can also define classes with type parameters.
Pair: {
   first A
   second B
}

// And the new type has a signature inferred as well, for explicity signature use #() 
Pair #(first A, second B) = {
  first A
  second B
}
```

```js
// Clases can be reused with `use`. Here is one of the definition of `Option` type in Yz, which is very common in provided examples.
Option: {
  T
}
Some: {
  use Option
}
None: {
  use Option
}
// But we can do without reusing code and provide "explicit named constructors" 

Option: {
  T
  Some(v T),
  None()
}
// Here the type Option can only be constructed with either `Some(T)` or `None`
x Option(String) // x is an Option of String
y Option(String) // y is an Option of String
// 
x = Option.Some("Hello") // use the Some constructor
y : Option.None()  // Use the None constructor
// check with x.Some or x.None that returns a boolean 
if x.Some , { 
  print("We've got a value `x.v`")
}
if x.None, {
  print("Trying to access would error `x.v`") // compilation error
}
// #to-do Finalize how the "explicit named constructors would work"


// The Option signature in this case is
Option #(T, Some(v T), None())

```

```js
// Similarly, we can define the `List` type
List: {
  T
  Cons(head T, tail List(T))
  Nil()
}
Cons : List.Cons
Nil : List.Nil
Cons(0, Cons(1), Cons(2, Nil()))
// To costruct list like other functional languages, we can define the concatenation opperator `++` by augmenting the `Int` type
Int: {
  ++ : {
    x T
    xs List(T)
    Cons(x, xs)
  }
}
0 ++ 1 ++ 2 ++ Nil()

```

```js
// We can define many functions that operate on list. For example, we can define a function that sum all elements in a list
sum: {
  xs List(T)
  if xs.Nil , {
     0 
  }, {
    // It wasn't Nil, which means it is Cons
    // thus the properties `head` and `tail` have values
    xs.head + sum(xs.tail)
  }
}
```

A bit of discussion on how to add "extension" methods: 

```js

// This doesn't work, it is Rustesque made up syntax
Int : { ... }
Hi: {
  sayHi #(String)
}
impl Hi for Int { 
  sayHi #(String) = { 
    "My name is `self`"
  }
}

one Hi = 1

// Our possible approach is to le the class be reopened
// which will be needed for file based source anyways 
Hi: {
  sayHi #(String)
}
Int : {
  sayHi #(String) = { 
    "My name is `self`"
  }
}
one Hi = 1

```