package main

import std "yz/runtime/rt"

type Person struct {
	std.Cown
	name std.String
}

func NewPerson(name std.String) *Person {
	return &Person{
		name: name,
	}
}

func (self *Person) greet() std.Unit {
	return std.Print(self.name)
}

func (self *Person) Greet() std.Unit {
	return std.LazyUnit(std.Schedule(&self.Cown, func() std.Unit {
		return self.greet()
	}))
}

func (self *Person) label() std.String {
	return self.name
}

func (self *Person) Label() std.String {
	return std.LazyString(std.Schedule(&self.Cown, func() std.String {
		return self.label()
	}))
}
