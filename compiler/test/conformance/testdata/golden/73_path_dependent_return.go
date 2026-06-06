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

type Graph interface {
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

type _makeNodeBoc struct {
	std.Cown
	g Graph
}

func (self *_makeNodeBoc) String() string {
	return "{ " + "g: " + std.StringifyRepr(self.g) + "; " + "call: {}" + " }"
}

func (self *_makeNodeBoc) Call(g Graph) *std.Thunk[any] {
	return std.Schedule(&self.Cown, func() any {
		self.g = g
		return NewUser(std.NewString("test"))
	})
}

var MakeNode = &_makeNodeBoc{
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
		var node *User
		var sg *SocialGraph
		std.Schedule(&self.Cown, func() std.Unit {
			sg = &SocialGraph{}
			_st0 := MakeNode.Call(sg)
			_bg0.Add(func() { node = _st0.Force().(*User) })
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		std.Print(node.name)
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
