#example

https://learnxinyminutes.com/docs/haskell/

```js
x: {
  r: []
  1 .to(5 ).each({ x Int

    x * 2 > 4 ? { r << x * 2}
  }
}
print("{x())}" // 2,4,6,7,10


// "tuples"
a b: {
  "haskell"
  1
}() // a: "haskell", b: 1


ctt: get(2, { "snd"; "can't touch this"; "da na na" })

// get {Int, <T>{...T}}
get: {
    index Int
    block <T>{...T}
    a b c = block()
    b
}

add: {a Int; b Int; a + b }
//
Adder: {
  + {Adder Adder}
}
Int: {
  + {other Int
  self + other // How to resolve self?
  }
}
add: <Adder>{ a Adder; b Adder; a + b }
add(1,2)

1 + 2

\ : { a Int; b Int; a div b }
35 \ 4 // 8

fib: {x Int
  x < 2 ? { <- 1 }
  fib(x - 1) + fib( x - 2)
}

fib: {x Int;
  when(
    { x == 1 } , { 1 }
    { x == 2 } , { 2 }
    { true }, { fib(x - 1 ) + fib( x - 2 ) }
  )
}
tuple: { "snd"; "can't touch this"; "da na na" }

_, b, _: tuple()
b

my_map: { array []A; f <A,B>{A;B}
  len(array) == 0 ? { <- []B }
  r: my_map(sub_array(array, 0), f)
  push(r, f(array[0]), 0)
}

(1 .. 5).collect({ x Int; x + 2 })

add: { a Int; b Int; a + b }
foo: { c Int; add(c, 10) }

foo(5) // 15

foo: { a Int; 10 + a }
foo(5) // 15

foo: {a Int; (10 + a) * 4 }
foo(5) // 60

// Types signatures


5 Int
'hello' String
true Bool

double {Int; Int} = { x Int; x * 2 }

// Control flow and if expressions
haskell: 1 == 1 ? {'awesome'}, {'awful'} // awesome

// Identation is not important
haskell: 1 == 1 ?
  {'awesome'},
  {'awful'}

args String
when(
  {args == 'help'}, {printHelp()},
  {args == 'start'}, {startProgram()},
  {true}, { print('bad args') }
)

(1 to 5).map({a Int; a *2 }) // [2,4,6,8,10]
for: { array []A; f <A,B>{A,B}
    array.map(f)
}
for( [1,2,3,4,5], {a Int; show(a)} )

for([1,2,3,4,5], show)

foldl({x Int; y Int; 2 * x + y}, 4, [1,2,3]) // 43
// Same as (2 * (2 * (2 * 4 + 1) + 2) + 3)

foldr({x Int; y Int; 2 * x + y}, 4, [1,2,3]) // 16
// Same as (2 * 1 + (2 * 2 + (2 * 3 + 4)))

Color: {}
red:Color{}
blue:Color{}
green:Color{}

say {Color,String} = {
    color Color
    when(
        {color == red}, {'You are Red!'},
        {color == blue}, {'You are Blue!'},
        {color == green}, {'You are Green!'},
    )
}

Point {Float;Float} = {
    x Float; y Float;
}
/*
Point: {
    x Float
    y Float
}
*/

distance {Point; Point; Float} = {
    a Point; b Point
    dx: (a.x - b.x) ** 2
    dy: (a.y - b.y) ** 2
    sqrt(dx + dy )
}

Name: {
    name: ""
    last_name: ""
    second_last_name: ""
}
m: Name{'Oscar'}
n: Name{'Oscar', 'Reyes'}
o: Name{'Oscar', 'Reyes', 'de la Cruz'}

CartesianPoint: {
    x Float
    y Float
}
PolarPoint2D: {
    r Float
    theta Float
}
my_point: CartesianPoint2D{7.0; 10.0}
x_of_my_point = my_point.x

my_point_2: CartesianPoint{ 9.0; my_point.y}
my_point_3: CartesianPoint{ 3.3; 4.0}

distance_from_origin: {
    p: {x Float;y Float}
    when(
        {meta_data(p).block_name == 'CartesianPoint2D'},
            { sqrt(p.x ** 2 + p.y ** 2)}
        {meta_data(p).block_name == 'PolarPoint2D'}, { p.r }
    )
}

Maybe: <T>{a T}
Nothing: {}
Just: <T>{a T; a}
Just{'Hello'}
Just{1}
Nothing{}

// TODO: type synonyms
Int: Integer // String = [Char]
Weight: Float
Height: Float
Point: (Float, Float)
get_my_height_and_weight { Person; Height; Weight }
find_center { Circle; Point }
some_person { Person }
some_circle {Circle}
distance { Point; Point; Point}

// 8  Typeclasses
Eq: <T>{
    == {Eq;Eq;Bool}
    != {Eq;Eq;Bool}
    != {Eq;Eq;Bool}
}


```
