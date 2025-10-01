#example
https://vosca.dev/p/0e4be7f921

```js
Size: {
    width Int
    height Int
    area: { width * height }
}

Button: {
    area: Size.area // what happens with width and height? they need a definition
    // to make it work it would need to define these two?
    width: Size.width
    height: Size.height 
    title String
}
button: Button {
    title: 'Click me'
    height: 2
}
button.width = 3
assert button.area() == 6
print(button)
```
