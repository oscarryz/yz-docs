#spec 
# 8. Concurrency

This chapter defines Yz's concurrency model: behaviour-oriented concurrency (BOC), transparent thunks, and structured concurrency.

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

## 8.4 Behaviour-Oriented Concurrency (BOC)

Yz's concurrency model is based on **Behaviour-Oriented Concurrency** (Cheeseman et al., OOPSLA 2023). Every singleton boc instance is a **cown** (concurrent owner) — a protected resource. Invocations that need a cown are called **behaviours**; they run only when they have acquired all required cowns.

### Cowns and Behaviours

- A **cown** is a singleton boc instance. It owns its fields exclusively.
- A **behaviour** is an invocation that must acquire one or more cowns before running.
- All required cowns are **acquired atomically** — a behaviour either gets all of them at once or waits.

```yz
counter: {
    count: 0
    increment: { count = count + 1 }
    get: { count }
}

counter.increment()   // behaviour — acquires counter's cown, runs body, releases
counter.increment()   // queued — waits until previous releases
print(counter.get())  // waits for previous to complete → prints 2
```

### Ordering Guarantee

Two behaviours that share at least one cown always run in **spawn order** — the one scheduled first runs first. Behaviours with no cown overlap run freely in parallel:

```yz
transfer(src, dst, 100)   // acquires src + dst atomically
check_balance(src)        // waits — shares src with transfer

transfer(s1, s2)          // runs in parallel with...
transfer(s3, s4)          // ...this (no shared cowns)
```

### Multi-Cown Acquisition

When a behaviour needs multiple cowns, it acquires them all at once. This prevents deadlock: there is no partial acquisition and no "acquire one, then wait for another":

```yz
sync #(b bank, l ledger) {
    b.balance = b.balance + 1    // atomic: both cowns held simultaneously
    l.total   = l.total   + 1
}
sync(bank, ledger)   // atomically acquires bank's and ledger's cowns
```

### Boc Forms and Cown Ownership

| Form | Owns a cown? | Behaviours serialized? |
|---|---|---|
| `foo: { field T; method: {...} }` | Yes — singleton cown | Yes |
| `Foo: { field T; ... }` (struct type) | Yes — one cown per instance | Yes per instance |
| `foo #(param T) { ... }` (boc declaration) | No | Fully parallel |

Boc declaration forms (`foo #(params) { body }`) are stateless — each call is an independent goroutine with no cown, no serialization.

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

## 8.6 Inter-Boc Communication

Bocs communicate by calling each other's methods. Each call is a behaviour scheduled on the target's cown:

```yz
producer: {
    consumer #(item String)    // Parameter: reference to consumer boc
    1.to(100).each({ i Int
        consumer.receive("item ${i}")
    })
}

consumer: {
    receive: {
        item String
        print("Got: ${item}")
    }
}

producer(consumer)   // producer schedules receives on consumer's cown
```

## 8.7 Concurrency Implementation (Go Backend)

| Yz Concept | Go Implementation |
|-----------|-------------------|
| Singleton boc (cown) | Struct with embedded `std.Cown` (lock-free atomic queue) |
| Method call (behaviour) | `std.Schedule(&self.Cown, func() T { ... })` |
| Multi-cown behaviour | `std.ScheduleMulti([]*std.Cown{...}, func() T { ... })` |
| Sub-boc call on held cown | `std.ScheduleAsSuccessor(&held.Cown, fn)` — inserted before external waiters |
| Thunk | `*std.Thunk[T]` — lazy wrapper forced at IO boundary |
| Materialization | `.Force()` call on `*Thunk[T]` |
| Structured concurrency | `std.BocGroup` + `WaitGroup` on child goroutines |
| Cown acquisition order | Atomic queue per cown; behaviours run when all queues grant |

## 8.8 Summary

```
Concurrency Model:
  Default:           All invocations are async (return *Thunk[T])
  Value wrapper:     Thunk (same type to the programmer, lazy internally)
  Materialization:   IO boundaries only (.Force())
  Ownership unit:    Cown (one per singleton boc instance)
  Behaviour:         Invocation that holds one or more cowns while running
  Acquisition:       Atomic — all cowns acquired at once or none
  Ordering:          Spawn order for behaviours sharing a cown
  Parallelism:       Behaviours with disjoint cowns run freely in parallel
  State safety:      Exclusive cown access — all field writes go through cown queue
  Deadlock:          Impossible — atomic multi-cown acquisition
  Data races:        Impossible — exclusive cown access
  Structured:        Parent waits for all spawned behaviours before completing
  Runtime:           Goroutines + lock-free cown queue (Go backend)
```

## 8.9 Theoretical Background

The concurrency model in Yz is a direct application of **Behaviour-Oriented Concurrency** developed at Imperial College London and Microsoft Research. In that model, protected resources are called **cowns** and asynchronous units of work are called **behaviours** — Yz uses the same terms and the same formal guarantees.

> Cheeseman et al. (2023). *When Concurrency Matters: Behaviour-Oriented Concurrency*.
> Proc. ACM Program. Lang. 7, OOPSLA2, Article 276.
> <https://marioskogias.github.io/docs/boc.pdf>
