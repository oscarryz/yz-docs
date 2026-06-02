#feature
# Build — Conditional File Inclusion

## Purpose

The `Build` compile-time boc decides which `.yz` files are included in a given compilation based on conditions (platform, feature flags, etc.). It operates on `_name.yz` companion files that declare `boc:` and a condition.

## Design

A `_name.yz` companion with a `boc:` field signals that `name.yz` is a conditional body for the named boc. `Build` reads the conditions and includes or excludes the file:

```
foo/
  _print_win.yz    ← Build reads this
  print_win.yz     ← included only on Windows
  _print_macos.yz
  print_macos.yz   ← included only on macOS
  _print_default.yz
  print_default.yz ← included when no other variant matches
```

```yz
// _print_win.yz
!:[Build]
boc: print
platform: Windows64
```

```yz
// _print_default.yz
!:[Build]
boc: print
// no condition — always matches if nothing else does
```

`Build` enforces that at most one variant matches per build target. Overlapping conditions are a compile error. A missing default with no matching variant is also a compile error.

## Disambiguation

A `_name.yz` without a `boc:` field is plain metadata — `Build` ignores it. Only `_name.yz` files containing `!:[Build]` and a `boc:` field are processed as variant declarations.

## Status

Design phase. Depends on YZC-0025 (infostrings), YZC-0028 (compile-time bocs), YZC-0085 (module system).
