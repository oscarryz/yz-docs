
This whole document was created using the `::` notion which was later rejected, it help to inform though how instances will be created. See `Re-Summary` at the bottom and feel free to ignore the whole doc.

In Yz you can assign a value to a variable if it matches the variable type structure, for instance

```javascript
a Int = 0 
b Int 
b = a // possible because they have the same type, hence same structure

block_1 { } // `block_1` type is an empty block with no variables 
block_1 = { print 'Hello' } // it still can do things, but those things cannot be acessible from the outside

block_2 { } //  `block_2` type is an empty block with no variables 
block_2 = block_1 // compatible

// Create an instance of block_1

//~b1: block_1{}~
//~b1() // Hello~

// b2: block_2{} // Hello 

block_3  {} // same as block_1 and block_2
//~b3: block_3{} // Error: block_3 is abstract? 
//~b4 block_1 // Error, there's no type block_1, define it with uppercase Block_1
Block_1 {} 
b5 Block_1 // Declared
//~b5() // Compilation error: Block_1 has not been defined. 

```

Block types can only be assigned to uppercase blocks, but regular types still can be cloned/copied, but they still have the raw type, in the example above `block_1, block_2, block_3` is `::{}`
whereas `b4` type is `Block_1` (which in turn is compatible with `{}`) 

```javascript
something: {
    x {}
    x()
}
other: {
    x Block_1
}
something b1
something b2
something b3
something b4

other b4
other b3 // yes, still works

yet_another: {
    x block_1 // Compilation Error, there's no type block_1, define it with uppercase e.g. Block_1
}

```

#answered 
Summary: 
- To create a new instance use ~~`{}`~~  `()`  the type of the new instance will be the same as the thing being instantiated. 
- Variables cannot be used as types ( `e : {}; f e // trying to make f of type e`)
- To create a new type use an uppercase
- For generic methods see [Generics without <>](Generics%20without%20<>.md)

#challenged How to distinguish  between new instances and a parameter of type block? 

```javascript
foo{} /*vs*/ foo {}
while { true } {

}
// vs 
while{}  // creating a copy of `while` ? 
```

Possible answer: use significant white space, but that might be very very subtle (`foo{} vs foo {}`)

#answered 
Only new types can create instances, so `Foo()` creates an instance of `Foo` while `foo {}` passes an empty block to the method foo. The type `foo :: {}` is context dependent, if there's no `foo` in scope it defines it of type empty block 

Re-summary: 
- Only upper-case can create user types
- Only upper-case can be instantiated
- If in case of ambiguity try first to look for a method, and then define a variable if absent
    - `foo {}`  if there's no `foo` in scope defines a variable of type empty block `{}`
    

Re-resummary 27/10/2023
- Create instances with `()` e.g. `Point(1 2)`
- Define block type with `::` e.g. `Point :: {Int Int} = {x:0 y:1}`  or `e :: {}`
- As before only uppercase starting identifiers can be new types 
### Related
[Define new types](../../Features/Define%20new%20types.md)
[Instances](../../Features/Replaced%20features/Instances.md)
[The block type](The%20block%20type.md)