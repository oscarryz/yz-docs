Instead of plain `()` we could use `#()` which makes it easier to differentiate with function calls. 

Consider: 

```js
// New syntax
f #( T, #(), #() )
// vs
f ( T, (), () )
f = // something
f ( "Hello" { print("world") } { 1 } ) 
```
