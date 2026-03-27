
```js

limit : 99 
aoc01_01 #(input [String], Int) {

	times_in_0 : 0
	current_position : 50 
	
	
	input.each({
		move String
		direction, times : parse(move)
		operation: match direction { 
			Left => int.+
		}, {
			Right => Int.-
		}
		new_position : operation(current_position, times) % (limit + 1)
		new_position == 0 ? {
				times_in_0 = times_in_0 + 1
			}
		}
	})
	times_in_0
}

parse #(move String, Direction, Int) {
	
	match { 
		move.at(0) == "L" => Direction.Left()
	}, {
		move.at(0) == "R" => Direction.Right()
	}, {
		Direction.Invalid()
	}
	
	int.parse(move.substring(1)).or(0)
}
Direction : {
	Left(),
	Right(),
	Invalid()
}

```