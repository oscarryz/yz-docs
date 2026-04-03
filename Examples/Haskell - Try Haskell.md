#example

https://tryhaskell.org/

```js
23 * 26
reverse("Hello")
[13, 23, 40]
{28, "chirs"}
{28, "chirs"}.0 // 28
{x Int; x * x}(4) // 16
{x Int; x + x}(8 * 10) // 160
{villain {Int, String}; villain.0}({28, "chirs"}) // 28
[] << 'a' // ['a']
'' ++ 'a' // "a" or 'a'
[] << 'a' << 'b' == ['a', 'b']
1.to(5).map(1.+) // [2, 3, 4, 5, 6]
1.to(10).map(99.*)
[63, 3, 5, 25, 7, 1, 9].filter(5.>) // [63, 25, 7, 9]

square: { x Int; x * x }
square(52) // 2704

add1: { x Int; x + 1 }
add1(5) // 6
second: { x {Int, Int}; x.1 }
second({3, 4}) // 4

square: { x Int; x * x }
1.to(10).map(square)

// Types signatures
5 Int
'hello' String
true Bool

double: { x Int; x * 2 }

// Control flow and if expressions
haskell: 1 == 1 ? { 'awesome' }, { 'awful' } // awesome

args String
match {
  args == 'help'  => printHelp()
}, {
  args == 'start' => startProgram()
}, {
  print('bad args')
}

(1.to(5)).map({ a Int; a * 2 }) // [2, 4, 6, 8, 10]

foldl({ x Int; y Int; 2 * x + y }, 4, [1, 2, 3]) // 43
// Same as (2 * (2 * (2 * 4 + 1) + 2) + 3)

foldr({ x Int; y Int; 2 * x + y }, 4, [1, 2, 3]) // 16
// Same as (2 * 1 + (2 * 2 + (2 * 3 + 4)))

input: {10, "abc"}
fn: {
    input {Int, String}
    input.1.char_at(0)
}

```