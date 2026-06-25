#solved

Define how a boc literal gets its type and how that type is used for compatibility checking across all use sites (struct, closure, interface).

---

## The problem

The compiler currently classifies boc literals during analysis — deciding early whether a literal is "a closure" or "a struct" — and branches on that classification throughout sema, the lowerer, and codegen. This creates fragile code paths and makes reasoning about boc typing non-uniform.

---

## Resolution

### BocLiteralType is a flat list of fields

A boc literal is a collection of fields. Each field has a name, a type, and optionally a value or body. No classification into Params / Methods / Fields / Returns during analysis — those are use-site concerns.

In Yz's own model, a field inside a boc is itself a boc — there is no separate "Field" concept:

```yz
// Conceptual model (Yz)
BocLiteralType {
    fields [Boc]   // each field is a boc: name, type, optional value/body
}
```

In the compiler implementation (Go), this maps to:

```go
// Compiler representation (Go)
type BocLiteralType struct {
    Fields []FieldNode   // FieldNode: Name string, Type Type, Value optional
}
```

The derived interface is mechanically computed from those fields:

```yz
{ x: 5, name: "foo", greet #() { ... }, y Int }
// derives: #(x Int, name String, greet #(), y Int)
```

### Field roles are use-site concerns

Whether a field acts as an input parameter, a stored value, a method, or a return value is not a property of the literal. It is determined at the use site by matching the derived interface against the expected type.

- Used as a struct → fields are stored state
- Used as a closure → params are inputs; last expression type is the return
- Used as an interface → fields are checked for presence and type

### Structural compatibility — one rule

> i1 satisfies i2 if i1 has **every field** in i2 with matching name and type. Extra fields in i1 are allowed. Missing fields are not.

Default values in the expected interface are **irrelevant** to structural compatibility. They only matter at direct call sites (see below).

### Default values — two distinct rules

**Rule 1 — Direct call site:** the compiler knows the full signature of the callee and fills in defaults for omitted arguments.

```yz
foo #(a Int, b Int = 1) { ... }
foo(1)   // b filled in by compiler — allowed
```

**Rule 2 — Structural compatibility:** a boc passed as a value must have every field of the expected interface, regardless of defaults.

```yz
bar #(baz #(a Int, b Int = 1)) { ... }
bar({ a: 1 })   // { a: 1 } derives #(a Int) — b is absent — incompatible
```

`{ a: 1 }` does not have `b`. The receiver may call `baz.b`; if `b` is absent that is a type error. Default values in the expected type do not fill the gap — the boc simply must have the field.

Default values in the **derived** interface (e.g. `{ a: 1 }` → field `a` has value `1`) carry no special weight for compatibility. Only the name and type are checked.

---

## Consequences

- `BocLiteralType` in `sema/types.go` becomes a flat field list — no Params / Methods / Fields / Returns subdivision
- All classification branches in sema and the lowerer (`hasInnerBocsOrMethods`, `bocLitHasParams`, `anonBocCache`, `anonDecls`) are deleted
- The lowerer dispatches on the **use-site expected type** rather than sema classification flags
- Structural compatibility is one function used uniformly across all use sites
- Implementation tracked in YZC-0080
