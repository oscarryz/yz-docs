#example

https://arturo-lang.io/playground?example=First-class%20functions

 ```js

`
!:[Example]
example: {
  fcf()
  output: [
  "sin/asin => 0.5"
  "cos/acos => 0.4999999999999999"
  "cube/croot => 0.5"
  ]
}  
`
fcf: {
    cube:  { x Int; x ^ 3 }
    croot: { x Int; x ^ (1 / 3) }

    names: ["sin/asin", 'cos/acos', 'cube/croot']

    func_list: [math.sin, math.cos, cube]
    inv_list: [math.asin, math.acos, croot]

    num: 0.5

    0.to(func_list.len()).each({
        i Int
        result: func_list[i](num)
        print('${names[i]} => ${inv_list[i](result)}')
    })
}
```