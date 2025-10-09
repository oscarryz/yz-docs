#example

https://www.anardil.net/2018/binary-tree-in-haskell.html

```js

binary_tree:{

  "deriving:['Read', 'Equal']"
  Tree : {
    T
    Empty(),
    Node(data T,
      left Tree(T),
      right Tree(T)
    )
    self Tree(T)
    // Although T is generic`
    // is expected to have the methods
    // == #(T, Bool), < #(T, Bool), > #(T,Bool)
    // The code would show something like
    //   insert #(a T)
    //       where T is #(
    //                 == #(T, Bool),
    //                  < #(T, Bool),
    //                  > #(T,Bool)
    //                )
    // But due the lack of "where" the compiler might
    // add the signature at compile time (but not in the code)
    // If explicit signature is needed, an explicit type should be created
    // Ord #(
    //   == #(T, Bool),
    //   < #(T, Bool),
    //   > #(T,Bool)
    // )
    // insert #(a Ord)
   insert #(a T) = {
      a T
      self.Empty ? {
      self = Node(a, Empty(), Empty())
        return
      }
      a == t.data ? {
        return //Node(a, t.left, t.right)
      }
      a < t.data ? {
        t.left.insert(a)
      }, {
        t.right.insert(a)
      }
    }
    height #(Int) = {
      self.Empty ? { 0 },{
        max(1, 1 + left.height(), 1 + right.height())
      }
    }
  }
  "Helper function"
  show #(t Tree(T), depth Int, width Int, String) = {
    t.Empty ? {" "},{
      leftside: show(t.left, depth + width, width)
      rightside: show(t.right, depth + width, width)
      center: depth.times.join(" ") ++ show(t)
      leftside ++ center ++ rightside
    }
  }
  widest_element #(t Tree(T), Int) = {
    t.Empty ? { 0 }, {
      maximum(
        widest_element(t.left)
        widest_element(t.center)
        show(t.data).length
      )
    }
  }
  make #(a [T], Tree(T)) = {
    root: Empty()
    a.each({e T, roo.insert(e)})
    root
  }
  random #(n Int) = {
    list : [Int]()
    n.times().do({
      list.add(random.int(1, 99))
    })
    make(list)
  }
}
```


```js
WhildMushroom: {
  name String
  is_edible Bool
}
"derive:['Debug']"
InedibleError:{}

mushrooms: [
  WildMushroom(
     name: "King Bolete",
     is_edible: true,
  ),
  WildMushroom(
     name: "Oyser",
     is_edible: true,
  ),
  WildMushroom(
     name: "Amanita",
     is_edible: false,
  ),
]

total_price [Result]()(String,InedibleError) =
  mushooms.map({
    m Mushroom,
    m.is_edible ? {
      Ok(m.name)
    }, {
      Err(InedibleError)
    }
  })
```
