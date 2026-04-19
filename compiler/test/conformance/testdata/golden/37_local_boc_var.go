package main

import std "yz/runtime/yzrt"

type _main_fooBoc struct {
}

func (self *_main_fooBoc) Call() *std.Thunk[std.String] {
	return std.Go(func() std.String {
		return std.NewString("hello")
	})
}


type _main_barBoc struct {
}

func (self *_main_barBoc) Call() *std.Thunk[std.String] {
	return std.Go(func() std.String {
		return std.NewString("world")
	})
}


func main() {
	_foo := &_main_fooBoc{}
	_bar := &_main_barBoc{}
	std.Print(_foo.Call().Force())
	std.Print(_bar.Call().Force())
}
