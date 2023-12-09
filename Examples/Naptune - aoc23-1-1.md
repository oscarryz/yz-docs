Using [Neptune impl](https://www.reddit.com/r/ProgrammingLanguages/s/ZJ4D36oZ2J)

```js
input.lines().map({
  line String
  digits: line
   .split()
   .filter(is_digit?)
   
  l: digits.len()
  l < 1 ? {
      0
  } {
    int.parse(
      digits[0] ++ digits[l - q]
    ).or { 0 }
}).sum()
```