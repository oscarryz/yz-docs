
```js
id: {
	x 
	x
}
id {`a `a }
id(42)
// 42 Int
id(true)
// true Bool
id('hi')
// 'hi' String
```


```js
upgrade: {
	op {Int Option(Int) }
	x  Option(Int)
    x >>= op
}


fs: std.fs

main: {
	fs.remove_dir_all('/').or_not()
	Ok()
}
```
