https://discuss.kotlinlang.org/t/confusion-about-whether-non-local-returns-work-or-not/2408


```js
bar #(Int) {
   list : [1,2,3,4]
   list.each({
       it Int
       match { 
          it >= 2 => it
	  }
	  println(it) 
   })
   println("must not get here")
   -1
}
main: {
	num : bar()
	println(num)
}
```