#example

```js
print('hello world')
```

```js
name: 'Voyager I'
year: 1977
antenna_diameter: 3.7
fly_by_objects: ['Jupiter' 'Saturn' 'Uranus' 'Neptune']
image: {
    'tags': ['saturn']
    'url': '//path/to/saturn.jpg'
}
```

```js
year >= 2001 ? {
    print('21st century')
} {year >= 1901 ?  {
        '20th century'

}}

// when function
when  [{ year >= 2001}: { print('21st Century'})
      { year >= 1901}: { print('20st Century'}])


fly_by_objects.each({ object String
    print('{object}')
}

// might need to use `1.to(12)` syntax instead
1 .to(12).each({ month Int
    print('{month}')
}
while { year < 2016} {
    year = year + 1
}

```

```js
fibonacci: { n Int
   (n == 0 || { n == 1 }) ? { n }

   fibonacci(n -1) + fibonacci(n - 2)
}
// might need to specify type
fibonacci: { n Int
   (n == 0 || { n == 1 }) ? { n }

   fibonacci(n -1) + fibonacci(n - 2)
}
result: fibonacci(20)

```

```js
// This is a normal one-line comment.
/*
    Multiline
*/
'Value of n is: {n}'
n: 0
info n  // Value of n is: 0
```

```js
Element: lib1.lib1.Element
lib2: lib2.lib2

element Element = Element{}
element2 lib2.Element = lib2.Element{}

```

```js
foo: lib2.lib2.foo
// all except foo
import lib2.lib2
```

[[../Questions/solved/concurrency/How to do concurrency]] TBD

```js

Television: {
    // Use [turn_on] to turn the power on instead
    `Example:
     yz> tv: Television{}
     ... tv.activate()


    @Deprecated('Use turn_on instead')
    Something else {something_else()}
    `
    activate: { turn_on() }

    // Turns the TV's power on
    turn_on: {
       ...
    }
}

// Use the element's info
tv: Television{}
info(tv.activate) // @Deprecated('Use turn_on instead')
process_annotations(info(tv.activate)) // Reads the `@` elements
process_doc(info(tv.activate)) // Read the `Examples:` section etc.
process_xyz(info(tv.activate)) // might read something else...
```

```js
// Write can have an I/O error usually retuns the number of bytes written
write: { data [Int]() r Int eh: {Int}
 ...
}
n: write([1 2 3] eh:{e Int;
     when [{ e == write.ERROR }:{ print('Error: {info(e))}'}]
 })
n: write([4 5 6])

```
