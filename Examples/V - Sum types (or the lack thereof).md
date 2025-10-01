#example
https://vosca.dev/p/fb2c261f00

```js
Node: {
    value Decimal
    left Option(Node) = None()
    right Option(Node) = None()
}
main: {
    left: Node(0.2)
    right: Node(0.3, None(), Node(0.4)}
    tree: Node(0.5, left, right)
    print(sum(tree))
}
sum: {
    tree Node
    when_eq tree [
        {None}: {0}
        {when.else}: {tree.value + sum(tree.left) + sum(tree.right)}
    ]
}

```