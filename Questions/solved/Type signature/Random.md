#example 

```js
Foo (
	w Int,
	abc Int,
	bar #(v Int, String),
	bar2 #(a Int, b Int)
) = {
	x Int
	y Int
	z Int
	bar : {
		v Int
		"value"
	}	
	bar2 : {
		a Int 
		b Int
	}
}
```