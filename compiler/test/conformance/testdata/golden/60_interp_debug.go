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

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) call() std.Unit {
	var p *Person = NewPerson(std.NewString("Alice"), std.NewInt(30))
	std.Print(std.NewString("Debug: ").Plus(std.NewString(std.Stringify(p))))
	std.Print(std.NewString("Name: ").Plus(std.NewString(std.Stringify(p.name))))
	std.Print(std.NewString("Expr: ").Plus(std.NewString(std.Stringify(std.NewInt(1).Plus(std.NewInt(2))))))
	return std.TheUnit
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		return self.call()
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
