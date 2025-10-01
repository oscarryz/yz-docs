#example


```JavaScript
input String
digits:'0123456789'.split()
is_digit?: {
  s String
  digits.contains(s)
}

// result
n : 0

// first and last characters
f String
l String

/*
  A block that returns 
  a block that takes 
  a `String`
*/
state {{String}}

// "Implements" `state`
first : {
  {
    s String
    if is_digit?(s) {
      f = s
      state = second
    }
  }
}

// "Implements" `state`
second: {
  {
    s String 
    if is_digit?(s) { 
        l = s
    }
    if s == '\n' {
        v : int.parse(f ++ l)
            .or{ 0 }
        n = n + v
        state = first
    }
  }
}

// Initialize state
state = first

// For each character in 
// Input, execute the 
// state machine 
_ : input.each(state())
print("Solution: `n`")
```