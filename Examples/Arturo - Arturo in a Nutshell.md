#example

https://arturo-lang.io/documentation/in-a-nutshell/

```js
// variables
a1: 2
a2: 3.14
a4: Complex( 1, 2.0 ) // 1.0 + 2.0i

// strings
c1: 'this is a string'
c2: `
        this is a multilie string
        that is verbatim
    `
// no characters, use string
ch: 'c' // still a string

// arrays
d: [1, 2, 3]

// dictionaries
e: ['name'   : 'John'
    'surname': 'Doe'
    'age'    : '34' // can't mix value types
    'likes'  : `['pizza', 'spaguetti']`]

// blocks (behave like closures/functions)
f: {
    x Int
    2 * x
}
f(2) // 4

// dates
h: date.now() // 2021-05-03T17:10:48+02:00

// core values
i1: true // std.true
i2: false // std.false



//---------------------------------------------------
// Aritmethic is implemented in the std.Number block
// eg. 1 + 1 is the same as 1.+(1) or 1.plus(1)
//---------------------------------------------------
1 + 1      // 2
8 - 1
4.2 - 1.1
10 * 2
35 / 4     // 8
35.0 / 4   // 8.75
2 ** 5     // 32
5 % 3      // 2

3 & 5      // 1
3 | 5      // 7
3 ^ 5      // 2

//------------
// constants
//------------
math.PI
meth.EPSILON

//---------------------
// Comparison operatos
//----------------------
// equality
1 == 1     // true
1 == 2     // false

// inequality
1 != 1     // false
1 != 2     // true

// more comparisons. They are still "methods"

1 < 10      // true
1 =< 10     // true
10 =< 10    // true
1 > 10      // false
1 >= 10     // false
11 >= 10    // true


//--------------
// Conditionals
//--------------

// logical operators
true  && true  // true
true  && false // false
true  || false // true
false || false // true

(1 == 2) && {2 < 3} // false
                    // the second block will not be evaluated

// simple "if" statements
2 > 1 ? { print('yes' }   // yes)
3 != 2 ? { print('true!'} // true!)

// if / else statements
2 > 3 ? { print('2 is greated than 3'},)
{ print('2 is not greater then 3'} // 2 is ot greater than 3)

// "switch" statements

when [
   {2 > 3}: { print('2 is greater than 3'})
   {true}:{ print('2 is not greater than 3'} // 2 is not greater than 3)
]

a: 2 > 3 ? {'yes'}, {'no'} // a: 'no'

// "loops"
arr: [1 4 5 3]
arr.each({ x Int
    print('x = `x`')
}
// x = 1
// x = 4
// x = 5
// x = 3

// with index
arr.for { i Int, x Int
    print('Item at position `i` => `x`')
}
// item at position 0 => 1
// item at position 1 => 4
// item at position 2 => 5
// item at position 3 => 3


// ranges
1.to(3).do{ x Int
    print('`x`')
}
// 1
// 2
// 3

'a' .to 'c' .each({ s String
    print(s)
}
// a
// b
// c

// looping through a dictionary
dict: ['name': 'John' 'surname': 'Doe']
dict.each({ key String, value String
    print('`key` -> `value`')
}

// name -> John
// surname -> Doe

//while loops
i: 0
while { 1 < 3 }, {
    print('i = `i`')
    i = i + 1
}
// i = 0
// i = 1
// i = 2

// Strings
a: 'tHis Is a stRinG'
print(a.upper()) // THIS IS A STRING

// concatenation
a: 'Hello ' ++ 'World!'
'hello'.split() // ['h' 'e' 'l' 'l' 'o']
'hello world'.split(' ') // ['hello' 'world']

// conversion
123.to_string() // "123"

int.parse!("123") // 123
int.parse("123") // std.option.Ok(123)
int.parse(":(") // std.option.Err("...")


// string interpolation `
x: 2
print('x = `x`' // x = 2)

//--------
// Blocks
//--------

// calculate block
sth: {print('Hello world'})
sth() // Hello world

// array indexing
arr: ['zero' 'one' 'two' 'three']
arr[0]      // zero
arr[arr.len()-1] // three
arr[3]           // three

x: 2
arr[x]            // two
arr[0] = 'nada'
print('`arr`' // ['nada' 'one' 'two' 'three'])

// adding elements
arr: [] String
arr << 'one'
arr << 'two'
print(arr // ['one', 'two'])
// also
arr.add('one')
arr.add('two')

// remove elements from array
arr: ['one' 'two' 'three' 'four']
arr.remove 'two'  // ['one', 'two', 'three', 'four']
arr.remove_at(0)   // ['three', 'four']

// getting the size fo an array
arr: ['one' 'two' 'three' 'four']
print("`arr.len())`") // 4

// getting a subarray
arr: ['one' 'two' 'three' 'four']
arr.subarray 0,1  // ['one', 'two']
arr.contains 'one' // true
arr.contains 'five' // false


// sorting
arr: [1 5 3 2 4]
arr.sort() // [1,2,3,4,5]
arr.sort { a Int, b Int, b - a} // descending order -> [5,4,3,2,1]

// mapping values
1 .to 10 .each({x Int, 2 * x}
Range(1,10).map{x Int, 2 * x} // [2,4,6,8,10,12,14,16,18,20]

n: [Int]()
1.to 10 .each({ x Int, n << x }
n.select { x Int, x % 2 == 0} // [2,4,6,8 10]
n.filter { x Int, x % 2 == 0} // [1,3,5,7,9]

// misc operations
arr: ['one' 'two' 'three' 'four']
arr.reverse() // ['four', 'three', 'two', 'one']

//------------------
// function / blocks
//------------------

f: {x Int, 2 * x }
f(10) // 20

// returning a value
g: { x Int
    x < 2 ? { 0 } // finishes the block early
    res: 0
    0 .to x .each({ z Int
        res = res + z
    }
    res
}
g(1) // 0
g(3) // 6

//-------------
// custom types
//-------------

Person: {
    name String
    surname String
    age Int
    print: {
        print('NAME: `name`, SURNAME: `surname`, AGE: `age`')
    }
    compare: {other Person, age - other.age}
    say_hello:{ print('Hello `ame`'})
}
// "constructor"
new_person: {name String, surname String, age Int
    Person( name.capitalize(), surname, age)
}

// create new instances of custom type
a: Person('John', 'Doe', 34)
b: Person(
    name:'Jane'
    surname:'Doe'
    age:34
)
a.say_hello() // Hello John
b.say_hello() // Hello Jane

// access block variables
print('First person name is `a.name`')
print('Second person name is `b.name`')

// modifing values
a.name = 'bob'
a.say_hello() // Hello bob

// verifying type
std.info(a).type // Person
std.is_a Person, a // true
std.info(Person).type // #(String, String,Int)

a.print() // NAME: John, SURNAME: Doe, AGE: 34

```
