---
type: example
generated: { by: "oscarryz", at: 2023-12-06T20:13:02-06:00 }
---
#example
```js
// This is a single comment
/*
   This is a multiline comment
*/

`
   Counter is a block that has a variable count of type Int,
   an increment method that increases count by 1,
   and a value method that returns the current count.
`
Counter: {
    count: 0
    increment: {
        count = count + 1
    }
    value: { count }
}

main: {
    c: Counter()
    print("c.value = ${c.value()}")
    c.increment()
    c.increment()
    print("after two increments: ${c.value()}")
}
```