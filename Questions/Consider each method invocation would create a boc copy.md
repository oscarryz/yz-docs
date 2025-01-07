
Currently every time a block is executed, the state is kept between calls: 

```js
foo: {
  a Int
  a = a + 1
}
foo.a = 1
foo() // a : 2
foo() // a : 3
foo() // a : 4

```

Consider an alternative where each invocation start from a clean slate (just like a regular function does in 99.9% of the programming languages)

```js
foo: {
  a Int
  a = a + 1
}
foo.a = 1
foo() // a : 2
foo() // a : 2
foo() // a : 2
```

We need to find a way to reset the execution but preserving the fact that blocks are stateful objects. 

_The example above is wrong, but the question remains..._
