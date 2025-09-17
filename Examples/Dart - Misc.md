#example

```javascript
print 'hello world'
```

```javascript
name: 'Voyager I'
year: 1977
antenna_diameter: 3.7
fly_by_objects: ['Jupiter' 'Saturn' 'Uranus' 'Neptune']
image: {
    'tags': ['saturn']
    'url': '//path/to/saturn.jpg'
}
```

```javascript
year >= 2001 ? {
    print '21st century'
} {year >= 1901 ?  {
        '20th century'

}}

// when function
when  [{ year >= 2001}: { print '21st Century'}
      { year >= 1901}: { print '20st Century'}]


fly_by_objects.for_each { object String
    print '{object}'
}

// might need to use `1.to(12)` syntax instead
1 .to 12 { month Int
    print '{month}'
}
while { year < 2016} {
    year = year + 1
}

```

```javascript
fibonacci: { n Int
   (n == 0 || { n == 1 }) ? { return n }

   fibonacci(n -1) + fibonacci(n - 2)
}
// might need to specify return type
fibonacci: { n Int
   (n == 0 || { n == 1 }) ? { return n }

   fibonacci(n -1) + fibonacci(n - 2)
}
result: fibonacci(20)

```

```javascript
// This is a normal one-line comment.
/*
    Multiline
*/
'Value of n is: {n}'
n: 0
info n  // Value of n is: 0
```

```javascript
Element: lib1.lib1.Element
lib2: lib2.lib2

element Element = Element{}
element2 lib2.Element = lib2.Element{}

```

```javascript
foo: lib2.lib2.foo
// all except foo
import lib2.lib2
```

[[../Questions/solved/concurrency/How to do concurrency]] TBD

```javascript

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

```javascript
// Write can have an I/O error usually retuns the number of bytes written
write: { data []Int r Int eh: {Int}
 ...
}
n: write([1 2 3] eh:{e Int;
     when [{ e == write.ERROR }:{ print 'Error: {info(e)}'}]
 })
n: write([4 5 6])

```
