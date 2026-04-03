package main

import std "yz/runtime/yzrt"

func main() {
	var score std.Int = std.NewInt(85)
	var grade std.String = func() std.String {
		if score.Gteq(std.NewInt(90)).GoBool() {
			return std.NewString("A")
		} else if score.Gteq(std.NewInt(80)).GoBool() {
			return std.NewString("B")
		} else if score.Gteq(std.NewInt(70)).GoBool() {
			return std.NewString("C")
		} else {
			return std.NewString("F")
		}
	}()
	std.Print(grade)
}
