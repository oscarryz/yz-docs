package main

import std "yz/runtime/rt"

type User struct {
	std.Cown
	name std.String
}

func NewUser(name std.String) *User {
	return &User{
		name: name,
	}
}

func (self *User) String() string {
	return "User(name: " + std.StringifyRepr(self.name) + ")"
}

type Graph struct {
	std.Cown
}

func NewGraph() *Graph {
	return &Graph{
	}
}

func (self *Graph) String() string {
	return "Graph()"
}

type SocialGraph struct {
	std.Cown
}

func NewSocialGraph() *SocialGraph {
	return &SocialGraph{
	}
}

func (self *SocialGraph) String() string {
	return "SocialGraph()"
}

type _acceptBoc struct {
	std.Cown
	g *Graph
	n any
}

func (self *_acceptBoc) String() string {
	return "{ " + "g: " + std.StringifyRepr(self.g) + "; " + "n: " + std.StringifyRepr(self.n) + "; " + "call: {}" + " }"
}

func (self *_acceptBoc) Call(g *Graph, n any) *std.Thunk[std.String] {
	return std.ScheduleMulti([]*std.Cown{&self.Cown, &g.Cown}, func() std.String {
		self.g = g
		self.n = n
		return std.NewString("ok")
	})
}

var Accept = &_acceptBoc{
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
		var s std.String
		var sg *SocialGraph
		var u *User
		std.Schedule(&self.Cown, func() std.Unit {
			sg = &SocialGraph{}
			u = NewUser(std.NewString("Alice"))
			std.GoStore(_bg0, Accept.Call(sg, u), &s)
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		std.Print(s)
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
