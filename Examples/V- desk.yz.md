#example
[twitter](https://twitter.com/preslavrachev/status/1467555292242186241/photo/1)
```js
import os

main: {
  vertical_direction: 0
  horizontal_direction: 0

  lines  err : os.read_lines 'input.txt'
  err != null ? { panic '`err`' }

  lines forEach {
      line String
      parts: line.split(' ')
      direction: parts[0]
      magnitude: numbers parse_int parts at 1
      magnitude: numbers.parse_int(parts.at(1))

      when_eq direction  [
          {'up'}: { vertical_direction = vertical_direction - magnitude }
          {'down'}:{ vertical_direction = vertical_direction + magnitude }
          {'forward'}: { horizontal_direction = horizontal_direction + magnitude }
      ]

      direction == 'up' ? {
        vertical_direction = vertical_direction - magnitude
      }
      direction == 'down' ? {
        vertical_direction = vertical_direction + magnitude
      }
      direction == 'forward' ? {
        horizontal_direction = horizontal_direction + magnitude
      }

  }
  print('`vertical_direction * horizontal_direction`')
}

```
