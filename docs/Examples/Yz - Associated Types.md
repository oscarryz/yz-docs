---
type: feature
generated: { by: "oscarryz", at: 2026-05-09T02:13:54+02:00 }
---
#feature 

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