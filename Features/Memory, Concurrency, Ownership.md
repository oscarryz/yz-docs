## Memory, Concurrency, Ownership
Variables declared in a block are owned by that block, they cannot be modified outside of the block directly but instead Yz  internally channels the modification to the owning block. 

From the end user perspective it looks like the value is being modified but internally the owning block queues or handles the update so there is only one writer.

If multiple external bocs try to update the same variable, the owner boc decides (internally) which one goes first, ideally in the order they show up in the code. 

e.g. 

```javascript
owner: {
    data Int
}
other: {
   thing : owner // thing has a reference to owner 
   thing.data = 1 // the block `{ data Int }` that is referenced by
   // both `owner` and `thing` is the one that modifies data and not 
   // "referenes" `owner` nor `other`
}
```

This allows to always have a single writer for a variable, the writer will be always the declaring block(`owner` in the example above).


