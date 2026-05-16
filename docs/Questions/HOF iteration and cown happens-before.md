#open-question

# HOF Iteration and Cown Happens-Before Semantics

## The Question

When a higher-order iteration method like `Range.do()` or `Int.times().do()` calls a user-supplied closure, and that closure contains boc method calls on external cowns, what ordering guarantee does the caller get?

```yz
5.times().do({
    counter.increment()
})
```

Inside the closure, `counter.increment()` acquires `counter`'s cown and runs. Cown acquisition is guaranteed correct regardless of how `do` is implemented. The open question is **iteration order vs spawn order**.

## Two Possible Semantics

### A — Sequential (force each closure before next)

`do` treats each closure call as synchronous: it forces the returned thunk before calling the next iteration. Result: iterations are fully ordered. `counter` ends up incremented exactly N times, in order, before `do` returns.

```yz
Range: {
    do #(body #()) {
        // pseudocode — body() returns a thunk; force it before continuing
        body().force()   // wait for this iteration before next
        do(body)         // recurse
    }
}
```

**Consequence:** closure side-effects on cowns are totally ordered. Slow closures block the iterator.

### B — Concurrent (fire-and-forget, BocGroup)

`do` spawns each closure call into a BocGroup and waits at the end. All closure invocations are queued on the cown concurrently; the cown's queue serializes them in spawn order, which equals iteration order (since spawns happen sequentially in the loop). `do` returns only after all have completed.

```yz
Range: {
    do #(body #()) {
        // spawn all, wait at end — structured concurrency
    }
}
```

**Consequence:** cown serialization preserves spawn order, so `counter` increments in iteration order anyway. Parallel closures on *different* cowns run concurrently — better throughput.

## Key Insight

Option B gives the same *observable result* as A for single-cown closures, because the cown queue serializes in spawn order. But for closures that touch *multiple cowns* or do *non-cown work* (e.g. I/O), the execution is overlapped in B and strictly sequential in A.

Option B is more consistent with the rest of the BOC model (structured concurrency, spawn-and-wait).

## Related Concerns

- **Closure body BocGroup**: the closure itself (lowered by `lowerClosureBody`) already wraps statement-position boc calls in an implicit BocGroup and waits before returning. So a single `counter.increment()` inside the closure completes before `fn()` returns regardless of which option `do` uses.

- **Non-native Int methods**: `times()`, `to()`, `clamp()`, and similar methods defined in Yz (after uppering) are subject to the same design choice. They are not themselves cowns, but the closures they accept may capture cowns.

- **`each` on Array**: the existing runtime `Array.Each` has the same question deferred — it currently calls the closure synchronously in a Go range loop.

## Open Sub-Questions

1. Should Yz expose both semantics (e.g. `.do()` = sequential, `.spawn()` = concurrent)?
2. Should the compiler detect when a closure captures a cown and automatically adjust the iteration strategy?
3. Does the answer change for `map()`, where the closure returns a value? (Map likely must be sequential to build the result array in order.)
