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
	std.Print(std.NewString(std.StringifyRepr(a.Plus(b))))
	std.Print(std.NewString(std.StringifyRepr(a.Minus(b))))
	std.Print(std.NewString(std.StringifyRepr(a.Star(b))))
	std.Print(std.NewString(std.StringifyRepr(a.Slash(b))))
	std.Print(std.NewString(std.StringifyRepr(a.Neg())))
	std.Print(std.NewString(std.StringifyRepr(a.Eqeq(b))))
	std.Print(std.NewString(std.StringifyRepr(a.Neq(b))))
	std.Print(std.NewString(std.StringifyRepr(a.Lt(b))))
	std.Print(std.NewString(std.StringifyRepr(a.Gt(b))))
	std.Print(std.NewString(std.StringifyRepr(a.Lteq(b))))
	std.Print(std.NewString(std.StringifyRepr(a.Gteq(b))))
	var x std.Decimal = std.NewDecimal(3.14)
	std.Print(std.NewString(std.StringifyRepr(x.Abs())))
	std.Print(std.NewString(std.StringifyRepr(x.Pow(std.NewDecimal(2)))))
	std.Print(x.ToStr())
	std.Print(std.NewString(std.StringifyRepr(std.NewDecimal(5).Slash(std.NewDecimal(2)))))
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
