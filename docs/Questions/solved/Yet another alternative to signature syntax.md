Instead of plain `()` we could use `#()` which makes it easier to differentiate with function calls. 

Consider: 

```js
// New syntax
f #( T, #(), #() )
// vs
f ( T, (), () )

f = // something

f ( "Hello", { print("world") }, { 1 } ) 

```


#answered Using `#` to denote "block" type, can be used in uppercase types too

```js
Person #(name String, greet #(String)) = {
  name String
  greet #(String) = { 
      "Hi, my name is `name`"
  }
}
// is the same as
Person :{
  name String
  greet : { 
      "Hi, my name is `name`"
  }
}
```
