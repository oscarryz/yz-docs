## Currently
<sup>Aug 7, 2025</sup>

Use `Type` `:` `{`...`}`

 ```js
// Type definition
Person : {
	name String
	address Address
}
// instantiation
p : Person("John", Address("Street", "City"))

```

## Proposed

Use `type Type {...}`

```js
type Person { 
	name String
	address Address
}
p : Person("John", Address("Street", "City"))
```


**Pros**: 
- Helps differentiate a type definition from a regular assignment, beyond the simple upper-case
**Cons:**
- Needs a new reserved word (which is not that bad)

usr.data and img.data is utilized just before the end of the code block, so the compiler/runtime would need to not wait until the end of the code block to synchonize.

What about combining lazy evaluation with async?

Async calls run and the result is pre-assigned to a variable, which can be passed around at will, and only first when it is used/needed for a computation it will block and wait to synchronize so that the result is actually available.

I think that would create the coloring problem again, because there wouldn't be a way to know when to await for the values to be used. 

```js
outer: {
   profile: enclosing()
   // should `outer`
   // wait here or exit?
}
main: {
   outer()
   // would main finish 
   // right away?
}
```