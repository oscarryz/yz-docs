When block name is non-word and has at least one parameter

```js

    Array: <T>{
    
        'Add element to array...'
        << : {
            element T
        }
        add: <<
        push: <<

        ++ : {
            other []T
        }
        concatenate: ++

        ! : {
            print 'hello'
        }
        
    }

    a: [1 2 3]
    b: [4 5 6]

    a << 7 // a.<<(7)
    a.push 7 
    a.add 7 
    
    a ++ b // a.++(b)
    a.concatenate b 

    a ! // invalid
    a.!() // ok 
```
