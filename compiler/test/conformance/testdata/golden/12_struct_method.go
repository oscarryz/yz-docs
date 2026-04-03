package main

import std "yz/runtime/yzrt"

type Named struct {
	name std.String
}

func NewNamed(name std.String) *Named {
	return &Named{
		name: name,
	}
}

func (self *Named) hi() *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		return std.Print(self.name)
	})
}
