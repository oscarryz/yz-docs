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
a [Int] = []Int // Is an empty array
// short declr + init
a : []Int // empty array of ints 

// Generic
a [T] = [1, 2, 3]
a : []T
```

```javascript
// Type
array [Int]
// init
array = []Int

//short decl + init 
array: []Int // identical to the init 
// or 
array: []T // no type specified yet until first usage
array << 1 // now type is []Int
```

Example

```javascript
a [String]// a is an array of string
a = []String
a << 'Hello' // or a.add 'Hello'
print(a[0]) // prints Hello

```
Also to consider

```javascript
a [String] // array of strings
// 28 agosto 2022. Maybe not
// Sometime in 2024. Yes!
```

Literal, elements separated by comma

```javascript
a: ['Hello', 'World']
print(a[0]) // prints Hello
```

