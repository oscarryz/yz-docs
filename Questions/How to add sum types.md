This is probably for some other version or some other language. 
Follow up exploration of [Type variants](../Features/Type%20variants.md)

Currently the approach to take will be use "abstract classes" as describe by munificient here: 
https://www.reddit.com/r/ProgrammingLanguages/comments/10jewgp/comment/j5m7jn9

But due the lack of sealed classes we can't do an exhaustive analysis of all the variants it might support. This would have to be runtime check with a panic / exception of unknow cases. 


Take for instance in [Boolean Blindness](https://shreevatsa.wordpress.com/2015/01/31/boolean-blindness/)


```haskell

data Num = Zero | Succ Num

plus :: Num -> Num -> Num
plus Zero     y = y
plus (Succ x) y = Succ (plus x y)

```

Haskell can be exhaustive because it can pattern match on `zero` or `Succ Num`

Similar can be achieved with Java (probably v21)

```java
sealed interface Num permits Zero, Succ {}
record Zero() implements Num { }
record Succ(Num pred) implements Num {}


Num plus(Num x, Num y) {
  return switch(x) {
    case Zero z ->  y;                        
    case Succ s-> new Succ(plus(s.pred(), y));
  };
}
```

The switch expression can match on `Zero` and `Succ`  and know that it is exhaustive because the `Num` interface is sealed and only allows those two variants. 
In Yz there's nothing like that. The approximation would be moving the `plus` as a method of `Num` and implementing in the "subclasses"

```js
// Yz
Num #(
   plus #(Num, Num)
)
Zero: {
  plus : {
    y Num
    y
   }
}
Succ: {
   pred Num
   sum: {
      y Num
      Succ(sum(pred, y))
   }
}
zero Zero = Zero()
one Succ = Succ(zero)
two: one.plus(one)
```

Which might be enough(ish)

Creating a variant using  [Pattern matching 1](Pattern%20matching%201.md)
and something like a hamcrest api would be 

```js 
//Yz
Num: {}
Zero: {}
Succ: {
  pred Num
}

sum: {x Num, y Num 
  match(x).type([
    {z Zero}:{ y }
    {s Succ}:{ Succ(plus(s.pred, y)) }
    { true }:{ panic() } // can't validate all the cases at compile time
  ])  
}
```
To do this version we would need to have a way to perform a type cast which we currently don't have. 
Options: 
1. `s Succ = (Succ) value` C version,  problem, `Succ` is a type and can be passed as param
2. `s Succ = value.(Succ)`  Go version, problem also as the first is not entirely type safe
3. `s Succ = Succ.(value)` Same semantic, but slightly more `Yz`er 
4. `s Succ = Succ(value)`  Wrong, because `Type()` is a constructor.
5. `info(value).instance_of(Succ, { s Succ, /* do something with s } }` I mean.. at what cost, and still would need to cast it
6. `switch value { Succ(pred) -> sum(pred,y),...}`  Built in pattern matching, but it needs a way to check for exhaustiveness (or to close the `Num` definition)
7. `s Succ = value` Can compile time check if value is either `Succ` or a super 

#answered  in [Type variants](../Features/Type%20variants.md)

The example above would be
```js
Num: {
  Succ(pred Num)
  Zero()
}
sum #(x Num, y Num) {
  x when 
  { Zero => y },
  { Succ => Succ(x.pred, y)},
}
```