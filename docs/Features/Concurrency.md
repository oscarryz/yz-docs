
#feature

# Concurrency

Concurrency in Yz language is built on a single, unified model: **every value is a concurrent owner (cown), and every function call is am asynchronous behaviour**. There are no locks, explicit threads, async/await annotations, nor function coloring. The runtime handles the synchronization automatically, and the compiler optimizes away the overhead in the common case.

## Core Concepts

### Everything is a Cown

A cown (concurrent owner) is a protected piece of data that provides the only entry point to that data in the program. Every value in the language — whether a simple integer, an object, or a function — is implicitly a cown.

```js
counter: { value Int }   // counter is a cown, automatically
```

A cown is either **available** or **acquired**. Only one behaviour (or boc invocation) can hold a cown at a time. When a cown is acquired, all other behaviours that need it wait in a queue.

### Every Call is a Behaviour

A behaviour is an asynchronous unit of work. When you call a boc, you are spawning a behaviour that will run when it has acquired all the cowns it needs. Crucially, **all required cowns are acquired atomically** — a behaviour either gets all of them at once, or waits until it can.

```js
transfer(src, dst, 100)   // spawns a behaviour requiring {src, dst}
check_balance(src)         // spawns a behaviour requiring {src}
```

This means Yz runtime infers concurrency from data dependencies, not from annotations.

### Happens-Before Ordering

The runtime guarantees that if two behaviours share at least one cown, and one was spawned before the other, the earlier one always runs first. This is called the **happens-before** relation.

```js
transfer(s1, s2)     // b1: acquires {s1, s2}
check_balance(s1)    // b2: must wait — s1 is taken by b1
```

`check_balance` will always see the result of `transfer`. No explicit synchronization needed.

Behaviours that share **no** cowns have no ordering constraint and run freely in parallel:

```js
transfer(s1, s2)    // b1: acquires {s1, s2}
transfer(s3, s4)    // b2: acquires {s3, s4} — runs in parallel with b1
```

## Safety Guarantees

The model provides four guarantees by construction, not by programmer discipline:

**Data-race freedom** — since each cown's data is isolated and only one behaviour holds it at a time, two behaviours can never concurrently mutate the same state.

**Deadlock freedom** — because all cowns are acquired atomically (all at once or not at all), the circular waiting that causes deadlock cannot form. There is no "grab one lock, then wait for another."

**Determinism** — for any two behaviours sharing a cown, their execution order is always the same: the one spawned first runs first. A program's observable behaviour is deterministic with respect to spawn order.

**Structured lifetimes** — a block does not complete until all behaviours it spawned have completed, regardless of whether return values have already been delivered to the caller. Child behaviours cannot outlive their parent scope. See [Structured Concurrency](#structured-concurrency) below.

## Performance

Every function call going through a queue might sound expensive. In practice it is not, for two reasons.

First, an **uncontended cown** — one with no other waiters — is acquired instantly. The queue exists but is empty, so acquisition is a single atomic operation with no waiting.

Second, the compiler and the runtime can **inline and elide** cown acquisition entirely for values that are probably local to a single behaviour. Correctness is guaranteed by the model; performance is a pure optimization problem, decoupled from correctness.

The result: you write simple, safe code, and the runtime figures out what can run in parallel.

## Return Values and Suspension

Return values in Yz are **concurrently lazy**. When you call a function, it starts running immediately as a behaviour, and a transparent thunk is returned to the caller. The thunk is forced — and the caller suspends — only when the value is actually used.

```js
w : load("world")   // load starts running, w is a thunk
other()             // runs immediately, load still in progress
print(w)            // forced here — suspends until load completes
```

The suspension is transparent. You write sequential-looking code; the runtime handles the interleaving. This gives you future/promise semantics without any `Future[T]` type or `.then()` chain.

Importantly, a thunk being forced does not mean the spawning block is done. Even after `print(w)` receives its value, if `load` spawned further child behaviours internally, the enclosing block waits for all of them. Return values and scope completion are distinct — this is the structured concurrency guarantee described in the next section.

### Multiple Return Values

Functions can return multiple values. All outputs are resolved simultaneously when the behaviour completes — there is no partial resolution.

```js
value, errors : parse(input)    // suspends until parse finishes
                              // then value AND errors are assigned together

errors.len() == 0 ? {
    print(value) // value is available here
}, {
   print(errors)
}
```

Because functions and objects are the same thing, multiple return values are equivalent to returning a single anonymous object. These two forms are identical in semantics:

```js
// explicit destructuring
value, err : parse(input)

// attribute access
parse(input)
value : parse.value
err   : parse.err
```

## Structured Concurrency

A block does not complete until **all behaviours it spawned have completed**, even if those behaviours were spawned indirectly or if all return values have already been delivered to the caller.

```js
outer: {
    done Bool
    inner: {
        time.delay(10)
        print("inner done")
    }
    done = true // assigned immediately 
}
outer()
print("outer done `outer.done`")    // guaranteed to print AFTER "inner done"
```

This is the **scope guarantee**: the lifetime of every behaviour is bounded by the block that spawned it. Child behaviours cannot escape their parent scope. This prevents the classic failure mode of fire-and-forget tasks that leak resources, propagate errors silently, or produce results no one is listening to.

### Thunks Do Not Escape

Return values (thunks) are child behaviours. Forcing a thunk delivers the value to the caller, but does not close the child's scope — any work still running inside that child continues, and the parent scope waits for it.

```js
main: {
    result : compute()    // compute starts, thunk returned
    print(result)         // thunk forced — main gets the value
                          // but main does NOT exit here if compute
                          // has child behaviours still running
}                         // main exits only when all descendants done
```

This is what distinguishes structured concurrency from plain futures. A future in other languages can outlive its creator. In this language, it cannot.

### Scope and Parallelism Together

Structured lifetimes compose naturally with parallel execution. Multiple child behaviours run in parallel within a scope, and the scope waits for all of them:

```js
process: {
    a : step_one()      // runs in parallel...
    b : step_two()      // ...with this
    combine(a, b)       // waits for both, then runs
}                       // exits only after combine and all descendants finish
```

The happens-before ordering ensures `combine` waits for `a` and `b`. The scope guarantee ensures `process` waits for `combine`. Both properties hold simultaneously with no extra annotation.

### Producer Lifetimes

Structured concurrency has a direct consequence for long-running or infinite producers: a producer spawned inside a block is bounded by that block's scope, so the block cannot exit until the producer finishes.

```js
main: {
    boring("sync")          // spawns an infinite producer
    5.times().do({
        print(boring.next())
    })
}                           // main cannot exit — boring never finishes
```

Cancellation is not yet part of the model. For infinite producers the recommended pattern is to use independent instances with explicit launcher and reader blocks so their lifetime can be managed separately from the calling scope. See [Blocks and Instances](/docs/blocks-and-instances) and the Producer-Consumer section below.

## Ordering and Scheduling

The happens-before relation is determined by **spawn order and cown overlap**. This means the programmer controls parallelism by controlling when behaviours are spawned.

Consider four philosophers each eating once, where each needs two forks:

```js
// Sequential — each call overlaps with the next, forced into a chain
eat(f1, f2)   // b1
eat(f2, f3)   // b2 — waits for b1 (shares f2)
eat(f3, f4)   // b3 — waits for b2 (shares f3)
eat(f4, f1)   // b4 — waits for b3 (shares f4)

// Parallel — non-overlapping pairs spawned first
eat(f1, f2)   // b1
eat(f3, f4)   // b3 — no overlap with b1, runs in parallel
eat(f2, f3)   // b2 — waits for both b1 and b3
eat(f4, f1)   // b4 — waits for both b1 and b3
```

The second schedule runs two things in parallel in the first phase and two in the second. **Incorrect scheduling hurts performance but never correctness** — the program is always data-race free and deadlock free regardless of spawn order.

## Blocks and Instances

Concurrency composes naturally with the language's block and instance model. Independent instances each carry their own cowns, so they share no state and run fully in parallel. Each instance also forms its own structured scope — its child behaviours are bounded by its own lifetime, not the parent's, which makes independent instances the natural tool for managing long-running concurrent work. For a full explanation see [Blocks and Instances](/docs/blocks-and-instances).

## Producer-Consumer Communication

Because every value is a cown and every access goes through a queue, producer-consumer communication falls out of the model naturally — no channels required as a language primitive.

A producer writing to a shared array and a consumer reading from it are automatically ordered:

```js
boring: {
    m String
    messages: [String]()
    next: { messages.pop() }
    i: 1

    while({ true }, {
        messages.push("`m` `i`")
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

The happens-before guarantee ensures `boring` writes at least once before `main` reads (boring is spawned first). After that, boring can race ahead freely — since it writes to an array rather than a single value, no messages are lost. The consumer always reads in order via `pop`.

Because `boring` contains an infinite loop, the structured concurrency guarantee means `main`'s scope will not close until `boring` finishes. Until cancellation is supported, infinite producers should be structured as independent instances with explicit launcher and reader blocks so their lifetime can be managed independently.

For bounded producer-consumer with backpressure, the standard library provides `Chan[T]` — a cown wrapping a bounded queue with blocking push and pop semantics. This is a library type, not a language primitive; it is built entirely from the same cown mechanics described above.

For block signatures and type annotations see [Block Signatures](/docs/block-signatures).

## Summary

| Concept | How it works |
|---|---|
| Shared state | Everything is a cown — protected by default |
| Parallelism | Behaviours with non-overlapping cowns run in parallel automatically |
| Ordering | Behaviours sharing a cown run in spawn order |
| Deadlock | Impossible — all cowns acquired atomically |
| Data races | Impossible — exclusive access guaranteed |
| Async | Every call is async — suspension is transparent at return value use |
| Structured lifetimes | A block waits for all descendant behaviours before completing |
| Thunk escaping | Impossible — return values are child behaviours, bounded by their scope |
| Channels | Fall out of cown semantics — available as `Chan[T]` in stdlib |
| Performance | Uncontended cowns are inlined/elided by the compiler |

## Further Reading

The concurrency model in this language is directly inspired by **Behaviour-Oriented Concurrency (BoC)**, a concurrency paradigm developed at Imperial College London and Microsoft Research. The formal model, semantics, and implementation details are described in:

> Cheeseman et al. (2023). *When Concurrency Matters: Behaviour-Oriented Concurrency*. Proc. ACM Program. Lang. 7, OOPSLA2, Article 276. [https://marioskogias.github.io/docs/boc.pdf](https://marioskogias.github.io/docs/boc.pdf)