package main

import std "yz/runtime/rt"

type Named struct {
	std.Cown
	name std.String
}

func NewNamed(name std.String) *Named {
	return &Named{
		name: name,
	}
}

func (self *Named) Hi() *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		return std.Print(self.name)
	})
}
