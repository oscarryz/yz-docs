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

## Boc Declaration

A boc variable can declare its signature and body together — this is the **boc declaration** form:

```yz
greet #(msg String, to_whom String, String) {
    "`msg` `to_whom`"
}
```

Labeled entries are input fields; the trailing unlabeled `String` is the output.

## Definite Assignment

A typed declaration with no default is **uninitialized**. The compiler requires it to be assigned on all control-flow paths before it is read:

```yz
result Int          // uninitialized
result = compute()  // assigned
print(result)       // OK — assigned before read
```

Reading an uninitialized variable is a compile error:

```yz
x Int
print(x)   // compile error: x used before assignment
```

This also applies across conditional branches — if any path skips the assignment, the read is an error:

```yz
value Int
flag ? { value = 1 }, { }  // one branch doesn't assign value
print(value)               // compile error: value may be uninitialized
```

The fix is either to assign on every path, or give the variable a default:

```yz
value: 0   // default — always initialized
flag ? { value = 1 }, { }
print(value)  // OK
```

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

See [Generics](docs/Features/Replaced%20features/Generics.md)