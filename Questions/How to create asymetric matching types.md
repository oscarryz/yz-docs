
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



