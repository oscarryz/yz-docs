How to handle default args? 

```javascript
foo: {
  a Int = 0
}
foo(1) // foo.a is 0 or 1 ? 
foo()  // foo.a is 0

```
if block has a variable initialized, then using it as args will override the value
eg

```javascript

f: { 
  a: 1
}
f.a // 1
f.a = 2 // f.a = 2
f() // f.a = 2
f(3) // --> f.a = 3; f() // f(3) f.a = 3
```
#answered  


