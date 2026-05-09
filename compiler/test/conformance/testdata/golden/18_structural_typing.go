package main

import std "yz/runtime/rt"

type Greeter interface {
	greet() *std.Thunk[std.Unit]
}


type Person struct {
	name std.String
	secret std.String
}

func NewPerson(name std.String, secret std.String) *Person {
	return &Person{
		name: name,
		secret: secret,
	}
}

func (self *Person) Greet() *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		return std.Print(self.name)
	})
}

func greet_all(g Greeter) *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		return g.Greet().Force()
	})
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		std.Schedule(&self.Cown, func() std.Unit {
			var p *Person = NewPerson(std.NewString("Alice"), std.NewString("my secret"))
			_bg0.Go(func() any {
				return greet_all(p).Force()
			})
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
