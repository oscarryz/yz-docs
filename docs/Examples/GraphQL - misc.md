---
type: example
generated: { by: "oscarryz", at: 2023-12-06T20:13:02-06:00 }
---
#example

```js
Project: {
    name String
    tagline String
    contributors [User]() // [User]
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
