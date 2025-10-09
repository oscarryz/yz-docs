#feature
Array

```js
// Jul 19 2024
// Type
[Type]
// e.g.
// declaration
a [Int]
// initialization
a = [1, 2, 3]

// decl + init
a [Int] = [1, 2, 3]

// short declr + init
a : [1, 2, 3]

// emtpy decl + init
a [Int] = [Int]() // Is an empty array
// short declr + init
a : [Int]() // empty array of ints 

// Generic
a [T] = [1, 2, 3]
a : [T]()
```

```javascript
// Type
array [Int]
// init
array = [Int]()

//short decl + init 
array: [Int]() // identical to the init 
// or 
array: [T]() // no type specified yet until first usage
array << 1 // now type is []Int
```

Example

```javascript
a [String]// a is declared as an array of strings
a = [String]() // a initialized as an empty array of strings
a << 'Hello' // or a.add 'Hello' // with non-word medthod invocation ()
print(a[0]) // access element 0 of the array

```
Also to consider

```javascript
a [String] // array of strings
```

Literal: elements are separated by comma

```javascript
a: ['Hello', 'World']
a[0] = "Hola"
print(a[0]) // prints Hola
```

