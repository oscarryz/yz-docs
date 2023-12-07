https://leetcode.com/problems/two-sum/solutions/

```javascript
two_sum: {
    array []Int
    target Int
    map: [Int]Int
    array.for_each { i Int n Int 
       map.contains v ? {
         return [i map.get n]
       } {
         map.put target - n i 
       }
    }
    []
}
two_sum: {
  nums []Int
  target Int

  result: []Int
  hash : [Int]Int
  nums.for_each { i Int; n Int 
    hash.contains target - n ? {
      result << i
      result << hash.get need
      return result
    } {
      hash.put n, i
    }
  }
  result
}

```
