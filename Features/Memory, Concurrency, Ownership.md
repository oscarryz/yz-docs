
Variables declared in a block are owned by that block, they cannot be modified outside of the block. Yz provides a mechanism to make it look like they can be modified from the outside, but everytime is the owning block who performs the action. 

e.g. 

```javascript
owner: {
    data Int
}
other: {
   thing: owner // thing has a reference to owner 
   thing.data = 1 // the block referenced by both owner and thing is the one that modifies data and not the block `other`
}
```

This allows to always have a single writer for a variable, the writer will be always the declaring block(`owner` in the example above).

