Using [Neptune impl](https://www.reddit.com/r/ProgrammingLanguages/s/ZJ4D36oZ2J)

```js
solution: input
  .lines()
  .map {
    line String
    digits: line
        .split()
        .filter(string.is_digit?)
   
    l: digits.len()
    l < 1 ? { 0 } 
    {
       f : digits[0]
       l : digits[l - 1]
       int.parse(f ++ l).or { 0 }
    }
  }
  .sum()
```