---
type: example
generated: { by: "oscarryz", at: 2023-12-06T20:13:02-06:00 }
---
#example
https://leetcode.com/problems/two-sum/solutions/

```js
two_sum: {
    array [Int]()
    target Int
    map: [[Int]: Int]()
    array.each({ i Int; n Int
       map.contains(target - n) ? {
         [i, map.get(target - n)]
       }, {
         map.put(target - n, i)
       }
    })
    []
}
two_sum: {
  nums [Int]()
  target Int

  result: [Int]()
  hash: [[Int]: Int]()
  nums.each({ i Int; n Int
    hash.contains(target - n) ? {
      result.push(i)
      result.push(hash.get(target - n))
      result
    }, {
      hash.put(n, i)
    }
  })
  result
}

```
