# Ticket registry

One file per ticket, named by its permanent `YZC-####` id. The file is the ticket's full
description: the problem, what was observed, and what is still undecided.

Every open ticket has a file here, and [`../tasks.md`](../tasks.md) remains the open index that
links to them. Tickets closed before this format have their detail in
[`../tasks-done.md`](../tasks-done.md). [`../tasks-detail.md`](../tasks-detail.md) holds only legacy
entries that are not in the open index.

## Minting a new id

1. Take the "Next available" number from the top of [`../tasks.md`](../tasks.md) and bump it. The
   number space is shared with closed tickets that have no file here. Never reuse a number.
2. Create `YZC-####.md` in this directory with OKF front matter: `type` (`"Bug"` or
   `"Feature Idea"`), `title`, `description`, `tags`, `status`, `implementation: "proposed"`,
   `generated`, and `sources` for every document or source file the ticket relies on.
3. Add a line to [`../tasks.md`](../tasks.md) under Open Tickets that links to the file. Do not
   add an entry to `tasks-detail.md`.
4. Add a row to the table below and an entry to [`log.md`](log.md).

## Tickets

| Id | Title | Implementation |
|---|---|---|
| [YZC-0009](YZC-0009.md) | Range iteration | proposed |
| [YZC-0013](YZC-0013.md) | Array `<<` append | proposed |
| [YZC-0014](YZC-0014.md) | Option/Result method chaining | proposed |
| [YZC-0016](YZC-0016.md) | String `++` concatenation | proposed |
| [YZC-0019](YZC-0019.md) | `break`/`continue`/`return` in loops | proposed |
| [YZC-0023](YZC-0023.md) | Cancellation / non-local return | proposed |
| [YZC-0024](YZC-0024.md) | `return`, `break`, `continue` (major) | proposed |
| [YZC-0031](YZC-0031.md) | Scalar Types in Yz Source (uppering) | proposed |
| [YZC-0039](YZC-0039.md) | Operators audit | proposed |
| [YZC-0042](YZC-0042.md) | `yz` tool: run, new, add, init (wraps yzc + yz fetch) | proposed |
| [YZC-0044](YZC-0044.md) | Producer-consumer example and golden test | proposed |
| [YZC-0058](YZC-0058.md) | GoSource: Go-backed type implementations | proposed |
| [YZC-0060](YZC-0060.md) | Implement `self` as a lexically-resolved compiler built-in | proposed |
| [YZC-0088](YZC-0088.md) | Codegen: attach compiled annotation boc to declaration metadata | proposed |
| [YZC-0096](YZC-0096.md) | `yz fetch`: dependency fetcher | proposed |
| [YZC-0099](YZC-0099.md) | `self` / native-builtins: confirm go_source calling convention covers built-in operators | proposed |
| [YZC-0112](YZC-0112.md) | `continue` in a cond-match arm does not fall through | proposed |
| [YZC-0116](YZC-0116.md) | A non-leaf held-cown call inside a match arm or closure body deadlocks | proposed |
| [YZC-0118](YZC-0118.md) | Macro run-output cache serves stale output after a macro edit | proposed |
| [YZC-0121](YZC-0121.md) | A call to an unresolvable identifier fails only at `go build` | proposed |
| [YZC-0122](YZC-0122.md) | Same-named file-wrapper boc in signature form fails at `go build` | implemented |
| [YZC-0123](YZC-0123.md) | A local declared without an initializer is dropped | proposed |
| [YZC-0124](YZC-0124.md) | Boc body is not atomic across inner calls | proposed |
| [YZC-0125](YZC-0125.md) | Sibling boc calls run sequentially | proposed |
| [YZC-0126](YZC-0126.md) | Constructor call at the start of a statement fails to parse in a type's methods | proposed |
| [YZC-0127](YZC-0127.md) | Boc-typed fields of a type are dropped | proposed |
| [YZC-0128](YZC-0128.md) | A type's methods are not visible on its instances inside its own methods | proposed |
| [YZC-0129](YZC-0129.md) | A singleton and a type whose names differ only in case fail to build | proposed |
