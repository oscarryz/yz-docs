---
type: example
generated: { by: "oscarryz", at: 2023-12-06T20:13:02-06:00 }
---
#example

```js
ToString: {
	to_string #(String)
}
Person: {
	name String
	to_string: {
		name
	}
}
List: {
  data [T]()
  ...
}
// # Here `ToString` is only available for instances
// # of `List` if `T` also implements `ToString`.
ListOfToStrings: {
	data [] ToString
	add: { ts ToString
		data.add(ts)
	}
}

```