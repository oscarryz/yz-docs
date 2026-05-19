#example


[Data structures](https://www.qdbplang.org/docs/examples#:~:text=str%20Print.-,Data%20Structures,-All%20of%20the)

```js
Stack: {
  T
  // this doesn't work because of the Self data 
  // type which doesn't exists
  data #(Some(#(val T, next Self)){ None() }
  push: {
    element T
    curr_data: data()
    Stack(
        data: {
            Some({
              val: { element }
              next: { curr_data }
            })
         }
    )

  }
  peek: {
    data()
  }

}

Stack(Int).push(3).push(2).peek().print()
```