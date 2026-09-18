---
type: example
generated: { by: "oscarryz", at: 2023-12-09T03:59:33-06:00 }
---
#example

Using [Neptune impl](https://www.reddit.com/r/ProgrammingLanguages/s/ZJ4D36oZ2J)


```js

solution: input
  .lines()
  .map({
    line String
    digits: line
      .split()
      .filter(string.is_digit)
      .to_list()
    int.parse(
        digits.first().or({ '0' }) ++
        digits.last() .or({ '0' })
    )
  })
  .sum()
```