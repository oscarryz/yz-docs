#example

```js
Project: {
    name String
    tagline String
    contributors []User // [User]
}

// Get tagline from project GraphQL
query: {
    where: 'project.name="GraphQL"'
    field: ['tagline']
}
// GraphQL
{
    project(name:'GraphQL') {
        tagline
    }
}
// Yz
{
    project: {
        name:'GraphQL'
        fields:['tagline']
    }
}
```
