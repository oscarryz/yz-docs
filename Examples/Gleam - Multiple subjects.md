https://tour.gleam.run/flow-control/alternative-patterns/

```js
io: std.io
int: std.int
debug: print

main: {
	number: int.random(10)
	debug(number)

	result: when [
		{  [2 4 6 8].contains(number) } : {"This is an even number"}
		{  [1 3 5 6].contains(number) } : {"This is an odd number"}
		{ when.else } : {"I'm not sure"}
	]
	debug(result)
}

... 

main: {
	debug(get_first_non_empty([[]Int [1 2 3] [4 5]]))
	debug(get_first_non_empty([[1 2] [3 4 5] []]))
	debug(get_first_non_empty([[]Int []Int []Int]))
}
get_first_non_empty: {
	lists [][]T

	when [	
	   { list.len() == 0 } : { []T }
           { list[0].len() == 0 } : { get_first_non_empty(lists.shift()) }
	   { when.else } : { list[0] } 
	]
}

..
main: {
	numbers: [1 2 3 4 5]
	debug(get_first_larger(numbers 3))
	debug(get_first_larger(numbers 5))
}
get_first_larger (list []Int; limit Int; Int) = {
	when [
		{list.len() == 0} : {0}
		{list[0]>limit}: {list[0]}
		{when.else} : { get_first_larger(list.shift()) }
	]
}
```
