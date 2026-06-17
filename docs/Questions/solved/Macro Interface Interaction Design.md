#solved 

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

4. **Uppercase type name dispatch** *(current leaning — see [Macros](docs/Features/Macros.md))* — uppercase names in the annotation body are the dispatch key. A bare uppercase name (`Derive`) triggers with no config; a named boc (`JSON: { ignore: false }`) provides config validated against the handler's `Schema`. Lowercase keys are passive metadata — readable by handlers but not themselves dispatch triggers. No explicit list needed; no key-name-to-handler mapping required. The handler's own type name resolves at parse time via normal name resolution, so FQN, imports, and typos are handled uniformly by the existing type system.

The hybrid is the pragmatic path for options 1–3: define the interface and target design now; implement the first few cases as compiler-internal special cases; promote them to real Yz definitions when bootstrapping allows. Option 4 sidesteps the key-name mapping problem entirely — dispatch is type resolution, which the compiler already does.

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


# Discussion

## Dispatch

Option 4 (uppercase type name dispatch) is the current direction. It resolves the open question about key-name → handler resolution by eliminating it: the name in the annotation IS the type name, resolved by the same mechanism as any other type reference. See [Macros](docs/Features/Macros.md) for the updated examples.

---

## Resolution

Three distinct annotation use cases are in scope for the current stage of Yz:

**Runtime metadata** — the annotation boc is validated against its inferred type at compile time and exposed as a typed value at runtime. Access is via `value.annotation.fieldName`. No dynamic map needed; the compiler infers the shape from the annotation boc body and the consumer reads it through normal field access. Libraries that need generic introspection can use boc reflection (`boc.fields[]`).

**Macro dispatch** — an uppercase type name in the annotation body is the dispatch key, resolved by normal type resolution. A bare name (`Debug`) invokes the macro with no config; a named boc (`Serialize: { pretty: true }`) provides config validated against the macro's `Schema`. The macro implements `run #(Boc, Boc)` and transforms the annotated declaration's AST. This covers code generation, derive-style patterns, and any other compile-time AST manipulation.

**Native extension** (`GoSource`) — a hardcoded first-pass in the compiler, distinct from macro expansion. The compiler processes all `Native` annotations before type resolution begins, loading the referenced Go source files so that the types and functions they define are available during normal compilation. This is the mechanism for binding to Go's standard library and third-party Go packages without rewriting them in Yz.

These three cover the meaningful use cases at this stage. The dispatch question is resolved: annotation keys are type names, not strings mapped to handlers. The `GoSource`/`Native` split from the macro system is intentional — it must run earlier than any macro can.

**External tooling (dependencies, project management)** — loading external sources such as dependencies and project configuration is not handled by macros. These are declared as passive annotation metadata (lowercase keys), and external tools such as the build tool read them directly. No macro dispatch, no AST transformation. This keeps compilation predictable: the compiler never fetches, resolves, or executes external processes as a side effect of annotation processing. The build tool owns that responsibility, running before the compiler is invoked.

---

## Future work

The following categories were considered and deferred:

**Conditional compilation** — annotations that cause the compiler to exclude a declaration on certain platforms or under certain build conditions. Mechanically this is a macro that returns an empty boc, but it requires macros to run before type checking (since excluded code may not type-check on the target platform). Not needed until Yz targets multiple backends or platforms.

**Compiler directives** — hints to the backend about optimization or linking behavior (`no_inline`, target selection, ABI). Can be passed as compiler flags in the interim. When the language is further along, these could be a reserved macro subtype with hardcoded backend handling.

**ABI / memory layout control** — `repr(C)`, alignment, packing. Out of scope; all FFI goes through Go, which handles layout on the Yz side's behalf.

**The hybrid case** — annotations that are simultaneously macro dispatch and runtime metadata (e.g., `@Cacheable("myCache")` in Java/Spring, where the string is both a codegen trigger and a runtime value embedded in the generated wrapper). This is expressible in the current model but has no motivating example in Yz yet.

**Root interface naming** — whether `Annotatable` is the right name for the common supertype of `Macro` and `GoSource`. Cosmetic until the interface hierarchy is implemented.
