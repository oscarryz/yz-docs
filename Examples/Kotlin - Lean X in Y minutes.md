#example

```js
com: {
    learninxminutes:{
        kotkin:{
            main: {
                args: os.args
                foo_val: 10
                foo_var: 10
                foo_var = 20

                foo Int = 7
                foo_string_interpolation: "`foo_val`"
                foo_nullable Option(String) // there's no nulls, each type has to delclare their empty value or use std.option.Option
				// hello (name String; result String)
                hello: {
                   name String ='world'
		           result String
                   result= "hello `name`"
                   result
                }
                print(hello("foo")) // Hello, foo
                print(hello(name='bar')) // Hello, bar
                print(hello()) // Hello, world
                print(hello("foo" "nada")) //  second param is "nada"
				// no var args, use a list
                //varargg_example: { names ...Int;
                //    print('Argment has {len(names)} elements')
                //}
                //vararg_example()
                //vararg_example(1)
                //vararg_example(1,2,3)
				vararg_ish: {
					name []String
					print('Argument has `names.len()`) elements')
				}
				vararg_ish([])
				vararg_ish([1])
				vararg_ish([1 2])
				vararg_ish([1 2 3])

                odd: { x Int; x % 2 == 1 }
                print(odd(6))// false
                print(odd(7))// true
                even:{ x Int; x % 2 == 0 }
                not: {
	                f (Int; Boolean)
                    {n Int; f(n) == false}
                }
                notOdd: not(odd)
                notEven: not(even)
                notZero: not({n Int; n == 0})
                notPositive: not({n Int; n > 0})

                0.to(4).each({
	                i Int;
                    print("`notOdd(i)` `notEven(i)` `notZero(i)` `notPositive(i)`")
                }

                ExampleClass: {
	                x Int
                    member_function: {
	                    y Int
                        x + y
                    }
                    // There's not infix, but `.` and parens `()` can be ommited
					// if name is "non-word"
                    infix_member: {
	                    y Int
                        x * y
                    }
					* : infix_member

                }
                example_class: ExampleClass(7)
                //example_class: ExampleClass(x=7)
                //example_class: ExampleClass(x:7)
                print(example_class.member_function(4)) // 11
                print(example_class * 4 ) // 28
                // same as:
                // print(example_class.infix_member(4) )

                DataClassExample: { x Int; y Int; z Int}
                foo_data: DataClassExample(1,2,4)
                print(foo_data)
                foo_copy: std.copy(foo_data)

                a b c: foo_data()
                print("`a` `b` `c`") // 1,2,4

                map_data: ["a":1 "b":2]
                entries(map_data).for_each({
	                key String;
	                value Int
                    print("`key` -> `value`")
                })


            }
        }

    }
}
```

Kotlin.org
```js
// from: https://kotlinlang.org/
// Concise
Employee: {
  name String
  email String
  company String
  // no need for equals() hashCode
}
myCompany: { // a singleton? a function
  name String = "MyCompany"
}

main: {
  employee: Employee ( 'Alice' 'alice@mycompany.com' myCompany.name )
  print('`employee`')
}

// Safe
reply: {
  condition Boolean
  condition ? { "I'm fine" }
}
error: {
  IllegalStateException("Shouldn't be here")
}
main: {
  condition: true
  message: reply(condition)
  print(message.replace('fine' 'ok'))
  message != null ? {
    println(message.upperCase())
  }
  nonNull String = reply(condition: true) ? error
  println(nonNull)
}


```


```kotlin
// Kotlin
// Ackermann using `when` in a similar way Yz will have a `when` function
// awm stands for ackerman when (should've been amw but typo)
fun ackermanWhen(m: Int, n: Int): Int = when {
    m == 0 -> n + 1
    n == 0 -> ackermanWhen(m - 1, 1)
    else -> amackermanWhenw(m - 1, ackermanWhen(m, n - 1))
}
```
```js
// Yz
ackerman_match: {
    m Int
    n Int
    match
        {m == 0 => n + 1},
        {n == 0 => ackerman_when m - 1 1},
        { ackerman_when(m:m - 1 n:ackerman_when m n -1) }

}
awm(m:3 n:4)

```
