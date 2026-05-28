package sema

import (
	"strings"
	"yz/internal/ast"
)

// FieldInitState tracks which fields of LOCALLY CONSTRUCTED struct variables
// are definitely assigned on the current control-flow path.
//
// Only variables created via ShortDecl (`b : Bar(...)`) are tracked.
// TypedDecl-no-value parameters (`b Bar`) are NOT tracked; isAssigned returns
// true for them (caller guarantees they're fully initialized).
type FieldInitState struct {
	// locals maps varName → set of definitely-assigned field names.
	// Only locally-constructed struct variables appear here.
	locals map[string]map[string]bool
}

func newFieldInitState() *FieldInitState {
	return &FieldInitState{locals: make(map[string]map[string]bool)}
}

// addLocalVar registers a new locally-constructed struct variable with no
// fields assigned yet.
func (s *FieldInitState) addLocalVar(varName string) {
	s.locals[varName] = make(map[string]bool)
}

// markAssigned marks varName.fieldPath as definitely assigned.
// fieldPath may be dotted ("inner.field"); every prefix is also marked so that
// isAssigned prefix-checks pass ("inner" is implicitly assigned when "inner.field" is).
// No-op if varName is not tracked (it's a parameter — already initialized).
func (s *FieldInitState) markAssigned(varName, fieldPath string) {
	fields, ok := s.locals[varName]
	if !ok {
		return
	}
	fields[fieldPath] = true
	// Mark all prefixes: "inner.field" → also mark "inner".
	parts := strings.Split(fieldPath, ".")
	for i := 1; i < len(parts); i++ {
		fields[strings.Join(parts[:i], ".")] = true
	}
}

// isAssigned reports whether varName.fieldPath is definitely assigned.
// For a dotted path every prefix must also be assigned (e.g. "inner" before "inner.field").
// Returns true if varName is not tracked (parameters are always initialized).
func (s *FieldInitState) isAssigned(varName, fieldPath string) bool {
	fields, ok := s.locals[varName]
	if !ok {
		return true // untracked = parameter or inherited field = always initialized
	}
	parts := strings.Split(fieldPath, ".")
	var prefix string
	for i, p := range parts {
		if i == 0 {
			prefix = p
		} else {
			prefix += "." + p
		}
		if !fields[prefix] {
			return false
		}
	}
	return true
}

// propagateInner copies the field-init state of innerFI[innerVar] into
// varName under fieldPrefix. Used when varName.fieldPrefix is assigned a
// struct value whose own init state is already known.
func (s *FieldInitState) propagateInner(varName, fieldPrefix string, innerFI *FieldInitState, innerVar string) {
	if _, ok := s.locals[varName]; !ok {
		return
	}
	for f, assigned := range innerFI.locals[innerVar] {
		if assigned {
			s.markAssigned(varName, fieldPrefix+"."+f)
		}
	}
}

// clone returns a deep copy of s for branch analysis.
func (s *FieldInitState) clone() *FieldInitState {
	c := &FieldInitState{locals: make(map[string]map[string]bool, len(s.locals))}
	for varName, fields := range s.locals {
		cf := make(map[string]bool, len(fields))
		for f, v := range fields {
			cf[f] = v
		}
		c.locals[varName] = cf
	}
	return c
}

// intersect keeps only the field assignments present in BOTH s and other.
// Variables that appear only in one side (declared inside one branch only)
// are removed — they're going out of scope at the merge point.
func (s *FieldInitState) intersect(other *FieldInitState) {
	for varName, fields := range s.locals {
		otherFields, ok := other.locals[varName]
		if !ok {
			// Only in this branch — going out of scope; remove.
			delete(s.locals, varName)
			continue
		}
		// In both branches — keep only the intersection of assigned fields.
		for f := range fields {
			if !otherFields[f] {
				delete(fields, f)
			}
		}
	}
}

// initLocalVar registers varName as a locally-constructed struct variable and
// marks the fields provided in the constructor call as definitely assigned.
// Fields with HasDefault=true (optional, have default values) are always marked.
func initLocalVar(fi *FieldInitState, varName string, st *StructType, call *ast.CallExpr) {
	fi.addLocalVar(varName)

	// Optional fields (HasDefault=true) always have a value — mark them assigned.
	for _, f := range st.Fields {
		if f.IsTypeField {
			continue // compile-time only; always "assigned"
		}
		if f.HasDefault {
			if _, isMethod := f.Type.(*BocType); !isMethod {
				fi.locals[varName][f.Name] = true
			}
		}
	}

	// Check whether any argument is named.
	hasNamed := false
	for _, arg := range call.Args {
		if arg.Label != "" {
			hasNamed = true
			break
		}
	}

	if hasNamed {
		for _, arg := range call.Args {
			if arg.Label != "" {
				fi.locals[varName][arg.Label] = true
			}
		}
		return
	}

	// Positional arguments: map to required (HasDefault=false, non-method) fields
	// in declaration order.
	var required []string
	for _, f := range st.Fields {
		if f.IsTypeField {
			continue // compile-time only; not a constructor value parameter
		}
		if !f.HasDefault {
			if _, isMethod := f.Type.(*BocType); !isMethod {
				required = append(required, f.Name)
			}
		}
	}
	for i := range call.Args {
		if i < len(required) {
			fi.locals[varName][required[i]] = true
		}
	}
}
