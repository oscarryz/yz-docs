Currently generics are good for single value that don't require calling methods


```js
Box {
   data T 
}
b : Box(String)
unknown Box(T)
b.data = "Hola"
unknown = Box(1) // instantiatate Box with an int
```

But is not so good with calling methods 

```javascript
add : { 
    a T 
    b T 

    a.plus(b)
}
```

The first solution is to declare an `Adder` type 
```js
Adder #( plus(Adder), Adder) 
add : { 
 a Adder
 b Adder
 a.plus(b)
}
add(1, 0)
add(1, 0.0) // nope, they are different types
```

But a better solution would be to support ~~fully parametric polymorphism~~ type constraints
so we can enumerate what kind of thing is T

```js
add : { 
  T : Int | Decimal // made up syntax
  a T 
  b T 
  a.plus(b)
}
add(1, 0.0)
```
This might open the possibilities to add union types 
```js
Adder {
  Int
  Decimal
  plus #(Adder, Adder)
}
add : { 
   a Adder
   b Adder
   a.plus(b)
}
add(1, 1.0) // both are adders, and they keep their internal type

```

We might push even further and access type-specific methods

```js
Twice : {
   i Int
   twice: { i * i } 
}
Half : {
   d Decimal
   half: { d / 2 }
}

TwiceOrHalf: { 
   Twice
   Half
}
// made up syntax
do_something : {
  e T : Twice | Half 
  switch e.type() // {
	Twice: print("Two times: `e.twice()`)
	Half : println("Halved: `e.half()`)
  } 
}
// Yz'ish syntax
do_something : {
  e TwiceOrHalf 
  print(e.info().type().when([
	{Twice}: {'Two times: `e.twice()`' }, 
	{Half}: {'Halved: `e.twice()`' }
  ])) 
}
// We might need first keyword
do_something: {
  e TwiceOrHalf 
  switch(e.type) { 
	Twice: { print("Two times: `e.twice()`) } 
	Half : { println("Halved: `e.half()`) } 
  }
}
```

Another example is haskell type classes 

```haskell
sort :: Ord a => [a] -> [a]
```

Where `Ord` ([link](https://hackage.haskell.org/package/base-4.20.0.1/docs/Data-Ord.html))defines the operations `a` can do

We need a way to say: _This T can respond to the method "compare"_

```js
sort #(arr [T], arr[T])
```

Currently we would have to declare a type first, and then use it
```js
Ord #(
  T
  compare #(T, Int)
  < #(T, Bool)
  <= #(T, Bool)
  > #(T, Bool)
  >= #(T, Bool)
)
sort #(arr [Ord], [Ord])

sort([1 2 4 2 6 7]) // Ok, Int implements Ord

```

But we wouldn't be able to use something like `Ord | Debug`  or something like that... 

```js
"#[Derive: Ord, Debug]"
OrdOrDebug:{  }
```
But I guess that's fine.

Solution: Generic types can have any number of methods, but the compiler will verify if the arguments have the needed methods. 

```js
sort #(a [T], [T]) 
sort([4 2 3 1]) // ok
sort([{},{},{}]) // compile error, [#()] doesn't have a '>=' method.
```

#answered 
