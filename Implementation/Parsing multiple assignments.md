#impl

[Reddit discussion](https://www.reddit.com/r/ProgrammingLanguages/comments/1h7741a/parsing_multiple_assignments/)


### Context

In Yz `=` and `:` operators return the RHS value, they are expressions (?)

So 
```
b = a = 1 // creates variables a and b with values 1
```

Also, Yz support assigning multiple values. 
```
a, b = foo()
```

How are these two interact? 

### Discussion
e.g. 


Matching number of terms.
```js
a, b, c = y, x, z
```

*`,` needs to have higher precedence than `=`*
*`()` needs to have higher precedence than `,`*

Different number of terms in RHS (needs to match number of terms too)

```js
a, b, c = foo() // foo has 3 return variables
```

The Reddit thread mentions to consider nested terms. 

```js
(a, b, (c, d)) = (10, 20, (30, 40))
```

Yz would need to specify how a different number of terms on the RHS would match the LHS

e.g. 

```js
a, b, c = (1,2,3,4,5,6,7,8)
// or 
a,b,c = f() // f returns 8 values
// a: 6, b: 7, c: 8

```
What would this be? 

```js
a, b, c = f(), g() // both f, and g return 4values

// Defining f and g as
f: {1,2,3,4}
g: {5,6,7,8}
a, b, c = f(), g()
// Would be the same as 
a, b, c = (1,2,3,4),(5,6,7,8)
/*
a: 6
b: 7
c: 8
*/
// Compilation error, f() not used  
```


To specify how the values should be unpacked parenthesis are needed

```js
a, (b, c) = f(), g()
/*
a: 4
b: 7
c: 8
*/
```

Because the parenthesis would group thing together, this part is the first group, that part is the second group. 


### Questions

Can we destructure arrays? 

```
a, b, c : [1, 2, 3] // No
```

Invalid,  this would assign the array to `c` and invalid for `a` and `b`

It would require parenthesis.
```
(a, b, c) : [1, 2, 3] // Yes
```

What about "tuples" ? 

```
a, b, c : {1, 2, 3} // No
```
Neither, that would assign the "tuple" to `c`

It would either require to execute it

```
a, b, c : {1, 2,3}()// executes the tuple
```

Or use parenthesis on the LHS

```
(a, b, c) : {1, 2, 3} // Yes
```


What about composing complex objects? 

```
p: {x: 0, y: 0}
p = (1, 2) // invalid
p = {1, 2} // invalid needs attributes x and y
```

What about nested LHS

```
// TODO analyze this further
// What is the required shape on the RHS to assing 
// values to ((a, b), c), d 
((a,b),c),d = ??
```


### Operator `=` can't result in an expression

This leads to the decision of making `=` operator to not be an expression because it would complicate a list of expressions. For instance, this would look like two variable declarations: 
```
a = 1, b = 2
```

But if assignment is an expression AND comma has **higher** precedence would lead to interpret it as: 

```
a = (1,b) = 2 
```

Which is an error.  To avoid it it would require to use parenthesis: 
```
(a = 1), (b = 2)
```

Which seems like a lot for just declaring two variables. 

If comma has a **lower** precedence would make things like 
```
a, b = 1 , 2 
```

Be interpreted as 
```
a, (b=1), 2
```
Which is fine, but is not the intention of the multiple assignments, use would have to write:
```
(a, b) = (1, 2)
```


If `=` is not an expression (and `,` still used to separate expressions instead of `;`) 

```
a=1, b=2 // two variables
a,b=1,2 // (a,b)=(1,2)
```


#### Update
It seems `,` can have a different parsing precedence in certain constructs like assignment, so the following `a,b = 1,2` doesn't evals to `a, (b=1), 2`

The simplest solution is to 

a) Make assignment (both `=` and `;` ) not expressions 
b) Use semicolon to separate statements, so `a = 1, b = 2` should be written as `a = 1; b = 2` if they have to go in the same line. 

For this to happen I would need to re-work the automatic comma insertion 


```js
// A block with two variables:
// `name` of type "String"
// `type_system` of type "Array of Strings"
language #(String,[String])
language= {
    name: "Yz", 
    type_system: ["static", "strong", "structural"],
}



main #(String,[Int],[String:Bool])
main = {
	msg: "Hello",
	array: [1, 2 ,3 ],
	dictionary: [
		"ready" : false,
	],
}

a Int
b Int

#(Int,Int)
some = {
  a = 1,
  b = 2,
}

#(Int,Int)
other = { 
  a = 1 + 2
  b = 3 - 4

  a = 1.+(2)
  b = 3.+(4)

  a = 3
  b = 7
  a, b
  a, b = 3 - 4
}

point: { x: 1, y : 2}

point = {3, 4} // no, it needs var names

anon_point: { 1, 2}
anon_point = { 3, 4 } // yes, have same signature #(Int,Int)
{a, b} = point // no, it needs var names
{ a, b } = anon_point // yes but looks weird
a, {b, c} = g(), f() // no, and probably no to `a, (b, c)` either



```


In Yz everything is a block of code (boc for short) and a boc is a list of expressions and statements. 

An statement is a variable declaration or one in the following keywords `use`, `return`, `break`, `continue`.

An assignment is an expressions 

```
{
  expr1
  expr2
  stmt1
  stmt2
  expr3
}
{ 
  1
  name: "Hi" // expression value is  "Hi"
  name = "bye" // {1, "Hi", "bye", "bye"}
  name
  age Int // (1 Hi bye bye; age Int; 2 3 4)
  2, 3, 4
}

{
a: 1
m String
b: 2
}
#(a Int, m String, b Int)
#(b Int) = {
  a : 1
  m String
  b : 2 
}

g #(Int) = {
  a : 1
  m String
  b : 2
}
g() // return 2 
f #(a Int) = { 
  b : a + 1
}
f() // compilation error 

 a, b, c = x, y, z // only assigns x to c, the rest are "normal" expressions

a : b : c = 1
x : y : z = 2
f #() = { 
    a
    b
    c = x
    y 
    z
}
f #() = {
  a, b, c = x, y, z
}
a, b, c = (x, y, z)
a, b, c = {x, y, z}()

a, b, c = x, (y, z)

```

Ways to access boc variables

1. Method invocation `(` _params_ `)`  (application?)
2. Dot notation `.` ( member access ) 

Discarded way to access boc variables

Destructuring

```
a,b, = point // trying to access point.x, point.y
pont = {a, b} // unless point signature is #(Int,Int) and a and b are Int,Int

```

Creating a new block trying to capture might not work as it would lose a reference to it
```
{a, b} = point // we lose the reference to {a, b}

```

It would need to define a new boc and then assign it
```
np : {a Int, b Int} = point // attempt to keep the reference

```

Which is equivalent to:
```
np : {a Int, b Int}
np = point
```


To be able to assign more than one variable they need to be in parenthesis

```js
// setup 
a: b: c: 1 // declares 3 vars
g #(Int,Int) // g returns 2 integers

// First regular 
a, b, c = g()
// Evaluates to, three expressions
a
b
c=g() // and c takes the last value from g

// Second what is intended
a, (b, c) = g()
// Evaluates to two expressions
a
b,c = g() // takes the g.-1, g.-2 values
```

Even if they are in different lines 

```
main: {
  1 + 1
  (a, b) = g()
}
```

Which looks rather odd. A possible solution to do an automatic semi colon insertion 
```
main: {
  1 + 1; 
  a, b = g() // a, b could then be considered a single expression 
}
```

But that would require tuples bocs to use it too

```
some: {
  1; 2; 3;
}
```

# Commas and semicolons

The aim is to use commas as much as possible while minimizing the need for semicolons

```
// Without commas or semicolons
Person : {
  name String
  dob  Date
}
print #(
  p Person
    String
)
print  = { 
   p Person 
   "This person name is `p.name`" 
}

// With both
Person : {
  name String
  dob Date
}
print #(
  p Person
    String
)
print = {
  p Person;
  "This person name is `p.name`";
}
```

If inside `()` or `[]` use commas
If inside `{}` use semicolons
```
{ 1, 2, 3 } /*would become*/ {1; 2; 3}
```


