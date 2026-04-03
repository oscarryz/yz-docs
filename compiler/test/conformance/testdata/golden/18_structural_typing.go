package main

import std "yz/runtime/yzrt"

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

func (self *Person) greet() *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		return std.Print(self.name)
	})
}

func greet_all(g Greeter) *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		return g.greet().Force()
	})
}

func main() {
	var p *Person = NewPerson(std.NewString("Alice"), std.NewString("my secret"))
	_bg0 := &std.BocGroup{}
	_bg0.Go(func() any {
		return greet_all(p).Force()
	})
	_bg0.Wait()
}
