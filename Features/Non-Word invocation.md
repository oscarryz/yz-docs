#feature 

When boc name is non-word, we can invoke it without `.` and `parenthesis` as long as it has at lest one parameter.

```js
Example: { 
   // the "<<" variable
   << : { 
	   n Int
	   printnln(n)
   }
}
e : Example()
e << 1 // same as e.<<(1)
```
