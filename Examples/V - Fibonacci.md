#example
https://vosca.dev/p/9cd071f5fc

```js
fib: {
    n Int
    f: [0,1]
    2.to(n).each{ i Int
        f[i] = f[i - 1] + f[i - 2]
    }
    f[n]
}
main:{
    31.times({
        print(fib(i))
    }
}
`output
0
1
1
2
3
5
8
13
21
34
55
89
144
233
377
610
987
1597
2584
4181
6765
10946
17711
28657
46368
75025
121393
196418
317811
514229
`
```