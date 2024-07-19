1. Syntax of the language and some example code snippets: 
2.``` // Comment /* muti line comment */ variable_name Type = expression // variable declaration and initialization (optional) variable_name : expression // variable declaration and type inferred from type // e.g. name String = "Alice" // defines variable named: `name` type `String` initialized with the string `Alice` age : 42 // defines variable named `age` intialized with value `42` and type `Int` (integer) inferred /* Block of code defined by using `{` `}` */ { message: "hello" times: 1 } // blocks are assignable to variables, the following // creates a variable named `block_1` with the block type inferred. block_1: { a : 1 } // The block type itself is defined in parenthesis, so the type of the variable `block_1` above would be: block_1 (a Int) // `block_1` is a variable of type block, containing a variable `a` of type `Int` // Blocks can be invoked by using `(` `)`, so to invoke the previous block we would write block() // array `[` `]` array: []Int // empty int array array_2: [1 2 3 4] // array with elements 1, 2, 3 and 4, type []Int inferred // dictionary `[` `:` `]` e.g the folowing defines a dictionary of type String to String `[String]String` // with the contents "name" -> "Alice" and "last_name" -> "Brown" dict: [ "name" : "Alice" "last_name": "Brown"] 2. Basic data types supported: - Int: unbounded integer type (like smalltalk, where there's no integer overflow) - Dec: unbonded floating point type - String : string literal delimited by pairs of " or ', always multiline and supporting interpolation with `$()` exemple: `name : "Bob"; s : "Hello $(name)"` - Array: []Type, array literal `[ val1 val2 val3]` - Dictionary: [KeyType]ValueType, dictionary literal `[ key_1: value_1 key_2: value_2]` 3. Control structures: None, only recursion 4. N/A




```
Tree: {
	Empty
	Leaf(Int)
	Node(Int Tree Tree)
	data T
	left Tree(T)
	right Tree(T)
}
Bool: {
	True
	False
}
enum Option 
Option: {
	Some(T)
	None
	data T
}
```