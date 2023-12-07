Should we make an special case for arrays and dictionaries for declaration and initialization? 

Usually the decl + init is like the following: 
```javascript
foo String = 'Hola'
foo: 'Hola'

```

Arrays however and specificially asociative arrays can't have empty value if no datatype is specified


```js
map: ["":0] // this can't be the empty as it has 1 value
map.len() // 1
map: [String]Int // this is the empty dictionary  which is the same as
map [String]Int   

map: ['foo':1, 'bar':2] // map [String] Int

```

Similar for Arrays

```js
array: [] // array of what? 
arrray:[]Int // ok

```
#answered in [Array](../../Features/Array.md) and [Associative arrays](../../Features/Associative%20arrays.md)
