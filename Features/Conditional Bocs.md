#feature
#pattern-matching 

A similar version of pattern matching
A special type of boc that can be used to check on a condition of a data type.


## Evaluate conditions 

The conditional boc has 2 parts, the condition and the action which is a list of expressions or statements, separated by `=>`

```js
// bool returning string 
{ 1 < 2 => "One is lower than two"} 

// bool returning nothing
{ is_monday() =>  
  println("wake up")
  println("have breakfast")
  println("go to work")
}```

They are evaluated by the `match` which takes a list of conditional bocs separated by comma and evaluates them in order. 
The result of the action is returned by the `match` so it can be used as an expression: 

```js
score: 85
grade: match {
 score >= 90 => "A" 
}, {
  score >= 80 => "B" 
}, {
  "C"
}
println(grade)
```

If more than one value is needed the same rules apply as regular bocs

Example 
```js
factorial #(n Int) {
  match { 
    n == 0 => 1 
  }, { 
    n > 0 => factorial(n -1)
  }
}
```


If cascade evaluation is needed you can use continue

```js
n Int
match {
   n % 3 == 0 => print("Fizz")
   continue // evaluated the following condition
}, {
	n % 5 == 0 => print("Buzz")
}, {
	print("`n`")
}
```

The last element of the conditional block chain is the default value and doesn't need a condition 

## Match type variant

The other way to use the conditional bocs is to specify a type variant, to use this version the `match` needs a subject from whom the constructor used will be evaluated

```js
x Option(String) // could be either Some or None

// executing statemtnts 
match x {
	Option.Some(String) => println("The value is `x.value`")
}, {
	Option.None() => println("There was no value")
}

// or returning a value
value : match x {
	Option.Some(String) => x.value
}, {
   // Last option becuase Option can only be Some or None
   // so the clause can be ommited
   "No value"	
}
```