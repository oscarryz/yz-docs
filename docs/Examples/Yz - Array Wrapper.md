---
type: example
generated: { by: "oscarryz", at: 2023-12-06T20:13:02-06:00 }
---
#example
[leetcode-2695](https://leetcode.com/problems/array-wrapper/)

```js
ArrayWrapper: {
    nums [Int]()
    +: {
        other ArrayWrapper
        r: 0
        nums.each({ n Int
            r = r + n
        })
        other.nums.each({ n Int
            r = r + n
        })
        r
    }
    string: {
        r: '['
        l: nums.length()
        nums.each({ i Int; e Int
            r = r ++ '${e}'
            i < l ? {
                r = r ++ ', '
            }, { }
        })
        r = r ++ ']'
    }
}
```
