#example

https://tryhaskell.org/

```js
23 * 26
reverse "Hello"
[13 23 40]
{28 "chirs"}
{28 "chirs"}.0 // 28
{x Int;x * x}(4) // 16
{x Int; x + x}(8 * 10) // 160
{villain {Int String}; villain.0 }({28 "chirs"}) // 28
[] << 'a' // ['a']
'' ++ 'a' // "a" or 'a'
[] << 'a' << 'b' == [ 'a' 'b']
1.to(5).map 1.+ // [2 3 4 5 6]
1.to(10).map 99.*
[63 3 5 25 7 1 9].filter 5.> // [63 25 7 9]

square: {x Int; x * x}
square(52) // 2704

add1: {x Int; x + 1} /*in*/; add1 5 // 6
second: {x {Int Int}; x.1} /*in*/; second {3 4} // 4

square: {x Int; x * x} /*in*/; 1.to(10).map square

square ::{Int} = {x * x}
square:{x; x*x}

// rejected
add_three ::{x Int y Int z Int} = {
    x + y + z
}
add_three:{
    x + y + z
}
// Yz 1.0 ? //see: [[Re-declare variables if there's type signature?]]
add_three { x Int y Int z Int } = {
    x + y + z
}

let (a,b) = (10, 12) in a * 2
fn: {a Int; a * 2 }
fn(10 20) // not possible because of 2 params

input: {10 "abc"}
fn: {
    input {Int String}
    input.1.char_at(0)
}




```
