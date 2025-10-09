#example

https://tour.gleam.run/flow-control/alternative-patterns/

```js
io: std.io
int: std.int
debug: print

main: {
	number: int.random(10)
  	printint(number)

	result: match  {
	  [2 4 6 8].contains(number) => "This is an even number"
	  }, {
	    [1 3 5 6].contains(number) => "This is an odd number"
	  }, {
	  "I'm not sure"
	}
	
	print(result)
}
// using hamcrest match library
is: hamcrest.is
match: hamcrest.match
else: hamcrest.else

result: match(number).conditions([
  is.either([2 4 6 0]): {"This is an even number"}
  is.either([1 3 5 6]): {"This is an even numer" }
  else : { "I'm not sure" }
])

main: {
	print(get_first_non_empty([[Int]() [1 2 3] [4 5]]))
	print(get_first_non_empty([[1 2] [3 4 5] []]))
	print(get_first_non_empty([[Int]() [Int]() [Int]()]))
}
get_first_non_empty: {
	lists [[T]]

	when [
      { list[0].len() == 0 } : { list[0] }
      { list[0].len() > 0 } : { get_first_non_empty(lists.shift()) }
      { list.len() == 0 } : { list }
	]
}
// hamcrests
match(list).conditions[{
  is(empty(list[0])) : { list[0] }
  is.not(empty) : { get_first_non_emtpy(list.shift()) }
  is.empty : {list}

}]

//..
main: {
	numbers: [1 2 3 4 5]
	print(get_first_larger(numbers 3))
	print(get_first_larger(numbers 5))
}
get_first_larger #(list [Int], limit Int, Int) = {
	when [
		{list[0] > limit}  : {list[0]}
		{list[0] <= limit} : { get_first_larger(list.shift()) }
		{list.len() == 0} : {0}
	]
}
```
