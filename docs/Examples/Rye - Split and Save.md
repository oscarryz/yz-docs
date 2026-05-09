#example

```lisp
; Rye
split-and-save: fn { apples people } {
  apples / people |^check { 101 "Can't split to no people" }
  |to-string |write* %apples.txt
}

split-and-save 10 5 ; returns 2 and saves it to a file
split-and-save 10 0 ; returns a nested failure that turns into error
; Error(101): Can't split to no people
;	Error: Can't divide by Zero.

; If we want to 1 in case of any error because we
; need it for further processing we can for example use fix
; function

split-and-save 10 0 |fix { 1 } ; returns 1
```


```js
// Yz
split_and_save: {
   apples Int
   people Int

   apples.checked_div(people).and_then({ result Int
     fs.write("apples.txt", "${result}")
     Result.Ok(result)
   }).or(
      Result.Err("Can't split to no people")
   )
}
split_and_save(10, 5) // returns Ok(2) and saves it to a file
split_and_save(10, 0) // returns Err("Can't split to no people")
// To handle it we can use `Result.or_else`
split_and_save(10, 0).or_else({ Result.Ok(1) })

// If we want to preserve the original error we can create a new one
// with handle
10.checked_div(0).map_err({ e Err
  Result.Err("Custom message. Cause: ${e}")
})
// Would print(Error("Custom message. Cause: Error(Divison by zero))")

```

```js
// Proposal for an `if` helper boc (exploration)
if: {
    cond Bool
    then #(V)
    else: #(V)
    cond ? { then() }, { else() }
}
```