How to allow a type to specify its variants while keeping the structural typing coherent


Ideal case, the type specify its generics and the variants it supports

```js
Option: {
	// Data held
	value V 
	// variants
	Some: {value V}
	None : {} 
	// methods below... 
    self Option
	or_else: {
		replacement #(V)
		match info(self).type [
			{Some}:{ value }
			{None}: {replacement()}
		]
	}
}
// use example
maybe_something #(Option(String)) = { 
  if random.next_int(10) [
    { 6 } : { Some("Bingo") }
    { _ } : { None }
  ]
}

optional_string : maybe_something()
certainly_a_string : optional_string.or_else { "Nothing" }
// problem, we need to setup `self` :/ 

```

But that would require patter matching and a way to set `self` to match on the type. 
Also It's kinda weird (albeit cool) how new types are inside the type. All the methods would require to test if is A or B to do something whereas I better impl would be to let the impl do 

```js
option : {
	// Option #(or_else #(r #(V)))
	Option #( 
		or_else (replacement #(V))
	)
	// Some ( V; or_else (r (V)))
	Some : {
		value V	
		or_else : {
			replacement #(V)
			value
		}
	}
	None: {
		error E
		or_else: {
			replacement #(V)
			replacement()
		}
	}
}
main: {
	User: {}
	find_user #(String, Option(User)) = {
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

The caveat with is we would have to implement all the methods and 