We don't use commas nor semicolons but sometimes it seems to be a good idea to have a way to say: "this ends here". 

```js
a: 1 b: 2 
// vs 
a: 1; b: 2

foo: { e String print e }
// vs
foo: { e String ; print e }

// It would make a difference in generic and return parameter type though 
bar { generic_variable SpecificType } // e.g. bar {a String}
// vs 
bar { generic_variable; SpecificType} // e.g. bar { a; String }
```

#answered Optional semicolons. They would be useful to delimit things but can be omitted most of the time