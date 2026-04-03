#example

https://www.softax.pl/blog/rust-lang-in-a-nutshell-1-introduction/

### Variables

```js
    x: 3
    y: 9.0

    s1: 'Hello world'
    s2: s1

    // Functions
    test: { x Int }
    add: { x Decimal; y Decimal; r Decimal }

    own_and_forget: { v List(String) }
    print: { v List(String) }
    change: { v List(String) }
    main: { result Result(Unit, std.io.Error) }

    foo: { x Int
        x > 0 ? { x }, { x + 1 }
    }

    foo: { x Int
        x > 0 ? {
          x
        }, {
          x + 1
        }
    }
    add: { x Decimal; y Decimal; r: x + y }
    f: add
    res: f(5, 7)
    f2 #(Decimal, Decimal, Decimal) = add

    Color: { Int; Int; Int }
    Nil: {}
    Foo: { text String }
    foo: Foo("Hi")
    foo2: Foo(text: "Hi")
    foo3: Foo()
    foo3.text = "Foo3"

    foo: {
        Foo: {
            self Foo
            text: "Hi"
            to_myself: {
                print('`self.text`')
            }
        }
        new: {
            s String
            f: Foo()
            f.self = f
            f.text = s
            f
        }
    }
    a: foo.new("Nada")
    a.to_myself() // Nada

    Node: {
        self Node
        data String
        to_string: {
            self.data
        }
    }
    s: Node()
    s.self = s
    s.data = 'Hola'

    HelloMacro: {
        hello_macro: {}
    }

    Pancakes: {
        algo: {}
    }

    //main.yz
    listener: tcpListener.bind("1237.0.0.1:7878")

    listener.incoming().each({ stream Stream
    })


    // main.yz
    "Sorts the args descending order and greets them"
    args: os.args.sort({ a String; b String
        a <=> b
    })
    args.each({ s String
        print('Hello `s`')
    })

    correct: "Pass123"
    guess_password: {
        tries: 1

        guess: input("What is the password?")
        guess == correct ? {
          print("That's correct")
        }, {
          tries >= 3 ? {
            print("Too many attempts")
          }, {
            guess_password(tries + 1)
          }
        }
    }

```

[[../Features/Replaced features/Private read only variables]]
```js
// This is how smalltalk does it, not necessarily possible with yz
BlockWithExit: {
    exit: {}
    block: {}
    on: {
        a_block: {}
        block = a_block
    }
    value: {
        exit: { break }
        block()
    }
}
// There's no block block but ok
Block: {
    withExit: {
        bwe: BlockWithExit()
        bwe.on(self)
    }
}

max_before_nil: {
    collection List(Int)
    max: 0
    supplier: collection.read_stream()
    while({ supplier.is_empty() == false }, {
        v: supplier.next()
        v == Option.None() ? { break }, { }
        max = max.max(v)
    })
    max
}
`
// Iterate an Int array and exit if value is null
// update max
yz>max_before_nan [1,2,3, core.int.NOT_A_NUMBER, 4, 5]
>> 3
`
max_before_nan: { numbers [Int]()
    max: { a Int; b Int; a > b ? { a }, { b } }
    m: 0
    numbers.each({ n Int
        n == core.ints.NOT_A_NUMBER ? {
            break
        }, { }
        m = max(m, n)
    })
}
```

Early in smalltalk
```smalltalk
    maxBeforeNil: aCollection
      | max supplier loop |
      max := 0.
      supplier := aCollection readStream.
      loop := [
        [supplier atEnd]
          whileFalse: [
            value := supplier next.
            value isNil ifTrue: [loop exit].
            max := max max: value]
       ] withExit.
      loop value.
      ^max
    ```
```js
    max: {
        a Int
        b Int
        a > b ? { a }, { b }
    }
    gah: {
        list List(T)
        block: { t T }
        list.is_empty() ? {
            return
        }, { }
        t: list.head()
        block(t)
        gah(list.tail(), block)
    }

    { 1 + 2 }() //-> 3

```