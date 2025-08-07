A type can be defined as another type, because the `:` is an expressions so `A:{}` defines the empty type `A` and `B:A` defines the same type, thus making it an alias. 
You cannot add new operations on the new type, but you can override them on  creation:

e.g 

```js
// A and B are the same thing
B: A: {
  say_hi: { "Hi" }
}
// C structurally matches but redefines `say_hi`
C: {
  say_hi : A.say_hi
  say_hi = { "Bye" }
}
// Compilation error because we can't write to `say_hi` outside the block
// C.say_hi = "Bye" // C changes the internal implemention, migth be useful to simulate inheritance

print(A().say_hi()) // Hi
print(B().say_hi()) // Hi
print(C().say_hi()) // Bye


```


This could be also useful [Enums](Replaced%20features/Enums.md)
