package main

import std "yz/runtime/rt"

type Named struct {
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

type Person struct {
	Named
	last_name std.String
}

func NewPerson(name std.String, last_name std.String) *Person {
	return &Person{
		Named: *NewNamed(name),
		last_name: last_name,
	}
}
