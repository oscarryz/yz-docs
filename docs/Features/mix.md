#feature 

_Might be replaced by Compile time bocs_


# Code Composition: `mix`

The `mix` keyword merges the fields and methods of one boc into another. It provides compositional reuse without inheritance.

## Basic usage

```yz
Named: {
  name String
  hi: {
    print("My name is ${name}")
  }
}

`
compile_time:[Mix]
mix: [Named]
`
Person: {
  last_name String
}

p: Person("Jon", "Doe")
p.hi()        // My name is Jon
p.name = "Jane"
p.hi()        // My name is Jane
```

`Person` gains `name String` and `hi` from `Named`, plus its own `last_name String`. All members are accessible without any prefix.

## What mix does

- **Flattens fields**: `Named`'s `name` field becomes a field of `Person`
- **Promotes methods**: `Named`'s `hi` method becomes callable directly on `Person` instances
- **Preserves binding**: `hi` still closes over `name` — it uses the `Person` instance's `name`, not some separate `Named` instance

## Conflict detection

Any name clash between the host boc and the mixed-in boc is a compile error:

```yz
Loggable: {
  level Int
}

Config: {
  mix Loggable
  level String   // Error: 'level' is already defined by Loggable
}
```

This prevents silent shadowing. Conflicts between two mixed-in bocs are also errors.

## Methods on mixed-in bocs

Methods defined in the mixed-in boc are available directly on the host:

```yz
Named: {
  name String
  hi: {
    print(name)
  }
}

Person: {
  mix Named
  last_name String
}

p: Person("Alice", "Smith")
p.hi()  // Alice
```

## Cross-file mix

When composing bocs defined in separate files, use the fully-qualified name:

```yz
mix house.front.Named
```

See [Code organization](Code%20organization.md) for how file paths map to boc names.
