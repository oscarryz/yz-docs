#example

https://inko-lang.org/


```js



// Explicit syntax
div1: {
	left Int
	right Int

	right == 0 ? {
		Result.Err('division by zero')
	}, {
	   Result.Ok(left / right)
	}
}
// Implicit syntax
div2: {
   left Int
   right Int

   right == 0 ? {
		Result.Err('Division by zero')
	}, {
	    Result.Ok(left / right)
	}
}
div3: {
   left Int
   right Int
   res: div1(left, right)
   res.is_err() ? { res }, {
	   res.is_ok_and({ v Int; v == 5 }) ? {
		   Result.Ok(50)
	   }, {
		   Result.Ok(res.get())
	   }
	}
}
main: {
	div(10, 2).or_else({ e Result.Err
		exit(e)
	})
    // we can also just "unwrap" the Ok value
    div(10, 2).get() // will exit if error
}
// result library
result: {
	Result: {
		is_ok #(Bool)
		is_ok_and #(predicate #(Bool))
		or_else #(action #(Err))
		get #(V)
	}
	Ok: {
		value // generic value
		is_ok: { true }
		is_ok_and: { predicate #(Bool); predicate() }
		or_else: { action #(Err) }
		get: { value }
	}
	Err: {
		message: ""
		is_ok: { false }
		is_ok_and: { p #(Bool); false }
		or_else: {
			action #(Err)
			action(Err(message))
		}
	}

}
```