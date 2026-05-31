package main

import std "yz/runtime/rt"

type _GraphNodeBound interface {
	Label() *std.Thunk[std.String]
}


type Graph interface {
}


type City struct {
	std.Cown
	name std.String
}

func NewCity(name std.String) *City {
	return &City{
		name: name,
	}
}

func (self *City) String() string {
	return "City(name: " + std.StringifyRepr(self.name) + ")"
}

func (self *City) label() std.String {
	return self.name
}

func (self *City) Label() *std.Thunk[std.String] {
	return std.Schedule(&self.Cown, func() std.String {
		return self.label()
	})
}

type CityGraph struct {
	std.Cown
}

func NewCityGraph() *CityGraph {
	return &CityGraph{
	}
}

func (self *CityGraph) String() string {
	return "CityGraph()"
}

type _describeBoc struct {
	std.Cown
	g Graph
	n _GraphNodeBound
}

func (self *_describeBoc) String() string {
	return "{ " + "g: " + std.StringifyRepr(self.g) + "; " + "n: " + std.StringifyRepr(self.n) + "; " + "call: {}" + " }"
}

func (self *_describeBoc) Call(g Graph, n _GraphNodeBound) *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		self.g = g
		self.n = n
		return std.Print(self.n.Label().Force())
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
		var cg *CityGraph
		var london *City
		std.Schedule(&self.Cown, func() std.Unit {
			cg = &CityGraph{}
			london = NewCity(std.NewString("London"))
			_bg0.GoWait(Describe.Call(cg, london))
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
