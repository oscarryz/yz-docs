#example

```js

arr: [1,2,3,4,5,6,7,8,9,10]
print(arr.select({ x Int; x % 2 == 0 })).to_str()


// Apply function to value
arr: [1, 2, 3, 4, 5]
print(arr.map(2.*).to_str()) // [2 4 6 8 10]

 ```

Explanation:
2 type is `Int`
`Int`'s `*` method is defined as:
```js
Int: {
  value T // 2
  * #(Int, Int) {
    other Int
    value * other
  }
}
```
Thus `2*` is the same as writing
```js
arr.map({other Int; 2 * other})
```

Array concatenation
```js
a = [1, 2, 3] ++ [4, 5, 6] // [1 2 3 4 5 6]
```

```js
// Arturo
print(select 1..10 => even?)
// Yz
print(1.to(10).select(even).to_str())
even #( n Int ) {
    n % 2 == 0
}
// or
even : {
	n Int
	n % 2 == 0
}
print(1.to(10).select({n Int; n % 2 == 0}).to_str())

Array: {
   select: {
       predicate #(Int,Bool)
       r: []Int

       n Int
       i: 0
       while { i < len() } {
          n: at(i)
          _: predicate n ? {
             r << n
          }
          i = i + 1
       }
       r
     }
   }
}

```

Ackerman
```js
ackermann: { m Int, n Int
  m == 0 ? {
    n + 1
  } {
    n == 0 ? {
      ackermann(m - 1, 1)
    } {
      ackermann(m - 1, ackermann(m, n - 1))
    }
  }
}
// or with match
ackermann : {
	m Int
	n Int
	match {
		 m == 0  =>  n + 1
		 n == 0  =>  ackermann( m - 1, 1)
		 default =>  ackermann(m - 1, ackermann(m, n - 1))
	}
}
ackermann: {
  m Int
  n Int
  match {
    m == 0 =>  n + 1
    n == 0 => ackermann(m - 1, 1)
    default => ackermann(m - 1, ackermann(m, n -1 ))
  }
}
```


```js
// 2.0
// With possible match pattern in the future
ackermann: {
  m Int, n Int
  match {
    m == 0 :  n + 1
    n == 0 : ackermann(m - 1, 1)
    _      : ackermann(m - 1, ackermann(m, n -1 ))
  }
}

```
