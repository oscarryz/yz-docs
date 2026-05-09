#example
https://zig.news/edyu/zig-if-wtf-is-bool-48hh

```js
// Basic conditional
true ? {
    print('Hello Ed')
}, {
    print('Hello World!')
}
...
dude_is_ed: {
    name String
    name == 'Ed'
}
// or just
dude_is_ed: { name String; name == 'Ed' }

say_hello: {
    name String
    dude_is_ed(name) ? {
        println('Hello ${name}')
    }, {
        println('Hello world')
    }
}
...
// ## Error-handling conditional
// Returns an error
// { Result } // Result<String>
dude_is_ed_or_error: {
    name String
    name == 'Ed' ? {
        Result.Ok(name)
    }, {
        Result.Err('Wrong person')
    }
}
say_hello: {
    name String
    dude_is_ed_or_error(name).and_then({ name String
      println('good seeing you ${name} again')
    }).or({ e Error
        println('got error ${e}')
    })
}
say_hello_ignore_error: {
    name String
    dude_is_ed_or_error(name).Ok ? {
        print('good seeing you ${name} again')
    }, {
        // Result.Ok? #( if_ok #(V), if_err #(Err(V)))
        print('Hello world')
    }
}
// ## Mixing boolean with error-handling conditional
dude_is_edish_or_error #(name String, Result(String)) = {
    name String
    match {
        name == 'Ed' => {
            print('Hello ${name}')
            Result.Ok(true)
        }
    }, {
        name == 'Edward' => {
            println('Hello again ${name}')
            Result.Ok(false)
        }
    }, {
        Result.Err('Wrong person')
    }
}
say_hello_edish: {
    name String
    dude_is_edish_or_error(name).and_then({ ok Bool
        println('ed? ${ok}')
        ok ? {
            println('Good seeing you ${name}')
        }, {
            println('Good seeing you again ${name}')
        }
    }).or({ err Error
        println('Got error ${err}')
        println('Hello world')
    })
}
// ## Optional conditional
dude_is_maybe_ed: {
    name String
    match {
        name == 'Ed' => {
            print('Hello ${name}')
            Option.Some(true)
        }
    }, {
        name == 'Edward' => {
            println('Hello again ${name}')
            Option.Some(false)
        }
    }, {
        println('Hello world')
        Option.None()
    }
}
say_hello_maybe_ed: {
    name String
    result: dude_is_maybe_ed(name)
    match {
        result.Some && result.v == true => println('Hello ${name}')
    }, {
        result.Some && result.v == false => println('Hello again ${name}')
    }, {
        println('Hello world')
    }
}
// ## Bonus
greeting: dude_is_ed("Ed") ? { 'Hello' }, { 'goodbye' }


```
