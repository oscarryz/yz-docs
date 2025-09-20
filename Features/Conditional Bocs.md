#feature
#pattern-matching 

A similar version of pattern matching
A boc can validate a condition or match a type using the keyword `when` and separating the condition and the action with an arrow `=>`

Example 
```js
factorial #(n Int) {
  match { 
    n == 0 => 1 
  }, { 
    n > 0 => factorial(n -1)
  }
}
```


