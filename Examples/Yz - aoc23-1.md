```JavaScript
input String
digits:'01234567890'.split() 

// result
n : 0

// first and last characters
f String
l String

/*
  state machine "signature":

  A block that has a
 `check` block that takes a `String`
*/
state {{String}}

/* 
  implements `state` structurally
*/
first : {
  {
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
   {
     s String 
     when[
      { digits.contains(s) }:{ ​l = s }
      {s == '\n'}:{
        n = n + int.parse(f ++ l)
                .or{ 0 }
        state = first 
      }]
   }
}

// initial state
state = first

input.each(state())
print("Solution: $(n)")
```
    
