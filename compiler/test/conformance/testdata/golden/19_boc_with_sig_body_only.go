package main

import std "yz/runtime/yzrt"

func greet(name std.String) *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		return std.Print(name)
	})
}

func shout(msg std.String) *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		return std.Print(msg)
	})
}

func main() {
	_bg0 := &std.BocGroup{}
	_bg0.Go(func() any {
		return greet(std.NewString("Alice")).Force()
	})
	_bg0.Go(func() any {
		return shout(std.NewString("hello")).Force()
	})
	_bg0.Wait()
}
