# YZC-0028 — Macro System (thin end-to-end slice)

Branch: `macros-sonnet`

## Context

Macros are the last gate before YZC-0031 (scalar types in Yz source) and YZC-0088. Design is settled in `docs/Features/Macros.md`: bocs satisfying `Macro : { Schema #(); run #(subject Boc, config Schema, Boc) }` run at compile time, triggered by uppercase-keyed entries in annotations (`` `Debug: {}` ``); their returned boc's slots merge into the annotated boc. Two-phase build: macro packages compile to native executables first, then the main compilation invokes them via subprocess.

**User-settled decisions:**
1. **Scope — thin slice**: trigger detection, bootstrap build, subprocess + merge, cycle detection, basic caching, one working `Debug` macro, conformance tests. Full reflection API / JSON-Derive-Validate macros deferred.
2. **Wire format — Yz source (non-executable data subset), both directions.** No JSON.
3. **Boc API — minimal**: `Boc : { name String, fields [Field], source String }`, `Field : { name String, type String }`.

Note: the YZC-0028 checklist in tasks-detail.md is stale (`Compile` interface, `macros: [...]` list) — the Features docs are authoritative; update the checklist as part of this ticket.

## Key design resolutions

- **`Boc`/`Field`/`NoConfig` + `generated(src)` helper = compiler-injected Yz prelude**, prepended only when bootstrap-compiling macro packages. No sema builtins, no rt Go types. Macro packages compile **unwrapped** (concatenated stmts, no `IsFileWrapper`) so prelude types resolve and `Debug` is a top-level Go type — mirrors `conformance_test.go`'s in-process `compile()`.
- **Macro detection is a structural AST scan** (no sema): `ShortDecl` with uppercase name whose BocLiteral contains a `Schema` entry (alias `Schema : NoConfig` or assoc-type form) and a `run` BocDecl with 3 sig params (`Boc`, schema type, unlabeled `Boc` return).
- **run→wire mapping**: only the returned Boc's `source` field is used outbound; the prelude helper `generated(src)` wraps it. Macro authors build generated code as source strings.
- **Merge is pre-sema**: after parsing, before `AnalyzeFile`, `expandMacros` walks decls; for each trigger it serializes the subject's *current* elements (progressive merge), invokes the executable, parses stdout, appends `sf.Stmts` to the subject `BocLiteral.Elements` (same pattern as `injectIntoBocLiteral`, build.go:360). Sufficient for this slice — Debug only needs parsed TypedDecl fields. Interleaving-with-inference is a documented follow-up.
- **Macro packages are excluded from the phase-2 app build** (deleted from `byDir`); macro defs in project root = error.
- **Sema needs no change for triggers**: `Debug: {}` inside an annotation is a benign ShortDecl declaration in the annotation's isolated scope. Bare `` `Debug` `` (no colon) form is deferred.
- **Subject annotated-decl form for this slice**: annotated `ShortDecl` (uppercase name + BocLiteral value). `TypedDecl`/`BocDecl` subjects deferred.
- **Verified hazard**: `Field.type` — `type` is not a Yz keyword but lowering emits field names verbatim into Go → invalid Go. Phase 0 fixes this with `goSafeName`.
- **Import chain is legal**: `macrowire` lives at `yz/runtime/macrowire` (module `yz`, non-internal) so the generated `yzapp` macro executable can import it, and it may itself import `yz/internal/parser`.

## Wire format

Inbound (compiler → macro stdin), data subset only (ShortDecl keys; String/Int/Decimal/Bool literals; BocLiterals; arrays):

```
subject: {
    name: "Person"
    fields: [
        { name: "name", type: "String" },
        { name: "age", type: "Int" }
    ]
}
config: { pretty: true }
```

Outbound (macro stdout): a **raw Yz boc body** of generated slots — real code, unquoted:

```
debug #(String) { "Person(" + "name: " + name.to_str() + ")" }
```

Compiler parses stdout with `parser.New(...).ParseFile()`; non-zero exit → compile error quoting stderr.

## Phases (TDD; `make test` green after each; `make test-full` at Phase 5)

### Phase 0 — Go-keyword-safe field names
- `compiler/internal/ir/lower.go`: `goSafeName(name string) string` (append `_` for Go keywords); apply at struct-field emission, ctor params, field-init sites.
- Golden test: struct with `type String` field.

### Phase 1 — `yz/runtime/macrowire` (wire codec)
- **New** `compiler/runtime/macrowire/macrowire.go`: `FieldSpec{Name, Type}`, `ConfigEntry{Key, Value ConfigValue}` (ordered), `Payload{SubjectName, Fields, Config}`, `(p *Payload) Encode() string`, `DecodePayload([]byte) (*Payload, error)` (uses `yz/internal/parser`).
- **New** `compiler/internal/ast/print.go`: `TypeExprString(te TypeExpr) string` (no printer exists today; switch over SimpleTypeExpr+TypeArgs, ArrayTypeExpr, DictTypeExpr, BocTypeExpr, MemberTypeExpr).
- Round-trip unit tests against the real parser (empty fields/config, all scalar kinds, quote escaping).

### Phase 2 — Macro scan + registry (pure functions)
- **New** `compiler/cmd/yzc/macro.go`: `macroDef{Name, RelDir, SchemaTypeName, SchemaFields, Decl}`, `macroRegistry{byName, byDir, binPath, state}`, `scanMacroDefs(stmts, relDir)`, `trigger{Name, Config, Pos}`, `annotationTriggers(ann)`, `validateConfig`, `encodeConfig`, `buildSubjectPayload`.
- Table tests: detection positives + near-misses (missing Schema, 2-param run, lowercase), trigger extraction, config validation, payload building.

### Phase 3 — Prelude + bootstrap build
- `macro.go`: `macroPrelude` const + memoized `preludeStmts()`; `bootstrapMacroPackage(projectDir, files, relDir, reg, stack)` — parse sorted files, concat stmts unwrapped + prelude, reject top-level `main`, recursive `expandMacros` (macro-on-macro), sema→lower→codegen, `genMacroMain(defs)` (Go main dispatching on `os.Args[1]`, decodes stdin via macrowire, builds `Boc`/config values with `std.New*`, calls `NewX().Run(subject, cfg).Force()`, prints `out.Source().GoString()`), then reuse `writeGeneratedGo` (gen: `target/macros/<pkgKey>/gen`) + `goBuild` (bin: `target/macros/<pkgKey>/bin/macros`).
- Caching: `source.hash` = sha256(sources + prelude + yzc version); skip build when hash matches and binary exists. Non-scalar schema field → scan-time error.
- Verify during impl: `NewDebug()` ctor arity (alias `Schema` field should be `IsTypeField`, excluded from ctor); `config Schema`-as-alias param resolution — fallback: author writes concrete type name (`config NoConfig`), log follow-up if alias is flaky.
- Tests: `genMacroMain` output assertions; hash-skip logic. (Real `go build` deferred to Phase 5 harness — too slow for fast gate.)

### Phase 4 — Expansion in the main build
- `macro.go`: `expandMacros(stmts, relDir, reg, stack, projectDir)` — walk top-level stmts + file-wrapper elements for annotated uppercase ShortDecls; per trigger in order: unknown-macro / same-package errors, `ensureMacroPackageBuilt(dir, stack)` (lazy bootstrap + cycle check via visiting stack), progressive re-serialize, `invokeMacro` (run cache at `target/macros/<pkgKey>/runs/<sha256(name+payload)>.out`, else `exec.Command(bin, name)` with payload on stdin), parse stdout, append to subject elements.
- `build.go`: `compileProject` — discovery pass after `byDir` grouping: parse each dir once, `scanMacroDefs`; on hits: root-dir error, duplicate-name error, register, `delete(byDir, dir)`. `compilePackageDir` gains `reg` param (nil-safe); call `expandMacros` before `AnalyzeFile`.
- Tests: cycle-stack + same-package with fake registry (no subprocess); merge test: append canned macro output to a parsed Person, run sema, assert method resolves.

### Phase 5 — Debug macro proof + conformance
- **New example** `compiler/examples/macro_debug/`: `macros/debug.yz` (Debug macro: `Schema : NoConfig`, `run` walks `subject.fields.each({ f Field ... })` building the `debug #(String)` source via `+`/`to_str()`, returns `generated(src)`); `main.yz` with `` `Debug: {}` `` on `Person : { name String, age Int }`; `main.output` — picked up by `TestExamples`.
- **New** `compiler/test/conformance/macro_test.go` — driver-level harness (skipped under `-short`, reuses `buildYzc`), cases under `testdata/macros/<case>/`: `debug_merge` (build+run+output+cache-reuse via bin mtime), `unknown_macro`, `cycle`, `same_package`, `root_macro` (each with `expected.error` substring).
- `make test-full`.

### Phase 6 — Docs + tracking
- `spec/12-macros.md` stub (trigger rule, interface, two-phase build, wire format, slice limits).
- Update `docs/Implementation/tasks.md` + `tasks-detail.md`: fix stale YZC-0028 checklist, mark slice done, record deferred list below as follow-ups. Move detail to `tasks-done.md` per convention.

## Error messages

- `unknown macro "Debgu": no Macro implementation with that name found` (with pos via `diagnostic.Format`)
- `macro cycle detected: macros_a.GenB -> macros_b.GenA -> macros_a.GenB`
- `macro "Debug" is defined in package "macros" and cannot be applied to a boc in the same package; macros must live in a separate package from the bocs they process`
- `macro "Debug" cannot be defined in the project root package; move it to a sub-package (e.g. macros/)`
- `macro "Debug" failed (exit 1):\n<stderr>` / `macro "Debug" returned invalid Yz source: <parse error>`
- `macro "JSON": config key "ignor" does not match any Schema field (expected: field_name String, ignore Bool)`

## Verification

1. Per-phase: unit tests + `make test` (fast gate).
2. End-to-end: `yzc run compiler/examples/macro_debug` prints `Person(name: Ada, age: 36)`; `target/gen/main.go` contains merged `func (self *Person) Debug()`; rebuild reuses cached macro binary.
3. All error cases produce the exact messages above.
4. `make test-full` at end (101+ golden, 25+ error, examples, runtime).

## Deferred (record as follow-ups)

1. Full Structural Reflection Boc API (methods, type_params, annotation slot, field-level annotations)
2. JSON / Derive / Validate std macros
3. Bare `` `Debug` `` trigger form (needs analyzeAnnotationBody exemption)
4. Macro expansion interleaved with inference (inferred field types, generated-constraint attribution, `q.Schema` subjects)
5. Serializing ShortDecl-default fields, methods, type params of the subject
6. Non-scalar/defaulted config; full sema Schema validation via assoc-type machinery
7. Mixed runtime+macro packages (currently excluded from app build)
8. Prelude → stdlib source root (converges with YZC-0031)
9. `name.info` companion triggers; `TypedDecl`/`BocDecl` subjects
10. Structure-hash run caching independent of formatting; cache eviction
11. Spec 12 full text
