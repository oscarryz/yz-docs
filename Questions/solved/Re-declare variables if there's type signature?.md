#answered Yes, it looks a bit redundant, but 90% of the time we won't be using type signature.

If a block has a type: 

```javascript
b { s String n Int }
// Or soon 
b ::{ s String n Int }
```

Does it make sense to re-declare the variables in the body? 

```javascript
// Explicit type signature
b :: { s String n Int } 
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
    n
    s
    n.times {
        s = s ++ s
    }
}
// b type would be
b :: { 
    n :: { times :: {}}
    s :: { ++    :: {s s}} // smh 
}
```