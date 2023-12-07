Array

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
	Array: <T>{
		elements ...T
		push: { e T 
			elements << e
		}
		<< : push
	}
	... 
strings: ['hola']
strings << 'adios'
```
