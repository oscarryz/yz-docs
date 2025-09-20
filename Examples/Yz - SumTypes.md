#example

1. Same data types on all the variants
  1. No values
```js
Switch:{
  On(),
  Off()
}
```
  
  2. Same values
  ```js
FileAccess:{ 
    Read(v:1),
    Write(v:1),
    Execute(v:1)
  }
```

  2. Different values
```js
Planet:{
    Mercury(mass: 123, radius: 456),
    Venus(mass: 9393, radius: 9881)
    ... 
  }
```
1. Different data types on all the variants
  2. All the variants differ
  ```js
Shape:{
    Empty(),
    Circle(radius Double),
    Rectangle(width Double, height Double),
    Point(x Int, y Int),
    Text(msg String, font_size Int)
  }
```
  4. Some of the variants differ
  ```js
Status:{ 
    Success(code Int),
    Warning(code Int),
    Error(code Int, msg String),
    Timeout(), 
    Pending(),
  }

classify_status #(status Status) { 
  match status
    {Success => print("Sucessful response (2xx): `status.code`")},
    {Warning => print("Redirection (3xx): `status.code`")},
    {Error   => print("Error ocurred: `status.mgs`(`status.code`)")},
    { println("Other status")},
}
```


Example from [Core-lang](https://core-lang.dev) landing page: 

```js
// Pet has 3 properties  defined by 2 variants: 
// `name String`, `lives Int` and `years Int`
Pet: {
  Cat(name String, lives Int),
  Dog(name String, years Int),
}

// Describe uses the `when` matching
// The variant is used to match but 
// no need to deconstruct there e.g. Cat(n, l)
// because now is safe to access pet.lives directly
// TODO: Decide if the parameter type is needed 
// e.g. Cat(String,Int), or even the `()` e.g Cat()
describe #(pet Pet, String) {
	match pet
    { Cat => "cat `pet.name` has `pet.lives` lives" },
    { Dog => "dog `pet.name` has `pet.years` of age" },
}

[Cat("Lila", 7), Dog("Fenton", 6), Cat("Molly", 9)]
   .filter({p Pet   ; p.name.len() > 4})
   .map   ({p Pet   ; describe(p) })
   .each  ({d String; print(d)})
// -> cat Molly has 9 lives
// -> dog Fenton is 6 years of age

```


