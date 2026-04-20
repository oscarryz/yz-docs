#spec 
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

**Stateful boc instances** (body form, no `#(...)`) are actors with message queues. **Stateless bocs** (BocWithSig form, `foo #(params) { ... }`) are not actors — each call creates a fresh independent goroutine with no shared queue.

| Form | Actor? | Queue? | Parallel calls? |
|---|---|---|---|
| `foo: { field T; ... }` | Yes | Yes | No — serialized |
| `Foo: { field T; ... }` (instance) | Yes | Yes | Yes — each instance has its own queue |
| `foo #(param T) { ... }` | No | No | Yes — fully parallel |

The compiler may further optimize:
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

### Timeout Pattern and Non-Local Return

Non-local `return` from a callback (see §7.5) exits the enclosing boc early, which is intentionally in tension with structured concurrency:

```yz
fetch: {
    id String
    time.sleep(10.seconds(), {
        return Option.None()   // non-local: exits `fetch` early
    })
    return find(id)            // also exits `fetch`
}
```

Both goroutines race to return from `fetch`. The first `return` wins; subsequent non-local returns from escaped callbacks are silently discarded. The losing goroutine runs to natural completion.

**Open design question**: how to cleanly cancel the losing goroutine to avoid goroutine leaks. See [Questions/How to cancel a running block](../Questions/How%20to%20cancel%20a%20running%20block.md).

## 8.6 Single-Writer Principle (SWMR)

Yz uses the **Single Writer, Multiple Reader** model for field access:

| Operation | Who | How |
|---|---|---|
| Read `a.field` | Any boc | Direct — no queue, no async overhead |
| Write `a.field = v` from inside `a` | `a` and nested bocs | Direct |
| Write `a.field = v` from outside `a` | Any boc | Queued through `a`'s actor |

Only one goroutine ever writes a field — the field's owner. Reads are unrestricted and synchronous. This avoids torn writes without making reads expensive.

```yz
counter: {
    count: 0
    increment: { count = count + 1 }  // direct write — inside counter
}

// From outside: reads are direct, writes are queued
x: counter.count          // direct read — any goroutine, no queue
counter.count = 5         // queued write — sent to counter's actor
counter.increment()       // queued call — also through actor
```

**Tradeoff**: reads may see slightly stale values (a queued write hasn't applied yet). Multi-field reads are not atomically consistent — use a method that reads both fields within one actor turn when consistency is needed.

**Note**: The "modify then call" pattern is not atomic:
```yz
hi.text = "Goodbye"       // queued write
hi.recipient = "everyone" // queued write
hi()                      // queued call — three separate queue entries, not atomic
hi("Goodbye", "everyone") // prefer this — one atomic queue entry
```

See [Features/Single Writer](../Features/Single%20Writer.md) for the full specification.

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
