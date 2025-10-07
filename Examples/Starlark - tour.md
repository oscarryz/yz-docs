
https://github.com/bazelbuild/starlark?tab=readme-ov-file#tour

```js
// Define a number
number : 10 

// Define a dictionary 
people: [
	"Alice": 22, 
	"Bob": 40, 
	"Charlie": 55,
	"Dave": 14
]

names : ", ".join(people.keys())  // Alice, Boc, Charlie, Dave

// Define a "function"
"Return a greeting"
greet: {
   names [String]
   "Hello `names`"
}

greeting : greet(names) 

above30 : people.filter({name String; age Int; age > 30})

println("`above30.len()` people are above 30.")

fizz_buzz : { 
	n String
	1.to(n).do({
		i Int
		s : ""
		println(match 
			{ i % 3 == 0 => s += "Fizz"; continue },
			{ i % 5 == 0 => s += "Buzz" },
			{ "`i`" })
		
	})
}
fizz_buzz(20)
```