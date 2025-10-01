#example


From [Crystal HTTP Server](https://crystal-lang.org/reference/1.5/getting_started/http_server.html)

```js
server: http.Server(
    get: {
        context http.Context
        context.response.context_type = 'text/plain'
        context.response.print('Hello world!  The time is `time.now()`')
    }
)
address: server.bind(8080)
print('Listening on https://`adress`')
server.listen()

```

```js
the_beatles: [
    "John Lennon",
    "Paul McCartney",
    "George Harrison",
    "Ringo Starr",
 ]
shout: false
say_hi_to: ''
strawberry: false

option_parser: optionParser.parse({
    parser Parser
    parser.banner: 'Welcome to the Beatles app'

    parser.on(['-v', '--version', 'Show version'], {
        print('version 1.0')
        return
    })

    parser.on ['-h'  '--help'  'Show help']  {
        print('`parser`')
        return
    }

    parser.on ['-t'  '--twist'  'Twist and SHOUT']  {
        shout = true
        return
    }

    parser.on ['-g Name'  '--goodbye_hello=NAME'  'Say hello to who ever you want']  {
        say_hi_to = name
    }

    parser.on ['-r'  '--random_goobye_hello'  'Say hello to one random member']  {
        print(the_beatles[random.Int(the_beatles.len))]
        return
    }

    parser.on ['-s'  '--strawberry'  'Strawberry fields forever mode ON']  {
        strawberry = true
    }

    parser.missing_option {
        option_flag String
        print(io.std_err('Error `option_flag` is missing something'))
        print(io.std_err(""))
        print(io.std_err('{parser}'))
        return
    }

    parser.invalid_option { option_flag String
        print(io.std_err  "Error: `option_flag` is not a valid option")
        print(io.std_err  '`parser`')
        return
    }
})

members: the_beatles
shout ? {
    members = members.map(strings.to_upper_case)
}
strawberry ? {
    print('Streawberry fields ferever mode ON')
}

print(`
    Group Members
    =============
`)
member.for_each({
    member String
    print('`strawberry ? {🍓} {-}` `member`')
})
say_hi.empty() == false ? {
    print("You say goodbye  and I say hello to `say_hi_to`")
}

```


Using a map as configuration:
```js
the_beatles: [
    "John Lennon",
    "Paul McCartney",
    "George Harrison",
    "Ringo Starr"
 ]
shout: false
say_hi_to: ''
strawberry: false
parser: OptionParser(
    configure: [
        ['-v',  '--version', 'Show version']: {
            print('version 1.0')
        },
        ['-h'  '--help'  'Show help']: {
             print('`parser`')
        },
        ['-t'  '--twist'  'Twist and SHOUT']: {
            shout = true
        }
        ['-g Name'  '--goodbye_hello=NAME'  'Say hello to who ever you want']: {
            say_hi_to = name
        },
        ['-r'  '--random_goobye_hello'  'Say hello to one random member']: {
            print(the_beatles[random.Int(the_beatles.len)])
        }
        ['-s'  '--strawberry'  'Strawberry fields forever mode ON']: {
            strawberry = true
        }
    ]
)


     missing_option: { option_flag String
        printerr 'Error {option_flag} is missing something'
        printerr ""
        printerr '{parser}'
    }
    invalid_option: { option_flag String
        printerr  "Error: {option_flag} is not a valid option"
        printerr  '{parser}'
    }
)
members: the_beatles
shout ? {
    members = members.map(strings.to_upper_case)
}
strawberry ? {
    print('Streawberry fields ferever mode ON')
}

print(`
    Group Members
    =============
`)
member.for_each({ member String
    print('{strawberry?{🍓} {-}} {member}')
})
say_hi.empty() == false ?({
    print(`)
    You say goodbye  and I say hello to {say_hi_to}
    `
})

```
