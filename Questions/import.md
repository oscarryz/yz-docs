#open-question 
The `import` keyword would import in the current scope the variables in the scope imported. 


It is mostly like a macro. 


```js
a : {x Int, y Int}
b : {
  mix a
  x = 0 // valid, x was created by the import above.
}
b.x = 1 // valid
b.y = 2 // valid
c : { 
  mix b // would be the same as import a
  x Int // invalid, x being redeclared
}
```

When instantiable blocks would work the same, but every instantiation will create a copy

```js
Animal : {
  talk: #(String)
  // 30, 40 other methods below 
  walk: #(String) = { "I'm walking" }
}
Dog: {
  mix Animal
  talk = {"Woof woof"}
}
d : Dog()
print(d.talk()) // Woof woof
d.walk = {"nope"}
print(d.walk()) // nope

d2: Dog()
// un affected.
print(d2.walk()) // I'm walking

```

It is useful to bring a lot of symbols into the scope. 

Probably `import` is not the best name, tbd what name to use. 

