# Implementation/tickets/ — update log

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
