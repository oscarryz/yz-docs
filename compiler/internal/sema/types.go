// Package sema implements semantic analysis for Yz: scope resolution,
// type inference, type checking, FQN assignment, and structural compatibility.
package sema

import (
	"fmt"
	"strings"
)

// ---------------------------------------------------------------------------
// Type interface
// ---------------------------------------------------------------------------

// Type represents any Yz type.
type Type interface {
	// typeName returns a human-readable representation for error messages.
	typeName() string
	// IsCompatibleWith reports whether this type is structurally compatible
	// with target (i.e., this type can be used where target is expected).
	// Width subtyping: a wider type (more fields) is compatible with a narrower target.
	IsCompatibleWith(target Type) bool
}

// ---------------------------------------------------------------------------
// Built-in scalar types
// ---------------------------------------------------------------------------

// BuiltinType is one of the five scalar built-in types.
type BuiltinType struct {
	name string
}

func (t *BuiltinType) typeName() string { return t.name }

func (t *BuiltinType) IsCompatibleWith(target Type) bool {
	// Built-ins are compatible only with the same built-in or Unknown (error recovery).
	switch u := target.(type) {
	case *BuiltinType:
		return t.name == u.name
	case *UnknownType:
		return true
	case *GenericType:
		return true // generics are resolved at use-site; always compatible for now
	}
	return false
}

func (t *BuiltinType) String() string { return t.name }

// Singleton built-in types — use these throughout the compiler.
var (
	TypInt     = &BuiltinType{name: "Int"}
	TypDecimal = &BuiltinType{name: "Decimal"}
	TypString  = &BuiltinType{name: "String"}
	TypBool    = &BuiltinType{name: "Bool"}
	TypUnit    = &BuiltinType{name: "Unit"}
)

// ---------------------------------------------------------------------------
// Boc type
// ---------------------------------------------------------------------------

// BocParam describes one parameter in a boc type signature.
type BocParam struct {
	Label      string // empty for anonymous / return-type-only entries
	Type       Type
	HasDefault bool
	IsReturn   bool // true for return-type-position entries
}

// BocType is the structural type of a boc (function/block).
// Parameters represent inputs; Returns are the types of the last expression(s).
type BocType struct {
	Params  []BocParam
	Returns []Type
}

func (t *BocType) typeName() string {
	var b strings.Builder
	b.WriteString("#(")
	for i, p := range t.Params {
		if i > 0 {
			b.WriteString(", ")
		}
		if p.Label != "" {
			b.WriteString(p.Label)
			b.WriteString(" ")
		}
		b.WriteString(p.Type.typeName())
		if p.HasDefault {
			b.WriteString(" = ...")
		}
	}
	for _, r := range t.Returns {
		b.WriteString(", ")
		b.WriteString(r.typeName())
	}
	b.WriteString(")")
	return b.String()
}

func (t *BocType) IsCompatibleWith(target Type) bool {
	switch u := target.(type) {
	case *BocType:
		// Structural boc compatibility: same number of required params and
		// return types, with compatible individual types.
		if len(t.Params) != len(u.Params) {
			return false
		}
		for i := range t.Params {
			if !t.Params[i].Type.IsCompatibleWith(u.Params[i].Type) {
				return false
			}
		}
		if len(t.Returns) != len(u.Returns) {
			return false
		}
		for i := range t.Returns {
			if !t.Returns[i].IsCompatibleWith(u.Returns[i]) {
				return false
			}
		}
		return true
	case *UnknownType:
		return true
	case *GenericType:
		return true
	default:
		_ = u
		return false
	}
}

func (t *BocType) String() string { return t.typeName() }

// ---------------------------------------------------------------------------
// Struct type (user-defined uppercase boc or anonymous structural type)
// ---------------------------------------------------------------------------

// StructField is one field in a struct type.
type StructField struct {
	Name string
	Type Type
}

// VariantCase is one constructor in a sum type (variant type).
// For `Pet: { Cat(name String, lives Int), Dog(name String, years Int) }`,
// there are two VariantCases: Cat and Dog.
type VariantCase struct {
	Name   string
	Fields []StructField
}

// StructType is the structural type of a user-defined boc or any named type.
// Two StructTypes are compatible if one has at least all the fields of the other.
// IsInterface is set for type-only declarations (`Name #(methods...)`) where all
// params are BocTypes — these generate Go interfaces, not structs.
// IsVariant is set for sum types: `Pet: { Cat(...), Dog(...) }`.
// IsSingleton is set for lowercase body-form bocs that have inner structure
// (inner bocs or BocWithSig methods). These are singletons, not constructor types.
// Returns holds the body's last-expression types for singleton bocs (the call return type).
// TypeParams holds formal type parameter names for generic types (e.g., ["V"] for Option[V]).
// TypeConstraints maps each type param to the methods inferred as required on it.
type StructType struct {
	Name            string                          // may be empty for anonymous structural types
	Fields          []StructField                   // in declaration order (merged across all variants)
	IsInterface     bool                            // true when declared as Name #(boc-params...) with no body
	IsVariant       bool                            // true for sum/variant types
	IsSingleton     bool                            // true for lowercase bocs with inner structure
	Returns         []Type                          // body return types (only when IsSingleton=true)
	Variants        []VariantCase                   // variant constructors (only when IsVariant=true)
	TypeParams      []string                        // formal type parameter names (non-nil for generic types)
	TypeConstraints map[string][]*GenericConstraint // typeParam → inferred method requirements
}

func (t *StructType) typeName() string {
	if t.Name != "" {
		return t.Name
	}
	var parts []string
	for _, f := range t.Fields {
		parts = append(parts, fmt.Sprintf("%s %s", f.Name, f.Type.typeName()))
	}
	return "{ " + strings.Join(parts, "; ") + " }"
}

// fieldMap returns a name→type map for fast lookup.
func (t *StructType) fieldMap() map[string]Type {
	m := make(map[string]Type, len(t.Fields))
	for _, f := range t.Fields {
		m[f.Name] = f.Type
	}
	return m
}

// IsCompatibleWith implements width subtyping:
// t is compatible with target if t has at least all of target's fields
// with compatible types.
func (t *StructType) IsCompatibleWith(target Type) bool {
	switch u := target.(type) {
	case *StructType:
		srcFields := t.fieldMap()
		for _, tf := range u.Fields {
			srcType, ok := srcFields[tf.Name]
			if !ok {
				return false // target has a field this type doesn't
			}
			if !srcType.IsCompatibleWith(tf.Type) {
				return false
			}
		}
		return true
	case *UnknownType:
		return true
	case *GenericType:
		return true
	}
	return false
}

func (t *StructType) String() string { return t.typeName() }

// ---------------------------------------------------------------------------
// Array type
// ---------------------------------------------------------------------------

// ArrayType is the type `[Elem]`.
type ArrayType struct {
	Elem Type
}

func (t *ArrayType) typeName() string { return "[" + t.Elem.typeName() + "]" }

func (t *ArrayType) IsCompatibleWith(target Type) bool {
	switch u := target.(type) {
	case *ArrayType:
		return t.Elem.IsCompatibleWith(u.Elem)
	case *UnknownType:
		return true
	case *GenericType:
		return true
	}
	return false
}

func (t *ArrayType) String() string { return t.typeName() }

// ---------------------------------------------------------------------------
// Dict type
// ---------------------------------------------------------------------------

// DictType is the type `[Key:Val]`.
type DictType struct {
	Key Type
	Val Type
}

func (t *DictType) typeName() string {
	return "[" + t.Key.typeName() + ":" + t.Val.typeName() + "]"
}

func (t *DictType) IsCompatibleWith(target Type) bool {
	switch u := target.(type) {
	case *DictType:
		return t.Key.IsCompatibleWith(u.Key) && t.Val.IsCompatibleWith(u.Val)
	case *UnknownType:
		return true
	case *GenericType:
		return true
	}
	return false
}

func (t *DictType) String() string { return t.typeName() }

// ---------------------------------------------------------------------------
// Generic constraints
// ---------------------------------------------------------------------------

// GenericConstraint records one method required on a generic type parameter.
// These are inferred by analyzing how T-typed values are used inside the bodies
// of methods defined on a generic type.
type GenericConstraint struct {
	TypeParam  string // which type param needs this (e.g. "T")
	MethodName string // required method name (e.g. "to_string", "eqeq")
	Line       int    // source line where this call appears in the generic body
	Col        int    // source column
	Context    string // "StructName.methodName" — which method requires this
}

// ---------------------------------------------------------------------------
// Generic type parameter
// ---------------------------------------------------------------------------

// GenericType represents an unresolved single-letter type parameter (T, E, K, V…).
// It is compatible with any type during the early analysis passes; generics are
// fully resolved in a later substitution pass.
type GenericType struct {
	Name string // single uppercase letter
}

func (t *GenericType) typeName() string { return t.Name }

func (t *GenericType) IsCompatibleWith(_ Type) bool { return true }

func (t *GenericType) String() string { return t.Name }

// ---------------------------------------------------------------------------
// Thunk type (lazy boc invocation result)
// ---------------------------------------------------------------------------

// ThunkType wraps the result type of a boc invocation. Every boc call is
// non-blocking; the result is a lazy thunk that materializes on first use.
type ThunkType struct {
	Inner Type
}

func (t *ThunkType) typeName() string { return "Thunk[" + t.Inner.typeName() + "]" }

func (t *ThunkType) IsCompatibleWith(target Type) bool {
	// A Thunk[T] is compatible with T (it materializes on use) and with
	// Thunk[U] when T is compatible with U.
	switch u := target.(type) {
	case *ThunkType:
		return t.Inner.IsCompatibleWith(u.Inner)
	case *UnknownType:
		return true
	case *GenericType:
		return true
	default:
		// A thunk is implicitly compatible with its inner type because it
		// materializes before use.
		return t.Inner.IsCompatibleWith(target)
	}
}

func (t *ThunkType) String() string { return t.typeName() }

// ---------------------------------------------------------------------------
// Unknown / error type
// ---------------------------------------------------------------------------

// UnknownType is used for error recovery — when a type cannot be determined,
// further type errors are suppressed to avoid cascading failures.
type UnknownType struct{}

func (t *UnknownType) typeName() string { return "<unknown>" }

func (t *UnknownType) IsCompatibleWith(_ Type) bool { return true }

func (t *UnknownType) String() string { return "<unknown>" }

// Unknown is the singleton error-recovery type.
var Unknown = &UnknownType{}

// ---------------------------------------------------------------------------
// Namespace type (directory hierarchy node)
// ---------------------------------------------------------------------------

// NamespaceType represents a directory-level namespace node.
// Children maps the next path segment name to a symbol whose type is either
// another NamespaceType (subdirectory) or a PackageType (leaf package).
type NamespaceType struct {
	Children map[string]*Symbol
}

func (t *NamespaceType) typeName() string             { return "<namespace>" }
func (t *NamespaceType) IsCompatibleWith(_ Type) bool { return false }

// ---------------------------------------------------------------------------
// Package type (leaf of namespace tree — imported Go package)
// ---------------------------------------------------------------------------

// PackageType represents a compiled Yz sub-package that has been imported.
// Exports maps each exported Yz name to its sema Symbol.
type PackageType struct {
	PkgAlias   string            // Go import alias, e.g. "front"
	ImportPath string            // full Go import path, e.g. "yzapp/house/front"
	Exports    map[string]*Symbol
}

func (t *PackageType) typeName() string             { return "<package:" + t.ImportPath + ">" }
func (t *PackageType) IsCompatibleWith(_ Type) bool { return false }
