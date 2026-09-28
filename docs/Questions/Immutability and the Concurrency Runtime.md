---
type: open-question
generated: { by: "oscarryz", at: 2026-09-28T00:00:00-06:00 }
sources:
  - id: rej-immutability
    resource: rejected/Immutability.md
    title: "Immutability (rejected)"
  - id: rej-consider-immutability
    resource: rejected/Consider Immutability.md
    title: "Consider Immutability (rejected)"
  - id: rej-immutable-outside
    resource: rejected/Immutable outside of block.md
    title: "Immutable outside of block (rejected)"
  - id: boc-interface
    resource: ../Features/Boc Interface.md
    title: "Boc Interface `#()`"
  - id: concurrency
    resource: ../Features/Concurrency.md
    title: Yz Concurrency
  - id: rej-stateless-bocs
    resource: rejected/Stateless bocs and pure functions.md
    title: "Stateless bocs and pure functions (rejected)"
  - id: yzc-0107
    resource: ../Implementation/tasks-done.md
    title: "YZC-0107 — Sibling call to a non-leaf held-cown callee deadlocks when its result is forced in the same method"
---
#open-question

### Problem

Immutability was previously proposed and dropped three times:

- Immutability[^rej-immutability] — `#rejected`/`#dropped`, tagged "No immutability." Rebinding (`s = 'new name'`) was judged sufficient; instance fields are treated as copy-on-write by convention rather than enforced.
- Consider Immutability[^rej-consider-immutability] — explored value semantics (Elixir-style: value immutable, reference reassignable). Rejected with the reasoning: "All things will be mutable; with Concurrency[^concurrency] each variable would have a single writer, thus, this is not strictly speaking needed."
- Immutable outside of block[^rej-immutable-outside] — proposed that a boc's fields be mutable inside, immutable from outside (mutation only through methods). Rejected with the same reasoning: "If we achieve '1 writer' with 'everything is an actor' this is not needed."

All three rejections rest on one argument: the concurrency runtime (single-writer-per-resource, described in Concurrency[^concurrency]) already gives race-free mutable state, so a separate immutability mechanism is redundant for *correctness*.

That argument holds for correctness. It does not cover two things Yz currently has no answer for.

### Gap 1 — boc interfaces give access control, not mutability control

Boc Interface `#()`[^boc-interface] narrows the *visible* surface of a boc — fields left out of `#(...)` aren't reachable from outside:

```yz
Person #(name String) {
    name String
    password String   // not in interface — hidden from callers
    greet #() { print(name) }
}
```

But every field that *is* in the interface is both readable and writable externally. There is no way to say "`name` is visible but not assignable." Narrowing the interface shrinks which fields you can reach, not what you can do once you've reached one. A caller holding `alice` can do `alice.name = "eve"` for any field the interface exposes — the "privacy" is a visibility filter, not a mutability guarantee. This is the gap the earlier rejections never addressed, because none of them were framed as an access-control question — they were framed as "do we need value semantics," and the answer each time was "the runtime already prevents races," which is a different question from "can a caller overwrite a field I only meant to expose for reading."

### Gap 2 — the concurrency runtime pays a cost that immutable data doesn't need to pay

Per Concurrency[^concurrency]: every value is a protected resource; "only one running boc can acquire a resource at a time" — including two bocs that only ever *read* it. There's no reader/writer distinction. Two concurrent readers of the same value still serialize through the resource's queue, because the runtime has no way to know neither of them will write.

An immutable value can never be the target of a write, by construction. If the compiler can prove a value is immutable, it can never be a data-race source, and none of the machinery in Concurrency[^concurrency] — atomic multi-resource acquisition, happens-before queueing, the resource being "acquired" at all — is needed to use it. Concurrent readers of immutable data don't need to queue behind each other; there's no writer to order against. This is a real, currently-unavailable optimization: it's not just cheaper uncontended acquisition (which the runtime already claims — "a single atomic operation with no waiting"), it's *zero* acquisition, and it also removes the serialization that today's single-writer model imposes even on read-only access patterns (e.g. many bocs sharing a read-only config/lookup value all queue today; they wouldn't need to if immutability were expressible and provable).

This connects to the same territory as Stateless bocs and pure functions[^rej-stateless-bocs], which found a way to opt *bocs* out of actor serialization (the Uppercase convention) for parallel *computation*. There's currently no analogous way to opt *data* out of resource-queue serialization for parallel *reads*.

Adjacent precedent, not a match: **YZC-0107**[^yzc-0107] fixed a reentrant-cown deadlock — `check() { countdown(3) == 0 }` forced a self-recursive sibling call's result while `check`'s own cown was still held. The fix, `hoistHeldCownCalls` (`compiler/internal/ir/lower.go:3670`), splices the nested call out into its own `name : call(...)` binding so it takes the already-correct spawn-then-wait path instead of being forced in place. That's real prior art for rewriting the AST to change *when* a value participates in cown scheduling — but it's a scheduling-order fix, not an immutability exemption; nothing in that ticket lets a value skip cown acquisition altogether. No such exemption mechanism exists in the codebase today — the closest thing to it is the question this doc is asking.

### What's being asked here

Given the surface-vs-mutability gap and the concurrency-runtime cost argument above — both absent from the three earlier rejections — is it worth reopening immutability, scoped narrowly to:

1. A way to mark an interface field (or a whole `#(...)`) read-only from outside the boc, independent of whether the underlying field is mutated internally — closing Gap 1 without touching value/reference semantics.
2. Whether "provably immutable after construction" is a case the compiler could detect (e.g. a struct boc with no internal mutation of a given field after its constructor) and use to elide that value's participation in the resource-acquisition protocol — addressing Gap 2.

This is deliberately not proposing value semantics again (Elixir-style copy-on-write, or the earlier "everything is a value" framing) — those were the parts of the prior proposals that were rightly rejected as unnecessary complexity for a single-writer runtime. The question is narrower: read-only *access* as a declared property of an interface, and its potential to let the compiler skip synchronization it currently can't prove is unnecessary.

[^rej-immutability]: Immutability (rejected)
[^rej-consider-immutability]: Consider Immutability (rejected)
[^rej-immutable-outside]: Immutable outside of block (rejected)
[^boc-interface]: Boc Interface `#()`
[^concurrency]: Yz Concurrency
[^rej-stateless-bocs]: Stateless bocs and pure functions (rejected)
[^yzc-0107]: YZC-0107 — Sibling call to a non-leaf held-cown callee deadlocks when its result is forced in the same method
