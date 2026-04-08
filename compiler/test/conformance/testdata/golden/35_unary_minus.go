package main

import std "yz/runtime/yzrt"

func main() {
	var x std.Int = std.NewInt(5)
	var neg_x std.Int = x.Neg()
	std.Print(neg_x)
	var a std.Int = std.NewInt(10)
	var b std.Int = std.NewInt(3)
	var result std.Int = a.Minus(b.Neg())
	std.Print(result)
}
