A type can be defined as another type, because the `:` is an expressions so `A:{}` defines the empty type `A` and `B:A` defines the same type, thus making it an alias. 
You cannot add new operations on the new type, but you can override them on  creation:

e.g 

```js
C: B: A: {
  say_hi: { "Hi" }
}
// A, B and C are the same thing
C.say_hi = "Bye" // C changes the internal implemention, migth be useful to simulate inheritance

print(A().say_hi()) // Hi
print(B().say_hi()) // Hi
print(C().say_hi()) // Bye


```


This could be also useful [Enums](Enums.md)
