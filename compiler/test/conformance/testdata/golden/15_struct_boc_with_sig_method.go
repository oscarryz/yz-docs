package main

import std "yz/runtime/yzrt"

type Person struct {
	name std.String
}

func NewPerson(name std.String) *Person {
	return &Person{
		name: name,
	}
}

func (self *Person) greet() *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		return std.Print(self.name)
	})
}

func (self *Person) label() *std.Thunk[std.String] {
	return std.Go(func() std.String {
		return self.name
	})
}
