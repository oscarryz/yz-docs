package main

import std "yz/runtime/yzrt"

func main() {
	a := std.Http.Get(std.NewString("https://httpbin.org/get"))
	b := std.Http.Get(std.NewString("https://httpbin.org/uuid"))
	std.Print(a.Force())
	std.Print(b.Force())
}
