
#answered  See [SumTypes](SumTypes.md)

Original discussion below 

I called them asymmetric matching types, but they are actually Sum types 
Also, how to create sum types (I don't think that's the same).

See for example: [Pomar - http request](Pomar%20-%20http%20request.md)

It needs a type with two possible states (similar to `Option { Some, None }`) 


How can we create a single type that can have two variants?


```js
Status {
}
Ok {}
Error {
	msg String
}
status Status = either_ok_or_error()
```

#### Possible solution
Don't rely on different types, use methods instead to tell the difference, this is what Java does
```js
optional: {
    
    Optional {
      T
      is_some {Bool}
    }
    Some {
        T
        data T
        is_some: { true }
    }
    None {
        is_some: { false }
    }
}
User:{}
...
x Optional(User) // x Optional<User>
x = select(User).where { u User ; u.name == "Alice"}
// could be Some(user)
// or None()

// Not instantiated generic type
y Optional(T)

// by the way, `select` here is generic
select: { 
    T 
    // returns an instance of `Where` which is generic 
    // on T
	Where(T) 
}
// `Where` has a generic method `where`
Where : {
	e E
	// it take a `clause` which is a block with a 
	// generic element `e` and returns a boolean
	where: {
		clause { e E ; Bool }
		clause(e) 
	}
}
```

So for status:
```js
http: {
   Status # (
       is_ok #(Bool)
   )
   Ok: {
       is_ok : {true}
   }
   ClienError: {
      msg String
      is_ok: { false }
    }
}
...
status http.Status = ok_err()
if status.is_ok {
   print("all good")
} {
   print("problems")
}

```


Ordering example: 

```rust
// rust
pub enum Ordering {
  Less = -1, 
  Equal = 0,
  Greater = 1,
}
// Yz
Ordering : Int 
Less : Ordering
Equal : Ordering
Greater : Ordering

compare #(Int, Int, Ordering) = { 
  a Int
  b Int
  a > b ? { return Greater() }
  a < b ? { return Less() }
  Equal
}
```

Another option is to create an inner type like this [Kotlin example](https://elizarov.medium.com/kotlin-and-exceptions-8062f589d07#:~:text=sealed%20class%20ParsedDate,errorOffset)

```js
Status: {
  Ok : {}
  Err: {msg String}
}
fetch #(Status) = {
  ... 
  Status.Ok()
  Status.Error()
}

```

But that doesn't go well with the language, instead a boc `status` with three types could be used
```js
status: {
  Status: {}
  Ok : Status
  Err: Status 
}
```
But the top level (`Status`) needs to hold all the data the others might need, like `message`

```js
status: {
  Status : {
    msg String = ""
    is_ok : { true }
    
  }
  Ok : Status
  Err : Status
  Err.is_ok = { false }
}
fetch #(Status) = { 
  ... 
  status.Ok()
  status.Err("No data available")
}
status Status = fetch()
status.is_ok ? { 
  "etc.etc"
}
```


If we have a macro system, the "inherited" type can be "derived"

```js
Status: {}

'derive: [Status]'
Err: {
  msg String
}

'derive: [Status]'
Ok: {
}
```


The problem right now, it is weird to have a type inside another type, because it's not clear (yet) how to create a new instance:
```js
Status: {
  Ok
  Error(reason String)
}
// is it:
s1 Status = Status.Ok()
// or 
s2 Status = Ok()

```
