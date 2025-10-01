#example
[FizzBuzz.st] https://gist.github.com/oscarryz/d54dd569fea585ec008c6f20af2e97ec 

Possible v0.0.1
```js
s: ''

1.to(100).each({ i Int
  by5: i % 5 == 0
  by3: i % 3 == 0

  by5 ? { s = s + 'Fizz' }
  by3 ? { s = s + 'Buzz' }
  (by5 && by3) == false ? { s = i.to_string() } 

  println s 
  s = ''
} 
```
Improved v0.0.1

```js
1.to(100).each({ i Int
  by5: i % 5 == 0
  by3: i % 3 == 0

  by5 ? { s = s + 'Fizz' }
  by3 ? { s = s + 'Buzz' }

   (by5 || by3) == false ? { s = i.to_string() }
  println s 
  s = ''
  
} 
```

Ideal v0.1.0
```js
1.to(100).each({ i Int 
  by5: i % 5 == 0
  by3: i % 3 == 0

  by5 ? { s = s + 'Fizz' }
  by3 ? { s = s + 'Buzz' }
  (by5 || by3) == false ? { s = i.to_string() } 

  println s 
  s = ''
} 

```

Ideal v0.1.0

```js

1.to(100).each({ i Int 

	by5: i % 5 == 0
	by3: i % 3 == 0

	by5 ? { s = s + 'Fizz'}
	by3 ? { s = s + 'Buzz'}
	(by5 || by3) == false ? { s = i.to_string() }
	print s 
	s = ''


}
```

Yz v1.0.0
```js
1.to(10).each({ i Int
    s: ''
    match {
        i % 5 == 0 => s = s + 'Fizz'
    }, {
        i % 3 == 0 => s = s + 'Buzz'
    }, {
        s = s + i.to_string()
    }
    print s    
}
```

Final v1.0 
```js
1.to(10).each({ i Int 
    i % 3 == 0 ? { print('Fizz') }
    i % 5 == 0 ? { print('Buzz') }
    (i % 3 != 0 && i % 5 != 0) ? { print('`i`') }
    println()
}
```

Another version 

```js
else: {true}
fizzbuzz: {
        n Int
        m3: n % 3 == 0
        m5: n % 5 == 0

        r: match {
           (m3 && m5) => 'FizzBuzz'
        }, {
           m3 => 'Fizz'
        }, {
           m5 => 'Buzz'
        }, {
           n.to_string()
        }
        print("`r`")
}
1.to(100).each({ i Int; fizzbuzz(i) })
```

Another one 
```js
else: {true}
fizz_buzz: {
  n Int
  m3: n % 3
  m5: n % 5
  print(match {
    (m3 == 0 && m5 == 0) => "FizzBuzz"
  }, {
    m3 == 0 => "Fizz"
  }, {
    m5 == 0 => "Buzz"
  }, {
    "`n`"
  })
}
1.to(100).each({ i Int; fizz_buzz(i) })
```

With match


```js
1.to(100).each({
  i Int
  print(fizz_buzz(i))
})

fizz_buzz #(i Int) {
  match {
   i % 3 == 0 => "Fizz"
  }, { 
    i % 5 == 0 => "Buzz"
  }, {
   "`i`"
  }
}
```