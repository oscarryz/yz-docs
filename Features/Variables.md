#feature

Variables are declared with an identifier followed by a type identifier or type signature.

```yz
// Declares a variable `msg` of type String
msg String

// Declares a variable `salute` of type "boc returning String"
salute #(String)
```

## Initialization

Variables can be initialized with `=`:

```yz
// Declares and initializes `msg` with value "Hi"
msg String = "Hi"

// Declares and initializes `salute` with a boc value
salute #(String) = {
    "Hello world"
}
```

The shorthand `:` declares and initializes a variable, inferring the type from the value:

```yz
msg: "Hi"
salute: {
    "Hello world"
}
```

## Boc signature shorthand

A boc variable can declare its type signature and body together in one statement:

```yz
greet #(msg String, to_whom String, String) {
    "`msg` `to_whom`"
}
```

Here the trailing `String` (unlabeled) is the return type; the labeled entries are inputs.

## Variables as parameters

Inside a boc, a typed declaration with no default value is a **required parameter** — it must be provided when the boc is invoked:

```yz
add: {
    a Int     // required parameter
    b Int     // required parameter
    a + b
}

result: add(3, 4)   // → 7
```

A declaration with a default value is an optional parameter:

```yz
greet: {
    name String = "World"   // optional, defaults to "World"
    "Hello, `name`!"
}

greet()           // → "Hello, World!"
greet("Alice")    // → "Hello, Alice!"
```

## Generic type variables

Variables can carry a generic type parameter:

```yz
item T   // T is a generic type variable
```

See [Generics](./Generics.md)