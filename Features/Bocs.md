#feature

In Yz, everything is a block of code (boc for short)

A block plays the same role packages, modules, functions, methods, procedures, closures, objects and classes play in other programming languages. 

The construct is extremely simple and flexible

To create a block put the code between `{` and `}`
```javascript
{
   print('Hi')
}
```

Assign a block to a variable and execute it

```javascript
hi: {
    print('Hi')
}
hi() // Hi
```

Blocks naturally can have variables

```javascript
hi: {
   text: 'Hello'
   recipient: 'World'
   print '`text`, `recipient`!'
}
```

And as we just saw, those variables can be other blocks

```javascript
hi: {
   text: 'Hello'
   recipient: 'World'
   action: {
       print '`text`, `recipient`!'
   }
   action()
}
hi() // Hello, World! // calls action() which prints the message
```

We can access the variables inside a block by using the `.`  notation. 
For instance we could've call the `action` variable directly like this:

```javascript
hi.action() // Hello, World!
```

Using the `.` notation we can modify the variables before executing the block and retrieve them after the execution

```javascript
// Changing the defaults
hi.text = 'Goodbye'
hi.recipient = 'everybody'
hi() // Goobye, everybody!
hi.text //'Goodbye'
```

As a convenience for the modification and retrieval, we can use a function call syntax. We can retrieve multiple values at once

```javascript

a, b, c, d : hi('Goodbye', 'everybody')  
// parenthesis are optional when arity >0
// Before execution
//   hi.text => Goodbye 
//   hi.recipient => everybody
// During
//   prints 'Goodbye, everybody!'
// After
//   a = 'Goodbye'
//   b = 'everybody'
//   c = #(String) // a block that returns a String
//   d = 'Goodbye, everybody', the result of the invocation, which is the last statement.
```

The assign order starts at the bottom to the top and assignees from the most right to the left.


```javascript
z : hi('Yz')             
// hi.text => 'Yz'
// z => 'Yz, world!' which is the last executed expression `hi.action`
y z : hi('Yz', 'blocks')  
// hi.text => 'Yz'
// hi.recipient =>'blocks'
// y => {String} the `action` String block
// z => 'Yz, blocks!'
```

Finally you can use the name of the variables as arguments to specify what variable you want to set:
```javascript
hi(action: {print('Nothing')}) 
// hi.action => a string block {String} that prints 'Nothing'
// prints: "Nothing"
hi(recipient: 'Yz world') 
// hi.recipient => 'Yz world'
// prints: "Hello, Yz world!"

```

Often times you need a copy of the block to have different variables. You can declare a block as "instantiable" by naming it with starting uppercase (but not all uppercase) and create a new instance with the "instantiable literal" construct.

```javascript
Hi: {
    text: 'Hello'
    recipient: 'World'
    action: {
       print('`text`, `recipient`!')
    }
    action()
}
a: Hi() // use default values
b: Hi('Hola' 'mundo') // reassinging variable in order, can ommit name  
c: Hi(recipient:'again') // reassigning directly by name 
d: Hi(action: {'Nothing'})

// You can still call them 
a() // Hello, World!
b() // Hola, mundo!
a() // Hello, World! // still hello world
c() // Hello, again!
d() // (doesn't do anything)

// But it would be more natural if you use them as objects
// each one would be a different instance
a.action() // Hello, World!
b.recipient // 'mundo' 
// etc
a, b, c, d, e : Hi(), Hi('Hola', 'mundo'),  Hi(recipient:'again'), Hi(action: {'Nothing'})
```
