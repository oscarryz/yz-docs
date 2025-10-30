## Code Composition: The `mix` Keyword

The three-letter keyword `mix` serves as the primary mechanism for **compositional reuse** in Yz, achieving seamless structural injection and scope flattening within a Block of Code (BOC) literal. Unlike typical `import` directives, which establish references requiring qualified access, `mix` performs a compile-time deep merge of the source block's structure and behavior directly into the host block's definition. This mechanism implements the structural **Mixin pattern** at the language level.  

### I. Architectural Semantics

The `mix` operation is a non-trivial structural transformation guided by Yz’s principles of static structural typing. When a source block is mixed in, the compiler enforces three key guarantees:  

1. **Structural Flattening:** The fields (state) and procedures (behavior) defined within the source block are physically integrated into the host block's structure, dissolving the source block's encapsulation boundaries at the point of definition.
    
2. **Unqualified Lexical Scope Promotion:** All merged members are promoted directly into the host block's unqualified local scope. This allows the host block and its consumers to access the mixed-in members directly, without needing a prefix or delegation boilerplate.
    
3. **Stateful Binding (Mixin Behavior):** The source block's procedures maintain their semantic integrity (i.e., closures) and correctly bind to the state variables (fields) now present in the host instance.  
    

This provides a powerful alternative to traditional delegation (explicit composition), as it allows for compositional reuse while maintaining the simplicity of single-name access for all structural components.

### II. Type System Integration

The behavior of `mix` depends critically on the host block's context:

- **In User-Defined Types (UDT):** When used in a UDT (which acts as a factory), `mix` ensures that every instance created via invocation (`Person()`) receives a unique, isolated, mutable copy of the mixed-in state. This facilitates object composition with correct state isolation, aligning with the definition of a stateful mixin.  
    
- **Resulting BOC Type:** The BOC type (`#(...)`) of the host block structurally conforms to the merged set of fields and procedures from the host's native definition and all mixed-in blocks.
    

### III. Conflict Resolution (The Structural Contract)

To preserve the integrity of Yz's structural typing, `mix` employs strict, explicit conflict resolution.

**Mandatory Compilation Error:** If any member name (field or procedure) in the source block conflicts with a member name already defined in the host block, or if conflicts arise between multiple mixed-in blocks, a compilation error is mandatory.

This contract prevents silent shadowing or implicit overrides, guaranteeing that structural integration is always transparent and predictable to the advanced programmer.

### IV. Syntax and Examples

The `mix` keyword is used inside a block literal (`{...}`), referencing the identifier of the block (or UDT) to be merged.

#### Example 1: Structural Flattening and State Binding

Here, the `Named` block defines reusable state (`name`) and behavior (`hi`). The `Person` UDT uses `mix` to absorb this structure.


```js
// --- Reusable Structure (e.g., in named.yz) ---
Named : {
    // State to be injected
    name String 
    
    // Behavior that operates on the injected state
    hi: {
        "My name is `{name}`"
    }
}

// --- Host Block (person.yz) ---
Person : {
    // Structural Merge: 'name' and 'hi' are now local to Person
    mix Named 
    
    // Host-specific state
    last_name String
}

// --- Usage & Resulting Scope ---

// Instance created via UDT:
p : Person("Jon", "Doe") 

// 1. Unqualified access to mixed-in behavior:
p.hi()          // Success: Resolves to Person's 'hi' using Person's 'name' state. 

// 2. Direct access to mixed-in state:
p.name = "Jane" 
p.hi()          // Prints: "My name is Jane"
```

#### Example 2: The Conflict Error

If the host block attempts to define a field that already exists in the mixed-in block, the compiler rejects the operation immediately.


```js
Loggable : {
    level Int // Conflict: 'level' is provided by Loggable
}

Config : {
    mix Loggable 
    
    // Error: Field 'level' is already defined by the Loggable mixin.
    level String // Compilation Error: Conflicting member 'level'
}
```

This strict conflict management ensures that `mix` is a high-integrity mechanism for structural reuse.