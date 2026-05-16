#example
[FizzBuzz.st] https://gist.github.com/oscarryz/d54dd569fea585ec008c6f20af2e97ec

Simple version using conditionals:
```js
1.to(100).each({ i Int
    s: ""
    i % 3 == 0 ? { s = s + "Fizz" }, { }
    i % 5 == 0 ? { s = s + "Buzz" }, { }
    (i % 3 != 0 && i % 5 != 0) ? { s = i.to_string() }, { }
    print(s)
})
```

Idiomatic version using match:
```js
1.to(100).each({ i Int
    m3: i % 3 == 0
    m5: i % 5 == 0
    s: match {
        (m3 && m5) => "FizzBuzz"
    }, {
        m3 => "Fizz"
    }, {
        m5 => "Buzz"
    }, {
        i.to_string()
    }
    print(s)
})
```

Best version using a boc declaration:
```js
fizz_buzz #(i Int, String) {
    m3: i % 3 == 0
    m5: i % 5 == 0
    match {
        (m3 && m5) => "FizzBuzz"
    }, {
        m3 => "Fizz"
    }, {
        m5 => "Buzz"
    }, {
        i.to_string()
    }
}

main: {
    1.to(100).each({ i Int; print(fizz_buzz(i)) })
}
```
