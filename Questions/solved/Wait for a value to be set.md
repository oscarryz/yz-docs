We usually have this way to wait until a block execution finishes: 

```js
// The following line blocks until fetch returns
thing: fetch('resource')

```
For non blocking we have 

```js

// This line doesn't block
fetch('resource')
// `thing` here might have garbage or the real value
thing: fetch.value 
```
How can we make waiting for a value to be ready like in Go channels `<-`

```
thing <- fetch.value
```

Possible answer is to rely on the block finishing... 

```js
fetch('resource')
```


The almost complete solution is to have a "double function callback ish?" (not sure what to call it).
On the current design (11/02/2023) we can pass a callback (similar to a channel) and we already have a mechanism to wait until a function completes

```js

long_function: {
	callback (String) // callback is a block that receives/returns a String
	r: random.int(100) // a random number between 1 - 100 
	s: 'Complex compuation: $(r)'
	callback(s) // sends the value `s` to the callback 
}
somewhere_else: {
    value: callback.wait_for_s_to_be_ready() // this would solve our problems...	
}
// And potentially this would work, it would just go on an infinite loop
// (until) the value is set
callback: {
	message String = strings.undefined 
	wait_for_message_to_be_ready: {
	    message == strings.undefined ? { 
		    wait_for_message_to_be_ready()
	    } { 
		    message
	    }
	}
}
```

A better implementation (or ideal) would block the method call until the value is ready, probably using channels themselves?  (if implemented in go) or lock/unlock

```js
callback: {
	message String
	wait_for_message: {
	   lock()
	}
}
```

Another source of inspiration is [Awaitability](https://github.com/awaitility/awaitility/blob/master/awaitility/src/main/java/org/awaitility/core/ConditionAwaiter.java#L96-L122) implementation which is indeed a for loop submitting executors.

```java
while (!timeout) {
    executor.submit(callback())
    if success -> break
}
return value
```

Another semi ideal is to put it on a library like `wait_for(reference)` and create a listener on it

```javascript
nac: {message String}
compute: {
	callback {String}
	m: 'Result is $(random.int(100))'
    callback(m)
}
somewhere_else: {
   callback {message String}
   result: wait_for(callback.message)
}
...
wait_for: { 
	reference {}
	notify_here: {return}
    observe(reference, notify_here)
    // or something like that
}
```


```js
c {String; get {String}}: {
	msg String
	set: { msg String }
	get: { 
		msg != string.undefined ? { msg } { get() }
	}
}
busy: {
   c {set {String}}
   c.set('Long copmputation $(int.random(100))')
}
consume: {
  c { get{String}}
  value: c.get()
}
main: {
   busy(c)
   consume()
}

```

It seems the answer is [to use "Promises"](https://www.reddit.com/r/ProgrammingLanguages/comments/17omnvi/is_there_a_feature_where_a_variable_can_be/)


It seems the default mechanism Yz has is already enough, sure, it won't look like channels but probably that's even better. 

The example above could be as follows: 

```javascript
long_function: {
	callback (String) // callback is a block that receives/returns a String
	r: random.int(100) // a random number between 1 - 100 
	s: 'Complex compuation: $(r)'
	callback(s) // sends the value `s` to the callback 
}
callback: (message String)
somewhere_else: {
	wait: long_function( callback )
	msg: callback.message
    print msg
}
main:{
	somewhere_else()
	print 'Hola' 
}

```
So, `a` calls `b` and waits for it to be completed. 

Another example is the producer consumer in Synchronous Concurrencyn


[Pony](https://rosettacode.org/wiki/Synchronous_concurrency#Pony)  and [Go](https://rosettacode.org/wiki/Synchronous_concurrency#Go)

```js
printer: {
   count: 0
   {
	   line String
	   print(line)
	   count = count + 1
   }
}
main: {
   p: printer()
   f: file.open('input.txt').or_panic()
   _: f.each_line { 
	   line String
       p(line)
   }
   print 'Number of lines read: $(p.count)'    
}
```

Yz has asynchronous calls, so the key is to force something to wait to be done by assigning it to a value, in the example above `_`   



#answered 
Keep using the *"assign to a variable to make it finish"*  procedure. If can be infinite use a callback. 


