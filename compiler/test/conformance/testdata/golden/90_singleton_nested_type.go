package main

import std "yz/runtime/rt"

type _roomwindow struct {
	std.Cown
	size std.Int
}

func New_roomwindow(size std.Int) *_roomwindow {
	return &_roomwindow{
		size: size,
	}
}

func (self *_roomwindow) String() string {
	return "_roomwindow(size: " + std.StringifyRepr(self.size) + ")"
}

type _roomBoc struct {
	std.Cown
}

func (self *_roomBoc) String() string {
	return "{ }"
}

var Room = &_roomBoc{}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	var w *_roomwindow = New_roomwindow(std.NewInt(3))
	std.Print(w.size)
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
