#answered 
#rejected  Keep the current way `Name: {}` and `Name #(){}`

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

