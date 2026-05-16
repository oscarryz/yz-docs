#rejected
#solved

# Stateless bocs and pure functions

## The problem

All bocs in Yz are actors — they have persistent fields and calls are serialized through their queue. This works well for objects and singletons. But it creates a fundamental problem for utility functions:

```yz
add: { a Int; b Int; a + b }
```

`add` is a singleton. Concurrent callers queue. Two goroutines cannot call `add` in parallel — they serialize. This is correct for an actor, but wrong for a mathematical function with no meaningful shared state.

## Why this question is rejected

The question assumed that "stateless" required a new syntactic construct or keyword to opt a boc out of the actor model. It doesn't. The existing design already provides a complete answer through two complementary mechanisms:

### 1. The Uppercase convention

Use an Uppercase boc for anything that should produce independent instances:

```yz
Add: { a Int; b Int; a + b }

result1: Add(3, 4)()   // fresh instance, independent
result2: Add(5, 6)()   // fresh instance, runs concurrently
```

Each `Add(...)` call creates an independent actor. No shared queue. No serialization. This is the language's built-in answer to "parallel-safe callable."

### 2. The thin dispatcher pattern for library code

For utility functions that need a familiar lowercase call site (`while`, `each`, `map`), a lowercase singleton dispatches to a fresh Uppercase worker:

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

`while` is a thin singleton — it processes each call in microseconds (non-blocking thunks), spawning an independent `While` instance for the real work. Concurrent callers don't contend for the actual execution.

### 3. Thunks prevent deadlock in recursion

The original deadlock concern ("a recursive singleton would block on itself") is resolved by thunks. `fact(n-1)` inside a singleton `fact` does not block — it returns a thunk and enqueues the next call. The queue drains naturally:

```
fact(5) → enqueues fact(4), returns thunk T5 = 5 * T4
fact(4) → enqueues fact(3), returns thunk T4 = 4 * T3
...
fact(0) → returns 1, T0 materializes
T1, T2, T3, T4, T5 all resolve
```

No deadlock. The thunk model was already solving this.

## What was actually right in the original question

The throughput concern is real: a singleton utility called from many concurrent contexts serializes all callers. The answer is not a new keyword but a usage pattern:

> **If you need parallel execution, use Uppercase. Singletons are for shared state, not parallel computation.**

Library authors writing hot utilities use the thin dispatcher pattern (lowercase + Uppercase worker). Application code that accidentally uses a lowercase singleton for concurrent computation gets serialization — noticeable in practice, fixable by switching to Uppercase.

## Options considered and why they were rejected

**Option A (stateless keyword)**: not needed — the Uppercase convention already provides fresh instances.

**Option B (compiler infers from field access)**: fragile — adding `add.a` anywhere silently changes the concurrency model of `add`.

**Option C (Uppercase = fresh instance)**: this IS the answer, already in the language.

**Option D (BocWithSig `#(...)` as the function form)**: partially right — `#(...)` params are local to each call — but it encodes the distinction as a naming convention inside the signature rather than at the boc level. Subsumed by the clearer Uppercase convention.

## Related

- [How to use utility functions](How%20to%20use%20utility%20functions.md)
- [Concurrency](../../Features/Concurrency.md)
- [Single Writer](Single%20Writer.md)
- [Bocs](../../Features/Bocs.md)
