How to allow a type to specify its variants while keeping the structural typing coherent

Answer: The base type would define all the possible value and the subtypes would implement them to be considered subtypes. 


e.g. `Result`, `Ok` and `Error`

```js
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