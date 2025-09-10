#example
```js
TreeNode: {
    val Int
    left TreeNode
    right TreeNode
    
}
max_depth_dfs: { root TreeNode
    stack: LinkedList<TreeNode>{}
    depths: LinkedList<TreeNode>{}
    current: root

    max_depth: 0
    current_depth: 0
    
    while {stack.is_empty().or({current != null})} {
        current == null ? {
            current_depth > max_depth ? {
                max_depth = current_depth
            }
            current = stack.pop().right
            current_depth = depths.pop() + 1 
        } {
            stack.push(current)
            current = current.left
            depths.push(current_depth)
            current_depth = current_depth + 1
        }
    }
    max_depth
}
```