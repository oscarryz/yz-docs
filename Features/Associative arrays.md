

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