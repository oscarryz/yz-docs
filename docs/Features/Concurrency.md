
#feature


# Concurrent by Default

Yz is concurrent by default. Every boc invocation runs asynchronously, sharing the same memory space.

```yz
// Both run concurrently
foo()
bar()
```

The result of an invocation can be assigned to a variable, passed as an argument, or stored in an array immediately — without blocking. The value is resolved lazily: the flow only waits when the value is actually needed, typically at an IO boundary.

```yz
outer: {
    // These three invocations launch concurrently
    rv: foo()
    bar(rv)
    s: rv.to_str()

    // Flow blocks here — waits for `s` (which waits for `foo` and `to_str`)
    print("The value is: `s`")
}
```

## Structured Concurrency

A boc does not complete until all bocs it launched have completed. This is structured concurrency: the parent boc's lifetime encloses all its children.

```yz
parent_boc: {
    foo()
    bar()
}
// parent_boc() completes only when both foo() and bar() have finished
parent_boc()
```

## Summary

- Every boc call is async.
- Assigning the result to a variable (or passing it as an argument) does not block.
- Using the value — typically through IO — waits until the value is ready.
- A boc completes only after all bocs it launched have completed.
- Calls to the same object are executed in the order they are received (single-writer principle).

If a specific ordering is needed, wait for the result explicitly by using it before proceeding.

## Example

```yz
handle: {
    the_user: find_user()
    the_order: find_order()
    // Both fetches run concurrently; Response() waits for both values
    Response(the_user, the_order)
}
```

## Timeout

Using `return` (non-local return) exits the enclosing boc. This can implement a timeout by launching a timer boc that calls `return` after a delay:

```yz
fetch: {
    id String
    // Launches immediately; after 10 seconds, exits `fetch` with None
    time.sleep(10.seconds(), {
        return Option.None()
    })
    // Also launches immediately; if it finishes first, returns the result
    return find(id)
}
```

See: [Go - Go Concurrency Patterns](../Examples/Go%20-%20Go%20Concurrency%20Patterns.md)
