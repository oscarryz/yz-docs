#example

https://tour.gleam.run/flow-control/alternative-patterns/

```js
io: std.io
int: std.int
debug: print

main: {
	number: int.random(10)
  	printint(number)

	result: match {
	  [2, 4, 6, 8].contains(number) => "This is an even number"
	}, {
	  [1, 3, 5, 7].contains(number) => "This is an odd number"
	}, {
	  "I'm not sure"
	}

	print(result)
}

main: {
	print(get_first_non_empty([[Int](), [1, 2, 3], [4, 5]]))
	print(get_first_non_empty([[1, 2], [3, 4, 5], []]))
	print(get_first_non_empty([[Int](), [Int](), [Int]()]))
}
get_first_non_empty: {
	lists [[T]]

	match {
      list[0].length() == 0 => list[0]
    }, {
      list[0].length() > 0  => get_first_non_empty(lists.shift())
    }, {
      list.length() == 0    => list
	}
}

main: {
	numbers: [1, 2, 3, 4, 5]
	print(get_first_larger(numbers, 3))
	print(get_first_larger(numbers, 5))
}
get_first_larger #(list [Int], limit Int, Int) {
	match {
		list[0] > limit  => list[0]
	}, {
		list[0] <= limit => get_first_larger(list.shift(), limit)
	}, {
		list.length() == 0 => 0
	}
}
```