package main

import std "yz/runtime/rt"

type Book struct {
	std.Cown
	title std.String
}

func NewBook(title std.String) *Book {
	return &Book{
		title: title,
	}
}

func (self *Book) String() string {
	return "Book(title: " + std.StringifyRepr(self.title) + ")"
}

func (self *Book) Title() std.String {
	return self.title
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	var books std.Dict[std.String, *Book] = std.NewDict[std.String, *Book]().Set(std.NewString("Dune"), NewBook(std.NewString("Dune")))
	if books.Has(std.NewString("Dune")).GoBool() {
		std.Print(std.NewString("has Dune"))
	}
	if books.Has(std.NewString("Foundation")).GoBool() {
	} else {
		std.Print(std.NewString("no Foundation"))
	}
	var b *Book = books.At(std.NewString("Dune"))
	std.Print(b.title)
	std.Print(std.NewString(std.StringifyRepr(books.Length())))
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
