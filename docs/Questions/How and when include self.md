#open-question 

Some data types need self, but this is not a keyword


When [Compile Time Bocs](docs/Features/Compile%20Time%20Bocs.md) [Info strings](docs/Features/Info%20strings.md) and [Native Type Annotations](docs/Questions/Native%20Type%20Annotations.md) is completed it can be used to generate the self type

e.g. 

```js
`!:[Native, Derive, Self]
self: "self"
`
Person: {
   name String 
   to_str #(String) {
	   self.name // self created by the "macro"
   }
}
```

This depends on [Compile time bocs Interface interaction design](docs/Questions/Compile%20time%20bocs%20Interface%20interaction%20design.md)
