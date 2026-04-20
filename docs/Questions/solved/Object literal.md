

#answered  Yz will use `()` notation to create instance e.g. `Person()`
#rejected

We could have a notation for object literal such that the Boc defined with
```js
Person: {
  name : ""
  age  : 0
}
```

Could be instantiated with: 

```js
p : Person {
  name: "Alice"
  age : 42
}
```

But it seems this look like the declaration way to much and would make the code harder to read (the only difference being the colon `:`  between the type name and the declaration in the first, and the lack of it during the initialization)

Thus the only way to create instance is with the `()` notation: 

```js
p : Person(
  name : "Alice"
  age  : 42
)
```

Furthermore, the object literal would need a `=` in the fields to denote the variable is not being declared but assigned

```js
p : Person {
  name = "Alice"
  age = 42
}
```

Creating a new variant of `{}` thing

Finally Yz having structural typing so we can use straight block literal instead

```js
p : {
  name: "Alice"
  age : 42
}

q Person = p // granted
```

