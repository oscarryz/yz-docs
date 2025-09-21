#answered Yes, it looks a bit redundant, but has its spot. 

Answer in : [[Signatures + Literals duplication]]

If a block has a type: 

```javascript
b #(s String, n Int)
```

Does it make sense to re-declare the variables in the body? 

```javascript
// Explicit type signature
b #( s String, n Int )
// initialization 
b = {
    n.times {
        s = s ++ s 
    }
}
// But this wouldn't be possible with type inference because we don't know what type is n or s only they have the methods `times` and `++`
// Compilation error: "No variable `n` found"  and "No variable `s` found"
b:{
    n.times {
        s = s ++ s 
    }
}
// Here we would declare them explicitly 
b: {
    n Int
    s String
    n.times {
        s = s ++ s
    }
}
```

 We could also use generics and constraint on their usage, see [Generics without <>](Generics%20without%20<>.md) 
 

```javascript
// We could also use generics and constraint on their usage: 
b: {
    n N
    s S
    n.times {
        s = s ++ s
    }
}
// b type would be
b # (
    n # (times #())
    s #(++ #(s S)) // smh 
)
```

What if during declaration and initialization we use another symbol? 

```js
b #(s String,n Int) => {
  n.times({
    s = s.++(s)
  })
}
```

Or no symbol at all (this is the accepted)
```js
b #(s String,n Int) {
  n.times({
    s = s.++(s)
  })
}
```

🤔