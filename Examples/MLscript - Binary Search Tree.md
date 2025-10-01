#example

https://ucs.mlscript.dev/ > Advanced Examples > Binary Search Tree


```js
Tree: {
  A
  Empty()
  Node(value A, left Tree(A), right Tree(A))

}

```

```js
// Binary Tree

type Tree = Node | Empty

// insert
// fromlist
// contains

class Tree {
  data
  left// = new Tree() // empty
  right// = new Tree() // empty

  insert(d) {
    if (this.data == undefined) {
      this.data = d
      this.left = new Tree()
      this.right = new Tree()
    } else {
      let branch = d < this.data ? this.left : this.right
      branch.insert(d)
    }
  }

  static fromList(list) {
    const root = new Tree()
    for (e of list ) {
      root.insert(e)
    }
    root

  }
  contains(element) {
    if ( this.data === element ) {
      true
    } else if (this.data == undefined) {
      false
    } else  {
      let branch = element < this.data ? this.left : this.right
      branch.contains(element)
    }
  }
/*
 root.contains("fin")
root/
  "hola"
  left/
    "adios"
    left
    right/
      "fin"
      left
      right
  right
  */

}


```


```go
package main

import "golang.org/x/exp/constraints"

type Tree[T constraints.Ordered] struct {
  data T
  left *Tree[T]
  right *Tree[T]
}
func (t *Tree[T]) insert(data T) {
  if t == nil {
    t = &Tree[T]{data}
    t.data  = data
    t.left =  &Tree[T]{}
    t.right = &Tree[T]{}
  }
  if data < t.data {
    t.left.insert(data)
  } else {
    t.right.insert(data)
  }
}
func (t *Tree[T]) contains(data T) {
  if t.data == data {
    true
  }
  if t == nil {
    false
  }
  if data < t.data {
    t.left.contains(data)
  } else {
    t.right.contains(data)
  }
}
```
