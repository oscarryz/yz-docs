package main

import std "yz/runtime/yzrt"

func main() {
	var name std.String = std.NewString("World")
	var n std.Int = std.NewInt(42)
	std.Print(std.NewString("Hello, ").Plus(std.NewString(std.Stringify(name))).Plus(std.NewString("!")))
	std.Print(std.NewString("Answer: ").Plus(std.NewString(std.Stringify(n))))
	std.Print(std.NewString("Sum: ").Plus(std.NewString(std.Stringify(n.Plus(std.NewInt(1))))))
}
