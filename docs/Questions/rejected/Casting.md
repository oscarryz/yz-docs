---
type: rejected
generated: { by: "oscarryz", at: 2023-12-06T20:13:02-06:00 }
---
#rejected 
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

 
