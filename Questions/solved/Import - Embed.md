Embedding instantiable blocks causes the enclosing obj to have the values as if they were declared previously.
They can be re-assigned, but no re-declared

```js

Base: {
    name String
}
Container: {
    embed Base
    to_string: { name }
}
// usage
b: Base { 'base' }
c: Container { 'container' }
print(c.to_string()) // 'container'



Base: {
    b   Int
    tag String
    describe_tag : {
        print 'Base tag is {tag}'
    }
}
Container: {
    Base
    c String
    //tag String // error: already declared in Base
    describe_tag = {
        print 'Container tag is `tag`'
    }
}
Component: {
    use Base
    name String
}

b: Base(b: 10, tag: "b's tag")
co: Container{c: 'foo', tag: "co's tag"} 
c: Component{ name: 'Comp'} // ce: need value for tag
print(b.describe_tag())// Base tag is b'stag
print(co.describe_tag()) // Container tag is co'stag
//print c.describe_tag() // Base tag is b'stag

```

#answered Feb 15 2023
No embed, but you can "hook", behavior from other blocks, see  [Inheritence](../../Features/Pending%20features/Inheritence.md)

Import is a simple name shortening see [Import](../../Features/Pending%20features/Import.md)























