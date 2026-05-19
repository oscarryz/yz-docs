#example

From: https://cuetorials.com/introduction/

```js

albums: [{
    artist: 'Led Zeppelin',
    album: 'BBC Sessions',
    date: '1997-11-11'
}]
```

Basically the same as in Cuetorials

```js
hierarchy: {

    Schema #(
        hello String,
        life Int,
        pi Decimal,
        nums [Int],
        struct #()   
    )
    `
    !:[Cue]
    cue-on: [Schema]
    `
    Constrained: {

        `cue:=~"[a-z]+"`
        hello String

        `cue:>0`
        life Int

        `cue: list.MaxItems(11)`
        nums [Int]
        struct #()
    }
    value: Constrained(
        hello: 'world'
        life: 42
        pi: 3.14
        nums: [1, 2, 3, 4, 5]
        struct: {
            a: 'a'
            b: 'b'
        }
    )
}
```

More straightforward, compacting schema, constraint and data

```js
// hierarchy.yz
`!:[Cue]`
Schema :{
	`cue: '=~ "[a-z]+"'`
	hello: "world"
	
	`cue: ">0"`
	life: 42
	
	pi: 3.14
	
	`cue: "list.MaxItems(10)"`
	nums: [1,2,3,4,5]
	struct: {
		a: "a"
		b: "b"
	}
}
```


```js
`!:[Cue]`
Server: {
    `cue:">5000 & <10_000"`
    port: 1
}
```

See also: [Compile Time Bocs](docs/Features/Compile%20Time%20Bocs.md)
