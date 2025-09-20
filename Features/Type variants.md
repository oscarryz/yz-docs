#feature

>[!NOTE]
> Similar to Sum Types, but not exactly the same. 
> 
>  Most up to date design in [sumtypes/Example](../Examples/Yz%20-%20SumTypes.md)



A type can specify the different constructor variants wrapping specific data used in that variant. 

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

// To know what constructor was used call: `variable.Constructor` e.g 

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


(_Under revision to check if enums like this are possible. They should be_)
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

There are cases when a constructor is not needed,  but a variant is required, in this case a variable should be enough and it is usually declared outside of the type

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

If still want to use the constructor
```js
planets: {
  Planet #(Mercury(), Venus(), Earth(), str #(String)) {
    mass Decimal
    radius Decimal
    Mercury(mass:3.303e+23, 2.4397e6)
    Venus(4.869e+24, 6.0518e6),
    Earth(5.976e+24, 6.37814e6),
    str: {
      "`magic`: Mass= `mass`, Radius= `radius`."
    }
  }
}

something_else #() = {

  // p Planet = Planet() // compilation error, can't create etc.
  p Planet = Planets.Mercury() // Weird but possibly this is what is needed. 

}
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
        self.Ok // .Ok means: was the constructor used Ok() ?
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

So, we have the following cases: 

1. Same data types on all the variants
  1. No values
  2. Same values
  3. Different values
2. Different data types on all the variants
  1. All the variants differ
  2. Some of the variants differ


### Examples

- [Yz - SumTypes](../Examples/Yz%20-%20SumTypes.md)
- [Rust - Enums](../Rust%20-%20Enums.md)





#todo Write a single [[SumTypes]] document.