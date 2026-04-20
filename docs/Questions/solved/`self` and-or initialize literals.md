#answered  if needed, add a `self` variable

eg in [[Arturo  - Examples]]  the int `2` is used: 

```javascript
[1 2 3 4 5].map 2.*
```

Assuming the `Int` class would have a `value`  variable

```javascript
Int: {
    value Int
    *: {
        n Int
        value * n
    }
}
```

Where `value` is the equivalent of `this` or `self`. 
Also, how to do built-ins? like `*` in this case


Possible answer, set self as a variable in a constructor

```javascript
person: {
    new: {
        name String
        p: Person{name: name}
        p.self = p
    }
    Person: {
        name String
        self Person
        to_string: {
            self.name
        }
    }
}
me: person.new "Yz"
print '{me}' // 

```

If follow the idea of [Private read only variables](../../Features/Replaced%20features/Private%20read%20only%20variables.md) + [Consider Immutability](Consider%20Immutability.md)  we could have a private self that will be overriden on the next instance
Update: this doesn't work becuase `self` will be overriden and will change. I thought I could get away with that if the instance kept a separate reference. 

```javascript

person: { 

    self Person
    new: {
        name Sring
        self = Person {name}
        self
    }
    Person: {
        name String
        to_string: {
            self.name
        }
    }
}

p: person.new 'One'
/*
internal state: 
  
   new.name = 'One'
   
   person.self = Person { name: new.name  } // 'One'
*/
q: person.new 'Two'
/*
    new.name = 'Two'
    person.self = Person { name: new.name} // 'Two'
    p.to_string == self.name // 'Two' 
*/
```

