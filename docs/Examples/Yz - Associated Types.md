```js
Graph #(
 Node #()
 Edge #()
 neighbors #(Node, [Edge])
)
print_neighbors #(g G Graph, n G.Node)


SocialGraph : {
   Node : Int
   Edge : String
   neighbors : {
      node Node
      ["The end"]
   }
}
...
so: SocialGraph()
print_neighbors(so, 1)
   
```