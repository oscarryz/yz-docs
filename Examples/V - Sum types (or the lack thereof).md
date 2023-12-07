https://vosca.dev/p/fb2c261f00

```javascript
empty: Node {}
Node: {
    value Decimal
    left empty
    right empty
}
main: {
    left: Node {0.2}
    right: Node {0.3 empty Node {0.4}}
    tree: Node{0.5 left right}
    print sum tree
}
sum: {
    tree Node
    when_eq tree [
        {empty}: {0}
        {when.else}: {tree.value + sum(tree.left) + sum(tree.right)}
    ]
}

```