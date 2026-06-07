package main

import std "yz/runtime/rt"

type _ResultVariant int

const (
	_ResultOk _ResultVariant = iota
	_ResultErr
)

type Result[T any, E any] struct {
	_variant _ResultVariant
	value T
	error E
}

func NewResultOk[T any, E any](value T) *Result[T, E] {
	return &Result[T, E]{
		_variant: _ResultOk,
		value: value,
	}
}

func NewResultErr[T any, E any](error E) *Result[T, E] {
	return &Result[T, E]{
		_variant: _ResultErr,
		error: error,
	}
}

func (self *Result[T, E]) String() string {
	switch self._variant {
	case _ResultOk:
		return "Result.Ok(value: " + std.StringifyRepr(self.value) + ")"
	case _ResultErr:
		return "Result.Err(error: " + std.StringifyRepr(self.error) + ")"
	}
	return "Result(?)"
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	var a *Result[std.Int, std.String] = NewResultOk[std.Int, std.String](std.NewInt(42))
	var b *Result[std.Int, std.String] = NewResultErr[std.Int, std.String](std.NewString("division by zero"))
	switch a._variant {
	case _ResultOk:
		std.Print(std.NewString(std.StringifyRepr(a.value)))
	case _ResultErr:
		std.Print(a.error)
	}
	switch b._variant {
	case _ResultOk:
		std.Print(std.NewString(std.StringifyRepr(b.value)))
	case _ResultErr:
		std.Print(b.error)
	}
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
