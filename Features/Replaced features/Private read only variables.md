#feature
 Private(ish) variables

In a block create an entity block with the access method e.g. 

```js
person: {
    Person: {
        get_age {Int}
    }
}
```

Add a factory method returning a new definition overriding the accessor

```js
person: {
    ... 
    new_with_age: {
        age Int
        r: Person {
            get_age: {
                age
            }
        }
    }
}
```
Use the factory method
```js
someone: person.new_with_age(42)
someone.get_age() // returns 42
someone.age // doesn't have a variable `age`
someone.get_age.age // not sure if can be visible though 
    
```

Jan 11 2023

Probably a better way would be to create a block "interface" exposing the "public" fields and have another structuarally implemeting it

In the example above we want the method `get_age` but not expose the variale `age`

```javascript
Person: {
    age Int
    get_age: {
        age
    }
}
GetAger: {
    get_age {Int} 
}

// Creates a `Person` instance 
// but the `print_age` can only see the `get_age` method
// from the `GetAger` "interface" 
print_age: {ager GetAger
    print '{ager.get_age()}'
}
print_age Person{age:42} 
```

Another way is to add a signature to the block and don't include the variables

```javascript
secret { result Int } = {
    a: 2
    b: 4
    result = a * b
}
secret() // secret.result = 8

x: secret.result // z = 8
y: secret.a // Compilation error: no variable named 'a' in block signature
z: secret.b // Compilation error: no variable named 'b' in block signature
```

Which is similar to the above, because the type is defined as `{result Int}` which is the interface and `GetAger` is also exposing a public interface as `{get_age{Int}}`
```javascript

```

#answered 

Add a block signature. Only variables included in the signature will be publicly accessible.

```javascript
Counter #( get #(Num), inc #() ) = {
    // not included in signature, not visible from the outside
    count : 0
    get: {count} // the value is returned
    inc: {count = count + 1} // and modified
    ++: inc // "alias" of `inc`
}
a: Counter()
a.get() // returns 0
a.inc() // sets a.count to 1
a.++()  // set a.count to 2
a.get() // 2
// compilation error, no variable `count` in type Counter
// a.count 
```

```js
Counter #(c Int) = {
  c Int
}
counter : Counter()
counter.c = 1 // compilation error, counter.c cannot assign value outside of the block, consider using a modifier method e.g. counter.set_c (or counter.c=)
...
Counter #(c= #(Int), c #(Int)) = { 
  _c Int
  c= { c Int; _c = c}
  c : { _c }
}
counter : Counter()
counter.c= counter.c() + 1 

```