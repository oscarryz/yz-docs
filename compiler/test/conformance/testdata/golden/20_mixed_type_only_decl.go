package main

import std "yz/runtime/yzrt"

type Named struct {
	name std.String
	greet func() *std.Thunk[std.Unit]
}

func NewNamed(name std.String, greet func() *std.Thunk[std.Unit]) *Named {
	return &Named{
		name: name,
		greet: greet,
	}
}

func (self *Named) Greet() *std.Thunk[std.Unit] {
	return self.greet()
}

func main() {
}
