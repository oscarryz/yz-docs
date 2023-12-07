Crazy idea, but how would it be like to use significan whitespace (like Python) to demark blocks

#rejected



```js
a: 1 
b: 2
c: a + b
d: 
    print "I'm a block"
d() // I'm a block
e: [1,2,3,4,5,6]
// full syntax
f: e.filter({item Int; item > 3}) // f:[4,5,6]
g: e.filter {item Int; item > 3}  // g:[4,5,6]

h: e.filter 
    item Int; item > 3  // h:[4,5,6]

// my original example
main:  
    names: io.args
    names.for_each
        name String
        print 'Hello {name}'
 

```

With generics using `+`

```js
compare: 
    a +T; b +T;
    result: 0
    a < b ? 
        result = -1
    a > b ?
        result = 1
    result 

// compare<Integer>
println compare 1,0  
// compare<String>
println compare 'a', 'b'  
// compare<Decimal>
println()

```

```js
hierarchy:  
    Schema:  
        hello String
        life Int
        pi Float
        nums [] Int
        struct   
     
    Constrained: 
        '[cue-constraint:=~"[a-z]+"]'
        hello String
        
        '[cue-constraint:>0]'
        life Int
        
        '[cue-constraint: list.MaxItems(10)]'
        nums []Int
 
     
    value: Constrained(
        hello: 'world',
        life: 42,
        pi: 3.14,
        nums: [1,2,3,4,5],
        struct:  
            a: 'a'
            b: 'b'
    )
        
         
     
 
```

```js
fs:  
    create:   
        path String
        // do something
        error_handler ErrorHandler
        file File = ... 
        file.noaccess ?  error_handler.handle(file) 
     
    ErrorHandler:  
       scenarios [File:Boolean] 
       handle:   
           file File
           scenarios.keys().for_each( k  File,Boolean 
                k(file).?( 
                    scenarios[k]();
                 )
            )
        
     
 
```

```
Person:  
  name String
  speak:   
    "Hi, my name is  name "
   
 
Computer:  
  model String
  speak:  
    "[ model ]: Beep bop"
   
 
do_speak:   speaker  name String; speak  String  
  print " speaker.speak() "
 
do_speak Person("Oscar")    // "Hi, my name is Oscar
do_speak Person("ABC-123") // [ABC-123]: Beep bop
do_speak 
    name:"Yz"; speak:  "Yup"// Yup


```

```js

input : `Lead Chef, Chipotle, Denver, CO, 10, 15

Stunt Double, Equity, Los Angeles, CA, 15, 25

Manager of Fun, IBM, Albany, NY, 30, 40

Associate Tattoo Artist, Tit 4 Tat, Brooklyn, NY, 250, 275

Assistant to the Regional Manager, IBM, Scranton, PA, 10, 15

Lead Guitarist, Philharmonic, Woodstock, NY, 100, 200

--JSON FORMAT BELOW--

{"name": "Spaceship Repairman", "location": {"city": "Olympus Mons", "state": "Mars"}, "organization": "Interplanetary Enterprises", "pay": {"min": 100, "max": 200}}

{"name": "State Park Ranger", "location": {"city": "Ray Brook", "state": "NY"}, "organization": "Adirondack Park Agency", "pay": {"min": 40, "max": 50}}

{"name": "Lead Cephalopod Caretaker", "location": {"city": "Atlantis", "state": "Oceania"}, "organization": "Deep Adventures", "pay": {"min": 10, "max": 15}}`

parse:   

    input String 

    a: input.split('\n')
    strategy: from_line
    a.for_each   
        line String 
        line == '--JSON FORMAT BELOW--' ?  
            strategy = from_json
         ,  
            strategy(line)
    a
 
from_line:  
    s String

    a: s.split(',')
    return  
        title: parts[0].trim(),
        organization: parts[1].trim(),
        location:  
            city: parts[2].trim(),
            state: parts[3].trim(),
         ,
        pay:  
            min: parts[4].trim(),
            max: parts[5].trim()
         
     
                
from_json:  
    json_string String

    o: json.parse(json_string)

    return  
      title: o.name,
      organization: o.organization,
      location:  
          city: o.city,
          state: o.state,
       ,
      pay:  
          min: o.pay.min,
          max: o.pay.max,
       
     
 
lookup_table [String][]Opportunities = []

filter_by_state:  
    state String               
    lookup_table[state]
 
init_lookup:  
    array []Opportunity 
    lookup_table = []
    array.for_each  
        o Opportinity
        entries: lookup_table[o.location.state]
        arrays.is_empty(entries) ?  
            entries = []
            lookup_table[state] = entries
         
        entries.add(o)

     
 

print_opportunities:  
    array []Opportunity; banner: 'All Opportunities'
    array.sort   
        a Opportunity; b Opportunity
        a.title.compare(b.title)
     
    print ' banner '
    array.for_each  
        o Opportunity

        print  `Title:  {o.title}
               Organization:  {o.organization}
               Location:  {o.location.city} {o.location.state}
               Pay:  {o.pay.min} -$ {o.pay.max}`
     
 

Opportunity:  
    title String
    organization String
    location  String, String 
    pay  {Int,Int} 
 

opportunities: parse(input)
init_lookup(opportunities)
print_opportunities(opportunities)
println()
prit_opportunities(filter_by_state(opportunities, 'NY'), 'NY Opportunities')




```

```js
1.to(100).each
  i Int 

  i % 5 == 0 ? 
      s = s + 'Fizz'
  i % 3 == 0 ? 
      s = s + 'Buzz'     

  (i % 5 == 0) || ({i % 3 == 0}) == false ?
      s = i to_string() 
  
  (i % 5 == 0) || ({i % 3 == 0})  == false ?
      s = i to_string() 

  println(s)
  s = ''
    
```


```js
HelloWorldModel: {
  saying String
}
model: HelloWorldModel(saying:'Hello world')
// or 
model: HelloWorldModel { 'Hello world' }
// maybe but maybe not: 
model: HelloWorldModel('Hello World') // maybe not because it creates an instance and not returns the last valu

win: Frame (
    title: '{model.saying} JavaFX'
    width: 200
    content: TextField (
        value: model.saying
    )
    visible: true
)
```


```js
    User:  
        username String
        name String
        last_name String
     
    
    o: User(name: 'Oscar', last_name: 'Reyes') 
    o2: User(username: 'oscarryz')
    o3: User()
    
    greet:
        user User
        match(o,
            {name:_, last_name:_}, {print('Hello {o.name} {o.last_name}')}
            {username:_}, {print('Hello {o.username}')}
            {}, {print('Hello new user')}
        )
  

    yz>greet(o)
    Hello Oscar Reyes
    
    yz>greet(o2)
    Hello oscarryz

    yz>greet(o3)
    Hello new user
```

```js

Vector: 
    x: 0
    y: 0
    z: 0
    add: 
        v Vector
        
```
#answered : No, looks cool but main design already considers braces.
