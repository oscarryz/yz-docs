In the current proposal most of the control structures are in the standard library, e.g. `if`, `when`, `while`, `some`, `none` etc. 

Using a block is an object and a function makes it almost impossible to do multithread because every concurrent invocation might trip into someone else's code. 

e.g. `if`  takes a boolean and a generic block 

```js
std: {
    if : {
        condition Bool          // 1
        action    {value}       // 2
        else      {value}       // 3
        condition ? action else // 4 invoke Bool.?(ifTrue ifFalse) method
    }
}
one: {
    if random.int() == 0 {
        print 'Hola'
    } {
        print 'Adois'
    }
}
two: {
    if false {
        print 'Never'
    }
}
one() 
two()
// Sequence of actions
/*
one.if.1
two.if.1
one.if.2 // hola
two.if.2 // never
one.if.3 // adois
two.if.3 // undefined
one.if.4 // execute random.int == 0
Because they are the same if two.if.2 would override one.if.2
*/
```

## Possible solution

### A. Execution creates a copy.

Make every execution a copy of the object. 
Pros: no stepping on each other
Cons: Can't keep state between executions

```js
person: {name:'Oscar'}
person.name // keep 
person() 
/*
creates a new instance of `person` (copies values)
executes the code, return value
returns name:'Oscar' -> 'Oscar'

*/ 
```

### B. have a `std.clone(block)` 

Call `std.clone` to create a copy and assign the copy to the new object
Pros: keep state separated
Cons: I don't want to have a special keyword `clone` or function

```javascript
person: {name: 'Oscar' }
juan: clone person
juan('Juan')
assert juan.name = 'Juan'
assert person.name == 'Oscar'
```

### C. Make them classes 

Make those method classes so they have to create new instances to operate 
Cons: The parameters are still linked to the original usage so they can still be changed 

```javascript
sdt: {
    If: {
        cond   Bool
        action {value}
        alt    {value}
        run: {
            cond.if_else action alt
        }
       
    }
    if: {
            cond   Bool
            action {value}
            alt    {value}
        body: {
            cond   Bool
            action {value}
            alt    {value}
            cond ? action alt
        }
        body() // creates a closure 
    }
}
```

### D. Make everything immutable 

See: [Immutability](Immutability.md)

### E. Create objects to handle the work

```javascript
condition: {
    boolean.true() // returns a new instance of true    
}
if condition() { 
    1
} {
    2
}
... 
True: Boolean {
    ? : { 
        tb {v}
        fb {v}
        tb()
    }
}
boolean: {
    true: {
        True{}
    }
}
while: {
    cond {Bool}
    body {v}
    w: While{}
    w.run cond body
}
While: {
    run: {
        cond {Bool}
        body {v}
        if cond() {
            body()
            run cond body
        }
    }
}
```


#answered 
Option E the utility functions ( if, when, while, etc ) will create a new instance of themselves ( `If`, `When`, `While`) which will set the variables individually so they are not shared across. It fits naturally with the language.  