package main

import std "yz/runtime/rt"

type _GraphNodeBound interface {
	Label() *std.Thunk[std.String]
}


type Graph interface {
}


type Point struct {
	std.Cown
	x std.Int
	y std.Int
}

func NewPoint(x std.Int, y std.Int) *Point {
	return &Point{
		x: x,
		y: y,
	}
}

func (self *Point) String() string {
	return "Point(x: " + std.StringifyRepr(self.x) + ", y: " + std.StringifyRepr(self.y) + ")"
}

func (self *Point) label() std.String {
	return self.x.ToStr().Plus(std.NewString(",")).Plus(self.y.ToStr())
}

func (self *Point) Label() *std.Thunk[std.String] {
	return std.Schedule(&self.Cown, func() std.String {
		return self.label()
	})
}

type PointGraph struct {
	std.Cown
}

func NewPointGraph() *PointGraph {
	return &PointGraph{
	}
}

func (self *PointGraph) String() string {
	return "PointGraph()"
}

type _describeBoc struct {
	std.Cown
	g Graph
	node _GraphNodeBound
}

func (self *_describeBoc) String() string {
	return "{ " + "g: " + std.StringifyRepr(self.g) + "; " + "node: " + std.StringifyRepr(self.node) + "; " + "call: {}" + " }"
}

func (self *_describeBoc) Call(g Graph, node _GraphNodeBound) *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		self.g = g
		self.node = node
		return std.Print(self.node.Label().Force())
	})
}

var Describe = &_describeBoc{
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
		var pg *PointGraph
		var p *Point
		std.Schedule(&self.Cown, func() std.Unit {
			pg = &PointGraph{}
			p = NewPoint(std.NewInt(1), std.NewInt(2))
			_bg0.GoWait(Describe.Call(pg, p))
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
