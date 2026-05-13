#feature 
# Yz Concurrency

## Overview

Yz is concurrent by default and is the only model. Every boc (block of code) invocation runs concurrently, lock free, data race free and in a deterministic way. There are no language constructs to drive concurrency like locks, threads, channels nor `async`/`await` annotations. The runtime handles
synchronisation automatically and the compiler optimises away overhead in the common
case.

---

## Concurrency Model

### Concurrent Resources

Every value in Yz — whether a simple integer, a boc, or a complex object — is
implicitly a protected concurrent resource. A resource is either **available** or
**acquired**. Only one running boc can acquire a resource at a time. When a resource is acquired,
all other bocs that need it wait in a queue.



### Every Invocation Is Concurent

When you invoke a boc, you are scheduling an asynchronous unit of work that will run
when it has acquired all the resources it needs. Crucially, **all required resources are
acquired atomically** — a running boc either gets all of them at once, or waits until
it can. There is no partial acquisition.

```js
transfer(src, dst, 100)   // acquires src and dst atomically before running
check_balance(src)        // waits if src is already acquired
```

### Happens-Before Ordering

The runtime guarantees that if two invocations share at least one resource, the one
spawned earlier always runs first.

```js
transfer(s1, s2)     // acquires {s1, s2}
check_balance(s1)    // must wait — s1 is taken by transfer
```

`check_balance` will always see the result of `transfer`. No explicit synchronisation
needed.

Invocations that share **no** resources have no ordering constraint and run freely
concurrently:

```js
transfer(s1, s2)    // acquires {s1, s2}
transfer(s3, s4)    // acquires {s3, s4} — runs in parallel with the first
```

---

## Safety Guarantees

The model provides four guarantees by design:

**Data-race freedom** — each resource is isolated and only one running boc holds it at a
time. Two invocations can never concurrently mutate the same state.

**Deadlock freedom** — because all resources are acquired atomically, the circular waiting
that causes deadlock cannot take place. There is no "acquire one resource, then wait for
another."

**Determinism** — for any two invocations sharing a resource, their execution order is
always the same: the one spawned first runs first. A program's observable behaviour is
deterministic with respect to invocation order.

**Structured concurrency** — a boc does not complete until all invocations it spawned
have completed, even if return values have already been delivered to the caller.
See [Structured Concurrency](#structured-concurrency) below.

---

## Performance

An **uncontended resource** — one with no other waiters — is acquired instantly. The queue
exists but is empty, so acquisition is a single atomic operation with no waiting.

 The runtime figures out what can run in parallel.

---

## Return Values

Return values in Yz are **transparently lazy**. When you invoke a boc, it gets scheduled to run as soon as its resources are acquired and a placeholder is returned to the caller. The placeholder
resolves — and the caller synchronizes — only when the value is actually used, typically in a I/O boundry. 

```js
u User = load("user:123")   // load scheduled to run immediately
other()                      // runs immediately, load still in progress
print(u)                     // synchronizes here until load completes
```

The type of `u` is always the declared return type — `User` in this case, not a
`Future[User]` or any wrapper type. The placeholder is an internal runtime concept,
not something the developer works with directly:

```js
load #(id String, User)   // returns User — always
```

The synchronization is transparent. You write sequential-looking code and the runtime handles
the interleaving.

This extends to boc calls that share resources with the caller. The return value is still
just the declared type — the runtime manages the ordering:

```js
balance #(acc Account) { acc.balance }

report #(acc Account) {
    b: balance(acc)      // b is Int — same resource, no special handling needed
    print("balance: ${b}")
}
```

---

## Structured Concurrency

A boc does not complete until **all invocations it spawned have completed**, even if
those invocations were spawned indirectly or if all return values have already been
delivered to the caller.

```js
outer : {
    done Bool
    inner : {
        time.delay(10)
        print("inner done")
    }
    done = true   // assigned immediately
}
outer()
print("outer done `outer.done`")   // guaranteed to print AFTER "inner done"
```

The lifetime of every invocation is bounded by the boc that spawned it. Child
invocations cannot outlive their parent scope. This prevents fire-and-forget tasks that
leak resources, propagate errors silently, or produce results no one is listening to.

### Return Values Do Not Close Child Scopes

Receiving a return value does not close the child's scope. Any work still running inside
that child continues, and the parent waits for it:

```js
main : {
    result : compute()   // compute starts, placeholder returned
    print(result)        // main gets the value here
                         // but main does NOT exit if compute has children still running
}                        // main exits only when all descendants are done
```
> _Note: Structured concurrency contexts for grouping and cancelling related invocations are under design._

---

## Scheduling And Parallelism

The ordering guarantee is determined by **invocation order and value overlap**. The
programmer controls parallelism by controlling when invocations are spawned.

Consider four philosophers each needing two forks:

```js
// Sequential — each invocation shares a fork with the next, forced into a chain
eat(f1, f2)
eat(f2, f3)   // waits for first  (shares f2)
eat(f3, f4)   // waits for second (shares f3)
eat(f4, f1)   // waits for third  (shares f4)

// Parallel — non-overlapping pairs spawned first
eat(f1, f2)   // runs in parallel...
eat(f3, f4)   // ...with this
eat(f2, f3)   // waits for both above
eat(f4, f1)   // waits for both above
```

The second schedule runs two invocations in parallel in the first phase and two in the
second.

> **Incorrect scheduling hurts performance but never correctness.** The program is
> always data-race free and deadlock free regardless of invocation order.

---

## Producer-Consumer

Because every value is protected and every access goes through a queue,
producer-consumer communication falls out of the model naturally.

A producer writing to a shared array and a consumer reading from it are automatically
ordered:

```js
boring : {
    m String
    messages : [String]()
    next : { messages.pop() }
    i : 1

    while({ true }, {
        messages.push("`m` `i`")
        i = i + 1
        time.delay(1)
    })
}

main : {
    boring("sync")
    5.times().do({
        print(boring.next())
    })
}
```

The ordering guarantee ensures `boring` writes at least once before `main` reads —
`boring` is spawned first. After that, `boring` can race ahead freely. Since it writes
to an array rather than a single value, no messages are lost.


---

## Summary

| Property | How Yz achieves it |
|---|---|
| Shared state safety | Every value is protected — exclusive access guaranteed |
| Parallelism | Invocations with non-overlapping values run in parallel automatically |
| Ordering | Invocations sharing a value run in spawn order |
| Deadlock | Impossible — all values acquired atomically |
| Data races | Impossible — exclusive access by construction |
| Async | Every invocation is async — suspension transparent at value use |
| Structured lifetimes | A boc waits for all descendant invocations before completing |
| Return value escaping | Impossible — return values are bounded by their scope |
| Cancellation | Via `cancel()` method on invocation handle |
| Performance | Uncontended values are inlined/elided by the compiler |

---

## Theoretical Background

The concurrency model in Yz is a variation of the Behaviour Oriented Concurrency model developed at Imperial College London and Microsoft Research. In that model, protected
resources are called **cowns** (concurrent owners) and asynchronous units of work are
called **behaviours**. In Yz, a cown is any value and a behaviour is any boc
invocation — the concepts map directly onto the language's  constructs rather and do not require separate primitives.

The formal model, semantics, and implementation details are described in:

> Cheeseman et al. (2023). *When Concurrency Matters: Behaviour-Oriented Concurrency*.
> Proc. ACM Program. Lang. 7, OOPSLA2, Article 276.
> <https://marioskogias.github.io/docs/boc.pdf>