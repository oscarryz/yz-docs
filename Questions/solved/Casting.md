
*answer*: No casting. need to create a new value with a constructor / factory method.

```javascript

   a : 1 // Int
   b Long
   // b = a error: invalid assigment
   // correct: 
   b = ints.to_long(a) 
   ints: {
       to_long: { n Int; l long
           // magic
       }
   }
```   

 #answered
