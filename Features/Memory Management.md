```js
greet: {
    msg: 'Hello'
}
Person: {
    what String
    say_hi: { print what }
}


a: Person {what: greet.msg }
a.say_hi() // 'Hello'

greeet.msg: 'Bye'
b: Person {what: greet.msg }

a.say_hi() // 'Hello' still referencing to the original 'Hello'
b.say_hi() // 'Bye'   referencing to the new 'Bye'
```


To avoid copying put in on an array or a box

```js
Box: <T>{
    value T
}
greet
```