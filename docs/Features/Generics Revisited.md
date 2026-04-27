# Yz Generics

## Core Philosophy

Yz generics are built on three principles:

- **Inference first** — constraints are never written by the developer, they are derived from usage and surfaced by the compiler and tooling
- **Minimalist** — no new constructs are introduced; generics fall out of existing language concepts
- **Statically typed** — all types are fully resolved at compile time despite being inferred

---

## Type Parameters

A type parameter is any **single uppercase letter** identifier: `T`, `U`, `V`, etc.

This is a language-level rule — a single uppercase letter is always and only a type parameter. There are no exceptions and no ambiguity.

```yz
// T and U are type parameters — unambiguous by identifier shape alone
transform : {
    thing T
    mapper U
    mapper.apply(thing)
}
```

This convention fits naturally with the existing identifier system:

| Shape | Meaning |
|---|---|
| lowercase | singleton |
| Uppercase multi-letter | instantiable boc |
| Single uppercase letter | type parameter |

---

## Constraint Inference

Constraints on type parameters are **never declared**. The compiler observes how a type parameter is used inside a boc body and derives the constraint automatically.

```yz
greet : {
    thing T
    thing.talk()    // compiler infers: T must have talk #()
}
```

The inferred constraint is the complete and precise description of what `T` must be. If a caller passes a boc that does not satisfy it, the error is reported at the call site:

```yz
greet(Person("Ann"))  // ok — Person has talk #()
greet(Animal())       // ok — Animal has talk #()
greet({})             // error: no talk method
```

---

## Operators As Methods

All operators in Yz are regular methods on boc. The `Int` type for example is:

```yz
Int : {
    + #(other Int, Int)
    > #(other Int, Bool)
    // ...
}
```

This means operator usage is inferred as a method constraint, exactly like any other method call. No special syntax or construct is needed:

```yz
add : {
    a T
    b T
    a + b    // inferred: T must have + #(T, T)
}
```

There are no special cases for operators anywhere in the generics system.

---

## Type Parameters Inside boc — Associated Types

A type parameter declared inside a boc acts as a **type slot** — a placeholder whose concrete value is fixed by each implementation. This is possible because in Yz a type is a boc like anything else, so holding a type in a slot is no different from holding any other value.

```yz
Converter : {
    T
    Output #()             // type slot — fixed per implementation
    convert #(T, Output)
}
```

Each implementation fixes `Output` concretely:

```yz
IntToString : {
    Output : String
    convert #(Int, Output) = { /* impl */ }
}

IntToFloat : {
    Output : Float
    convert #(Int, Output) = { /* impl */ }
}
```

The compiler tracks the relationship between a boc and its type slots. When used in a generic context, the output type is fully resolved through inference:

```yz
process : {
    thing T
    converter U
    result = converter.convert(thing)  // result is inferred as U.Output
    result.doSomething()               // adds to constraint on U.Output
}
```

No additional type parameter is needed to name the output type — the compiler derives and tracks it from `U`'s definition.

---

## Constraint Propagation

When a generic boc calls another generic boc, constraints propagate upward through the call chain automatically.

```yz
greet : {
    thing T
    thing.talk()       // constraint: T has talk #()
}

greetAll : {
    things []T
    things.forEach({
        thing T
        greet(thing)   // greet's constraint propagates up to greetAll's T
    })
}
```

`greetAll` ends up with the same constraint as `greet` — `T must have talk #()` — even though `talk` is never mentioned in `greetAll`'s body directly.

Chains of arbitrary depth are followed. The final constraint surfaced to the developer is always **flattened** — the full set of requirements on `T` regardless of how deep they originated.

---

## Compile-Time boc

A boc may contain a compile-time block that executes during compilation. This block has access to the inferred type information of the boc it annotates and can emit errors, synthesized types, or documentation.

This is the mechanism for metaprogramming in Yz. It is consistent with the rest of the language — a compile-time block is a boc like anything else, just executed at a different phase.

```yz
add : {
    compile: {
        // runs at compile time
        // has access to inferred constraints on T
        // can validate, emit docs, or synthesize types
    }
    a T
    b T
    a + b
}
```

The design of compile-time boc is ongoing.

---

## Module Boundaries and Tooling

The compiler always has full knowledge of inferred constraints. At module boundaries — where a developer or user of a boc cannot see its body — the compiler and tooling surface the flattened inferred constraint automatically.

A generic boc like:

```yz
add : {
    a T
    b T
    a + b
}
```

Is surfaced by tooling as approximately:

```
add #(a T, b T, T)
  T requires: + #(T, T)
```

The language itself does not need syntax to express this. It is the tooling's responsibility, not the language's.

---

## What Is Not In Yz Generics

The following features present in other languages are deliberately absent, with reasoning:

| Feature | Why absent |
|---|---|
| Explicit constraint declaration | Inference covers this completely |
| Variance annotations (`in`/`out`) | No inheritance means no variance problem |
| Type elements / union constraints | No primitive/object split — everything is a boc |
| Higher-kinded types | Contradicts minimalism; inference covers common cases |
| Explicit where clauses | Tooling's job, not the language's |
| Default method implementations | Not yet needed; under consideration |

---

## Comparison With Other Languages

| Feature | Yz | Go | Rust | Kotlin | Java | Haskell |
|---|---|---|---|---|---|---|
| Constraint declaration | Inferred | Explicit | Explicit | Explicit | Explicit | Inferred |
| Operator constraints | Methods on boc | Type elements | Operator traits | Operator methods | Not supported | Type classes |
| Associated types | Type slots (inferred) | Not supported | Explicit | Not supported | Not supported | Explicit |
| Variance | N/A — no inheritance | N/A | Automatic | Declaration-site | Use-site (PECS) | Automatic |
| Primitive/object split | None | Partial | None | JVM split | Full split | None |
| Higher-kinded types | Not supported | Not supported | Not supported | Not supported | Not supported | Supported |
| Monomorphization | TBD | GCShape stenciling | Full | JVM: erasure / Native: full | Erasure | Dictionary |
| Metaprogramming | Compile-time boc | None | Macros | Annotations | Annotations | Template Haskell |
| Type surfacing | Tooling (flattened) | Explicit in code | Explicit in code | Explicit in code | Explicit in code | Inferred signature |