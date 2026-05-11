#example 

The [Compile Time Bocs](Compile%20Time%20Bocs.md) provides code generation such as embedding (or mixing) a boc content into another, similar to Go's embedding but through code generation

```js

// The Named type
Named: {
    name String
    hi: {
        print("My name is ${name}")
    }
}

`
compile-time:[Mix]
mix: [Named]
`
Person: {
    last_name String
}

/* 
// The compile-time Mix generates:
Person: {
    name String
    last_name String
    hi: {
        print("My name is ${name}")
    }
}
*/
p: Person("Alice", "Smith")  
p.hi()        // My name is Alice
p.name = "Jane"  
p.hi()        // My name is Jane

```

