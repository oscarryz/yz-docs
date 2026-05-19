#open-question

# While Loop Yield and External Caller Interleaving

## The Contradiction

`docs/Features/Concurrency.md` contains a producer-consumer example that relies on `while` iterations interleaving with an external caller:

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

The doc's timing diagram shows `boring.next()` interleaving between while iterations — `boring` releases the resource after one iteration and `main` pops a message before the next iteration runs.

This cannot happen with the current implementation.

## Why It Cannot Happen

`while` is implemented as boc recursion. The next iteration is spawned *inside* the current iteration's body, which makes it a **successor** in the cown queue (ScheduleAsSuccessor, Phase E.1).

Successors are inserted at the head of a cown's wait queue, ahead of any externally-spawned callers. For `boring`'s cown the queue looks like:

```
boring("sync")           ← main spawns this (position 0)
  └─ iteration 1         ← spawned inside boring("sync"), promoted ahead of external callers
       └─ iteration 2    ← spawned inside iteration 1, same
            └─ iteration 3
                 ...
boring.next()            ← main's first do-block, queued after boring("sync")
boring.next()            ← main's second do-block
...
```

Because `{ true }` always returns true, iteration N spawns iteration N+1 before returning. `boring.next()` is always behind all pending iterations. For an infinite while loop it never runs.

## Why Transfer Is Correct, Producer-Consumer Is Not

The transfer example relies on the same successor semantics to get atomicity:

```yz
transfer(src, dst, amount) {
    src.balance-=(amount)   // successor on src
    dst.balance+=(amount)   // successor on dst
}
```

`balance-=` and `balance+=` run as successors, ensuring a transfer completes atomically before another transfer touching the same account can run. Correct behavior.

For producer-consumer the desired behavior is the opposite: after one push, the consumer should be able to interleave before the next push. Successor semantics prevent that.

## The Core Design Question

Does `while` (and iterative structures in general) need a mechanism to yield to external callers between iterations?

If yes, what is that mechanism?

## Possible Approaches

### A — Tail scheduling for while body

The next iteration is scheduled at the **tail** of the cown's queue, not as a successor. External callers spawned before the next iteration's spawn-point run first (FIFO).

Consequence: `boring.next()` (spawned by main before any iteration is spawned) would run between iterations. Infinite loops become interruptible.

Risk: changes the atomicity model for while — operations inside two iterations of the same while body are no longer ordered with respect to each other.

### B — Explicit yield statement

A `yield` expression inside a boc body suspends the current invocation, releases the cown, and re-enqueues the continuation at the tail of the queue. Other callers run before the continuation resumes.

```yz
while({ true }, {
    messages.push("${m} ${i}")
    yield   // release cown, let next() in
    i = i + 1
})
```

Consequence: programmer controls interleaving explicitly. More expressive but requires a new language construct.

### C — Producer-consumer via a separate channel/mailbox primitive

Producer-consumer is a pattern that needs explicit synchronization primitives (channel, mailbox, ring buffer) rather than direct boc calls. `boring.next()` would block on a channel read; `boring` pushes to the channel. The while loop runs independently on the channel without needing to release `boring`.

Consequence: changes the example significantly; the doc's point about "falling out of the model naturally" no longer holds.

### D — Two-phase while: iteration spawning vs execution

The while driver runs outside the iterated boc — it does not hold the resource while spawning the next iteration. Each iteration runs as an independent invocation in queue order, not as a nested successor.

This is the most consistent with the rest of the model but requires redesigning how `while` is lowered.

## Relationship to YZC-0010

YZC-0010 asks whether HOF iteration (`Range.do`, `Int.times().do`) calls closures sequentially or concurrently. That question is about HOF method semantics on externally-supplied closures.

YZC-0036 is a deeper structural question: does any iterative structure (user-defined while loops included) block external callers permanently? The answer to YZC-0010 may depend on the answer to YZC-0036 — if "tail scheduling" (Option A) is adopted, HOF methods likely follow the same rule.

## Related

- `docs/Features/Concurrency.md` — producer-consumer example (currently incorrect)
- `docs/Questions/HOF iteration and cown happens-before.md` — YZC-0010
- `examples/transfer/main.yz` — correct use of successor semantics
