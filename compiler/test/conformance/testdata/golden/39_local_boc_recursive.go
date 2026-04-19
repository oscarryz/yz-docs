package main

import std "yz/runtime/yzrt"

func main() {
	var f any = func(n std.Int) std.Unit {
		if n.Eqeq(std.NewInt(0)).GoBool() {
			std.Print(std.NewString("fin"))
		} else {
			std.Print(n)
			std.Go(func() std.Unit {
				return f(n.Minus(std.NewInt(1)))
			})
		}
		return std.TheUnit
	}
	std.Go(func() std.Unit {
		return f(std.NewInt(3))
	})
}
