#open-question

Define the exact way the `Macro` interface will work and how the many implementations would implement it, aside from the definition in [Macros](../docs/Features/Macros.md).

---

## What belongs in an annotation

An annotation (backtick boc body before a declaration) serves multiple purposes:

- **Macro configuration** — data consumed by a code-generating macro (`json:`, `cue:`)
- **Build configuration** — data consumed by build tools (`go_source:`, `dependencies:`, `source_paths:`)
- **Passive metadata** — data read by external tools, linters, doc generators

All three share the same syntax (Yz boc body). The question is how each kind is dispatched and validated.

---

## Unified Schema interface

Everything that interacts with annotations declares a `Schema` associated type. This gives compile-time annotation validation uniformly, whether the consumer is a macro, a build tool, or a linter.

The proposed hierarchy:


```yz
// NOTE: the Foo : Bar { ... } is merely illustrative
// Root: anything that can appear in an annotation
Annotatable : {
    Schema #()
}

// Subtypes

Macro : Annotatable {
    run #(Boc, Boc)   // additionally transforms AST; returns modified boc
}

GoSource : Annotatable {
    Schema #(go_source String)   // links a Yz type to a Go implementation file
}

Project : Annotatable {
    Schema #(dependencies [Dependency], source_paths [Path])
}
```

`Macro` is a subtype of `Annotatable` — it has a Schema AND transforms AST. `GoSource` and `Project` are siblings: they validate annotation content but do not transform AST. They affect the build graph instead.

This reflects the fundamental distinction:
- **Macros** — code transformation; run during compilation; input/output are both `Boc`
- **Tools** (`GoSource`, `Project`, `yz fetch`) — build environment; include files, fetch deps; no AST output

---

## Three-tier annotation content (for scalar types and library bindings)

For built-in scalar types (`Int`, `String`, etc.) and Go library wrappers, the annotation and method body work together in three tiers:

**Tier 1 — compiler built-in** (stays hardcoded)
Literals and essential operators (`+`, `-`, `==`). The compiler handles these by name. No annotation needed. Never changes.

**Tier 2 — Go-backed** (`go_source`)
Methods that require external Go code (`parse`, `format`, library calls). Declared in Yz (body-less), implemented in the linked `.go` file. The Go file handles all adaptation — wrapping Go errors into Yz `Result`, mapping Go types to Yz types, etc.

**Tier 3 — pure Yz**
Higher-level methods written entirely in Yz, built on tier 1 and tier 2. Portable across backends.

```yz
`go_source: "yz/std/int.go"`
Int: {
    + #(other Int)             // tier 1 — compiler built-in
    parse #(s String, Int)     // tier 2 — implemented in int.go
    times #(n Int, Range) {    // tier 3 — pure Yz
        Range(0, n)
    }
}
```

For Go library wrappers, the same mechanism generalises:

```yz
`go_source: "yz/k8s/controller.go"`
Controller: {
    deploy #(name String, replicas Int)
    scale  #(name String, count Int)
}
```

The `.go` file is a normal Go file that imports the external library and adapts it to Yz types and semantics. The compiler includes it in the build; naming convention (`deploy` → `Deploy`) links the declaration to the implementation.

---

## Dispatch — open question

When the compiler sees a key in an annotation body (e.g. `go_source:`, `json:`, `dependencies:`), how does it find the handler?

**Candidate models:**

1. **Explicit list** — `macros: [JSON, Cue]` names what to invoke. `go_source:` and `deps:` are compiler-reserved keys, not dispatched through the macro system. Clean; slightly verbose; FQN verbosity solved later by an import system.

2. **Schema-based lookup** — the compiler finds all `Annotatable` implementations in scope, reads their `Schema` associated types, and maps key names to handlers. `GoSource.Schema` has `go_source String` → any annotation with `go_source:` dispatches to `GoSource`. Uniform; requires the type system and scope resolution to be available before annotation processing.

3. **Hybrid** — compiler-reserved keys (`go_source`, `dependencies`) are special-cased initially and gradually promoted to real `Annotatable` implementations as the language becomes self-hosting.

The hybrid is the pragmatic path: define the interface and target design now; implement the first few cases as compiler-internal special cases; promote them to real Yz definitions when bootstrapping allows.

---

## Bootstrap problem

To process a `GoSource.Schema` annotation, the compiler needs `String` to work. `String` is defined via `go_source`. Circular.

Same issue as Rust's `#[lang]` items: the compiler must hardcode the scalar types at a bedrock level, and the Yz definitions of those types are aspirational until the language is self-hosting. The interface hierarchy documents the target design; the compiler implementation reaches it incrementally.

This is not a reason to avoid the design — it is a reason to implement it in stages.

---

## Open questions

- What is the correct name for the root interface? (`Annotatable`, `CompilerTool`, `Processor`, `Handler`?)
- Schema-based dispatch: how does key name → handler resolution work in the presence of imports and FQNs?
- Can an annotation have multiple handlers for different keys? (e.g. `go_source:` AND `macros: [JSON]` in the same annotation)
- How does `GoSource` interact with the macro system — can a macro also have `go_source`?
- Stage-0 plan: which types must remain hardcoded forever, and which can be promoted to real Yz definitions?
