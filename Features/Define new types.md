If we have a new type `Person`  


```
Person {}
```

How can differentiate it from instantiating it? 
```
Person {}
p: Person{}
```

#answered 
Use parenthesis 
```
Person{}
p: Person()
```

The current option suggested by [Instances](../Features/Instances.md) (I think) is to see if it was defined earlier, but that doesn't seem to be enough. 

A special syntax could be needed ( `::` was suggested earlier )  but I want it to keep it cleanish, probably a `type` keyword if it is not initialized?

```
type Person {} // no
```
And then  on initialization 

```
type Person {} = {} // no

```
And inference: 

```
Person: {} // yes
```

How would it be used as variable? 

```
a Person // yes

```

Probably `::` reserved to declare it? Nah

```
Person :: {}
method {}
method {} = {}
method: {}
method{}
```

Possible answer, use `::` only to declare the type of a type ... (type alias?) 

So if a `Person` has a name of type string
```javascript
Person :: { 
    name String
}
p: Person{} // not possible... Person has not been initialized 
p Person = { name = "Alice"} // valid, the block is giving a value that satisfies Person type signature.. 
// this is slighly similar to an abstract class? 

// We can then initialize `Person`
Person :: {
    name String
} = {
    name String // seems silly
}
// That's why you use type inference
Person : {
    name String
}

```

This is only reserved for uppercase names, so the following is an error `person :: {}`

Using `::` should be rare and only needed when a type signature is needed to a new type and most likely to create a hierarchy 

```javascript
// Profession declares a type with a `name` method that returns `String`
Profession :: { 
    profession_name :: {
        String
    }
}
// Declared and inferred
Teacher: {
    profession_name: {
        'A teacher'
    }
}
Director: {
    profession_name: {
        'A director'
    }
}
// `Teacher` and `Director` "implement" `Profession`, we can't instantiate Profession directly but we can use it as type
print_profession_name: {
    person Profession
    print person.profession_name()
} 
print_profession_name(Teacher{}) // 'A teacher'
print_profession_name(Director{}) // 'A director'

```

Is not necessarily this is a block, it could be another type creating a [Type Alias (wip)](Type%20Alias%20(wip).md) 
You can't add new methods until [Extension methods](../Questions/solved/Extension%20methods.md) are supported.
 
```
Time :: Int
```
## Summary

- Use `::` to define the type of a new type.
    - Could be a block or an existing type
- Only allowed in uppercase starting identifiers
- Use `:` to infer a new type

#answered 
### Solves: 
[The block type](solved/The%20block%20type.md)
[Type Alias (wip)](../Features/Type%20Alias%20(wip).md)
[Instances](../Features/Instances.md)

### Related:
[Blocks](../Features/Blocks.md)
[Block type](../Features/Block%20type.md)

