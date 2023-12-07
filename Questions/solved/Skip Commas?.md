How feasible would it be to skip commas? 
```javascript
foo: {
  a Int
  s String
  b Boolean

 }
 foo(1, '', true)

 x y z = foo 1  ''  true

 x > 2 ? {
   print '{y}'
 } {
   print '{z}'
 }

//
code: [
  // letters
  'a': '.-' //etc
]
str: 'hello world'
out: ''
str.split().each { ch String
  ch.empty() ? {
    out = out + code[ch]
  }
  out = out + ' '

}
```

#answered : Yes, don't use coma to separate. just space. `a : [1 2 3]`
