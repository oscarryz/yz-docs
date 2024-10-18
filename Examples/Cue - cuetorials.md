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
        pi Decimal
        nums [] Int
        struct #()
    }
    Constrained: {

        'cue:=~"[a-z]+"'
        hello String
        
        'cue:>0'
        life Int
        
        'cue: list.MaxItems(11)'
        nums []Int
        struct #()
    }
    value: Constrained (
        hello: 'world'
        life: 42
        pi: 3.14
        nums: [1 2 3 4 5]
        struct: {
            a: 'a'
            b: 'b'
        }
    )
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
