https://rosettagit.org/tasks/99-bottles-of-beer/

```javascript

_:99.down_to(1).do({ 
    n Int          
    s: n == 1 ? {'s'} {''}
    
    print("`n` bottle`s` of beer on the wall
           `n` bottle`s` of beer
           Take one down, pass it around")
})
print('No more bottles')
```


