package ast

import "strings"

// TypeExprString returns the Yz source representation of a TypeExpr.
// Used by macrowire to serialize field types in the wire payload.
func TypeExprString(te TypeExpr) string {
	if te == nil {
		return ""
	}
	switch t := te.(type) {
	case *SimpleTypeExpr:
		if len(t.TypeArgs) == 0 {
			return t.Name
		}
		args := make([]string, len(t.TypeArgs))
		for i, a := range t.TypeArgs {
			args[i] = TypeExprString(a)
		}
		return t.Name + "(" + strings.Join(args, ", ") + ")"
	case *ArrayTypeExpr:
		return "[" + TypeExprString(t.ElemType) + "]"
	case *DictTypeExpr:
		return "[" + TypeExprString(t.KeyType) + ":" + TypeExprString(t.ValType) + "]"
	case *BocTypeExpr:
		if len(t.Params) == 0 {
			return "#()"
		}
		parts := make([]string, len(t.Params))
		for i, p := range t.Params {
			if p.Label != "" {
				parts[i] = p.Label + " " + TypeExprString(p.Type)
			} else {
				parts[i] = TypeExprString(p.Type)
			}
		}
		return "#(" + strings.Join(parts, ", ") + ")"
	case *MemberTypeExpr:
		return t.Object + "." + t.Member
	}
	return ""
}
