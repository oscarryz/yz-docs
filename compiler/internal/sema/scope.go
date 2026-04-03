package sema

import (
	"fmt"
	"strings"

	"yz/internal/ast"
)

// ---------------------------------------------------------------------------
// Symbol
// ---------------------------------------------------------------------------

// Symbol is a named entity in a scope.
type Symbol struct {
	Name string
	Type Type
	FQN  string   // fully-qualified name (empty for locals without a global FQN)
	Node ast.Node // declaration site; nil for built-ins
}

// ---------------------------------------------------------------------------
// Scope
// ---------------------------------------------------------------------------

// Scope is one frame in the lexical scope chain.
// Inner scopes point to their enclosing parent.
type Scope struct {
	parent *Scope
	syms   map[string]*Symbol
}

// newScope creates a child scope of parent (nil for the root).
func newScope(parent *Scope) *Scope {
	return &Scope{parent: parent, syms: make(map[string]*Symbol)}
}

// Define adds a new symbol to this scope.
// It does NOT check for redeclaration — the analyzer does that.
func (s *Scope) Define(sym *Symbol) {
	s.syms[sym.Name] = sym
}

// Lookup searches for name starting in this scope and walking up the parent chain.
// Returns nil if not found.
func (s *Scope) Lookup(name string) *Symbol {
	for frame := s; frame != nil; frame = frame.parent {
		if sym, ok := frame.syms[name]; ok {
			return sym
		}
	}
	return nil
}

// LookupLocal searches only in this scope frame (not parents).
func (s *Scope) LookupLocal(name string) *Symbol {
	return s.syms[name]
}

// ---------------------------------------------------------------------------
// Built-in scope
// ---------------------------------------------------------------------------

// builtinMethods maps type name → map of method name → return type.
// This is the minimum needed for type-checking the first milestone.
// Non-word method names use the symbol-name convention (plus, qm, eqeq, etc.).
var builtinMethods = map[string]map[string]Type{
	"Int": {
		"plus":    TypInt,     // +
		"minus":   TypInt,     // -
		"star":    TypInt,     // *
		"slash":   TypInt,     // /
		"percent": TypInt,     // %
		"lt":      TypBool,    // <
		"gt":      TypBool,    // >
		"lteq":    TypBool,    // <=
		"gteq":    TypBool,    // >=
		"eqeq":    TypBool,    // ==
		"neq":     TypBool,    // !=
		"to":      &StructType{Name: "Range"}, // 1.to(10) → Range
		"to_string": TypString,
	},
	"Decimal": {
		"plus":    TypDecimal,
		"minus":   TypDecimal,
		"star":    TypDecimal,
		"slash":   TypDecimal,
		"lt":      TypBool,
		"gt":      TypBool,
		"lteq":    TypBool,
		"gteq":    TypBool,
		"eqeq":    TypBool,
		"neq":     TypBool,
		"to_string": TypString,
	},
	"String": {
		"plus":    TypString,   // + (concatenation)
		"eqeq":    TypBool,
		"neq":     TypBool,
		"length":  TypInt,
		"to_string": TypString,
	},
	"Bool": {
		"ampamp":    TypBool,  // &&
		"pipepipe":  TypBool,  // ||
		"qm":        Unknown,  // ? — return type depends on boc args; resolved later
		"eqeq":      TypBool,
		"neq":       TypBool,
	},
}

// nonWordToMethodName maps a non-word operator literal to its Go symbol name.
var nonWordToMethodName = map[string]string{
	"+":   "plus",
	"-":   "minus",
	"*":   "star",
	"/":   "slash",
	"%":   "percent",
	"<":   "lt",
	">":   "gt",
	"<=":  "lteq",
	">=":  "gteq",
	"==":  "eqeq",
	"!=":  "neq",
	"&&":  "ampamp",
	"||":  "pipepipe",
	"?":   "qm",
	"<<":  "ltlt",
	">>":  "gtgt",
	"!":   "bang",
	"++":  "plusplus",
	"--":  "minusminus",
}

// NonWordMethodName returns the Go-safe method name for a non-word operator.
// If the operator is not in the built-in table, it falls back to a mechanical
// character-by-character name (e.g. "!=>` → "bangEqGt").
func NonWordMethodName(op string) string {
	if name, ok := nonWordToMethodName[op]; ok {
		return name
	}
	// Mechanical fallback: name each character.
	var b strings.Builder
	charNames := map[rune]string{
		'+': "plus", '-': "minus", '*': "star", '/': "slash",
		'%': "percent", '<': "lt", '>': "gt", '=': "eq",
		'!': "bang", '&': "amp", '|': "pipe", '?': "qm",
		'~': "tilde", '^': "caret", '@': "at", '#': "hash",
		'$': "dollar", '.': "dot", ':': "colon", ';': "semi",
	}
	first := true
	for _, ch := range op {
		name, ok := charNames[ch]
		if !ok {
			name = fmt.Sprintf("u%04x", ch)
		}
		if first {
			b.WriteString(name)
			first = false
		} else {
			// CamelCase the subsequent parts
			if len(name) > 0 {
				b.WriteString(strings.ToUpper(name[:1]) + name[1:])
			}
		}
	}
	return b.String()
}

// newBuiltinScope creates the root scope pre-seeded with built-in identifiers.
func newBuiltinScope() *Scope {
	s := newScope(nil)

	// Primitive type symbols (used as type-name identifiers in declarations).
	for _, bt := range []*BuiltinType{TypInt, TypDecimal, TypString, TypBool, TypUnit} {
		s.Define(&Symbol{Name: bt.name, Type: bt})
	}

	// true / false — constants of type Bool.
	s.Define(&Symbol{Name: "true", Type: TypBool})
	s.Define(&Symbol{Name: "false", Type: TypBool})

	// print — #(value String) — accepts any type via structural widening in codegen;
	// for type-checking purposes we model it as accepting Unknown (anything).
	s.Define(&Symbol{Name: "print", Type: &BocType{
		Params:  []BocParam{{Label: "value", Type: Unknown}},
		Returns: []Type{TypUnit},
	}})

	// while — #(cond #(Bool), body #()) — loop primitive
	s.Define(&Symbol{Name: "while", Type: &BocType{
		Params: []BocParam{
			{Label: "cond", Type: &BocType{Returns: []Type{TypBool}}},
			{Label: "body", Type: &BocType{}},
		},
		Returns: []Type{TypUnit},
	}})

	// info — #(value T, InfoResult) — deferred; modeled as unknown
	s.Define(&Symbol{Name: "info", Type: &BocType{
		Params:  []BocParam{{Label: "value", Type: Unknown}},
		Returns: []Type{Unknown},
	}})

	// http — built-in singleton backed by runtime/yzrt.Http.
	// Methods: get #(uri String) String, post #(uri String, body String) String.
	s.Define(&Symbol{Name: "http", Type: &StructType{
		Name: "_httpBoc",
		Fields: []StructField{
			{Name: "get", Type: &BocType{
				Params:  []BocParam{{Label: "uri", Type: TypString}},
				Returns: []Type{TypString},
			}},
			{Name: "post", Type: &BocType{
				Params:  []BocParam{{Label: "uri", Type: TypString}, {Label: "body", Type: TypString}},
				Returns: []Type{TypString},
			}},
		},
	}})

	return s
}
