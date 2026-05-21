package main

import std "yz/runtime/rt"

type Point struct {
	std.Cown
	x std.Int
	y std.Int
}


func (self *Point) String() string {
	return "Point(x: " + std.Stringify(self.x) + ", y: " + std.Stringify(self.y) + ")"
}
