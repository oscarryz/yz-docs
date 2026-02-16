# 11. Conformance Tests

This chapter defines the conformance test suite — a collection of canonical Yz programs with expected behavior that validates compiler correctness.

## 11.1 Purpose

Conformance tests serve as:

1. **Living specification** — executable examples that complement the prose spec
2. **Compiler validation** — automated tests that any Yz compiler must pass
3. **Regression prevention** — catch unintended behavior changes

## 11.2 Test Format

Each test is a `.yz` file with a companion `.expected` file:

```
spec/11-conformance-tests/
  01-literals/
    int_literal.yz
    int_literal.expected
    decimal_literal.yz
    decimal_literal.expected
    string_literal.yz
    string_literal.expected
  02-declarations/
    ...
```

### `.yz` File

A complete, self-contained Yz program that prints its output:

```yz
// int_literal.yz
print("42")
print("`42`")
print("`0`")
```

### `.expected` File

The exact expected stdout output:

```
42
42
0
```

## 11.3 Test Categories

### 01 — Literals

| Test | Validates |
|------|-----------|
| `int_literal` | Integer literal parsing |
| `decimal_literal` | Decimal literal parsing |
| `string_single_quote` | Single-quoted strings |
| `string_double_quote` | Double-quoted strings |
| `string_interpolation` | Backtick interpolation |
| `string_escape` | Escape sequences (`\n`, `\t`, etc.) |

### 02 — Declarations

| Test | Validates |
|------|-----------|
| `short_decl` | `: ` declaration with inference |
| `typed_decl` | Explicit type declaration |
| `typed_decl_init` | Typed declaration with default |
| `multiple_decl` | Multiple declarations (commas) |

### 03 — Expressions

| Test | Validates |
|------|-----------|
| `arithmetic` | `+`, `-`, `*`, `/`, `%` on Int |
| `left_to_right` | `1 + 2 * 3 = 9` (no precedence) |
| `parens` | `1 + (2 * 3) = 7` (explicit grouping) |
| `comparison` | `<`, `>`, `<=`, `>=` |
| `equality` | `==`, `!=` on various types |
| `logical` | `&&`, `||` on Bool |
| `conditional` | `?` method on Bool |
| `unary_neg` | Unary `-` |
| `string_concat` | `+` on String |

### 04 — Bocs

| Test | Validates |
|------|-----------|
| `simple_boc` | Boc creation and invocation |
| `boc_params` | Required and optional parameters |
| `boc_return` | Return value (last expression) |
| `boc_multi_return` | Multiple return values |
| `boc_nested` | Nested bocs and scoping |
| `boc_closure` | Variable capture by reference |
| `boc_as_value` | Passing bocs as arguments |
| `boc_named_args` | Named argument invocation |

### 05 — Types

| Test | Validates |
|------|-----------|
| `type_create` | User-defined type instantiation |
| `type_fields` | Field access and mutation |
| `structural_compat` | Width subtyping assignment |
| `structural_method` | Structural compatibility for method params |
| `variant_create` | Variant constructor usage |
| `variant_match` | Match on variant discriminant |
| `option_some_none` | Option type usage |
| `result_ok_err` | Result type usage |
| `generic_infer` | Generic type parameter inference |
| `mix_basic` | Mix composition |
| `mix_conflict` | Mix conflict = compile error |

### 06 — Control Flow

| Test | Validates |
|------|-----------|
| `conditional_true` | `?` with true condition |
| `conditional_false` | `?` with false condition |
| `match_condition` | Condition-based match |
| `match_variant` | Variant-based match |
| `match_default` | Default branch |
| `match_continue` | Fallthrough with `continue` |
| `each_range` | `1.to(n).each(...)` |
| `each_array` | Array iteration |
| `while_loop` | `while(cond, body)` |
| `break_in_loop` | `break` exits iteration |
| `continue_in_loop` | `continue` skips iteration |
| `return_early` | `return` exits boc |

### 07 — Collections

| Test | Validates |
|------|-----------|
| `array_literal` | Array creation |
| `array_access` | Index access |
| `array_append` | `<<` operator |
| `array_map` | Transform |
| `array_filter` | Filter |
| `array_reduce` | Fold |
| `dict_literal` | Dictionary creation |
| `dict_access` | Key access |
| `dict_set` | Set key-value |
| `dict_iteration` | Iterate entries |

### 08 — Concurrency

| Test | Validates |
|------|-----------|
| `async_basic` | Basic async invocation returns thunk |
| `async_chain` | Thunk dependency chain |
| `async_materialize` | IO triggers materialization |
| `actor_sequential` | Actor processes messages in order |
| `structured_concurrency` | Parent waits for children |

### 09 — Modules

| Test | Validates |
|------|-----------|
| `multi_file` | Cross-file access (same directory) |
| `subdirectory` | Namespace from subdirectory |
| `smart_nesting` | Namespace flattening (case-sensitive) |
| `access_control` | Explicit signature hides internals |

### 10 — Error Cases (Compile Errors)

| Test | Expected Error |
|------|----------------|
| `type_mismatch` | Assigning incompatible types |
| `undefined_var` | Using undeclared variable |
| `undefined_method` | Calling non-existent method |
| `mix_conflict` | Conflicting field names in mix |
| `missing_param` | Required parameter not provided |
| `return_type_mismatch` | Return type doesn't match signature |

## 11.4 Running Tests

```sh
yz test spec/11-conformance-tests/
```

The test runner:
1. Compiles each `.yz` file
2. Runs the compiled program
3. Captures stdout
4. Compares to the `.expected` file
5. Reports pass/fail

Error case tests expect a compilation error (non-zero exit) with a specific error message pattern.
