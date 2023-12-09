```JavaScript
input String
digits:'01234567890'.split() 

n:0
a String
b String
/*
 `state` is a block that has a
 `check` variable of type block that takes a `String`

*/
state {check{s String}}

first : {
  check: {
     s String
     if digits.contains(s) {
        a = s
        state = second
     }
  }
}

second: {
   check: {
     s String 
     when[
      {digits.contains(s):{ ​
         b = s
      }
      {s == '\n'}:{
        n = n + int.parse(a ++ b)
        state = first 
      }]
   }
}

state = first

input.each(state.check)
print("Solution: $(n)")
```
    
