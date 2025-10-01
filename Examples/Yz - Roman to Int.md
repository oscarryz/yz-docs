#example
```js

roman_to_int: {
    roman String
    roman_t : [
        "I": 1
        "V": 5
        "X":10
        "L":50
        "C":100
        "D":500
        "M":1000
    ]
    num: 0
    pre : roman_t[roman.chat_at(0)]

    actual Int
    roman.each({ 
        cursor Int
        c String
        actual = roman_t[cursor]

        pre < actual ? {
          num = num - pre * 2 
        }
        num = num + actual
        pre = actual
    }
    num
}
in: input("Roman number: " )
print((roman_to_int(in)))

```