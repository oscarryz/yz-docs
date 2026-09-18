---
type: example
generated: { by: "oscarryz", at: 2024-03-07T16:00:59-06:00 }
---
#example
```js
apple: 'apple'
fruits: ['apple', 'pear', 'orange']

fruits.filter({ f String; f != apple })
predicate: {
    f String
    f != apple
}
predicate('apple') // false

r: [String]()
fruits.each({ f String
  predicate(f) ? {
    r.push(f)
  }
})
Predicate: {
    test: {
        f String
        f != apple
    }
}
p: Predicate()
p() // test
p.test('orange') // true
```

