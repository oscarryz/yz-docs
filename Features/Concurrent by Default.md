# Concurrent by Default 
Yz is concurrent by default, every method invocation will run in its own coroutine asynchronously, although sharing the same memory space.
When a variable is assigned to a method call Yz will wait until that method completes, turning the code into a synchronous call. 
The block will finish until all the method calls it created finishes (structural concurrency)

- Every method call is async
- Assigned as variable (or used as argument for another method) the execution will wait
- Once all the calls finishes, the method is self will finish. 

### Synchronous example
Here's a regular synchronous program. All the method invocations are assigned to a value and thus executed in sequential order
```javascript
'Will be executed synchronized, because `restaurant` cannot take orders or close before it opens'
main: {
    restaurant: open_restaurant()
    restaurant: restaurante.take_orders()
    restaurant.close()
}
// same program 
main: {
   open_restaurant()
   .take_orders()
   .close()
   // main would finish at this point
}
```
### Asynchronous example

We can use the `open_restaurant` method as an example of concurrent execution
```javascript
'Create a `Restaurant` instance and get everything ready'
open_restaurant: {
    r: Restaurant()
    r.start_kitchen()
    r.prepara_tables()
    r.open_register()
    r
}
```
In the previous example all the actions happen concurrently as they don't depend on each other. 


### Syncing asyncs 
Most of the time we require coordination of async calls after they complete. 
We can let the methods run and then sync by requesting their values

```javascript
retrieve_order: {
    customer_id  Int
    product_id Int
    fetch_customer( customer_id ) // no assign, so will run async
    fetch_product( product_id )  // no assign, so will run async
    Order {
       customer: fetch_customer.customer  // assigning, thus 'awaits' for methos completion
       product : fetch_product.product    // assigning, thus 'awaits' for methos completion
    }
}    
```

Example modified from [Structured Concurrency in Java](https://openjdk.org/jeps/428#:~:text=For%20example%2C%20in%20this%20single%2Dthreaded%20version%20of%20handle()%20the%20task%2Dsubtask%20relationship%20is%20apparent%20from%20the%20syntactic%20structure%3A)
```javascript
//Yz sequential
handle: {
    the_user: find_user()
    the_order: find_order()
    Response { the_user the_order}
}
// Yz concurrent
handle: {
    find_user()
    find_order()
    if some_logic() {
        other_thing() // maybe this is executed first
    }
    Response { find_user.user find_order.order }
}
```

Timeout: 
Create a block that exits the current block by calling `return` after the specified amount of time. How to clean up resources is still TBD

```javascript
fetch: {
    id String
    // creates a timeout
    time.sleep 10.seconds() {
        // after 10 seconds just return finilizing the `fetch` execution
        return
    }
    data: find(id)
    return // force return to avoid waiting for the sleep
}
```

Error handling
A couple of options:
1. Have the async function return an Option/Result kind of data 
2. Check the data to see if it's an error type
3. Check the data to see if it's valid

```javascript
fetch:{
    id String
    timeout(10 {return})
    // option 1
    data: find(id).or { "Couldn't find data"} 
    // option 2
    data err: find(id)
    err ? { "Couldn't find data" } { data } 
    // option 3
    data: find(id)
    data == something.invalid ? { "Couldn't find data" } {data}
    
}
```


See: [Go - Go Concurrency Patterns](Go%20-%20Go%20Concurrency%20Patterns.md)
