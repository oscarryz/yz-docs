#example

https://vosca.dev/p/669817e77a

```js
sqr:  {n Int; n * n}
cube: {n Int; n * n * n}
run: {
    value Int
    op #(n Int, Int)
    op(value)
}
main: {
    {
        print('Anonymous functions')
    }()

    print(run, 5, srq) // 25

    // Anounymous function can be declared inside other functions:
    double_fn: { n Int ; n + n }

    // Functions can be passed around without assigning them to variables
    res: run(5, {n Int; n + n})
    print('`res`')// 10

    // You can even have array/map of functions
    fns: [sqe cube]
    r: fns[0](10) // 100
    print("`r`")

    fns_map: [
        'sqr': sqr
        'cube': cube
    ]
    println(fns_map['cube'](2))
}

```
