#example

https://hyperpolyglot.org/ml

```js
// blocks
{}
//
/*
*/
a:2
a Int = 2
// coalesce
a Optional(String) = something()
a.ifPresent { do_something_else() }
when [
    { a.isPresent }: {use a}
    { a.isPresent == false }: {use default_value }
]
// better
r: a.or{default_value}
foo == -999 ? { None() } { Some(n) }
r: foo == -999 ? {Err()},{Ok(n)}

literal: [1 2 3]
empty_list: []Int
empty_list_test: ([]Int).len() == 0
// or maybe
empty_list_test: len([]) == 0
cons: [1] << 2 // [1 2] 2 >> [1]
employee >> list // ?? what is this?
head: list[0]
tail: list.sub_array(0)
l:   list.len()
nth: list[n]
n:   [7 8 9].index_of(8) // 1
//iterate
list.each({ e Int; print('`e`'})


doubled: [1 2 3 4].map { x Int; x * 2}
all_gt_2 : true
[1 2 3 4].each({x Int; all_gt_2 = all_gt_2 && {x > 2}}
[1 2 3 4].all { x Int; x > 2 } // false
[1 2 3 4].filter( 2.> )
//
Array: {
  T
...
    all: {
        predicate #(T,Bool)
        result: true
        data.each({ e T
           result = result && {block(e)}
        }
        result
    }
}

functions: {
    average: { a Int; b Int; a + b / 2.0 }
    // invoke
    average(1, 0)
    named_params: {
        m Int
        s Int
    }
    named_params(s: 1 m: 2)
    to_s: { x Color
        when_eq x [
          {Red}: {'red'}
            {Green}: {'green'}
            {Blue}: {'blue'}
        ]
    }
    recursive: {x Int;
        x > 0 ? { recursive x - 1 }
    }
    // anon funcs
    { x Int; x + 1}
    // currying ?
    plus: {a Int; b Int; a + b}
    plus2: { a Int; plus(a, 2) }
    plus2(5) // returns 7
    // Also
    plus2: 2.+
    plus2(5) // returns 7
    // composition
    f: {x Int; x + 2}
    g: {x Int; x * 3}
    f(g(4))
   // f g 4
    // lazy eval
    arg1: {x Int; y #();  x }
    arg(7,{panic()})

    // type synonym
    Name: String

}
