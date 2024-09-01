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

But a better solution would be to support fully parametric polimorphism 
so we can enumarate what kind of thing is T

```js
add : { 
  T : Int | Decimal // made up syntax
  a T 
  b T 
  a.plus(b)
}
add(1, 0.0)
```
This might open the posibilities to add union types 
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

