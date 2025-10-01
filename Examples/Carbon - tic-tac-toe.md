#example

https://github.com/shadowninja55/carbon/blob/master/examples/tic-tac-toe.cb

```js
	board: "".repeat(3).repeat(3)

	display_board: {
    	board.split('\n').for_each({
          \row String
      	print(row.join("|"))
		})
	}
	won_game: {
    player String
		win: player.repeat(3)
		[board, transpose(board)].each({
			board.each({ row String
				row == win ? { true }
			}
		}
		diag: { board Board
			map( {i Int
				board[i][i]
			})
		}
		diag(board) == win || { diag(board.map(reverse)) == win}
	}

	 arr: [1 2 3]
	 arr.for_each ({ i Int
		 println(i)
	 })

	 i: arr.lenght - 1
	 while({ i >= 0 }{
		 print(arr[i])
		 i = i - 1
	 })

	 meaning_of_life: 42
	 false ? {
		 print('false')
	 } {
		 meaning_of_life == 42 ? {
			 print('the meaning of life is 42')
		 }   {
			 print('mmmm')
		 }
	 }


	print_and_increment: {
		x: 0
		{
			print(x)
		   x = x - 1
		}
	}
	printer: print_and_increment()
	printer()
	printer()
	printer()


	Button: {
		make: { count Int
			times: when_eq count [
				{1}: {"once"}
				{2}: {"twice"}
				{when.else}: { "{count} times"}
			]
			msg: 'Click me `times`'
		}
	}
```


```js
//https://twitter.com/preslavrachev/status/1467555292242186241/photo/1
main: {
	vd: 0
	hd: 0
	lines: os.read_lines('input.txt').or {
		e Error
		panic ('`e`')
	}
	lines.each({
		line String
		p : line.split(' ')
		d : p[0]
		m : int.parse(p[1]).or { e Error panic('`e`}) }
		vd = vd + when_eq d [
			{'up'}: { m.neg() }
			{'down'}: { m }
			{when.else}: {0}
		]
		hd = hd + d == 'forward' ? { m } { 0 }
	}
	print("`vd * hd`")
}

	main: {
		vertical_direction: 0
		horizontal_direction: 0
		lines: os.read_lines('input.txt').or {
			err Error
			panic('{err}')
		}

		lines.for_each({
			line String
			parts: line.split(' ')
			direction: parts[0]
			magnitude: numbers.parse_int(parts[1]).or {panic()}
			when_eq direction [
				{'up'}: {vertical_direction = vertical_direction - magnitude}
				{'down'}: {vertical_direction = vertical_direction + magnitude}
				{'forward'}: {horizontal_direction = horizontal_direction + magnitude}
			]

		})
		print('`vertical_direction * horizontal_direction*`')
	}
```
