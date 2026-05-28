#impl-plan
# YZC-0066 Implementation Plan — Path-Dependent Types

Design doc: [`docs/Features/Path Dependent Types.md`](../Features/Path%20Dependent%20Types.md)

---

## What the compiler already does

Before any changes, the compiler already:

- Handles bare `T` (GENERIC_IDENT) in a struct body: appends to `st.TypeParams`, registers a `GenericType` symbol in scope
- Emits `Box[T any]` in Go for structs with TypeParams
- Emits generic Go functions `[V any]` for BocDecl with type params
- In golden tests 29–33, `Box: { T; value T }` already produces `type Box[T any] struct { value T }` — T is **not** emitted as a runtime field

The changes are therefore mostly **additive**. The existing TypeParams mechanism continues to drive Go generic output; the new `IsTypeField` flag is layered on top.

---

## Phasing

### Phase A — Explicit T field + constructor syntax + `b.T` (Steps 1–4, 7–8)
No new AST nodes. No parser changes. Makes `Box: { T; value T }` / `Box(Int, 42)` / `Box(value: 42)` / `b.T` work correctly.

### Phase B — `g.Node` in signatures + type alias inside bocs (Steps 5–6, 9)
Requires a new `MemberTypeExpr` AST node and a parser lookahead change. Enables `process #(g Graph, n g.Node)` and `Node : User` inside `SocialGraph`.

### Phase C — Call-site type variable unification (Step 10)
Full inference for `map #(collection List(A), fn #(A, B), List(B))`. Complex; can stay deferred post-Phase B.

---

## Step 1 — Add `MetaType` and `IsTypeField` to the type system

**File:** `compiler/internal/sema/types.go`

**Changes:**

Add `MetaType` — the type of every type (`#()`). Every UpperCase boc satisfies it:

```go
type MetaType struct{}
func (*MetaType) typeName() string { return "#()" }

var TypMeta = &MetaType{}
```

Compatibility rule: any `*BuiltinType`, `*StructType`, or `*GenericType` is compatible with `*MetaType` (every type satisfies the metatype). Add a `case *MetaType: return true` to `IsCompatibleWith` on those types.

Add `IsTypeField bool` to `StructField`:

```go
type StructField struct {
    Name        string
    Type        Type
    HasDefault  bool
    IsTypeField bool  // field holds a type value (bare GENERIC_IDENT or `Node : User`)
}
```

**Risk:** Low — purely additive.

**Test:** Unit test: `TypMeta.typeName() == "#()"` and each built-in type is compatible with `TypMeta`.

---

## Step 2 — Bare GENERIC_IDENT → IsTypeField in struct body

**File:** `compiler/internal/sema/analyzer.go`

**Current behavior** (around line 984): bare `T` in a struct body appends to `st.TypeParams` and registers a `GenericType` in scope. No `StructField` is created for T.

**New behavior:** additionally append a field:

```go
st.Fields = append(st.Fields, StructField{
    Name:        ident.Name,  // "T"
    Type:        TypMeta,
    IsTypeField: true,
})
```

Keep the TypeParams append and scope registration unchanged — Go generic output still relies on them.

**Risk:** Medium. Every site that iterates `st.Fields` for constructor purposes must be updated to skip `IsTypeField` entries (Steps 7–8). Missing one site silently generates invalid Go. Audit list:

- `lowerStructArgs` in `lower.go`
- `checkGenericConstraints` in `analyzer.go`
- `initLocalVar` in `definite_assign.go`
- All `for _, f := range st.Fields` loops in lower.go and codegen.go

**Also affects variant types:** `Option: { V; Some(value V); None() }` — V must also get IsTypeField. Verify golden tests for variants (tests 40–50 range) still pass.

**Test:** Sema unit test: after analyzing `Box: { T; value T }`, the StructType has two fields — one with `IsTypeField=true, Name="T"` and one regular `Name="value"`.

---

## Step 3 — Constructor syntax for explicit type fields

**File:** `compiler/internal/sema/analyzer.go`, `compiler/internal/ir/lower.go`

**Scenario A — positional: `Box(Int, 42)`**

In `analyzeCall`, when the callee is a StructType with IsTypeField entries, split positional args into type-args (filling IsTypeField fields) and value-args (filling regular fields). Build a type-substitution map from type-args.

**Scenario B — named: `Box(value: 42)` (T inferred)**

When no type arg is provided but named value args are, infer T from the value arg type. Concretely: match `value: 42` → `value` field has type `T` → T must be `Int`.

**Practical first-pass approach:** Emit the Go call without explicit type args (`NewBox(value)`) and let Go's own type inference fill in T. No explicit unifier needed in Phase A. Full unification deferred to Phase C.

**Error case:** `Box(42)` with positional and no named args — 42 would map to the T (IsTypeField) position. Since 42's type is `*BuiltinType{Int}` and IsTypeField expects `*MetaType`, report:

```
error: YZC-00XX: Box(42) — first argument maps to type field T; 42 is a value not a type.
       Use Box(Int, 42) or Box(value: 42)
```

**Risk:** Medium-high for the error detection; low for the code generation (Go infers).

**Test:** Golden test `66_box_explicit.yz`: `Box: { T; value T }` / `b: Box(value: 42)` / `print(b.value)` → compiles, prints 42. Error test: `Box(42)` on explicit-T Box → correct error.

---

## Step 4 — `b.T` as valid member access

**File:** `compiler/internal/sema/analyzer.go` — `fieldType` function (around line 1389)

**Change:** `fieldType` already looks up field by name. Since T is now in `st.Fields` (Step 2), `b.T` returns a field type of `TypMeta`. Mark the expression with type `MetaType`.

In `analyzeMember`: skip the definite-assignment check for IsTypeField fields (they are always initialized at construction time via the type arg).

In `lower.go` — `lowerMember`: when the field is an IsTypeField, emit... nothing useful at runtime (T is compile-time only). For now, emit a compile-time constant string of the type name — or better, skip it and mark as a type-only expression. Practically: since Go generic structs don't store T, `b.T` in Yz source has no direct Go equivalent in Phase A. Emit a TODO comment and an empty string placeholder, or restrict `b.T` to type-annotation positions only (defer value-position use to Phase B/C).

**Risk:** Medium. The interesting case — `b.T` used as a type annotation elsewhere (`n b.T`) — requires path-dependent resolution (Step 9). For Phase A, support `b.T` only in value-expression position returning `MetaType`; full use in type position is Phase B.

**Test:** Sema test: `b.T` on an explicit-T Box has type `*MetaType`.

---

## Step 5 — Implicit T → auto-collect into TypeParams

**File:** `compiler/internal/sema/analyzer.go`

**Scenario:** `Box: { value T }` — T appears in a field type but is never declared on its own line.

**Current behavior:** `resolveTypeExpr` for GENERIC_IDENT returns `&GenericType{Name: "T"}` immediately. TypeParams is empty, so the emitted Go struct is non-generic (`Box` with a field `value T` — invalid Go since T is not defined).

**Change:** After processing all elements of a struct body in `analyzeStructBoc`, scan all field types for any `*GenericType` whose Name is not yet in `st.TypeParams`. Append it. This auto-populates TypeParams from implicit use, making the Go output `Box[T any]` without an IsTypeField field for T.

**Invariant:** IsTypeField entries are only created by explicit bare-identifier declarations (Step 2). Implicit type variables (discovered via field type scan) go into TypeParams only.

**Risk:** Low.

**Test:** Sema unit test: `Box: { value T }` → TypeParams=["T"], zero IsTypeField fields. Compare with `Box: { T; value T }` → TypeParams=["T"], one IsTypeField field.

---

## Step 6 — Type alias inside a boc (`Node : User`)

**File:** `compiler/internal/sema/analyzer.go`

**Scenario:** `SocialGraph: { Node: User; Edge: Relationship; ... }` — Node and Edge are ShortDecl where the RHS is a TYPE_IDENT resolving to a struct type.

**Detection:** In `analyzeShortDecl`, when:
- The LHS name is uppercase (TYPE_IDENT convention), AND
- The RHS is a bare `*ast.Ident` (not a CallExpr), AND
- The resolved type of the RHS is `*StructType`

→ This is a type alias binding. Create the StructField with `IsTypeField: true`.

**Effect:** `SocialGraph.Node` becomes an IsTypeField holding `User`'s StructType. `g.Node` (when g is a SocialGraph) resolves to `User`.

**Risk:** Medium. The heuristic (bare uppercase ident = type reference) needs care: `x: MyVar` where MyVar happens to be uppercase but is a value, not a type. Guard: check that the symbol's `Node` in the scope is a type symbol (has no constructor call args, its type is `*StructType` with `IsVariant: false` etc.).

**Test:** Sema test: `SocialGraph: { Node: User; ... }` → Node field IsTypeField=true, field type is `*StructType{Name: "User"}`.

---

## Step 7 — IR: skip IsTypeField in constructor loops

**File:** `compiler/internal/ir/lower.go`

All sites that iterate `st.Fields` to build Go constructor parameters must skip IsTypeField entries:

| Function | Location | Change |
|---|---|---|
| `lowerStructArgs` | ~line 462 | Add `if f.IsTypeField { continue }` |
| `checkGenericConstraints` (in analyzer.go) | ~line 1689 | Same guard — already skips BocType; add IsTypeField |
| `initLocalVar` in `definite_assign.go` | | Same guard |
| Any other `for _, f := range st.Fields` that builds arg lists | Audit needed | Same guard |

**Risk:** Medium — requires thorough audit. Approach: after changes, run `go test ./...` and check for Go compile errors in golden test output files (invalid generated code fails the compile step in conformance tests).

**Test:** Run all 65 golden + 18 error tests — all must still pass.

---

## Step 8 — Codegen: skip IsTypeField in Go struct emission

**File:** `compiler/internal/codegen/codegen.go`

In `emitStructDecl` (around line 95), when emitting the Go struct body fields, skip IsTypeField entries:

```go
for _, f := range sd.Fields {
    if f.IsTypeField { continue }  // compile-time only; not a runtime field
    // emit: f.Name + " " + goType(f.Type)
}
```

Same skip in the Go constructor function parameter list (~line 135).

**Practical note:** The existing golden tests already show `type Box[T any] struct { value T }` — T is not emitted as a runtime field. This step formalises that existing behavior behind the IsTypeField flag instead of relying on the "T matches a TypeParam entry" coincidence.

**Risk:** Low.

**Test:** Golden test: `Box: { T; value T }` → Go struct has exactly one field (`value T`), not two.

---

## Step 9 — `g.Node` in type position (Phase B)

**Files:** `compiler/internal/ast/ast.go`, `compiler/internal/parser/parser.go`, `compiler/internal/sema/analyzer.go`, `compiler/internal/ir/lower.go`

### 9a — New AST node

```go
// MemberTypeExpr represents a type expression of the form `g.Node`
// appearing in a type-annotation position.
type MemberTypeExpr struct {
    Pos
    Object Expr   // the instance whose type field is accessed
    Member string // the type field name
}
func (*MemberTypeExpr) typeNode() {}
```

### 9b — Parser change

In `parseTypeExpr` (around line 894): after parsing a `SimpleTypeExpr` (bare IDENT), if the next token is `.` followed by a TYPE_IDENT, consume both and return a `MemberTypeExpr`. This requires a one-token lookahead — use the existing `save`/restore pattern already used in `parseArrayOrDict`.

### 9c — Sema: `resolveTypeExpr` case for `MemberTypeExpr`

```go
case *ast.MemberTypeExpr:
    sym := a.currentScope.Lookup(t.Object.(*ast.Ident).Name)
    if sym == nil { /* error */ }
    st, ok := sym.Type.(*StructType)
    if !ok { /* error */ }
    for _, f := range st.Fields {
        if f.Name == t.Member && f.IsTypeField {
            // return the type stored in this field
            return f.Type  // MetaType if abstract; concrete type if known
        }
    }
    // error: no type field named t.Member on this type
```

For the simple case where `g : SocialGraph()` (concrete local), `sym.Type` is `*StructType{Name:"SocialGraph"}` and `Node`'s field type is `*StructType{Name:"User"}` → return that directly.

For the abstract case where `g Graph` (parameter), `sym.Type` is `*StructType{Name:"Graph"}` and `Node`'s field type is `MetaType` → return `MetaType` (or a new `PathDependentType`). At this stage, emit `any` in Go and defer full checking.

### 9d — IR/lowerer

In `goType`, add a case for `*PathDependentType` (if introduced) → emit `any` for now.

**Risk:** High. The parser lookahead and the `Object` being an IDENT vs something more complex (what if `Object` is itself a member expression?) need careful scoping. For Phase B, restrict to single-level `ident.TypeName` only.

**Test:** Golden test `66_graph_node.yz`:
```yz
Graph: { Node #(); neighbors #(Node, [Node]) }
SocialGraph: { Node: Int; neighbors #(Node, [Node]) = { n Node; [n] } }
process: { g Graph; n g.Node; g.neighbors(n) }
main: { sg: SocialGraph(); process(sg, 42) }
```

---

## Step 10 — Call-site type variable unification (Phase C, deferred)

**Scenario:** `map #(collection List(A), fn #(A, B), List(B))` — A and B are unbound type variables in a BocDecl signature. At the call site, infer A from `collection`'s element type and B from `fn`'s return type.

**Deferral rationale:** Go's own generic type inference handles the generated output correctly. The sema gap (unresolved A/B in return-type expressions) manifests as `GenericType` values flowing into variables, which get `any` in Go. This is incorrect for strict type checking but generates working code.

**When to implement:** After Phase B is stable and if type errors at call sites become a practical problem.

---

## Dependency graph

```
Step 1 (MetaType + IsTypeField)
    │
    ▼
Step 2 (GENERIC_IDENT → IsTypeField field)
    │
    ├──▶ Step 5 (implicit T auto-collect into TypeParams)
    │
    ├──▶ Step 7 (IR: skip IsTypeField in constructor loops)
    │        │
    │        ▼
    │    Step 8 (codegen: skip IsTypeField in Go struct)
    │
    ├──▶ Step 3 (constructor syntax: Box(Int,42) / Box(value:42))
    │
    └──▶ Step 4 (b.T member access)
              │
              ▼
          Step 9 (g.Node in type position)  ← Phase B
              │
              ▼
          Step 6 (Node : User inside boc)   ← Phase B
              │
              ▼
          Step 10 (call-site unification)   ← Phase C, deferred
```

---

## Tests per phase

| Phase | Tests |
|---|---|
| A complete | All 65 golden + 18 error still pass; new golden `66_box_explicit.yz`; error test `Box(42)` on explicit-T Box |
| B complete | New golden `66_graph_node.yz` (Graph/SocialGraph/process); `66_social_graph_alias.yz` (Node : User binding) |
| C complete | Golden for `map #(List(A), #(A,B), List(B))` call site |

---

## Known risks and edge cases

1. **~6–8 sites in lower.go iterate `st.Fields` for constructor args.** Missing any one silently generates invalid Go. Systematic fix: add a helper `dataFields(st *StructType) []StructField` that filters IsTypeField entries, and use it everywhere.

2. **Variant types** (`Option: { V; Some(value V); None() }`) — V must also get IsTypeField under the new model. Verify golden tests 40–50 range still pass.

3. **`TypeSignature` in types.go** iterates `st.TypeParams` separately then fields. With IsTypeField fields that match TypeParams names, the output might duplicate. Audit carefully.

4. **`initLocalVar` in `definite_assign.go`** iterates `st.Fields` to mark optional fields as assigned. IsTypeField fields are always "assigned" (filled at construction time via type arg). Add `if f.IsTypeField { fi.locals[varName][f.Name] = true; continue }`.

5. **Pre-existing issue (golden test 33 — `makePair`):** Generated Go uses out-of-scope type params at the call site. The new implementation should not make this worse. Full fix requires Step 10.

6. **`g.Node` parser lookahead** must not consume a `.` that belongs to a method call chain. Scope: only trigger `MemberTypeExpr` when parsing in a type-expression context (not in expression context).

7. **`Box: { T; value T }` already works in practice** for Go output — do not regress it. The change in Step 2 must remain a no-op for existing Go output while adding the IsTypeField metadata.
