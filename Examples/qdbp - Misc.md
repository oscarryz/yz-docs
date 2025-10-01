#example


Interesting links:

[tomprimozic/type-systems](https://github.com/tomprimozic/type-systems)
[tomprimozic/type-systems - extensible_rows/infer.ml ](https://github.com/tomprimozic/type-systems/blob/master/extensible_rows/infer.ml)

For type inference
[Extensible records with scoped labels](https://www.microsoft.com/en-us/research/publication/extensible-records-with-scoped-labels/)
[An introduction to Object-Oriented Calculus](https://pages.cpsc.ucalgary.ca/~robin/FMCS/FMCS2014/An%20introduction%20to%20OOC.pdf)



```js

/*
 Multiline comment
*/
// Single line
three: 3
hello_world: "Hello, world!"

// We can bundle data togher using blocks
block: {
    method1: {3} // a method named `method1` returns 3
    // method2 takes and argument and returns it
    method2: {
       arg T;
       arg
    }
}

// We can invoke ethods within a block
empty_block: {}
// We can invoke methods within a block
three = block.method1()
// Methods with arguments might have the arugument name specified
four: block.method2(4) // or
four = block.method2 4 // or
four = block.method2(arg:4) // if specified parenthesis are needed

// Blocks are similar to functions and to invoke it use `()` as you see above
// Also can have inner blocks (like method1 and 2), if the method name is
// a non-word name the `.` can be skipped during the call.
// e.g. of non-word names ! ? >> + - << >>= etc
identity_fn: {
    ! : [T]{x T; x}
}
five: identity_fn ! 5

// Methods withing a block cannot be overloarded by paramter names, but they can have default values
print_args: {
    !: <T>{
        arg1 T // although probably shouldn't be T but Printer
        arg2 T = Printer{}
        arg1.print()
        arg2.print()
    }
}
Printer: {
    print: {}
}

// To evaluate muitiple expressions in the same line separate them with `;`
print_args "Hello world" // assuming String has a `print` method
print_args "Hello" "world"
print_args "Hello world" ; print_args "Hello" "world"
print_args(arg1: something_with_print{})
print_args(arg2: something_with_print{}) // not possible, arg1 doesn't have a value
// However all of these would execute concurrently, if sequence is needed, assign the result to a variable
// assinging to `_` will make the execution to wait until the call is completed
_: print_args "Hello world"
_: print_args "Hello" "world"

print_args 3 // print_args is generic

Int: {
    + {other Int; ...}
}
five: 3.+(2)
// Because methods with at lest one parameter can ommit parenthesis and methods named with non-words can omit `.` the following
// is the same
five: 3 + 2

// Defines a circle
small_circle:{
    radius: {3}
    diameter: {radius() * 2}
}

// There's no inheritance but to reuse code,
// new variables can refer to existing variables.
// So the following
big_circle: {
    radius: {6}
    diameter: small_circle.diameter
}
twelve: big_circle.diameter()
bigger_circle: {
    radius: big_circle.radius
    area: {radius() * radius() * pi}
}
one_hundred_eight: bigger_circle.area()

```
