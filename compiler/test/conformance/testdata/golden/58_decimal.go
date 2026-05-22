package main

import std "yz/runtime/rt"

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	var a std.Decimal = std.NewDecimal(10)
	var b std.Decimal = std.NewDecimal(3)
	std.Print(a.Plus(b))
	std.Print(a.Minus(b))
	std.Print(a.Star(b))
	std.Print(a.Slash(b))
	std.Print(a.Neg())
	std.Print(a.Eqeq(b))
	std.Print(a.Neq(b))
	std.Print(a.Lt(b))
	std.Print(a.Gt(b))
	std.Print(a.Lteq(b))
	std.Print(a.Gteq(b))
	var x std.Decimal = std.NewDecimal(3.14)
	std.Print(x.Abs())
	std.Print(x.Pow(std.NewDecimal(2)))
	std.Print(x.ToStr())
	std.Print(std.NewDecimal(5).Slash(std.NewDecimal(2)))
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
