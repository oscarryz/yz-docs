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

type _ColorVariant int

const (
	_ColorCircle _ColorVariant = iota
	_ColorSquare
)

type Color struct {
	_variant _ColorVariant
	hue std.Int
	side std.Int
}

func NewColorCircle(hue std.Int) *Color {
	return &Color{
		_variant: _ColorCircle,
		hue: hue,
	}
}

func NewColorSquare(side std.Int) *Color {
	return &Color{
		_variant: _ColorSquare,
		side: side,
	}
}

func (self *Color) String() string {
	switch self._variant {
	case _ColorCircle:
		return "Color.Circle(hue: " + std.StringifyRepr(self.hue) + ")"
	case _ColorSquare:
		return "Color.Square(side: " + std.StringifyRepr(self.side) + ")"
	}
	return "Color(?)"
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
		var s *Shape
		var c *Color
		std.Schedule(&self.Cown, func() std.Unit {
			std.GoStore(_bg0, NewShapeCircle(std.NewInt(5)), &s)
			std.GoStore(_bg0, NewColorCircle(std.NewInt(180)), &c)
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		std.Print(s.radius)
		std.Print(c.hue)
		var s2 *Shape = NewShapeCircle(std.NewInt(10))
		std.Print(s2.radius)
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
