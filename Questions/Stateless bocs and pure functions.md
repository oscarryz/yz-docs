#open-question

# Stateless bocs and pure functions

## The problem

All bocs in Yz are actors — they have persistent fields and calls are serialized through their queue. This works well for objects and singletons. But it creates a fundamental problem for utility functions:

```yz
add: { a Int; b Int; a + b }
```

`add` is a singleton. `add(3, 4)` sets `add.a = 3`, `add.b = 4`, then runs the body. Concurrent callers queue. Two goroutines cannot call `add` in parallel — they serialize.

This is correct behavior for an actor, but wrong for a mathematical function that has no meaningful shared state.

## The tension

The power of Yz is that bocs are BOTH callable AND field-accessible — the same construct serves as function, object, and actor. Making bocs "stateless" would remove the field-accessible part, which removes the duality that makes Yz unique.

The question: **can a boc opt out of the actor model when it has no meaningful persistent state?**

## What the language already provides

The syntax already encodes the distinction — it's just not obvious:

**Stateful actor** (body form):
```yz
add: { a Int; b Int; a + b }   // singleton, fields persist, calls serialize
```
After `add(3, 4)`, `add.a == 3`. Concurrent calls queue.

**Stateless function** (BocWithSig shorthand form):
```yz
add #(a Int, b Int, Int) { a + b }   // each call is a fresh goroutine, no persistent state
```
`add.a` does not exist. Concurrent calls run fully in parallel. Each invocation is a separate goroutine with local copies of `a` and `b`.

The BocWithSig form is Yz's existing answer to "I want a pure function." The params are local to each call. There is no singleton state. Calls don't serialize.

## The remaining question: can you ever have BOTH?

The tension point: what if you want a boc that:
- Has accessible fields (`person.name`)
- AND can be called concurrently without serializing

This is not possible in the current model — persistent fields imply actor, which implies serialized calls.

**Option A: explicit `stateless` keyword**

```yz
stateless add: { a Int; b Int; a + b }
```

A stateless boc creates a fresh instance for every call (like Uppercase types), but does not persist state. `add.a` is undefined (or returns the last call's value, which is a data race). Concurrent calls run in parallel.

Problem: "stateless" but with accessible fields is incoherent. Which call's `a` does `add.a` return?

**Option B: the field-access distinction determines the semantics**

If you write `add.a`, you've declared that `add` is an object (stateful). If you only write `add(3, 4)` and never access fields, the compiler could detect this and treat it as stateless — no queue, parallel calls.

Problem: this is implicit and surprising. Adding a `add.a` read anywhere in the codebase changes the entire concurrency model of `add`.

**Option C: types (Uppercase) are always fresh instances**

```yz
Add: { a Int; b Int; a + b }
result: Add(3, 4)   // creates fresh instance, runs, result is 7
```

Each `Add(3, 4)` call creates a new instance, runs the body, and the instance is discarded. No shared singleton state. Concurrent calls run in parallel (each has its own instance).

The tradeoff: naming convention (`Add` vs `add`) carries the entire semantic weight of "parallel-safe function" vs "serialized actor." Renaming changes behavior.

**Option D: accept the limitation, document it, and use BocWithSig for functions**

The rule becomes:
- Use `name #(params) { body }` for functions (parallel-safe, no persistent state)
- Use `name: { field_decls; body }` for actors/objects (persistent state, serialized)
- Use `Name: { ... }` for instantiable types (fresh instance per call)

"Utility functions" in libraries should always use the BocWithSig form. The body form is for objects.

This is the least novel approach but the clearest semantics.

## Implications for `person.name`

If `person` is defined with the body form (`Person: { name String; age Int }`), then `person.name` works — it's a stateful instance. No conflict.

If `greet` is defined with the BocWithSig form (`greet #(name String) { print(name) }`), then `greet.name` doesn't exist — there is no persistent state. This is intentional: you traded field access for parallel safety.

You cannot have both simultaneously with the same boc.

## Status

No design decision made. Option D (document the distinction clearly, lean on BocWithSig for functions) is the current de-facto behavior. Options B and C are worth exploring. Option A is probably not viable.

## Related

- [How to use utility functions](solved/How%20to%20use%20utility%20functions.md)
- [utility functions](solved/utility%20functions.md)
- [Concurrency](../Features/Concurrency.md)
- [Single Writer](../Features/Single%20Writer.md)
- [Bocs](../Features/Bocs.md)
