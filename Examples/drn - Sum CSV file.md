#example

[DRN Example](https://www.reddit.com/r/ProgrammingLanguages/comments/1fey7qm/rate_my_syntax/)

## Example

This was supposed to be a "simple" example of

> *Read a CSV file, sum the rows, filter those greater than zero*
> *create a new CSV and write it*

_(There's a short syntax reference at the bottom)_

```js
// Opens a file and reads it into a String
read_file: {
  file_name String
  os.open( file_name ).and_then {
    file File
    string.from_array(file.read_all(), "UTF-8")
  }
}

// Reads a String, and sums each row and format it back to a string
sum_rows: {
  content String
  Ok(content
    .split("\n")
    .map {
      line String
      line.split(',').map {
        column String
        int.parse_int(column).or(0)
      }
    }
    .reduce(int.+)
    .filter(0.>)
    .join(",")
  )
}
// Creates a block bounded to a file name.
// This block  receives a string
// and writes it to a file and returns a Result of nothign
write_result: {
  output_file String
  {
    output String
    os.open(output_file) .and_then {
      file File
      file.write_all(output.to_array())
    }
  }
}

main: {
  read_file("input.csv")
    .and_then(sum_rows)
    .and_then(write_result("output.csv"))
    .or_else {
      e Error
      println("Coulnd't process file. Error: `e`")
    }
}

```


Now, as it can be seen, a lot of the code relies on standard library functions.

Error handling is done similar to Rust, where there is a type with two possible outcomes: `Ok` or `Error`

## std functions signatures

Here's made up signature of all the standard library functions signatures used above:

```js
// The block `std` has the blocks: `result`, `int`, `string`, `array` and  `os` in this example
// each one of them declares functions and types
// ... means: ommited for brevety
std: {
  result: {
    T
    Result: {
      and_then #(action #(T, Result(T)), Result(T))
      or  #(e Error)
      ...
    }
    Ok : {
      data T
      ...
    }
    Error: {
       ...
    }
  }
  // int defines the Int type nd operations on it
  int: {
    parse_int #(input Int, Result(Int))
    + #(Int,Int,Int)
    Int : {
        > #(Int, Bool) // e.g. 1 > 2 returns false
        + #(Int, Int)  // e.g. 1 + 2 returns 3
    }
    ...
  }
  string: {
       String: {
         data [Int]
         split #(String, [String])
        }

     ...
  }
  array: {
    T
    // This is the "magic" implementation type
    // under the `[]` syntax
    Array: {
      map: #(#(T,U), [U])
      reduce #( #(T, T, T), [T])
      filter #( #(T, Bool), [T])
    }
    ...
 }

  os: {
    open #(file_name String, Result(File))
    File : {
      read_all #(Result([Int]))
      write_all #([Int], Result(Int))
    }
    ...
  }
  ...
}
```


## Syntax

Some syntax short explaination:

- `{ stuff }` a block of code containing _stuff_ . Variables inside blocks work as both parameter and values (and attributes).
  - e.g. `{ print("Hello world") }`
- `TypeName` uppercase name denotes a type
  - e.g. `String`
- `TypeName(params)` creates an instance
  - e.g. `String("hello")`
- `identifier:  expr`  declares and initializes a variable of type `expr`
  - e.g.: `n : 1` `n` has type `Int`
- `identifier TypeName` declares a variable of type `TypeName`
  - e.g. `s String` declares `s` of type `String`.
- `[TypeName]` array of that type
  - e.g. `a [String]` is an array of strings
- `expr.identifier(params)` invoked the method `identifier` on the resulting type of `expr`
  - e.g. `"hi".len()` (can ommit parenthesis if more than one argument)
- `expr.identifier` access the field `identifier` on the resulting type of `expr`
  -  e.g. `{ a : 1 }.a` access the variable `a` declared in the block.
-  `#(TypesName)` block signature
  -  e.g. `#(Int,Int)` is a block that takes an `Int` and returns an `Int`
- `T` lonely single uppercase letter, is a generic.
  - e.g. `Box(T)` generic box, `Box(String)` instantiates a box with type String
- String interpolation uses backtick
  - "Hello `whom`" would be `"Hello world"` if `person` is `"world"`
