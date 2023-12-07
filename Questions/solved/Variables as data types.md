

If a block has a type: 

    a {}

And it's intialized 

```
a {} = {}
```

Can return or take variables or expressions

```
a {Int} = {1}
```

But they need a variable name if meant to be used outside

```
a {i Int} = {i Int}
a.i = 0
```
And those variables can be block themselves 

```
b {b1 {} } = {b {}}
b.b1 = {}
```

And they can also have variables

```
b {b1 {c Int}} = { b1 {c Int} }
b.b1 = {c Int}
```

What's the type of a block returned by executing that inner block? 
```
... 
c: b.b1(1)
answer: 'c is of type Int'
```
Is the following valid
```
...
b.b1 = { 
    c Int
    c
}
aswer: 'Yes, it's taking and returning 1 value'
questin: 'Shouldn't be `{c Int Int}` as in `{c Int; c * 2}`' //tbd
```

What type is `b1`
```
b2: b.b1
answer: '{c Int}'
//or
b2 {c Int} = b.b1

```

Can we do this?

```
b2 b.b1
answer: 'No, that would be calling b1 with b.b1 as param'
```

Can we do this?

```
b2 {c Int}
answer: 'Yes please.'
// or 
b2: {c Int} // yes that too
```

Can we create an alias for a type like `{c Int}` ? Yes, what about upper case thing: `C:{c Int}` and then we can declare `b2 C`
Can you redefine the thing about using `C`?
```
b {b1 {c Int}} = { b1 {c Int} }
b.b1 = {c Int}

C: {c Int}
b {b1 C} = {
    b1 C
}
// Can I do this?
b.b1 = C // A: No, C is not a value but a type alias
D: {b1 C}
b3 D = b1 // Yes, it's possible and matches.

```

Solving the problem exposed in [this reddit thread](https://www.reddit.com/r/ProgrammingLanguages/comments/12rmkjc/comment/jgzvr2y/?utm_source=share&utm_medium=web2x&context=3)

Given: 
```
counter: { 
    count Int 
    { 
        increment: {count = count + 1}
        get: {count}
    }
}
```
What's the type for `counter`

```
counter { count Int { increment {} get{Int} } } = ... 
c1: counter()
// what type is c1?
c1 { increment {} get {Int} }
c2: counter()
// same type 
// Can I declare a type: `{ increment {} get {Int} }`
// Yes
Counter: { increment {} get {Int} } // o_O that's a body declaration
Counter { increment {} get {Int} } // that's  better
// Again using `Counter`
counter { count Int Counter } = {
    count Int
    Counter{}
}

```

The above doesn't work because would use the creation literal. 
But if instead of that literal we use a function call? 

```javascript
Counter { c Int get {Int}} = {
    c: 0
    get: {
       c 
    }
}
c1: Counter()
c2: Counter(1)
```

```
counter: {
    count Int
    {
        get: {
            count
        }
    }
}
c1: counter()
// c1 type
c1 { get {Int} }
Counter { get {Int} }
// c1 type
c1 Counter
// counter redefined using Counter
counter: {
    count Int
    Counter(
        get: {
            count
        }
    )
}
```


#answered Can't use variables as types, but we can create "aliases" for the types and then use those aliases, we would need to stop using the Object instantiation literal and use it as a function

```
Point{x Int y Int} = {x Int; y Int}
p1: Point()
p0: {x:0; y:0}
p0()

```
#challenged That would create an inconsistency because calling the alias ( `Foo()` ) would behave different from calling the block (`{a Int}()`)  might need to use `new(Point)` as in Go.

Hm does it though? When I call the alias `Foo()` I'm calling the body, not the signature... It might stil be the righ way to go


```javascript
// Problem
counter: {
    count Int
    {
        increment: {
            count = count + 1
        }
    }
}
c: counter(0)   // c.count == 0
c.increment()  // c.count == 1
d: counter(10) // d.count == 10, c.count == 10
// Solution
counter: {
   a_count Int
   {
      count: a_count
      increment: {
          count = count + 1
      }
    }
   }
}
c: counter(0)   // c.count == 0
c.increment()   // c.count == 1
d:counter(10)   // d.count ==0, c.count == 1
// c type definition
c { count Int increment {} }
// Alias using the `::` operator
Count::{ count Int increment {}} = {
    count Int
    increment {} = {
        count = count + 1
    }
}
// Type inferred form
Count: {
    count Int
    increment: {
        count = count + 1 
    }
}
// Usage from previous example
count: {
    a_count Int
    Count { a_count }
}
c: counter(0) 
c.increment()
d: counter(10) // c.counter == 1, d.counter == 10

```

#answered  Will use a `::` operator to create an alias to a type, so if the type is `String`  we can define: 

```
Word :: String
words []Word 
```
Would be equivalent to `[]String`

This way we can create custom types: 

```
Person::{
    name String
    last_name String
}
p Person = Person{'Yz' 'Language'}
// same as 
p {name String last_name String} = {'Yz' 'Language'}


```
Would be the same as

```
p { name String last_name String} = {name:'Yz'; last_name:'Language'}

```
The explanation for `Person{}` can be found in [Instances](../../Features/Instances.md)

#answered  We'll only use uppercase variables as datatypes e.g. `Point` but not `point`

- What's the type for person though? A: `{name String last_name String}`  
- What's the type for Word ?  A: `String`



