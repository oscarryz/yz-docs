
Should creating a new instance initialize all the uninitialized fields? Or should the compiler validate before calling? 


e.g. 

```js
// Uninitialized fields
Person : { 
  name String
  age Int
}
p : Person() // Should this be a compilation error? 
// vs
p : Person(name: "Alice", age:42) 
// or 
p : Person("Alice", 42)

// Or should be only an error if I try to use it? (which could be harder to detect)
...
foo: {
  q Person
  print(`q.name`)
}
foo(p) // will use p
...
```

If that's not a compilation error: How what is the program behavior when the variable is accessed. 
If that's a compilation error: How can we create a boc whose values are know later on? 

Possible answer: 
- Force to define a default value. 
- Define an `Option` value e.g. `Person: { name Option(String), age Option(Int) } `

 #open-question 

