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

type _greetBoc struct {
	std.Cown
	name std.String
	greeting std.String
}

func (self *_greetBoc) String() string {
	return "{ " + "name: " + std.StringifyRepr(self.name) + "; " + "greeting: " + std.StringifyRepr(self.greeting) + "; " + "call: {}" + " }"
}

func (self *_greetBoc) Call(name std.String, greeting std.String) *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		self.name = name
		self.greeting = greeting
		std.Print(self.greeting)
		return std.Print(self.name)
	})
}

var Greet = &_greetBoc{
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		var p *Person
		std.Schedule(&self.Cown, func() std.Unit {
			p = NewPerson(std.NewString("Alice"), std.NewInt(30))
			std.Print(p.name)
			std.Print(std.NewString(std.StringifyRepr(p.age)))
			_st0 := Greet.Call(std.NewString("Bob"), std.NewString("Hello"))
			_bg0.Add(func() { _st0.Force() })
			_st1 := Greet.Call(std.NewString("Carol"), std.NewString("Hi"))
			_bg0.Add(func() { _st1.Force() })
			_st2 := Greet.Call(std.NewString("Dave"), std.NewString("Hello"))
			_bg0.Add(func() { _st2.Force() })
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
