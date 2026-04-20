#answered in [Associative arrays](../../Features/Associative%20arrays.md) array (aka Dictionary, aka Map)

```javascript
dict [String] Int // declaration
dict [String:Int] // declaration this would make more sense 
dict []String:Int
empty: #String:Int[] // initialization ? 
empty: [String]Int []
empty: []String:Int



```
Akin to [Array initialization](Array%20initialization.md) if empty it needs to specify types... 

```javascript

map [String:Int]// map is a asociative array whose keys are `String and values are `Int`
map = [:] // if already declared this initializes it 
map [String:Int] = [:]

map2: <String:Int>[:]
array: <Int>[]

// Empty array
array: []Int
// Empty map
map  : [:]String:Int

// Declaration
array []Int
map   []String:Int
x:  ['name':'Oscar']
x['name'] // Oscar
x['name': 'reyes']
x['name']='reyes'


```

Strong proposal: both array and dictionaries only have declaration and that initialized them as empty.
e.g.

```javascript
array: []Int // `array` is an empty array of ints
array << 1 // adds `1` to array

map:   [String]Int // `map` is an empty dictionary whose key are Strings and values are objects
map['hola'] = 1 // `1` is associated with the key `hola`
```


If previously declared, can initialize as empty: 

```js
dict []Int:String
dict = []

empty_dict: []Int:String // needs type

```