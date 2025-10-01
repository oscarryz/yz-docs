#example


[qdbp Examples](https://www.qdbplang.org/docs/examples)

```js
// Functions
// With generics
add : {
	a P
	b P
	a + b // a must support the `+ {n}` method
}
add 1 2
// without generics
add: {
	a Int; b Int
	a + b
}
// With type
add #( a Int, b Int ) = {
	a + b
}
// idea... when adding a type we can define inplace using | or ->
//add {a Int b Int |  a + b }
//add {a Int b Int ->  a + b }
// not needed
add {a Int; b Int; a + b }

// Generics
print: {that that.print()}
print(3 // as long as `Int` has a `print` method)
print('Hello' // same for `String`)

// If/Then/Else
if: {
	cond Bool
	then #(v V) // generic
	else #(v V) // generic
	cond.if_true_if_false(then, else)
	// or
	cond ? then, else
}
// named args require use parenthesis
if (cond: 1 > 2
	then: {
		"true".print()
	}
	else: {
		"false".print()
	})
// no named args
if 1 > 2 {
	"true".print()
} {
	"false".print()
}

// kinda weird
Bool {
	val Bool
	if: {
		then {v}
		else {v}
		val ? then else
	}
}
true.if{
	"t".print()
} {
	"f".print()
}
// switch
switch: {
	v T
	this : {
		val : {
			v T
		}
		result: {
			None()
		}
		case: {
			v (T)
			then (T)
			val() == v ? {
				result = then()
				this.result = { Some(result) }
			} {
				this
			}
		}
		default: {
			then
			this.result().some?{ v T v }
			this.result().none(then())
		}
	}
	this
}
str: switch(5)
	.case(1 {"one"})
	.case(2 {"two"})
	.case(3 {"three"})
	.default({"none of the above"})

str: switch 5
	.case 1 {"one"}
	.case 2 {"two"}
	.case 3 {"three"}
	.default {"none of the above"}

// with a mapn in Yz
switch: {
    value
	conds [(Bool)](v V)
	default {v}

	conds.each({
		condition {Bool}
		then {v}
		value == condition() ? {
			v() // breaks the cycle
		} {
			continue; // tricky part
		}
	}
	// still here?
	default()
}

str: switch 5 [
	{ 1 }: {"one"}
	{ 2 }: {"two"}
	{ 3 }: {"three"}
]  { "none of the above"}
print(str)


// Data structures
stack: {
	data: {
		None()
	}
	push: {
		e
		current_data = this.data()
		{
			this.data = {Some({
				value: {e}
				next: {current_data}
			})}
		}
	}
	peek: {
		d: this.data()
		d.is_none({Err("empty stack")})
		d.is_some({v v})
	}
}
stack.push(3).push(2).peek().print


// Clases
circle: {
	radius Float
	{
		radius: {
			radius
		}
		print: {
			r.print()
		}
	}
}
my_circle: circle(radius:3)
my_circle.print()

basic_circle: {
	radius: {3}
	print: {radius().print()}
}
colored_circle: {
	radius: basic_circle.radius
	color: {"red"}
	print: {
		 color().print()
		 radius().print()
	}
}



```
