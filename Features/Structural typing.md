#feature
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
person Person = Person (
  name: 'ET',
  age: 5
)

alien Alien = person /// Error: property 'planet' is missing 
...
person Person = {
  name: 'ET'
  age: 5
}
print: { 
  
  person Person 
  person.info.properties().for_each({ 
    property Property 
    print '`property`: `refection.get(person, property)`'
  })
}
 ```


```javascript
Person: {
  name String
  speak: { 
    "Hi, my name is `name`"
  }
}
Computer: {
  model String
  speak: {
    "[`mode`}]: Beep bop"
  }
}
do_speak: { speaker #(String, speak #(String))
  print "`speaker.speak()`"
}
do_speak(Person("Oscar"))   // "Hi, my name is Oscar
do_speak(Person("ABC-123")) // [ABC-123]: Beep bop
do_speak({name:"Yz"; speak: {"Yup"}})  // Yup
```

Both names and types are checked if variable name is included.
 

```javascript

//'Example without variable names in the parameter'
hello: {
    message {Int String} // accepts a block that has an Int and a String
    print '`message`'    // prints: {a:1 s:''}
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
hello(msg) // executes the method, same as hello(msg)
Employee: {
    id Int
    name String
}
e: Employee(1, 'John')
hello(e) // prints '{id:1 name:'John'}' and thend the rest

// 'Example with names in parameter'
hi: {
   message #(number Int, label String)
   // can be accessed by name 
   print('`message`') // prints {number:0 label:''}
   message(1, 'Nothing')
   message.number // 1
   message.label  // Nothing
   n, s: message(2, 'Something') // n == 2, s == 'Something'
}
msg:{
   a:1
   s:''
}
hi(msg)// comp errm `msg` doesn't have `number Int` and `label String` properties
Pa: {
   number Int
   label String
}
p: Pa()
hi p // works
```
