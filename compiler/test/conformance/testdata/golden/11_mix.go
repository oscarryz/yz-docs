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
