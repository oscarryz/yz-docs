# Smart Nesting and Namespace Flattening

## Overview

Yz uses **filesystem-based module organization** where file paths determine block namespaces. The **smart nesting** feature eliminates redundant namespace nesting by automatically flattening blocks whose names match their containing filename.

This feature reduces boilerplate while maintaining clear, predictable code organization.

## Core Principle

When a top-level block name matches its filename, the block is placed directly in the file's namespace rather than nested beneath an intermediate block.

### Example

Without smart nesting:

```javascript
// util/math/sqrt.yz
sqrt: { x Int
  x * x  // Simplified for demo
}
// Creates: util.math.sqrt.sqrt ❌ Redundant nesting
```

With smart nesting:

```javascript
// util/math/sqrt.yz
sqrt: { x Int
  x * x
}
// Creates: util.math.sqrt ✅ Clean namespace
```

## Namespace Flattening Rules

### Rule 1: Name-Filename Matching

A top-level block is flattened if and only if its name matches the filename (without the `.yz` extension).

```javascript
// util/math/sqrt.yz
sqrt: { ... }  // Matches filename → util.math.sqrt (flattened)

// util/math/advanced.yz
sqrt: { ... }  // Does NOT match filename (advanced) → util.math.advanced.sqrt (nested)
```

### Rule 2: Case Sensitivity

Matching is **case-sensitive**. The block name must exactly match the filename.

```javascript
// util/math/Sqrt.yz
sqrt: { ... }  // 's' ≠ 'S' → util.math.Sqrt.sqrt (nested)

// util/math/Sqrt.yz
Sqrt: { ... }  // Exact match → util.math.Sqrt (flattened)
```

### Rule 3: Multiple Blocks in One File

A file can define multiple blocks. Blocks matching the filename are flattened; others are nested beneath the filename.

```javascript
// util/math/sqrt.yz
sqrt: { x Int
  x * x
}

helper: { x Int
  x + 1
}

// Creates:
// util.math.sqrt        (flattened, matches filename)
// util.math.sqrt.helper (nested, does not match)
```

### Rule 4: Subdirectories Create Intermediate Namespaces

Nested directories create intermediate namespaces in the path, regardless of smart nesting.

```javascript
// util/math/impl/sqrt.yz
sqrt: { ... }
// Directory path creates intermediate namespace:
// util.math.impl.sqrt (flattened at file level, but nested in directory)
```

## Compiler Warnings

To prevent mistakes from typos or misnamed blocks, the compiler emits warnings in the following cases:

### Warning 1: Filename-Name Mismatch

When a top-level block name does not match its filename, a warning is issued.

```javascript
// util/math/sqrt.yz
sqrtt: { ... }  // Typo or intentional?
```

**Compiler output:**

```
⚠️  Warning (util/math/sqrt.yz, line 1):
  Block 'sqrtt' at top level does not match filename 'sqrt'.
  This creates the namespace: util.math.sqrt.sqrtt
  
  Did you mean: sqrt?
  Use 'sqrtt' if this is intentional (creates a nested block).
```

**Rationale:** Typos can silently create unexpected namespaces, making them hard to catch at compile time. The warning alerts developers to potential mistakes.

### Warning 2: Case Mismatch

When block and filename differ only in casing, a warning is issued.

```javascript
// util/math/Sqrt.yz
sqrt: { ... }  // Lowercase 's' instead of uppercase 'S'
```

**Compiler output:**

```
⚠️  Warning (util/math/Sqrt.yz, line 1):
  Block 'sqrt' does not match filename 'Sqrt' (case difference detected).
  This creates the namespace: util.math.Sqrt.sqrt
  
  Did you mean: Sqrt (with uppercase 'S')?
```

### Warning 3: Multiple Top-Level Blocks with Same Name

If a file defines multiple blocks with the same name, an error is raised (duplicate definitions are not allowed).

```javascript
// util/math/sqrt.yz
sqrt: { x Int
  x * x
}

sqrt: { x Float
  x * x
}
```

**Compiler output:**

```
❌ Error (util/math/sqrt.yz, line 5):
  Duplicate definition of 'sqrt'. Block 'sqrt' is already defined at line 1.
```

## Common Patterns

### Single-Purpose Files

Most files define one block matching the filename:

```javascript
// util/math/sqrt.yz
sqrt: { x Int
  x * x
}
// Creates: util.math.sqrt ✅
```

### Multi-Block Files with Primary Block

A file can have one block matching the filename (flattened) and other helper blocks (nested):

```javascript
// util/math/sqrt.yz
sqrt: { x Int      // Primary block, matches filename
  helper(x)
}

helper: { x Int    // Helper block, nested beneath sqrt
  x * x
}

// Creates:
// util.math.sqrt         (flattened)
// util.math.sqrt.helper  (nested)
```

### Interface and Implementation Separation

Signatures can be declared in a parent file, with implementations in child files. Both must use the same name for matching:

```javascript
// util/math.yz
sqrt #(Int, Int)          // Signature: util.math.sqrt

// util/math/sqrt.yz
sqrt: { x Int, Int
  x * x
}                         // Implementation: util.math.sqrt (matches, merges with signature)
```

In this case, the signature from the parent and the implementation from the child **merge** into a single block `util.math.sqrt`. If both define the full block (not just a signature), an error is raised.

### Grouping Related Functions

Use a file with a generic name (not matching any block) to group related blocks:

```javascript
// util/math/impl.yz
sqrt: { x Int
  x * x
}

sin: { x Int
  // sin implementation
}

cos: { x Int
  // cos implementation
}

// Creates:
// util.math.impl.sqrt
// util.math.impl.sin
// util.math.impl.cos
```

## Disabling Smart Nesting (Optional Strict Mode)

For projects that want stricter control, an optional **strict naming mode** can be enabled in the project configuration:

```javascript
// project.yz
version: '0.1.0'
strict_naming: true  // Enforce exact filename-block matching
```

When `strict_naming` is enabled:

- **Warning becomes error**: Any filename-name mismatch results in a compilation error, not a warning.
- **Recommended for large codebases**: Ensures consistent naming conventions across the project.

Example with strict mode enabled:

```javascript
// util/math/sqrt.yz
sqrtt: { ... }  // Typo

// ❌ Compilation Error:
// Block 'sqrtt' does not match filename 'sqrt'.
// strict_naming is enabled: mismatches are not allowed.
```

## Implementation Notes for Compiler Authors

1. **Parsing Phase**: When building the namespace tree, check if the top-level block name matches the filename.
    
2. **Flattening Logic**:
    
    - Extract filename from path (remove `.yz` extension)
    - Compare against top-level block names (case-sensitive)
    - If match found, place block directly in namespace
    - If no match, nest block beneath filename
3. **Merge Detection**: If a parent file defines a signature and a child file defines the implementation (same name), merge them. Raise an error if both define full blocks.
    
4. **Warning Emission**:
    
    - Detect filename-name mismatches during parsing
    - Suggest the correct name in the warning message
    - Respect `strict_naming` configuration: convert warnings to errors if enabled
5. **Casing Sensitivity**: Treat filenames and block names as case-sensitive; emit specific warnings for case-only differences.
    

## Examples

### Example 1: Simple Utility Module

```javascript
// util/math/sqrt.yz
sqrt: { x Int
  x * x
}

// util/math/sin.yz
sin: { x Int
  // sin implementation
}

// Usage in program.yz
main: {
  result1: util.math.sqrt(4)
  result2: util.math.sin(0)
}

// Namespaces created:
// ✅ util.math.sqrt
// ✅ util.math.sin
```

### Example 2: Multi-Block File

```javascript
// util/math/advanced.yz
matrix_multiply: { m1 Matrix, m2 Matrix
  // Implementation
}

matrix_invert: { m Matrix
  // Implementation
}

// Usage in program.yz
main: {
  result: util.math.advanced.matrix_multiply(m1, m2)
}

// Namespaces created:
// ✅ util.math.advanced.matrix_multiply
// ✅ util.math.advanced.matrix_invert
```

### Example 3: Interface + Implementation

```javascript
// util/math.yz
sqrt #(Int, Int)      // Signature only

// util/math/sqrt.yz
sqrt: { x Int, Int    // Full implementation
  x * x
}

// Usage in program.yz
main: {
  result: util.math.sqrt(4)  // Calls the implementation
}

// Merged namespace:
// ✅ util.math.sqrt (signature from parent + implementation from child)
```

### Example 4: Warning Example

```javascript
// util/string/lenght.yz  (typo: should be 'length')
lenght: { s String
  // implementation
}

// Compiler output:
// ⚠️  Warning (util/string/lenght.yz, line 1):
//   Block 'lenght' at top level does not match filename 'lenght'.
//   This creates the namespace: util.string.lenght.lenght
//   Did you mean: lenght (to match filename)?
```

## Migration from Explicit Package Declarations

If migrating from a language with explicit package declarations, Yz's smart nesting eliminates the need for `package` or `module` statements:

### Before (with explicit declarations):

```javascript
package util.math

sqrt: { x Int
  x * x
}
```

### After (with filesystem organization):

```javascript
// util/math/sqrt.yz
sqrt: { x Int
  x * x
}
// Package automatically determined by file path
```