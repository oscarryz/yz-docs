#feature
# Boc Interface `#()`

The `#(...)` notation defines the **boc interface** — the type and public contract of a boc. It describes what a boc accepts as inputs and what it produces as outputs.

## Signature entries

The entries inside `#(...)` follow the same rules as variable declarations:

- Named with a type: `#(a Int)`
- With a default value: `#(a Int = 1)`
- Short decl: `#(a : 1)`
- Another boc interface: `#(a #())`
- Generic: `#(a T)`
- Generic with constraint: `#(o O Ord)`

## Inputs and outputs

**Labeled entries** (`name Type`) are input fields — the caller provides them. **Unlabeled entries** are outputs — the last N expressions in the body:

| Signature | Meaning |
|---|---|
| `#()` | no inputs, returns nothing |
| `#(String)` | no inputs, returns a String |
| `#(n Int, String)` | input field `n`, returns a String |
| `#(Int, String)` | no inputs, returns an Int and a String |
| `#(f #(A,B), items [A], [B])` | two input fields, returns an array |

```yz
#()              // returns nothing
#(Int)           // returns an Int
#(v Int)         // input field v, returns nothing
#(v Int, Int)    // input field v, returns an Int
#(T, U)          // returns a T and a U
```

## Boc declaration and expanded form

Two ways to attach a body to a boc interface.

The **boc declaration** form — signature and body together; the body uses the params directly:

```yz
greet #(name String) {
  print(name)
}
```

The **boc expanded form** — `=` separates signature from body; the body re-declares all params:

```yz
greet #(name String) = {
  name String
  print(name)
}
```

## The Interface

A boc interface `#(...)` serves two purposes simultaneously.

**Structural typing** — any boc whose shape matches the signature satisfies it; no `implements` needed:

```yz
Greeter #(greet #())           // any boc with a greet field qualifies
Runner  #(run #(), stop #())   // any boc with both run and stop qualifies
```

**Access control** — only fields declared in the interface are visible to external callers:

```yz
Person #(name String) {
    name String
    password String   // not in interface — hidden from callers
    greet #() { print(name) }
}

alice: Person("Alice", "secret")
alice.name      // "Alice" — accessible
alice.password  // error — not accessible
```

These two concerns are fused by design. Writing an interface simultaneously declares the type constraint and narrows the public surface.

## Inference and encapsulation

When no interface is written, a **short boc declaration** (`name : { }`) infers the interface from the body — all declared variables become public:

```yz
person : {
    name String
    password String
}
person.name      // accessible
person.password  // accessible — everything exposed
```

Choosing a **boc declaration** makes the interface explicit and keeps everything else hidden:

```yz
Person #(name String) {
    password String   // internal — not in interface
}

alice: Person("Alice", "secret")
alice.name      // accessible
alice.password  // not accessible
```

This is Yz's encapsulation model: the form you choose determines visibility. Short boc declaration = everything public. Boc declaration = explicitly controlled interface.
