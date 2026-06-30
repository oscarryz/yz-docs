package main

import std "yz/runtime/rt"

type Node struct {
	std.Cown
	value std.Int
	next *Node
}

func NewNode(value std.Int, next *Node) *Node {
	return &Node{
		value: value,
		next: next,
	}
}

func (self *Node) String() string {
	return "Node(value: " + std.StringifyRepr(self.value) + ", next: " + std.StringifyRepr(self.next) + ")"
}

func (self *Node) Value() std.Int {
	return self.value
}

func (self *Node) Next() *Node {
	return self.next
}

type _first_valueBoc struct {
	std.Cown
	n *Node
}

func (self *_first_valueBoc) String() string {
	return "{ " + "n: " + std.StringifyRepr(self.n) + "; " + "call: {}" + " }"
}

func (self *_first_valueBoc) Call(n *Node) std.Int {
	return std.LazyInt(std.ScheduleMulti([]*std.Cown{&self.Cown, &n.Cown}, func() std.Int {
		self.n = n
		return self.n.value
	}))
}

var First_value = &_first_valueBoc{
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	std.Print(std.NewString("ok"))
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
