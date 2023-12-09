```JavaScript
input String
digits:'01234567890'.split() 

// result
n : 0

// first and last characters
f String
l String

/*
  state machine signature.

 `state` is a block that has a
 `check` variable of type
  block that takes a `String`
*/
state {check{s String}}

/* 
  implements `state` structurally
*/
first : {
  check: {
     s String
     if digits.contains(s) {
        f = s
        state = second
     }
  }
}

/* 
  implements `state` structurally
*/
second: {
   check: {
     s String 
     when[
      {digits.contains(s):{ ​
         l = s
      }
      {s == '\n'}:{
        n = n + int.parse(f ++ l)
        state = first 
      }]
   }
}

state = first

input.each(state.check)
print("Solution: $(n)")
```
    
