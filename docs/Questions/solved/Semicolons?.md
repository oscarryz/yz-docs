We could ommit semicolons but code starts getting messy. 

Analyze and determine if they can be avoided. 

e.g. 

```javascript
// In declarations
a {Bool String}  // block that has Bool and String variables
// vs
b {Bool; String}

// In initializations
a = {b Bool s String}
// exactly the same as above but with new lines
a = {
    b Bool
    s String
}

//vs.
a = {b Bool; s String}
a = {
    b Bool; // ; less "needed" here
    s String
}


// In expressions
1 + 2  foo() 
// vs 

1 + 2; foo()

// In method call 

foo 1 bar() /*vs*/ foo 1; bar()
foo 1
bar()

```

#answered keep semicolons for multiple expressions in same line but skip it for declarations
