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




```js
inc (Int, Int, Int )
inc { Int Int Int }
inc :: Int , Int, Int 

inc Int -> Int -> Int 

describable_function ( description String, (some_arg Number, Boolean))


first_element :: ( <Type>, arr []Type, Optional<T> ) 
first_element: {
	<Type>
	arr []Type
	arr.len() == 0 ? {
		optional.none()
	} {
		optional.some(arr[0])
	}
}
first_element( [1 2 3 4] )


combine:: (<TypeA>, <TypeB>, a []TypeA, b []TypeB, []TypeA)
combine: {
	<TypeA>
	<TypeB>
	a []TypeA
	b []TypeB
	a.concat(b)
}

```
