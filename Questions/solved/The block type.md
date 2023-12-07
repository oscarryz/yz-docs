In the [Block type](../../Features/Block%20type.md)  we discuss how the block type can be defined, but later there was a problem with the instantiation because declaring a variable of type block `a {}` looks exactly the same as instantiating a block named `a{}`  (using the space or the uppercase would be too error prone). ~~As alternative we can use `::`  for types and plain `{}` for instantiation~~. This also frees up the `{}` to instantiate any block. 

Example, creating an instance of an existing block with `{}`

```javascript
// Example, creating an instance of an existing block with `{}`
// Defining `say_hi`
say_hi: {
    to String
    print 'Hi {to}'
}
say_hi 'World' // Hi World
print say_hi.to // World // because it's the same instance 

// ERR this is not possible 
greet_me: say_hi{to: 'Me'} // creates a copy of `say_hi` with the value `to:` as me 
greet_me() // 'Hi Me'

// What's the type of `say_hi` ? 
say_hi { to String String } // a block with a variable `to` of type `String` and a value String (the result of print)
```


Problem: the type alias vs block type vs instantiation, vs type of a variable? 

#answered 
Answers: 
- There's no type alias
- Block type use ~~`::{}`~~ `{}` but can only contain types
- Instantiation use ~~`{}`~~ `()`  and can contain code
- Type of an instantiated thing is the same as the thing type 
    - e.g 
        - ~~`a :{}; a1:a{}` then `a1` type is `{}`  ~~
    - `Be{}; b1:Be()` then `b1` type is `Be` which is compatible with `{}`
- If the thing type is upper case that's the type of the instantiated thing


See also: 
[Generics without <>](Generics%20without%20<>.md)
[Variables as data types](Variables%20as%20data%20types.md)
[Type Alias (wip)](../../Features/Type%20Alias%20(wip).md)
[Creating instances](Creating%20instances.md)
