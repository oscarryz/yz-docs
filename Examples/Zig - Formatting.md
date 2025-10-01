#example
[From zig](https://ziglearn.org/chapter-2/#formatting:~:text=over%20string%20printing.-,const%20Person%20%3D%20struct%20%7B,%7D,-JSON)

```js
Person: {
    name String
    birth_year Int
    death_year Int
    format: {
        dys: death_year != int.nan?{
                '{death_year}'
            } {
                ''
            }
         '{name} ({birth_year}-{dys})'    
    }
}
test 'custom fmt' {
    claude: Person {
       name: "Claude Shannon"
       birth_year: 1916
       death_year: 2002
    }
    claude_string: claude.format()
    assert equals claude_string "Claude Shannong (1916-2001)"
}
```