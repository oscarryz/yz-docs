#example
```js
// This is a single comment  
/*  
   This is a multiline comment
*/  
  
`  
   Counter is a block that has a variable count of type Int
   and a variable increment of type block that takes nothing and returns nothing
` 
Counter {  
   counter : 0
   increment : {   
      counter = counter + 1  
   }
}  
main: {  
    c: Counter()
    print("c.count = `c.count`")
    c.increment()   
    c.count > 10 ? {
          break   
    } {  
      continue
    }
    0
}
```
