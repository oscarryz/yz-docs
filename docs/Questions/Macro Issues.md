#implementation #open-question 

### Package Isolation & Circular Dependencies

- Add a dedicated callout (and a specific compiler error message) explaining the separate-package requirement for custom macros
- Package Isolation: Validate at build-time that macros live in separate packages to prevent circular deps.

### Subprocess Latency & Serialization Optimization

- Investigate an in-process fast path (shared library) for Compile execution to mitigate cold-build subprocess latency
- Serialization Cost: Benchmark subprocess serialization; add lazy evaluation or in-process fallback option if overhead is significant.
- Tooling: Add lint warnings for large bocs prone to serialization overhead, with best-practice guidance.

### Sandboxing, Determinism, & Network/IO Control

- Add a note that compile-time IO (e.g. http.get) is the implementer's responsibility to make deterministic/cacheable
- Implement a capability-based sandbox to strictly control and track side-effects like file system access.
- Enforce a "No-Network" constraint for Macros to guarantee offline-first, reproducible build hermeticity.

### annotation Validation, Scope, & Fallbacks

- Specify the fallback behavior when a field's annotation is missing a variable expected by a macro (default, compile error, etc.)
- Validation: Enforce compile-time validation of annotations against their Schema — fail early on mismatches.
- annotation Scope Docs: Clarify that only boc/object-level annotations trigger compilation; field annotations are passive metadata.

### Error Propagation & Diagnostic APIs

- Error Propagation: Standardize error messages for failed macros with boc/field context; always halt compilation.
- Provide a diagnostic API within run to emit compiler errors mapped directly to the user’s annotation source lines.

### Constraint Enforcement & Conflicts

- Upgrade constraint declaration from "strong convention" to a compiler warning when a macro impl generates calls on a type parameter without declaring constraints on S
- Constraint Conflicts: Add detection and clear errors for unsatisfiable/conflicting constraints from multiple macros.

### Name Mapping, Derivation, & Lints

- Add a lint/warning for name collisions between single-field Schema implementations sharing the same derived key
- Name Derivation Docs: Document rules for annotation field mapping; require explicit name: when Schema has multiple fields.

### Caching & Build Dependencies

- Cache Invalidation: Document robust strategy covering external resources, dependency versions, and forced rebuilds.
- Automatically register external files read during run as build dependencies for the incremental caching engine.

### Slot Generation (Visibility, Shadowing, & Tooling)

- Require an explicit override modifier when a Compile block intentionally shadows an existing slot to prevent bugs.
- Ensure IDE "Jump to Definition" support can trace generated slots back to the specific Compile logic that created them.
- Standardize a synthetic flag for generated slots so documentation tools can distinguish them from handwritten code.
- Define clear "Visibility" rules so Compile blocks can specify if generated slots are public or private to the parent.

### Compilation Ordering & Auditing Tools

- Provide a compiler "Expand" tool to export the final merged Boc structure for auditing and debugging.
- Add an optional priority field to the Macro interface to help the compiler auto-sort the macros array.