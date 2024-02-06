We can define a data is generic by not specifying the data type and bind it on first usage, however we cannot enforce instances of the same type to be the same generic data

```javascript
Node {
	 data // generic
	 left Node //
	 right Node 
}
a = Node(1)
a.left = Node("B")
```

### Exploration

We can receive the generic type as argument to to block 

```js
b: {
	T // generic data type 
	data T // data is of tye T
}
// invocation
b(String 'hola')
// can be inferred? 
b('hola')

```

Could have the reserved word "type" and use it in similar to [Zig](https://ziglang.org/documentation/master/#toc-comptime)

```js
b: { 
	T type
	data T
}
b(String 'Hola')
b('hola')


HashMap { 
	K
	V
	put: { 
		key K
		value V
	}
}
map: HashMap{String Int}
map.put('64' 64)

list: { 
	T
	{ 
		items []T
		len Int
	}
}
list = List(String)

Node {
	T
	data T 
	left Node{T}
	right Node{T}
}
n: Node{String}
```


```js
none: std.option.none
tree: {
	Node {
	   T
	   data T
	   left Node = None(T) // Type Node { T; data T; left Node{T}; right Node}
	   right Node = None(T)
	}
	builder: {
		{
			T
			data T // generic data
			Node(T data)
		}
	}
}
node: tree.builder(Int)
root: node(1) // the builder `node` is now bound to `Int`
root.left = node(2)
root.right = node(3)
node_str: tree.builder(String)
a : node_str('hi')
a.left = node_str('Bye')
// this would be still valid however
a.right = node(1) // a.right of type Node{String} is not compatible with Node{Int}


```


```js
something: {
	T { len {Int} }
	data T
	data.len() // data is generic but needs a `len()` method
}
something(1) // comp err, Int doesn't have a `len()` method
something("Hi") // ok
something([1 2 3) // ok.. we wanted soemthing with `len()` after all.
```



```js
stack: {

	new: {
		Stack() 
	}
	Stack { 
		T 
		values []T // generic array
		push: {
			v T // generic value
			values << v
		}
		<< : { v T  push v }
		pop: {
			values.len() == 0 ? {
			   option.None(T)
			} { 
				option.Some(value.pop())
			}
		}
	}
}
main: {
	s: stack.new(Int)
	s.push 42
	s << 43 
	// This is an error, as the stack inner `values` is now bound to `Int` types
	stack.push 'Oh no!'

    // Trying to use a generic type before is bound would return an error
    t: stack.new(Int)
    x: t.pop() // not known yet, probably compile error
	10 + x.get() // yeah... probably not... or maybe, there's no way to know


    as_param: {
	    t stack.Stack{T}
		10 + t.pop().get() // would require t.value to be bound to Int // :? maybe
    }

}
```



```js
User {
    id Number
    name String
}
Post {
    id Number
    user_id Number
    title String
    body String
}
Repo {
	T
    db DB
    find_by_id: { id Number
        table_name: info(T).name
        db.query_one('select * from {table_name} where id = ?', id)
    }
}
new_repo: {
	T
    db DB
    Repo{T db}
}
db = some_library.new_db()
user_repo = new_repo(User db)
post_repo = new_repo(Post db)

user: user_repo.find_by_id(1)
post: post_repo.find_by_id(1)
```

```js
compare: {
	T
    a T
    b T
    if ({a < b }, { <- -1 })
    if ({a > b }, { <- 1 })
    0
}

compare: {
	T
	a T
	 b T
    a > b ? { return -1 }
    a > b ? { return 1 } 
    0            
}
// compare<Integer>
println(compare(Int 1 0))
print compare Int 1  0
// compare<String>
println(compare(String 'a' 'b'))
// compare<Decimal>
print compare Decimal 'a' 'b'
println()

```



```js
// Reduce reduces a []T1 to a single value using a reduction function.
reduce: {
	T1
	T2
    s []T1
    initializer T2 
    f {T2 T1 T2}
	
    r: initializer
    s.for_each {
        index Int
        v T1
        r = f r v 
    } 
    r
}

// Filter filters values from a slice using a filter function.
// It returns a new slice with only the elements of s
// for which f returned true.
filter: <T> { s []T; f {T Bool}
    r []T
    r.for_each {_ Int; v T1
        f v ? { r << v }
    } 
    r
}
```