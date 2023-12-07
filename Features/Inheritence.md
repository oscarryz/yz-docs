Not like this, is not super obvious and still not decided if import would do something else that just use the names like in Java
```javascript
Named: {
    name {String}
    upper_case: {
        name().upper_case()
    }
}
Boo: {
    // Dropped, see import
    import Named
    // name: {'Boo'} // ce: name already defined
    name =  { 'Boo' } 
}
b: Boo{}
b.upper_case() // 

```

Feb 15 2023

Current accepted: in the new Block declaration hook up members from other blocks.


```javascript
Named: {
    name {String}
    upper_case: {
        name().upper_case()
    }
}
Boo: {
    name: {'Boo'}
    upper_case: Named.upper_case
}
//
n: Boo{}
n.upper_case() // BOO
```

If more are needed 

```javascript
Animal: {
    tell {}
    pet  {}
    feed {}
}
Cat: {
    tell { println 'Meow' }
    pet  { println 'purr' }
    feed { println 'lick' }
}
Lion: {
   tell: { println 'Roar' }
   pet: Cat.pet
   feed: Cat.feed
}
```

### Alternatives

It is possible to create an "anonymous" instance with some methods changed

```javascript
lion: Cat {
    tell : { println 'Roal' } 
}
```
Which is completely valid and the main way to have different behaviour. This is valid because in Yz, that's how new values are assigned: 

```javascript
Person : {
    name String
}
me: Person { name: 'Yz' }
```

Still open question: [Reuse context? Or embedding?](../Questions/Reuse%20context?%20Or%20embedding?.md)
