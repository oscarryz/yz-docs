How to allow a type to specify its variants while keeping the structural typing coherent


Ideal case, the type specify its generics and the variants it supports

```js
Option: {
	// Data held
	value V 
	error E
	// variants
	Some(V)
	None(E)
	// methods below... 
	or_else: {
		replacement (V)
		match Self {
			Some -> value
			None -> replacement()
		}
	}
}

```

But that would require patter matching to match on the type. Also It's kinda weird (albeit cool) how new types are inside the type. All the methods would require to test if is A or B to do something whereas I better impl would be to let the impl do 

```js
option : {
	// Option (or_else (r (V)))
	Option ( 
		or_else (replacement (V))
	)
	// Some ( V; or_else (r (V)))
	Some : {
		value V	
		or_else : {
			replacement (V)
			value
		}
	}
	None: {
		error E
		or_else: {
			replacement (V)
			replacement()
		}
	}
}
main: {
	User: {}
	find_user (String; Option(User)) = {
		 ... 
		 name == 'test' ? {
			 Some(User("test"))
		 } {
			 None('User is not test')
		 }
	}
	uo : find_user('test').or_else {
		User('Not test')
	} 
}
```

I