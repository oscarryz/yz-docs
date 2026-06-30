package main

import std "yz/runtime/rt"

type Person struct {
	std.Cown
	name std.String
	age std.Int
}

func NewPerson(name std.String, age std.Int) *Person {
	return &Person{
		name: name,
		age: age,
	}
}

func (self *Person) String() string {
	return "Person(name: " + std.StringifyRepr(self.name) + ", age: " + std.StringifyRepr(self.age) + ")"
}

func (self *Person) Name() std.String {
	return self.name
}

func (self *Person) Age() std.Int {
	return self.age
}
