# Ticket registry

One file per ticket, named by its permanent `YZC-####` id. The file is the ticket's full
description: the problem, what was observed, and what is still undecided.

This is a newer format than [`../tasks-detail.md`](../tasks-detail.md). Tickets filed before
`YZC-0124` still have their detail there, or in [`../tasks-done.md`](../tasks-done.md) once closed.
[`../tasks.md`](../tasks.md) remains the open index for both formats, and links here for tickets
that have a file.

## Minting a new id

1. Take the "Next available" number from the top of [`../tasks.md`](../tasks.md) and bump it. The
   number space is shared with tickets that have no file here. Never reuse a number.
2. Create `YZC-####.md` in this directory with OKF front matter: `type` (`"Bug"` or
   `"Feature Idea"`), `title`, `description`, `tags`, `status`, `implementation: "proposed"`,
   `generated`, and `sources` for every document or source file the ticket relies on.
3. Add a line to [`../tasks.md`](../tasks.md) under Open Tickets that links to the file. Do not
   add an entry to `tasks-detail.md`.
4. Add a row to the table below and an entry to [`log.md`](log.md).

## Tickets

| Id | Title | Implementation |
|---|---|---|
| [YZC-0124](YZC-0124.md) | Boc body is not atomic across inner calls | proposed |
