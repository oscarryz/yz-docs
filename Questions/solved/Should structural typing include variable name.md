Should structural typing match include variable name or only types are needed. 


Example 1:  Only needing types

```javascript
Person { String String }

oscar: {
   name: 'Oscar'
   last_name: 'Ryz'
}

p Person = oscar // yes, it has two strings

thing: { 
   serial_number: '123109823'
   make: 'acme'
}

p = thing // yes, it has two strings 
```

Example 2: Variable names + types
```javascript
Person { name String  last_name String }
oscar: {
   name: 'Oscar'
   last_name: 'Ryz'
}

p Person = oscar // yes, it has two strings 

thing: { 
   serial_number: '123109823'
   make: 'acme'
}

p = thing // compilation error, thing doeesn't have 'name' nor 'last_name'
```

Possible solution: 
Let the type define what it wants, if it specifies variable name, then it's needed. If it doesn't then don't ...

```javascript

Resource { name String String }
r Resource
o: { name: 'Oscar' last_name 'Ryz' }
r = o // yes, it has a `name` of type string and another string
t: { name 'xs_100' serial: '12301823' make: 'ACME' }
r = t // yes, it was a variable `name` of type string and at least other variable (order matters)
```


What about generics?  How to prevent a "typeless" variable be separated from a return type

```javascript
GenericData { data ; String }
p GenericData = { 1 'Hello'}
```


#answered Yes, if name is specified, name is required. If only type, then type. Works for generics too, in case of generics, use `;` to separate members. 