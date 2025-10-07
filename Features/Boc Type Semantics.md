# Boc Type Semantics

This document describes the semantics of boc types in the Yz compiler, including type inference, assignment behavior, and code generation patterns.

## Core Concepts

### Boc Types
A **boc type** is the type of a block of code (boc). It defines the structure and behavior of a boc, including:
- **Parameters**: Input values the boc accepts
- **Fields**: Internal state variables
- **Return Type**: What the boc returns (explicit or inferred)

### Boc Type Syntax
```yz
// Explicit boc type with parameters and return type
greet #(name String, age Int, String) {
    message: "Hello, " + name
    message  // return the message
}

// Implicit boc type (inferred from contents)
greeter: {
    name: "World"
    println("Hello, " + name)
}
```

## Type Inference

### Explicit Return Types
When a boc has an explicit return type signature:
```yz
greet #(name String, String) {
    name String
    "Hello, " + name
}
```
- The boc returns the **actual return value** (String)
- `greet("Alice")` returns a String, not a boc instance
- Used when you want to return a value directly

### Implicit Return Types
When a boc has no explicit return type:
```yz
greeter: {
    name: "World"
    println("Hello, " + name)
}
```
- The boc returns **Unit** if the last statement returns Unit
- The boc returns the **boc instance** for field access
- Used when you want to access internal state

## Assignment Semantics

The behavior of boc calls depends on whether the result is assigned to a variable:

### Assigned Bocs (Return Instance)
```yz
g: greeter()  // Returns boc instance for field access
println(g.name)  // Access the name field
```
- Returns the boc instance
- Allows field access (`g.name`, `g.age`, etc.)
- Used for data structures

### Unassigned Bocs (Execute Logic)
```yz
say("Hello")  // Executes the boc's logic
```
- **In Yz**: Executes the boc's logic directly
- **In Generated Go**: Calls the generated `run()` method
- Returns Unit (void)
- Used for side effects and execution

## Unified Boc Semantics

All bocs follow the same pattern in the generated Go code:

### Run Method (Generated Code Only)
- **The compiler generates a `run()` method for each boc that returns Unit (void)**
- The `run()` method is NOT part of the Yz language - it's generated for Go compilation
- In Yz, bocs are just blocks of code with fields and logic
- The `run()` method initializes fields and executes the boc's logic in Go

### Code Generation Pattern
```go
// Generated Go pattern: x.run(); x.value
instance.run()        // Execute logic (generated)
instance.fieldName    // Access fields
```

**Note**: The `run()` method is an implementation detail of the Go code generation, not a Yz language feature.

### Behavior Matrix

| Boc Type | Assignment | Behavior | Return Value |
|----------|------------|----------|--------------|
| Explicit Return Type | Any | Return actual value | String/Int/etc. |
| Implicit Return Type | Assigned | Return instance | Boc instance |
| Implicit Return Type | Not Assigned | Call run() | Unit (void) |

## Examples

### Example 1: Explicit Return Type
```yz
// Boc with explicit return type
greet #(name String, String) {
    name String
    "Hello, " + name  // return this string
}

main: {
    message: greet("Alice")  // message is String
    println(message)         // prints "Hello, Alice"
}
```

### Example 2: Implicit Return Type - Assigned
```yz
// Boc without explicit return type
greeter: {
    name: "World"
    println("Hello, " + name)
}

main: {
    g: greeter()      // g is boc instance
    println(g.name)   // prints "World"
}
```

### Example 3: Implicit Return Type - Not Assigned
```yz
// Boc without explicit return type
say #(s String) {
    println(s)
}

main: {
    say("Hello")  // calls run(), prints "Hello"
}
```

### Example 4: User-Defined Types
```yz
// User-defined boc type (uppercase name)
Person #(name String, age Int)

main: {
    p: Person("Alice", 30)  // Creates new instance
    println(p.name)         // prints "Alice"
    println(p.age)          // prints "30"
}
```

## User-Defined Type Instantiation

User-defined types (types with uppercase names) have special instantiation semantics that differ from regular boc calls.

### Instantiation Patterns

#### Positional Arguments
```yz
One: { a String }

main: {
    o: One("hi")  // Positional argument
    println(o.a)  // prints "hi"
}
```

**Generated Go Pattern**:
```go
o := func() *OneImpl { 
    inst := NewOne().(*OneImpl); 
    inst.a = "hi"; 
    return inst 
}()
```

#### Named Arguments
```yz
Person: { name String, age Int }

main: {
    p: Person(name: "Alice", age: 25)  // Named arguments
    println(p.name)  // prints "Alice"
}
```

**Generated Go Pattern**:
```go
p := func() *PersonImpl { 
    inst := NewPerson().(*PersonImpl); 
    inst.SetName("Alice"); 
    inst.SetAge(25); 
    return inst 
}()
```

### Key Semantic Rules

1. **Uppercase Detection**: Types starting with uppercase letters (Unicode-aware) are user-defined types
2. **Constructor Pattern**: Always use constructor pattern, never struct literals
3. **New Instance**: Each call creates a new instance, never reuses global instances
4. **Field Assignment**: Positional arguments map to fields in declaration order
5. **Setter Methods**: Named arguments use setter methods for interface-based types

### Type Resolution

#### Interface-Based Generation
User-defined types generate both interfaces and implementations:

```go
type Boc_One interface {
    Run()
    GetA() string
    SetA(string)
}

type OneImpl struct {
    a string
}

func NewOne() Boc_One {
    return &OneImpl{a: ""}
}
```

#### Field Access
Field access uses getter methods:
```go
o.GetA()  // Instead of o.a
```

#### Field Assignment
Field assignment uses setter methods:
```go
o.SetA("hi")  // Instead of o.a = "hi"
```

### Multi-Pass Processing

The compiler uses multi-pass processing for user-defined types:

1. **First Pass**: Identify all user-defined boc types
2. **Store Types**: Populate `typeFieldNames` map with field information
3. **Second Pass**: Process calls with known type references
4. **Generate Code**: Use constructor pattern with proper type assertions

### Error Handling

#### Common Issues
1. **Undefined Type**: `One` not recognized as user-defined type
2. **Missing Fields**: Positional arguments don't match field count
3. **Type Mismatch**: Arguments don't match expected field types
4. **Interface Access**: Direct field access instead of getter methods

#### Debugging
- Check if type name starts with uppercase letter
- Verify field names in `typeFieldNames` map
- Ensure proper constructor pattern generation
- Validate type assertions in generated code

## Implementation Details

### Multi-Pass Processing
1. **First Pass**: Identify all user-defined boc types
2. **Second Pass**: Process boc types with known references
3. **Prevents**: Duplicate type generation

### Type Tracking
- `typesWithExplicitReturnType`: Bocs with explicit return types
- `typesReturnUnit`: Bocs that return Unit
- `isBocCallAssigned`: Whether current boc call is assigned

### Generated Code vs Yz Language
- **Yz Language**: Bocs are blocks of code with fields and logic
- **Generated Go Code**: Bocs become structs with generated `run()` methods
- **The `run()` method is never written in Yz - it's generated by the compiler**

### Code Generation Patterns

#### Global Instance Pattern
```go
var _greeter_instance = func() *greeter { 
    inst := Newgreeter(); 
    inst.run(); 
    return inst 
}()
```

#### IIFE Pattern (Immediate Invoked Function Expression)
```go
func() interface{} { 
    instance := Newgreet(); 
    instance.run(); 
    return instance.result 
}()
```

## Built-in Operators

In Yz, operators are bocs on built-in types:
```yz
a: 5
b: 3
c: a + b  // Equivalent to a.+(b)
```
- `+` is a boc on the `Int` type
- `a + b` returns the result type (Int), not a boc instance
- All operators follow this pattern

## Field Naming

### Named Parameters
```yz
greet #(name String, age Int)  // name, age
```

### Unnamed Parameters
```yz
greet #(String, Int)  // result1, result2
```

### Result Fields
- Generated automatically for return values
- Named `result`, `result1`, `result2`, etc.
- Used for explicit return types

## Error Handling

### Common Issues
1. **Type Mismatch**: Boc returns instance but String expected
2. **Missing Return Type**: Implicit boc used where explicit needed
3. **Assignment Context**: Wrong behavior for assigned vs unassigned calls

### Debugging
- Check if boc has explicit return type
- Verify assignment context
- Ensure proper field naming

## Future Considerations

### Potential Improvements
1. **Better Type Inference**: Automatic detection of return types
2. **Interface Support**: Boc types as interfaces
3. **Generic Types**: Parameterized boc types
4. **Type Checking**: Stricter compile-time validation

### Current Limitations
1. **Complex Type Inference**: Explicit return types need better support
2. **Error Messages**: Could be more descriptive
3. **Performance**: Multi-pass processing could be optimized
