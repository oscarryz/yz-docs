### 1. Simple Generics (Parametric Polymorphism)

Definition: Types are passed as inputs to a function or structure (e.g., List<T>).
Decision Power: The caller (user) decides what the type is at the moment of use.
The Problem: It can lead to "parameter explosion." If a Graph needs a Node, Edge, and Weight, you must carry all three as external labels every time you reference the Graph (e.g., Graph<N, E, W>).
### 2. Associated Types
Definition: Types are defined as properties inside a module or trait (e.g., type Node inside a Graph signature).
Decision Power: The implementer decides the type once.
The Benefit: It creates a 1-to-1 mapping. If the compiler knows you are using a SocialGraph, it automatically "looks up" that the Node is a User. This encapsulates internal details and keeps function signatures clean.
### 3. The "Path-Dependent" Link
The Challenge: In most languages, if you pass a generic g, the compiler loses the link to its internal types unless you "anchor" it.
The Solutions:
Rust: Uses a type-level placeholder: fn process<G: Graph>(g: G, n: G::Node).
Scala/YZ: Uses the value itself as the path: fn process(g Graph, n g.Node).
Go: Lacks this link, so you are forced back into simple generics: func process[N any](g Graph[N], n N).
### 4. Implementation Styles: Structural vs. Nominal
Nominal (Rust/ML/Java): You must explicitly state implements or impl. This is intentional and prevents accidental matches (e.g., a Gun.draw() vs. a Shape.draw()).
Structural (Go/YZ): If the
"block" has the right members/methods, it satisfies the requirement automatically. This is highly flexible and reduces boilerplate.
### 5. Compilation Strategies
Monomorphization (Rust/C++): The compiler "copies and pastes" code for every specific type used. It's the fastest (zero-cost) but makes the binary larger.
Uniform Representation (OCaml/Java): The compiler writes one version of the function that treats everything as a pointer. It's more compact but requires "boxing/unboxing" (wrapping data in pointers).
### 6. Your "Yz" Language Insights
Blocks as Objects/Functions: By making blocks dual-purpose, you can treat a Type as both a Signature (for validation) and a Constructor (for instantiation).
Closing the Gap: Your use of G.Node in function signatures effectively implements Path-Dependent Types, giving you the architectural power of Scala/ML with a syntax that avoids the "noise" of traditional generics.