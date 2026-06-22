#feature
# Dependencies

## Overview

Dependencies are declared as passive annotation metadata in a `project.info` companion file. The compiler validates the annotation shape but never fetches or resolves anything — that is the responsibility of `yz fetch`, a separate tool that runs before the compiler is invoked. This keeps compilation predictable: no external processes are triggered as a side effect of annotation processing.

---

## Declaration format

Dependencies live inside the `project:` boc in `project.info` at the project root:

```yz
// project.info
project: {
    name: "my_project"
    dependencies: [
        yz_core:  { url: "https://github.com/oscarryz/yz-core" }
        some_lib: { url: "https://github.com/example/lib", ref: "v1.2.0" }
        local_db: { path: "../local-db" }
    ]
}
```

Each entry in `dependencies` is a named boc. The name (`yz_core`, `some_lib`, `local_db`) is a handle used by the tooling — it appears in error messages and in the lock file. It is **not** an import alias in code; users alias explicitly:

```yz
// main.yz
MyLib: org.foo.bar.my_lib
```

### Dependency fields

| Field  | Required | Description |
|--------|----------|-------------|
| `url`  | yes (unless `path`) | Git repository URL |
| `ref`  | no | Branch or tag name. If omitted, `yz fetch` resolves the default branch HEAD on first fetch. |
| `path` | yes (unless `url`) | Local filesystem path, relative to `project.info`. For development of sibling packages. |

`url` and `path` are mutually exclusive.

---

## Resolution and the lock file

`yz fetch` resolves the manifest to exact content hashes and writes `yz.lock` alongside `project.info`. The lock file is committed to version control and makes builds reproducible.

**First fetch** (no lock entry yet):
- `url` only → resolves default branch HEAD to a commit SHA
- `url` + `ref` → resolves the ref to a commit SHA
- Downloads the source, computes a content hash, writes both to `yz.lock`

**Subsequent fetches** (lock entry exists):
- Uses the pinned SHA from `yz.lock`; ignores `ref`
- Skips download if already cached

### Lock file format

`yz.lock` is a Yz array — machine-written, human-readable, diff-friendly:

```yz
[
    {
        name:         "yz_core"
        url:          "https://github.com/oscarryz/yz-core"
        sha:          "abc123def456..."
        content_hash: "sha256:9f86d08..."
    }
    {
        name:         "some_lib"
        url:          "https://github.com/example/lib"
        ref:          "v1.2.0"
        sha:          "def456abc123..."
        content_hash: "sha256:3f4a2b1..."
    }
]
```

`ref` is preserved in the lock file for human readability. `sha` is what governs resolution. `content_hash` verifies the downloaded source has not been tampered with.

Local path deps (`path:`) are not written to `yz.lock` — they are always resolved from disk.

---

## Cache

`yz fetch` stores downloaded sources in a global cache: `~/.yz/cache/<sha>/`. The compiler receives cached source roots as extra arguments to `yzc build`, consistent with multi-root compilation (YZC-0022).

---

## Standard library

The Yz standard library is not a dependency — it is implicit, like Go's stdlib. It is never declared in `project.info`.

---

## Small programs without a project

A `project.info` file is not required for programs that have no external dependencies. When `project.info` is absent, `yz fetch` is a no-op and the compiler is invoked directly. Dependency support for single-file programs without a `project.info` is deferred.

---

## Tooling

| Tool | Role |
|------|------|
| `yz fetch` | Resolves, downloads, and caches dependencies; writes `yz.lock` |
| `yz add <url>` | Adds a dependency entry to `project.info` and runs `yz fetch` |
| `yzc build` | Receives cached source roots as extra arguments; knows nothing about deps |

See also: [Annotations](Annotations.md) · [Code organization](Code%20organization.md)
