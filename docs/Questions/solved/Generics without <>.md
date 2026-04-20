
This was #answered 
Update: 03/12/2024: No, there's no easy way to [How to enforce data types in generics](How%20to%20enforce%20data%20types%20in%20generics.md)  so a type will be used. Single upper case letter, no <> and will use () for instantiation. 


Everything below is historic:

TLDR; Yes, variables without type are generic. They'll bind to the first usage. If a "tag" is needed maybe pass it as arg, but in general just don't pass it.

In qdbp, [generics](https://www.qdbplang.org/docs/examples#:~:text=as%20a%20function.-,Generics,-Methods%20are%20generic) don't specify the type and also infer them from the usage e.g. 

```smalltalk
"print is generic and needs an object that respond to the `Print` method"
print := {that | that Print.}
```

Can we do something similar in Yz? 

```javascript
print: {
    that // that is generic and should respond to the `print()` method
    that.print() 
}
```

The type inference would declare the `print`  block as 

```javascript

print { that { print {} } }
// Explained
print {           // `print` is a block 
    that {         // with a variable `that` of type  block  
        print {}   // with a variable `print` of type block of nothing
    } 
}
```

The problem is to know what to call a generic thing of a specific type ( e.g. list of strings vs list of ints)

```javascript
// print a specific type of thing
print: <E> { that E; /*do something with that*/}
pi: print<Int>{}
ps: print<String>{}
pi 1
ps 'Hi'
```

Without it print would not validate the type

```javascript
print: {that; /*something with that.. but what is thay? is any*/}
pi: print{}
pi 1
pi 'S' // valid because String is an instance of  `::{}`
```

It could probably be enforced to the first usage?(a: Yes) In this case `Int` but's it's kind of weird. How do you receive as a param? 

```javascript
execute: {
    action { that  } // action receives a `that` of type "any"
    action 'String'    // invokes `action` with the value `"String"` 
}
print: {that ; /*...*/} // defines a variable `print` of type block with a generic variable of tyep "any"
execute print  // calls `execute` with an `print` as "action"
// at this point `execute` "solidifies" into `execute { action { that String } } ` 
```

Problem: generic method vs generic object of specific type. 

```javascript
// generic blocks
Boo: {
    that
}
boo: {
    that
}
a: Boo{}
b: boo // boo{} wouldn't be valid
execute a // well you can't really execute Boo 'String' so this shouldn't be allowed.
execute b

```


How to tag a method where the there's no parameter and yet is generic? See example in [Generics](Questions/solved/Generics.md)

Probably do receive a parameter of type "type" by name it uppercase but it doesn't have a type and doesn't collide with a any single letter identifier: 
```javascript
Repo {
    T
    db DB
    find_by_id: { id Number
        table_name: info(t).name
        db.query_one('select * from {table_name} where id = ?', id)
    }
}
new_repo: {
    T
    db DB
    Repo{T db}
}

db = new_db() 
user_repo = new_repo(User db) 
post_repo = new_repo(Post db)
user: user_repo.find_by_id(1)
post: post_repo.find_by_id(1)
```

Ideal:

```javascript
Repo {
    table_name String
    db DB
    find_by_id: { id Int
        db.query_one 'select * from {table_name} where id = ? ' id
    }
}
new_repo: {
    table_name String
    db DB
    Repo(table_name db)
}
db = new_db()
user_repo = new_repo "user" db 
post_repo = new_repo "post" db
user: user_repo.find_by_id 1 
post: post_repo.find_by_id 1
```


Constraints: the generic is single letter uppercase (that leave us 26 options)


Another example: 
```javascript
compare: { a; b
    a < b ? { return -1 }
    a > b ? { return 1 }
    0
}
// T would be inferred as a block that has a method `<` that receives another item of the same type? and returns a bool? 
compare  { 
    a { < {{} Bool} > {{} Bool}}
    b {}
}

// compare<Int>
println(compare(1 0))
print compare 1  0
// compare<String>
println(compare('a' 'b'))
// compare<Decimal>
print compare 1.0 2.0
```


### Map/Reduce/Filter

Here is an example of how to write map, reduce, and filter functions for slices. These functions are intended to correspond to the similar functions in Lisp, Python, Java, and so forth.


```js
// Reduce reduces a []T1 to a single value using a reduction function.
reduce: {
    s []T; 
    initializer U; 
    f {U T U}
    r: initializer
    s.for_each {
        _ Int
        v T
        r = f r v 
    } 
    r
}

// Filter filters values from a slice using a filter function.
// It returns a new slice with only the elements of s
// for which f returned true.
filter: { s []T; f {T Bool}
    r []T
    r.for_each {_ Int; v T1
        f v ? { r << v }
    } 
    r
}
```
	

Here are some example calls of these functions. Type inference is used to determine the type arguments based on the types of the non-type arguments.

```js
	s: [1 2 3]
	//decimals: slices.map s, {i Int; new_decimal i } 
	// Now decimals is []Decimal = [1.0, 2.0, 3.0]

	sum: slices.reduce s 0 {i Int; j Int; i + j} 
	// Now sum is 6

	evens: slices.filter s {i Int; i % 2 == 0 } 
	// Now evens is []Int = [2]
```

#challenged using the generic inside the block would make it required which might be noisy, e.g. 

```javascript
List {
    T
    << { t T }
}

new_list:{
    a []T
    l : List{T}
    a.each{ item T
        l << item 
    }
}
l1: new_list [1 2 3] // could be like this but it would require: 
l2: new_list Int [1 2 3]
```

The ideal would be
```js
List {
    << { t }
}
new_list: {
    a []
    l: List{}
    a.each { item 
        l << item
    }
}
```


### Summary: 

It could work:
- ~~Just omit the type and the inference would create an interface in the block signature with the used methods~~
- ~~If the type is needed use a single letter type `e.g. T U S` etc~~
- ~~For "tagged" methods it would require to pass the type as parameter.~~

