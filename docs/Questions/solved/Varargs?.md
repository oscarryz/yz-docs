#open-question  How to handle varargs? 

```javascript

x {...Int} = {1 2 3 4 5} 
a b c: x() // a => 3, b => 4, c => 5
x 5 6 7 // still returns 1 2 3 4 5 because they were not assigned

y { ... Int} = {a:1 b:2 c:3}
l m n: y 3 4 5 6 7  // l => 3 m => 4 n => 5

y.a = 3
y.b = 4
y.c = 5
6
7
l: y.a //3
m: y.b //4
n: y.c //5
```

#answered  No var args, use an array (or a tuples)