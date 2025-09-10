#example


[Data structures](https://www.qdbplang.org/docs/examples#:~:text=str%20Print.-,Data%20Structures,-All%20of%20the)

```js
Stack {

  data : { None() }
  push : {
    element
    curr_data : data()
    Stack(
        data: {
            Some({
              val: {element}
              next: {curr_data}
            })
         }
    )

  }
  peek: {
    data()
  }

}

Sack().push(3).push(2).peek().print()
```
