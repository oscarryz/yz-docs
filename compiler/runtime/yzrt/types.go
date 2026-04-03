// Package yzrt is the Yz runtime library imported by generated Go code.
// Generated code refers to these types as std.Int, std.String, etc.
// (the "std" alias is set in the generated import block).
package yzrt

import (
	"fmt"
	"math"
	"strconv"
)

// ---------------------------------------------------------------------------
// Int
// ---------------------------------------------------------------------------

// Int is the boxed integer type.
type Int struct{ val int64 }

// NewInt constructs an Int from a Go int64.
func NewInt(v int64) Int { return Int{val: v} }

// GoInt returns the underlying Go int64 value.
func (i Int) GoInt() int64 { return i.val }

func (i Int) String() string { return strconv.FormatInt(i.val, 10) }

// Arithmetic
func (i Int) Plus(other Int) Int    { return Int{i.val + other.val} }
func (i Int) Minus(other Int) Int   { return Int{i.val - other.val} }
func (i Int) Star(other Int) Int    { return Int{i.val * other.val} }
func (i Int) Slash(other Int) Int   { return Int{i.val / other.val} }
func (i Int) Percent(other Int) Int { return Int{i.val % other.val} }

// Unary negation
func (i Int) Neg() Int { return Int{-i.val} }

// Comparison
func (i Int) Lt(other Int) Bool   { return Bool{i.val < other.val} }
func (i Int) Gt(other Int) Bool   { return Bool{i.val > other.val} }
func (i Int) Lteq(other Int) Bool { return Bool{i.val <= other.val} }
func (i Int) Gteq(other Int) Bool { return Bool{i.val >= other.val} }
func (i Int) Eqeq(other Int) Bool { return Bool{i.val == other.val} }
func (i Int) Neq(other Int) Bool  { return Bool{i.val != other.val} }

// To produces a half-open range [i, end).
func (i Int) To(end Int) Range { return Range{from: i.val, to: end.val} }

// ToStr converts to String.
func (i Int) ToStr() String { return String{val: strconv.FormatInt(i.val, 10)} }

// ---------------------------------------------------------------------------
// Decimal
// ---------------------------------------------------------------------------

// Decimal is the boxed floating-point type.
type Decimal struct{ val float64 }

// NewDecimal constructs a Decimal from a Go float64.
func NewDecimal(v float64) Decimal { return Decimal{val: v} }

// GoFloat64 returns the underlying Go float64.
func (d Decimal) GoFloat64() float64 { return d.val }

func (d Decimal) String() string { return strconv.FormatFloat(d.val, 'g', -1, 64) }

// Arithmetic
func (d Decimal) Plus(other Decimal) Decimal  { return Decimal{d.val + other.val} }
func (d Decimal) Minus(other Decimal) Decimal { return Decimal{d.val - other.val} }
func (d Decimal) Star(other Decimal) Decimal  { return Decimal{d.val * other.val} }
func (d Decimal) Slash(other Decimal) Decimal { return Decimal{d.val / other.val} }

// Unary negation
func (d Decimal) Neg() Decimal { return Decimal{-d.val} }

// Power (extra utility)
func (d Decimal) Pow(exp Decimal) Decimal { return Decimal{math.Pow(d.val, exp.val)} }

// Comparison
func (d Decimal) Lt(other Decimal) Bool   { return Bool{d.val < other.val} }
func (d Decimal) Gt(other Decimal) Bool   { return Bool{d.val > other.val} }
func (d Decimal) Lteq(other Decimal) Bool { return Bool{d.val <= other.val} }
func (d Decimal) Gteq(other Decimal) Bool { return Bool{d.val >= other.val} }
func (d Decimal) Eqeq(other Decimal) Bool { return Bool{d.val == other.val} }
func (d Decimal) Neq(other Decimal) Bool  { return Bool{d.val != other.val} }

// ToStr converts to String.
func (d Decimal) ToStr() String {
	return String{val: strconv.FormatFloat(d.val, 'g', -1, 64)}
}

// ---------------------------------------------------------------------------
// String
// ---------------------------------------------------------------------------

// String is the boxed string type.
type String struct{ val string }

// NewString constructs a String from a Go string.
func NewString(v string) String { return String{val: v} }

// GoString returns the underlying Go string.
func (s String) GoString() string { return s.val }

func (s String) String() string { return s.val }

// Plus concatenates two strings.
func (s String) Plus(other String) String { return String{s.val + other.val} }

// Comparison
func (s String) Eqeq(other String) Bool { return Bool{s.val == other.val} }
func (s String) Neq(other String) Bool  { return Bool{s.val != other.val} }

// Length returns the number of Unicode code points.
func (s String) Length() Int { return Int{int64(len([]rune(s.val)))} }

// ToStr returns the string itself.
func (s String) ToStr() String { return s }

// ---------------------------------------------------------------------------
// Bool
// ---------------------------------------------------------------------------

// Bool is the boxed boolean type.
type Bool struct{ val bool }

// NewBool constructs a Bool from a Go bool.
func NewBool(v bool) Bool { return Bool{val: v} }

// GoBool returns the underlying Go bool.
func (b Bool) GoBool() bool { return b.val }

func (b Bool) String() string {
	if b.val {
		return "true"
	}
	return "false"
}

// Logical
func (b Bool) Ampamp(other Bool) Bool   { return Bool{b.val && other.val} }
func (b Bool) Pipepipe(other Bool) Bool { return Bool{b.val || other.val} }

// Comparison
func (b Bool) Eqeq(other Bool) Bool { return Bool{b.val == other.val} }
func (b Bool) Neq(other Bool) Bool  { return Bool{b.val != other.val} }

// Qm is the conditional operator: flag ? { trueCase } , { falseCase }
// Both branches are passed as zero-argument functions so only the selected
// branch is evaluated.
func (b Bool) Qm(trueCase, falseCase func() any) any {
	if b.val {
		return trueCase()
	}
	return falseCase()
}

// ToStr converts to String.
func (b Bool) ToStr() String { return NewString(b.String()) }

// ---------------------------------------------------------------------------
// Unit
// ---------------------------------------------------------------------------

// Unit is the empty/void type — the result of a boc that returns nothing.
type Unit struct{}

// TheUnit is the singleton Unit value.
var TheUnit = Unit{}

func (u Unit) String() string { return "()" }

// ---------------------------------------------------------------------------
// Stringer helper (for Print / Info)
// ---------------------------------------------------------------------------

// Stringify returns a human-readable string for any yzrt value.
// It handles all boxed types plus raw Go values via fmt.Sprint.
func Stringify(v any) string {
	switch x := v.(type) {
	case fmt.Stringer:
		return x.String()
	default:
		return fmt.Sprint(v)
	}
}
