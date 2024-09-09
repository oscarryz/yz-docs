In the examples I interchange between using and not using parenthesis, commas, semicolons. 

Could it be that these punctuations symbols be optional? Could it be too much?


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
