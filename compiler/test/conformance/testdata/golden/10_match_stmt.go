package main

import std "yz/runtime/yzrt"

func main() {
	var x std.Int = std.NewInt(0)
	if x.Gt(std.NewInt(0)).GoBool() {
		std.Print(std.NewString("positive"))
	} else if x.Lt(std.NewInt(0)).GoBool() {
		std.Print(std.NewString("negative"))
	} else {
		std.Print(std.NewString("zero"))
	}
}
