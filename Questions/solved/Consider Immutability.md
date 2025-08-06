#dropped  All things will be mutable; with [Concurrency](../../Features/Concurrency.md)  each variable would have a single writer, thus, this is not strictly speaking needed
....

Probably in Yz can use values instead of variables so they are always final. 
This might complicate things for a language the is meant to be simple. 

Or could be immutable ala Elixir, where the value is immutable but the references can still change. 

```javascript
Point: {
  x Int
  y Int
}
p1: Point{1 2}
p2: p1
p1.x = 3
print '{p1}' // Point{3 2}
print '{p2}' // Point{1 2} Only the first one was modified and a new point was created
             // The other reference `p2` is unaffected
             // if there's another mutation, it will ressign to a new variable

p1.y = 5     // p1: Point{3 5}
```


This could work, but I wonder if that's a feature Yz needs. 

Pros: 

- Immutability ftw!
- Things will be easier to debug? 

Cons: 

- What happens with shared data? Will always be changing? 
    - yes, but you can keep a reference to the outer name.. 

```javascript
// Examples 
arr: [ 1 2 3 4 ]
same_array: arr
print '{arr}'        // [1 2 3 4]
print '{same_array}' // [1 2 3 4]
arr << 5 
print '{arr}'        // [1 2 3 4 5]
print '{same_array}' // [1 2 3 4]   <- same array didn't change

```

**answer**: Yes, it will change, if the original is needed, don't reference it directly, reference the enclosing instead.

```javascript
counter: {

    count: 0 // re-bind when `increase() is called`
    increase: {
        count.++() // modifies `Int.self` to 1
    }
    print_count: {
        print 'Count: {count}'
    }
}
counter.increase()
print_count() // Count: 1
Int: {
    self Int
    ++: {
        self = self + 1
    }

}
```
Other example:
```javascript
enclosing :{
   // Examples 
   arr: [ 1 2 3 4 ]
}

same_array: enclosing.arr
while { enclosing.arr.len() < 10 }, {
   // don't use same_array, modify enclosing.arr directly
   enclosing.arr << rand.int()
}


```
#answered
See [Immutability](Immutability.md)


