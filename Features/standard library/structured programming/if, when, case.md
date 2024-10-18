
```js
std: {

  Bool: {
   if_true_then_else #(then #(V), else #(V))
  }
  True: {
    if_true_then_else #(then #(V), else #(V)) = {
      then()
    }
  }
  False: {
    if_true_then_else #(then #(V), else #(V)) = {
      else()
    }
  }
true: True()
false: False()
  ... 
  if #(condition Bool, then #(V), else #(V)) = {
    condition.if_true_then_else(then, else)
  }
  if_any #( condition_action_map [#(Bool)]#(V)) = {
    cam = condition_action_map
    cam.for_each {
      condition #(Bool)
      action    #(V)
      condition().if_true_then_else(action)
    }
  }
  
}
```