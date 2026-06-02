#feature
# Deps — Dependency Declarations

## Purpose

`Deps` processes dependency declarations from a `_project.yz` companion file. It fetches and caches dependencies, wires them into the build as source roots, and keeps the project code free of manifest boilerplate.

## Design

```
my-project/
  _my-project.yz   ← Deps reads this
  my-project.yz    ← code stays clean
```

```yz
// _my-project.yz
!:[Deps]
dependencies: [
    { name: "FasterXML"; version: {5,1,0}; uri: "git@github.com/FasterXML/v5" }
    { name: "Other";     version_str: "v2026-1"; uri: "https://example.org/other" }
]
```

```yz
// my-project.yz
FastXML: org.fastxml.FasterXML   // explicit alias — visible, greppable

main: {
    f: FastXML()
    r: f.parse("...")
    ...
}
```

`Deps` fetches each dependency into a cache directory (`~/.yz-cache/` or `./vendor/`), then adds the fetched source as a source root so its bocs are available by FQN.

## Alias convention

Dependencies are explicitly aliased in code (`FastXML: org.fastxml.FasterXML`), not auto-imported. This keeps all names visible and greppable in the source file.

## Status

Design phase. Depends on YZC-0025, YZC-0028, YZC-0085. Supersedes the `yz.toml` approach sketched in YZC-0041.
