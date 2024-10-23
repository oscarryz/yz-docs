
Instead of current: 

```js
say_hi #(message String) = {
  message String
  print("I have a message for you: `message`")
}
```

We could omit the `=` at time of declaration and avoid redeclaring the parameter


```js
say_hi #(message String){
  print("I have a message for you: `message`")
}
```
