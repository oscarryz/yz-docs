#answered 

There are a lot of good questions here: 
https://www.reddit.com/r/ProgrammingLanguages/comments/17g6ny1/comment/k6h0lo5/?context=3


When a generic method binds, what happens to subsequent invocations? 


```
a : {  data }
a('hola') // data is String
a(1) // compilation error
```

How can we create reusable things if they'll bind and will never let go?

For instance boolean `?`

```js
Boolean: {
	? :  {
		true_block: { true_actions  {data} }
		false_block: {false_actions  {data} }
	}
}
(1 == 1).? { 'hola' } // true.? is bound to String 
(2 == 2).? { 1 }      // true.? is not attempted to bind to `Int` 
```

