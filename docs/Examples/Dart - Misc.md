#example

```js
print('hello world')
```

```js
name: 'Voyager I'
year: 1977
antenna_diameter: 3.7
fly_by_objects: ['Jupiter', 'Saturn', 'Uranus', 'Neptune']
image: {
    'tags': ['saturn']
    'url': '//path/to/saturn.jpg'
}
```

```js
year >= 2001 ? {
    print('21st century')
}, {
    year >= 1901 ? {
        '20th century'
    }
}

fly_by_objects.each({ object String
    print('${object}')
})

// might need to use `1.to(12)` syntax instead
1.to(12).each({ month Int
    print('${month}')
})
while({ year < 2016 }, {
    year = year + 1
})

```

```js
fibonacci: { n Int
   (n == 0 || { n == 1 }) ? { n }

   fibonacci(n - 1) + fibonacci(n - 2)
}
// might need to specify type
fibonacci: { n Int
   (n == 0 || { n == 1 }) ? { n }

   fibonacci(n - 1) + fibonacci(n - 2)
}
result: fibonacci(20)

```

```js
// This is a normal one-line comment.
/*
    Multiline
*/
// Value of n is: ${n}
n: 0
// no longer info(n)  // Value of n is: 0
```

```js
Element: lib1.lib1.Element
lib2: lib2.lib2

element: Element = Element()
element2: lib2.Element = lib2.Element()

```

```js
foo: lib2.lib2.foo
// module system is the file system; no import needed
```

[[../Questions/solved/concurrency/How to do concurrency]] TBD

```js

`!:[Example, Deprecate]`
Television: {
    // Use [turn_on] to turn the power on instead
	`
	example: {
		tv: Television()
		tv.activate()
	}
	deprecate: "Use turn_on instead"
    `
    activate: { turn_on() }

    // Turns the TV's power on
    turn_on: {
       ...
    }
}

// Use the element's info
tv: Television() // deprecate notice on build time
// generates Example docs
```

```js
// The example below is broken beyond repair, is kept for historical purposes
// Write can have an I/O error usually retuns the number of bytes written
write:  data [Int](); r Int; eh: {Int}
 ...
}
n: write([1, 2, 3], eh: { e Int
     match { e == write.ERROR => { print('Error: ${info(e)}') } }
 })
n: write([4, 5, 6])

// What was intended was this: 
// A write boc that takes an array of ints
// and returns a result with int as the number of writtent bytes if Ok
// and a strcuture with `written` and `message` if there is n error
write #(data [Int], Result(Int,#(written Int, message String))) {
     ... 
}
// write it
res: write([1,2,3])
match res {
	Ok => /* do nothing */
}, {
	Err => print("Error message: ${res.message}, bytes written ${res.written}") 
}
// also 

written : write([1,2,3]).or({
	e Result.Err(#(written Int, message String))
	print("Error message: ${e.message}, bytes written ${e.written}")
	e.written
})
// written has either way the number of bytes written


```