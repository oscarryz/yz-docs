#open-question 

Some data types need self, but this is not a keyword


When [Macros](docs/Features/Macros.md) and [GoExtensions](docs/Features/GoExtensions.md) is completed it can be used to generate the self type

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
