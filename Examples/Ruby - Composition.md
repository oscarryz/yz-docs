
https://www.reddit.com/r/ProgrammingLanguages/comments/1oe31os/should_object_fields_be_protected_or_private/


```js
Point: {
	x Int
	y Int
	to_string #() {
		"(`x`, `y`)"
	}
}
MovablePoint : {
	mix Point
	move_right #() {
		x = x + 1
	}
}

mp : MovablePoint(3,4)
println("`mp`")
mp.move_rigth()
puts(mp)

... 
// After `mix Point` the structure looks like this: 
MovablePoint : {
	x Int
	y Int
	to_string #() {
		"(`x`, `y`)"
	}
	...
}
```

```js
Point #(x Int, y Int) {
   z Int // not exposed 
}
MovablePoint: {
   mix Point 
   move #() {
      z = 0 // err, no z declared
   }
}

```
