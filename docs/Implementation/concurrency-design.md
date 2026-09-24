---
type: implementation
generated: { by: "oscarryz", at: 2026-05-09T04:23:18+02:00 }
---
#implementation 
# Yz Concurrency: BOC Model and Implementation Design

## 1. Current State

The compiler today emits code that uses Go goroutines for concurrency:

- Every boc method call → `std.Go(func() T { ... })` → goroutine launched immediately
- Every boc invocation in `main` → `_bg0.Go(func() any { ... })` inside a `BocGroup`
- The parent boc waits for all children via `BocGroup.Wait()` (structured concurrency)
- Return values are `*Thunk[T]`; `.Force()` blocks until the goroutine finishes

**The correctness gap.** The current model does not protect shared state. Two concurrent calls to the same singleton's method run as independent goroutines and both access `self.count` without any coordination:

```yz
counter: {
    count: 0
    increment: { count = count + 1 }
}
main: {
    counter.increment()   // goroutine A: read count, add 1, write back
    counter.increment()   // goroutine B: same — concurrent access — DATA RACE
}
```

The Go race detector will flag this. The two increments may both read the same initial value, losing one write. The Yz language specification guarantees this cannot happen: `counter` is a protected resource, and `increment` calls sharing it are serialized in spawn order.

---

## 2. BOC Concepts and Their Yz Mapping

The Behaviour-Oriented Concurrency model (Cheeseman et al., OOPSLA 2023) centers on two primitives:

| BOC term | Definition | Yz mapping |
|---|---|---|
| **Cown** | A concurrent owner: a protected resource | Every singleton boc instance (e.g. `Counter`) |
| **Behaviour** | A unit of work that runs when it holds all its cowns | A boc method invocation or top-level boc call |
| **Atomic acquisition** | A behaviour acquires all required cowns at once, or waits | Method dispatch blocks on the receiver cown |
| **Spawn order** | Two behaviours sharing a cown run in the order they were spawned | Methods on the same receiver always execute in call order |
| **Parallel execution** | Behaviours with no shared cowns run in parallel | Calls to different singletons still run concurrently |

### Primitive values are not cowns

In the paper, every object can be a cown. In Yz's implementation, **all boc instances** are cowns — both singleton bocs (lowercase, one global instance per FQN) and struct boc instances (uppercase, e.g. `Account`). Each instance carries an embedded `std.Cown`. Primitive values (`std.Int`, `std.String`, etc.) are plain Go values that are not cowns themselves; they are protected by the cown of the boc instance that owns them. You only read or write a primitive field from inside a `Schedule`/`ScheduleMulti` body that holds the enclosing boc's cown, or after all concurrent work on that boc is complete (after `BocGroup.Wait()`). Therefore no extra protection is needed at the primitive level.

### Thunks and cowns coexist

A thunk is the *return value* of a boc invocation — a placeholder that resolves when the invocation completes. Cowns govern *when the invocation is allowed to run*. These are orthogonal:

- Cown acquired → behaviour runs → produces a value → thunk resolves → cown released
- Caller gets the thunk immediately; cown serialization is invisible to the caller

The thunk model (lazy, transparent, `.Force()` at IO boundaries) does not change.

---

## 3. Structured Concurrency — What Stays the Same

The "spirit" of structured concurrency is unchanged:

- A parent boc does not complete until all invocations it spawned (directly or indirectly) have finished, even if their return values have already been delivered as resolved thunks.
- This is still implemented via `BocGroup` + `sync.WaitGroup`: each spawned child increments the group; `Wait()` at the end of the parent's body blocks.
- The only change: inside the spawned closure, the actual method call goes through the cown scheduler before running, not straight into execution. From `BocGroup`'s perspective this is transparent.

```
spawn(counter.increment())
    │
    ├─ BocGroup.Go(...)          ← spawns goroutine (structured concurrency tracking)
    │       │
    │       └─ Counter.Increment()  ← goroutine enters cown queue for Counter
    │               │
    │               └─ (waits for Counter cown if another behaviour holds it)
    │                       │
    │                       └─ runs, produces Thunk, cown released
    │
BocGroup.Wait()              ← parent blocks here until ALL children done
```

---

## 4. What Needs to Change

### Runtime: `Cown` type and `Schedule` function

Each singleton boc struct must carry a cown. The simplest correct first implementation is a mutex per instance:

```go
// runtime/rt/cown.go (Phase A)

type Cown struct {
    mu sync.Mutex
}

func (c *Cown) acquire() { c.mu.Lock() }
func (c *Cown) release() { c.mu.Unlock() }

// Schedule runs fn exclusively while holding cown c, then releases it.
// Returns a Thunk that resolves once fn completes.
func Schedule[T any](c *Cown, fn func() T) *Thunk[T] {
    return Go(func() T {
        c.acquire()
        defer c.release()
        return fn()
    })
}
```

### Generated singleton structs

```go
// Before (current):
type _counterBoc struct {
    count std.Int
}

// After (Phase A):
type _counterBoc struct {
    std.Cown             // embeds the mutex-based cown
    count std.Int
}
```

### Generated method dispatch

```go
// Before (current):
func (self *_counterBoc) Increment() *std.Thunk[std.Unit] {
    return std.Go(func() std.Unit {
        self.count = self.count.Plus(std.NewInt(1))
        return std.TheUnit
    })
}

// After (Phase A):
func (self *_counterBoc) Increment() *std.Thunk[std.Unit] {
    return std.Schedule(&self.Cown, func() std.Unit {
        self.count = self.count.Plus(std.NewInt(1))
        return std.TheUnit
    })
}
```

The call site in `main` does not change — the cown locking is fully inside the method.

### Multi-cown behaviours (Phase B)

When a boc call passes multiple singleton boc instances as arguments, the call must atomically acquire all of their cowns:

```yz
transfer #(src Account, dst Account, amount Int) {
    src.balance = src.balance - amount
    dst.balance = dst.balance + amount
}
```

`transfer` is a singleton boc (lowercase = one shared instance). Its fields `src`, `dst`, `amount` are set by the caller. The call `transfer(src, dst, 100)` must acquire both `src` and `dst` atomically before running — otherwise another behaviour could modify `dst` between the two field assignments.

Phase B introduces `ScheduleMulti` which atomically acquires all cowns in canonical order (sorted by pointer address to prevent deadlock, as in the BOC paper):

```go
func ScheduleMulti[T any](cowns []*Cown, fn func() T) *Thunk[T] {
    // Sort cowns by address — canonical order prevents deadlock.
    sort.Slice(cowns, func(i, j int) bool {
        return uintptr(unsafe.Pointer(cowns[i])) < uintptr(unsafe.Pointer(cowns[j]))
    })
    return Go(func() T {
        for _, c := range cowns { c.acquire() }
        defer func() {
            for _, c := range cowns { c.release() }
        }()
        return fn()
    })
}
```

The lowerer identifies which arguments in a boc call are singleton boc instances (cowns) and emits `ScheduleMulti` with those arguments instead of a single-cown `Schedule`.

---

## 5. Phased Implementation Plan

### Phase A — Mutex cowns (correct, simple)

Goal: eliminate data races. Non-overlapping singletons still run in parallel because they have separate mutexes.

1. Add `Cown` struct (embeds `sync.Mutex`) and `Schedule[T]` to `runtime/rt/cown.go`
2. Codegen: emit `std.Cown` as embedded field in every singleton boc struct
3. Codegen: change method thunk from `std.Go(...)` to `std.Schedule(&self.Cown, ...)`
4. Regenerate golden tests; run `go test -race ./...`

Limitation: does not yet implement true spawn-order determinism for same-cown behaviours beyond what mutex FIFO provides (Go's `sync.Mutex` is not strictly FIFO, though in practice FIFO under contention). The correctness guarantee (no data races) holds.

### Phase B — True BOC queue-based scheduler

Goal: strict spawn-order guarantee; lock-free implementation matching the paper's algorithm.

Each `Cown` holds an atomic pointer to the tail of a linked list of pending behaviours. A behaviour decrements an atomic counter as each of its cowns grants it a token; it runs when the counter reaches zero. This is the algorithm described in section 3 of the paper.

1. Replace mutex-based `Cown` with the atomic-queue implementation
2. `ScheduleMulti` registers the behaviour on all cowns simultaneously
3. Ordering is guaranteed by the queue structure, not by lock ordering
4. No sorting needed (the queue handles ordering per-cown)

#### Re-entrancy: sub-behaviours needing the parent's cown

The BOC paper does not define re-entrancy (a behaviour calling into a cown it already holds). In Yz this arises naturally: a method body spawns a child boc that needs the same receiver cown.

```yz
foo: {
    baz: { /* uses foo's fields */ }
    run: {
        baz()    // spawns a child that needs foo's cown
    }
}
```

With `BocGroup.Wait()` inside the Schedule closure this deadlocks: `run` holds foo's cown, `baz` queues on foo's cown, `run` waits for `baz` — cycle.

The fix: **`BocGroup.Wait()` must happen after the cown is released**, not inside the Schedule closure.

```go
// Wrong — deadlock if baz also needs self's cown:
func (self *_fooBoc) Run() *Thunk[std.Unit] {
    return std.Schedule(&self.Cown, func() std.Unit {
        _bg := &std.BocGroup{}
        _bg.Go(func() any { return self.Baz().Force() })
        _bg.Wait()   // ← holding cown here, Baz() can't acquire it
        return std.TheUnit
    })
}

// Correct — cown released before waiting:
func (self *_fooBoc) Run() *Thunk[std.Unit] {
    _bg := &std.BocGroup{}
    bodyDone := std.Schedule(&self.Cown, func() std.Unit {
        _bg.Go(func() any { return self.Baz().Force() })
        return std.TheUnit
    }) // cown released here
    return std.NewThunk(func() std.Unit {
        bodyDone.Force()
        _bg.Wait()   // ← cown free, Baz() can now acquire it
        return std.TheUnit
    })
}
```

This preserves structured concurrency (parent waits for all children) while avoiding the deadlock.

#### Re-entrancy, take two: sync-rewriting instead of deferring (2026-09-24, YZC-0008)

The deferred-`Wait()` fix above only covers a re-entrant call in **statement position** — the
lowerer only scanned top-level statement expressions for calls needing this treatment. A held-cown
call reached through a *nested* expression position (a call argument, e.g. `print(baz())`, or a
string-interpolation part, e.g. `"${baz}"`) was invisible to that scan and still deadlocked.

The fix isn't another deferred-`Wait()` variant: since the caller already holds the cown a nested
call needs, that call can just be executed as a **plain synchronous call** — no `Schedule`, no
goroutine, no `BocGroup` — with identical observable behavior. This is cheaper than deferral and
covers any AST shape uniformly (nested boc-param call, struct-instance method call, sibling-method
call all reduce to the same check: is the target's cown already held?).

**The one thing this is not safe for: a callee that itself has pending sub-work on the same
cown.** If the callee is self-recursive (or otherwise spawns a further call that needs the same
cown), converting the *outer* call to sync does not make the *inner* one safe — the inner call
still needs the cown released to run, and it now can't be, because the sync rewrite kept everyone
inside the same still-running, cown-holding call. Concretely: `f() { ...; f(n-1) }` calling itself
reentrantly must stay on the async/deferred-`Wait()` path above; only a call to a **leaf** callee
(no boc calls of its own) may be rewritten to sync. This isn't a conservative choice — rewriting
the recursive case to sync reintroduces the original deadlock: the inner call needs the cown that
the outer (now-synchronous) call cannot yet release. See `tasks-done.md`'s YZC-0008 entry for the
code-level details.

### Phase C — Closures capturing cowns

Nested bocs close over their enclosing boc's fields, which include potential cown references:

```yz
foo: {
    bar Account    // bar is an Account singleton — a cown
    baz: {
        bar.withdraw(10)   // baz closes over bar; this needs bar's cown
    }
}
```

In Phase A/B, `baz` is a method on `_fooBoc` and accesses `self.bar`. If `bar` is a struct field (a value), it's protected by `foo`'s cown — no issue. If `bar` is a reference to a standalone singleton boc (a separate cown), then `baz`'s call to `bar.withdraw(10)` needs `bar`'s cown in addition to `foo`'s cown.

Phase C analysis:
- Track which fields of a boc are themselves singleton boc instances (cowns)
- When a nested boc method accesses such a field via a method call, the callee's cown handles protection (method dispatch acquires the callee's cown)
- The open question: does the enclosing boc method (`baz`) need to declare `bar` as a needed cown upfront (for atomic acquisition), or is sequential per-cown acquisition (Phase A) sufficient?

Phase C is deferred until Phase A and Phase B are stable. It requires the boc uniformity work (all bocs lower to structs) to be complete first, since closures-as-boc-methods only works correctly once all bocs have structs.

---

## 6. What Does NOT Change

| Aspect | Status |
|---|---|
| Yz syntax | Unchanged — no new keywords, no annotations |
| `*Thunk[T]` return type | Unchanged — every boc call still returns a lazy thunk |
| `Force()` semantics | Unchanged — blocks until the thunk resolves |
| Structured concurrency (parent waits for children) | Unchanged — `BocGroup` + `WaitGroup` still used |
| `thunkVars` tracking in the lowerer | Unchanged — auto-inserts `.Force()` on use |
| Primitive types (`std.Int`, `std.String`, etc.) | Unchanged — protected by their enclosing singleton cown |
| `BocGroup.Go(...)` in `main` bodies | Unchanged — still spawns goroutines for structured concurrency tracking |

---

## 7. Invariants That Hold After Phase A

- **Data-race freedom**: two calls to any method of the same singleton run sequentially (mutex serializes them); different singletons run in parallel.
- **No deadlock**: single-cown acquisition cannot deadlock. Multi-cown (Phase B): canonical ordering prevents circular waiting.
- **Structured lifetimes**: a parent boc's `BocGroup.Wait()` ensures all children complete before the parent exits, regardless of cown scheduling.
- **Thunk transparency**: callers never observe the cown mechanism; they only see the resolved `*Thunk[T]` value.
