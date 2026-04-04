// Package ir defines the Intermediate Representation used between the
// semantic analysis phase and Go code generation.
//
// The IR models Go concepts directly (structs, methods, goroutines, thunks)
// rather than Yz concepts. The lowerer translates the Yz AST+sema into IR;
// the codegen emits Go source from IR.
package ir

// ---------------------------------------------------------------------------
// Top-level file
// ---------------------------------------------------------------------------

// File is the IR for one .yz source file → one .go output file.
type File struct {
	PkgName string   // Go package name ("main" for the entry file)
	Imports []string // import paths, in addition to the standard yzrt import
	Decls   []Decl
}

// ---------------------------------------------------------------------------
// Declarations
// ---------------------------------------------------------------------------

// Decl is a top-level declaration in an IR file.
type Decl interface{ irDecl() }

func (*StructDecl) irDecl()    {}
func (*SingletonDecl) irDecl() {}
func (*FuncDecl) irDecl()      {}
func (*InterfaceDecl) irDecl() {}

// StructDecl represents a user-defined Go struct type (from an uppercase Yz boc).
// Each call to the type creates a new instance via the generated constructor.
// When NoConstructor is true (type-only declaration: `Name #(params)`), the
// struct type is emitted without a constructor — instances cannot be created
// via `Name(args)` until a body is attached.
// When IsVariant is true, the struct is a sum/discriminant type: codegen emits
// a discriminant enum, flat merged struct, and per-variant constructors.
type StructDecl struct {
	Name          string
	Fields        []*FieldSpec
	Methods       []*MethodDecl
	NoConstructor bool
	IsVariant     bool
	Variants      []*IRVariantCase
}

// IRVariantCase is one constructor arm of a variant struct.
type IRVariantCase struct {
	Name   string       // e.g. "Cat"
	Fields []*FieldSpec // fields belonging to this variant (with types)
}

// SingletonDecl represents a lowercase Yz boc — a single, persistent instance.
// It generates:
//   - a private struct type  (_<name>Boc)
//   - a package-level var    (var <name> = &_<name>Boc{...})
//   - methods on that struct
type SingletonDecl struct {
	TypeName string // e.g. "_counterBoc"
	VarName  string // e.g. "counter"
	Fields   []*FieldSpec
	Methods  []*MethodDecl
}

// FuncDecl is a standalone Go function (used for the main entry point and
// top-level helpers).
type FuncDecl struct {
	Name    string
	Params  []*ParamSpec
	Results []string // Go type strings
	Body    []Stmt
}

// MethodDecl is a Go method attached to a struct.
type MethodDecl struct {
	RecvType string // e.g. "*_counterBoc"
	RecvName string // e.g. "self"
	Name     string
	Params   []*ParamSpec
	Results  []string // Go type strings (usually one *std.Thunk[T])
	Body     []Stmt
}

// InterfaceDecl represents a structural interface type (from `Name #(methods...)`
// with no body, all params being BocTypes). Generates a Go interface.
type InterfaceDecl struct {
	Name    string
	Methods []*InterfaceMethod
}

// InterfaceMethod is one method entry in a Go interface declaration.
type InterfaceMethod struct {
	Name       string
	Params     []*ParamSpec
	ResultType string // the T in *std.Thunk[T]
}

// FieldSpec is one field in a struct with an optional initializer.
// When Embedded is true, this field was created from a `mix` statement:
//   - the Go struct body emits just the type name (Go embedding syntax)
//   - EmbeddedFields lists the mixed-in type's own fields, used to generate
//     the constructor's parameter list and the NewMixedType(...) call.
type FieldSpec struct {
	Name           string
	Type           string // Go type string, e.g. "std.Int"
	Init           Expr   // may be nil (zero value)
	Embedded       bool
	EmbeddedFields []*FieldSpec // only set when Embedded=true
}

// ParamSpec is a function/method parameter.
type ParamSpec struct {
	Name string
	Type string // Go type string
}

// ---------------------------------------------------------------------------
// Statements
// ---------------------------------------------------------------------------

// Stmt is any IR statement.
type Stmt interface{ irStmt() }

func (*DeclStmt) irStmt()    {}
func (*AssignStmt) irStmt()  {}
func (*ReturnStmt) irStmt()  {}
func (*ExprStmt) irStmt()    {}
func (*ForStmt) irStmt()     {}
func (*IfStmt) irStmt()      {}
func (*WaitStmt) irStmt()    {}
func (*SwitchStmt) irStmt()  {}

// DeclStmt declares a local variable.
// If Type is empty, codegen uses `:=` (Go type inference).
type DeclStmt struct {
	Name string
	Type string // may be empty for := inference
	Init Expr
}

// AssignStmt mutates an existing variable or field: Target = Value.
type AssignStmt struct {
	Target Expr
	Value  Expr
}

// ReturnStmt returns a value from a function or closure.
// Value nil → return std.TheUnit.
type ReturnStmt struct {
	Value Expr
}

// ExprStmt is a statement-level expression (side-effect call, spawn, etc.).
type ExprStmt struct {
	Expr Expr
}

// ForStmt is a condition-controlled loop (compiles Yz `while`).
type ForStmt struct {
	Cond Expr   // must produce std.Bool
	Body []Stmt
}

// IfStmt is a conditional branch.
type IfStmt struct {
	Cond Expr
	Then []Stmt
	Else []Stmt
}

// WaitStmt emits `<groupVar>.Wait()` — used at end of scopes with spawned goroutines.
type WaitStmt struct {
	GroupVar string // local *std.BocGroup var name
}

// SwitchStmt is a discriminant match lowered to a Go switch statement.
// Subject is the variable being matched; TypeName is the Go discriminant enum type.
type SwitchStmt struct {
	Subject  Expr
	TypeName string // e.g. "_PetVariant" (used in switch subject._variant)
	Cases    []*SwitchCase
}

// SwitchCase is one case arm of a SwitchStmt.
type SwitchCase struct {
	ConstName string // e.g. "_PetCat"
	Body      []Stmt
}

// ---------------------------------------------------------------------------
// Expressions
// ---------------------------------------------------------------------------

// Expr is any IR expression.
type Expr interface{ irExpr() }

func (*IntLit) irExpr()      {}
func (*DecimalLit) irExpr()  {}
func (*StringLit) irExpr()   {}
func (*BoolLit) irExpr()     {}
func (*UnitLit) irExpr()     {}
func (*Ident) irExpr()       {}
func (*MethodCall) irExpr()  {}
func (*FuncCall) irExpr()    {}
func (*FieldAccess) irExpr() {}
func (*IndexExpr) irExpr()   {}
func (*ThunkExpr) irExpr()   {}
func (*ForceExpr) irExpr()   {}
func (*ClosureExpr) irExpr() {}
func (*SpawnExpr) irExpr()   {}
func (*NewGroupExpr) irExpr(){}
func (*MatchExpr) irExpr()    {}
func (*SwitchExpr) irExpr()   {}

// Literal nodes — codegen boxes these into std.NewXxx(...) calls.
type IntLit struct{ Val int64 }
type DecimalLit struct{ Val float64 }
type StringLit struct{ Val string }
type BoolLit struct{ Val bool }
type UnitLit struct{} // → std.TheUnit

// Ident is an identifier reference in the current scope.
type Ident struct{ Name string }

// MethodCall is recv.Method(args...).
// Method is already the Go name (e.g. "Plus" for the Yz operator "+").
type MethodCall struct {
	Recv   Expr
	Method string
	Args   []Expr
}

// FuncCall is func(args...) — for free functions and package-level calls.
type FuncCall struct {
	Func Expr
	Args []Expr
}

// FieldAccess is object.Field — used for receiver field access in methods.
type FieldAccess struct {
	Object Expr
	Field  string
}

// IndexExpr is object[index].
type IndexExpr struct {
	Object Expr
	Index  Expr
}

// ThunkExpr wraps a body in std.Go(...) (spawned goroutine) or
// std.NewThunk(...) (synchronous lazy), returning *std.Thunk[ResultType].
type ThunkExpr struct {
	ResultType string // e.g. "std.Int" or "std.Unit"
	Body       []Stmt // the closure body
	Spawn      bool   // true → std.Go, false → std.NewThunk
}

// ForceExpr materializes a thunk: Thunk.Force().
type ForceExpr struct {
	Thunk Expr
}

// ClosureExpr is an anonymous func literal (for passing boc args to builtins
// like While, or for the Qm conditional).
type ClosureExpr struct {
	Params     []*ParamSpec
	ResultType string // return type of the closure; "any" for Qm branches
	Body       []Stmt
}

// SpawnExpr is g.Go(func() any { body }) — launch a goroutine in a BocGroup.
// The result is typed as *std.Thunk[any]; the caller force-casts if needed.
type SpawnExpr struct {
	GroupVar string // the *std.BocGroup local var
	Body     []Stmt
}

// NewGroupExpr creates a new BocGroup: &std.BocGroup{}.
type NewGroupExpr struct{}

// MatchExpr is a condition match lowered to an immediately-invoked closure
// (IIFE) for expression position. Use tryLowerMatch to emit as an IfStmt
// chain in statement position instead.
type MatchExpr struct {
	ResultType string
	Arms       []*MatchArm
}

// MatchArm is one arm of a MatchExpr.
// Cond nil means this is the default (else) branch.
type MatchArm struct {
	Cond Expr
	Body []Stmt
}

// SwitchExpr is a discriminant match in expression position (emitted as an IIFE).
type SwitchExpr struct {
	Subject    Expr
	ResultType string
	Cases      []*SwitchCase
}
