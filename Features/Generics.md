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
