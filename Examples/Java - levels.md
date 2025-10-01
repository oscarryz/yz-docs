#example

From https://stackoverflow.com/a/73874369/20654


```js
Node: {
  id String
  level Int
  name String
  list []Node = []
}

input: [
  Node{"1" "1" "Item1"}
  Node{"2" "2" "Item2"}
  Node{"3" "2" "Item3"}
  Node{"4" "3" "Item4"}
  Node{"5" "2" "Item5"}
  Node{"6" "1" "Item6"}
  Node{"7" "2" "Item7"}
  Node{"8" "3" "Item8"}
  Node{"9" "4" "Item9"}
]
map [Int][]Node
map[1]=[]Node
input.each: { node Node
  list: map[node.level]
  list << node
  map[node.level + 1] = node.list
}
println('{map[1]}')

```
