package main

import std "yz/runtime/yzrt"

func main() {
	var n std.Int = std.NewInt(0)
	for n.Lt(std.NewInt(3)).GoBool() {
		n = n.Plus(std.NewInt(1))
	}
	std.Print(n)
}
