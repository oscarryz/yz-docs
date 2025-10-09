#example
```js
hanoi: { 
    s [Int]()
    h [Int]()
    d [Int]()
    ha: {
        n Int
        s [Int]()
        h [Int]()
        d [Int]()
        n > 0 ? {
            ha n - 1 s d h
            move s d
            ha n - 1 h s d
        }
    }
    ha s.len() s h d
}
move: { 
    from [Int]()
    to [Int]()
    to << from.pop()      
}
hanoi [3 2 1] [] []

```
