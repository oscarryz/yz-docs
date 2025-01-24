
Variables
 ```js
// Use variablle_name Type = expr optional initialization
msg String = "Hello"
s String // declared but not initialized

// Use `:` for shortcut declaration + initalization
t : "Sweet"
```

Blocks of code ( bocs )
```js
// A boc is enclosed by `{` `}`
// e.g.
{
  // statements or expressions 
   msg : "Hello"
   n Int = 42
}

// A boc can be assigned to a variables too
greet : { 
  msg: "Hello world"
  print(msg)
}

// Bocs have types too. A boc type its signature and 
// the syntax is `#(` boc_variable(s) `)`
greet #(msg String) = { 
  msg :"Hello world" // msg has to be defined here and match the signature
  print(msg)
}
// A boc type followed by a boc literal defines it immediately
greet #(msg String) { 
  msg = "Hello world" // Uses the `msg` declared in the signature
  print(msg)
}

// More signatures examples
#() // An empty boc
#(String) // with a String
#(n Int)  // with a variable `n` of type Int
#(g T)    // with a variable `g` of generic type
#(fn #(T, E))  // with another boc that has two generics `T` and `E`
// More about generics below
```

Execute bocs
```js
// Invoke by using `(` args `)` as you would in a function on some programming languges
greet() // prints "Hello world"
```

User defined types and creating instances
```js
// Uppercase names that are not single letter define a new type
Greet { 
  msg: "Hello world"
}
g Greet //Declares variable g of type Greet
// When the type name is invoked it creates a new instance of that type
h : Greet() // creates a different instance

// User defined types also have a type signature, this represents the programmable interface the type offers
// For Greet above the type signature would be
Greet #(msg String = "Hello world")



```

Reading bocs variables
```js
// Use `.` to access boc variables 
f : { 
  m : "Hi"
}
print(f.m) // prints "Hi"
f() // executes `f` which just assign "Hi" to f.m (again)
g Greet = Greet() // creates an instance of Greet
print(g.msg) // prints "Welcome home"

// The main difference between lowercase named bocs and upper case, is lowercase
// are singletons, whereas uppercase create instances and can be used where a type can be used
```

Writing boc variables
```js
Greet : { 
  msg String // The boc Greet defines `msg`
  update_msg #(new_mesage String) {
    msg = new_message // inner boc modifies it
  }
}
g : Greet("Hola") // assigns "Hola" to message
// g.msg = "Adios" // Compilation error
g.update_msg("Bye") // Also ok
```

Generics
```js
// A Single Uppercase letter denotes a generic type
x T // x data is generic
x = "hi" // data is bound to String
// x = 1 // Compilation error, x was already bound to String
// Can be used in user defined types too
Box : { 
  T // Box has a generic type T
  data T // has a variable `data` of that generic 
  set_data #(new_data T) {
    data = new_data
  }
}
// T is a type parameter, a new instance can be created passing a type as argument: 
bs Box(String) = Box(String) // bs will be a Box of String
// or simply bs : Box(String)
bs.set_data("value") // sets data to "value"
bi : Box(Int) // bi will be a Box of Int
bs.set_data(42)

// If specific usage is needed for the generic value, then the compiler will validate, e.g. the method `len` is called, then the instantiation must have a len boc 
send #(value T) {
  lenght : value.len() // T has to have the method len #(Int) to be instantiated
  print("Sending `lenght` elements")
}
send("Hola") // compiles
send([1,2,3]) // compiles
send(0.0) // Doesn't compile, Decimal doesn't have `len()` boc 

```

Constructor variants (similar to Enum, Unions and SumTypes, but they are none of those)
```js
// Inside a type, name the variant with the parameters it will use.
// These variables will  become variables of the type itself.
// Overlapping variables will be consolidated into one.

// e.g. With the same variables (`msg`)
Greet : {
  Hello(msg:"Hello")
  Bye(msg:"Goodbye")
}
g Greet = Greet.Hello()
h : Greet.Bye() // h type is Greet

// With different variables 
Messge : {
  Text(text String="")
  Empty()
}
// `Message.Text` has to be used with a String
t : Message.Text("We're calling you for your auto policy")
// `Message.Empty` has to be used without parameters
e : Message.Empty() 


// The variant is part of the type signature
// For `Greet`
Greet #(
  Hello(msg String="Hello"),
  Bye(msg String="Goodbye")
)

// For `Message`
Messgage #(
  Text(text String=""),
  Empty()
)

// Variants can have generic data too
Option : {
  Some(value T) // defines the variant with a variable `value` of type T
  None() // defines a varian with no variables
}
maybe Option(T) = Option.Some("thing")
maybe_not Option(T) = Option.None()

// Can have different generic types
Result:{
  Ok(value T),
  Err(e E)
}
result Result(Int, String) = Result.Ok(1)
bad_result Result(Int,String) = Resutlt.Err(":(" ))

```

Conditional blocks of code using `when { cond => action }`
```js
// To know what kind of Greet we have.
// (back ticks ` are used for string interpolation)
g Greet = get_greet()
g when
{ Greet.Bye => "You say `g.msg`" },
{ Greet.Hello => "and I say `g.msg`"},

// To know what kind of message we have
m : get_message()
m when
{ Message.Text => "The message content is `m.text`" },
{ Message.Empty => "There is no message" },

// Accessing generic data
maybe : Some("something") // 
maybe when
{ Some => print("We have `maybe.value`)"},
{ None => print("We have nothing")}

// Conditionals can use boolean expressions too: 
when
{ a > b  => "A is greater than B"},
{ a == b =>  "A is equals to B"),
    { => "A si lower than B"),
```


Arrays and Dictionaries

```js
[Type] // array type
[]Type // Empty array of type Type
[T] // generic array


// Array literal
[1,2,3] // type [Int]
[ { print("Hello")}, {print("Goodbye")}] // type [#()]

// Array access returns Result(T, OutoOBoundError)
a : [ 1, 2, 3 ]
a[0] // Ok(1)
a[5] // OutOfBoundError(5)



[Key:Value] // Dictionary type
[Key]Value // Empty dictionary with key of type Key and values of type Value
[K:V] // generic dictionary

// Dictionary literal of type [String:String]
["name": "Foo",
 "last_name": "Bar"] 

// Dictionary access returns Optional(V)
d : [ 1 : 2, 3: 4] // [Int: Int]
d[1] // Ok(2)
d[5] // None()

// A dictionary Language to Bocs that take a String and return a String
greet [String : #(String, String)]
greet = [ 
  "English"   : { n String; "Hello `n`" }, 
  "Español"   : { n String; "Hola `n`"  }, 
]

say_hi #(to_whom String, language String) = { 
  boc : greet[language]
  boc when 
  { Some => boc.value(to_whom) },
  { None => "Language not found" },
}
say_hi("Alice", "English")// Hello Alice
say_hi("Bob", "Español")   // Hola Bob
say_hi("Charly", "Brown")   // Language not found

```

Concurrency

Control flow 
