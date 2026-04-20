#answered Use them all. Each one preferred for a different use

`#(){}` For variable declaration
`{}` For lambdas
`#()={}` For anything else

--- 


A boc type defines the signature for it, and a literal implicitly defines the signature too


```js
// type 
#(s String)

// declaration 
foo #(s String)

// delcaration + initialization
foo #(String) = { 
   s: "fooo"
}

// shorcut decl + init with literal
foo : {
	s: "fooo"
}

// block literal with signature 
foo #(String) {
   s: fooo"
}
```

The problem is the latest, the `#(){}` syntax used to avoid having to repeat all the variables 
the type defines, creates a "preferred" way that might look in consistent with others. 


Example where it helps: 

```js
// you don't have to re-declare each part of the signature in the body
address #(street Strign, number Int, post_code String, country String, String) {
   "`street` no. `number`. ZipCode: `post_code`, `to_3la(country)`"
}
// "Corect" style where the boc literal needs to declare the var to match 
address #(street Strign, number Int, post_code String, country String, String) = {
	street String
	number Int
	post_code String
	country String 
   "`street` no. `number`. ZipCode: `post_code`, `to_3la(country)`"
}
```

But when using both it seems to be inconsistent: 


```js
// streamlined version, the variable declared on the type
// is availale on the body
foo #(s String) {
	s.len().to_str()
}
// note the `=`
bar #(s String) =  {
    s String // variable has to be declared to match the type
}

baz: {
   s String // poc literal 
}
```

#answered Use them all. Each one preferred for a different use:

### Blocks of Code (Boc)

In this language, a block of code, or "boc," is a first-class citizen. It's a reusable block of logic and data that you can define, assign to a variable, and pass as an argument to other functions. Bocs are defined within `{}` braces and come in three main flavors, each with a specific purpose.

---

### 1. Function Bocs

Use this style to declare a boc that acts like a function, taking parameters and returning a value. The `#()` syntax defines the boc's **type signature**, including its parameters and return type.

The `:` operator is used for **inferred declaration**, making the syntax concise.

**Syntax:**

```
// variable_name : #(parameter_1 type_1, ..., return_type) { ... boc_body ... }

// Example:
// A boc that takes a string `s` and returns a string.
reverse_boc : #(s String, String) {
     s.reverse()
}
```

---

### 2. Anonymous Bocs

Use this style when you need to pass a boc directly as an argument to a function. The syntax is concise because the boc's type is **inferred** from the function you're passing it to.

**Syntax:**

```
// function_name(..., { ... boc_body ... })

// Example:
// The `filter` function expects a boc with a single `n` parameter and a boolean return type.
[1, 2, 3].filter({ n Int; n > 2 })
```

---

### 3. Named Bocs

This style is for explicitly declaring a reusable boc that you want to refer to by name. This is useful for creating modular code and assigning a clear name to a piece of logic.

**Syntax:**

```
// variable_name #(parameter_1 type_1, ..., return_type) = { ... boc_body ... }

// Example:
// Define a reusable `reverse` boc.
reverse_named #(s String, String) = {
     s.reverse()
}

// Now you can pass the named boc to a function.
greet("Hello World", reverse_named)
```