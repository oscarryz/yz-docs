# Yz Compiler — Implementation Decisions

Resolved 2026-03-04, prior to starting compiler implementation.

---

## Syntax & Lexer

| # | Decision | Resolution |
|---|----------|------------|
| 1 | ASI rules | Only insert `;`, never commas. Commas are explicit in `()` and `[]`. |
| 2 | Multi-line strings | Strings are multi-line by default. No raw/heredoc syntax needed. |
| 3 | `=` semantics | `=` is **not** an expression — it is a statement. |
| 4 | `:` semantics | `:` (short declaration) is a statement, not an expression. |

## Assignment & Destructuring

| # | Decision | Resolution |
|---|----------|------------|
| 5 | Multiple assignment scope | Support both simple multi-return (`a, b = foo()`) **and** nested destructuring with parentheses (`a, (b, c) = foo(), bar()`). |

## Boc Semantics

| # | Decision | Resolution |
|---|----------|------------|
| 6 | Lowercase boc state sharing | State is **global** (singleton). New invocation does **not** reset state. The boc body runs once and state persists. Subsequent invocations share the same state. Example: `counter(); counter.increment(); counter()` → last call returns `1`. |

## Implementation Strategy

| # | Decision | Resolution |
|---|----------|------------|
| 7 | Compiler language | **Go** |
| 8 | Compiler architecture | Emit **Go source code** → then `go build` (not AST, not interpreter) |
| 9 | v0.1 milestone scope | **Full** — all spec features |
| 10 | Repository | `/Users/oscar/code/github/oscarryz/yz` — **new branch**, start from scratch (no reuse of existing code). Git: `https://github.com/oscarryz/yz` |

## Runtime Semantics

| # | Decision | Resolution |
|---|----------|------------|
| 11 | Unhandled errors | **Panic** (crash) |
| 12 | `nil` | **No `nil` concept.** `nil` is a valid identifier but has no special meaning. Use `Result`/`Option` instead. |
| 13 | `&&`/`||` short-circuit | They are regular methods, **not** compiler-special-cased. Short-circuit behavior is natural because the argument is a **boc** (lazy). e.g., `a \|\| { expensive() }` — the boc is only called if `a` is false. |
| 14 | Dependency config | No TOML. Use a **`.yz` file** (e.g., `project.yz`) implementing a project interface. Minimal format TBD. |
