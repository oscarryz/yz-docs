#example
https://vosca.dev/p/fb2c261f00

```js
Node: {
    value Decimal
    left: Option.None()
    right: Option.None()
}
main: {
    left: Node(0.2)
    right: Node(0.3, Option.None(), Node(0.4))
    tree: Node(0.5, left, right)
    print(sum(tree))
}
sum: {
    tree Node
    match {
        tree == Option.None() => 0
    }, {
        tree.value + sum(tree.left) + sum(tree.right)
    }
}

```
