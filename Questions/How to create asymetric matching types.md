
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
Don't rely on different types, use methods instead
```js

optional: {
    
    Optional {
      <T>
      is_some {Bool}
    }
    Some {
        <T>
        data T
        is_some: { true }
    }
    None {
        is_some: { false }
    }
}
User{}
...
x Optional(User)
x = select(User()).where { u User ; u.name == "Alice"}
// could be Some(user)
// or None()
```

So for status:
```js
http:{
   Status{
       is_ok {Bool}
   }
   Ok {
       is_ok : {true}
   }
   ClienError {
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
