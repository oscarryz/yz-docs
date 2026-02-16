# 8. Concurrency

This chapter defines Yz's concurrency model: async-by-default execution, the actor model, thunk materialization, and structured concurrency.

## 8.1 Design Principle

Yz is **concurrent by default**. Every boc invocation is asynchronous — it immediately returns a **thunk** that represents the future result. The value is materialized (forced) only when needed for **IO**.

This eliminates the need for `async`/`await` keywords, thread management, or explicit concurrency annotations.

## 8.2 Async Invocation

Every method/boc invocation returns immediately with a thunk:

```yz
result: expensive_computation(data)
// result is a thunk — computation is running concurrently
// execution continues immediately
```

The thunk is a value of the **same type** as the boc's return type. It is indistinguishable from a resolved value in most contexts.

### Chaining

Thunks compose naturally:

```yz
x: fetch_user("alice")         // thunk User
y: x.name                      // thunk String (depends on x)
z: y == "Alice"                // thunk Bool (depends on y)
r: z ? { "yes" }, { "no" }     // thunk String (depends on z)
```

No value is computed until materialization is triggered.

## 8.3 Materialization

A thunk is **materialized** (its computation is forced to completion and the result extracted) only at **IO boundaries**:

| IO Operation | Example |
|--------------|---------|
| Print to console | `print(value)` |
| Write to file | `file.write(data)` |
| Network send | `connection.send(message)` |
| Program exit | Top-level boc completes |

### Materialization Cascade

When an IO operation needs a value, all thunks in the dependency chain are materialized recursively:

```yz
x: fetch_user("alice")         // thunk
y: x.name                      // thunk (depends on x)
print(y)                        // MATERIALIZES: y → x → fetch_user
```

### Non-Materializing Operations

The following do **not** trigger materialization:

- Assignment (`z = x`)
- Passing as argument (`foo(x)`)
- Storing in a collection (`list << x`)
- Field access (`x.name`) — returns a new thunk
- Method calls (`x.to_string()`) — returns a new thunk
- Comparison (`x == y`) — returns a thunk Bool

## 8.4 The Actor Model

Every boc instance is an **actor**:

- It has its own **message queue** (channel)
- Messages (method calls) are processed **in order**
- Each actor runs on its own **green thread** (goroutine in Go backend)
- Actors do not share mutable state — communication is via message passing

### Actor Lifecycle

```
1. Boc is instantiated → goroutine + channel created
2. Method calls are sent as messages to the channel
3. Actor processes messages sequentially from its queue
4. Actor completes when its body finishes and all inner actors complete
```

### Message Ordering

Messages to a single actor are processed in **FIFO order**:

```yz
counter: {
    count: 0
    increment: { count = count + 1 }
    get: { count }
}

counter.increment()   // Message 1
counter.increment()   // Message 2
counter.increment()   // Message 3
print(counter.get())  // Message 4 → prints 3
```

### Actor Granularity

**Every** boc instance is an actor — regardless of whether it has mutable state, methods, or complexity. The compiler may optimize away actor overhead for:

- **Pure bocs** (stateless, no side effects) — inline execution
- **Literal bocs** (`{ 42 }`) — constant folding
- **Bocs invoked and immediately materialized** — synchronous execution

These optimizations are transparent to the programmer.

## 8.5 Structured Concurrency

A boc does not complete until **all of its inner bocs have completed**:

```yz
main: {
    a: slow_operation_1()    // Starts concurrently
    b: slow_operation_2()    // Starts concurrently
    c: slow_operation_3()    // Starts concurrently
    // main does not complete until a, b, and c are all done
}
```

This provides:

1. **Automatic cleanup** — no orphan goroutines
2. **Error propagation** — errors in inner bocs propagate to the parent
3. **Predictable lifetimes** — a boc's scope determines its children's lifetimes

### Timeout Pattern

To limit execution time, use the race pattern:

```yz
result: match {
    data: fetch_with_timeout() => data
}, {
    time.sleep(5000) => "timeout"
}
```

## 8.6 Single-Writer Principle

Since each boc is an actor with a sequential message queue, mutable state is inherently safe:

- Only one message is processed at a time per actor
- No concurrent mutation of the same field
- No locks, mutexes, or synchronization primitives needed

```yz
// This is safe — increment messages are processed sequentially
counter: {
    count: 0
    increment: { count = count + 1 }
}

// Even with concurrent callers:
1.to(1000).each({ i Int
    counter.increment()   // Each is a message, processed in order
})
```

## 8.7 Inter-Actor Communication

Actors communicate by calling each other's methods. These calls are messages sent to the target actor's queue:

```yz
producer: {
    consumer #(item String)    // Parameter: reference to consumer
    1.to(100).each({ i Int
        consumer.receive("item `i`")
    })
}

consumer: {
    receive: {
        item String
        print("Got: `item`")
    }
}

producer(consumer)   // producer sends messages to consumer
```

## 8.8 Concurrency Implementation (Go Backend)

| Yz Concept | Go Implementation |
|-----------|-------------------|
| Boc instance | Goroutine + channel |
| Method call | Message sent to channel |
| Thunk | `chan T` or lazy wrapper |
| Materialization | `<-chan` (channel receive) |
| Structured concurrency | `sync.WaitGroup` on child goroutines |
| Actor message queue | Buffered channel with sequential processing loop |

## 8.9 Summary

```
Concurrency Model:
  Default:           All invocations are async
  Value wrapper:     Thunk (same type, lazy)
  Materialization:   IO boundaries only
  Actor:             Every boc instance
  Message ordering:  FIFO per actor
  State safety:      Single-writer (sequential message processing)
  Structured:        Parent waits for all children
  Runtime:           Green threads (goroutines)
  Scheduler:         Preemptive
```
