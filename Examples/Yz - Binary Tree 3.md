#example
https://benchmarksgame-team.pages.debian.net/benchmarksgame/description/binarytrees.html#binarytrees

Implementation based on: [V](https://github.com/hanabi1224/Programming-Language-Benchmarks/blob/main/bench/algorithm/binarytrees/1.v), [Pony](https://github.com/hanabi1224/Programming-Language-Benchmarks/blob/main/bench/algorithm/binarytrees/1.pony) and [Go](hhttps://github.com/hanabi1224/Programming-Language-Benchmarks/blob/main/bench/algorithm/binarytrees/1.go)

```js

// SumTypes version
node: {
  Tree: {
      Empty()
      Node(
        left Tree = Empty()
        right Tree = Empty()
      )
      // Check returns an Int
      check #(Int) = {
          // If the constructor was Empty returns 0
          1 + .Empty ? { 
            0 
          } , { 
           // else runs left and right checks in parallel
           left.check(),
           right.check(),
          // adds the results
           left.0 + right.0
        }
      }
  }

  create #(depth Int, Tree) = {
    depth : 0 
    depth == 0 ? { 
      Empty() 
    }, {
      Node( create(n - 1), create(n - 1) )  
    }
  }
}
main: {
    // min, max, and stretch depths
    n: int.parse_int(os.args 1).or{ 4 }
    min_depth: 4
    max_depth: min_depth + 2 > n ? { min_depth + 2 }, { n }
    stretch_depth: max_depth + 1

    stretch_tree: node.create(stretch_depth)
    _: print('stretch tree of depth `stretch_depth`\t'++
             ' check: `stretch_tree.check()`')

    long_lived: node.create(max_depth)
    
    depth: min_depth
    _: while {depth <= max_depth} {
        iterations: 1 << max_depth - depth + min_depth
        // sum = ... makes it wait until the iteration finishes (sync) so `depth = depth + 2` runs
        sum: 0
        sum = iterations.times({
            sum = sum + node.create(depth).check()
        }
        depth = depth + 2
        print('`iterations`\t trees of depth `depth`\t check: `sum`')
    }
    _: print('long lived tree of depth `max_depth`\t check `long_lived.check())`'
}
```
