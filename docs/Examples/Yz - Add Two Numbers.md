#example
https://leetcode.com/problems/add-two-numbers/

```js

Node: {
    val: 0
    next: Option.None()
    is_next: { next.Some }
}
solution: {
    l1 Node
    l2 Node
    carry: 0
    result: Node(0)
    current: result

    while({ l1.is_next() || l2.is_next() || carry > 0 }, {
        sum: l1.val + l2.val + carry
        carry = sum % 10
        current.next = Option.Some(Node(sum / 10))
        current = current.next
        l1 = l1.next
        l2 = l2.next
    })
    result.next
}

// aux method to create a list from an array
list_from: {
    array [Int]()
    r: Node()
    c: r
    array.each({ item Int
        c.next = Option.Some(Node(item))
        c = c.next
    })
    r.next
}
l1: list_from([2, 4, 3])
l2: list_from([5, 6, 4])
r: solution(l1, l2) // Node{7 Node{0 Node{8}}}
```
