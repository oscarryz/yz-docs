#feature
# Boc type 

A block of code `boc` has a type. Just like a number literal `1` has a type `Int` or `1.0` a type `Decimal` a boc literal has a boc type defined by `#(` optional variables or expression types  `)`

The variables of a boc type follow the same rules as regular variables: 

- Named with a type: 
  `#(a Int)`

- Assigned a default value: 
   `#(a Int = 1)`

- Short decl:
   `#(a : 1)`

- Other boc types
   `#(a #())`
- Generic:
   `#(a T)`

In adition to variables, the boc type can have pure expression types, these are always at the end of the list. They can be thought as input and output parameters:

- A boc that just returns a String 
  `#(String)`

- A boc that takes an integer and returns a string
   `#(Int, String)`

- For instance the type of a variable `map` that takes a mapping function from A to B, an arrays of A and retuns an array of B would be:  
   `#(#(A,B), [A], [B])`


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


The assignment (`=`) can be omitted if the body follows immediately. in that case the body doesn't need to declare the variables and can use them directly 

```js
#(v Int) {
  print("`v + v`")
}
```

Type of the variable can be generic

```js

// Block that takes/returns a T and take/returns a U
#(T, U)


```

### Note 

For v0.1 The type `Unit` means the boc returns nothing e.g. 

```js
println #(String, Unit) 
```
instead of 

```js
println #(String)
```
