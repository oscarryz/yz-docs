#feature
By default the blocks return all the last values computed. It is possible to return earlier by using the `return` keyword

```js
check: {
   age Int
   age < 21 ? {
       message: 'You have to be over 21 to access this site'
	   return
   }
   message: 'Welcome to the site'
}
m: check(20)
print m // the last computed value was the `age <21` which in turn returned `message: 'You have to be... '`
```

