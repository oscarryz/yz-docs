```javascript
// Current: 

Node: <T>{
    data T
    left Node<T>
    rigth Node<T>
}

// Alternative
Node: {
    data T+
    left Node
    rigth Node
}
node: Node{'Hola' Node{'Adios'} Node{'Gracias'}}

// Yet another
Node: {
   data // generic like a' 
   left Node
   right Node
}
//
Foo: {
   bar
   baz
}
```

#answered for now keep using the `<T>` syntax
#challenged It makes messy almost everytime I want to use it. 

```javascript
Bool: {
    ? : {
       tb { $T }
       fb { $T }
    }
}
value: is_it_true() ? {
   'Yes'
} {
   'No'
}
// value is String ? should bind to String


box_from_array: {
    array []a
}
```

This is examined in [Generics without <>](Generics%20without%20<>.md)
