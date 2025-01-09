
1. Same data types on all the variants
  1. No values
```js
Switch {
  On()
  Off()
}
```
  
  2. Same values
  ```js
FileAccess { 
    Read(v:1)
    Write(v:1)
    Execute(v:1)
  }
```

  2. Different values
```js
Planet {
    Mercury(mass: 123, radius: 456)
    Venus(mass: 9393, radius: 9881)
    ... 
  }
```
1. Different data types on all the variants
  2. All the variants differ
  ```js
Shape {
    Empty()
    Circle(radius Double)
    Rectangle(width Double, height Double)
    Point(x Int, y Int)
    Text(msg String, font_size Int)
  }
```
  4. Some of the variants differ
  ```js
Status { 
    Success(code Int),
    Warning(code Int),
    Error(code Int, msg String),
    Timeout(), 
    Pending(),
  }

classify_status #(status Status) { 
  status >_ {s Success; print("Sucessful response (2xx): `s.code`"),
            {w Warning; print("Redirection (3xx): `w.code`"},
            {e Error; print("Error ocurred: `e.mgs`(`e.code`)") },
            { println("Other status")}
}
```


