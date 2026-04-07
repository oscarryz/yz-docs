package main

import std "yz/runtime/yzrt"

func main() {
	var list std.Array[std.Int] = std.NewArray(std.NewInt(1), std.NewInt(2), std.NewInt(3))
	doubled := std.ArrayMap(list, func(item std.Int) std.Int {
		return item.Star(std.NewInt(2))
	})
	doubled.Each(func(item std.Int) std.Unit {
		return std.Print(item)
	})
}
