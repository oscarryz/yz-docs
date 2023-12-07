https://zig.news/edyu/zig-if-wtf-is-bool-48hh

```javascript
## Basic _if_ statement
if true { 
    println 'Hello Ed'
}  {
    println 'Hello World!'
}
... 
dude_is_ed: { 
    name String 
    if name == 'Ed' { 
        true()
    } { 
        false()
    }
}
// or just
dude_is_ed: { String;  name == 'Ed' }

say_hello: {
    name String
    if dude_is_ed {
        println 'Hello $(name)'
    } {
        println 'Hello world'
    }
}
...
## Error-handling _if_ statement
// Returns an error
// { Result } // Result<String>
dude_is_ed_or_error: {
    name String
    if name == 'Ed' { 
        ok(name)
    } {
        err 'Wrong person'
    }
}
say_hello: { 
    name String
    dude_is_ed_or_error(name) ? {
        println 'good seeing you $(name) again' 
    } {
        // Result.? { ok_body {v}  error_body {Error} }
        err(Error)
        println 'got error $(err)'
    }
}
say_hello_ignore_error { 
    name String
      dude_is_ed_or_error(name) ? {
        println 'good seeing you $(name) again' 
    } {
        // Result.? { ok_body {v}  error_body {Error} }
        _ Error
        println 'Hello world'
    }
}
## Mixing boolean with error-handling _if_ statement
// {name String Result}
dude_is_edish_or_error : {
    name String
    when [
        {name == 'Ed'}: {
            print 'Hello $(name)'
            ok(true())
        }
        {name == 'Edward'}: {
            println 'Hello again $(name)'
            ok(false())
        }
        {when.default}: {
            err('Wrong person')
        }
    ]
}
say_hello_edish: {
    name String
    dude_is_edish_or_error .or_err {
        ok Bool
        println 'ed? $(ok)'
        ok ? { 
            println 'Goog seeing your $(name)'
        } { 
            println 'good sseeing you again $(name)'
        }
        
    } {
        err Error
        println 'Got error $(err)'
        println 'Hello world'
    }
}
## Optional _if_ statement
null_boolean Boolean{}
dude_is_maybe_ed: {
    name String
    when [
        {name == 'Ed'}: {
            print 'Hello $(name)'
            true()
        }
        {name == 'Edward'}: {
            println 'Hello again $(name)'
            false()
        }
        {when.default}: {
            println 'Hello world'
            null_boolean
        }
    ]
}
say_hello_maybe_ed: {
    name String
    when_eq dude_is_mabe_ed [
        {true}: {println 'Hello $(name)'}
        {false}: {println 'Hello again $(name)'}
        {null_boolean}: {println 'Hello world'}
    ]
}
// also
say_hello_maybe_ed: {
    name String
    println when_eq dude_is_maybe_ed(name) [
        {true}: {'Hello {name}'}
        {false}: {'Hello again $(name)'}
        {null_boolean}: {'Hello world'}
    ]
}
## Optional _if_ statement with error-handling
dude_is_maybe_ed_or_error: {
    name String
    when_eq [
        {'Ed'}:{ok(true())}
        {'Edward'}:{ok(false())}
        {'Eddie'}:{err('Wrong person')}
        {when.default}: {ok(null_boolean)}
    ]
}
say_hello_maybe_ed_or_error: {
    name String
    dude_is_maybe_ed_or_error .or_err {
        ok Bool
        println "ed? {ok}"
        when_eq ok [
            {true}: { println 'Hello $(name)' }
            {false}: { println 'Hello again $(name)' }
            {null_boolean}: { println 'Good bye $(name)' }
        ]
    
    } {
        err Error
        println 'Got error $(err)'
        println 'Hello world'
    }
}
## Bonus
greeting dude_is_ed("Ed") ? { 'Hello'} { 'goodbye'}


```
