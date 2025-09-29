Now that [Instances (again)](solved/Instances%20(again).md)  specifies how to create instances and define types, the type won't necessarily have an extra type but it would structurally match to any other block wit the same types.


```js
// new type
Point{ x Int y Int}
print_point: {
	p Point
}
print_xy: {
	xy {Int Int} // something with two Int's
}
p : Point(1 1)
print_point(p) // valid
print_xy(p) // valid

```

The question here is if `{Int Int}` if a enough to know this is the type of the block or we need something like `:: {Int Int}`  there's probably ambiguity when the block type defines names

```js
print_thing: {
   Thing { x Int y Int }  // defines a new type
   thing {x Int y Int}    // declares a variable 
   // could it be
   thing :: {x Int y Int}
   // or
   thing :: x Int y Int
}

Point {
   x Int
   y Int
   to_string {String} // a block that returns a string
}
p Point = { 1 2 {'(1 2)'}} //COMPERR: not possible because the block is specifiing values
// but Point needs names too 
Pointy:{ 
  Int 
  Int 
  #(String)
}
p1 Pointy = { 78 89 {'78,89'}} // yes! Pointy only requires values 
p = Point() // default values 
p2 Pointy = p // Yes! p has all the types required by Pointy 
p2.0 // 78
p2.1 // 89
p2.3 // '78,89'
```


#answered Define block with `Name {}` , it doesn't have a type by itself, but is structural compatible with other structures with the same members.

#updated New types do use `:`  the block type (signature) used `()`
So updated and final version would be:

```js
Point:{ x Int y Int}
print_point: {
	p Point
}
print_xy: {
	xy (Int Int) // something with two Int's
	xy.0 // first element
	xy.1 // second element
}

print_xy_ish : {
	xy (x Int ; y Int ) // named variables
	xy.x
	yx.y
}
p : Point(1 1)
print_point(p) // valid
print_xy(p) // valid
print_xy_ish(p) // still valid

```