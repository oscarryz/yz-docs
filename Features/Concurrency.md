# Concurrent by Default 
Yz is concurrent by default. Every method invocation will run in its own coroutine/thread asynchronously, although still sharing the same memory space.
```js
// Run concurrently 
foo()
bar()
```

If the result of a method call is assigned to a variable or used as argument to another boc, the code will be synchronous and will wait until the call completes

```
// Execution will wait for `fr` to have a value from the completion of `foo()`
fr: foo()
// `bar` will start execution when `baz()` call completes.
bar(baz())
```

The block enclosing the bocs calls, will complete when all the internal bocs have completed (similar to structural concurrency without scopes)

```
parent_boc : { 
   foo()
   bar() 
}
// executing `parent_boc` will launch `foo()` and `bar()` one of them might finish earlier than the other
// `parent_boc()` will complete only when both have completed
parent_boc()
```

To sum up: 
- Every method call is async
- Assigned as variable (or used as argument for another method) the execution will wait
- Once all the calls finishes, the method is self will finish. 

**Note** executing multiple calls on the same boc will effectively be sequential too, as they are executed as received

```
// `two` will run until `one` completes 
foo.one()
foo.two()
// but bar will trigger immediately
bar() 
```

If order of execution is important, synchronize it by waiting for the result. 

# Example

Because all the bocs synchronize at the end of the enclosing bloc, you can retrieve the result by accessing its state there.
```
{
  foo() // puts the result in a variable `r`
  // other calls
  foo.r // at this point foo.t has a value 
}
```

Example modified from [Structured Concurrency in Java](https://openjdk.org/jeps/428#:~:text=For%20example%2C%20in%20this%20single%2Dthreaded%20version%20of%20handle()%20the%20task%2Dsubtask%20relationship%20is%20apparent%20from%20the%20syntactic%20structure%3A)
```javascript
//Yz synchronous
synchronous_example : {
    the_user: find_user()
    the_order: find_order()
    Response(the_user, the_order)
}
// Yz concurrent
concurrent_example: {
    find_user()
    find_order()
    if some_logic() {
        other_thing() // maybe this is executed first
    }
    Response(find_user.user, find_order.order)
}
```

Timeout: 
Create a block that exits the current block by calling `return` after the specified amount of time. How to clean up resources is still TBD

```javascript
fetch: {
    id String
    // This will make the parent boc return after 10 seconds
    time.sleep(10.seconds(), {
        // after 10 seconds just return finilizing the `fetch` execution
        return
    })
    // Will execute right away
    return find(id) // explicity `return` will ignore the wait for `time.sleep` to complete. 
}
```

Error handling
A couple of options:
1. Have the async function return an Option/Result kind of data (Yes)

```javascript
fetch:{
    id String
    timeout(10, {return})
    // option 1
    data: find(id).or_else({ "Couldn't find data"})
}
```


See: [Go - Go Concurrency Patterns](Go%20-%20Go%20Concurrency%20Patterns.md)

This might be replaced with: [lazy eval pointers](Questions/Async%20%2B%20Lazy%20eval%20%2B%20Structured%20Concurrency.md)

