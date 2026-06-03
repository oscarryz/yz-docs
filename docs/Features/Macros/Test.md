#feature
# Test — Test Fragment Companion

## Purpose

`Test` marks a `.yz` file as a test fragment of its target boc, granting it access to the target's internal fields (white-box testing). Tests live next to their source, Go-style.

## Design

```
net/
  http_test.info   ← Test reads this
  http_test.yz     ← test code for net.http
  http.yz          ← net.http source
```

```yz
// http_test.info
!: [Test]
boc: http
test_engine: JUnit4Yz
```

`http_test.yz` is compiled as a fragment of `net.http` — it can see `http`'s internal (non-public) fields. The `Test` macro wires it into the test runner.

## Test isolation

Files with `!: [Test]` in their annotation are excluded from production builds. The test runner source root (or a build condition) separates them.

## Status

Design phase. Depends on YZC-0025, YZC-0028, YZC-0085.
