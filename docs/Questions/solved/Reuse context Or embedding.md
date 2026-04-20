There's not inheritance (because blocks are not classes)  nor embedding, instead to reuse you "just" plug the functionality you want in the new block. 
e.g. 

```javascript
talk: {
    print "Hola, hola, 1,2,3 testing"
}
person: {
    speak: talk  // variable `speak` is a reference to the block referenced by `talk`
}
person.speak() // 'Hola, hola, 1,2,3 testing'

```

Another example
```javascript
// These are singletons so they work just fine
small_circle:{
    radius: {3}
    diameter: {radius() * 2}
}
// So the following 
big_circle: {
    radius: {6}
    diameter: small_circle.diameter // this calls `small_circle.radius` instead of `big_circle.radius`, might not be what we need.
}
...
big_circle: {
    radius: {6}
    diameter: small_circle{radius:radius}.diameter // copies `small_circle` with a new radius definition (6) and then assigns diameter.
 }
```
Because blocks are closures, they capture the environment and they can re-use variables

```javascript
world: {
    times: 0
    talk: {
        times = times + 1
        print 'Hola'
    }
}
person: {
    speak: world.talk
}
person.speak() // incrementes world.times
world.times == 1 // 

```

What happens if the referenced code is part of an instantiable block? 
```javascript
Animal {
    name String
    sound String
    talk: {
        "{name} says: {sound}"
    }
}
cat: Animal( 'Cat' 'meow' )
cat.talk() // "Cat says: meow"
// This shouldn't work: 
Person {
    speak: Animal.talk 
}

Person().speak() //Because `name` and `sound` are not initialized... 
// probably we would need an instance? 
Person: {
    animal: Animal('Person' 'Hello') // valid
    speak:  animal.talk // valid
}
p: Person() // creates animal
p.speak()   // yes, but is that what we want? We would need to have all the state from Animal

```

#answered You can reuse other blocks as long as they are not instance blocks. If you want to use instance blocks members, you need to have an instance.