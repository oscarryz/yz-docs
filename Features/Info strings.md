# Info Strings

An info string is a string literal placed **before** a declaration. It attaches documentation metadata to the declaration, retrievable at runtime via `info(element)`.

```yz
"A message"
message: "Hello"

info(message).text  // "A message"
```

Info strings use regular string delimiters (`'` or `"`). They are compiled alongside the element they describe.

## Inline info strings

Any element can have its own info string:

```yz
"Prints a personalized greeting"
say_hello: {
  "What message to display"
  what: "Hello"

  "The recipient name"
  who: "World"

  print("`what`, `who`!")
}
```

## Structured info strings

Info strings can contain structured data (interpreted by tooling):

```yz
"
Prints the classic Hello, World! program.

variables: {
  what String = 'Hello'
  who  String = 'World'
}

tests: {
  assert say_hello() == 'Hello, World!'
  assert say_hello('Hola') == 'Hola, World!'
  assert say_hello(who: 'there') == 'Hello, there!'
}
version: 1.0
author: 'Yz developers'
"
say_hello: {
  what: "Hello"
  who: "World"
  print("`what`, `who`!")
}
```

## Retrieving info at runtime

```yz
i: info(say_hello)
print(i.text)    // prints the info string text
i.tests()        // runs the test examples
i.version        // 1.0
```

Info strings are a foundation for self-documenting code, doctest-style testing, and serialization annotations — all without a separate documentation language.
