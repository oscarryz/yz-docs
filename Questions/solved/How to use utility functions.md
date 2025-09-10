#answered Use them, and try to make them short lived so others can run.
#todo  Document this further 


Preferred is to create new types and instances of those types, but if using utility functions shouldn't be a problem because their usage will be sequential 

Long lived executions like `while` will be either yielding all the time, or pushing themselves into the execution queue on recursive calls. 

```js
while #(cond#(Bool), action#()) {
	match {
		cond() => action() // calling `cond` and `action` might yield the thread
		while(cond, action) // this would put then next `while` call in the queue
	}, {
	  // nothing
	}
}
```

See also: https://tutorial.ponylang.io/gotchas/scheduling#long-running-behaviors


--- 
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

#related [Consider each method invocation would create a boc copy](Consider%20each%20method%20invocation%20would%20create%20a%20boc%20copy.md)
