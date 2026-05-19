#example

```js

test: std.test
main: {
   tests: test.new()
   tests.test('Adding two integers', {
	   t test.Test
	   t.equal(10 + 5, 15)
	   t.equal(1 + -1, 0)
   })
   tests.run()
}
```


Using annotations?

```js

test: std.test
assert: std.test.assert

`!:[Test]`
main_test: {
	`test-name: "Should work"`
	exploration: {
		result: 2 + 2
		assert.equal(result, 4)
		assert.gte(result, 2)
	}
}
```

Using annotations?

```js

test: std.test
assert: std.test.assert

`!:[Test]`
main_test: {
	`test-name: "Should work"`
	exploration: {
		result: 2 + 2
		assert.equal(result, 4)
		assert.gte(result, 2)
	}
}
```


Using test docs

```js
`
!:[Test]
test-desc: "Some descripton"
test: {
	result: sum(2, 2)
    assert: result == 4
}
`
sum: { a Int; b Int; a + b }
```