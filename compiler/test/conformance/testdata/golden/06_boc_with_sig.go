package main

import std "yz/runtime/yzrt"

func greet(name std.String) *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		return std.Print(name)
	})
}

func main() {
	_bg := &std.BocGroup{}
	_bg.Go(func() any {
		return greet(std.NewString("Alice")).Force()
	})
	_bg.Go(func() any {
		return greet(std.NewString("Bob")).Force()
	})
	_bg.Wait()
}
