[DRN Example](https://www.reddit.com/r/ProgrammingLanguages/comments/1fey7qm/rate_my_syntax/)

## Example

This was supposed to be a "simple" example of 

> Read a CSV file, sum the rows, filter those greater than zero
> create a new CSV and write it

_(There's a short syntax reference at the bottom)_ 

```javascript
main: {
   os.open("input.csv").and_then {
     file File 
     file.read_all()
   } .and_then { 
     content [Int]
     
     output: String(content)
       .split("\n")
       .map {
          line String
          line.split(",").map {
            column String
            int.parse_int(column).or(0)
          }
       .reduce(int.+)
       .filter(0.>)
       .join(",")
          
      }     
      os.open("output.csv").and_then {
         file File
         file.write_all(output.to_array())
      }
  }.or { 
     e Error
     println("Coulnd't process file. Error: `e`")
  }
}
```


 
    

Now, as it can be seen, a lot of the code relies on standar library functions. Error handling is done similar to Rust, where theres a type with two possible outcomes: Ok or Error

## std functions signatures

Here's made up signature of all the standard library functions signatures used above:

```javascript
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
      map: #(#(T,E), [E])
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

- `{ stuff }` a block of code containing _stuff_ . Variables inside blocks work as both parameter and return values (and attributes).
  - e.g. `{ print("Hello world") }` 
- `TypeName` uppercase name denotes a type
  - e.g. `String`
- `TypeName(params)` creates an instance
  - e.g. `String("hello")`
- `identifier:  expr`  declares and initializes a variable of type `expr`
  - e.g.: `n : 1` `n` has type `Int`
- `identifier TypeName` declares a variable of type `TypeName`
  - e.g. `s String` declares `s` of type `String`.
- `[TypeName]` array of that tyep
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
