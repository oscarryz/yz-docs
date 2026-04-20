Because functions are singletons, we can really have a library without each caller overriding each other

```javascript
if 1 == 2 {}{}
if true {}{}
```
The internal would change from false to true. Same goes for longer lived things like `while`

```javascript
while {cond()} {}
while {cond2()} {}
```
The second call would override the  first parameter

#idea
Create and instance each time


```javascript
control: {
   While {
      apply: {
         cond {Bool}
         action {v}
         cond() ? { 
            action()
            apply cond action
         }
      }
    while: {
       cond {Bool}
       action { v }
       While().apply(cond action)
    }
}
```
It's a little bit ugly that it requires an `apply` all the time.

Maybe it can be anonymous?
```javascript
control: {
   while : {
      cond {Bool}
      action {v}
      i: {
       
         c : cond
         a : action
         self
         cond() ? { 
            action()
            self(cond action)
         }
      }
      i.self = i
      i
   }
}
```

alternate syntax for blocks (testing)
```javascript
control: {
   while : {
      cond :: Bool
      action :: v
      i: {
       
         c : cond
         a : action
         self
         cond() ? { 
            action()
            self(cond action)
         }
      }
      i.self = i
      i
   }
}
```

## Possible answer

Every function creates a world and is isolated from other worlds (a world is a coroutine/thread)

They communicate by returning assigning values, but common functions (like `while` ) are separate copies 

## Possible answer

just use the existing mechanism of creating either a class or an anonymous function. Take advantage of the type signatures and have the local create its own variables


```js
while { 
   cond { Bool }
   action { value }
} = {
   {
      cond { Bool }
	  action { value }
      c: cond
      a: action
      c() ? {
         action()
         while( c a )
      } 
   }()
}
```

 Point example 

```js
Point {
    x Int
    y Int
}
// same as 
point: {
    x Int 
    y Int
    { x: x y: y }
}
p: point( 1 2 ) 
q: Point( 1 2 )

```


#answered While not very practical the solution is to design your functions to return copies and/or new instances so they don't step on each other.