#feature
# Type Alias

`NewType : ExistingType` declares a type alias — `NewType` is another name for `ExistingType`. They are structurally identical; the compiler treats them as the same type.

```yz
Bar : Foo     // Bar and Foo are the same type
```

This is the same `:` short-declaration syntax used for values (`x : 42`). Applied at the type level, the compiler copies `Foo`'s full signature into `Bar`.

Expanded form:

```yz
// If Foo is: { name String }
Bar #( name String ) = Foo
```

The RHS is the implementation source; the LHS signature is inferred from it. In generated Go this becomes `type Bar = Foo`.

---

## Generic instantiation

A type alias where the RHS is a parameterized type creates a concrete specialization:

```yz
StringList : List(String)   // StringList.T = String
IntPair    : Pair(Int, Int)  // IntPair.K = Int, IntPair.V = Int
```

The expanded form substitutes the type arguments into the signature:

```yz
StringList #( add #(String), remove #(String), size #(Int) ) = List(String)
```

---

## Associated types via alias

Inside a concrete boc, type aliases bind the abstract type fields declared in an interface:

```yz
Graph : {
    Node #()
    Edge #()
    ...
}

SocialGraph : {
    Node : User         // type alias — Node is User in this boc
    Edge : Relationship
    ...
}
```

This is the mechanism behind associated types. See [Path Dependent Types](Path%20Dependent%20Types.md) for the full picture.

---

## Overriding behavior on construction

A type alias shares the implementation of its source. To use the same structure but replace a method, construct with an override:

```yz
A : {
    say_hi : { "Hi" }
}

B : A    // B is an alias for A

C : {
    say_hi : A.say_hi
    say_hi = { "Bye" }   // overrides say_hi in C's body
}

print(A().say_hi())   // Hi
print(B().say_hi())   // Hi
print(C().say_hi())   // Bye
```

---

## See also

- [Path Dependent Types](Path%20Dependent%20Types.md) — the unified model for aliases, generics, and associated types
- [Generics — Type Parameters](Generics%20-%20Type%20Parameters.md) — type parameters in boc declarations
