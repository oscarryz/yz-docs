---
type: rejected
generated: { by: "oscarryz", at: 2023-12-06T20:13:02-06:00 }
---
#rejected 


#idea 

Allow extension methods where a method can be defined for an existing type

```
Suit :: Int
// how to add "hello" method to `Suit` ? 
```

Possible solution just re-define it

```
Suit: {
    hello: { "hi" }
}
```
But that's very confusing as won't be a way to know if it's a new type or an extension. 

#answered Not for the moment




