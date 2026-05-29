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

type Resolver struct {
	std.Cown
	sg *SocialGraph
}

func NewResolver(sg *SocialGraph) *Resolver {
	return &Resolver{
		sg: sg,
	}
}

func (self *Resolver) String() string {
	return "Resolver(sg: " + std.StringifyRepr(self.sg) + ")"
}

func (self *Resolver) Resolve() *std.Thunk[*User] {
	return std.ScheduleMulti([]*std.Cown{&self.Cown, &self.sg.Cown}, func() *User {
		return NewUser(std.NewString("Alice"))
	})
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
		var u *User
		var sg *SocialGraph
		var r *Resolver
		std.Schedule(&self.Cown, func() std.Unit {
			sg = &SocialGraph{}
			r = NewResolver(sg)
			std.GoStore(_bg0, r.Resolve(), &u)
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		std.Print(u.name)
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
