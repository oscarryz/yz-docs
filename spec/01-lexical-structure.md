#spec 
# 1. Lexical Structure

This chapter defines the lexical grammar of Yz: how source text is divided into tokens.

## 1.1 Source Text

Yz source code is Unicode text encoded in **UTF-8**. Each source file has the extension `.yz`.

## 1.2 Character Categories

```
letter        = unicode_letter | '_'
unicode_letter = /* any Unicode code point classified as "Letter" */
digit         = '0' … '9'
```

## 1.3 Comments

Yz supports two forms of comments:

```yz
// Single-line comment — extends to end of line

/* Multi-line comment
   can span multiple lines */
```

Comments are treated as whitespace by the lexer and produce no tokens.

## 1.4 Tokens

The lexer produces the following token categories:

1. **Identifiers**
2. **Keywords**
3. **Integer literals**
4. **Decimal literals**
5. **String literals**
6. **Non-word identifiers** (operator-like method names)
7. **Delimiters and punctuation**

## 1.5 Identifiers

An identifier is a sequence of letters and digits beginning with a letter.

```
identifier = letter { letter | digit }
```

Identifiers are **case-sensitive**. The casing of the first character carries semantic meaning:

| First Character | Meaning | Example |
|-----------------|---------|---------|
| Lowercase letter | Variable or block name | `name`, `greet`, `fetch_user` |
| Uppercase letter (not all uppercase single letter) | User-defined type | `Person`, `NetworkResponse` |
| Single uppercase letter | Generic type parameter | `T`, `E`, `K`, `V` |

### Examples

```yz
// Variable / block names
message
count
fetch_user
to_whom

// User-defined types
Person
NetworkResponse
Tree

// Generic type parameters
T
E
```

## 1.6 Keywords

The following identifiers are reserved and cannot be used as variable or type names:

```
break
continue
return
match
mix
```

> **Note:** `true` and `false` are not keywords — they are constants of type `Bool` defined in the standard library.

## 1.7 Integer Literals

```
int_literal = digit { digit }
```

Negative integers are expressed using the unary `-` method (see §1.9).

### Examples

```yz
0
42
1000
```

## 1.8 Decimal Literals

```
decimal_literal = digit { digit } '.' digit { digit }
```

### Examples

```yz
3.14
0.5
100.0
```

## 1.9 Non-Word Identifiers

Non-word identifiers are method names composed of non-alphanumeric, non-delimiter, non-whitespace characters. Any sequence of such characters forms a valid non-word identifier — there is no fixed set.

```
non_word_identifier = non_word_char { non_word_char }
non_word_char       = /* any character that is NOT:
                         - a letter or digit
                         - a delimiter: { } ( ) [ ] : , ; . #
                         - whitespace
                         - a quote: ' "
                         - the backtick ` (infostring delimiter)
                         - the lone character =  (assignment)
                       */
```

> **Important:** `=` alone is NOT a non-word identifier — it is the assignment operator (the only true operator in Yz). However, `==` IS a non-word identifier (a method on all types). Similarly, `=>` is a delimiter (fat arrow), not a non-word identifier.

### Common Non-Word Identifiers on Core Types

The standard library defines these non-word methods on built-in types. They are not special to the parser — they are regular methods like any other:

| Name | Defined On | Description |
|------|-----------|-------------|
| `+` | `Int`, `Decimal`, `String` | Addition / concatenation |
| `-` | `Int`, `Decimal` | Subtraction (binary), negation (unary) |
| `*` | `Int`, `Decimal` | Multiplication |
| `/` | `Int`, `Decimal` | Division |
| `%` | `Int` | Modulo |
| `<` | `Int`, `Decimal` | Less than → `Bool` |
| `>` | `Int`, `Decimal` | Greater than → `Bool` |
| `<=` | `Int`, `Decimal` | Less or equal → `Bool` |
| `>=` | `Int`, `Decimal` | Greater or equal → `Bool` |
| `==` | All types | Structural equality → `Bool` |
| `!=` | All types | Structural inequality → `Bool` |
| `&&` | `Bool` | Logical AND → `Bool` |
| `\|\|` | `Bool` | Logical OR → `Bool` |
| `?` | `Bool` | Conditional — takes two boc arguments |
| `<<` | (user-defined) | Example: append to collection |

User-defined types can define their own non-word methods with any valid non-word identifier name (e.g., `<~>`, `|>`, `@@`, etc.).

### Non-Word Invocation

When a boc's name is a non-word identifier and it has at least one parameter, it can be invoked without `.` and parentheses:

```yz
a + b        // desugars to a.+(b)
x == y       // desugars to x.==(y)
flag ? { "yes" }, { "no" }  // desugars to flag.?(block1, block2)
list << item // desugars to list.<<(item)
```

## 1.10 String Literals

Yz has two string delimiters — single quotes and double quotes. They are equivalent; the closing delimiter must match the opening one.

```
string_literal = single_quoted | double_quoted
single_quoted  = "'" { string_char | '"' | interpolation } "'"
double_quoted  = '"' { string_char | "'" | interpolation } '"'
string_char    = /* any Unicode character except the matching quote */
               | escape_sequence
```

### String Interpolation

Embed any expression inside a string using `${}`:

```yz
name: "Alice"
greeting: "Hello, ${name}!"    // "Hello, Alice!"
math: 'Result: ${1 + 2}'       // "Result: 3"
```

```
interpolation = '$' '{' expression '}'
```

Braces inside the expression are balanced — `${foo({x})}` works correctly.

### Escape Sequences

```
escape_sequence = '\' ( 'n' | 't' | 'r' | '\\' | '\'' | '"' | '0' )
```

| Sequence | Meaning |
|----------|---------|
| `\n` | Newline |
| `\t` | Tab |
| `\r` | Carriage return |
| `\\` | Backslash |
| `\'` | Single quote |
| `\"` | Double quote |
| `\0` | Null character |

> **Open question:** Multi-line strings — whether strings can span multiple lines or if a raw/heredoc syntax is needed is TBD.

## 1.11 Delimiters and Punctuation

| Token | Name | Usage |
|-------|------|-------|
| `{` | Left brace | Block start |
| `}` | Right brace | Block end |
| `(` | Left paren | Invocation, grouping |
| `)` | Right paren | Invocation, grouping |
| `[` | Left bracket | Array / dictionary literal, type |
| `]` | Right bracket | Array / dictionary literal, type |
| `:` | Colon | Short declaration + initialization; dictionary key-value separator |
| `=` | Equals | Assignment (the only true operator in Yz) |
| `,` | Comma | Expression separator inside `()` and `[]` |
| `;` | Semicolon | Statement separator inside `{}` |
| `.` | Dot | Member access |
| `#` | Hash | Boc type signature prefix |
| `=>` | Fat arrow | Condition-action separator in conditional bocs |

## 1.12 Automatic Semicolon Insertion (ASI)

Like Go, Yz uses newline-based automatic semicolon insertion. A semicolon is automatically inserted after a line's final token if that token is one of:

- An **identifier** (word or non-word)
- An **integer literal**, **decimal literal**, or **string literal**
- One of the keywords: `break`, `continue`, `return`
- A closing delimiter: `)`, `]`, `}`

### Examples

```yz
// Source code:
name: "Alice"
age: 30
print("`name` is `age`")

// After ASI:
name: "Alice";
age: 30;
print("`name` is `age`");
```

A semicolon is **NOT** inserted after:

- `{` (block continues on next line)
- `,` (expression list continues)
- `.` (member access continues on next line)
- `=` (assignment continues on next line)
- `:` (short declaration continues on next line)
- `=>` (action continues on next line)
- `#` (signature continues on next line)
- Binary non-word identifiers at end of line (`+`, `-`, `*`, etc.) — the expression continues

> **Note:** The ASI rules may be refined during grammar development. The goal is: "if a newline comes after a token that could end a statement, insert `;`."

## 1.13 Non-Word Method Evaluation Order

All non-word method invocations have **equal precedence** and are evaluated **left-to-right** (like Smalltalk). There is no built-in operator precedence table. Use parentheses to control evaluation order.

The only exception is **unary `-`** (negation), which binds to the immediately following expression.

### Examples

```yz
1 + 2 * 3
// Evaluates left-to-right: (1 .+ 2) .* 3 = 9
// For mathematical precedence, use parentheses:
1 + (2 * 3)   // = 7

x > 0 && x < 10
// Evaluates left-to-right: ((x .> 0) .&& x) .< 10
// Use parentheses for the intended grouping:
(x > 0) && (x < 10)

x == 0 ? { "zero" }, { "nonzero" }
// Evaluates left-to-right: (x .== 0) .? (block1, block2)
// This works naturally because == produces the Bool that ? needs

-n + 1
// Unary - binds first: (n.-()) .+ 1
```

> **Note:** Parentheses are required to express mathematical precedence. `1 + 2 * 3` is `9`, not `7`.

## 1.14 Info Strings

An info string is a string literal that appears immediately before any element (variable, block, type). It attaches metadata that can be retrieved at runtime via `std.info()`.

```yz
`A greeting message`
greeting: "Hello, World!"

// Retrieve:
info(greeting).text  // "A greeting message"
```

Info strings use backtick delimiters at the block level (not to be confused with backtick interpolation *inside* strings). Multi-line info strings use double-quote delimiters:

```yz
"
Description of this block
version: 1.0
author: 'Yz developers'
"
say_hello: { ... }
```

> **Note:** The content inside info string blocks does not need to be valid Yz code.

## 1.15 Token Summary

```
Tokens:
  identifier         : [a-z_][a-zA-Z0-9_]*
  type_identifier    : [A-Z][a-zA-Z0-9_]+
  generic_identifier : [A-Z]
  non_word_identifier: any sequence of non-alphanumeric, non-delimiter, non-whitespace, non-quote chars (excluding lone '=' and '=>')
  int_literal        : [0-9]+
  decimal_literal    : [0-9]+ '.' [0-9]+
  string_literal     : "'" ... "'" | '"' ... '"'
  keyword            : 'break' | 'continue' | 'return' | 'match' | 'mix'

Delimiters:
  { } ( ) [ ] : = , ; . # =>
```
