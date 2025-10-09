#example
```js
`Given an integer array `nums` and an integer `k`, return _the_ `kth` _largest element in the array_.

Note that it is the `kth` largest element in the sorted order, not the `kth` distinct element.

You must solve it in `O(n)` time complexity.
`

kth: { a [Int]()  k Int

    l: [Int]()
    a.each({ e Int 
        add k e l
    }
    l[0]
}
add: { k Int e Int l [Int]() 

    i: 0
    while { i < l.len() } {
        e < l[i] ? {
            break
        }
        i++
    }
    l << i e
    l.len() > k ? {
       l.remove 0
    }
}
```
