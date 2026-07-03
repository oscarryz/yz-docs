package main

import std "yz/runtime/rt"

type Token struct {
	std.Cown
	type_ std.String
	value std.String
}

func NewToken(type_ std.String, value std.String) *Token {
	return &Token{
		type_: type_,
		value: value,
	}
}

func (self *Token) String() string {
	return "Token(type_: " + std.StringifyRepr(self.type_) + ", value: " + std.StringifyRepr(self.value) + ")"
}

func (self *Token) Type_() std.String {
	return self.type_
}

func (self *Token) Value() std.String {
	return self.value
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	var t *Token = NewToken(std.NewString("keyword"), std.NewString("if"))
	std.Print(t.type_)
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
