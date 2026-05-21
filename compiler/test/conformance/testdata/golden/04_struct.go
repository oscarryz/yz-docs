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
	return "Person(name: " + std.Stringify(self.name) + ", age: " + std.Stringify(self.age) + ")"
}
