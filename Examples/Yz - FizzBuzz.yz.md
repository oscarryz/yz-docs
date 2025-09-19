#example
[FizzBuzz.st] https://gist.github.com/oscarryz/d54dd569fea585ec008c6f20af2e97ec 

Possible v0.0.1
```javascript
s: ''

1.to 100 .do { i Int
  by5: i % 5 == 0
  by3: i % 3 == 0

  if  {by5}, {
    s = s + 'Fizz'
  } 
  if  {by3}, { s = s + 'Buzz'} 
  if  {by5 && {by3}} == false, {s = i to_string} 

  println s 
  s = ''
} 
```
Improved v0.0.1

```javascript
1.to 100 .do {i Int
  by5: i % 5 == 0
  by3: i % 3 == 0

  by5 ? { s = s + 'Fizz' }
  by3 ? { s = s + 'Buzz' }

   by5 || {by3}  == false ? { s = i.to_string   }
  println s 
  s = ''
  
} 
```

Ideal v0.1.0
```javascript
1 to 100  do { i Int 
  by5: i % 5 == 0
  by3: i % 3 == 0

  if {by5}, { s = s + 'Fizz'} 
  if {by3}, { s = s + 'Buzz'} 
  if {by5 || {by3}} == false, {s = i to_string} 

  println s 
  s = ''
} 

```

Ideal v0.1.0

```javascript

1 to 100 do { i Int 

	by5: i % 5 == 0
	by3: i % 3 == 0

	by5 ? { s = s + 'Fizz'}
	by3 ? { s = s + 'Buzz'}
	by5.or {by3}  == false ? { s = i.to_string()  }
	print s 
	s = ''


}
```

Yz v1.0.0
```javascript
1.to 10  { i Int
    s: ''
    when [
        { i % 5 == 0 }: {s = s + 'Fizz'}
        { i % 3 == 0 }: {s = s + 'Buzz'}
        { true }      : {s = s + i.to_string()  }
    ]
    print s    
}
```

Final v1.0 
```javascript
1.to(10).each { i Int 
    if i % 3 == 0 ? { print 'Fizz'}
    if i % 5 == 0 ? { print 'Buzz'}
    if i % 3 != 0 && {i % 5 != 0}? { print '{i}'}
    println()
}
```

Another version 

```javascript
else: {true}
fizzbuzz: {
        n Int
        m3: n % 3 == 0
        m5: n % 5 == 0

        r: when [
           { m3 && {m5} } : { 'FizzBuzz' }
           { m3 }: { 'Fizz' }
           { m5 } : { 'Buzz' }
            else : { ints.to_string(n) }
        ]
        print("`r`")
}
1.to(100).do fizzbuzz
```

Another one 
```js
else: {true}
fizz_buzz: {
  n Int
  m3: n % 3
  m5: n % 5
  print(match_eq [
    {{0, 0}} : {"FizzBuzz"}
    {{0, any(Int)}}: {"Fizz"}
    {{any(Int), 0}}: {"Buzz"}
    else: {'`n`'}
  ])
}
1.to(100).do(fizz_buzz)
```

With match


```js
fizz_buzz: {
  1.to(100).do({
    i Int
   print(
      match {
         i % 3 == 0 => "Fizz"
         continue
      },{ 
         i % 5 == 0 => "Buzz"
      }, {
         `i`
      }
    )
}
```