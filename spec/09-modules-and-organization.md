#spec
# 9. Modules and Code Organization

This chapter defines how Yz source code is organized into files, namespaces, and projects.

## 9.1 Overview

Yz uses the **file system** as its module system. There are no explicit module declarations, package keywords, or import statements. The directory structure **is** the namespace hierarchy, governed by five invariants.

## 9.2 Core Invariants

### Invariant 1 — File is boc body

The content of a `.yz` file is the **boc body** of the boc named after the file. The file name (without `.yz`) is the boc name; its content is everything inside `{ }`.

```
// net/http.yz — defines the boc net.http
Server: { port Int }
get: { uri String; ... }
```

This is equivalent to writing `http: { Server: { port Int }; get: { ... } }` inside a `net` boc.

**Uppercase file names** define a struct type. The file name becomes the type name; the content becomes the struct's fields and methods.

```
// Pet.yz — defines the struct type Pet
name String
age  Int
```

Equivalent to `Pet: { name String; age Int }`. Constructable anywhere in the package as `Pet(name: "Rex", age: 3)`.

**Explicit same-named inner boc** — if a root file's content includes a boc with the same name as the file, that inner boc is treated as the active boc body. Other declarations in the same file are package-level peers, not members of the wrapper.

```
// main.yz
counter: { count: 0; increment: { count = count + 1 } }

main: {
    counter.increment()
    print("${counter.count}")
}
main()
```

Here `counter` and `main` are both package-level declarations; `main()` invokes the entry boc.

### Invariant 2 — Directory is boc namespace

A directory defines the boc for its name. Files inside the directory **compose** its boc body — each file contributes a named sub-boc. No two files can conflict because the filesystem enforces name uniqueness.

```
net/
  http.yz    → net.http  (sub-boc of net)
  tcp.yz     → net.tcp   (sub-boc of net)
```

Result: `net: { http: { ... }; tcp: { ... } }`

A directory with no matching `.yz` file is a valid boc — it exists as a namespace with only sub-bocs and no own body.

### Invariant 3 — Source root is not part of the FQN

The source root is a namespace anchor, not a boc. `src/net/http.yz` defines `net.http`, not `src.net.http`.

Multiple source roots with the same FQN path merge into one boc:

```
src/net/http.yz    ↘
lib/net/http.yz    → net.http  (merged)
```

### Invariant 4 — `name.info` is an annotation companion, never a boc

A `name.info` file is not a boc. Its content is an annotation body attached to the boc sharing its base name. The compiler parses it, attaches it to that boc's annotation slot, and triggers any declared macros. What the annotation content means is up to those macros.

```
net/
  http.info   ← annotation companion for net.http (never a boc itself)
  http.yz     ← defines net.http
```

### Invariant 5 — File and directory with same stem coexist

`net/http.yz` defines `net.http`'s own body. `net/http/` contains sub-bocs of `net.http`. Both can exist together:

```
net/
  http.yz       ← net.http body (own fields, methods)
  http/
    client.yz   ← net.http.client (sub-boc)
    server.yz   ← net.http.server (sub-boc)
```

Result: `net.http` has both its own body and the sub-bocs `client` and `server`.

---

## 9.3 Simple Projects

For small projects, all `.yz` files in a single directory form the project:

```
my-project/
  main.yz
  utils.yz
  config.yz
```

All bocs are accessible to each other by name — no import statements needed.

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

```yz
// In main.yz
server: net.http.Server(port: 8080)
parsed: data.json.parse(raw_text)
```

## 9.5 Access Control

### Explicit Interface = Public Interface

A boc declaration with an explicit `#(...)` interface exposes only the fields listed in it:

```yz
Counter #(increment #(), get #(Int)) {
    count: 0
    increment #() { count = count + 1 }
    get #(Int) { count }
}

c: Counter()
c.increment()    // OK — in interface
c.count          // ERROR — not in interface
```

### No Interface = Everything Public

A short boc declaration infers its interface from the body — all declared variables become public:

```yz
Point: { x Int; y Int }
p: Point(1, 2)
p.x              // OK — public
```

## 9.6 Dependencies

External dependencies are fetched and placed in a vendor or cache source root. The `my-project.info` annotation companion (processed by the `Deps` macro — see YZC-0041) declares what is needed. The project code stays clean:

```
my-project/
  my-project.info   ← dependency declarations (processed by Deps macro)
  my-project.yz     ← code
  src/
    ...
```

> **Note:** Dependency resolution is TBD (YZC-0041/0042).

## 9.7 Entry Point

The entry point is the `main` boc, conventionally in `main.yz`. All declarations in the entry file are available at program start.

```yz
// main.yz
print("Hello, World!")
```

## 9.8 Summary

```
Invariants:
  File content    = boc body named after the file
  Directory       = boc namespace; files compose sub-bocs (no conflicts)
  Source root     = namespace anchor, not part of FQN
  name.info       = annotation companion for the named boc; never a boc itself
  File + dir      = can coexist: file = own body, dir/ = sub-bocs

Access:
  Explicit #()    = only listed fields/methods are public
  No signature    = all fields public (inferred interface)

Entry point:
  main boc in main.yz
```
