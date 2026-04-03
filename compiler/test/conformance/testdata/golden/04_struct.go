package main

import std "yz/runtime/yzrt"

type Person struct {
	name std.String
	age std.Int
}

func NewPerson(name std.String, age std.Int) *Person {
	return &Person{
		name: name,
		age: age,
	}
}
