How is memory GC'ed?
How to avoid out of memory?
Can it be collected per "actor"
All the time I think the main actor / block will handle the method invocation to keep one writer, what about the `=` operator? Would that be also handled by the creating actor?
Could there be problems because there are more than one reference? What if the creator actor goes out of scope? 

```js
X {
   name String
}
a:{
   x: X() // x created here
}
b: a.x // b is a reference to the X created in line 2
a = null // not possible, but it means goes out of scope.

```


#answered The memory management won't be as efficient as it can be and the GC per actor is still far away. For v1 relying on go's GC should be enough. 
