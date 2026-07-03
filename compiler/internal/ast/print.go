package ast

import "strings"

// TypeExprString renders a TypeExpr back to its Yz source form.
// Used by the macro wire format (YZC-0028) to serialize subject field types
// as written in source (pre-sema, no resolution).
func TypeExprString(te TypeExpr) string {
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
		params := make([]string, len(t.Params))
		for i, p := range t.Params {
			params[i] = bocParamString(p)
		}
		return "#(" + strings.Join(params, ", ") + ")"
	case *MemberTypeExpr:
		return t.Object + "." + t.Member
	default:
		return ""
	}
}

// bocParamString renders one boc signature parameter. Default values are
// omitted — only the label/type/constraint shape matters for serialization.
func bocParamString(p *BocParam) string {
	var sb strings.Builder
	if p.Label != "" {
		sb.WriteString(p.Label)
		if p.Type != nil {
			sb.WriteString(" ")
		}
	}
	if p.Type != nil {
		sb.WriteString(TypeExprString(p.Type))
	}
	for _, c := range p.Constraints {
		sb.WriteString(" ")
		sb.WriteString(TypeExprString(c))
	}
	return sb.String()
}
