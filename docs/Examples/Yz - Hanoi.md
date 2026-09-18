---
type: example
generated: { by: "oscarryz", at: 2023-12-06T20:13:02-06:00 }
---
#example
```js
hanoi: {
    s [Int]()
    h [Int]()
    d [Int]()
    ha: {
        n Int
        s [Int]()
        h [Int]()
        d [Int]()
        n > 0 ? {
            ha(n - 1, s, d, h)
            move(s, d)
            ha(n - 1, h, s, d)
        }
    }
    ha(s.length(), s, h, d)
}
move: {
    from [Int]()
    to [Int]()
    to << from.pop()
}
hanoi([3, 2, 1], [], [])

```