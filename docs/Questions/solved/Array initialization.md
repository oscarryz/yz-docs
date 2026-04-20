#answered in [Array](../../Features/Array.md)

```javascript
a []Int // array of ints
a = [] // empty array
b: [] // error: No data type
b: []Int // How is this different from declaration? 

b:[] // can this be an array of anything and then set the type when used? 

b << 1 // b []Int

```

What if is a a parameter?

```javascript
do_with_array: { b: []Int // it would have a type or be re-initialized... 
...
}
c: [1 2 3]
d: ['a' 'b' 'c']
do_with_array c 
do_with_array d 
```

#answered
Initialization literal needs a type if empty making it equivalent to just declaration. 


```javascript
// These are basically the same
a []Int 
a: []Int
```

More examples

```javascript
a [] Int // a is an array of ints
a: []Int // a in an empty array of ints basically the same as above
a: [1 2 3] // a in an array of ints initialized with 1 2 3
block: { a []Int //block can receive an array of ints as params
}
block [1 2 3] // ok
block [true] // error

block_string: { a: []String 

}
block_string ['a']
// The `numbers` variable is initialized
block_default_args: {
  numbers: [1 2 3 4]
  print '{numbers}' 
}
block_default_args() // [1 2 3 4]
// when parameters are passed they override the default value
block_default_args [0 0 0] // [0 0 0]
```

See [Default Args](Default%20Args.md)


