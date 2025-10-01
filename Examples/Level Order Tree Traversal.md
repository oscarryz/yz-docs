#example

Description

https://www.geeksforgeeks.org/level-order-tree-traversal/

Go

```go
package main

import (
        "fmt"
)

type Node[T any] struct {
        data  T
        left  *Node[T]
        right *Node[T]
}

func node[T any](data T, l *Node[T], r *Node[T]) *Node[T] {
        return &Node[T]{data, l, r}
}
func leaf[T any](data T) *Node[T] {
        node(data, nil, nil)
}

func main() {

        root := node("hello",
                node("germanic",
                        leaf("hello"),
                        leaf("hallo"),
                ),
                node("latin",
                        leaf("hola"),
                        leaf("ola"),
                ),
        )
        printLevel(root)
}

func printLevel[T any](n *Node[T]) {
        q := []*Node[T]{}
        if n != nil {
                q = append(q, n)
        }
        for len(q) > 0 {
                c := q[0]
                q = q[1:]
                if c.left != nil {
                        q = append(q, c.left)
                }
                if c.right != nil {
                        q = append(q, c.right)
                }
                fmt.Printf("%v\n", c.data)
        }
}
```


Java

```java
import java.util.ArrayList;

record Node<T>(
  T data,
  Node<T> left,
  Node<T> right
){}

void main() {
        var root = new Node<String>(
                "Hello",
                new Node<String>(
                        "Hallo", null, null ) ,

                new Node<String>(
                        "Hola",
                        new Node<String>("Adios", null, null) ,
                        null)
                );
        printLevel(root);
}

<T> void printLevel(Node<T> node) {
        if (node == null ) {
                return;
        }
        var q = new ArrayList<Node<T>>();
        q.add(node);
        while ( ! q.isEmpty()  ) {
                var c = q.get(0);
                q.remove(0);
                System.out.println(c.data());
                if (c.left() != null ) {
                        q.add(c.left());
                }
                if (c.right() != null ) {
                        q.add(c.right());
                }
        }
}
```


Rust

```rust
struct Node<T> {
    value T,
    left Option<Box<Node<T>>>,
    right Option<Box<Node<T>>>,
}

fn main() {
    let mut root = Node<String> {
        "Hola",
        Some(Box::new(Node<String> { "Adios", None, None } ) ),
        None
    };
    print_level( &root ) ;
}


fn level_print<T: Debug>(n : &Node<T>) {
    let mut q: Vec<&Node<T>> = vec![n];

    while !q.is_empty() {
        let c = q.remove(0);

        println!("{:#?", c.value);

        if let Some(l) = c.left.as_ref() {
            q.push(l)
        }

        if let Some(r) = c.right.as_ref() {
            q.push(r)
        }
    }
}
```

Yz

```js
// Node is a generic data structure
// using "named constructors" (similar to sum types)
Node: {
    T  // T is the generic type
    // `None` is the constructor with no arguments
    None()
    // `Tree` is the construtor with data, left and right defaulting to None
    Tree(
        data T,
        left Node(T) = None(),
        right Node(T) = None()
    )
}

main: {
    // Creating an instance with the `Tree` constructor
    root : Tree( 0
            Tree( -2
                Tree( -3 ),
                Tree( -1 )
            ),
            // `left` is ommited here, thus, named attributes are needed
            Tree( data:  2,
                right = Tree(3)
            )
        )
    print_level(root)
}

print_level #(Node(T)) = {
    r Node(T)
    q : [q]

    // for control structure
    for q.is_empty != false {

        c : q.remove(0)
        println("`c.data`")

        // if match control structure
        if c.left match Node.Tree {
            q.add(c.left)
        }
        if c.right match Node.Tree {
            q.add(c.right)
        }
    }
}
```
