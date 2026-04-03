package main

import std "yz/runtime/yzrt"

func main() {
	var x std.Int = std.NewInt(5)
	if x.Gt(std.NewInt(3)).GoBool() {
		std.Print(std.NewString("big"))
	} else {
		std.Print(std.NewString("small"))
	}
}
