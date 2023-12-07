From: https://cuetorials.com/introduction/

```js

albums: [{
    artist: 'Led Zeppelin'
    album: 'BBC Sessions'
    date: '1997-11-11'
}]
```

```js
hierarchy: {
    Schema: {
        hello String
        life Int
        pi Float
        nums [] Int
        struct {}
    }
    Constrained: Schema {

        '[cue-constraint:=~"[a-z]+"]'
        hello String
        
        '[cue-constraint:>0]'
        life Int
        
        '[cue-constraint: list.MaxItems(10)]'
        nums []Int
    }
    value: Constrained {
        hello: 'world'
        life: 42
        pi: 3.14
        nums: [1 2 3 4 5]
        struct: {
            a: 'a'
            b: 'b'
        }
    }
}
```
```js
t: a
b: a
a=t
```

```js
Server: {
    'cue:[(cue.val)= ">5000 & <10_000"]'
    port: 1 
}

```