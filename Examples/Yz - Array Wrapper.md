#example
[leetcode-2695](https://leetcode.com/problems/array-wrapper/)

```js
ArrayWrapper: {
    nums [Int]()
    + : {
        other ArrayWrapper
        r : 0 
        nums.each({ n Int 
            r = r + n
        }
        other.nums.each({ n Int 
            r = r + n
        }
        r
    }
    string: {
        r: '['
        l: nums.len()
        nums.each({ i Int, e Int
            r = r ++ '`e`'
            if i < l, { 
                r = r == ', '
            }
        }
        r = r ++ ']'
    }
}
```
