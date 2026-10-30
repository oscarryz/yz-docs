---
type: impl
generated: { by: "oscarryz", at: 2026-06-05T23:56:26+02:00 }
---
#impl
Legacy detail entries for closed tickets. Open tickets live in
[tickets/](tickets/index.md).

---

## Infrastructure

- [x] **[YZC-0093] Uppercase root file (`Foo.yz`) always-wrap: example + spec §9 clarification**

  YZC-0092 implemented always-wrap for lowercase root files only. Uppercase root
  files (e.g. `Foo.yz`) follow the same invariant but have no conformance test.

  Two sub-cases:
  - `Foo.yz` with free-floating fields (`name String; age Int`) → wraps to
    `Foo: { name String; age Int }` — the wrapper IS the struct type.
  - `Foo.yz` with an inner `Foo: {}` → wraps to `Foo: { Foo: {} }` — inner
    becomes an associated type (struct-outer, YZC-0082); `fileWrapperHasInnerBoc`
    should trigger unwrap so the inner `Foo` becomes top-level.

  Work:
  - Add a small root-level uppercase example (e.g. `examples/root_type/`) with
    a `Foo.yz` struct file and a `main.yz` that constructs it
  - Add spec §9 Invariant 1 clarifying note covering both sub-cases
  - Add or promote conformance test

- [x] **[YZC-0092] Always-wrap root files; `main()` as explicit entry invocation**

  Remove the `hasTopLevelBocNamed` guard in `build.go` so all root files are
  always wrapped in a boc named after the file — consistent with spec §9
  Invariant 1 ("file content = boc body named after the file").

  Consequences:
  - `main.yz` with `main: {}` wraps to `main: { main: {} }`. The outer `main`
    executes its body; to run the inner boc the file must call `main()` explicitly.
  - `main.yz` with free-floating statements (no inner `main: {}`) works as-is —
    the outer `main` body just executes them directly.
  - `Foo.yz` with `Foo: {}` wraps to `Foo: { Foo: {} }` — `Foo` inner becomes
    an associated type (YZC-0082). `Foo.yz` with `name String; age Int` wraps
    to `Foo: { name String; age Int }` — struct type, constructor works.

  Work:
  - Remove `hasTopLevelBocNamed` guard and helper from `build.go`
  - Update all existing examples: drop explicit same-name wrapper, or add
    `main()` call at end of `main.yz`
  - Update conformance tests accordingly
  - Add clarifying example to spec §9 Invariant 1 covering `Foo.yz` and
    `main.yz` with inner `main: {}`

- [x] **[YZC-0059] Design: macro interface interaction** — [resolved](../Questions/solved/Macro%20Interface%20Interaction%20Design.md)

  Settle the taxonomy and interaction model for macros before implementing
  YZC-0028 or any macro-dependent feature (YZC-0041, YZC-0058, YZC-0060).
  Depends on: ~~YZC-0025~~.

  **Questions to resolve:**

  **1. Macro taxonomy — two distinct concepts, not one**

  - *Code-generating macros*: receive an annotated Boc, return a transformed
    Boc merged back into the AST. Closest to Rust proc-macros. Declare
    `schema #()` so the compiler can validate the annotation body at
    compile time. Run during compilation.
  - *Build-hook programs*: `yz fetch`, `Native`, Go generator wrappers. Act
    on the environment (network, filesystem), not on code. No schema, no AST
    output. Run before or alongside compilation. `yz fetch` is the canonical
    example — a standalone program that reads annotations and populates a cache.

  The key differentiator: **schema**. Macros declare `schema #()`; the
  compiler validates the annotation body against it at compile time. Programs
  parse annotations ad hoc with no compile-time guarantees.

  **2. Macro identification — three candidate models**

  - *Explicit list only*: `!:[JSON]` / `macros: [JSON]` is the authority.
    Unambiguous; slightly verbose when the field name already implies the macro.
  - *Attribute name implies macro*: `json: {...}` in an annotation body
    implicitly triggers the `JSON` macro. Concise; risks silent failure on typo
    or when macro is not in scope.
  - *Hybrid (recommended starting point)*: explicit list (`!:[...]`) is the
    authority and always required. Macros find their config by convention —
    a macro named `JSON` looks for a `json:` field in the annotation body.
    No implicit triggering; no redundant naming. Explicit list = what runs;
    field name = where config lives.

  **3. Annotation format for tooling (gates YZC-0041)**

  `yz fetch` needs to locate dependency declarations in annotations. The
  agreed format (to be finalised here):

  ```
  `!:[Deps]
  dependencies: [
      my_lib: { version: "1.0.0", url: "https://..." }
  ]
  `
  ```

  - `!:[Deps]` (or `macros: [Deps]`) = explicit trigger
  - `dependencies:` = config field, routed to `Deps` by convention
  - Deps can appear inline in source or in a `name.info` companion file
  - Multiple declarations across files are merged; one `yz.lock` per project root
  - Lock file pins SHA (git commit SHA preferred — immutable by definition)
  - `Deps` macro = thin schema-only macro for compile-time validation;
    actual fetching is done by `yz fetch` (a standalone program, not a macro)

  **4. Programs vs macros — clear boundary**

  Programs (`yz fetch`) are not macros. They scan annotations for markers,
  act on the environment, and pass results to the compiler (e.g. extra source
  roots via `yzc build proj/ ~/.yz/cache/lib/`). They do not satisfy the
  `Macro` interface and do not transform AST. The `Deps` macro and `yz fetch`
  are complementary: the macro validates, the tool acts.

---

## Tooling

- ~~[ ] **[YZC-0041] `Deps` macro — compile-time dependency validation**~~

  ~~Cancelled.~~ Dependencies are passive annotation metadata, not macros.
  Superseded by YZC-0097.

- [x] **[YZC-0097] Annotation metadata contract for project and dependency configuration** *(replaces YZC-0041)*

  Define the annotation format for project-level and dependency metadata so
  that external tools (`yz fetch`, build tool) have a stable, documented
  contract to read from. This is passive metadata — lowercase keys, no macro
  dispatch, no AST transformation. The compiler validates the annotation shape
  at compile time (typed metadata) but never fetches or resolves anything.
  Compilation remains predictable: no external processes are triggered by the
  compiler as a side effect of annotation processing.

  - Define `dependencies:` annotation format (fields: `version`, `url`, `sha`)
  - Define `project:` annotation format (name, source paths, description)
  - Specify where these annotations may appear (project root `.yz` file or
    `project.info` companion)
  - Document the contract between annotation metadata and external tools

  ```yz
  `dependencies: [
      my_lib: { version: "1.0.0", url: "https://..." }
  ]
  project: { name: "my_project", source_paths: ["src/"] }`
  ```

  Depends on: ~~YZC-0059~~.

---

## Major Features

### ~~YZC-0076 — Existential associated types: opaque-token / path-identity tracking~~ — CLOSED

**Closed.** The original motivation was to support an array of heterogeneous macros, each carrying a different `Schema` shape, validated at compile time:

```yz
macros: [Debug, Serialize, Ord]   // each has a different Schema
```

Iterating over that list and validating each schema would require knowing the specific `Schema` type per element — which requires either path-dependent types (`m.Schema` where `m` is a variable) or first-class existential types.

The dispatch mechanism changed (YZC-0059): macros are now triggered by uppercase type name resolution in the annotation body. The compiler always dispatches to a specific concrete type (`Debug`, `Serialize`, etc.) — it never holds a generic `Macro` variable and dispatches through it. Therefore `m.Schema` on a runtime variable never arises, and the cross-root rejection problem does not occur.

**What Yz has:**
- Associated types (YZC-0066) — compile-time, concrete shape known when the specific type is in hand
- Structural compatibility — checks annotation body against the concrete `MacroImpl.Schema`

**What Yz does not have:**
- Runtime path-dependent types — `m.Schema` where `m` is a variable of interface type `Macro`
- First-class existential types — holding a heterogeneous collection and unpacking each element's specific associated type shape

These gaps would matter if macros were ever runtime-dispatched objects or if a user wanted a heterogeneous list of macros with per-element schema introspection. Neither is needed under the current dispatch model.
