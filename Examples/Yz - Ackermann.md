#example
[Ackermann function](https://en.wikipedia.org/wiki/Ackermann_function)

```js
ackermann: { m Int n Int
    when [
        {m == 0}: {n + 1}
        {n == 0}: { ackermann m - 1 1 }
        {when.else}:{ackermann m - 1 ackermann m n - 1 }
    ]
}
```

```kotlin
fun amw(m: Int, n: Int): Int = when {  
    m == 0 -> n + 1  
    n == 0 -> ackermann(n - 1, n)  
    else -> ackermann(m - 1, ackermann(m, n - 1))  
}
```

```js
ackermann: { m Int; n Int
  m == 0 ? { 
    n + 1
  } {
    n == 0 ? {
      ackermann m - 1  1
    } {
      ackermann m - 1 ackermann m n - 1
    }
  }
}
```

Is it the same as the following?

```js
ackermann: { m Int; n Int
  m == 0 ? { 
    n + 1
  }
  n == 0 ? {
    ackermann m - 1  1
  } {
    ackermann m - 1 ackermann m n - 1
   }
}
```

```js
`
function ack(m, n) {  
 m === 0 ? n + 1 : ack(m - 1, n === 0  ? 1 : ack(m, n - 1));  
}
`
ack: {m Int; n Int
    m == 0 ? { n + 1 } {ack m - 1 n == 0 ? {1}{ack m n - 1}}     
}


ack: {m Int; n Int
    m == 0 ? {
     n + 1 
    } {
      ack m - 1 n == 0 ? {
        1
      }{
         ack m n - 1
      }
    }     
}

0 .to(3).each({ m Int 
    0 .to(4).each({ n Int 
        print('ackermann `m` `n` => `ackermann m n`')
    }

}
```
