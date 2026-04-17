package main

import std "yz/runtime/yzrt"

func main() {
	foo := func() *std.Thunk[std.String] {
		return std.Go(func() std.String {
			return std.NewString("hello")
		})
	}
	bar := func() *std.Thunk[std.String] {
		return std.Go(func() std.String {
			return std.NewString("world")
		})
	}
	std.Print(foo().Force())
	std.Print(bar().Force())
}
