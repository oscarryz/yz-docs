#feature 
# Error Handling

Yz uses `Option` and `Result` variant types for error handling — no exceptions, no null references.

## Option

`Option` represents a value that may or may not be present:

```yz
Option: {
  V
  Some(value V)
  None()
}
```

Use `match` to handle both cases:

```yz
result: find_user(42)
match result
  { Some => print("Found: `result.value`") },
  { None => print("Not found") }
```

Or use `or` to provide a default value:

```yz
name: find_name(id).or("anonymous")
```

## Result

`Result` represents either success or failure, with typed error information:

```yz
Result: {
  T, E
  Ok(value T)
  Err(error E)
}
```

Use `match` to handle both outcomes:

```yz
divide #(a Int, b Int, Result) {
  b == 0 ? {
    Err("division by zero")
  }, {
    Ok(a / b)
  }
}

r: divide(10, 2)
match r
  { Ok  => print("Result: `r.value`") },
  { Err => print("Error: `r.error`") }
```

## Chaining

Methods like `and_then` and `or_else` allow composing fallible operations:

```yz
process_file #(filename String, Result) {
  read_file(filename)
    .and_then { content String
      parse_content(content)
    }
    .or_else { error String
      print("Failed: `error`")
    }
}
```

See also [Type variants](Type%20variants.md) for the general variant/sum type mechanism.
