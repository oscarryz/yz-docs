#example
Where X = Yz

```js

// Single line comments
/*
Multiline comments
*/
"A string on top of a variable can be retrieved by calling `info(thing)` at runtime. They are multi line"
thing Int

// Statements terminated by ;
do_stuff()

// but they don't have to be, as semicolons are automatically inserted
// wherever there's a new line, except in certain cases
do_stuff()


// Numbers, String and Operators
// One number type
1 // = 1
1.5 // = 1.5

1 + 1 // = 2
0.1 + 0.2 // 0.30000000000000004
5 / 2 // 2.5

// Boolean
true
false

'String'
"String"
`String`

// Equality:
1 == 1 // true
1 == 2 // false
1 != 1 // false
1 != 2 // true

'String' ++ 'concatenation'



first_elem: { a [T](); a[0] }
second_elem: { a [T](); a[1] }
`"
fun evenly_positioned_elems (odd::even::xs) = even::evenly_positioned_elems xs
  | evenly_positioned_elems [odd] = []  (* Base case: throw away *)
  | evenly_positioned_elems []    = []  (* Base case *)
"`
evenly_positioned_elements: { a [T]()
    a.length() == 0 ? { a }, {
        a.length() == 1 ? { [T]() }, {
            a[1..].each({ // even positions
            })
        }
    }
}

```

https://learnxinyminutes.com/docs/go/

```js
// single line
/*
  multi line
*/

"meta information" // Calling `info(variable)` returns 'meta information'
variable: 1

// No imports needed but you can "alias" a module
m: math // alias for math module
http: net.http // you can use http now directly

print("Hello world")

beyond_hello()
beyond_hello: {
    x Int // Var declaration
    x = 3 // var assignment
    y: 4  // Short declaration with type infer, declare and assign
    sum, prod: learn_multiple(x, y) // take the last two variables "returned"
    print("sum: `sum` prod: `prod`") // simple output
    learn_types() // < y minutes, learn more!
}

/*
 Args and return values are the internal block variables.
 For convenience a "function()" syntax can be used
 Order of assignment is:
 Top down for args
 Top down starting at -n
*/
learn_multiple: {
    x Int
    y Int
    sum Int
    prod Int

    sum = x + y
    prod = x * y
}

// So you can call
learn_multiple(1, 2, 3, 4)
// x: 1, y:2, sum: 3, prod: 4
// then the body is executed and sum and prod are overridden and replaced with 3 and 2
print(learn_multiple.sum) // 3
print(learn_multiple.prod) // 2

// some built-in types and literals.
learn_types: {
    str: 'Learn Yz!'
    s2: 'A "raw" string literal
    can include line breaks'
    g: 'Σ' // non ascii literal
    f: 3.14195

    u: 7
    pi: 22.0 / 7

    n: Byte('\n')
    a4: Array.new_size([Int](), 4)
    a5: [3, 1, 5, 10, 100] // Array initialized with 5 elements
    a4_copy: Array.copy(a4)
    a4_copy[0] = 25
    print("`a4_copy[0] == a4[0]`") // false

    s3: [4, 5, 9]
    s4: Array.new_size([Int](), 4)
    d2: [][Decimal]()
    s3_copy: s3
}
// learn memory
learn_memory: {
    p: Int()
    s: array.new([Int](), 20)
    s[3] = 7
    r: -2
}
expensive_computation: {
    math.exp(10)
}
learn_flow_control: {
    true ? {
        print("told ya")
    }, { }
    false ? {
        // pout
    }, {
        // gloat
    }
    x: 42
    match {
        x > 2 => print('x is gt 2')
    }, {
        x < 2 => print('x is lt 2')
    }, {
        print('true')
    }
}

```
