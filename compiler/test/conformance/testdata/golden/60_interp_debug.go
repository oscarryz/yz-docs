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

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	var p *Person = NewPerson(std.NewString("Alice"), std.NewInt(30))
	std.Print(std.NewString("Debug: ").Plus(std.NewString(std.StringifyRepr(p))))
	std.Print(std.NewString("Name: ").Plus(std.NewString(std.StringifyRepr(p.name))))
	std.Print(std.NewString("Expr: ").Plus(std.NewString(std.StringifyRepr(std.NewInt(1).Plus(std.NewInt(2))))))
	return std.TheUnit
}

func (self *_mainBoc) Call() std.Unit {
	return std.LazyUnit(std.Schedule(&self.Cown, func() std.Unit {
		return self.call()
	}))
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
