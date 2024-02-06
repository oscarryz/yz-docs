In [Inko - Generic Data types](Inko%20-%20Generic%20Data%20types.md) the bound is determined by the use of the variable


e.g. 

```javascript
f: { v }
1 + f('a') // would fail because Int.+() expects a number  which is weird .. 
```

Do we need that? Can we get away without it?

#answered yes that's how it will work
#challenged Makes it hard to know [[How to enforce data types in generics]]  , we will use type as a parameter without name
