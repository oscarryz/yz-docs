#feature

In Yz, sum types (discriminated unions) are expressed as **type variants**: a boc that lists named constructors, each carrying its own data. A value of the type holds exactly one constructor at a time, and `match` dispatches on which one it is.

```yz
Option: {
    V
    Some(value V)
    None()
}

x: Some("hello")
match x
    { Some => print(x.value) },
    { None => print("nothing") }
```

For the full design — syntax, constructor rules, exhaustiveness, and more examples — see [Type variants](Type%20variants.md).