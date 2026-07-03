// Package macrowire implements the Yz macro wire protocol: the compiler encodes
// a subject boc + config into Yz-source text that is piped to a macro
// executable on stdin, and reads the macro's generated boc body from stdout.
package macrowire

import (
	"fmt"
	"strconv"
	"strings"

	"yz/internal/ast"
	"yz/internal/parser"
)

// FieldSpec describes one data field of the subject boc.
type FieldSpec struct {
	Name string
	Type string // Yz type expression (e.g. "String", "Int", "[String]")
}

// ConfigValue is a scalar value from the macro trigger's config boc.
type ConfigValue interface{ configVal() }

type ConfigString struct{ V string }
type ConfigInt struct{ V int64 }
type ConfigBool struct{ V bool }

func (ConfigString) configVal() {}
func (ConfigInt) configVal()    {}
func (ConfigBool) configVal()   {}

// ConfigEntry is one key-value pair in the config payload (ordered).
type ConfigEntry struct {
	Key   string
	Value ConfigValue
}

// Payload is the full inbound wire message sent to a macro executable on stdin.
type Payload struct {
	SubjectName string
	Fields      []FieldSpec
	Config      []ConfigEntry
}

// Encode serialises p as Yz-source text per the wire format:
//
//	subject: {
//	    name: "Person"
//	    fields: [
//	        { name: "name", type_: "String" },
//	    ]
//	}
//	config: { pretty: true }
func (p *Payload) Encode() string {
	var b strings.Builder
	b.WriteString("subject: {\n")
	b.WriteString("    name: ")
	b.WriteString(strconv.Quote(p.SubjectName))
	b.WriteString("\n    fields: [")
	for i, f := range p.Fields {
		if i > 0 {
			b.WriteString(",\n        { name: ")
		} else {
			b.WriteString("\n        { name: ")
		}
		b.WriteString(strconv.Quote(f.Name))
		b.WriteString("; type_: ")
		b.WriteString(strconv.Quote(f.Type))
		b.WriteString(" }")
	}
	b.WriteString("]\n}\n")
	if len(p.Config) > 0 {
		b.WriteString("config: {\n")
		for _, e := range p.Config {
			b.WriteString("    ")
			b.WriteString(e.Key)
			b.WriteString(": ")
			b.WriteString(encodeConfigValue(e.Value))
			b.WriteString("\n")
		}
		b.WriteString("}\n")
	}
	return b.String()
}

func encodeConfigValue(v ConfigValue) string {
	switch vv := v.(type) {
	case ConfigString:
		return strconv.Quote(vv.V)
	case ConfigInt:
		return strconv.FormatInt(vv.V, 10)
	case ConfigBool:
		if vv.V {
			return "true"
		}
		return "false"
	}
	return `""`
}

// DecodePayload parses a wire-format payload produced by Encode.
// It uses the Yz parser to read the source text and extracts the subject
// name, fields, and config entries.
func DecodePayload(src []byte) (*Payload, error) {
	sf, err := parser.New(src).ParseFile()
	if err != nil {
		return nil, fmt.Errorf("macrowire: parse payload: %w", err)
	}
	p := &Payload{}
	for _, node := range sf.Stmts {
		sd, ok := node.(*ast.ShortDecl)
		if !ok || len(sd.Names) == 0 {
			continue
		}
		switch sd.Names[0].Name {
		case "subject":
			if err := decodeSubject(sd, p); err != nil {
				return nil, err
			}
		case "config":
			if err := decodeConfig(sd, p); err != nil {
				return nil, err
			}
		}
	}
	return p, nil
}

func decodeSubject(sd *ast.ShortDecl, p *Payload) error {
	if len(sd.Values) == 0 {
		return nil
	}
	boc, ok := sd.Values[0].(*ast.BocLiteral)
	if !ok {
		return fmt.Errorf("macrowire: subject value is not a boc literal")
	}
	for _, elem := range boc.Elements {
		kv, ok := elem.(*ast.ShortDecl)
		if !ok || len(kv.Names) == 0 || len(kv.Values) == 0 {
			continue
		}
		key := kv.Names[0].Name
		switch key {
		case "name":
			s, err := decodeString(kv.Values[0])
			if err != nil {
				return fmt.Errorf("macrowire: subject.name: %w", err)
			}
			p.SubjectName = s
		case "fields":
			arr, ok := kv.Values[0].(*ast.ArrayLiteral)
			if !ok {
				return fmt.Errorf("macrowire: subject.fields is not an array")
			}
			for _, elem := range arr.Elements {
				f, err := decodeFieldSpec(elem)
				if err != nil {
					return err
				}
				p.Fields = append(p.Fields, f)
			}
		}
	}
	return nil
}

func decodeFieldSpec(n ast.Expr) (FieldSpec, error) {
	boc, ok := n.(*ast.BocLiteral)
	if !ok {
		return FieldSpec{}, fmt.Errorf("macrowire: field entry is not a boc literal")
	}
	var f FieldSpec
	for _, elem := range boc.Elements {
		kv, ok := elem.(*ast.ShortDecl)
		if !ok || len(kv.Names) == 0 || len(kv.Values) == 0 {
			continue
		}
		switch kv.Names[0].Name {
		case "name":
			s, err := decodeString(kv.Values[0])
			if err != nil {
				return FieldSpec{}, fmt.Errorf("macrowire: field.name: %w", err)
			}
			f.Name = s
		case "type_":
			s, err := decodeString(kv.Values[0])
			if err != nil {
				return FieldSpec{}, fmt.Errorf("macrowire: field.type_: %w", err)
			}
			f.Type = s
		}
	}
	return f, nil
}

func decodeConfig(sd *ast.ShortDecl, p *Payload) error {
	if len(sd.Values) == 0 {
		return nil
	}
	boc, ok := sd.Values[0].(*ast.BocLiteral)
	if !ok {
		return fmt.Errorf("macrowire: config value is not a boc literal")
	}
	for _, elem := range boc.Elements {
		kv, ok := elem.(*ast.ShortDecl)
		if !ok || len(kv.Names) == 0 || len(kv.Values) == 0 {
			continue
		}
		key := kv.Names[0].Name
		val, err := decodeConfigVal(kv.Values[0])
		if err != nil {
			return fmt.Errorf("macrowire: config.%s: %w", key, err)
		}
		p.Config = append(p.Config, ConfigEntry{Key: key, Value: val})
	}
	return nil
}

func decodeConfigVal(n ast.Expr) (ConfigValue, error) {
	switch v := n.(type) {
	case *ast.StringLit:
		s, err := strconv.Unquote(v.Value)
		if err != nil {
			return nil, err
		}
		return ConfigString{V: s}, nil
	case *ast.IntLit:
		i, err := strconv.ParseInt(v.Value, 10, 64)
		if err != nil {
			return nil, err
		}
		return ConfigInt{V: i}, nil
	case *ast.Ident:
		switch v.Name {
		case "true":
			return ConfigBool{V: true}, nil
		case "false":
			return ConfigBool{V: false}, nil
		}
	}
	return nil, fmt.Errorf("unsupported config value type")
}

func decodeString(n ast.Expr) (string, error) {
	lit, ok := n.(*ast.StringLit)
	if !ok {
		return "", fmt.Errorf("expected string literal")
	}
	return strconv.Unquote(lit.Value)
}
