TLDR: 

Use generic Option/Result to wrap and unwrap values. See also [Type variants (similar to SumTypes)](Type%20variants%20(similar%20to%20SumTypes).md)

```javascript
compute: {
    if boolean.random() {
        Ok(1)
    } {
        Err('No luck')
    }
}
// Some examples
v: compute.or {
    println 'Failed'
}
v: compute.or_err {
    e Err
    print 'Error {e}'
}
compute ? {
    pritnln 'Yes'
} {
    pritnln 'No'
}
compute .ok_or_err {
    num Int
    println 'Result was {num}'
} {
    err Error
    println 'Result was {err}'
}
// etc. etc more methods to be added on `Option/Result` to handle the error and unwrap
```

Errors are "just" another value, the programmer is responsible for checking it (like in Go)
Possible way to handle it though is to have an errro handler call back, check against a value or have if / else blocks with an "optional class"

Error handling
```zig
// zig
const std = @import("std");

pub fn main() void {
    const file = std.fs.cwd().openFile("does_not_exist/foo.txt", .{}) catch |err| label: {
        std.debug.print("unable to open file: {}\n", .{err});
        const stderr = std.io.getStdErr();
        break :label stderr;
    };
    file.writeAll("all your codebase are belong to us\n") catch return;
}
```

```javascript
fs: std.fs
main: {
   // option 1 compare file to an error file
   file: fs.cwd().open_file "does_not_exists/foo.txt"
   file == fs.no_exists ? {
       print 'unable to open file {info(file).error}'
   } {
       file.write_all 'all your codebase are belong to us\n'
   }
   
   // option 2 pass an error handler
   file_2: fs.cwd().open_file "doesn_not_exists/foo.txt" {
       error Error
       print 'unable to open file {error}'
       break // or return
   }
   // caveat, no way to skip this section
   file.write_all 'all your codebase are belong to us\n'


   // option 3 receive a file and an error
   file_3, error: fs.cwd().open_file 'does_not_exists/foo.txt'
   if error == fs.error {
       print 'unable to open file {info(error)}'
   } {
      file_3.write_all 'all your codebase are belong to us\n'
   }
   // option 4 use Option + "pattern matching"
   of: fs.cwd().open_file 'does_not_exist/foo.txt'
   when[
       {of.some()}:{
          file_4: of.get()
          file_4.write_all 'all your codebase are belong to us\n'
       }
       {of.none()}:{
           print 'unable to open file {info(error)}'
        }
    ]
    // option 5 use Option with 'if/else' logic
   optional_file: fs.cwd().open_file 'does_not_exist/foo.txt'
   optional_file ? {
        file File
        file.write_all 'all your codebase are belong to us\n'
   } {
       error Error
       print 'unable to open file {error}'
   }
   // option 6, Option (or Result) can have an `or` method to unwrap the value or handle the error
   optional_file: fs.cwd().open_file 'invalid.txt'
   file: optional_file.or { 
       error Error;
       print 'unable to open file {error}'
       return // will exit the current block
    }
    file.write_all 'all your codebase are belong to us\n'
    // option 7: Using Result with an API similar to Rust's
    fs.cwd()
	.open_file("does_not_exists/foo.txt")
	.and_then({ f File; f.write_all("all your codebase are belong to us") })
        .or_else({ e Error; print("We found an error `e`)})
   
   

}
```
I like the last approach better with result chaining to `and_the` so we 
act only if the result was Ok.
Definition of `Option`

```javascript
option: {
    T 
    E
    Option: {
        element T
        error   E
        
        ?  { 
            when_present { e T }
            when_absent { e E }
        }
        or {
            when_absent { e E }
        }
    }
    some: Option {
        ? : {
            when_present { e T }
            when_absent { e E }
            when_present(element)
        }
        or : { 
            when_absent { e E }
            element
        }
    }
    none: Option {
        
        ? : {
            when_present { e T }
            when_absent { e Error }
            when_absent(error)
        }
        or : { 
            when_absent { e E }
            when_absent(error)
        }
    }
    of: {
        e T
        e E
        if something_magical {
           some(e t) 
        } {
           none(e t)
        }
    }
}
...
result: option<File fs.Error>.of(valid_file fs.no_error)


```
See [Hare - Introduction](Hare%20-%20Introduction.md) for an example

```js
risky_method {Params Value}
// regular
value: risky_method( params  {error Value; print_error('Oops: {info(error)}') })

// optional parenthesis
value: risky_method  params  {error Value; print_error 'Oops: {info(error)}'  } 

// same optional parenthesis, different 
value: risky_method  params  {error Value
      print_error 'Oops: {info(error)}'  
} 

// no handler, will handle on return
value: risky_method(params)
value == risky_method.error ? {
  print_error('Oops: {info(error)}')
}
// no handler, will handle on return
value: risky_method(params)
if value == risky_method.error, {
  print_error('Oops: {info(error)}')
}

// Option 
value: risky_merhod(params).or {
   error Error
   print_error 'Ooops {error}'
   break
}
```

Other example:

```RED
split-and-save: fn { apples people } { 
  apples / people 
  |^check { 101 "Can't split to no people" }  
  |to-string |write* %apples.txt 
}

split-and-save 10 5 ; returns 5 and saves it to a file
split-and-save 10 0 ; returns a nested failure that turns into error
; Error(101): Can't split to no people 
;	Error: Can't divide by Zero.

; If we want to return 1 in case of any error because we
; need it for further processing we can for example use fix
; function

split-and-save 10 0 |fix { 1 } ; returns 1
```
```javascript
// Assuming the type `Result` exist and has two implementations: `Ok` and `Err`
split_and_save: {
   apples Decimal
   people Decimal
   
   result: apples / people // returns a `Result` which is either `Ok` or `Err`
   // use `or` to unwrap or handle the error
   actual_number = result.or { aritmetic_error Error 
       // Exit early with a custom Error
       return Err<Decimal>{ 
          message: "Can't split to no people" 
          cause: aritmetic_error
       }
   }
   // happy path 
   fs.write actual_number.to_str() "apples.txt"
   result
}
split_and_save 10 5 // returns Ok{5} and saves
split_and_save 10 0 // returns 
                    // Err{
                    //    message: "Can't split to no people" 
                    //    cause: Err{
                    //        message:"Can't divide by Zero."
                    //     }
                    // }
                        
split_and_save( 10 0 ).or_else { 1 } // Calls `or_else` which returns the alternative value.

```


There are no errors or exceptions, only panic that terminates the program. The preferred way is to handle errors is with `option` & `result` blocks: 

```javascript
Repo: {
    users []User
    find_user_by_id: { 
        id Int
        r.users.each { u User 
            u.id == id ? {
                return option.some(u)
            }
        }
        option.none<User>() // or `user: { none: option.none<User>() }` return `user.none`
    }
}
r Repo
u : r.find_user_by_id( 10 ).or_else { User{'default'}}

```

See summary in:  [Error Handling#Summary](../Questions/solved/Error%20Handling.md#Summary)
