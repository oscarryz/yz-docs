#example

https://pyret.org/index.html

```js
BinTree: {
  value Int
  left: leaf
  right: leaf
}
leaf = BinTree()

`'
Calculate the sum of node values
Examples:
  >> tree_sum(leaf)
  0
  >> node4: BinTree{4}
  >> tree_sum(BinTree{5, node4})
  9
'`
welse: when.else
tree_sum: { t BinTree
  when_eq t [
    { leaf } : { 0 }
    {welse}: { t.value + tree_sum(t.left) + tree_sum(t.right) }
  ]
}
// also
tree_sum: {t BinTree
    t == leaf ? { 0 } {
        t.value + tree_sum(t.left) + tree_sum(t.right)
    }
}

// As "module"
tree_sum: {
    doc: 'Calculate the sum of node values'
    tests: {
        assert tree_sum(leaf) == 0
        node_4: BinTree{4}
        node_5: BinTree{5 node_4}
        assert tree_sum node_5  == 9
    }
    leaf: BinTree()
    BinTree: {
        value Int
        left: leaf
        right: leaf
    }
    f: {
        t BinTree
        t == leaf ? {
            0
        } {
            t.value + tree_sum(t.left) + tree_sum(t.right)
        }
    }
    f
}
```

```js
make_counter: {
    counter: 0
    {
        counter = counter + 1
        counter
    }
}
l1: make_counter()
l1() == 1
l1() == 2
l2: make_counter()
l2() == 1
l1() == 3
l2() == 2
```

```js
[1 2 3].map{ n * n }      //[ 1 4 6]
[1 2 3].filter { n >= 2 } //[2 3]
[4 5 6].fold { sum Int n Int
    sum + n
}
// 15
```
