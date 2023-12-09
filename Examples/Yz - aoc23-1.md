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

  A block that has a block 
  that takes a `String`
*/
state {{String}}

/* 
  implements `state` structurally
*/
first : {
  {
    s String
    digits.contains(s) ? {
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
    when [
      { digits.contains(s) }:{
        l = s 
      }
      {s == '\n'}:{
        v : int.parse(f ++ l)
            .or{ 0 }
        n = n + v
        state = first 
      }
    ]
  }
}

// initial state
state = first

input.each(state())
print("Solution: $(n)")
```
    
