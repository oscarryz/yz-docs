We treat the last variables declared as return values

```javascript
    foo: {
        bar: 1
        baz: 2
    }

    a, b = foo() // a -> 1, b -> 2

```
We should consider returning the expression too as value so the variable is not needed. 

```javascript
    foo: {
        1
        2
    }
    a, b = foo() // a = 1, b = 2
```
#answered
Yes, everything will be an expression except for variable declaration / initialization
But acutally there's nothing in the 'all' part, the rest are method invocations
#challenged... why couldn't they be expressions too `a: b: 1` //a and b are 1 

```javascript
a: 1
{a}() // invokes the block and returns a 

```

