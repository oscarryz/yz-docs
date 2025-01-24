
Instead of current: 

```js
say_hi #(message String) = {
  message String
  print("I have a message for you: `message`")
}
```

We could omit the `=` at time of declaration and avoid redeclaring the parameter


```js
say_hi #(message String) {
  print("I have a message for you: `message`")
}
```

But the following would have to be legal too

```js
n Int 1
```

So `varname` + `type` + `literal`  is valid across the board

Other way could be using `=>`

```js
hi #(msg String) => {
  print("Msg: `msg`")
}

n Int => 1
n Int = 1
```

#answered  
Yes,  when the boc follows the signature immediately, the variables in the signature are automatically declared in the boc

```
greet #(msg String) {
  print(msg) // no need to declar `msg` in the boc
}

greet = { 
  msg String // to override still need to declare it
}
greet = {
  print(msg) // compilation error, there is no print
}
```