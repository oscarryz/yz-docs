https://doc.rust-lang.org/rust-by-example/custom_types/enum.html


```js
WebEvent : { 

	PageLoad(),
	PageUnload(),
	KeyPress(c String),
	Paste(s String),
	Click(#(x Int, y Int))
	
}

inspect #(even WebEvent) {
   match event {
	   WebEvent.PageLoad => println("page loaded")
   }, {
	   WebEvent.PageUnload => println("page unloaded")
   }, {
	   WebEvent.KeyPress(String) => println("pressed `event.c`.")
   }, {
	   WebEvent.Paste(String) => println("pasted \"`event.s`\"")
   }, {
	   WebEvent.Click(#(Int,Int)) => println("clicked at `event.x`, `event.y`.")
   }
}
main: {
	pressed : WebEvent.KeyPress("x")
	pasted : WebEvent.Parse("my text")
	click : WebEvent.Click({x: 20, y: 80})
	load : WebEvent.PageLoad()
	unload : WebEvent.PageUnload()
	
	inspect(pressed)
	inspect(pasted)
	inspect(click)
	inspect(load)
	inspect(unload)
}

```
