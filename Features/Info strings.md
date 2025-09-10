#feature
Information Strings. You can add a string before any element and will be available by calling  `std.info(element)` 

This text will be available as the `text` property of the returned block: 
```javascript

`A message`
message: 'Hello'
info(message).text // A message
```

You can add blocks there too, these block don't need to be valid yz code, but different tools might require them to be. e.g. `yz test` will require the `tests` variable 
to be valid or a serialization tool would require a valid string etc.

Example, the following block, has a info string, describing the block, what are the variables the block has, what tests can execute, version and author
```javascript
`
Prints the classics "Hello, World!" program to the screen

variables: {
  what String = 'Hello' // what message to display
  who  String = 'World' // what to say
}

tests: {
  assert say_hello()            == "Hello, World!"                   // Uses defauls      
  assert say_hello 'Hola'       == "Hola, World!"                    // Overrides first variable 'what'
  assert say_hello who: 'there' == "Hello, there!"                   // Explicitly overrides variable 'who'
  assert say_hello who: 'home' what: 'Welcome'  == "Welcome, home!"  // "Named parameters"
}
version: 1.0
author: 'Yz developers'

As you can see, this text right here (and the one at the begining) are not
valid programs while the other variables are. 
`
say_hello: {
   // Any element can can an info string
  'What message to display'
  what: 'Hello' 

   // Can have validation info
   // Or serialization info
  `
   validation: "\w*"
   json_field: 'xyz'
  `
  who:  'World' 
  // COuld be also used as running examples
  // that will be validated with yzc test  
  `
  Example: print 'Hello' 'world'
  Output: Hello, world!
  `
  print '{what}, {who}!'
}
```

To retrieve it use the `std.info` block and pass the element

```javascript
info: std.info say_hello
print info.text  // Prints the classics "Hello, World!"... etc.etc
info.tests()     // Runs the tests
info.version     // 1.0
info.examples()  // run the examples 

```
To have this information available the info string is compiled along with the element

