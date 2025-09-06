
[Go - Pipeline 3](Examples/Go%20-%20Pipeline%203.md)

https://github.com/adonovan/gopl.io/blob/master/ch8/pipeline3/main.go

```js
// Version with a parameter to chain the pipeline

// Counter takes a transformer ()
counter #(transformer #(Int)) {
	1.to(100).do({
		i Int
		transformer(i)
	})
}
create_squarer #(consumer #(Int)) {
	squarer #(i Int) {
		consumer( i * i)
	}
	squarer
}
printer #(n Int) {
	print(n)
}
main #() {
	counter(create_squarer(printer))
}
```

```js
// no signature version
counter : {
	consumer #(Int)
	1.to(100).do({
		i Int
		consumer(i)
	})
}
squarer : { 
	consumer #(Int)
	#(n Int) {
		consumer(n * n)
	}
}
printer #(n Int) {
	print("`n`")
}
main: {
	counter(squarer(printer))
}
```

```js
// Direct version relying on closures, so no 
// function parameters are needed. 
counter #(Int)
squarer #(Int, Int)
printer #(Int)

counter = { 
  1.to(100).do({
	  i Int
     squarer(i)
  })
}
squarer = {
	i Int
	print(i * i)
}
printer = {
	print(i)
}
```

