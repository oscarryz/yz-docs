package main

import std "yz/runtime/rt"

type _ShapeVariant int

const (
	_ShapeCircle _ShapeVariant = iota
	_ShapeRectangle
)

type Shape struct {
	_variant _ShapeVariant
	radius std.Int
	width std.Int
	height std.Int
}

func NewShapeCircle(radius std.Int) *Shape {
	return &Shape{
		_variant: _ShapeCircle,
		radius: radius,
	}
}

func NewShapeRectangle(width std.Int, height std.Int) *Shape {
	return &Shape{
		_variant: _ShapeRectangle,
		width: width,
		height: height,
	}
}

func (self *Shape) String() string {
	switch self._variant {
	case _ShapeCircle:
		return "Shape.Circle(radius: " + std.StringifyRepr(self.radius) + ")"
	case _ShapeRectangle:
		return "Shape.Rectangle(width: " + std.StringifyRepr(self.width) + ", height: " + std.StringifyRepr(self.height) + ")"
	}
	return "Shape(?)"
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() any {
	var s *Shape = NewShapeCircle(std.NewInt(5))
	if std.NewBool(s._variant == _ShapeCircle).GoBool() {
		std.Print(std.NewString("is a circle"))
	}
	if std.NewBool(s._variant == _ShapeCircle).GoBool() {
		std.Print(s.radius)
	}
	if std.NewBool(s._variant == _ShapeRectangle).GoBool() {
		std.Print(s.width)
	} else {
		std.Print(std.NewString("not a rectangle"))
	}
	var is_circle std.Bool = std.NewBool(s._variant == _ShapeCircle)
	if is_circle.GoBool() {
		std.Print(std.NewString("still a circle"))
	}
	return std.TheUnit
}

func (self *_mainBoc) Call() *std.Thunk[any] {
	return std.Schedule(&self.Cown, func() any {
		return self.call()
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
