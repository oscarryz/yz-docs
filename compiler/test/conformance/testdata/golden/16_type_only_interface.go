package main

import std "yz/runtime/rt"

type Greeter interface {
	greet() std.Unit
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

func (self *Person) greet() std.Unit {
	return std.Print(self.name)
}

func (self *Person) Greet() std.Unit {
	return std.LazyUnit(std.Schedule(&self.Cown, func() std.Unit {
		return self.greet()
	}))
}
