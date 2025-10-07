#feature
Update: Oct 19 2024
Update: Oct 7 2025 Will use the keyworkd `mix`

#challenged [[Questions/import]]

It might be hard to assign multiple variables by hand. 

No import keyword, just assign to a new variable: 
#answered  No `import` keyword, to "import" all assign the top level block e.g. `math: core.math; math.max(1 2)`

```js
util: org.apache.something.util
util.compare 1 2 // 
```
~~Similar to Java, import will only help to avoid typing the full name. Assigning blocks from other places still will work just as any other variable assignment~~


Let's say we have a math package

```javascript
core: {
    math: {
        max: {
            a Int
            b Int
            a - b > 0 ? { a } {b}
        }
        min: {
            a Int 
            b Int
            a - b > 0 ? {b} {a}
        }
    }
}
// Used without import
r1: core.math.max 1 2
r2: core.math.min 1 2
// Could also be written as
mix core.math
//import core.math.max
r1: max(1,2)
//import core.math.min
r2: min(1, 2)

// or 
max: core.math.max
min: core.math.min
r1: max 1 2
r2: min 1 2
```

Is not possible however to import all the items at once :/, would need to use it as name space

```javascript
import core.math
math.max 1 2
math.min 1 2
```

But with `mix`  now we can
```js
mix core.math

```

So, there's no real benefit, because we can just use variable assignation

```javascript
math: core.math
r1: math.max 1 2
r2: math.min 1 2
```

Thinking .... what are the thread safe implications for this? What if two threads call the same function? What value is returned. 
A: use instantiable bocs

Possible answer: each thread creates its own world and the only way to sync with through [Channels](../../Questions/solved/concurrency/Channels.md) NOPE

#to-do Review with [Concurrency](Concurrency.md) All good. 

#open-question How to deal with concurrency and libraries? e.g
```javascript
foo:{
   x: max(1,2)
}
// somewhere else 
bar: {
   y: max(42,100)
}
```
If both are executed, is there a chance they get the wrong answer?  Or that one blocks the other? 
A: They will queue the request and response one at a time. Which obviously can create bottle necks but  the recommendation is to create new instances: 
```js
// private implementation
Max: {
      a Int
      b Int
      ? : { a > b ? { a }, { b } }
}
 max #(Int,Int,Int) = {
    a Int
    b Int
    Max(a, b).?() // creates a new object and returns 
}
```
