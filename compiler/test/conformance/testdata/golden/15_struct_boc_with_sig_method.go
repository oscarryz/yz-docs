package main

import std "yz/runtime/rt"

type Person struct {
	name std.String
}

func NewPerson(name std.String) *Person {
	return &Person{
		name: name,
	}
}

func (self *Person) Greet() *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		return std.Print(self.name)
	})
}

func (self *Person) Label() *std.Thunk[std.String] {
	return std.Go(func() std.String {
		return self.name
	})
}
