---
type: resolved
generated: { by: "oscarryz", at: 2026-05-20T01:00:05+02:00 }
---
#resolved

# While Loop Yield and External Caller Interleaving

## The Problem

`docs/Features/Concurrency.md` contains a producer-consumer example that relies on `while`
iterations interleaving with an external caller:

```yz
boring: {
    m String
    messages: [String]()
    next: { messages.pop() }

    i: 1
    while({ true }, {
        messages.push("${m} ${i}")
        i = i + 1
        time.delay(1)
    })
}

main: {
    boring("sync")
    5.times().do({
        print(boring.next())
    })
}
```

The timing diagram shows `boring.next()` interleaving between while iterations. This cannot
happen with the current implementation.

## Why It Cannot Happen

`while` is implemented as boc recursion. The next iteration is spawned *inside* the current
iteration's body, which registers it as a **successor** in the cown queue (ScheduleAsSuccessor,
Phase E.1). Successors run before any externally-queued callers.

```
boring("sync")           ← main spawns this first
  └─ iteration 1         ← successor, runs before external callers
       └─ iteration 2    ← successor of 1, same
            └─ iteration 3
                 ...
boring.next()            ← queued by main after boring("sync"), never reached
boring.next()
...
```

For an infinite `while({ true }, ...)` the external callers never run.

## Why ScheduleAsSuccessor Exists

The BOC paper's formal model appends all spawned behaviours to the **tail** of the pending
queue — the opposite of ScheduleAsSuccessor. That works in the paper because the underlying
language is synchronous: a running behaviour can read and write its cowns directly (e.g.
`src.balance -= amount` is a plain field mutation). No async dispatch is needed for inner
operations.

Yz is different: **every call is async**. `src.balance-=(amount)` inside `transfer` is itself
a boc invocation. Without ScheduleAsSuccessor, an external `transfer(bob, alice)` spawned
by main could jump ahead of the deposit that the first transfer was in the middle of, breaking
atomicity.

ScheduleAsSuccessor was added to preserve the atomicity the paper gets for free from
synchronous execution. It is correct for the transfer case. The producer-consumer failure is
a consequence of applying it uniformly to recursive calls too.

## Resolution: Recursive → Tail, Non-Recursive → ScheduleAsSuccessor

The rule that resolves the conflict:

- **Non-recursive inner call** (callee FQN ≠ caller FQN) → ScheduleAsSuccessor  
  Preserves atomicity. External callers cannot jump in between related operations.

- **Recursive inner call** (callee FQN = caller FQN) → tail enqueue  
  Allows external callers to interleave between iterations. Infinite loops become
  interruptible.

Detection is compile-time: the lowerer compares the callee's FQN against the enclosing
boc's FQN at each call site.

### Why this works for transfer

`balance-=` and `balance+=` are different bocs from `transfer` — non-recursive →
ScheduleAsSuccessor. They remain successors of transfer, preserving atomicity.

### Why this works for producer-consumer

`while` calls itself — recursive → tail enqueue. The queue on `boring`'s cown becomes:

```
boring("sync")  →  while_iter1  →  boring.next()  →  boring.next()  →  while_iter2  →  ...
```

`while_iter2` is tail-enqueued after `boring.next()` calls that were already in the queue.
Producer and consumer naturally alternate.

## Design Principle

> **Direct recursion = producer pattern. HOF / indirect recursion = processor pattern.**

A boc that uses `while` (direct recursion) to produce values yields between iterations,
allowing consumers to interleave. A HOF like `each` or `map` uses `a→b→a` indirect
recursion — each call to the closure body and the next iteration are both
non-recursive from their call sites — so they use ScheduleAsSuccessor and run
sequentially without external interleaving. That is correct behaviour for processing:
you want all iterations to complete without interruption.

If a user wants to produce values that external callers can consume between iterations,
they must use direct recursion (`while` or an explicit recursive boc).

## Scenario Analysis

### 1. Direct recursion — `while`, recursive boc

```yz
while({ true }, { messages.push(...) })   // while calls while
```

`while → while`: recursive → tail enqueue. External callers interleave between
iterations. Correct for producers, event loops, and any pattern that needs to yield.

### 2. HOF iteration — `each`, `map`, `filter`

```yz
[1, 2, 3].each({ i Int; process(i) })
```

`each` calls the closure body (non-recursive → ScheduleAsSuccessor), then calls itself
for the next element (recursive → tail enqueue). The closure runs as a successor before
the next external caller, but the next iteration is tail-enqueued. Each closure body
completes before external callers can interleave with the iteration sequence. Sequential
processing behaviour — correct for non-producing loops.

### 3. Indirect recursion — ping-pong, `a → b → a`

```yz
ping: { pong() }   // ping calls pong — non-recursive → ScheduleAsSuccessor
pong: { ping() }   // pong calls ping — non-recursive → ScheduleAsSuccessor
```

Both calls are non-recursive from their individual call sites. The chain
`ping → pong → ping → pong → ...` hogs the cown: external callers never interleave.

This is acceptable because a ping-pong protocol is not producing values for external
consumers — it is doing coordinated internal work that should run atomically from the
outside's perspective. The sequential behaviour is correct.

### 4. Fan-in / fan-out — closure calling back into enclosing boc

```yz
coordinator: {
    results [Result]
    workers.each({ w Worker
        r: w.process()
        results.append(r)   // calls back into coordinator's cown — indirect recursion
    })
}
```

`coordinator → closure → coordinator.results.append` — each call is non-recursive from
its own site → ScheduleAsSuccessor at every step. The chain hogs `coordinator`'s cown
while all workers report back. External callers on `coordinator` wait until all results
are collected. This is correct: a partial results list is not a valid observable state.

### 5. Server accept loop

```yz
server: {
    while({ true }, {
        client: accept()     // fresh cown per connection
        handle(client)       // handler needs only client's cown, not server's
    })
}
```

`while → while`: recursive → tail enqueue. `handle(client)` shares no cown with
`server` — ScheduleAsSuccessor has nothing to insert it on for the server cown, so
`handle` is independently scheduled on `client`. `while_iter2` is tail-enqueued on
`server`. Both run in parallel. The server keeps accepting while handlers run
concurrently. Correct behaviour for an async server loop.

### 6. Indirect recursion through shared cown (known limitation)

```yz
a: { x Int; b(x) }   // a calls b on cown x — non-recursive → ScheduleAsSuccessor
b: { x Int; a(x) }   // b calls a on cown x — non-recursive → ScheduleAsSuccessor
```

If `a` and `b` share cown `x` and call each other indefinitely, external callers on `x`
never get in — same as ping-pong. Compile-time detection requires full call-graph cycle
analysis, which is not implemented. This is a known limitation of the FQN-comparison
approach. In practice, indirect cycles on the same cown that run indefinitely are rare
and are typically expressing coordinated internal work (acceptable hogging) rather than
a producer pattern.

## Stack Overflow

Not a concern. Each boc invocation is a separate goroutine. `while` calling `while`
(tail enqueue) means:

- Iteration N runs in goroutine G₁, finishes, G₁ exits
- Iteration N+1 runs in a new goroutine G₂

No call stack accumulates across iterations. Each goroutine has a shallow stack covering
only one iteration's work. An infinite `while` loop creates an unbounded number of
goroutines over time but never a deep stack. Go's goroutine scheduler handles this
efficiently.

## Impact on Concurrency.md

The producer-consumer example in `docs/Features/Concurrency.md` becomes valid once this
rule is implemented. The timing diagram is correct under the new scheduling semantics —
`boring.next()` does interleave between while iterations as shown.

## Implementation Note

The lowerer must detect recursive call sites before emitting Schedule vs
ScheduleAsSuccessor. Indirect cycles (`a → b → a`) are not detected by FQN comparison
alone — full call-graph cycle analysis would be required to handle them. That analysis
is deferred; indirect cycles are accepted as a known limitation.

## Related

- [Concurrency](docs/Features/Concurrency.md) — producer-consumer example (valid once resolved)
- [HOF iteration and cown happens-before](docs/Questions/HOF%20iteration%20and%20cown%20happens-before.md)  — YZC-0010 (resolved by this: HOF uses indirect recursion → sequential)
- [compiler/examples/transfer](compiler/examples/transfer/main.yz) — correct use of ScheduleAsSuccessor
- `compiler/runtime/rt/` — ScheduleAsSuccessor, Schedule implementations
