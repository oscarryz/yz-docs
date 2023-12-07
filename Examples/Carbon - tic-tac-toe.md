https://github.com/shadowninja55/carbon/blob/master/examples/tic-tac-toe.cb

	board: repeat(3, repeat(3, " "))
	
	display_board: {
		board.split('\n').for_each({
			row
			print(join(" | ", row))
		})
	}
	won_game: { player Player
		win: repeat(3, player)
		[board, transpose(board)].for_each {
			board.for_each { row String
				row == win ? { <- true }
			}
		}
		diag: { board Board
			map( {i Int
				board[i][i]
			})
		}
		diag(board) == wing || 
	}

	 arr: [1,2,3]
	 arr.for_each ({ i Int
		 println(i)
	 })

	 i: arr.lenght - 1
	 while({ i >= 0 },{
		 print(arr[i])
		 i = i - 1
	 })

	 meaning_of_life: 42

	 false ? {
		 print('false')
	 }, {
		 meaning_of_life == 42 ? {
			 print('the meaning of life is 42')
		 } ,  {
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
			times: match(count, 
				{1}, {"once"},
				{2}, {"twice"},
				{count}, { "{count} times"}
			)
			msg: 'Click me {times}'
		}![[desk.yz]]
	}

https://twitter.com/preslavrachev/status/1467555292242186241/photo/1


	main: {
		vertical_direction: 0
		horizontal_direction: 0
		lines, err: os.read_lines('input.txt')
		err != null ? { panic('{err}') }

		lines.for_each({ line String
			parts: line.split(' ')
			direction: parts[0]
			magnitude: numbers.parse_int(parts[1])
			direction == 'up' ? {
				vertical_direction = vertical_direction - magnitude
			}
			direction == 'down' ? {
				horizontal+direction = horizontal_direction + magnitude
			}
			direction == 'forward' ? {
				horizontal_direction = horizontal_direction + magnitude
			}
		})
		print('{vertical_direction * horizontal_direction*}')
	}

