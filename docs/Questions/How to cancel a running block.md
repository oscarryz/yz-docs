#open-question

# How to cancel a running boc?

## The core tension

Non-local `return` inside a callback is intended to exit the enclosing boc early — most importantly for timeout patterns:

```yz
fetch: {
    id String
    time.sleep(10.seconds(), {
        return Option.None()   // non-local return: exits `fetch`
    })
    return find(id)            // also exits `fetch`
}
```

Both the timer callback and `find(id)` fire concurrently. Whichever `return` executes first wins. This is a **"race to return"** pattern.

But this conflicts with structured concurrency: "a boc does not complete until all bocs it launched have completed." Non-local early return exits the parent before its children finish.

### Three problems this creates

1. **Goroutine leak**: if `find(id)` wins at 5s, the timer goroutine is still sleeping for 5 more seconds with nobody waiting for it.

2. **Escaped non-local return**: if the timer wins and `fetch` completes, `find(id)` eventually finishes and tries to non-local return into a boc that is already done. The runtime needs a policy for this (silent discard is probably correct).

3. **Structured concurrency violation**: the parent boc exits before all children complete.

## Approaches explored

### Cooperative cancellation (current direction)

The boc being cancelled must periodically check a flag:

```yz
f: {
    keep_running: true
    inner: {
        while { keep_running }, {
            print("Working")
            time.sleep(1)
        }
    }
    cancel: {
        keep_running = false
    }
    inner()
}
f()
...
f.cancel()
```

Pro: simple to understand. Con: only works if the boc cooperates (checks the flag). Long-running operations that don't check it cannot be cancelled.

External stop signal (cooperative, injected):

```yz
f: {
    s #(stop_requested #(Bool))
    while { true }, {
        s.stop_requested() ? { return }
        print("Working")
        time.sleep(1)
    }
}
```

See also: [Hylo cancellation](https://docs.hylo-lang.org/language-tour/concurrency#cancellation)

### Implicit nursery / context (not yet designed)

An implicit cancellation token (similar to Go's `context.Context` or a structured concurrency nursery) could be threaded through all calls automatically. When a boc completes via non-local return, the nursery signals all sibling goroutines to stop.

- Pro: goroutines don't leak; structured concurrency is preserved
- Con: requires cooperative checking at every step; adds runtime complexity; long-running operations must be written to be cancellable; implicit passing is "magic"

This was considered but never fully designed. It is the most principled solution.

### First-return-wins with silent discard (minimal)

The runtime atomically marks the parent boc as "returned." Subsequent non-local returns from escaped callbacks are silently discarded. Goroutines continue running until natural completion.

- Pro: simple runtime implementation
- Con: goroutines leak (bounded, but real); structured concurrency not preserved

## Status

**No design decision made.** The cooperative pattern works for explicit cancellation scenarios but does not solve the timeout/race-return case cleanly. The implicit nursery approach is the most correct but requires significant design work.

## Related

- [Concurrency](../Features/Concurrency.md) — structured concurrency semantics
- [Single Writer](../Features/Single%20Writer.md) — actor model and queue semantics
- [return, break, continue](../Features/return%2C%20break%2C%20continue.md) — non-local return semantics
- [Transactions or atomic units of work](Transactions%20or%20atomic%20units%20of%20work.md)
