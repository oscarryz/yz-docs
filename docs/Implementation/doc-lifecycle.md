---
type: impl
generated: { by: "oscarryz", at: 2026-09-24T19:00:00+02:00 }
verified: { by: "oscarryz", at: 2026-09-24T17:17:00+02:00 }
sources:
  - resource: https://raw.githubusercontent.com/GoogleCloudPlatform/open-knowledge-format/refs/heads/main/SPEC.md
    title: Open Knowledge Format (OKF) specification
---
#impl
# Doc Lifecycle for `docs/Implementation/` (OKF)

This directory is an OKF bundle: `index.md` is the table of contents, `log.md` is the change
history, and every other file is a doc with OKF front matter. This note is the rule set for
keeping those three kinds of file correct — read it before adding, changing, or removing anything
here.

---

## Plan docs are temporary

A plan doc (`<ticket>-plan.md`, `type: impl-plan`) exists to get a ticket implemented. It is not
the permanent record of the ticket.

- **Create** one when a ticket needs a design decided before coding starts (see
  `yzc-0066-plan.md` for the shape: root cause / fix design / test plan / risks).
- **Delete** it once the ticket is implemented, verified, and closed into `tasks-done.md`. The
  code, the tests, and the `tasks-done.md` entry are what future readers should trust — a plan doc
  left in place next to them only invites drift (it will describe intent, not what shipped).
- **If implementation can't fully close the ticket**: file a new ticket in `tasks.md` /
  `tasks-detail.md` for the remaining gap, and either trim the plan doc down to the outstanding
  work or delete it and let the new ticket's own entry carry the gap forward. Never leave a plan
  doc that describes work already done — trim it to what's still open, or remove it.
- **If a durable design insight surfaces while implementing** (true regardless of the ticket, not
  tied to its specifics): it does not belong in the plan doc either. Fold it into whichever living
  design doc already owns that topic (e.g. `concurrency-design.md` for cown/scheduling behavior),
  as a short dated addendum — not a new file, unless no existing doc owns the topic. A rule or a
  fact, not a narrative of how it was found (see `conformance-golden-tests.md` for the target
  shape: one trap, one rule, no session trace).

## `index.md`

- One line per file: `[name](path) — one-sentence description of what the file covers.`
- Update it in the same change that adds, renames, or removes a file in this directory. A stale
  or missing entry is a bug in the index, not a follow-up.

## `log.md`

- Audience: someone looking at this directory **after a merge**, trying to understand what
  changed — not someone replaying the session that produced it.
- One entry per changed file, per day (ideally per PR-sized change), not one entry per edit.
  Several touches to the same file in one session collapse into a single line describing the net
  change, not each intermediate step.
- State the *what* only — a phrase, not a paragraph. The *why* and *how* live in the file itself,
  or in `tasks-done.md` for a closed ticket. If a log line is restating content the target file
  already has, cut it.
- Never log a file that isn't in the tree at merge time. A plan doc created and deleted the same
  day gets no `log.md` line at all — it was scratch work, and the log reflects the result, not the
  process.
