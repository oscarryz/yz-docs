There's no enums, but following the convention instances of a given type can be created and use them as enumeration types


eg. 


```js
veggies: {
    Veggie : {}
    SQUASH: Veggie{}    
    CABBAGE: Veggie{}    
    BROCCOLI: Veggie{}    
}
... 
Veggie: veggies.Veggie

eat: { 
    what Veggie
    when_eq what [
        {veggies.SQUASH}:{ print 'eating squash'}
        {veggies.CABBAGE}:{ print 'eating cabbage'}
        {veggies.BROCCOLI}:{ print 'eating brocoli'}
    ]   
}
veggies: {
    Veggie: {
        name String
        carbs Int
    }
    SQUASH:Veggie{name:'squash'; carbs:1}
    // also
    //SQUASH:{name:'squash'; carbs:1} // due to structural typing
    CABBAGE:Veggie{name:'cabbage'; carbs:3}
    BROCCOLI:Veggie{name:'broccoli';carbs:0}
}


```

If used without any other information a new type can be created and aliased to another: 
```js
veggies: { 
  Squash: 
  Cabbage: 
  Broccoli: 
  Veggie: {}
}
random_salad #(Veggie) = { 
  Broccoli()
}
v Veggie = random_salad()
```

If different values are needed, then initialize them 
```js
Planet: {
  mass Decimal
  redius Decimal
  G = 6.67300E-11;
    

  surfaceGravity #(Decimal) = {
     G * mass / (radius * radius);
  }
    
  surfaceWeight #(otherMass Decimal, Decimal) = {
    otherMass * surfaceGravity();
  }
}
MERCURY: Planet(3.303e+23, 2.4397e6)
VENUS : Planet(4.869e+24, 6.0518e6)
EARTH : Planet(5.976e+24, 6.37814e6)
MARS : Planet(6.421e+23, 3.3972e6)
JUPITER : Planet(1.9e+27, 7.1492e7)
SATURN : Planet(5.688e+26, 6.0268e7)
URANUS : Planet(8.686e+25, 2.5559e7)
NEPTUNE : Planet (1.024e+26, 2.4746e7)
```

With params
```javascript
// enums idea

// also
veggies: {
    Veggie (String, Int) // a veggie is a String+Int block
    // An instance is structurally matched with a block like this
    // which might be enough if no access to the fields is needed
    SQUASH:{'squash', 1} // do we need `;` e.g. {'squash'; 1}? a: yes
    CABBAGE:{'cabbage'; 3}
    BROCCOLI:{'broccoli'; 0}
    name: {v Veggie; v.0}
    carbs:{v Veggie; v.1}
}
when_eq someVeggie [
    {veggies.SQUASH}:{ print 'Eating squash'}
    {veggies.CABBAGE}:{ print 'Eating cabbage'}
    {veggies.BROCCOLI}:{ print 'Eating broccoli'}
]
name: veggies.name(someVeggie) 
name,carbs: someVeggie()
```

```javascript
// Jun 1st 2023: Yes, just construct instances and assign them to uppercase variables, there's nothing special about them
veggies: {
    Veggie : {}
    SQUASH: Veggie()    
    CABBAGE: Veggie()    
    BROCCOLI: Veggie()    
}
...
Veggie: veggies.Veggie

eat: { 
    what Veggie
    when_eq what [
        {veggies.SQUASH}:{ print 'eating squash'}
        {veggies.CABBAGE}:{ print 'eating cabbage'}
        {veggies.BROCCOLI}:{ print 'eating brocoli'}
    ]
}
    
```    

```javascript
// enums idea
// Jun 2st 2023: Naaaah, uppercase thing followed by a block. It's not clear that that is.
Veggie: {
    SQUASH{'squash' 1}
    CABBAGE{'cabbage' 3}
    BROCCOLI{'broccoli' 0}
    name String
    carbs Int
}

spinach: Veggie{'spinach' 0}
LETTUCE: Veggie{'lettuce' 0}

eat: {veggie Veggie
     
    when_eq veggie, [
        {spinach}: {print 'Eating spinach'}
        {SQUASH}: {print 'Eating squash'}
        {spinash}: {print 'Eating broccoli'}
    ]
}
```
