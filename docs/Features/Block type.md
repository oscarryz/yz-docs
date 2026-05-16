#feature
# Boc Signature / Interface `#()`

The boc type is defined by its signature, and this signature is the interface for the boc.

Every block of code `boc` has a type. Just like a number literal `1` has type `Int` or `1.0` has type `Decimal`, a boc literal has a boc type defined by `#(...)` where `...` represents optional variables or expression types.

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
- Generic with constraint:   
   `#(o O Ord)`

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


The `=` can be omitted if the body follows the signature immediately — this is the **boc declaration** form. The body can use the signature's variables directly without re-declaring them:

```js
#(v Int) {
  print("`v + v`")
}
```

When `=` is kept, it is the **boc expanded form** — the body must re-declare all parameters:

```js
#(v Int) = {
  v Int
  print("`v + v`")
}
```

Type of the variable can be generic

```js

// Block that takes/returns a T and take/returns a U
#(T, U)


```

## The Signature Is the Interface

A boc signature `#(...)` is intentionally dual-purpose:

**1. Type constraint (structural typing).** Any boc with the right shape satisfies the signature. No explicit `implements` needed.

```yz
Greeter #(greet #())     // interface: requires greet #()
```

**2. Access control (encapsulation).** Only the fields and methods declared in the signature are 
visible to external callers. Fields omitted from the signature are hidden. Inferred boc types 
will always expose all the variables

```yz
Person #(name String, greet #())   // only name and greet are visible
Person = {
    name String      // public — also in signature
    password String  // internal — not in the signature, not accessible externally
    greet #() { print(name) }
}
```

These two concerns are fused by design. Yz defaults to **public-by-default**: if no signature is written, all fields are accessible. Writing a signature simultaneously narrows the public interface AND declares the type constraint. If full field visibility with a type constraint is desired, include all fields in the signature.
