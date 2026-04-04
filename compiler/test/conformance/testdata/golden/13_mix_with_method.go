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

func (self *Named) Hi() *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
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
