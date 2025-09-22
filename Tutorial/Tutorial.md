
Variables
 ```js
// Use variable_name Type = expr optional initialization
msg String = "Hello"
s String // declared but not initialized

// Use `:` for shortcut declaration + initalization
t : "Sweet"
```

Blocks of code ( bocs )
```js
// Bocs: Functions and Objects Combined

// In this language, a boc is a unique construct that acts as both a function (callable code) and an object (containing data). This means a boc can be invoked like a function and its members can be accessed like an object's properties.

greet : { // Declares a boc named 'greet'
    msg : "Hi" // A member (variable) of the boc
    say_hello #(name String) { // A method (function) of the boc
        print("`msg`, `name`!") // Accessing the boc's member 'msg' within the method
        msg = "Hello again" // Modifying the boc's internal state
    }
}

// Calling the boc's method:
greet.say_hello("Alice") // Output: Hi, Alice!
greet.say_hello("Bob")   // Output: Hello again, Bob! (State persists)

// Accessing the boc's member directly:
print(greet.msg) // Output: Hello again (Modified by the method call)

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

// Bocs have types too. A boc type is its signature and 
// the syntax is `#(` boc_variable(s) `)`
greet #(msg String) = { 
  msg :"Hello world" // msg has to be defined here and match the signature and be assignable
  print(msg)
}
// A boc type followed by a boc literal defines it immediately (no need to add `=` between signature and body )
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

Executing bocs
```js
// Invoke by using `(` args `)` as you would in a function call 
// on some programming languges
greet() // prints "Hello world"

x : foo() // The return value of foo is assigned to x. More details below
```

Parameters and return values 
```js
// The boc variables are not strictly parameters nor return values, they are both
// The `()` notation is a shortcut to assign or read values eg
greet("hi")  // assign "hi" to `msg`
greet() // but then when executed, it will override it with "Hello world"

// Likewise, return values are the last expression executed, in the greet example that is `print(msg)` which prints to the console and returns the printed value, which is "Hello world"
// So the following assigns "Hello world" to `pm`
pm : greet() 

// There could be more than one argument and they can be used 
// as more than one return value(s).
// e.g 1 return value
f #(Int) { 42 } // f is a boc with an Int value, which is the last executed and can be used as return value
n : f() // n is not 42

id : #(n Int) {  // fg names the variable `n` of type Int and simply return it
  n 
}
one : id(1) // one gets the value `1`

// eg. 2 return values
// The following boc takes two String variables
// and returns them in reverse order
swap #(a String, b String) {
   b
   a 
}
a : "hello"
b : "world"
// The last value in the boc is treated as return value in order
a, b = swap(a,b)
// a is now "world"
// b is now "hello"

// If a boc has no expression at the end the previous expression is used 
// if none exists the type is #() 
noop #()
// x : noop() // compilation error, `noop` doesn't return any value
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

// For `Greet` above the type signature would be
Greet #(msg String)
// It reads `Greet` is a type with a String variable named `msg`



```

Reading bocs variables
```js
// Use `.` to access boc variables 
f : { 
  m : "Hi"
}
print(f.m) // prints "Hi"
f() // executes `f` which just assign "Hi" to f.m (again)
// Note, bocs retains state between execution, unless of course the execution assings a new value. In the example above, f.m gets reassigned (to the same value "Hi")

g Greet = Greet() // creates an instance of Greet
print(g.msg) // prints "Welcome home"

// The main difference between lowercase named bocs and upper case, is lowercase
// are singletons, whereas uppercase create instances and can be used where a type can be used
```

Writing boc variables
```js
GreetUpdater : { 
  msg String // The boc Greet defines `msg`
  update_msg #(new_message String) {
    msg = new_message // inner boc modifies it
  }
}
g : GreetUpdater("Hola") // assigns "Hola" to message
// Assignment only works if the variable assigned is inside the scope where is attempted to be assigned. 
// For instance the following
// g.msg = "Adios" 
// Would cause a compilation error, because the variable `msg` is inside the boc `g` (defined in `GreetUpdater`) but the following: 

g.update_msg("Bye") 
// Would work, becuase the write happens inside the scope where `msg` was defined, because the inner boc `update_msg` can see the parent scope

// g.update_msg is also valid, because is invoking the boc, but is not using the `=` operator directly.

```

Generics
```js
// A Single Uppercase letter denotes a generic type
x T = "hi" // x data is generic and initialized to String

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
bi.set_data(42)

// If specific usage is needed for the generic value, then the compiler will validate, e.g. the method `len` is called, then the instantiation must have a len boc 
send #(value T) {
  length : value.len() // T has to have a boc: `len #(Int)` to be instantiated
  print("Sending `length` elements")
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
// NOTE: These varians are an exception to structural typing and behave like nominal typing whe comparing against them inside the `when` conditional checking (see below)

// e.g. With the same variables (`msg`)
Greeter : {
  Hello(msg:"Hello") // using `:` short syntax declaration
  Bye(msg:"Goodbye")
}
g Greet = Greeter.Hello()
h : Greeter.Bye() // h type is Greet

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
// For `Greeter`
Greeter #(
  Hello(msg String="Hello"), // using explicit type `=` value declaration
  Bye(msg String="Goodbye")
)

// For `Message`
Message #(
  Text(text String=""),
  Empty()
)

// Variants can have generic data too
Option : {
  Some(value T) // defines the variant with a variable `value` of type T
  None() // defines a varian with no variables
}
maybe Option(T) = Option.Some("thing")
maybe_not Option(T) = Option.None() // Safe to infer as the value won't be used 

// Can have different generic types
Result:{
  Ok(value T),
  Err(e E)
}
result Result(Int, String) = Result.Ok(1)
bad_result Result(Int,String) = Result.Err(":(" )

```

Conditional blocks of code using `match { cond => action }`
```js
// To know what kind of Greeter we have we can use `when`
// `when` is an special case where the variant can be checked making use of nominal typing whereas the rest of the language uses structural typing

// To check what kind of `Greeter` variant was used
g Greeter = get_greet()
g when
{ Greeter.Bye => "You say `g.msg`" },
{ Greeter.Hello => "and I say `g.msg`"},

// To know what kind of message we have
m : get_message()
m when
{ Message.Text => "The message content is `m.text`" },
{ Message.Empty => "There is no message" },

// Accessing generic data
maybe : Some("something")
maybe when
{ Some => print("We have `maybe.value`")},
{ None => print("We have nothing")}

// `when` conditionals can use boolean expressions too if no operand is used
when
{ a > b  => "A is greater than B"},
{ a == b =>  "A is equals to B"},
{ true => "A is lower than B"},
```


Arrays and Dictionaries

```js
[Type] // array type of type `Type` e.g
a [Int] // declares an array of Ints
[]Type // Empty array literal (value) of type Type, e.g 
a [Int] = []Int

[T] // Generic array, needs a concrete initializaton e.g. 
ga [T] = []Int // the generic array `ga` gets bound to []Int

// same is a boc that takes a string  s and returns it
same #(s String) {s}
// Array literal
[1,2,3] // type [Int]
[{ same("Hello")}, {same("Goodbye")}] // type [#(s String)]

// Array access returns Result(T, OutofBoundsError)
a : [ 1, 2, 3 ]
a[0] // Result.Ok(1)
a[5] // Result.OutOfBoundsError(5)



// Dictionary type:
[Key:Value] // is a dictionary type form Key to Value 
// eg. 
m [Int:String] // declares a dictionary of Int to Strings

// Empty dictionary literal:
[Key]Value  // e.g. 
m [Int:String] = [Int]String // m is now an empty dictionary of Int to String


[K:V] // generic dictionary type
gm [K:V] = [Int]String // the generic dictionary `gm` is bound to `[Int:String]`

// Dictionary literal of type [String:String]
["name": "Foo",
 "last_name": "Bar"] 

// Dictionary access returns Optional(V)
d : [ 1 : 2, 3: 4] // [Int: Int]
d[1] // Some(2)
d[5] // None()

// A dictionary from String to Bocs with two String values
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
wip

Control flow 
wip