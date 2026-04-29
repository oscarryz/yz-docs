#feature

# Infostrings

An infostring is a string literal placed **before** a boc or field declaration. It
attaches passive metadata to the declaration — describing it, annotating it, or
providing parameters to `Compile` implementations that read it.

Infostrings describe. They do not execute.

See also: [Compile-Time Bocs](./yz-compile-time-bocs.md) ·
[Structural Reflection](./yz-structural-reflection.md)

---

## Syntax

A string literal immediately before a declaration is its infostring:

```yz
"json:user_name"
name String
```

```yz
"deprecated: use new_name instead"
old_name String
```

Both `'` and `"` delimiters are valid. The infostring must appear on the line
immediately before the declaration it describes.

---

## Boc-Level Infostrings

An infostring before a boc declaration attaches to the whole boc:

```yz
"Represents a registered user in the system"
Person : {
    name String
    email String
}
```

---

## Field-Level Infostrings

An infostring before a field attaches to that field only:

```yz
Person : {
    "json:first_name"
    name String

    "json:ignore"
    internal_id String
}
```

---

## Multiple Infostrings

A declaration can have only one infostring. If multiple annotations are needed they are
combined in a single string, with interpretation left to the `Compile` implementation
that reads them:

```yz
"json:dob validate:required"
dob Date
```

---

## How Infostrings Are Read

Infostrings are stored in the `Boc` instance the compiler creates for every
declaration. They are accessible as `Boc.infostrings` — a `[String]` — inside any
`Compile` implementation:

```yz
JSON : {
    run #(Boc, Boc) = { parent Boc
        parent.fields.forEach({ f Boc
            f.infostrings.forEach({ s String
                s.starts_with("json:") ? {
                    // use the json field name
                    json_name : s.split(":").last()
                }, {}
            })
        })
    }
}
```

The infostring format is a plain string. There is no enforced structure — each `Compile`
implementation defines and documents the format it expects.

---

## What Infostrings Are Not

Infostrings are not code. They do not execute. They have no effect on the program unless
a `Compile` implementation reads and acts on them.

They are not a documentation system — though a `Compile` implementation could use them
to generate documentation.

They are not a test framework — though a `Compile` implementation could use them to
generate tests.

They are parameters to `Compile` implementations, nothing more and nothing less.