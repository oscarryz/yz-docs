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
Manually create and instance each time


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
      w: {
       
         cond #(Bool)
         action #(v)
         cond() ? { 
            action()
            w: (cond, action)
         }
      }
      w()
   }
}
```

#related [Consider method execution to create a copy](Consider%20method%20execution%20to%20create%20a%20copy.md)
