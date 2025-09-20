#feature 
Also known as dictionaries or hash maps
```js
// Type
[Key_Type : Value_Type] 

// e.g. 
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

// conditions is a dictionary that takes blocks that return Bool  as keys
// and blocks that return a generic A as values, in the case below
// they return Strings
conditions [ #(Bool) : #(A) ] = [
	{is_it_monday()} : {"Monday"},
	{is_it_tuesday()} : {"Tuesday"},
	{true} : {"Not Monday nor Tuesday"}
]
```