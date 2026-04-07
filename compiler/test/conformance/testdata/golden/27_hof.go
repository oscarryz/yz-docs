package main

import std "yz/runtime/yzrt"

func main() {
	var list std.Array[std.Int] = std.NewArray(std.NewInt(1), std.NewInt(2), std.NewInt(3), std.NewInt(10), std.NewInt(20))
	var filtered std.Array[std.Int] = list.Filter(func(item std.Int) std.Bool {
		return item.Gt(std.NewInt(10))
	})
	filtered.Each(func(item std.Int) std.Unit {
		return std.Print(item)
	})
}
