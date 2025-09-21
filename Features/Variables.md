#feature

Variables are declared with an identifier followed by a type identifier or type signature.

e.g. 
```js
// declares a variable `msg` of type String
msg String

// declares a variable `salute` of type block that returns a String
salute #(String)
```

Variables can be initialized with `=`

```js
// declares and intializes `msg` with value "Hi"
msg String = "Hi"

// declared and intializes `salute` with a boc value
salute #() = {
   "Hello world"
}
```

The shorthand `:` can also be used to declare and initilize the variable inferring the type from 
the value 


```js
msg : "Hi"
salute: {
   "Hello world"
}
```

For easy of use bocs can also declare the type and add a block next to with with the implementation

```js
greet #(msg String, to_whom String,String) {
   "`msg` `to_whom`"
}
```

See: 
[Signatures + Literals duplication](../Questions/solved/Signatures%20+%20Literals%20duplication.md)
