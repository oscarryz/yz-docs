#open-question

Tracked as **YZC-0060**.

### Problem

Struct bocs and singleton bocs currently have no way for a method body to refer to "the current instance" as a whole value — only to individual fields by name. Two earlier attempts at solving this are on record and both fail the same way:

- [`self` and-or initialize literals](solved/`self`%20and-or%20initialize%20literals.md) — tried a stored `self` field, wired up by a constructor (`p.self = p`) or by a `new` boc assigning to a shared `self` var. The doc's own conclusion: *"this doesn't work because `self` will be overridden and will change"* — every new instance clobbers the same slot.
- [How and when include self](How%20and%20when%20include%20self.md) — proposed generating `self` via a macro (`` `!:[Native, Derive, Self]` ``). Superseded by this proposal — see "Why not a macro" below.

Both attempts share the same root cause as JavaScript's `this` problem: they try to make `self` a *value assigned at some point in time* (construction, or a macro-injected initializer) rather than *an expression resolved lexically wherever it's written*.

### Why not a macro

This was explored in detail before writing this proposal. Two independent reasons a macro can't do this:

1. **Category error.** Macros run at compile time over a structural snapshot of a type *declaration* (`Boc{name, fields}`), before sema, before any instance exists. There is no live object for a macro to reference — a macro can only emit source text once, for the type, not per-instance. Whatever a macro generates has to be the same for every instance, so it can't "point back" at a specific one.
2. **Even a hypothetical macro with write-access to an enclosing/"parent" scope is a worse problem than the one it solves.** Today macro expansion is a pure one-way function: `subject → generated source appended to that same subject`. That's what keeps it tractable — reading a boc's own declaration (plus its annotation) tells you everything a macro will do to it. Letting a macro reach into a parent or sibling scope to wire up a self-reference reopens ordering hazards (which macro runs first), cycles (parent mutates child while child's macro expects the mutated parent), cache-key blowup (the run cache is keyed on subject-payload hash; parent-dependence would invalidate it constantly or leave it stale), and hygiene collisions (two macros picking the same generated slot name). Real macro systems that avoid this class of bug (e.g. Rust proc-macros) deliberately scope expansion to "the tokens you're attached to, in, tokens back out" for exactly this reason.

### Proposal: `self` as a lexically-resolved compiler built-in

Add `self` as a reserved identifier, valid inside any method body of a struct boc or singleton boc, that denotes the current receiver instance. Resolution is **lexical, at compile time** — fixed by which boc-method body textually encloses the reference — never by how or when the method is called.

This is not new machinery: the lowerer already tracks exactly this internally. `internal/ir/lower.go` sets `RecvName: "self"` and a `receiverState{name, fields, structType}` per boc-method body via `setReceiver`/`restoreReceiver`, and field references inside a method are already rewritten to `self.field` in the generated Go. This proposal exposes that existing internal receiver binding to Yz source as a first-class expression, instead of only using it to prefix field names.

```yz
Person: {
    name String
    print_self #() Unit {
        print(self.name)   // ordinary field access, works today
        print(self)        // NEW — the whole instance
    }
    identity #() Person {
        self                // NEW — return the current instance
    }
}
```

#### Semantics

- `self` is valid only inside a method body belonging to a struct boc or singleton boc. Using it anywhere else (top-level statements, a boc literal that never becomes a method, a free function) is a **compile-time error**, not a nil/undefined value at runtime.
- **Nested boc-literal closures lexically capture the enclosing `self`**, the same way they already capture other outer variables (see `list.filter`/`list.each` HOF closures). A closure boc-literal is not itself a fresh receiver scope, so:
  ```yz
  Person: {
      name String
      fields [Field]
      describe_each #() Unit {
          fields.each({ f Field
              print(self.name + "." + f.name)   // self = the enclosing Person instance
          })
      }
  }
  ```
  This is the exact case that breaks in JavaScript (a nested/callback function silently gets a *different* `this`). It doesn't break here because `self` was never rebound per call — the closure has no receiver of its own to rebind it to.
- No detached-method-reference form exists in Yz (no `.bind`/`.call`/`.apply`-equivalent, no passing a bare method off its instance) — so the mechanism JS uses to *produce* the wrong `this` (multiple call-time invocation styles for the same function body) has no counterpart to build. There is exactly one way `self` gets its value: lexical position.
- Interaction with `go_source`-backed methods (YZC-0058) — `self` inside a body-less, Go-delegated method is out of scope for this ticket; tracked separately as **YZC-0099**. Working hypothesis: `go_source`-backed (tier-2) methods have no Yz body at all, and the bound Go function already receives the instance explicitly as its first parameter by the existing `//yz:bind` convention — so there's likely nothing left to design here, only to confirm. `self` used inside a *pure-Yz* method on a `go_source`-annotated type (e.g. a `times` helper built on top of a native `+`) is expected to just work under this ticket's ordinary rules.

#### Open implementation questions

- **Reserved word vs. context-sensitive keyword.** `self` is not currently reserved anywhere in the lexer/parser/sema (confirmed by grep — it's purely the hardcoded Go receiver name in the lowerer). Making it reserved means deciding what happens to any existing field/param/var literally named `self` — hard reserved word (compile error to declare one) vs. context-sensitive (only special inside expression position, shadowable as a field name like any other). Leaning toward hard-reserved for simplicity, since `self.self` would be confusing either way.
- Whether `self` is settable (`self = ...`) or read-only. Read-only is simpler and avoids reintroducing a "self gets reassigned" failure mode akin to the rejected `solved/self and-or initialize literals` approach.
- Exact grammar: is `self` a primary expression (like `Ident`) or its own AST node? Given sema already special-cases receiver-scoped field lookups, likely cheapest as a special `Ident` case resolved during sema against `recvStructType`, mirroring what the lowerer already does — see `setReceiver`/`recvFields` in `internal/ir/lower.go`.
