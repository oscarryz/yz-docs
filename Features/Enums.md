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

With params
```javascript
// enums idea

// also
veggies: {
    Veggie (String, Int) // a veggie is a String+Int block
    // An instance is structurally matched with a block like this
    // which might be enough if no access to the fields is needed
    SQUASH:{'squash' 1} // do we need `;` e.g. {'squash'; 1}? a: no
    CABBAGE:{'cabbage' 3}
    BROCCOLI:{'broccoli' 0}
    name: {v Veggie v.0}
    carbs:{v Veggie v.1}
}
when_eq someVeggie [
    {veggies.SQUASH}:{ print 'Eating squash'}
    {veggies.CABBAGE}:{ print 'Eating cabbage'}
    {veggies.BROCCOLI}:{ print 'Eating broccoli'}
]
name: veggies.name(someVeggie) 
name carbs: someVeggie()
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

```js
// Enum idea Jun 2024 :P
tree: {

	Empty (_) = {
	}
	Leaf (T) = {
	}
	Node (T Tree(T) Tree(T)) {
	}
	Tree ( a T )
}
// Jun 2024 Rejected: Looks like something is initialized but then assigned something else...
```
