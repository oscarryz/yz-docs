

```JavaScript
input String
digits:'0123456789'.split()
digit?: {
  s String
  digits.contains(s)
}

// result
n : 0

// first and last characters
f String
l String

/*
  state machine "signature":

  A block that returns a block 
  that takes a `String`
*/
state {{String}}

// "Implements" `state`
first : {
  {
    s String
    if digit?(s) {
      f = s
      state = second
    }
  }
}

// "Implements" `state`
second: {
  {
    s String 
    when [
      {digit?(s)}:{ 
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

// For each character in input
// execute the state machine 
_ : input.each(state())
print("Solution: $(n)")
```