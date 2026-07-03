// Package macrowire implements the wire format between the Yz compiler and
// compiled macro executables (YZC-0028).
//
// The format is Yz source restricted to the non-executable data subset
// (the same subset annotations allow): ShortDecl keys with String / Int /
// Decimal / Bool literals, boc literals, and array literals.
//
// Inbound (compiler → macro stdin):
//
//	subject: {
//	    name: "Person"
//	    fields: [
//	        { name: "name", type: "String" },
//	        { name: "age", type: "Int" }
//	    ]
//	}
//	config: { pretty: true }
//
// Outbound (macro stdout) is a raw Yz boc body of generated slots; it is
// parsed by the compiler directly and needs no support here.
//
// This package lives outside internal/ so the generated macro executable
// (module yzapp) can import it; being part of module yz, it may itself use
// the real Yz parser for decoding.
package macrowire

import (
	"fmt"
	"strconv"
	"strings"

	"yz/internal/ast"
	"yz/internal/parser"
)

// FieldSpec is one field of the subject boc: name and type as written in
// the Yz source (pre-sema, unresolved).
type FieldSpec struct {
	Name string
	Type string
}

// ConfigKind discriminates the scalar kinds allowed in a macro config block.
type ConfigKind int

const (
	ConfigString ConfigKind = iota
	ConfigInt
	ConfigDecimal
	ConfigBool
)

// ConfigValue is one scalar config value.
type ConfigValue struct {
	Kind ConfigKind
	Str  string
	Int  int64
	Dec  float64
	Bool bool
}

// ConfigEntry is one key/value pair of the config block. Entries are ordered
// as written in the annotation.
type ConfigEntry struct {
	Key   string
	Value ConfigValue
}

// Payload is the inbound message: the subject boc's structure plus the
// macro's config block.
type Payload struct {
	SubjectName string
	Fields      []FieldSpec
	Config      []ConfigEntry
}

// ---------------------------------------------------------------------------
// Encoding
// ---------------------------------------------------------------------------

// Encode renders the payload in the inbound wire format.
func (p *Payload) Encode() string {
	var sb strings.Builder
	sb.WriteString("subject: {\n")
	sb.WriteString("    name: " + quoteYz(p.SubjectName) + "\n")
	if len(p.Fields) == 0 {
		sb.WriteString("    fields: []\n")
	} else {
		// One entry per line inside each boc literal: a comma after `name: v`
		// would be consumed as a multi-value ShortDecl continuation, so inline
		// `{ a: 1, b: 2 }` does not parse. Newlines separate entries instead.
		sb.WriteString("    fields: [\n")
		for i, f := range p.Fields {
			sep := ","
			if i == len(p.Fields)-1 {
				sep = ""
			}
			sb.WriteString("        {\n")
			sb.WriteString("            name: " + quoteYz(f.Name) + "\n")
			sb.WriteString("            type: " + quoteYz(f.Type) + "\n")
			sb.WriteString("        }" + sep + "\n")
		}
		sb.WriteString("    ]\n")
	}
	sb.WriteString("}\n")
	sb.WriteString("config: {")
	if len(p.Config) > 0 {
		sb.WriteString("\n")
		for _, e := range p.Config {
			sb.WriteString("    " + e.Key + ": " + e.Value.encode() + "\n")
		}
	}
	sb.WriteString("}\n")
	return sb.String()
}

func (v ConfigValue) encode() string {
	switch v.Kind {
	case ConfigString:
		return quoteYz(v.Str)
	case ConfigInt:
		return strconv.FormatInt(v.Int, 10)
	case ConfigDecimal:
		return strconv.FormatFloat(v.Dec, 'g', -1, 64)
	case ConfigBool:
		return strconv.FormatBool(v.Bool)
	}
	return `""`
}

// quoteYz renders s as a double-quoted Yz string literal. The escapes mirror
// unquoteYz exactly so encode/decode round-trips are lossless. `$` is escaped
// because the parser treats a bare `${` as interpolation, but skips escape
// pairs when splitting.
func quoteYz(s string) string {
	var sb strings.Builder
	sb.WriteByte('"')
	for _, r := range s {
		switch r {
		case '"':
			sb.WriteString(`\"`)
		case '\\':
			sb.WriteString(`\\`)
		case '$':
			sb.WriteString(`\$`)
		case '\n':
			sb.WriteString(`\n`)
		case '\t':
			sb.WriteString(`\t`)
		case '\r':
			sb.WriteString(`\r`)
		default:
			sb.WriteRune(r)
		}
	}
	sb.WriteByte('"')
	return sb.String()
}

// unquoteYz reverses quoteYz: strips the surrounding double quotes and
// unescapes in a single pass (unlike naive ReplaceAll chains, `\\n`
// round-trips as backslash+n, not newline).
func unquoteYz(raw string) (string, error) {
	if len(raw) < 2 || raw[0] != '"' || raw[len(raw)-1] != '"' {
		return "", fmt.Errorf("not a quoted string: %q", raw)
	}
	inner := raw[1 : len(raw)-1]
	var sb strings.Builder
	for i := 0; i < len(inner); i++ {
		c := inner[i]
		if c != '\\' {
			sb.WriteByte(c)
			continue
		}
		i++
		if i >= len(inner) {
			return "", fmt.Errorf("dangling escape in %q", raw)
		}
		switch inner[i] {
		case '"':
			sb.WriteByte('"')
		case '\\':
			sb.WriteByte('\\')
		case '$':
			sb.WriteByte('$')
		case 'n':
			sb.WriteByte('\n')
		case 't':
			sb.WriteByte('\t')
		case 'r':
			sb.WriteByte('\r')
		default:
			// Preserve unknown escapes verbatim.
			sb.WriteByte('\\')
			sb.WriteByte(inner[i])
		}
	}
	return sb.String(), nil
}

// ---------------------------------------------------------------------------
// Decoding
// ---------------------------------------------------------------------------

// DecodePayload parses an inbound wire message using the real Yz parser.
func DecodePayload(src []byte) (*Payload, error) {
	sf, err := parser.New(src).ParseFile()
	if err != nil {
		return nil, fmt.Errorf("parsing payload: %w", err)
	}
	p := &Payload{}
	var haveSubject, haveConfig bool
	for _, stmt := range sf.Stmts {
		sd, ok := stmt.(*ast.ShortDecl)
		if !ok || len(sd.Names) != 1 {
			return nil, fmt.Errorf("payload: unexpected statement %T", stmt)
		}
		switch sd.Names[0].Name {
		case "subject":
			bl, ok := sd.Values[0].(*ast.BocLiteral)
			if !ok {
				return nil, fmt.Errorf("payload: subject is not a boc literal")
			}
			if err := decodeSubject(bl, p); err != nil {
				return nil, err
			}
			haveSubject = true
		case "config":
			bl, ok := sd.Values[0].(*ast.BocLiteral)
			if !ok {
				return nil, fmt.Errorf("payload: config is not a boc literal")
			}
			cfg, err := decodeConfig(bl)
			if err != nil {
				return nil, err
			}
			p.Config = cfg
			haveConfig = true
		default:
			return nil, fmt.Errorf("payload: unknown key %q", sd.Names[0].Name)
		}
	}
	if !haveSubject || !haveConfig {
		return nil, fmt.Errorf("payload: missing subject or config")
	}
	return p, nil
}

func decodeSubject(bl *ast.BocLiteral, p *Payload) error {
	for _, el := range bl.Elements {
		sd, ok := el.(*ast.ShortDecl)
		if !ok || len(sd.Names) != 1 {
			return fmt.Errorf("payload subject: unexpected element %T", el)
		}
		switch sd.Names[0].Name {
		case "name":
			s, err := decodeString(sd.Values[0])
			if err != nil {
				return fmt.Errorf("payload subject name: %w", err)
			}
			p.SubjectName = s
		case "fields":
			arr, ok := sd.Values[0].(*ast.ArrayLiteral)
			if !ok {
				return fmt.Errorf("payload subject fields: not an array literal")
			}
			for _, el := range arr.Elements {
				fbl, ok := el.(*ast.BocLiteral)
				if !ok {
					return fmt.Errorf("payload subject field: not a boc literal")
				}
				fs, err := decodeField(fbl)
				if err != nil {
					return err
				}
				p.Fields = append(p.Fields, fs)
			}
		default:
			return fmt.Errorf("payload subject: unknown key %q", sd.Names[0].Name)
		}
	}
	return nil
}

func decodeField(bl *ast.BocLiteral) (FieldSpec, error) {
	var fs FieldSpec
	for _, el := range bl.Elements {
		sd, ok := el.(*ast.ShortDecl)
		if !ok || len(sd.Names) != 1 {
			return fs, fmt.Errorf("payload field: unexpected element %T", el)
		}
		s, err := decodeString(sd.Values[0])
		if err != nil {
			return fs, fmt.Errorf("payload field %s: %w", sd.Names[0].Name, err)
		}
		switch sd.Names[0].Name {
		case "name":
			fs.Name = s
		case "type":
			fs.Type = s
		default:
			return fs, fmt.Errorf("payload field: unknown key %q", sd.Names[0].Name)
		}
	}
	return fs, nil
}

func decodeConfig(bl *ast.BocLiteral) ([]ConfigEntry, error) {
	var entries []ConfigEntry
	for _, el := range bl.Elements {
		sd, ok := el.(*ast.ShortDecl)
		if !ok || len(sd.Names) != 1 {
			return nil, fmt.Errorf("payload config: unexpected element %T", el)
		}
		v, err := decodeValue(sd.Values[0])
		if err != nil {
			return nil, fmt.Errorf("payload config %s: %w", sd.Names[0].Name, err)
		}
		entries = append(entries, ConfigEntry{Key: sd.Names[0].Name, Value: v})
	}
	return entries, nil
}

func decodeString(e ast.Expr) (string, error) {
	lit, ok := e.(*ast.StringLit)
	if !ok {
		return "", fmt.Errorf("expected string literal, got %T", e)
	}
	return unquoteYz(lit.Value)
}

func decodeValue(e ast.Expr) (ConfigValue, error) {
	switch v := e.(type) {
	case *ast.StringLit:
		s, err := unquoteYz(v.Value)
		if err != nil {
			return ConfigValue{}, err
		}
		return ConfigValue{Kind: ConfigString, Str: s}, nil
	case *ast.IntLit:
		n, err := strconv.ParseInt(v.Value, 10, 64)
		if err != nil {
			return ConfigValue{}, fmt.Errorf("bad int literal %q", v.Value)
		}
		return ConfigValue{Kind: ConfigInt, Int: n}, nil
	case *ast.DecimalLit:
		d, err := strconv.ParseFloat(v.Value, 64)
		if err != nil {
			return ConfigValue{}, fmt.Errorf("bad decimal literal %q", v.Value)
		}
		return ConfigValue{Kind: ConfigDecimal, Dec: d}, nil
	case *ast.Ident:
		switch v.Name {
		case "true":
			return ConfigValue{Kind: ConfigBool, Bool: true}, nil
		case "false":
			return ConfigValue{Kind: ConfigBool, Bool: false}, nil
		}
		return ConfigValue{}, fmt.Errorf("unsupported config value %q", v.Name)
	default:
		return ConfigValue{}, fmt.Errorf("unsupported config value %T", e)
	}
}
