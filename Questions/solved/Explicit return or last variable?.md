Should we return explicitly? 



    f: {
        1
    }
    a: f() // a is 1 

 Or should we need a variable and look for that variable (maybe when calling `()` will be "auto" scanned)

     f: {
         r: 1 
     }
     a: f() // a is 1 because  r was the last variable
     g: {
         1
     }
     b: g() // nothing

Alternative we can get all the value from the method either named or unnamed

    f: {
        1
        2
        3
    }
    a,b,c = f() // a=1, b=2, c=3
    max: { a Integer; b Integer;
        a > b ? { a }, { b }
    }
    /*
    max(a,b) {
        a > b ? a : b
    }
    */
    
    m: max(1,2) // last expresion is `?{{},{}}`

    Boolean: {
        ? : <T>{
            ifTrue {T}
            ifFlase {T}
        }
    }
    true: Boolean{}

Aug 16 2022

Will return last variable(s) from len - n

```js

foo: {
    a: 1
    b: 2
    c: 3
}
x,y,z = foo() // x: 1, y:2, z:3
```

And will use `return`  for non-local returns

```js
check: { value Int
    value == 0 ? {
        return  // will exit the `check` block
    }
    do_more_stuff()
}
other: {
    x: {
        return 1
    }
    x()
    0
}
other() // 1
```
#answered
