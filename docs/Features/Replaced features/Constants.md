---
type: replaced-feature
generated: { by: "oscarryz", at: 2023-12-06T20:13:02-06:00 }
---
#replaced-feature
#feature
There's no constants, but the convention is to use all upper case. 

The compiler doesn't enforce it, so just don't change it 

#open-question  Should the compiler enforce immutable with all upper case? #answered No

```js
math: {
    PI: 3.141516
}
...
Circle: {
    r Real
    area: {
        math.PI * r ** 2
    }
}
```
