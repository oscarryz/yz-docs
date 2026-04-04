// Package ast defines the Abstract Syntax Tree node types produced by the Yz parser.
//
// Every node embeds a Pos for source location. The tree mirrors the grammar
// defined in spec/02-grammar.ebnf.
package ast

import "yz/internal/token"

// ---------------------------------------------------------------------------
// Position
// ---------------------------------------------------------------------------

// Pos records the source location of a node.
type Pos struct {
	Line int // 1-based
	Col  int // 1-based byte offset from line start
}

func (p Pos) Position() Pos { return p }

// ---------------------------------------------------------------------------
// Node interfaces
// ---------------------------------------------------------------------------

// Node is the base interface for all AST nodes.
type Node interface {
	Position() Pos
}

// Stmt is a node that does not produce a value usable in an expression.
type Stmt interface {
	Node
	stmtNode()
}

// Expr is a node that produces a value.
type Expr interface {
	Node
	exprNode()
}

// TypeExpr represents a type annotation.
type TypeExpr interface {
	Node
	typeNode()
}

// ---------------------------------------------------------------------------
// Source file
// ---------------------------------------------------------------------------

// SourceFile is the root of the AST for one .yz file.
//
// It contains a flat list of top-level statements separated by semicolons
// (after ASI). Each element is either a Stmt or an Expr (expressions can
// appear as statements; their value is the boc's last-expression return
// candidate).
type SourceFile struct {
	Pos
	Stmts []Node // Stmt | Expr
}

// ---------------------------------------------------------------------------
// Statements
// ---------------------------------------------------------------------------

// ShortDecl is the `:` declaration: `name: "Alice"`.
// The LHS may be a list of identifiers for multiple assignment from a
// multi-return call: `a, b: swap(x, y)`.
type ShortDecl struct {
	Pos
	Names  []*Ident // at least one
	Values []Expr   // at least one (len == 1 for single decl, may be multi-return call)
}

func (s *ShortDecl) stmtNode() {}

// TypedDecl is an explicit-type declaration: `age Int` or `age Int = 30`.
// When Value is nil the declaration is uninitialized (acts as a required parameter).
type TypedDecl struct {
	Pos
	Name  *Ident
	Type  TypeExpr
	Value Expr // nil if no initializer
}

func (s *TypedDecl) stmtNode() {}

// Assignment assigns to an existing variable or member: `name = "Bob"`.
// For multiple assignment, Names has more than one element and Values
// has the matching expressions: `a, b = swap("x", "y")`.
type Assignment struct {
	Pos
	// Single-target assignment: Target is set, Names is nil.
	// Multi-target: Names is set, Target is nil.
	Target Expr    // used for single assignment (including member/index access)
	Names  []*Ident // used for multi-assignment LHS
	Values []Expr   // RHS; usually one expression (possibly multi-return call)
}

func (s *Assignment) stmtNode() {}

// ReturnStmt is `return [expr]`.
type ReturnStmt struct {
	Pos
	Value Expr // nil for bare `return`
}

func (s *ReturnStmt) stmtNode() {}

// BreakStmt is `break`.
type BreakStmt struct{ Pos }

func (s *BreakStmt) stmtNode() {}

// ContinueStmt is `continue`.
type ContinueStmt struct{ Pos }

func (s *ContinueStmt) stmtNode() {}

// MixStmt is `mix Identifier` — flattens another boc's fields into the current boc.
type MixStmt struct {
	Pos
	Name *Ident
}

func (s *MixStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// Expressions
// ---------------------------------------------------------------------------

// Ident is an identifier: word (lowercase), type (uppercase multi-char), or
// generic (single uppercase letter).
type Ident struct {
	Pos
	Name    string
	TokType token.Type // IDENT | TYPE_IDENT | GENERIC_IDENT
}

func (e *Ident) exprNode() {}

// IntLit is an integer literal: `42`.
type IntLit struct {
	Pos
	Value string // raw text
}

func (e *IntLit) exprNode() {}

// DecimalLit is a decimal literal: `3.14`.
type DecimalLit struct {
	Pos
	Value string // raw text
}

func (e *DecimalLit) exprNode() {}

// StringLit is a string literal (either quote style, with optional interpolation).
// Value includes the surrounding quotes as they appear in source.
type StringLit struct {
	Pos
	Value string // raw text including delimiters and escape sequences
}

func (e *StringLit) exprNode() {}

// UnaryExpr is a unary negation: `-expr`.
// In Yz, unary `-` is the only unary operator (desugars to `expr.-()` in sema).
type UnaryExpr struct {
	Pos
	Op      string // always "-"
	Operand Expr
}

func (e *UnaryExpr) exprNode() {}

// BinaryExpr is a non-word method invocation: `left op right`.
// All non-word methods have equal precedence, evaluated left-to-right.
// Desugars to `left.op(right)` in sema/codegen.
type BinaryExpr struct {
	Pos
	Left  Expr
	Op    string // the non-word identifier literal, e.g. "+", "==", "?"
	Right Expr
}

func (e *BinaryExpr) exprNode() {}

// CallExpr is a method or boc invocation: `greet("Alice")` or `person.greet()`.
// The callee is an Expr (may be Ident, MemberExpr, etc.).
type CallExpr struct {
	Pos
	Callee Expr
	Args   []*Argument
}

func (e *CallExpr) exprNode() {}

// Argument is a single call argument, optionally named.
type Argument struct {
	Pos
	Label string // "" for positional
	Value Expr
}

// MemberExpr is member access: `person.name`.
type MemberExpr struct {
	Pos
	Object Expr
	Member *Ident
}

func (e *MemberExpr) exprNode() {}

// IndexExpr is index access: `array[0]` or `dict["key"]`.
type IndexExpr struct {
	Pos
	Object Expr
	Index  Expr
}

func (e *IndexExpr) exprNode() {}

// GroupExpr is a parenthesised expression: `(expr)`.
// A single-element group is just for precedence; it does not create a tuple.
type GroupExpr struct {
	Pos
	Expr Expr
}

func (e *GroupExpr) exprNode() {}

// BocLiteral is a `{ BocBody }` literal.
// Elements are the statements/expressions inside the body.
// InfoString (if present) immediately precedes the enclosing declaration in
// the source; it is attached here for AST consumers that need it.
type BocLiteral struct {
	Pos
	Elements   []Node // Stmt | Expr | *VariantDef
	InfoString *StringLit // nil if no info string precedes this boc
}

func (e *BocLiteral) exprNode() {}

// BocWithSig is `name #(params) [body]` in two forms:
//
//   - Shorthand: `name #(params) { body }` — params auto-scoped into body;
//     unlabeled types at the end of the param list are return-type annotations.
//
//   - Body-only: `name #(params) = { body }` — body redeclares its own params
//     as the first N TypedDecl statements; ALL sig params are inputs (none are
//     return-type annotations); return type is inferred from body's last expr.
//
// The sema phase enforces the matching rules; the parser just captures both.
type BocWithSig struct {
	Pos
	Name     *Ident
	Sig      *BocTypeExpr
	Body     *BocLiteral // nil if signature-only declaration (no body)
	BodyOnly bool        // true when `= { body }` form (body redeclares params)
}

func (e *BocWithSig) stmtNode() {}

// InterpPart is one segment of an InterpolatedStringExpr.
// Either a literal text fragment or an embedded expression.
type InterpPart struct {
	IsExpr bool
	Text   string // raw text content (no outer quotes) for text parts
	Expr   Expr   // non-nil for expression parts
}

// InterpolatedStringExpr is a string with backtick-embedded expressions:
// `"Hello, `name`!"` desugars to a Plus chain at the IR level.
type InterpolatedStringExpr struct {
	Pos
	Parts []InterpPart
}

func (e *InterpolatedStringExpr) exprNode() {}

// ConditionalExpr is `cond ? {trueCase}, {falseCase}` — a Bool conditional.
// TrueCase and FalseCase are typically BocLiterals.
type ConditionalExpr struct {
	Pos
	Cond      Expr
	TrueCase  Expr
	FalseCase Expr
}

func (e *ConditionalExpr) exprNode() {}

// ArrayLiteral is `[expr, expr, ...]` or the empty form `[Type]()`.
type ArrayLiteral struct {
	Pos
	Elements []Expr   // nil for empty-type form
	ElemType TypeExpr // non-nil for empty-type form `[Type]()`
}

func (e *ArrayLiteral) exprNode() {}

// DictLiteral is `[key: value, ...]` or the empty form `[K:V]()`.
type DictLiteral struct {
	Pos
	Entries  []*DictEntry // nil for empty-type form
	KeyType  TypeExpr     // non-nil for empty-type form
	ValType  TypeExpr     // non-nil for empty-type form
}

func (e *DictLiteral) exprNode() {}

// DictEntry is one `key: value` pair inside a DictLiteral.
type DictEntry struct {
	Pos
	Key   Expr
	Value Expr
}

// MatchExpr is either a variant match or a condition match.
//
//   Variant match:   `match expr { Variant.Case() => ... }, ...`
//   Condition match: `match { cond => expr }, ...`
//
// Subject is nil for condition match.
type MatchExpr struct {
	Pos
	Subject Expr             // nil for condition match
	Arms    []*ConditionalBoc
}

func (e *MatchExpr) exprNode() {}

// ConditionalBoc is one `{ [condition =>] body }` arm of a match expression.
type ConditionalBoc struct {
	Pos
	Condition Expr   // nil for the default (no-condition) arm
	IsVariant bool   // true if condition is a variant constructor pattern
	Body      []Node // Stmt | Expr
}

// InfoString is a string literal that immediately precedes a declaration and
// attaches metadata. The compiler stores it in the AST but generates no runtime
// code for it (tooling support deferred).
type InfoString struct {
	Pos
	Value string // raw text including delimiters
}

func (e *InfoString) exprNode() {}

// VariantDef is a variant constructor declaration inside a type boc body:
// `Ok(value T)` or `Err(error E)`.
type VariantDef struct {
	Pos
	Name   string // TYPE_IDENT
	Params []*BocParam
}

func (v *VariantDef) stmtNode() {}

// ---------------------------------------------------------------------------
// Type expressions
// ---------------------------------------------------------------------------

// SimpleTypeExpr is a bare type name: `Int`, `String`, `T`, or a generic
// application: `Option(T)`.
type SimpleTypeExpr struct {
	Pos
	Name     string
	TokType  token.Type // TYPE_IDENT | GENERIC_IDENT
	TypeArgs []TypeExpr // non-nil for generic application
}

func (t *SimpleTypeExpr) typeNode() {}

// ArrayTypeExpr is `[T]`.
type ArrayTypeExpr struct {
	Pos
	ElemType TypeExpr
}

func (t *ArrayTypeExpr) typeNode() {}

// DictTypeExpr is `[K:V]`.
type DictTypeExpr struct {
	Pos
	KeyType TypeExpr
	ValType TypeExpr
}

func (t *DictTypeExpr) typeNode() {}

// BocTypeExpr is `#(params)` — a boc type signature.
type BocTypeExpr struct {
	Pos
	Params []*BocParam
}

func (t *BocTypeExpr) typeNode() {}

// BocParam is one entry in a boc signature parameter list.
// It covers three forms:
//
//   1. Named parameter:          `name Type`            (Name set, Type set)
//   2. Anonymous parameter:      `Type`                 (Name empty, Type set)
//   3. Named with default:       `name Type = default`  (all fields set)
//   4. Variant constructor:      `Ok(value T)`          (Variant set)
type BocParam struct {
	Pos
	Label   string   // parameter label; empty for anonymous or variant
	Type    TypeExpr // nil if Variant is set
	Default Expr     // nil if no default value
	Variant *VariantDef // non-nil for a variant constructor param
}
