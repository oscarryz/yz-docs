#answered 
#accepted 
How a block interface is defined


The block signature is a list of variables / types / generics separated by `;`  enclosed in parenthesis: 

```javascript
// Add is a block that takes two `Int` and returns an `Int`
add ( Int ; Int ; Int ) 
add = { 
	a Int
	b Int
	a + b
}
// Can also have names and initialize right away.
// When the initalization follows the signature variable redefinition is needed.
add ( a Int ; b Int ; Int) = {
	a Int
	b Int
	a + b 
}
// Can also be inferred
//  ( Int ; Int ; Int )
add : {
	a Int 
	b Int 
	a + b
}
// if a type has a single capital letter, then in generic
node  ( data T ; left T; right T )  

// if the name capital then the block is instantiable
Node  ( data Option(T) ; left Option(Node(T)); right Option(Node(T)))

root : Node( Some(0) Some(Node(-1)) Some(Node(1)))


node  ( data T ; left T; right T )  = {

}
```


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
Node (data T, left Node(T), right Node(T) ) 
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


[Wikipedia - Type Signature](https://en.wikipedia.org/wiki/Type_signature)


```js
function_name ( T; U; var1 Type1 ; var2 Type2; outvar OutType);
function_name ( T, U, var1 Type1 , var2 Type2, outvar OutType)
is_even (n Int, Bool)

```


```js
factorial: {
	n Int
	n.gt(0).ifelse({
		n.times(factorial(n.minus(1)))
	}{
		n
	})
}
factorial: {
	n Int
	n > 0 ? {
		n * factorial(n - 1 )
	}{
		n
	} 
}

factorial :: ( n Int ; Int ) = { 
	n > 0 ? {
		n * factorial (n - 1)
	} {
		n
	}
}

```




```js
json: std.json
payload String =   `
        "vals": {
            "testing": 1,
            "production": 42
        },
        "uptime": 9999
        `
Config ( Vals(testing Int, production Int) , uptime Int) = {
   Vals: {
       testing Int
       production Int
   }
   uptime Int
}
config (stream json.TokenStream, res json.JSON) = {
   stream: json.TokenStream(payload)
   res = json.parse(Config stream {})
}
main (Result) = {
    if config.vals.production > 50 {
         error "only up to 50 supported"
    }
    print "up={config.uptime}"
}
```



```javascript
import std.collections.map

'derive: [Show]'
Group: {
  group_size Int
  greoup_name String
}

double_group #( Group; Group ) = { 
    g Group
    Group ( 2 * g.size, g.name )
 }

names [String]
name = ['Sheldon' 'Leonard' 'Penny' 'Rajesh' 'Howard']

groups #([]String)
groups = {map new_group names ++ map double_group groups}

nth #(Int, []Group, String)
nth = {
  n Int
  gs []Groups
  when [
    {gs.is_empty}: {strings.error}
    {n < gs[n].size}: {gs[n].name}
    {when.else}:{ nth (n - gs[n].size  gs.rest())
  ]
}

```
