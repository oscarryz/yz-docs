```javascript
; Rye
split-and-save: fn { apples people } { 
  apples / people |^check { 101 "Can't split to no people" }  
  |to-string |write* %apples.txt 
}

split-and-save 10 5 ; returns 5 and saves it to a file
split-and-save 10 0 ; returns a nested failure that turns into error
; Error(101): Can't split to no people 
;	Error: Can't divide by Zero.

; If we want to return 1 in case of any error because we
; need it for further processing we can for example use fix
; function

split-and-save 10 0 |fix { 1 } ; returns 1
```    


```javascript
split_and_save: {
   apples Int
   people Int
   r: apples / people
   int.is_error(r) ? {
       option.fail "Can't split to no people"
   } {
      bytes: fs.write r.to_str() 'apples.txt'
      option.ok bytes
   }
}
split_and_save(10 0).or_else { 1 }
```

```javascript
if: <T> {
    cond Bool
    true_block <T>{}
    false_block: <T>{}
    cond ? true_block false_block
}
```
