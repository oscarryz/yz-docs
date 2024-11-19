>[!NOTE]
> This is not complete Sum type, but rather how to provide a similar functionality


A type can specify the different constructor names wrapping specific data, the used constructor for a given variable can be accessed with the form `.Contructor` which returns a `Bool`

For instance, the type `NetworkResponse` has the attributes: `data T`, `error String`, `timeout Bool`

```js
// Regular type 
NetworkResponse: {
  data T
  error String
  timeout Bool
}
```

A value can be created using the different values: 
```js
get_response #(NetworkResponse) = {
  NetworkResponse(timeout: true)
}
nr NetworkResponse = get_response()
// query nr to see how to handle the different data.. 
```

Using a constructor can help to know how data was created

```js
// Using named constructors
NetworkResponse: {
  Sucess(data T),
  Failure(error String),
  Timeout(timeout Bool),
}
nr NetworResponse = ... 
// nr still can access all the attributes:
nr.data 
nr.error
nr.timeout

// but can't create NetworkResponse directly, has to use one of the constructors
get_response #(NetworkResponse) = { 
  Success([1,2,3])
}
nr: get_response()
// to know what constructor was used call variable.Constructor e.g 

nr.Sucess ? { print("It was a sucess, the data is: `nr.data`")}
nr.Failure ? { 
    print("We did everything we could but we failed with: `nr.error`"))
}
nr.Timeout ? { wait().and_then(try_again) }



```

The type signature would contain the constructor. For the above data the signature is:
```js
#(
  Sucess(data T),
  Failure(error String),
  Timeout(timeout Bool)
  ... other methods 
)
```


These constructors are suited to group different variants of the type
Constructors can also be used with the same data, e.g. 

```js
Token: {
  Invalid,
  Eof,
  Comment,
  // Identifier oand basic type literals
  Identifier(name String),
  Int,
  Float,
  Img,
  Char,
  String,
  // ... etc
}
a Token.Int()
```

There are cases when a constructor is not needed but a variant, in this case a variable should be enough and it is usually declared outside of the type

```js
planets: {

  Planet: {
    mass Decimal
    radius Decimal
  }
    // Upper case just to denote it behaves like a constant
    // Can still be modified though 
    MERCURY :Planet(3.303e+23, 2.4397e6),
    VENUS   :Planet(4.869e+24, 6.0518e6),
    EARTH   :Planet(5.976e+24, 6.37814e6),
    MARS    :Planet(6.421e+23, 3.3972e6),
    JUPITER :Planet(1.9e+27,   7.1492e7),
    SATURN  :Planet(5.688e+26, 6.0268e7),
    URANUS  :Planet(8.686e+25, 2.5559e7),
    NEPTUNE :Planet(1.024e+26, 2.4746e7),
    PLUTO   :Planet(1.27e+22,  1.137e6)
}
// To make `mass` and `radius` private, create an empty signature for Plant
planets: {

  MERCURY Planet
  VENUS Planet
  Planet #() = {
    mass Decimal
    radius Decimal
    // and initialize the outer variables
    MERCURY = Planet(3.303e+23, 2.4397e6)
    VENUS = Planet(4.869e+24, 6.0518e6),
    
  }
}
planet Planet = planets.VENUS
```


How to allow a type to specify its variants while keeping the structural typing coherent


```js
std: {
  result: {
    Result: {
      T, E,
      Ok(data T),
      Err(err E),
      is_ok: { 
        .Ok // .Ok means: was the constructor used Ok() ?
      }
      get: { data }
      cause: { err }
    }
  }
  // Signature
  // Result #(T, E, data T, err E, is_ok #(Bool), get #(T), cause #(E)) 
  Result #(T, E, Ok(data T), Err(err E), is_ok #(Bool), get #(T), cause #(E)) 
}

// result signature
  result #(
    Result,
    Result.Ok,
    Result.Err,
 )
```


*Previous information below, might be outdated*

Answer: The base type would define all the possible value and the subtypes would implement them to be considered subtypes. 

e.g. `Result`, `Ok` and `Error`

```javascript
std: {
  result: {
    // Define the base type
    Result: {
      T 
      E
      and_then #(fn #(Result(U,E))) 
      is_ok #(Bool)
      is_err #(Bool)
      get #(V)
      cause #(E)
    }
    // The subtype would have to implement all the methods
    Ok: {
      T 
      E
      v T
      e E
      and_then : {
        fn #(T, Result(U,E))
        fn()
      }
      is_ok : { true }
      is_err: { false }
      get: { v }
      cause: { e }
    }
    // Err would immpelment that too
    Err: {
      v T
      e E
      and_then: {
        self // ??
      }
      is_ok: { false }
      is_err: {true}
      get: {v}
      cause:{e}
      
    }
  }
}
// The we can use either Ok or Err where Result is expected
fetch_data #(id String, Result(String,String)) 
fetch_data("123").and_then({
  data String
  print("Data has been retrieved: `data`")
  Ok(data) 
})
d Result(String,String) = Err("", "Unable to fetch stuff")

```


To define default values a sum type can define the absence of value, it would need to rely on another sum type though. 


```js
Node: {
	v T // Generic value
	next Option(Node(T)) // `next` is an optional node of T
}
root Node(Int) = Node(1, Some(Node(2, Some(Node(3, None())))))


```

This approach doesn't need pattern matching, although it can be added later. 
The drawback is it needs to define all the methods and variables in the upper type and those need to be repeated below even if not needed. 