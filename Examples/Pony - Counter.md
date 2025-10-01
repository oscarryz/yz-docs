#example

[counter/main.pony](https://github.com/ponylang/ponyc/blob/54225e41bd141b8b7a64f0e65d3dbe46bc5317d4/examples/counter/main.pony)

```js
Counter: {
    count Int = 0
    increment: {
        count = count + 1
    }
    get_and_reset: {
        display #(Int)
        display(count)
        count = 0
    }
}
display: {
        result Int
        print("`result`")
}
main:{
    n: int.parse(os.args[1]).or(Ok(10)).get()
    counter: Counter()

    _: n.times(counter.increment)

    counter.get_and_reset(display)
}

```
