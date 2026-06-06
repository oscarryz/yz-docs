package rt

// ThunkString, ThunkInt, ThunkDecimal, ThunkBool wrap *Thunk[T] and expose
// forwarding methods that return new cold thunks. This lets boc-call results
// participate in expression chains without forcing the thunk eagerly.
//
// Go generics do not allow type-specific methods on a generic receiver
// (func (th *Thunk[String])), so each scalar gets its own named wrapper.
// Codegen emits these types for boc calls that return scalars; the thunks
// propagate lazily through the expression tree and are forced by BocGroup.Wait.

// ---------------------------------------------------------------------------
// ThunkString
// ---------------------------------------------------------------------------

type ThunkString struct{ t *Thunk[String] }

// GoStringThunk launches fn in a goroutine and returns a hot ThunkString.
func GoStringThunk(fn func() String) ThunkString { return ThunkString{Go(fn)} }

// NewStringThunk returns a cold ThunkString that evaluates fn on first Force.
func NewStringThunk(fn func() String) ThunkString { return ThunkString{NewThunk(fn)} }

func (th ThunkString) Force() String { return th.t.Force() }

func (th ThunkString) Eqeq(other String) ThunkBool {
	return ThunkBool{NewThunk(func() Bool { return th.t.Force().Eqeq(other) })}
}
func (th ThunkString) Neq(other String) ThunkBool {
	return ThunkBool{NewThunk(func() Bool { return th.t.Force().Neq(other) })}
}
func (th ThunkString) Lt(other String) ThunkBool {
	return ThunkBool{NewThunk(func() Bool { return th.t.Force().Lt(other) })}
}
func (th ThunkString) Gt(other String) ThunkBool {
	return ThunkBool{NewThunk(func() Bool { return th.t.Force().Gt(other) })}
}
func (th ThunkString) Lteq(other String) ThunkBool {
	return ThunkBool{NewThunk(func() Bool { return th.t.Force().Lteq(other) })}
}
func (th ThunkString) Gteq(other String) ThunkBool {
	return ThunkBool{NewThunk(func() Bool { return th.t.Force().Gteq(other) })}
}
func (th ThunkString) Plus(other String) ThunkString {
	return ThunkString{NewThunk(func() String { return th.t.Force().Plus(other) })}
}
func (th ThunkString) Length() ThunkInt {
	return ThunkInt{NewThunk(func() Int { return th.t.Force().Length() })}
}
func (th ThunkString) Contains(sub String) ThunkBool {
	return ThunkBool{NewThunk(func() Bool { return th.t.Force().Contains(sub) })}
}
func (th ThunkString) HasPrefix(prefix String) ThunkBool {
	return ThunkBool{NewThunk(func() Bool { return th.t.Force().HasPrefix(prefix) })}
}
func (th ThunkString) HasSuffix(suffix String) ThunkBool {
	return ThunkBool{NewThunk(func() Bool { return th.t.Force().HasSuffix(suffix) })}
}
func (th ThunkString) ToUpper() ThunkString {
	return ThunkString{NewThunk(func() String { return th.t.Force().ToUpper() })}
}
func (th ThunkString) ToLower() ThunkString {
	return ThunkString{NewThunk(func() String { return th.t.Force().ToLower() })}
}
func (th ThunkString) Trim() ThunkString {
	return ThunkString{NewThunk(func() String { return th.t.Force().Trim() })}
}
func (th ThunkString) ToStr() ThunkString {
	return ThunkString{NewThunk(func() String { return th.t.Force().ToStr() })}
}

// ---------------------------------------------------------------------------
// ThunkInt
// ---------------------------------------------------------------------------

type ThunkInt struct{ t *Thunk[Int] }

func GoIntThunk(fn func() Int) ThunkInt    { return ThunkInt{Go(fn)} }
func NewIntThunk(fn func() Int) ThunkInt   { return ThunkInt{NewThunk(fn)} }
func (th ThunkInt) Force() Int             { return th.t.Force() }

func (th ThunkInt) Plus(other Int) ThunkInt    { return ThunkInt{NewThunk(func() Int { return th.t.Force().Plus(other) })} }
func (th ThunkInt) Minus(other Int) ThunkInt   { return ThunkInt{NewThunk(func() Int { return th.t.Force().Minus(other) })} }
func (th ThunkInt) Star(other Int) ThunkInt    { return ThunkInt{NewThunk(func() Int { return th.t.Force().Star(other) })} }
func (th ThunkInt) Slash(other Int) ThunkInt   { return ThunkInt{NewThunk(func() Int { return th.t.Force().Slash(other) })} }
func (th ThunkInt) Percent(other Int) ThunkInt { return ThunkInt{NewThunk(func() Int { return th.t.Force().Percent(other) })} }
func (th ThunkInt) Neg() ThunkInt              { return ThunkInt{NewThunk(func() Int { return th.t.Force().Neg() })} }
func (th ThunkInt) Abs() ThunkInt              { return ThunkInt{NewThunk(func() Int { return th.t.Force().Abs() })} }
func (th ThunkInt) Lt(other Int) ThunkBool     { return ThunkBool{NewThunk(func() Bool { return th.t.Force().Lt(other) })} }
func (th ThunkInt) Gt(other Int) ThunkBool     { return ThunkBool{NewThunk(func() Bool { return th.t.Force().Gt(other) })} }
func (th ThunkInt) Lteq(other Int) ThunkBool   { return ThunkBool{NewThunk(func() Bool { return th.t.Force().Lteq(other) })} }
func (th ThunkInt) Gteq(other Int) ThunkBool   { return ThunkBool{NewThunk(func() Bool { return th.t.Force().Gteq(other) })} }
func (th ThunkInt) Eqeq(other Int) ThunkBool   { return ThunkBool{NewThunk(func() Bool { return th.t.Force().Eqeq(other) })} }
func (th ThunkInt) Neq(other Int) ThunkBool    { return ThunkBool{NewThunk(func() Bool { return th.t.Force().Neq(other) })} }
func (th ThunkInt) To(end Int) ThunkRange      { return ThunkRange{NewThunk(func() Range { return th.t.Force().To(end) })} }
func (th ThunkInt) ToStr() ThunkString         { return ThunkString{NewThunk(func() String { return th.t.Force().ToStr() })} }

// ---------------------------------------------------------------------------
// ThunkDecimal
// ---------------------------------------------------------------------------

type ThunkDecimal struct{ t *Thunk[Decimal] }

func GoDecimalThunk(fn func() Decimal) ThunkDecimal  { return ThunkDecimal{Go(fn)} }
func NewDecimalThunk(fn func() Decimal) ThunkDecimal { return ThunkDecimal{NewThunk(fn)} }
func (th ThunkDecimal) Force() Decimal               { return th.t.Force() }

func (th ThunkDecimal) Plus(other Decimal) ThunkDecimal  { return ThunkDecimal{NewThunk(func() Decimal { return th.t.Force().Plus(other) })} }
func (th ThunkDecimal) Minus(other Decimal) ThunkDecimal { return ThunkDecimal{NewThunk(func() Decimal { return th.t.Force().Minus(other) })} }
func (th ThunkDecimal) Star(other Decimal) ThunkDecimal  { return ThunkDecimal{NewThunk(func() Decimal { return th.t.Force().Star(other) })} }
func (th ThunkDecimal) Slash(other Decimal) ThunkDecimal { return ThunkDecimal{NewThunk(func() Decimal { return th.t.Force().Slash(other) })} }
func (th ThunkDecimal) Neg() ThunkDecimal                { return ThunkDecimal{NewThunk(func() Decimal { return th.t.Force().Neg() })} }
func (th ThunkDecimal) Abs() ThunkDecimal                { return ThunkDecimal{NewThunk(func() Decimal { return th.t.Force().Abs() })} }
func (th ThunkDecimal) Pow(exp Decimal) ThunkDecimal     { return ThunkDecimal{NewThunk(func() Decimal { return th.t.Force().Pow(exp) })} }
func (th ThunkDecimal) Lt(other Decimal) ThunkBool       { return ThunkBool{NewThunk(func() Bool { return th.t.Force().Lt(other) })} }
func (th ThunkDecimal) Gt(other Decimal) ThunkBool       { return ThunkBool{NewThunk(func() Bool { return th.t.Force().Gt(other) })} }
func (th ThunkDecimal) Lteq(other Decimal) ThunkBool     { return ThunkBool{NewThunk(func() Bool { return th.t.Force().Lteq(other) })} }
func (th ThunkDecimal) Gteq(other Decimal) ThunkBool     { return ThunkBool{NewThunk(func() Bool { return th.t.Force().Gteq(other) })} }
func (th ThunkDecimal) Eqeq(other Decimal) ThunkBool     { return ThunkBool{NewThunk(func() Bool { return th.t.Force().Eqeq(other) })} }
func (th ThunkDecimal) Neq(other Decimal) ThunkBool      { return ThunkBool{NewThunk(func() Bool { return th.t.Force().Neq(other) })} }
func (th ThunkDecimal) ToStr() ThunkString               { return ThunkString{NewThunk(func() String { return th.t.Force().ToStr() })} }

// ---------------------------------------------------------------------------
// ThunkBool
// ---------------------------------------------------------------------------

type ThunkBool struct{ t *Thunk[Bool] }

func GoBoolThunk(fn func() Bool) ThunkBool  { return ThunkBool{Go(fn)} }
func NewBoolThunk(fn func() Bool) ThunkBool { return ThunkBool{NewThunk(fn)} }
func (th ThunkBool) Force() Bool            { return th.t.Force() }

func (th ThunkBool) GoBool() bool           { return th.t.Force().GoBool() }
func (th ThunkBool) Ampamp(other Bool) ThunkBool   { return ThunkBool{NewThunk(func() Bool { return th.t.Force().Ampamp(other) })} }
func (th ThunkBool) Pipepipe(other Bool) ThunkBool { return ThunkBool{NewThunk(func() Bool { return th.t.Force().Pipepipe(other) })} }
func (th ThunkBool) Eqeq(other Bool) ThunkBool     { return ThunkBool{NewThunk(func() Bool { return th.t.Force().Eqeq(other) })} }
func (th ThunkBool) Neq(other Bool) ThunkBool      { return ThunkBool{NewThunk(func() Bool { return th.t.Force().Neq(other) })} }
func (th ThunkBool) ToStr() ThunkString            { return ThunkString{NewThunk(func() String { return th.t.Force().ToStr() })} }

// Qm is the conditional: flag ? { trueCase }, { falseCase }
// Both branches are closures so only the selected branch is evaluated.
// Returns a *Thunk[any] so BocGroup.Add can force it.
func (th ThunkBool) Qm(trueCase, falseCase func() any) *Thunk[any] {
	return NewThunk(func() any { return th.t.Force().Qm(trueCase, falseCase) })
}

// ---------------------------------------------------------------------------
// ThunkUnit
// ---------------------------------------------------------------------------

type ThunkUnit struct{ t *Thunk[Unit] }

func GoUnitThunk(fn func() Unit) ThunkUnit  { return ThunkUnit{Go(fn)} }
func NewUnitThunk(fn func() Unit) ThunkUnit { return ThunkUnit{NewThunk(fn)} }
func (th ThunkUnit) Force() Unit            { return th.t.Force() }

// ---------------------------------------------------------------------------
// ThunkRange
// ---------------------------------------------------------------------------

type ThunkRange struct{ t *Thunk[Range] }

func (th ThunkRange) Force() Range { return th.t.Force() }

// ---------------------------------------------------------------------------
// WrapXThunk constructors — wrap an existing *Thunk[T] into the concrete ThunkX
// type. Needed because ThunkX.t is unexported; callers outside this package
// (i.e. generated code) must use these constructors.
// ---------------------------------------------------------------------------

func WrapStringThunk(th *Thunk[String]) ThunkString   { return ThunkString{th} }
func WrapIntThunk(th *Thunk[Int]) ThunkInt             { return ThunkInt{th} }
func WrapBoolThunk(th *Thunk[Bool]) ThunkBool          { return ThunkBool{th} }
func WrapDecimalThunk(th *Thunk[Decimal]) ThunkDecimal { return ThunkDecimal{th} }
func WrapUnitThunk(th *Thunk[Unit]) ThunkUnit          { return ThunkUnit{th} }
func WrapRangeThunk(th *Thunk[Range]) ThunkRange       { return ThunkRange{th} }
