
#answered  No need for macros nor pattern matching atm, see [Type variants (similar to SumTypes)](Type%20variants%20(similar%20to%20SumTypes).md) and [Type Alias](Features/Type%20Alias.md)



How to make Extensibility, Enums, Union Types, and Specialization in a way that goes with the spirit of Yz which is: everything is a block of code, and things are publicly accessible while keeping the syntax minimum.

### Extensibility
We want to add more functionality to a type without having to rewrite it all. 
We can assign attributes from an existing type directly

```js
Person : { 
  name String
  say_hi : { print('My name is `name`') }
}
PolitePerson : {
  name : Person.name
  say_hi: Person.say_hi
  say_bye: {print 'This is `name` saying good bye'}
}

```

This could be automated in the future with macros
```js
'[derive:Person]'
PolitePerson: {
  say_bye: {print 'This is `name` saying good bye'}
}
```

The default way is to use `import`
```js
Person: {
  name String
  say_hi: {print("My name is `name`")}
}
PolitePerson: {
  import Person
  say_bye: {print('This is `name` saying good bye')}
}
```
### Enums

Not really Enums in the whole sense, but just a way to know if a given type can be used when other is expected, more like aliases

```js
// Yes
DayOfWeek : {
  name String
}
Monday: DayOfWeek
Tuesday : DayOfWeek
Wednesday: DayOfWeek
// Same as above with different style

DayOfWeek : {
  name String
  ... 
}
// Strictly speaking this would work 
Monday: 
Tuesday:
Wednesday: 
  DayOfWeek 

// Other way, is to create types inside
DayOfWeek: { 
  Monday: {}
  Tuesday : {}
  Wednesday: {}
}
// Sugaring 
DayOfWeek: { 
  Monday
  Tuesday 
  Wednesday
}
// nah.. it's ok to have inner types, but not like that, as it should require to have an instance to create the inner instance


```
### Union Types
```js
Result: {}
Ok: {}
Error: {}
result: {
  Result: {
    T 
    E
    and_then #(U, #(T, Result(U,E))) 
  }
}
```

### Specialization
This could be done by creating a type alias which takes all the feature of the original type. 

```javascript

// Case 1
// Creating type, and then instances that implement the method
// In this case `Bool` has the method `if #(#(T),#(T),T)` 
Bool: { 
   if #(is_true #(T), else #(T), T) 
   ? : if
}
// The instances `true` and `false` assing the given block to the 
// if attribute
true: Boolean()
false: Boolean()

// `true` evaluates the `is_true` branch
true.if = {
  is_true #()
  else #()
  is_true()
}
// `false` evaluates the `is_false` branch
false.if = {
  is_true #()
  else #()
  else()
}
compare_numbers: { 
  condition Bool
  //condition.if(is_true:{"It was true"}, else:{"It was false"})
  condition ? {
    print("It was true")
  } { 
    print ("It was false")
  }
}
compare_numbers( 1 > 2 ) // prints "It was false"

Bool: {
  ... 
  
}
True: { 
  import Bool
  if = { /*ovewrite if*/ }
}
False: { 
  import Bool
  if = { /*ovewrite if*/ }
}
true: True()// the singleton 
false: False() // the singleton


// Case 2
// Creating a type an synomums
Result : {
  T
  e E
  and_then #(operation #(T), Result(T,E))
}

Ok : Result(T)
OK.and_then = {
  operation #(T)
  operation()
}
Error: Result(E)
Error.and_then = {
  operation #(T)
  Error(e)
}
```