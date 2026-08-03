#open-question

~~Superseded~~ — see [Self keyword](Self%20keyword.md) (YZC-0060). The macro-generated approach sketched below was considered and rejected: macros run at compile time over a type declaration, before any instance exists, so they can't produce a value that points back at a specific instance. `self` is proposed as a compiler built-in resolved lexically instead.

---

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
