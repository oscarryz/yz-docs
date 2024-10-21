# Block type 

Example of a block that takes nothing and returns nothing

```js
#() // empty block type
#() // empty block 
#() = {} // initialized as empty block 
```

Example block that returns an `Int`

```javascript
#(Int) // a block that takes (or returns) an Int 
#(Int)={1} // initialized with a block that returns 1
#(Int)={1}() // and invoked
```

Example of a block that takes or has a variable `v`

```javascript
#(v Int) // Block with a variable v 
#(v Int)={v Int} // initialized with a variable int 
#(v Int)={v Int}(2) // invoked with param 2
```

if the `=` is omitted the variable can be used: 

```js
#(v Int) {
  print("`v`")
}
```

## Type alias

[Type Alias](Features/Type%20Alias.md)

