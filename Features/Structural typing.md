```javascript
Person: {
  name String
  age Int
}
Alien: {
  name String
  age Int
  planet String
}

alien: Alien { 
  name: 'ET'
  age: 5
}

person Person = alien 

... 
person Person = Person {
  name: 'ET'
  age: 5
}

alient Alien = person /// Error: property 'planed' is missing 
...
person Person = {
  name: 'ET'
  age: 5
}
print: { person Person 
  person.info.properties().for_each { property Property 
    print '{property}: {refection(person, property)}'
  }
}
 ```


```javascript
Person: {
  name String
  speak: { 
    "Hi, my name is {name}"
  }
}
Computer: {
  model String
  speak: {
    "[{model}]: Beep bop"
  }
}
do_speak: { speaker {String {String}}
  print "{speaker.speak()}"
}
do_speak Person{"Oscar"}   // "Hi, my name is Oscar
do_speak Person{"ABC-123"} // [ABC-123]: Beep bop
do_speak {name:"Yz"; speak: {"Yup"}}  // Yup
```


#open-question Does the structure checks the names or just the values? 

#answered  Just the types, no the variable name
```javascript
//Are these two the same? 
{1 2 3} 
{a: 0 b: 0 c:0}
// yes

//They certainly have the same structure
{Int Int Int}
// To access the properties of an unamed block, use indexes
b {Int Int Int} = { a: 1; b:2; c:3}
// b type doesn't have names
b.0 // b.a -> 1
b.1 // b.b -> 2
b.2 // b.c -> 3
```

Rebuke: how to access the members though
```javascript
a {Int Int Int}
a.a
a.b
a.c  
```

Posible solution (April 9 2023)
Type declaration do include names and if ommited can use index like tuples:
```javascript
a {Int Int Int}
a.0 // a[0]
a.1 // a[1]
a.2 // a[2]
```

which is problematic because the type won't have the name of the variables

```javascript
block {Int Int Int} = {a:1; b:2; c:3}
block.a // how do we know there's an `a` ? 
```

Possible soution (April 9 2023), do include names
```javascript
block {a Int b Int c Int}= {a:1; b:2; c:3}
block.a // 1
block.b // 2
```

Possible solution: Include both, no names and named. If named, no semicolon and only in declarations; access by name. If no named, no commas nor semi colons and access by index. 

```javascript

//'Example without variable names in the parameter'
hello: {
    message {Int String} // accepts a block that has an Int and a String
    print '{message}'    // prints: {a:1 s:''}
    message(42 'Hello')   // executes it
    message.1 // 42
    message.2 // 'Hello'
    n s : message(-1 'Bye') // n == -1, s == 'Bye'
}
// invokes the above
msg: {
    a:1
    s:''
}
hello msg // executes the method, same as hello(msg)
Employee: {
    id Int
    name String
}
e: Employee{1 'John'}
hello(e) // prints '{id:1 name:'John'}' and thend the rest

// 'Example with names in parameter'
hi: {
   message {number Int label String}
   // can be accessed by name 
   print '{message}' // prints {number:0 label:''}
   message(1 'Nothing')
   message.number // 1
   message.label  // Nothing
   n s: message(2 'Something') // n == 2, s == 'Something'
}
msg:{
   a:1
   s:''
}
hi msg // comp errm msg doesn't have `id Int` and `name String` properties
P: {
   number Int
   label String
}
p: P{}
hi p // works
```
