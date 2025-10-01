#example

```js
bubble_sort: {
    values [] Int
    _: values.each({
        i Int
        _ Int
        j: 0
        _:while {j < values.len() - i - 1 } {
            values[j] > values[j + 1] ? {
                values[j] values[j + 1] = {values[j + 1] values[j]}()
            }
            j = j + 1
        }
    }
    values
}
main: {
    v: [25 13 8 1 9 22 50 2]
    v = bubble_sort v
    v.each({i Int
        print("{i}")
    }
}
```
