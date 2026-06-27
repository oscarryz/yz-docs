package main

import std "yz/runtime/rt"

type _utilsExtraBoc struct {
	std.Cown
}

func (self *_utilsExtraBoc) String() string {
	return "{ " + "help: {}" + " }"
}

func (self *_utilsExtraBoc) help() std.String {
	return std.NewString("extra help")
}

func (self *_utilsExtraBoc) Help() std.String {
	return std.LazyString(std.Schedule(&self.Cown, func() std.String {
		return self.help()
	}))
}


type _utilsBoc struct {
	std.Cown
	extra *_utilsExtraBoc
}

func (self *_utilsBoc) String() string {
	return "{ " + "extra: " + std.StringifyRepr(self.extra) + " }"
}

var Utils = &_utilsBoc{
	extra: &_utilsExtraBoc{},
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	std.Print(Utils.extra.Help())
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
