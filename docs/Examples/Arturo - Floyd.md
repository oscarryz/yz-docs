---
type: example
generated: { by: "oscarryz", at: 2023-12-06T20:13:02-06:00 }
---
#example

https://arturo-lang.io/playground/?8RJD5F


```js
width: {
    rows Int
    col Int
    floor(2 + log(col.+(1, (row * row - 1 / 2)), 10))
}
floyd: {
    rows Int
    n: 1
    row: 1
    col: 0
    while({ row <= rows }, {

        print(pad("${n}", width(rows, col)))
        col = col + 1
        n = n + 1
        col == row ? {
          print("")
            col = 0
            row = row + 1
        }
    })
}
floyd(5)
print("")
floyd(14)
```