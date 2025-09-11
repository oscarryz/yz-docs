#replaced-feature

Replaced with [Creating instances](../../Questions/solved/Creating%20instances.md)

A new type can be created by creating an alias of a type signature. 

First let's see what a type signature is and how to use it
```javascript
// Type signature: Is a block with a variable count of type Int
#(count Int)

// Can declare a var with that type signature
b1 #(count Int)
// Can initialize a block matching that signature
b1 = {
    x: 1
    y: 2
    count: x + y // count of type Int declared and signature satisfied
}
b1() // returns count, 3
```

Usually the type signature can be inferred from the content
```javascript
b1: {
    count: 10 // b1 type signature is {count Int}
}
b2: {
    10 // b2 type signature is {Int}
}

```

To create copies from a block create a ~~type signature alias using the `::` operator and then assign a block to that alias~~  (there's no alias)

```javascript
b1: {count: 10}
Counter #( count Int )
Counter = b1 // assign the body of b1 `{count:10}` to the type alias `Counter`
// To create a copy use the instantiation construct
c1: Counter() // creates an instance
b1.count // 10
c1.count // 10
c1.count = 12 // 12
b1.count // still 10 
```
Usually you assign the block directly to the alias

```javascript
Counter: { 
    count: 10
}
c1: Counter()
b1: Counter()
b1.count // 10
c1.count //10
c1.count = 12 // 12
b1.count // still 10
```


This is another way, by creating anonymous blocks that capture the environment of the enclosing block, but this style might be too cumbersome, yet, valid.

```javascript
point: {
   x Int
   y Int
   // returns a block that has variable `x` with the value of x
   // and a variable `y` with the value of y 
   {x:x, y:y}
}
p: point(10,20) // returns the last executed which is the anon block
p.x  // 10
p.y  // 20

thing: {
   a_name String
   {
      name: a_name
   }
}
duck: thing("Duck")
duck.name // 

```


```javascript
// This is valid, but it would be better to define the type outside
person: {
    Person : {
        _name String
        name #(String)
        greet #(Person, String)
        builder #(String, #(String))
        call_me #(#())
        ok      #(#(#()))
    }
}
// Yes
p: person()(
    _name: 'Oscar'
    name: {s String}
    greet: {p Person; "Hola `p.name()`"}
    builder: {s String; {s}}
    call_me: {{print('OK')}}
    ok     : { do {} { print 'Ok'; do()}}
)
```



### Good

```javascript
a:thing{}
a.name 
// a type 
a thing // not possible no differentiation from calling a(thing)
a Thing // Yes, this is what I've been using
```



[Type Alias](Features/Type%20Alias.md)