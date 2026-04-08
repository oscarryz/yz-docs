# return, break, and continue

## return

By default, a boc returns the value of its last expression. Use `return` to exit early:

```yz
check: {
  age Int
  age < 21 ? {
    return "You have to be over 21 to access this site"
  }
  "Welcome to the site"
}

print(check(20))   // You have to be over 21 to access this site
print(check(25))   // Welcome to the site
```

`return` exits the boc where it is lexically defined. Inside a nested boc (e.g. a callback), `return` exits that inner boc, not the outer one.

## break

`break` exits the innermost enclosing loop:

```yz
max_from_list: {
  list [Int]
  m: 0
  list.each { item Int
    item < 0 ? {
      break    // stops iteration immediately
    }
    m = max(m, item)
  }
  m
}
```

## continue

`continue` skips the rest of the current loop iteration and moves to the next:

```yz
list.each { n Int
  n % 2 == 0 ? {
    continue   // skip even numbers
  }
  print(n)     // prints only odd numbers
}
```

## In match expressions

`continue` inside a `match` arm evaluates the next arm (fall-through):

```yz
n Int
match {
  n % 3 == 0 => { print("Fizz"); continue }
}, {
  n % 5 == 0 => print("Buzz")
}, {
  print("`n`")
}
```
