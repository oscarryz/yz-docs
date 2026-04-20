#feature

To define a new type, we use upper case names, followed by the type body, the following would define a new type `Person`

```js
Person : {
   name String
   to_str : {
	   "My name is `name`"
   }
}
```

A type can also have a type signature, just like a regular boc

```js
Person #(name String, to_str #(String)) 
```

## Summary

- Only allowed in uppercase starting identifiers
- Use `:` to infer a new type

#answered 
### Solves: 
[The block type](solved/The%20block%20type.md)
[Type Alias](Features/Type%20Alias.md)
[Instances](Replaced%20features/Instances.md)

### Related:
[Bocs](Bocs.md)
[Block type](../Features/Block%20type.md)

