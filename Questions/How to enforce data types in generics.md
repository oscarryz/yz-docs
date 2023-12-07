We can define a data is generic by not specifying the data type and bind it on first usage, however we cannot enforce instances of the same type to be the same generic data

```javascript
Node {
	 data // generic
	 left Node //
	 right Node 
}
a = Node(1)
a.left = Node("B")
```
