package codegen

// goKeywords is the set of Go reserved words. Yz identifiers are unrestricted
// by this list (e.g. a field named `type`), so any Yz name emitted as a Go
// identifier must pass through goSafeName. Display positions (homoiconic
// String() labels) keep the original Yz name.
var goKeywords = map[string]bool{
	"break": true, "case": true, "chan": true, "const": true, "continue": true,
	"default": true, "defer": true, "else": true, "fallthrough": true, "for": true,
	"func": true, "go": true, "goto": true, "if": true, "import": true,
	"interface": true, "map": true, "package": true, "range": true, "return": true,
	"select": true, "struct": true, "switch": true, "type": true, "var": true,
}

// goSafeName escapes a Yz identifier that collides with a Go keyword by
// appending an underscore. All other names pass through unchanged.
func goSafeName(name string) string {
	if goKeywords[name] {
		return name + "_"
	}
	return name
}
