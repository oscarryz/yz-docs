# Implementation/tickets/ — update log

## 2026-10-30

**Creation** — [`YZC-0009`](YZC-0009.md): range iteration. Moved from `tasks-detail.md`.

**Creation** — [`YZC-0013`](YZC-0013.md): array `<<` append. Moved from `tasks-detail.md`.

**Creation** — [`YZC-0014`](YZC-0014.md): Option/Result method chaining. Moved from `tasks-detail.md`.

**Creation** — [`YZC-0016`](YZC-0016.md): string `++` concatenation. Moved from `tasks-detail.md`.

**Creation** — [`YZC-0019`](YZC-0019.md): `break`/`continue`/`return` in loops. Moved from `tasks-detail.md`.

**Creation** — [`YZC-0023`](YZC-0023.md): cancellation / non-local return. Moved from `tasks-detail.md`.

**Creation** — [`YZC-0024`](YZC-0024.md): `return`, `break`, `continue` (major). Moved from `tasks-detail.md`.

**Creation** — [`YZC-0031`](YZC-0031.md): Scalar Types in Yz Source (uppering). Moved from `tasks-detail.md`.

**Creation** — [`YZC-0039`](YZC-0039.md): operators audit. Moved from `tasks-detail.md`.

**Creation** — [`YZC-0042`](YZC-0042.md): `yz` tool: run, new, add, init (wraps yzc + yz fetch). Moved from `tasks-detail.md`.

**Creation** — [`YZC-0044`](YZC-0044.md): producer-consumer example and golden test. Moved from `tasks-detail.md`.

**Creation** — [`YZC-0058`](YZC-0058.md): GoSource: Go-backed type implementations. Moved from `tasks-detail.md`.

**Creation** — [`YZC-0060`](YZC-0060.md): implement `self` as a lexically-resolved compiler built-in. Moved from `tasks-detail.md`.

**Creation** — [`YZC-0088`](YZC-0088.md): codegen: attach compiled annotation boc to declaration metadata. Moved from `tasks-detail.md`.

**Creation** — [`YZC-0096`](YZC-0096.md): `yz fetch`: dependency fetcher. Moved from `tasks-detail.md`.

**Creation** — [`YZC-0099`](YZC-0099.md): `self` / native-builtins: confirm go_source calling convention covers built-in operators. Moved from `tasks-detail.md`.

**Creation** — [`YZC-0112`](YZC-0112.md): `continue` in a cond-match arm does not fall through. Moved from `tasks-detail.md`.

**Creation** — [`YZC-0116`](YZC-0116.md): a non-leaf held-cown call inside a match arm or closure body deadlocks. Moved from `tasks-detail.md`.

**Creation** — [`YZC-0118`](YZC-0118.md): macro run-output cache serves stale output after a macro edit. Moved from `tasks-detail.md`.

**Creation** — [`YZC-0121`](YZC-0121.md): a call to an unresolvable identifier fails only at `go build`. Moved from `tasks-detail.md`.

**Creation** — [`YZC-0122`](YZC-0122.md): same-named file-wrapper boc in signature form fails at `go build`. Moved from `tasks-detail.md`.

**Creation** — [`YZC-0123`](YZC-0123.md): a local declared without an initializer is dropped. Moved from `tasks-detail.md`.

**Update** — [`index.md`](index.md): every open ticket now has a file; added the 22 moved tickets.

## 2026-10-27

**Creation** — [`YZC-0125`](YZC-0125.md): sibling boc calls run sequentially.

**Creation** — [`YZC-0126`](YZC-0126.md): constructor call at the start of a statement fails to parse in a type's methods.

**Creation** — [`YZC-0127`](YZC-0127.md): boc-typed fields of a type are dropped.

**Creation** — [`YZC-0128`](YZC-0128.md): a type's methods are not visible on its instances inside its own methods.

**Creation** — [`YZC-0129`](YZC-0129.md): a singleton and a type whose names differ only in case fail to build.

**Update** — [`YZC-0124`](YZC-0124.md): links YZC-0125 under Related; adds the thin `while` dispatcher to what has to keep working.

**Update** — [`index.md`](index.md): added YZC-0125 to YZC-0129.

## 2026-10-13

**Creation** — [`YZC-0124`](YZC-0124.md): boc body is not atomic across inner calls. Filed while writing a concurrency demo; describes the problem and an undecided proposal.

**Creation** — [`index.md`](index.md): the registry for per-ticket files, with the rule for minting an id.
