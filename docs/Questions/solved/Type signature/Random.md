---
type: example
generated: { by: "oscarryz", at: 2024-03-21T13:44:01-05:00 }
---
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