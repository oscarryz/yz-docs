How a block interface is defined



```js
add :: Int Int Int  = {
	a Int
	b Int
	a + b
}
add (Int, Int, Int) = {
	a Int
	b Int
	a + b
}
add ( a Int , b Int, Int ) => { 
	a + b 	
}
add { a Int; b Int ; Int} = {
	a Int
	b Int
	a + b
}
add : { 
	a Int 
	b Int
	a + b
}
person (name String, last_name String, to_string (String)) = {
	to_string = { 
		name + last_name 
	} 
}
person: {
	name String
	last_name String
	to_string: {
		name + last_name 
	}
}
person: {
	name String
	last_name String
	to_string (String) = {
		name + last_name 
	}
}
alice : person('Alice', 'Blanc') // alice is `(String)` because of `to_string`
bob : person
bob('Bob' 'Ito') 
bob.name // Bob
bob.last_name // Ito
// same goes for person
Person : {
	name String
	last_name String
	to_string (String) = {
		name + last_name 
	}
}
p: Person('Alice' 'Blanc' )

//...
Person ( name String , last_name String, to_string (String))
Person = person
alice : Person('Alice' 'Blanc')
bob : Person('Bob', 'Ito')


// Generics 
Node ( data T, left Node(T), right Node(T) ) 
// or 
Node : {
	data T
	left Node(T)
	right Node(T)
}

// Rule, single upper case letter is a generic type
a Node(T) = Node('a string') // T parameter bound to String
a.left = Node('left')
b.right = Node(1) // compilation error 

```