```js
Lion: {
  food_consumed : Set{}

  roar: {
    'roar'
  }
  eat: { food String; 
    validate(food)
    food_consumed.add(food)
  }

  get_food_consumed: {
    Array.of(food_consumed)
  }

  validate: {foo String
    reversed: food.split('').reverse().join('')
    reverse == food ? {
      error('palindromes are discusting!'
    }
  }

  get_favorite_food: {
    get_foods_consumed().filter { food String 
      food.lenght % 2 == 1 
    }
  }
}
```

No braces ( `{}` )

```js
Lion:
   consumed: Set{}
   
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

