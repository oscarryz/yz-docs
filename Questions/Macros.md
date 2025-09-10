#open-question
#v2

- Regular Yz code
- Creates AST nodes at compile time

```js
macro! say_hello {
	print("Hello!")
}
main: {
	say_hello!()
}
macro! create_block { 
	$func_name Ident
	$func_name {
		println!("You called $(stringify!($func_name))()")
	}
}
create_block!(foo)
create_block!(bar)

macro! def { 
	v []Ident
	t Type
	e Expr
	{
		v t = e
	}
}
def x y Int 0 
```

#open-question 