#example
```js
apple: 'apple'
fruits: ['apple' 'pear' 'orange']

fruits.filter { f String ; f != apple } 
let predicate = {
    f String
    f != apple 
}
predicate('apple') // false 

r: []String
fruits.each({ f String 
  if predicate(f) {
    r.push(f)
  }
}
Predicate { 
    test: {
        f String
        f != apple
    }
}
p : Predicate()
p() // test 
p.test('orange') // true 
```

