#answered Use them, and try to make them short-lived (thin dispatchers) so others can run.

# How to use utility functions

## The core rule

**Uppercase = fresh independent execution. Lowercase = one shared thing.**

Use Uppercase bocs for anything called multiple times or concurrently. Use lowercase singletons for shared state you intend to be singular.

```yz
// Do: Uppercase boc — each call is a fresh, independent actor
Foo: {
    bar: {}
}
f: Foo()
f.bar()

// Do: create and call immediately
f: Foo()
f()

// Don't: lowercase singleton called expecting parallel execution
foo: {
    print("a")
}
foo()  // called a million times — all serialize through one queue
```

## Why this works

Every `Foo()` call creates a fresh actor instance with its own queue. Concurrent callers never contend. Each execution is fully isolated.

Every `foo()` call on a lowercase singleton sends to the same actor queue. Callers serialize. This is correct for shared state (a counter, a connection pool) and wrong for pure utility work.

The naming convention already communicates intent. Uppercase signals "I produce instances." Lowercase signals "I am the one shared thing."

## Long-lived utilities: the thin dispatcher pattern

For truly shared library code (`while`, `each`, `map`, `filter`), a lowercase singleton that dispatches to a fresh Uppercase worker solves the contention problem:

```yz
while: {
    cond #(Bool)
    action #()
    While: {
        cond #(Bool)
        action #()
        cond() ? { action(); While(cond, action)() }
    }
    While(cond, action)()
}
```

`while` is a singleton but does almost no work — it receives the call and immediately spawns a fresh `While` instance. With thunks (non-blocking calls), `while` finishes processing each message in microseconds and moves to the next caller. Both `While` instances run concurrently and independently.

No deadlock: each recursive `While(cond, action)` call creates a fresh instance, so there is no singleton queue to deadlock against.

No corruption: each `While` instance owns its own `cond` and `action` fields. Concurrent callers never see each other's state.

Note: there is no need for a named `apply` method — `w()` calls the boc body directly. The boc is the apply.

## Why lowercase `w: { ... }` inside `while` doesn't work

The naive approach:

```yz
while: {
    w: {
        cond #(Bool)
        action #()
        cond() ? { action(); w(cond, action) }
    }
    w()
}
```

`w` is a field of the `while` singleton. When `foo` and `bar` call `while` concurrently:

- `foo`'s call sets `while.w` to foo's boc, calls `w()` (non-blocking thunk)
- `bar`'s call sets `while.w` to bar's boc — **overwrites foo's `w`**
- foo's loop recursion now references bar's boc

The Uppercase worker pattern avoids this because each `While(...)` call produces a distinct object, not a reassignment of a shared field.

## The honest trade-off

`Foo()` creates a goroutine and channel per call. For utilities called millions of times in a tight loop, that overhead is real compared to a plain function call. The `fn` keyword (when available) avoids this by running in the caller's context with no allocation. For application code the overhead is negligible; for hot library utilities it matters.

## Related

- [Stateless bocs and pure functions](Stateless%20bocs%20and%20pure%20functions.md)
- [Concurrency](../../Features/Concurrency.md)
- [Single Writer](Single%20Writer.md)
