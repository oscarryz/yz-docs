# Use space for non words method only? 


Currently (Feb 11, 2022) I'm thinking of letting method invocation use either `.` or empty space: 

    object method(1,2,3)

As long as it has parameters. 

    object method() // would still require . 

Also I'm thiking of making the parenthesis optional

    object method 1,2,3

That would allow for neat interactions with methods with non-word characters for instance "add" `+`

    object + one
    list ++ other

But it might might be easier to the eye to see the method with a `.`  always or most of the time 

e.g 
```javascript
while({supplier.is_empty() == false}, {
    v: suplier.next()
    v.==(nil).?{loop.exit()}
    max = max.max(v)
})
//vs
while({supplier.is_empty() == false}, {
    v: suplier.next()
    v == nil ? {loop.exit()}
    max = max max(v)
})
// vs v0.1
max: {a Int; b Int; a > b ? {a},{b}}
while {supplier.is_empty() == false}, {
    v: supplier.next()
    v == nil ? { break }
    m = max(m,v)
}
```


Aug 16 2022
#answered
Will use `.`  for all word methods and space will be optional for non word

eg. 
```javascript
one + two 
one.plus two

One: {
    + : { other T }
    pus : { other T }
}
```

```javascript
foo: {
    ==: { other Int
       false
    }
    eq: { other Int 
        false
    }
    #: {
        'hash'
    }
    
}

x: 1
foo == x  // yes
foo.== x  // yes
foo.==(x) // yes
foo eq x  // compilation error 
foo.eq x  // yes
foo.eq(x) // yes

foo #    // compilation error
foo.#    // gets the method reference
foo.#()  // 'hash'

```
