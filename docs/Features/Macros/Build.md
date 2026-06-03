#feature
# Build — Conditional File Inclusion

## Purpose

The `Build` macro decides which `.yz` files are included in a given compilation based on conditions (platform, feature flags, etc.). It operates on `name.info` annotation companions that declare `boc:` and a condition.

## Design

A `name.info` companion with a `boc:` field signals that `name.yz` is a conditional body for the named boc. `Build` reads the conditions and includes or excludes the file:

```
foo/
  print_win.info    ← Build reads this
  print_win.yz      ← included only on Windows
  print_macos.info
  print_macos.yz    ← included only on macOS
  print_default.info
  print_default.yz  ← included when no other variant matches
```

```yz
// print_win.info
!: [Build]
boc: print
platform: Windows64
```

```yz
// print_default.info
!: [Build]
boc: print
// no condition — always matches if nothing else does
```

`Build` enforces that at most one variant matches per build target. Overlapping conditions are a compile error. A missing default with no matching variant is also a compile error.

## Disambiguation

A `name.info` without a `boc:` field is plain metadata — `Build` ignores it. Only `name.info` files containing `!: [Build]` and a `boc:` field are processed as variant declarations.

## Status

Design phase. Depends on YZC-0025 (annotations), YZC-0028 (macros), YZC-0085 (module system).
