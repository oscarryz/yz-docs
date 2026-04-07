package main

import std "yz/runtime/yzrt"

func identity[V any](value V) *std.Thunk[V] {
	return std.Go(func() V {
		return value
	})
}

func main() {
	x := identity(std.NewString("hello"))
	std.Print(x.Force())
}
