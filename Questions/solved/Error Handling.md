[Error handling](../../Features/Error%20handling.md) provides some options to handle errors, most promising is `Optional` with `some` and `none`  instances, similar to `Boolean` with `true` and `false` but still needs some work. 

The simplest thing and original idea is to define specific errors for a type and then compare against that: 

```javascript
std: {
   fs: {
       File: {
          write_all: {s String; ... }
       }
       'This is an error'
       file_not_exists: File{}
       open: { path String;
          // etc etc
          // return
          file_not_exist
       }
   }
}
...
file: std.fs.open 'non_existing'
if file == std.fs.file_not_exists {
    file.write_all 'Content'
} {
    std.err.print 'Cannot open file because it does not exists'
}

```


#answered 

# Summary
There's no errors, so use classes to mark return values as erroneous. 
Initially I thought comparing with a given instance `if file == files.not_found`  where `not_found` is a `File`, but that was when I had one return value. Now we can return more than one we can either do like Go: 
```
value error : stuff()
error ? {
    handle
}
```

Or like Rust with `Option/Result`
```
result : stuff()
value : result.or {
    default_value
}
```

Note: We can return more than one value because all the variables are accessible. Also we can return an expression block

```
foo: {
    v: logic()
    if v == -1 {
        { v error('Invalid') }
    } {
        { v ok() }
    }
}
v e: foo()
e ? { 
    handle_error()
}
```
See more examples in [Zig - If](Zig%20-%20If.md)



