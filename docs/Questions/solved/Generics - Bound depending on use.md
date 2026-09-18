---
type: solved
generated: { by: "oscarryz", at: 2023-12-06T20:13:02-06:00 }
---
#solved in [Generics - Type Parameters](docs/Features/Generics%20-%20Type%20Parameters.md)


In [Inko - Generic Data types](Inko%20-%20Generic%20Data%20types.md) the bound is determined by the use of the variable


e.g. 

```javascript
f: { v }
1 + f('a') // would fail because Int.+() expects a number  which is weird .. 
```

Do we need that? Can we get away without it?

#answered yes that's how it will work
#challenged Makes it hard to know [[How to enforce data types in generics]]  , we will use type as a parameter without name
