#example

```js
Lion: {
  food_consumed: Set()

  roar: {
    'roar'
  }
  eat: { food String
    validate(food)
    food_consumed.add(food)
  }

  get_food_consumed: {
    Array.of(food_consumed)
  }

  validate: { food String
    reversed: food.split('').reverse().join('')
    reversed == food ? {
      error('palindromes are disgusting!')
    }, { }
  }

  get_favorite_food: {
    get_food_consumed().filter({ food String
      food.length() % 2 == 1
    })
  }
}
```

No braces ( `{}` ) — exploration of alternative syntax (not standard Yz):

```
Lion:
   consumed: Set()

   roar:
     'roar'

   eat:
      food String
      validate food
      consumed.add food

    validate:
       food String
       reversed: food.split '' .reverse().join ''
       reversed == food ?
           error "I don't like palindromes"

    get_favorite_food:
       consumed.filter
          item String
          item.length() % 2 == 1

```
