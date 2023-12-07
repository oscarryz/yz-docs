	```javascript
PersonId {}
PetId {}

Pet {
  id PetId
  owner PersonId
}
a: PetId()
Pet ( a a ) // valid
```
Currently creating a pet would accept in both the same thing (and that's good)
Do we need a way to make them different to the type system can enforce it?

#answered That's a drawback from Yz it's heavily structure based and not nominal. We might do [what Typescript does](https://www.typescriptlang.org/play#example/nominal-typing) to a certain degree: 

```javascript
ValidatedString {
	value String
	__brand: 'User input post validation'
}
validate_user_input: {
	input String
	simple: input.replace('<', "<=")
	ValidatedString(simple)
}
print_name: {
	name ValidatedString
	print name
}
input: "alert('bobby tables')"
validatd_input: validate_user_input(input)
print_name validated_input
```

