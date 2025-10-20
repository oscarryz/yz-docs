
https://github.com/octalide/mach?tab=readme-ov-file#simple-examples

```js

// hello.yz
main : {
   println("Hello, World!")
}

// fib.yz
fib #(n Int, Int) {
   match { 
	   n < 2 => n
   }, { 
	   fib(n - 1) + fib( n - 2)
   } 
}
main: {
    max : 10 
    println("`fib(max)`")
}
// fac.yz
factorial #(n Int, Int) {
   match {
	   n == 0 => 1
   }, {
	   n * factorial( n - 1)
   }
}
main: {
	max : 10 
	println("`fact(max)`")
}


```