#example
https://benchmarksgame-team.pages.debian.net/benchmarksgame/description/binarytrees.html#binarytrees

Implementation based on: [V](https://github.com/hanabi1224/Programming-Language-Benchmarks/blob/main/bench/algorithm/binarytrees/1.v), [Pony](https://github.com/hanabi1224/Programming-Language-Benchmarks/blob/main/bench/algorithm/binarytrees/1.pony) and [Go](https://github.com/hanabi1224/Programming-Language-Benchmarks/blob/main/bench/algorithm/binarytrees/1.go)

```js


node: {
    empty: Node {
        check: {0}
    }
    Node:{
        left  : empty
        right : empty 
        check: {
            check_left:  { check(left) } 
            check_right: { check(right) } 
            // execute in parallel 
            check_left()
            check_right()
            // when done use the result
            result: 1 + check_left.result + check_right.result
        }
    }

    create: {
        depth: 0
        n == 0 ? {
            Node(empty, empty)
        } {
            Node(create(n - 1), create(n - 1))
        }
        
    }
}
main: {
    // min, max, and stretch depths
    n: int.parse_int(os.args[1]).or_else({ 4 })
    min_depth: 4
    max_depth: min_depth + 2 > n ? { min_depth + 2 } { n }
    stretch_depth: max_depth + 1

    stretch_tree: node.create stretch_depth
    print("stretch tree of depth `stretch_depth`\t check: `stretch_tree.check()`")

    long_lived: node.create max_depth
    
    depth: min_depth
    while({ depth <= max_depth }, {
        iterations: 1 << max_depth - depth + min_depth
        sum: 0
        // sum = ... makes it wait until the iteration finishes (sync) so `depth = depth + 2` runs
        sum = iterations.times({
            sum = sum + node.create(depth).check()
        })
        depth = depth + 2
        print("`iterations`\t trees of depth `depth`\t check: `sum`")
    })
    print("long lived tree of depth `max_depth`\t check `long_lived.check()`")
}
```
