
```js
ToString {
	to_string {String}
}
Person  {
	name String
	to_string: {
		name 
	}
}
List {
  data []
  ... 
}
//# Here `ToString` is only available for instances
// # of `List` if `T` also implements `ToString`.
ListOfToStrings {
	data [] ToString
	to_string: {
	   ... 
	}
}
// hmmm 
  
```