#example

#concurrency 

[A boring function](https://go.dev/talks/2012/concurrency.slide#12)
```js
boring: {
  msg String
  i: 0
  while({true}: {
    print("`msg`, `i`")
    time.delay(random.int(3))
    i = i + 1
  }
}
main: {
  boring("boring!") // launches
  // and waits at the end of the main block because of structural concurrency
}
```


[Channels](https://go.dev/talks/2012/concurrency.slide#19)

There are no channels, instead use a boc to communicate.

To avoid race ahead the following example uses an array to store the generated messages, then retrieves them in order by calling `messages.pop`

```js
boring: {
    m String

    messages: [String]()
    i: 1
    
    while({ true }, {
        
        messages.push("`m` `i`")    // acquires {messages}, writes, releases
        i = i + 1
        time.delay(1)
    })                      // ← natural yield point here!
}

main: {
    boring("sync")          // launches boring, doesn't block
    5.times().do({
        print(boring.messages.pop()) // acquires {messages}, reads, releases
    })
}
```

Because the BoC model, the boc in the `while` and the boc in the `do` loop have to gain access to the `messages` boc, and due to the _happens-before_ trait, the first message is written before the first read takes place.


[Generator](https://go.dev/talks/2012/concurrency.slide#24)

Similar principle, share the messages queue, but start the loop from outside.

The boring last two expressions are bocs, one has an infinite loop writing to the array, the other extracts an element from the array 

```js
boring: {
  s String
  messages: [String]()
  
  {
    while({true}, {
      messages.push("`s` `i`")
      i = i + 1
      time.delay(1) 
    })
  }
  { messages.pop() }
}
main: {
  // assigning the bocs
  msg, l : boring("sync")
  l() // start the loop
  5.times().do({
    print("you said `msg()`")
  })
}
```

Now each can handle a service

```js
main:{
   joe, jl: boring("joe")
   ann, al: boring("ann")
   jl()
   al()
   5.times().do({
      print(joe())
      print(ann())
   })
}
```

[Multiplexing][26]

```js
fan_in : {
   a #(String)
   launcher_a #()
   b #(String)
   launcher_b #()

   messages: [String]()
   
   read: { src #(String

   {{
       launcher_a()
       while_true({
         messages.push(a())
       })
   }()
   {
       launcher_b()
       while_true({
         messages.push(b())
       })
   }()}
   { messages.pop() }
}
main:{
   m, l = fan_in(boring("joe"), boring("ann"))
   l()
   while_true(m)
}
```

26: https://go.dev/talks/2012/concurrency.slide#26