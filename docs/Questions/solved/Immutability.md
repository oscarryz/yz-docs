
Blocks don't change after they have been created, but you can still re-assing values to them updating that reference. 


```javascript
s: 'name'
t: s // s and t have the string 'name'
s = 'new name' // only s was changed, t still is 'name'
print s // new name
print t // name
```

For instantiable blocks, it might seem they mutate, but they are copied?

```javascript
Person: {
    name String
    age Int
}
a: Person{'no name' 0}
b: a // both have `Person {name: 'no name' age: 0}`
a.name = 'Named' 
// now `a` has 
// `Person {name: 'Named' age:0} and whereas b still have 
// `Person {name:'no name' age:0}`
a.age  = 42      // a has `Person{name:'Named' age:42}`

// Internally a new instance is created (the original is immutable) the data copied and it is assigned to the same 
// variable (`a`)
a = new Person {'Named' a.age}
a = new Person {a.name  42}
```



#dropped No immutability 
#challenged we might need it to solve [Instances for libraries](Instances%20for%20libraries.md)

```javascript
true: Boolean ( 
    ?: {
        tb #(V)
        fb #(V)
        true_block()
    }
)
{true ? { a } { b }}()
{true ? { c } { d }}()

result: true
result ? {1} {2}
true.?.tb

// true.?.true_block == {a}
//
```

[How to deep copy?](How%20to%20deep%20copy?.md)

