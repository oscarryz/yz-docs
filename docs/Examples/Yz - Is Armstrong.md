---
type: example
generated: { by: "oscarryz", at: 2023-12-06T20:13:02-06:00 }
---
#example
[Reddit](https://www.reddit.com/r/ProgrammingLanguages/comments/176it3o/showcase_your_lang_by_sharing_an_armstrong_number)

```js
is_armstrong #(Int, Bool) {
   n Int
   n == digits(n).map(3.*).sum()
}
digits #(Int, [Int]()) {
  n Int
  "${n}".split("").map({ c String; strings.parse_int(c) })
}
```




