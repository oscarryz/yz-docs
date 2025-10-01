#example
 
[Strong static typing, a hill I'm willing to die on...](https://www.svix.com/blog/strong-typing-hill-to-die-on/)

```js
birthday_greeting: {
    params [] 
    "{params[0] is params[1]!"
}
birthday_greeting_2: {
   name
   age
   "`name` is `age`!"
}


birthday_greeting_3: {
   name String
   age  Int
   "`name` is `age`!"
}
type_of (birthday_greeting_3) // { String Int String } 


// 
Person : { 
   name String
   age  Int
}
birthday_greeting_2: { 
   person // generic variable
  '`person.name` will turn `person.age + 1 ` next year!' 
   // person will need to have a variable `name` of type any, and a `age` of type int (or something with `+` method that takes a int
}

birthday_greeting_3: {
   person Person
  '`person.name` will turn `person.age + 1 ` next year!' 
  // person is not generic, and would only take `Person`'s or something that matches 
}

main: {
   person Person = { name: 'Hello' age: 12 }
   birthday_greeting_2 person
   birthday_greeting_3 person
}


//  There's a bug below
Person : {
   id String
   name String
   age Int
}
Pet : {
   id String
   owner String
}
   
id: 'p123'
person: Person { 'John' 20 }
cache.set 'person-{id}' person 
pet Pet = cache.get 'preson-`id`' // failure over failure... 

// Getting better
PersonCacheKey: { String } 
Person: {
   id String
   name String
   age Int
}
PetCacheKey {}
Pet  {
   id String
   owner String
}
id : 'p123'
person : Person { id 'John' 20 }
cache.set PersonCacheKey{id} person 
pet Pet = cache.get PersonCacheKey{id} 

// Even better
PersonId {String}
PetId {String}

Person {
  id  PersonId
  name String
  age Int
}
// In the example this would ensure each thing is different, but in Yz they are the same because they structural match 
Pet { 
  id PetId 
  owner PersonId 
}

p1: new_person()
p2: new_person()
child: make_child p1 p2 
```
