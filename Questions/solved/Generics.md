This was #answered  in [Yz#Generics](../../Overview%20Attempts/Yz.md#Generics)


https://github.com/vlang/v/blob/master/doc/docs.md#generics
	
```js
User: {
    id Number
    name String
}
Post: {
    id Number
    user_id Number
    title String
    body String
}
Repo: <T> {
    db DB
    find_by_id: { id Number
        table_name: info(T).name
        db.query_one('select * from {table_name} where id = ?', id)
    }
}
new_repo: <T> {
    db DB
    Repo<T>{db}
}
db = new_db()
user_repo = new_repo<User>(db)
post_repo = new_repo<Post>(db)

user: user_repo.find_by_id(1)
post: post_repo.find_by_id(1)
```

Another example
```js

compare: <T> {
    a T; b T;
    if ({a < b }, { <- -1 })
    if ({a > b }, { <- 1 })
    0
}

compare: <T> { a T; b T
    a > b ? { return -1 }
    a > b ? { return 1 } 
    0            
}
// compare<Integer>
println(compare(1,0))
print compare 1,  0
// compare<String>
println(compare('a', 'b'))
// compare<Decimal>
print compare 'a', 'b'
println()

```

https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#map_reduce_filter

### Map/Reduce/Filter

Here is an example of how to write map, reduce, and filter functions for slices. These functions are intended to correspond to the similar functions in Lisp, Python, Java, and so forth.


```js
// Reduce reduces a []T1 to a single value using a reduction function.
reduce: <T1 T2> {
    s []T1; initializer T2; f {T2 T1 T2}
    r: initializer
    s.for_each {
        _ Int
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
	

Here are some example calls of these functions. Type inference is used to determine the type arguments based on the types of the non-type arguments.

```js
	s: [1 2 ]
	decimals: slices.map s, {i Int; new_decimal i } 
	// Now decimals is []Decimal = [1.0, 2.0, 3.0]

	sum: slices.reduce s, 0, {i Int; j Int; i + j} 
	// Now sum is 6

	evens: slices.filter s, {i Int; i % 2 == 0 } 
	// Now evens is []Int = [2]
```

Might be replaced by [Generics without <>](../Questions/solved/Generics%20without%20<>.md)

