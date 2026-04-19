# Boc Uniformity — Design Gap Analysis

## The Intended Design

In Yz there is **one construct: the boc**. Not "bocs and functions." Not "bocs and typed bocs." One thing.

### What a boc is

A boc is a named sequence of declarations and expressions enclosed in `{ }`. Its variables are its **fields**. Its nested bocs are its **children**. Calling a boc runs its body. Fields persist between calls.

```yz
counter: {
    count: 0
    increment: { count = count + 1 }
    value: { count }
}
counter.increment()   // runs the body of `increment`
counter.increment()
print(counter.value()) // 2
```

`count`, `increment`, and `value` are all fields of `counter`. Calling `counter.increment()` runs the body `{ count = count + 1 }`. Calling `counter()` would run counter's own body top-to-bottom, reinitializing its fields.

### Calling a boc with arguments

Positional args set fields in declaration order before running the body:

```yz
greet: {
    name String
    message String
    print("`name`: `message`")
}

greet("Alice", "hello")   // sets name="Alice", message="hello", then runs body
greet.name = "Bob"        // set field directly
greet()                   // runs body again with current fields
```

Fields persist after the call. `greet.name` is `"Bob"` after the second call.

### BocWithSig is syntactic sugar

The expanded form of a typed boc declaration:

```yz
foo #(a Int, Int) = { a Int; a }
```

means: declare `foo` with signature `#(a Int, Int)`, assign a boc whose body redeclares `a` as a field and returns it. The `a` in the signature and the `a Int` in the body are the same field.

The shorthand form (no `=`) avoids redeclaring the fields in the body:

```yz
foo #(a Int, Int) { a }   // a is injected into the body scope automatically
```

Both are equivalent to:

```yz
foo: { a Int; a }         // body form — exactly the same thing
```

`a` is a field of `foo`. It persists between calls. `foo.a` is accessible from outside. Calling `foo(3)` sets `foo.a = 3` and runs the body.

The `#(params)` notation exists for two purposes:
1. To annotate a type signature without a body (`Point #(x Int, y Int)` declares a struct type)
2. As a convenience to avoid redeclaring params in the body

It does **not** create a different kind of boc. There is no such thing as a "stateless boc" or "function boc." There is only the boc.

### Boc literals create instances

A boc literal `{ ... }` without a name creates a **new anonymous boc instance** each time it appears in code. It is not a singleton.

```yz
[{}, {}, {}]           // three distinct instances, structurally identical
list.filter({ item Int; item > 10 })  // creates one new instance, passes it to filter
```

This is identical to:

```yz
pred: { item Int; item > 10 }   // a named singleton
list.filter(pred)               // same instance every call to filter
```

vs.

```yz
list.filter({ item Int; item > 10 })  // a new instance on each filter call
```

Calls from `filter` to the boc argument are serialized through that instance's queue — but since the instance is scoped to one `filter` call and `filter` is itself async, this is not a bottleneck.

### lowercase = singleton, Uppercase = type

- Lowercase name: one shared instance for all callers at that scope. `counter.increment()` always hits the same `counter`.
- Uppercase name: each call creates a fresh independent instance. `Person("Alice")` and `Person("Bob")` are separate objects.
- This rule applies at all nesting levels without exception.

### Source roots

A project has one or more **source root directories**. The source root is a mount point — not itself a boc. Everything inside it is.

```
project/
  src/          ← source root 1 (app code)
    foo/
      bar.yz
  lib/          ← source root 2 (third-party / stdlib)
    baz.yz
```

FQN trees:

```
src root:  foo → bar → (declarations in bar.yz)
lib root:  baz → (declarations in baz.yz)
```

`src/foo/bar.yz` contributes `foo.bar`, `foo.bar.something`, etc. — not `src.foo.bar`. The source root is stripped from the FQN. Multiple roots produce independent FQN forests with no cross-root collision.

### Bocs form a forest

```
project/
  xyz/
    foo.yz:
      bar: {}
      baz: {
        qux: {}
      }
```

FQN tree (single source root at project root):

```
xyz → foo → bar
          → baz → qux
```

`bar` is at "file scope" and `qux` is "nested" — but that is only a position in the tree. The same rules apply:

- Lowercase → singleton at that scope
- Uppercase → type at that scope
- `#(params)` → same boc, fields declared upfront as sugar
- All calls → async, return a thunk
- Fields persist between calls
- Access via FQN from outside, local name from inside

---

## Where the Misunderstanding Crept In

Understanding this section is important to avoid repeating the same confusion.

### The confusion: BocWithSig looked like a function

The `#(params) { body }` syntax looks like a function signature and body in most languages:

```yz
add #(a Int, b Int, Int) { a + b }
```

This looks exactly like:
```go
func add(a int, b int) int { return a + b }
```

The visual similarity pulled the implementation toward treating BocWithSig as a Go function. Once that decision was made, the params (`a`, `b`) became local variables (Go function arguments) rather than fields of a struct. This is where field persistence was lost: after `add(3, 4)`, `add.a` should be `3`, but the Go function lowering makes `a` a stack-local that disappears on return.

### The confusion compounded: "stateless" was applied to BocWithSig

Once BocWithSig was being lowered to Go functions, the word "stateless" appeared in documentation and design discussions — "BocWithSig bocs have no persistent fields." But this is a description of the implementation, not the design. The design has no stateless bocs. All bocs have fields. The `#(params)` form is just a way to declare those fields upfront.

This led to treating BocWithSig as a **different concept** with different semantics, rather than as syntactic sugar over the same concept. Design questions like "when should you use BocWithSig vs body form?" arose from this confusion — but the answer is simply: "BocWithSig is body form with declared-upfront field names."

### The confusion persisted into HOF discussion

When discussing higher-order bocs, "boc vs function" framing led to questions like "does `#(String, Int)` accept only functions or only bocs?" The real question is simpler: `#(String, Int)` describes the shape of a boc (takes a String field, produces an Int). Any boc with that shape satisfies it, whether declared with BocWithSig syntax or body syntax.

### What should have been understood from the start

```
foo: { a Int; a }          ← boc. a is a field.
foo #(a Int, Int) { a }    ← same boc. sugar: a declared in sig, not body.
foo #(a Int, Int) = { a Int; a }  ← same boc. expanded sugar form.
```

The three are identical in semantics. The compiler should lower all three the same way.

---

## Current Implementation Gaps

The compiler has three separate lowering paths — all of which deviate from the true model in different ways.

### Path 1: File-scope body-form boc (`counter: { ... }`)

**Closest to correct.** Generates a Go struct + package-level var + goroutine methods.

```go
type _counterBoc struct { count std.Int }
func (self *_counterBoc) Increment() *std.Thunk[std.Unit] { return std.Go(...) }
var Counter = &_counterBoc{count: std.NewInt(0)}
```

**Gap:** Calling the boc itself — `counter()` — is not supported. Per the design, calling `counter()` runs counter's body, reinitializing its fields. Currently `Counter` is just a struct instance; there is no generated method for "run the top-level body." The `increment` and `value` variables-that-are-bocs are wired as methods, but the outer body is not.

### Path 2: BocWithSig (`countdown #(n Int) { ... }`)

**Semantically wrong.** Generates a Go function with a local parameter:

```go
func countdown(n std.Int) *std.Thunk[std.Unit] { return std.Go(...) }
```

This works as a computational approximation — the countdown runs correctly — but it loses the field semantics:

- `n` is a Go stack variable, not a struct field. It disappears on return.
- `countdown.n` is inaccessible after a call. Per the design it should be `3` after `countdown(3)`.
- The boc has no struct, so no inner bocs or additional fields could be added to `countdown` later without breaking the model.

This path exists because the implementation confused BocWithSig with Go functions. It is a semantic shortcut that happens to produce correct computational output for pure recursive/stateless uses, but it diverges from the design.

### Path 3: Local boc inside a boc body (`f: { n Int; ... }` inside `main:`)

**Broken.** Generates a plain Go function literal in an `any` variable:

```go
var f any = func(n std.Int) std.Unit { ... }
std.Go(func() std.Unit { return f(std.NewInt(3)) })  // no BocGroup, program exits
```

- `f` has no struct, no fields, no methods
- Inner bocs inside `f` are dropped
- Recursive calls fire goroutines with no coordination
- `f(3)` in `main:` has no `BocGroup` — the program exits immediately

This path was left unimplemented because the BocWithSig shortcut (Path 2) was handling the "function-like" use case, and local body-form bocs without a use case weren't prioritized.

### Gap Summary

| Aspect | File-scope body | BocWithSig | Local body |
|---|---|---|---|
| Struct emitted | Yes | **No — Go func** | No — func literal |
| Fields persist between calls | Yes | **No — stack locals** | No |
| `.field` access from outside | Yes | **No** | No |
| Inner bocs as fields | Yes | N/A | Dropped |
| Methods on the boc | Yes | N/A | Dropped |
| Calling the boc itself (`f()`) | **Not generated** | Yes (as Go func call) | Sort-of (wrong type) |
| Goroutine-wrapped calls | Yes | Yes | No |
| BocGroup coordination | Yes (in main:) | Yes (in main:) | No |
| Recursive self-calls | Yes (thunk) | Yes (thunk, fixed) | Fire-and-forget |
| Literal creates instance | N/A | N/A | No — reuses same var |
| FQN in sema | Yes | Yes | Partial |

### HOF callbacks — intentional special case (pending design decision)

`list.filter({ item Int; item > 10 })` — the boc literal argument is currently lowered as a synchronous Go `func(std.Int) std.Bool`, not as a goroutine-returning boc. This is a deliberate compiler shortcut: `filter` needs to call the predicate and get back a Bool synchronously to do its job.

Per the design, this literal IS a boc instance. Calls to it from `filter` would go through its queue — sequentially for that instance, but the `filter` call itself is async from the caller's perspective. This is semantically correct but would be slower than a direct function call for large collections.

This is tracked as a pending design decision before touching goldens 27, 34, and 05 (while).

---

## Desired State

Every boc declaration at any nesting level should produce:

1. A **struct type** (lifted to package level for Go compatibility) named by FQN
2. A **"call" method** that runs the body
3. Additional **named methods** for each inner boc that is itself a named boc
4. An **instance** created at the point of declaration in the enclosing scope
5. All calls to the boc go through goroutine + thunk

```yz
main: {
    boc: {
        n: 0
        bar: { n = n + 1 }
        baz #() { print(n) }
    }
    boc.bar()
    boc.baz()
}
```

Desired Go output:

```go
type _main_bocBoc struct {
    n std.Int
}

func (self *_main_bocBoc) Call() *std.Thunk[std.Unit] {
    return std.Go(func() std.Unit {
        self.n = std.NewInt(0)   // body re-runs reinitializes
        return std.TheUnit
    })
}

func (self *_main_bocBoc) Bar() *std.Thunk[std.Unit] {
    return std.Go(func() std.Unit {
        self.n = self.n.Plus(std.NewInt(1))
        return std.TheUnit
    })
}

func (self *_main_bocBoc) Baz() *std.Thunk[std.Unit] {
    return std.Go(func() std.Unit {
        std.Print(self.n)
        return std.TheUnit
    })
}

func main() {
    _boc := &_main_bocBoc{}
    _bg0 := &std.BocGroup{}
    _bg0.Go(func() any { return _boc.Bar().Force() })
    _bg0.Go(func() any { return _boc.Baz().Force() })
    _bg0.Wait()
}
```

For BocWithSig (`countdown #(n Int) { ... }`), the same model applies:

```go
type _countdownBoc struct {
    n std.Int
}

func (self *_countdownBoc) Call() *std.Thunk[std.Unit] {
    return std.Go(func() std.Unit {
        if self.n.Eqeq(std.NewInt(0)).GoBool() {
            std.Print(std.NewString("done"))
        } else {
            std.Print(self.n)
            self.n = self.n.Minus(std.NewInt(1))
            Countdown.Call().Force()   // recursive: same singleton
        }
        return std.TheUnit
    })
}

var Countdown = &_countdownBoc{}
```

Note: for the recursive case, the singleton model means all recursive calls share the same `Countdown` instance. The thunk model ensures no deadlock — `Call()` fires a goroutine and returns immediately, so the queue drains between recursive steps.

---

## Implementation Approach

### Pass 1 — Sema: record struct shape for all boc declarations

- `analyzeBocDecl` should call `analyzeStructBoc` for ALL lowercase boc bodies (not just Uppercase), producing a `StructType` with the correct fields at all nesting levels
- `analyzeBocWithSig` should no longer be a fundamentally separate path — it resolves to the same `StructType` model, with params as fields
- FQN registration should work the same at all depths

### Pass 2 — Lowerer: lift all boc structs to package level

- Pre-pass to collect all boc declarations at any nesting depth
- For each, emit: `_fqnBoc` struct, `Call()` method (the body), named methods for inner bocs
- In the enclosing function body, emit instance creation (`&_fqnBoc{}`) at the declaration point
- File-scope bocs get package-level `var Boc = &_bocBoc{}`; local bocs get local `_boc := &_bocBoc{}`

**Acceptance criterion**: `UPDATE_GOLDEN=1` on `39_local_boc_recursive` and `37_local_boc_var` should produce lifted structs and `BocGroup` coordination instead of `var f any` and inline `.Force()`.

### Pass 3 — Unify lowering paths

- Merge `lowerTopLevel`, `lowerBocBody`, `lowerClosureBody`, `lowerBocAsStmts` into one path
- Remove `localBocVars` tracking and `var f any` hacks
- BocWithSig lowering becomes struct emission, not Go function emission

**Note**: goldens 27, 34, 05 (HOF callbacks, while loop bodies) should NOT be updated until the design question about synchronous callback bocs is resolved.

---

## Risks and Unknowns

- **Go local type restriction**: Go does not allow methods on locally-declared types. All boc structs must be lifted to package level. FQN-based naming avoids collisions.
- **Uppercase outer boc + lowercase inner boc**: `Foo: { bar: { count: 0 } }` — `bar` is a per-`Foo`-instance singleton. `_foo_barBoc` must be instantiated inside `NewFoo()`, not as a package-level var. This is the main structural difference from top-level singletons.
- **Calling the boc itself**: `counter()` vs `counter.increment()` — the generated `Call()` method runs the full body and reinitializes fields. This is correct by design but has no golden test yet. The current struct model just doesn't generate `Call()` at all.
- **BocWithSig field persistence**: `countdown(3)` should leave `countdown.n == 3` after the call. The current Go-function lowering loses this. Test coverage for this semantic does not yet exist.
- **Actor queue / SWMR**: concurrent calls to the same singleton from different goroutines currently race on struct fields. This is a known deferred issue — struct emission can land without fixing it.

## What NOT to Change Yet

- **Directory and file bocs** — defer until in-file nesting works correctly
- **Actor queue** — separate concern; struct emission lands first
- **HOF callback semantics** — resolve the design question before touching goldens 27/34/05
- **FQN cross-file references to nested bocs** — defer until single-file case is solid
- **BocWithSig field persistence tests** — add golden tests for `foo.a` access after call once struct model is in place
