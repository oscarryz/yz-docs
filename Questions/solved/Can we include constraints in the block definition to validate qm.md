#answered  Keep the constraints with sumtypes or doc strings

Idea for another version perhaps. 

In [CUE lang types are values](https://cuelang.org/docs/tour/basics/types-are-values/) and thus, you can specify a value as the type of a field

For instance: 

```json
// CUE
#non_empty : !~ "^$"
non_empty_string : "" & #non_empty
```

Which allows to define an schema where a value cannot be empty

```json
// CUE
// The book schema
book: {
  title: string & !~ "^$"
}
// A validation that fails:
hp: book & {
  title: ""
}
//hp.title: invalid value "" (out of bound !~"^$"):
```

Can we do something like that? 

```js
Book #(
  title String !~ "^$"
)
hp: Book("") // would result in compilation error?
```

We currently have the idea to use the info string for that
```js
Book #(
  'cue:!~"^$'
  title String
)
hp: Book("") // should be a compilation error
```

But relies on putting a validation in the info string.

#idea We could expand the values literals to include the validation e.g 

```js
  title: '!~"^$"'
  id:  >0 
// title is a string but the string represents a regexp !~^$ non empty
// id is a number and the number is >0 greater than 0

// Or have values as types
title != "^$" // type string != ""
```

No: it clashes with types start with Uppercase

Another idea is to expand the type definition to include constraints

```js
id Int.(> 0 && < 100) // type number gt 0 and lt 100 
title String.(!= "")
weekend DayOfWeek.(Saturday|Sunday)

```
No: starts getting weird

Probably the best solution is to stick to the info string and have a preprocessor / macro

```js
'constraints: > 0 && < 100'
PositiveInt : Int 

id PositiveInt = -1 // compilation error

'constraints: != "" && len() < 30'
NonBlankString : String
title NonBlankString = "" // compilation erors

```


This constraint can be used to defined exhaustive values maybe? 

```js
Mon:Tue:Wed:Thu:Fri:Sat:Sun:DayOfWeek:{foo:"bar"}

'constraints: Sat | Sun'
WeekEnd: DayOfWeek

ss DayOfWeek = Sat // Ok
day DayOfWeek = Mon // compilation error

```

Status: Probably would be better to just stick to runtime validation and augment the type system to make better API but remain in the current status of `Some|None`  and abstract classes as sum types

```js
//
DayOfWeek : {
	Mon(),
	Tue(),
	Wed(),
	Thu(),
	Fri(),
	Sat(),
	Sun()
	
	is_weekend #(Bool) {
		match {
			Sat() => true
		}, {
			Sun() => true
		}, {
			false
		}
	}
}
```


```js
// Using `Option(NonBlank)` to enforce non-blankness
NonBlank: String 
non_blank #(String, Option(NonBlank)) = { 
  s String
  match {
	  s.is_empty() => None()
  }, {
	  Some(s)
  }
}
title : non_blank("")
title.and_then({ s String ;  print("`s` was a good book") }).or({panic()})
title >>= { s String ;  print("`s` was a good book") }.or({panic()})

```






