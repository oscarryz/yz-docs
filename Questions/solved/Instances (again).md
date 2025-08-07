It seems using `{}` for object literal would be overloading the `{}` too much.

Alternative is to use `()`  for instantiation

```js
p: Point(1 2) // instead of p: Point{1 2}
```

That would also free the need to use `Foo:{}` to define new types to differentiate it from `Foo{}` and the need it to have a `type`
```js
// new possible syntax to define types
Point {
   x Int
   y Int
}
// if needed we can define a keyword, wouldn't be terrible
type Point { 
   x Int 
   y Int
}
```

However how can we match anonymous structures? 
e.g.

```js
p Point
p = {x:1 y:2} // this should be valid because they have he same shape.
```

The arguments will follow the rules for method invocation, they can provide only values or use names 
```js
Point (y: 2 x: 1 ) // x and y specified in opposite direction 
Point (y: 2 1 ) // probably not
```

What would happen with anonymous blocks? 

```js
Point { 
  x Int 
  y Int
  to_string : {
     'Point(x: $x y: $y)'
  }
}
x: {x:1 y:2}
point Point = x // invalid ... x doesn't have `to_string`
valid_point : Point(0 0)
x = valid_point // valid. `valid_point` satifies {x Int yInt}
```


[Instances](../../Features/Replaced%20features/Instances.md)

#answered Two ways, using the object literal and creating new types: 
```javascript
Foo {
	name String
}
f: Foo("I'm froot")
g: {
	name("I'm groot")
}
e { name String }
e = f // yes 
e = g // yes
```
