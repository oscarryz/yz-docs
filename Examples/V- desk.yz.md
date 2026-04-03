#example
[twitter](https://twitter.com/preslavrachev/status/1467555292242186241/photo/1)
```js
main: {
  vertical_direction: 0
  horizontal_direction: 0

  lines: os.read_lines('input.txt')

  lines.each({ line String
      parts: line.split(' ')
      direction: parts[0]
      magnitude: numbers.parse_int(parts[1])

      match {
          direction == 'up'      => vertical_direction = vertical_direction - magnitude
      }, {
          direction == 'down'    => vertical_direction = vertical_direction + magnitude
      }, {
          direction == 'forward' => horizontal_direction = horizontal_direction + magnitude
      }
  })
  print('`vertical_direction * horizontal_direction`')
}

```
