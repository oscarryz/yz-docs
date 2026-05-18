#solved

## Decision

- `name Type` (uninitialized typed declaration) is a **required field** — the compiler enforces it must be assigned before first use (definite assignment analysis).
- `Bar()` with unassigned required fields is a compile error **unless** every unassigned field is provably initialized on all control-flow paths before it is read.
- Crossing a boc boundary: the compiler requires the value to be fully initialized at the call site.
- `Option(T)` is for values that are **semantically optional** (absent is a valid state). It is not a workaround for uninitialized fields.
- Default values (`f String = ""`) are for fields with a **clear meaningful zero**. They also allow `Bar()` with no arguments.

See: `spec/03-expressions-and-statements.md` §3.2, `docs/Features/Variables.md`. Tracked as YZC-0034.

---

## Original Discussion

Should creating a new instance initialize all the uninitialized fields? Or should the compiler validate before calling? 


e.g. 

```js
// Uninitialized fields
Person : { 
  name String
  age Int
}
p : Person() // Should this be a compilation error? 
// vs
p : Person(name: "Alice", age:42) 
// or 
p : Person("Alice", 42)

// Or should be only an error if I try to use it? (which could be harder to detect)

p.name // error 
p.age  // error 
...
foo: {
  q Person
  print("${q.name}")
}
foo(p) // will use p -> error, name is not intialized 
...
```

If that's not a compilation error: How what is the program behavior when the variable is accessed. 
If that's a compilation error: How can we create a boc whose values are know later on? 

Possible answer: 
- ~~Force to define a default value.~~
- ~~Define an `Option` value e.g. `Person: { name Option(String), age Option(Int) }`~~
- (Settled) Definite assignment: compiler tracks assignment before use within a scope; construction requires all required fields to be provided or provably assigned before first read.
