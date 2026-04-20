#feature

# Gotchas

Common mistakes and the patterns that fix them.

---

## Calling a singleton expecting parallel execution

**The mistake:**

```yz
add: {
    a Int
    b Int
    a + b
}

// From two concurrent bocs:
x: add(3, 4)
y: add(5, 6)   // queues behind x — not parallel
```

`add` is a lowercase singleton. Every caller queues through the same actor. `x` and `y` do not run concurrently — they serialize.

This rarely causes correctness problems (thunks prevent deadlock, results are still correct), but it is a throughput bottleneck when concurrency matters.

**The fix: use an Uppercase boc**

```yz
Add: {
    a Int
    b Int
    a + b
}

x: Add(3, 4)   // fresh instance
y: Add(5, 6)   // fresh instance — runs concurrently with x
```

Each `Add(...)` call produces an independent actor. No shared queue. No serialization.

**The rule:**

> Lowercase = one shared thing. Uppercase = fresh independent execution.

If you find yourself calling a lowercase boc from many concurrent contexts and expecting parallelism, switch to Uppercase.

---

## Using a singleton for long-running work

**The mistake:**

```yz
worker: {
    while({ jobs.has_next() }, {
        process(jobs.next())
    })
}

worker()   // one queue — all work serializes
worker()   // queues behind the first call
```

**The fix: create instances**

```yz
Worker: {
    while({ jobs.has_next() }, {
        process(jobs.next())
    })
}

w1: Worker()
w1()   // independent actor
w2: Worker()
w2()   // independent actor — runs concurrently
```

---

## Implementing a reusable utility (the thin dispatcher pattern)

Library utilities (`while`, `each`, `map`) need a lowercase call site but must not serialize concurrent callers. The pattern: a thin singleton dispatcher that spawns a fresh Uppercase worker per call.

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

`while` is a singleton but does almost no work — it receives the call and immediately spawns a fresh `While` instance. With non-blocking thunk calls, `while` finishes each message in microseconds. Concurrent callers each get their own `While` actor running independently.

The recursive step `While(cond, action)()` creates a fresh instance for each iteration — no shared singleton to deadlock against.

Note: there is no need for a named method like `apply` or `run`. Calling `While(cond, action)()` invokes the boc body directly. The boc is its own entry point.

---

## Expecting recursion to deadlock

**Not a problem.** Thunks make boc calls non-blocking. A recursive singleton does not deadlock:

```yz
fact: {
    n Int
    n > 0 ? { n * fact(n - 1) }, { 1 }
}
```

`fact(n-1)` enqueues the next call and returns a thunk immediately — it does not wait. The singleton's queue drains naturally (n=5, n=4, ..., n=0), then the thunk chain resolves in reverse. Correct result, no deadlock.

Recursion through Uppercase bocs is even cleaner (fresh instance per call, no shared queue at all), but singleton recursion works too.

---

## Modifying a boc's fields from outside and then calling it

**The mistake:**

```yz
hi.text = "Goodbye"
hi.recipient = "everybody"
hi()
```

These are three separate messages in `hi`'s queue. Another actor can interleave between them. The call may see stale field values if something else writes `hi.text` between your write and your call.

**The fix: pass arguments atomically**

```yz
hi("Goodbye", "everybody")   // one message — atomic
```

One queued operation. No interleaving possible.

See [Single Writer](Single%20Writer.md) for the full model.
