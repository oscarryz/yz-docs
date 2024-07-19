Also known as dictionaries or hash maps
```js

// Jul 19 2024
// Type
[Key_Type : Value_Type] 

// e.g. 
// declaration
d [String:Int] 
// initialization 
d = [ "one": 1 "two": 2]
// decl + init 
e [String:Int] = ["one":1 "two":2]

// short decl + init
f : ["one":1 "two":2 ]

// empty 
g2 [String:Int] = [:]
// short decl + init empty
g1 : [String]Int

// generic
g3 [K:V] = [String]Int 
g4 [K:V]
g4["hello":1]

// conditions is a dictionary that takes blocks that return Bool  as keys
// and blocks that return a generic A as values, in the case below
// they return Strings
conditions [ (Bool) : (A) ] = [
	{is_it_monday()} : {"Monday"}
	{is_it_tuesday()} : {"Tuesday"}
	{true} : {"Not Monday nor Tuesday"}
]
```

 
 _Historic analysis below, might be outdated_
 
Type
```javascript
// Dictionary type:  
[Type : Type]
// literal 
[ key1: value1 key2: value2]

// Empty
// generic
[:]
// typed (same as type)
[Type:Type]

// Access
// Read
[key]
// Write
[key: new_value]

// Examples: 
dict1: [String: Int] // empty and ready to use 
dict2: ['name': 4 'last_name': 9]
dict3: [:] // empty and genric 
l: dict2["name"] // `l` is 4
dict2["last_name": 80] // dict2["last_name"] now returns 80
```

You can ignore everything below here... 

Proposal 1: `[key]value`

Declaration: 
```js
dict [String] Int
```
Empty literal is also generic 
```js
empty: [:]
```
Declaration and initialization
```js
dict [String] Int = [ "Hello": 5 "world!": 6 ]
```

Definition (literal, and type inference) 
```js
dict: [ "Hello":5  "world!" : 6 ] /// [String]Int inferred
```
Access: 

```js
n: dict["hello"] // read 
dict["hello"] = 6 // write
```

When using empty `[:]` the types are generic and bound on first usage: 


```js
map: [:]
map["name"] = {} // `map` is a dictionary with `String` key and empty block ( `{}` ) value
map["data"] = 1 // compilation error, dictionary value must be block `{}` 
```

Proposal 2 `[key:value]`

Declaration: 
```js
dict [String:Int]
```
Empty literal is also generic 
```js
empty: [:]
```
Declaration and initialization
```js
dict [String:Int]  = [ "Hello": 5 "world!": 6 ]
```

Definition (literal, and type inference) 
```js
dict: [ "Hello":5  "world!" : 6 ] /// [String]Int inferred
```
Access: 

```js
n: dict["hello"] // read 
dict["hello":6] // write
```

When using empty `[:]` the types are generic and bound on first usage: 


```js
map: [:]
map["name"] = {} // `map` is a dictionary with `String` key and empty block ( `{}` ) value
```

Proposal 3 `{key=value}` 


 

```js
dict [String] Int // array String index int as value
dict = ['key':1  'a': 2] // initialization and literal
empty: [String] Int
// or 
empty: [:] // empty dict, will bind oh the first usage
empty['key'] = 0 // now is [String]Int

dict['key'] // access to 'key' returs 1
dict['a'] = 4 // set new value for key 'a'

other_dict: [:] // no longer error
```

 Example: 

 ```js
map [String] Integer
// how to initialize it is tbd because [] is empty array
map: [:] // fails to declare the type
map: #[]  // same as this
map [String] Int = #[] 
map['count':10]
//or
map['count'] = 10
map['count'] // returns 10
print(map['count']) // prints 10

```

Literal
```js
map: ['count': 0]
print(map['count']) // prints 10
```

```js
a [Int]; a = [1]
a []Int; a = [1]
m [Int]String; m=[1]='hello' // no but, `m = [:]; m[1] = 'hello' // yes`
m [Int:String]; m=[1:'hello'] // yes

a2 []Int; a2 = [1,2,3,4,5]
m2 [Int]String; m2 = [1:'Uno',2:'Dos',3:'Tres',4:'Cuatro']
m2[5:'Cinco',6:'Seis', 7 : 'Siete', 8 : 'Ocho', 9 : 'Nueve']


...
m3 [Int]String
1.to 10, {i Int
  m3[i:'Hola']
} 
print m3  // [1:'Hola', 2:'Hola'...]
```

 ## Open question 

#open-question 

- Empty hash literal, that doesn't collide with empty array ? `[]` or is empty block `{}` 
```javascript
empty_dict: [String]Int 
```
    - Maybe `#[]`  (it's a bit weird)
    - Maybe `[:]`  both would need to dict to be declared beforehand
```js
loop_up [String]Employee = #[]    
```

- How to check if key is present without using `null`'s ?
    - A: will be initialized to the block's defualt value (``''``, `false`, `0` , etc)

~~Option1: Add a "magic" built-in method `Yz.has_key(dict, key)`~~

```js
//map [Int] String = #[]

//has_key(map, 1)? {
//   return map[1] // return map[1]
//} 

```

Strong contender: 

```javascript
map: [String]Int
```
