Array

```js
// Jul 19 2024
// Type
[Element_Type]
// e.g.
// declaration
a [Int]
// initialization
a = [1 2 3]

// decl + init
a [Int] = [1 2 3]

// short declr + init
a : [1 2 3]

// emtpy
a [Int] = [] // 
// short declr + init
a : []Int // empty array of ints 

// Generic
a [T] = [1 2 3]
a : []T
a : [] (String,String)
a [(String,String)]
```

```javascript
// Type
array []Int
// init
array = []

//short decl + init 
array: [] Int // identical to the type
// or 
array: [] // no type specified yet until first usage
array << 1 // now type is []Int
```

Example

```javascript
a [] String // a is an array of string
// ~how to initilize it is tbd~ answer: [] or []String for short decl
a = []
a << 'Hello' // or a.add 'Hello'
print(a[0]) // prints Hello

```
Also to consider

```javascript
a [String] // array of strings
// 28 agosto 2022. Maybe not
```
	

Literal

```javascript
a: ['Hello' 'World']
print(a[0]) // prints Hello
```

Array block declaration
```javascript

	'special'
	Array: {
		T
		elements T
		push: { e T 
			elements << e
		}
		<< : push
	}
	... 
strings: ['hola']
strings << 'adios'
```
