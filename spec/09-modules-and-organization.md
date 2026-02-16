# 9. Modules and Code Organization

This chapter defines how Yz source code is organized into files, namespaces, and projects.

## 9.1 Overview

Yz uses the **file system** as its module system. There are no explicit module declarations, package keywords, or import statements. The directory structure **is** the namespace hierarchy.

## 9.2 File = Boc

Each `.yz` file is an implicit top-level boc. All declarations at the top level are the file's exports.

```yz
// file: math.yz
pi: 3.14159

circle_area: {
    radius Decimal
    pi * radius * radius
}
```

## 9.3 Simple Projects

For small projects, all `.yz` files in a single directory form the project:

```
my-project/
  main.yz
  utils.yz
  config.yz
```

All top-level declarations across all files in the directory are accessible to each other without import statements.

## 9.4 Namespaces from Directory Structure

In larger projects, subdirectories create namespaces:

```
my-project/
  main.yz
  net/
    http.yz
    tcp.yz
  data/
    json.yz
    csv.yz
```

Files in subdirectories are accessed via their path:

```yz
// In main.yz
server: net.http.Server(port: 8080)
parsed: data.json.parse(raw_text)
```

## 9.5 Smart Nesting (Namespace Flattening)

When a boc's name **exactly matches** its file name (case-sensitive), the namespace is **flattened** — the file-level namespace is removed and the boc is promoted one level up.

### Example — Flattening

```
net/
  http.yz     ← contains: http: { ... }
```

Without smart nesting: `net.http.http` (redundant)
With smart nesting: `net.http` (flattened — the boc **becomes** the file namespace)

### Example — No Flattening (Case Mismatch)

```
net/
  http.yz     ← contains: Http: { ... }
```

`Http` ≠ `http` (different case) → **no flattening**. Access is: `net.http.Http`

### Rules

1. **Exact case match**: Flattening only occurs when the boc name **exactly matches** the filename (without `.yz`). `http.yz` matches `http`, NOT `Http` or `HTTP`
2. **Single match**: If a file contains exactly one boc whose name matches the filename, the file-level namespace is flattened
3. **Multi-block files**: If a file has multiple top-level declarations, flattening applies only to the matching one; others remain under the file namespace
4. **No match**: If no declaration matches the filename, the file namespace is preserved as-is

### Examples

```
// File: net/http.yz
http: {
    server: { port Int; ... }
    client: { ... }
}
helper: { ... }

// Access:
net.http.server(port: 8080)     // http flattened (matches http.yz exactly)
net.http.helper()               // helper also accessible under net.http
```

```
// File: data/Json.yz
Json: { ... }

// Json ≠ Json.yz filename? Filename is "Json" → matches exactly
// Access: data.Json  (flattened)
```

```
// File: data/json.yz
Json: { ... }

// Json ≠ json → NO flattening (case mismatch)
// Access: data.json.Json
```

## 9.6 Access Control

### Explicit Signature = Public Interface

A boc with an explicit `#(...)` signature exposes only the parameters listed in the signature:

```yz
Counter: {
    count: 0                        // Internal
    increment #() {                 // Public
        count = count + 1
    }
    get #(Int) {                    // Public, returns Int
        count
    }
}

c: Counter()
c.increment()    // OK — public
c.get()          // OK — public
c.count          // ERROR — not in any public signature
```

### No Signature = Everything Public

When no explicit signature is given, a synthetic signature is created that includes all internal variables:

```yz
Point: {
    x Int
    y Int
}
// Synthetic: #(x Int, y Int)  — both fields are public

p: Point(1, 2)
p.x              // OK — public
```

## 9.7 Dependencies

External dependencies are declared in a project configuration file (e.g., `yz.toml` or similar):

```toml
[project]
name = "my-app"
version = "1.0.0"

[dependencies]
http = { git = "github.com/yz-std/http", version = "0.2.0" }
json = { git = "github.com/yz-std/json", version = "0.1.0" }
```

Dependencies are accessed by their declared name as a namespace root:

```yz
server: http.Server(port: 8080)
data: json.parse(raw)
```

> **Note:** The exact configuration format and dependency resolution mechanism are TBD for v0.1.

## 9.8 Entry Point

The program entry point is the top-level boc in the file designated as the entry point (conventionally `main.yz`). All top-level expressions in this file are executed at program start.

```yz
// main.yz
print("Hello, World!")
```

## 9.9 Compiler Warnings

The compiler issues warnings for:

1. **Namespace mismatch**: A file contains a type whose name doesn't match the filename (if strict naming mode is enabled)
2. **Unused imports**: A namespace is referenced but no declarations from it are used
3. **Shadowing**: A local variable shadows a top-level or outer-scope variable

## 9.10 Summary

```
Organization:
  File              = Implicit top-level boc
  Directory          = Namespace
  Smart nesting      = Flatten when boc name matches filename
  Access control     = Explicit #() signature hides internals
  No imports         = File system is the module system
  Dependencies       = Configuration file (yz.toml)
  Entry point        = main.yz (top-level expressions)
```
