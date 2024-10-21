
Compare two pieces of code (two blocks) and use some special notation to help the comparison

Example from: 
[https://flatironschool.com/blog/pattern-matching-elixir/](https://flatironschool.com/blog/pattern-matching-elixir/)

Display full name if user has both `name` and `last name`
Display `username` if only has that
Display generic message otherwise

```javascript

    User: {
        username String
        name String
        last_name String
    }
    
    o: User (
        name: 'Oscar'
        last_name: 'Reyes'
    )
    
    o2: User (
        username: 'oscarryz'
    )
    o3: User()

    // using when
    `
    tests: {
         o: User {
            name: 'Oscar'
            last_name: 'Reyes'
        }
        
        o2: User {
            username: 'oscarryz'
        }
        o3: User{}
        
        assert greet o == 'Hello Oscar Reyes'
        assert greet o2 ==  'Hello oscarryz'
        assert greet o3 == 'Hello new user'
    }
    `
    greet: {
        user User
        when [
            { user.name != '' && {user.last_name != ''}}: {
                print 'Hello {user.name} {user.last_name}'
            }
            { user.username != '' }: {
                print 'Hello {o.username}': 
            }
        ] { print 'Hello new user' }
    }
    
    // with `is_set` func
    greet: {
       user User
       when [
         { is_set [user.name user.last_name]}: { print 'Hello {user.name} {user.last_name}' }
         { is_set [user.username]           }: { print 'Hello {user.username}               }
        ] /*else*/ {print 'Hello new user'}
    }
    // Pattern matching similar to hamcrest
    greet: {
        user User
        // might not be possible
        match(o,
            {name:any(); last_name:any()}: {print('Hello {o.name} {o.last_name}')}
            {username:any()}             : {print('Hello {o.username}')}
            /* else */                     {print('Hello new user')}
        )
    }

//    yz>greet o 
//    Hello Oscar Reyes
//    
//    yz>greet o2 
//    Hello oscarryz
//
//    yz>greet o3 
//    Hello new user
//


```

Wihout (current)

```javascript

    greet: {
        o Person
        o.name != string.undefined && {o.last_name != string.undefined} ? {
            print('Hello {o.name} {o.last_name}')
        } {
            o.user_name != string.undefined ? {
                print('Hello {o.username}')
            } {
                print('Hello new user')
            }
        }
    }

```
Another from [Simple password check](https://gist.github.com/oscarryz/4b182cf7a8a696f4acdee38159f003e6)

Guess a password in 3 attempts 

```javascript
    // Simple strightfoward solution
    correct: 'pass123'
    guess_password: {
        attempts: 0
        guess: input('Guess the password')

        guess == correct ? {
            print("You got it right!")
        }, {
            attempts > 3 ? {
                print("To many attemps, account is locked")
            }, {
                guess_pasword(attempts + 1)
            }
        }
    }
    
```

 // Now with "pattern matching"
```javascript
    correct: 'pass123'
    guess_password: {
        attempts: 0

        when [ 
            {attempts < 3} : {
                guess: input('Guess the password')
                guess == correct ? {
                    print("You got it right!")
                } {
                    guess_pasword(attempts + 1)
                }
            } 
            {attempts >= 3}: {
                print("To many attemps, account is locked")
            }
        ]
    }

```

in Yz v1.0 
```javascript
correct: 'pass123'
guess_password: {
    attempts: 0

    // Calls `match` passing the current attemps 
    // and a dictionary of blocks to test and actions
    
    match attempts [
        // 3 > attempts ? input
        3.> : {
            input 'Guess the password' == correct ? {
                print 'You got it right!'
            } {
                guess_password attempts + 1
            }
        }
        // 3 <= attempts
        3.<= : {
            print 'Too many attempts, account is locked'
        }
        
    ]

}

match: <T>{
    value T
    actions: [{T Bool}]{T}

    actions.for_each { key {T Bool} action {T}
    
        key value ? {
            return action()
        } 
    }
}
```

Possible definition of `match`  or `when`

v1.0
```js
when: <T>{
  conditions [{Bool}]{T} // dictionary with key block of Bool `{Bool}` and block of T `{T}` as valuej
  conditions.for_each { 
    cond   {Bool}
    action {T}
    cond() ? { return action()}
  }
  panic() // should match something, use 'true' for a default case
} 
```

[case/when](https://arturo-lang.io/playground/?example=conditional%20structures%20-%20case%20when) in action
```js
1.to 5 { it Int
    when [
        {it < 2} : {print "{it}: it's less than 2"}
        {it == 2}: {print "{it}: it's 2"}
        {it == 3}: {print "{it}: it's 3"}
        { true  }: {print "{it}: the number is too big"}
    ]
}
```

With match 
```javascript
1.to 5 { it Int
    match it [
        2.>= : {print "{it}: it's less than 2"}
        2.== : {print "{it}: it's 2"}
        3.== : {print "{it}: it's 3"}
        when.else:  {print "{it}: the number is too big"}
    ]
}
```


```js
when_equals animal [
    {'cat'}: { 'meow'}
    {'dog'}: { 'woof'}
    {'mouse'}: { 'sequeak'}
]

// with match
match animal [
    'cat'.== : {'meow'}
    'dog'.== : {'woof'}
    'mouse'.== : {'sequak'}
] {'shrug'}


// with match (named params)
match(subject:animal cond:[
    'cat'.== : {'meow'}
    'dog'.== : {'woof'}
    'mouse'.== : {'sequak'}
] else:{'shrug'})

// when_equals implementation

when_equals: <K,V>{ subject K  comps [{K}]{V}

    comps.for_each { k K v V
        subject == k() ? {return v()}
    }
}
```


TODO: verify with [Pattern Matching](Questions/solved/Pattern%20Matching.md)