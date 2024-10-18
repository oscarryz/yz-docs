Revision of: [Pattern matching](Features/Pattern%20matching.md)

```js
status : {
  Status #()
  InProcess
  Completed #(date Date)
  Rejected #(cause String)
}
Order #(
  description String
  status Status
)
...
order Order
switch order { 
   case { desc , InProcess } :  { print("Your `desc` is almost ready") }
   case { _ , Completed(date) }: { print("Your order is has been fufilled at: `date`") }
   case { _, Rejected(cause) } :  { print("Your order was rejected because: `cause`") }
}
```


1. Extracting 
  1. elements from tuples
  2. elements from lists
  3. elements from dictionaries
2. Conditional based on object values
  1. Checking for type
  2. Matching multiple patterns
3. Destructuring
  1. Assigning value from tuple
  2. Destructuring nested eelements
4. Guard claouses and defaults
  1. Using guard 
  2. Handling defaults

// Exploratory syntax

1.1
```js
data: {1; "hello"; 3.14}
list: [1, 2, 3]

match data [
  {x, _, z}: {print("x: `x`, z: `z`")}
]
match list [
  {x, _, z} : {print("x: `x`")}
]
```

1.2
```js
dict: ["name": "Alice" "age": 30 "city": "New York"]
match dict [
  {"name": name, "age": age, "city": city} : print("Name: `name`, age: `age`, city: `city`")
]

```

2.1

```js
value: 42
match value [
  {Int} : print("value is an integer")
  {String}: print("value is a string")
  {_}  : print("values is something else")
]
```

2.2

```js
x: 10
match x {
  {1|2|3}: print("x is 1,2 or 3")
}

```

3.1
```js
point : {3; 5}
match point [
  {x, y}: print("x: `x`, y: `y`")
]

```
3.2

```js
person = {"name": "Bob", "address": {"street": "Main St", "number": 123}}
Person : {
  name String
  address Address
}
Address: {
  street String
  number Int
}
p: Person("Alice", Address("Main St", 456))

match person [
  {"name": name, "address": {"street": street}} : print... 
]
match p [
  {Person(name, Address(street, _))}: print("`name` lives in `street`)
]
```

4.1
```js
value : 10 

match info(value).type, value [
  { Int, value > 0  } :  {print("positive integer")}
  { Int } : {print "Non positive integer" }
  { _ }  : {print "Not an integer" }
 ]
```

4.2

```js
data: "abc"
match data [
  { x, y, z; x == y } : {print("First two elements are equal")}
  { _ } : { print("Data doesn't match the pattern")} 
]
```


Current map of functions

1.1
```js
data: {1; "hello"; 3.14}
list: [1, 2, 3]

match data [
  {x, _, z}: {print("x: `x`, z: `z`")}
]
match list [
  {x, _, z} : {print("x: `x`")}
]
dict: ["name": "Alice" "age": 30 "city": "New York"]
match dict [
  {"name": name, "age": age, "city": city} : print("Name: `name`, age: `age`, city: `city`")
]

value: 42
match value [
  {Int} : print("value is an integer")
  {String}: print("value is a string")
  {_}  : print("values is something else")
]
x: 10
match x {
  {1|2|3}: print("x is 1,2 or 3")
}

point : {3; 5}
match point [
  {x, y}: print("x: `x`, y: `y`")
]

person = {"name": "Bob", "address": {"street": "Main St", "number": 123}}
Person : {
  name String
  address Address
}
Address: {
  street String
  number Int
}
p: Person("Alice", Address("Main St", 456))

match person [
  {"name": name, "address": {"street": street}} : print... 
]
match p [
  {Person(name, Address(street, _))}: print("`name` lives in `street`)
]

match info(value).type, value [
  { Int, value > 0  } :  {print("positive integer")}
  { Int } : {print "Non positive integer" }
  { _ }  : {print "Not an integer" }
 ]
data: "abc"
match data [
  { x, y, z; x == y } : {print("First two elements are equal")}
  { _ } : { print("Data doesn't match the pattern")} 
]
```

With special keyword `match`
```js
match(items) {
  { cond } : { action }
}


data: {1; "hello"; 3.14}
list: [1, 2, 3]

match (data) {
  {x, _, z}: {print("x: `x`, z: `z`")}
}
match (list) {
  {x, _, z} : {print("x: `x`")}
}
dict: ["name": "Alice" "age": 30 "city": "New York"]
match dict {}
  {"name": name, "age": age, "city": city} : print("Name: `name`, age: `age`, city: `city`")
}

value: 42
match (value) {}
  {Int} : print("value is an integer")
  {String}: print("value is a string")
  {_}  : print("values is something else")
}
x: 10
match( x ){
  {1|2|3}: print("x is 1,2 or 3")
}

point : {3; 5}
match (point) {}
  {x, y}: print("x: `x`, y: `y`")
}

person = {"name": "Bob", "address": {"street": "Main St", "number": 123}}
Person : {
  name String
  address Address
}
Address: {
  street String
  number Int
}
p: Person("Alice", Address("Main St", 456))

match( person) {}
  {"name": name, "address": {"street": street}} : print... 
}
match (p ){}
  {Person(name, Address(street, _))}: print("`name` lives in `street`)
}

match (info(value).type, value) {}
  { Int, value > 0  } :  {print("positive integer")}
  { Int } : {print "Non positive integer" }
  { _ }  : {print "Not an integer" }
 }
data: "abc"
match (data) {}
  { x, y, z; x == y } : {print("First two elements are equal")}
  { _ } : { print("Data doesn't match the pattern")} 
}
```

With special keyword `match`  and "free range" patterns
```js
match(items) {
   cond  : { action }
}


data: {1; "hello"; 3.14}
list: [1, 2, 3]

match (data) {
  x, _, z: {print("x: `x`, z: `z`")}
}
match (list) {
  x, _, z : {print("x: `x`")}
}
dict: ["name": "Alice" "age": 30 "city": "New York"]
match(dict) {
  "name": name, "age": age, "city": city : {
    print("Name: `name`, age: `age`, city: `city`")
  }
}

value: 42
match (value) {
  Int : {print("value is an integer")}
  String: {print("value is a string")}
  _  : {print("values is something else")}
}

x: 10
match( x ){
  1|2|3: {print("x is 1,2 or 3")}
}

point : {3; 5}
match (point) {
  x, y: {print("x: `x`, y: `y`")}
}

person = {"name": "Bob", "address": {"street": "Main St", "number": 123}}
Person : {
  name String
  address Address
}
Address: {
  street String
  number Int
}
p: Person("Alice", Address("Main St", 456))

match(person) {
  "name": name, "address": {"street": street} : {
    print("`name` lives in `street`")
  }
}
match (p ){}
  Person(name, Address(street, _)): {
    print("`name` lives in `street`")
  }
}

match (info(value).type, value) {}
   Int, value > 0   :  {print("positive integer")}
   Int, _ : {print "Non positive integer" }
   _   : {print "Not an integer" }
 }
 
data: "abc"
match (data) {
   x, y, z, x == y  : {print("First two elements are equal")}
   _ : { print("Data doesn't match the pattern")} 
}
```

Using `if` as keyword

```js
if(items) {
   cond  : { action }
}


data: {1; "hello"; 3.14}
list: [1, 2, 3]

if (data) {
  x, _, z: {print("x: `x`, z: `z`")}
}
if (list) {
  x, _, z : {print("x: `x`")}
}
dict: ["name": "Alice" "age": 30 "city": "New York"]
if(dict) {
  "name": name, "age": age, "city": city : {
    print("Name: `name`, age: `age`, city: `city`")
  }
}

value: 42
if (value) {
  Int : {print("value is an integer")}
  String: {print("value is a string")}
  _  : {print("values is something else")}
}

x: 10
if( x ) {
  1|2|3: {print("x is 1,2 or 3")}
}

point : {3; 5}
if (point) {
  x, y: {print("x: `x`, y: `y`")}
}

person = {"name": "Bob", "address": {"street": "Main St", "number": 123}}
Person : {
  name String
  address Address
}
Address: {
  street String
  number Int
}
p: Person("Alice", Address("Main St", 456))

if(person) {
  "name": name, "address": {"street": street} : {
    print("`name` lives in `street`")
  }
}
if (p ){
  Person(name, Address(street, _)): {
    print("`name` lives in `street`")
  }
}

if (info(value).type, value) {
   Int, value > 0   :  {print("positive integer")}
   Int, _ : {print "Non positive integer" }
   _   : {print "Not an integer" }
 }
 
data: "abc"
if (data) {
   x, y, z, x == y  : {print("First two elements are equal")}
   _ : { print("Data doesn't match the pattern")} 
}
```

// UCS

```js
if x is a and y is b
then a + b 
else 
0

if x is Some(y) and y is Some(z)

if x is Some(y) and f(y) is Some(z)

if let size = size(args) {
  0 : "null"
  abs(size) : 
    100.< : "large"
    10.> : "small"
      _   : "medium"
}

if x is 
  Right(None) : default_value
  Right(Some(cached)) : f(cached)
  Left(Input) && compute(input) is
    None : default_value
    Some(result): f(result)

when x [
  { Right(None) } : {default_value}
  { Right(Some(cached))} : {f(cached)}
  { _ }: { when input [
    {None} : {default_value}
    {Some(result)}: {f(result)}
  ]}
]


if (name) { 
  .starts_with("_") && name.tail_option == Some(".txt")
  && name_postfix.to_int_option == Some(0)
  && 0 <= 0 and 0 < arity 
}
```


Back to map is things
```js

data: {1; "hello"; 3.14}
list: [1, 2, 3]

match_is_some data [
  {
    data.0
    data.1
  } : {print("x: `data.0`, z: `data.1`")}
]


match_is_some list [
  {
    list[0]
    list[1]
  } : {print("x: `list[0]`, z: `list[1]`")}
]
dict: ["name": "Alice" "age": 30 "city": "New York"]
match_is_some dict [
  { 
    dict["name"]
    dict["age"] 
    dict["city"]  
  } :  print("Name: `dict["name"]`, age: `dict["age"]`, city: `dict["city"]`")
]

value: 42
match_eq type_of(value) [
  {Int} : print("value is an integer")
  {String}: print("value is a string")
  {_}  : print("values is something else")
]
x: 10
match_either x {
  { 1; 2; 3 }: print("x is 1,2 or 3")
}

point : {3; 5}
match_is_some point [
  {point.0; point.1}: print("x: `point.0`, y: `point.1`")
]

person = {"name": "Bob", "address": {"street": "Main St", "number": 123}}
Person : {
  name String
  address Address
}
Address: {
  street String
  number Int
}
p: Person("Alice", Address("Main St", 456))

match_eq info(person).type [
  {
    #(name String, address #(stree String))
  }:  print("`person.name` lives in `person.street`)
]
match_like p [
  {Person(some(String), Address(street, "" ))}: 
  print("`p.name` lives in `p.street`)
]

match value [
  { info(value).type == Int && value > 0  } :  {print("positive integer")}
  { info(value).type == Int } : {print "Non positive integer" }
  { true }  : {print "Not an integer" }
 ]
data: "abc"
match data[
  { data.at(0) == data.at(1) } : {print("First two elements are equal")}
  { _ } : { print("Data doesn't match the pattern")} 
]



 ```

Using a `Matcher` library similar to [Hamcrest](https://hamcrest.org/)


```js

data: {1; "hello"; 3.14}
list: [1, 2, 3]

is: has:  std.hamcrest.is
match: std.hamcrest.match
else : std.hamcrest.else

match(data).conditions([
  has.element_at([0,2]) : {print("x: `data.0`, z: `data.1`")}
  else : {""}
])


match(list).conditions([
  has.indexes([0,2]) : {print("x: `list[0]`, z: `list[2]`")}
  else : {""}
])

dict: ["name": "Alice" "age": 30 "city": "New York"]
match(dict).conditions([
  has.keys(["name", "age", "city"]) :  {
    print("Name: `dict["name"]`, age: `dict["age"]`, city: `dict["city"]`")
  }
  else : {""}
]

value: 42
match(value).conditions([
  is.of_type(Int): { 
    print("value is an integer")
  }
  is.of_type(String): {
    print("value is a string")
  }
  else: {print("values is something else")}
])

x: 10
match(x).conditions([
  is.either([1,2,3]) : { print("x is 1,2 or 3") }
  else: { print("x is neither 1,2 nor 3" )}
])

point : {3; 5}
match(point).conditions([
  has.elements_at([0,1]): {print("x `point.0`, y: `point.1`")}
  else: {""}
])

person : {name: "Bob"; address: {street: "Main St", number: 123}}
Person : {
  name String
  address Address
}
Address: {
  street String
  number Int
}
p: Person("Alice", Address("Main St", 456))

match(person).conditions([
  {person.name == "Bob"}: {print("Hello Bob")}
  else: {""}
])

match(person).conditions([
  has.type_of(#(name String, address #(street String))): {
    print("`person.name` lives in `person.address.street`)
  }
  else: {""}
]

match_like p [
  {Person(some(String), Address(street, "" ))}: 
  print("`p.name` lives in `p.street`)
]

value: 42
match(value).conditions([
  is.of_type(Int).and({value > 0 }): { print ("positive integer" )}
  is.of_type(Int) : { print ("non positive integer" )}
  else : {print "Not an integer" }
])
data: "abc"
match(data).conditions([
  { data.at(0) == data.at(1) } : {print("First two elements are equal")}
  is.same_at([0,1]) : {print("First two elements are equal")}
  else : { print("Data doesn't match the pattern")}
])

is: match: std.hamcrest.is
fib: {
  match(n).conditions([
    is.either([0,1]) : {return n}
    else : { return fib( n - 1) + fib (n - 2) }
  ])
}
```


##  std.hamcrest  library

The hamcrest library contains matchers that can be used in a dictionary of predicates/actions to perform program logic.


e.g.  an object `is` would be set the data subject, that its methods would evaluate it

```js
// `is` partial definition
match
is: {
  data T
  // Takes a dictionary of predicate to actions
  // if a predicate matches, the action is invoked
  conditions #(p_a [#(Bool):#(V)])
  either #(list [T], #(Bool))
}
// Usage
x: 42
is(x) // setup the subject in `is.data`
is.conditions([
  {x == 42 || {x == 4}} : { 100 }
  { true } : { 1000 }
])
```
In the example above, the method `is.either` can be used to return the block that evaluates if x is either 42 or 4. 
Aliases can be setup to improve readability: 
```js
is: std.hamcrest.is
x: 42
match(42).conditions([
  is.either([42, 4]) : {100}
  else : { 1000 }
])
```





```js
hamcrest {
  else #(#(Bool)) = {
    { true }
  }
  data T
  match #(d T) = {
     data = d 
      is
  }
  is: {
    data T
    conditions #(predicates_and_actions [#(Bool):#(V)], V) = {
        predicates_and_actions [#(Bool):#(V)]
        predicates_and_actions.for_each({
          predicate #(Bool)
          action #(V)
          predicate() ?  
            return action()
          })
        })
        panic("Should've returned something")
    }
    either #(list [T], #(Bool)) = {
      list.each { 
        e T 
        e == data ? { return { true } }
      } 
      { false }
    }
  }
}

```