#feature 
# Single Writer

Every boc is the **sole writer** of its own variables. External code can read another boc's fields freely, but cannot write to them directly — writes from outside are queued through the owner's actor.

This is the **Single Writer, Multiple Reader (SWMR)** model.

## Reads and writes

| Operation | Who can do it | How |
|---|---|---|
| Read `a.b` | Any boc | Direct — no queue, no thunk overhead |
| Write `a.b = v` from inside `a` | `a` and its nested bocs | Direct |
| Write `a.b = v` from outside `a` | Any boc | Queued through `a`'s actor |

```yz
counter: {
    count: 0
    increment: { count = count + 1 }  // direct write — inside counter
}

other: {
    counter.count          // direct read — fine from anywhere
    counter.count = 5      // queued write — sent to counter's actor, not applied immediately
    counter.increment()    // also queued — standard actor call
}
```

## Why SWMR (not full queue for reads)

If reads also went through the actor queue, every `a.b` access would be async, returning a `*Thunk[B]` instead of a `B`. That means:

- `print(counter.count)` would require queueing a read, waiting for a thunk, then forcing it — enormous overhead for trivial access
- `x: a.b; y: a.c` would be two separate async operations with no ordering guarantee between them

SWMR avoids this: reads are cheap and synchronous. The tradeoff is that a read might see a slightly stale value if a queued write hasn't been applied yet. This is the standard eventual-consistency tradeoff of lock-free concurrent systems, and is acceptable for nearly all use cases.

## What this solves

**No torn writes.** A reader never sees a value that is mid-write (partially updated). The owner applies writes atomically from its own goroutine.

**No data races on mutation.** Two goroutines can never simultaneously write the same field. All writes are serialized through the owner.

**The field-mutation + call pattern is safe (but not atomic).** Even when external code does:

```yz
hi.text = "Goodbye"        // queued write
hi.recipient = "everybody" // queued write
hi()                       // queued call
```

All three operations are ordered through `hi`'s queue. However, they are **not atomic as a group** — another actor can interleave operations between them. When atomicity matters, prefer the single-call form:

```yz
hi("Goodbye", "everybody")  // atomic: one queued operation
```

## Multi-field reads are not atomic

SWMR does not guarantee consistency across multiple reads:

```yz
x: rect.width    // direct read
y: rect.height   // direct read — another goroutine may have written between these two
area: x * y      // potentially inconsistent
```

For a consistent multi-field snapshot, expose a method that reads both fields within the same actor execution:

```yz
rect: {
    width Int
    height Int
    dimensions #(Int, Int) {
        width, height   // both read in the same actor turn — consistent
    }
}
x, y = rect.dimensions()
```

## Connection to structured concurrency

The single-writer guarantee is enforced by the runtime's actor queue, not by the compiler. The concurrency model (see [Concurrency](Concurrency.md)) ensures calls to the same boc are serialized. Field writes from outside are treated as implicit calls and join the same queue.
