#feature 
# Creating Instances

Uppercase-named bocs are types. Call them with `()` to create a new instance:

```yz
Person: {
  name String
  age Int
}

alice: Person("Alice", 30)       // positional
bob: Person(name: "Bob", age: 25) // named arguments
```

Named arguments can appear in any order:

```yz
bob: Person(age: 25, name: "Bob")  // same as above
```

## Nested initialization

Instances can be composed inline:

```yz
Frame: {
  title String
  width Int
  content TextField
  visible Bool
}

TextField: {
  value String
}

model: HelloWorldModel(saying: "Hello, World!")

win: Frame(
  title: "`model.saying` App",
  width: 200,
  content: TextField(value: model.saying),
  visible: true
)
```

## Default values

Fields with defaults are optional during construction:

```yz
Config: {
  host String = "localhost"
  port Int = 8080
}

c1: Config()                    // host="localhost", port=8080
c2: Config(port: 9090)          // host="localhost", port=9090
c3: Config("example.com", 443)  // both overridden
```

See [Bocs](Bocs.md) for the full rules on parameters vs fields.
