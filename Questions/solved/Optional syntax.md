#answered 

Style guide:
- Parenthesis for word invocations
- Parenthesis optional for non-word invocations
- Prefer to omit semicolon but use it for two expressions in the same line
- Commas mandatory to separate arguments and, in the multiple assignment list



In the examples I interchange between using and not using parenthesis, commas, semicolons. 

Could it be that these punctuations symbols be optional? Could it be too much?


```js
factorial : {
   n Int
   if n == 0 { 
     1 
   } { 
     n * factorial n -1 
   }
}
// Full punctuation
factorial : {
   n Int;
   if (n == 0 { 
      1; 
   }, { 
      n * factorial( n - 1)  
   });
}
// Following the recommendations
factorial : {
	n Int
	if (n == 0, {
		1
	}, {
		n * factorial( n - 1)
	})
}
```
