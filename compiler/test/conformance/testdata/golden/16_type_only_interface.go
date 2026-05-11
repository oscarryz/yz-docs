package main

import std "yz/runtime/rt"

type Greeter interface {
	greet() *std.Thunk[std.Unit]
}


type Person struct {
	std.Cown
	name std.String
}

func NewPerson(name std.String) *Person {
	return &Person{
		name: name,
	}
}

func (self *Person) Greet() *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		return std.Print(self.name)
	})
}
