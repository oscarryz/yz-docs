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

func (self *User) Name() std.String {
	return self.name
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

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	var u *User = NewUser(std.NewString("Alice"))
	std.Print(u.name)
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
