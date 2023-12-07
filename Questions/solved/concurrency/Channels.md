Any block when used within a 'nursery' becomes a channel 


```javascript
a_nursery <> {
    a Int = 0
    go fish a
    go create_fish a
}
fish: {
    a Int
    value: <- a
}
create_fish: {
    a Int
    a <- 42
}

```