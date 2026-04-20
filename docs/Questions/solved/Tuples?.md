When a block has no variables but has values, is that a tuple?

```javascript
a: {1 2 3} // type is `{Int Int Int}`
// same as above
a {Int Int Int} = {1 2 3}
```

Is this a tuple? We said in [Structural typing](../../Features/Structural%20typing.md) we can access their values by index 

```js
a: {1 2 3}
a.0 // 1
a.1 // 2
a.2 // 3
```

As args: 

```javascript
b: {
   a {Int Int Int}
}
b {1 2 3}
```
#answered  Yes, these blocks are tuples. 