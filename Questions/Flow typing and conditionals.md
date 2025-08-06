#pattern-matching 
#answered in  [Conditional Bocs](../Features/Conditional%20Bocs.md)

In the last example of  [SumTypes/Example](../Features/sumtypes/Example.md) we use `>_` to match with one of the types. 

Here is a more complex example that shows how could it interact with other features like  booleans ? , breaks, and assigments

The example comes from the [core-lang](https://github.com/core-lang/core/blob/eafdaba9ddf318a9c28c7101759038f164648465/bench/splunc/splunc.dora#L48-L107)repository

```js
splay #(tree Option(Node), node Node, Option(Node) ) = {
    tree match { 
        _ None 
        return tree
    }, {
        some Some
        tree = some.value

        node.left = None()
        node.rigth = None()
        l Node = node
        r Node = node

        while { true } , { 
            node.value <= tree.value ? { 
                tree.left match {
                    _ None(Node)
                    break
                }, {
                    some Some(Node)
                    tl : some.value
                    tree.left = tl.right
                    tl.right = Some(tree)
                    tree = tl
                    tree.left match { _ None; break }
                }
                // link right 
                r.left = Some(tree)
                r = tree
                tree = tree.left.or_panic()
            }
            node.value >  tree.value ? { 
                tree.right match {
                    _ None(Node)
                    break
                }, {
                    some Some(Node)
                    tr : some.value
                    tree.right = tr.left
                    tl.right = Some(tree)
                    tree = tr
                    tree.right match { _ None; break }
                }
                // link left 
                l.right = Some(tree)
                l = tree
                tree = tree.left.or_panic()
            }
        }
        l.right = tree.left
        r.left = tree.right
        tree.left = node.right
        tree.right = node.left
        Some(tree)
    }

}
```


Same program with the "traditional"  Option.is_some method

```js
splay #(tree Option(Node), node Node, Option(Node) ) = {
    tree.is_none() ? {
      return tree
    }
    tree = tree.or_panic()
    node.left = None()
    node.rigth = None()
    l Node = node
    r Node = node

    while { true },  { 
      node.value <= tree.value ? { 
        tree.left.is_none() ? { 
          break;
        }
        tl = tree.left.value
        tree.left = tl.right 
        tl.right = Some(tree)
        tree = tl 
        tree.left.is_none() ? {
          break
        }
      }
      node.value > tree.value ? { 
        tree.right.is_none() ? { 
          break;
        }
        tr = tree.right.value
        tree.right = tr.left 
        tr.left = Some(tree)
        tree = tr 
        tree.right.is_none() ? {
          break
        }
      }
      l.right = tree.left
      r.left = tree.right
      tree.left = node.right
      tree.right = node.left
      Some(tree)
   }
```

Yet the same example with a different formatting;  uses `>_` symbol instead of `match` ,fat arrow (`=>`) to separate the action from the test and aligns matches blocks with the symbol

```js
splay #(tree Option(Node), node Node, Option(Node) ) = {
  tree when 
    {  _ None  => return tree },
    { some Some => 
      tree = some.value
    
      node.left = None()
      node.rigth = None()
      l Node = node
      r Node = node
    
      while { true } , { 
        node.value when
          { <= tree.value => 
            tree.left when 
              { None(Node) => break },
              { some Some(Node) =>
                tl : some.value
                tree.left = tl.right
                tl.right = Some(tree)
                tree = tl
                tree.left >_ { _ None; break }
              }
              // link right 
              r.left = Some(tree)
              r = tree
              tree = tree.left.or_panic()
          }, 
          { > tree.value => 
          tree.right when
             { None(Node) =>  break },
             { some Some(Node) => 
               tr : some.value
               tree.right = tr.left
               tl.right = Some(tree)
               tree = tr
               tree.right >_ { _ None; break }
              }
              // link left 
              l.right = Some(tree)
              l = tree
              tree = tree.left.or_panic()
            }
          }
          l.right = tree.left
          r.left = tree.right
          tree.left = node.right
          tree.right = node.left
          Some(tree)
    }
}
```



Applying this syntax to another example from [DeltaScript](https://jbunke.github.io/blog/when)

```js
c Color
pfx String = "The color is "

c when 
  { c.alpha == 0 => print("`pfx` transparent") },
  { Color("#00000") => print("`pfx` black") },
  { Color("#fffff") => print("`pfx` whilte" },
  { c.r == c.g && { c.r == c.b} && {opaque(c)} => 
      print("`pfx` a shade of grey") } 
  { #ff0000, #00ff00, #0000ff => print("`pfx` an RBG primary color") }
  { else => printf("`pfx` not a match") }


```


Another example from [Core-lang](https://core-lang.dev/tutorial#toc-lang-12)
```js
  describe_float #(float Decimal, String) = { 
    float when
    { Decimal.NaN => "not a number"},
    { float.is_finite() == false => "infinity"},
    { float >= 0.0 => "positive"}, 
    { "negative"}
  }
  describe_float(123.4)
  describe_float(1.0/0.0)
```


### Decision

Flow typing will use the `when` keyword which will have a list of bocs with the match and the action separated by `=>` . 

When matching variants, only the name is needed for the match and the original variable can use the inner value directly (no need to deconstruct) as now is safe.

```js
Option: {
  Some(value T)
  None()
}
opt : maybe_something()
opt when 
  {Some => print("The value is `opt.value`")},
  {None => print("No value")}
```


When matching conditions, a boolean expression will be matched  ( #to-do  define if more than one)  It might not even need the subject, just go straight into `when` (akin to [Golang switch with no condition](https://go.dev/tour/flowcontrol/11) )

```js
color Color = some_color()
description: when 
  {color.alpha == 0 => "tranparent"},
  {color.hex == "#ffffff" => "white"},
  {color.hex == "#00000" => "black"},
  {color.r == color.g && {color.r == color.b} && {opaque(color)} =>
    "a shade of gray"
  },
  { color.hex.in(["#ff0000", "#00ff00", "#0000ff"]) => "an RGB primary color"},
  { color == Color.brigth_opaque => "bright"},
  {"not a match"}
 print("The color is `description`")
  
```
