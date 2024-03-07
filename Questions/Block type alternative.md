What if the block type could use `::` and `;`

eg

```javascript
// current
walkie_talkie: {
  thing {talk {String} walk {String}} 
  print thing.walk()
  print thing.talk()
}
// this suggestion
walkie_talkie: {
  thing :: talk :: String; walk :: String
  print thing.walk()
  print thing.talk()
}
```
Another option, similar to Haskell

```js

thing  (walk (String); talk (String))

thing (String; String) = {
	name String
	last_name String 
}
```


```js
// https://en.wikipedia.org/wiki/Higher-order_function

// twice :: (Int -> Int) -> (Int -> Int)
// twice :: f::(Int) ; (Int Int)
// twice { {Int} {Int; Int} } 
twice: {
  f { Int } 
  { 
      x Int 
      f(f(x))
  }
}
plus_three: { 
    i Int 
    i + 3
}
g: twice(plus_three)
println(g(7))


// new block signature 
twice :: (Int ; Int ) ; (Int; Int)
twice :: (Int ) ; (Int)
twice = { 
    f :: Int 
    { i Int; f ( f ( x ) ) }
}
plus_three :: Int 
plus_three = { 
    i Int
    i + 3
}
g :: Int 
g =  twice(plus_three)
g(7)



twice: { 
    f {Int;Int}
    {
        i Int
        f(f(i))
    }
}
g: twice(3.+)
```

