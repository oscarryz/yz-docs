#wip

When a value type is omitted, the variable is generic and will bind the the first usage.

```javascript
something: {
	data
	data.len() // data is generic but needs a `len()` method
}
something(1) // comp err, Int doesn't have a `len()` method
something("Hi") // ok
something([1 2 3) // ok.. we wanted soemthing with `len()` after all.

```

Although it might work, this is still an open question in: [[How to enforce data types in generics]]


```js
none: std.option.none
tree: {
	Node {
	   data
	   left Node = none()
	   right Node = none()
	}
	builder: {
		{
			data // generic data
			Node(data)
		}
	}
}
node: tree.builder()
root: node(1) // the builder `node` is now bound to `Int`
root.left = node(2)
root.right = node(3)
node_str: tree.builder()
a : node_str('hi')
a.left = node_str('Bye')
// this would be still valid however
a.right = node(1) // :(


```

#open-question 