#example

https://inko-lang.org/


```js



// Explicit syntax
div1: {
	left Int
	right Int

	(right == 0).ifTrue(then:{
		result.Error('divison by zero')
	}
	else: {
	   result.Ok( left / right)
	})
}
// Implicit syntax, no parenthesis and `.` ommited on `?`
div2: {
   left Int
   right Int

   	right == 0  ? {
		result.Error 'Division by zero'
	}
	{
	    result.Ok left / right
	}
}
div3: {
   left Int
   right Int
   res: div1 left right
   res.is_err() ? { res } {
	   res.is_ok_and {v Int v == 5} ?  {
		   result.Ok(50)
	   } {
		   result.Ok(res.get())
	   }
	}
}
main: {
	div(10 2).or_else { e result.Err
		exit(e)
	}
    // we can also just "unwrap" the Ok value
    div(10 2 ).get() // will exit if error
}
// result library
result: {
	Result {
		is_ok { Bool }
		is_ok_and { predicate {Bool}}
		or_else { action {Err}}
		get {v}
	}
	Ok {
		value // generic value
		is_ok: {true}
		is_ok_and: {predicate {Bool} predicate()}
		or_else: {action {Err} }
		get: {value}
	}
	Err {
		message: ""
		is_ok: {false}
		is_ok_and: {p {Bool} false}
		or_else: {
			action {Err}
			action Err(message)
		}
	}

}
```
