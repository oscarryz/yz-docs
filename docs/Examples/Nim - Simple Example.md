https://nim-lang.org/

```js
strformat: std.strformat

// Using the Cue macro to add validation
`!:[Cue]`
Person: {
	name String
	`cue: ">0"`
	age Int
}
people: [
	Person(name:"John", age: 45),
	Person(name:"Kate", age: 30),
]

people.each({person Person; 
 // Type-safe string interpolation, 
 // evaluated at compile time.
 print("${person.name} is ${person.age} years old")
})

// The odd_numbers doesn't work like this
// but I think that could be alternatives... 
odd_numbers #(T Mod, a [T], buffer[T]) {
	a.each({
		x % 2 == 1 ? {
			buffer.add(x) // yield is not possible in Yz
		}
	})
}
buffer [Int]()
odd_numbers(Int, [3,6,9,12,15,18], buffer)
buffer.each({ odd Int
	print("`odd`")
})

...
// Using Yz macro system to transform a dense
// data-centric description of x86 instructions
// into a lookup table (array of strigs?) that 
// can be used in the code :/

// Macro that transforms 
// a `;` string into an array :/ 
ToLookUpTable : {
    Schema : #(data String)
    perform #(input Boc, Boc) {
	   data: input.annotation.data
	   result: [String])()
	   data.split(";").each({item String; result.add(item)})
	   boc.assign(result) 
    }
}

// Annotation to make use of the macro
`!:[ToLookUpTable]
 data: "mov;btc;cli;xor"`
opcodes [String] // initialized now. 

```