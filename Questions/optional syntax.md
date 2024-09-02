In the examples I interchange between parenthesis and no parenthesis,betwee comma and no comma. 

Could it be that parenthesis, comma (and semicolon while at that) are optional? Would it be too much? 


```js
factorial : {
   n Int
   if n == 0 { 
      1 
   } { 
      n * factorial n -1 
   }
}
//
factorial : {
   n Int;
   if (n == 0 { 
      1; 
   }, { 
      n * factorial( n - 1)  
   });
}
```
