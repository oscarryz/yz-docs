How do we force an object to have a default value?
How would this work with the parameters? 


```
point: {x Int y Int}
```


We could force everything to have a value, even if is a parameter, for that would be its default value. 
For things that might not have a value we can use `Optional`

```js
// default value
point: {x Int = 0; y Int = 0} 
x1 y1 = point() // x1: 0 y1: 0
// no value
Person: {
    name String = "" // default name is empty
    middle_name optional.None // generic None of String
}
```


We could also allow don't have a default value, but if attempted to use the compiler would check and throw a compilation error. 


```js
say_hi: {message String print message } // no default value for message

say_hi() // compilation error, `message` requires a value
say_hi('Hello, world!') // there's a value
```

Same would happen with classes, generics, methods etc. 
When using named parameters, you can use up to the one that has no deefault

```js
foo: {
	a Int = 0 
	b String = ''
	c Bool
}
f() // compilation error, `c` requires a value
f(c:false) // ok
f(b: 'Hi') // compilation error, `c` requires a value
f(c: true  a: 0 b: 'GG') // all good 


```

#answered  Can ommit default value but if attempted to use, compilation error will arise.
