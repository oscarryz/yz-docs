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

## Required vs optional fields

Fields with no default are **required**. You can construct an instance without providing them, but reading an uninitialized field is a compile error:

```yz
Person: {
  name String
  age Int
}

p: Person()    // OK — instance created, fields uninitialized
p.name         // compile error: name used before assignment
```

Assign them first:

```yz
p: Person()
p.name = "Alice"
p.age = 30
p.name         // OK
```

Or provide values at construction:

```yz
p: Person("Alice", 30)   // positional — both fields initialized
```

Before passing an instance to another boc, all its required fields must be assigned:

```yz
greet #(p Person) { print(p.name) }

p: Person()
greet(p)          // compile error: p.name not assigned

p.name = "Alice"
p.age = 30
greet(p)          // OK
```

Fields with a default value are **optional**:

```yz
Config: {
  host String = "localhost"
  port Int = 8080
}

c1: Config()                    // host="localhost", port=8080
c2: Config(port: 9090)          // host="localhost", port=9090
c3: Config("example.com", 443)  // both overridden
```

Use `Option(T)` only when a field may legitimately have no value (absence is meaningful). Don't use it as a workaround for required fields you haven't filled in yet — that's what default values and definite assignment are for (see [Variables](Variables.md)).

See [Bocs](Bocs.md) for the full rules on parameters vs fields.
