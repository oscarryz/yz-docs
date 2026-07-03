package macrowire

import (
	"reflect"
	"testing"
)

func roundTrip(t *testing.T, p *Payload) *Payload {
	t.Helper()
	encoded := p.Encode()
	got, err := DecodePayload([]byte(encoded))
	if err != nil {
		t.Fatalf("DecodePayload failed: %v\npayload:\n%s", err, encoded)
	}
	return got
}

func TestRoundTripBasic(t *testing.T) {
	p := &Payload{
		SubjectName: "Person",
		Fields: []FieldSpec{
			{Name: "name", Type: "String"},
			{Name: "age", Type: "Int"},
		},
		Config: []ConfigEntry{
			{Key: "pretty", Value: ConfigValue{Kind: ConfigBool, Bool: true}},
		},
	}
	got := roundTrip(t, p)
	if !reflect.DeepEqual(p, got) {
		t.Errorf("round trip mismatch:\nwant %+v\ngot  %+v", p, got)
	}
}

func TestRoundTripEmpty(t *testing.T) {
	p := &Payload{SubjectName: "Empty"}
	got := roundTrip(t, p)
	if got.SubjectName != "Empty" || len(got.Fields) != 0 || len(got.Config) != 0 {
		t.Errorf("round trip mismatch: %+v", got)
	}
}

func TestRoundTripAllScalarKinds(t *testing.T) {
	p := &Payload{
		SubjectName: "Server",
		Config: []ConfigEntry{
			{Key: "host", Value: ConfigValue{Kind: ConfigString, Str: "localhost"}},
			{Key: "port", Value: ConfigValue{Kind: ConfigInt, Int: 8080}},
			{Key: "timeout", Value: ConfigValue{Kind: ConfigDecimal, Dec: 1.5}},
			{Key: "secure", Value: ConfigValue{Kind: ConfigBool, Bool: false}},
		},
	}
	got := roundTrip(t, p)
	if !reflect.DeepEqual(p, got) {
		t.Errorf("round trip mismatch:\nwant %+v\ngot  %+v", p, got)
	}
}

func TestRoundTripEscaping(t *testing.T) {
	tricky := []string{
		`plain`,
		`has "quotes"`,
		`back\slash`,
		`double\\backslash`,
		"newline\nand\ttab",
		`interp ${x} marker`,
		`literal \n not newline`,
	}
	for _, s := range tricky {
		p := &Payload{
			SubjectName: s,
			Config: []ConfigEntry{
				{Key: "v", Value: ConfigValue{Kind: ConfigString, Str: s}},
			},
		}
		got := roundTrip(t, p)
		if got.SubjectName != s {
			t.Errorf("subject name %q round-tripped as %q", s, got.SubjectName)
		}
		if got.Config[0].Value.Str != s {
			t.Errorf("config value %q round-tripped as %q", s, got.Config[0].Value.Str)
		}
	}
}

func TestRoundTripComplexTypes(t *testing.T) {
	p := &Payload{
		SubjectName: "Container",
		Fields: []FieldSpec{
			{Name: "items", Type: "[Int]"},
			{Name: "lookup", Type: "[String:Int]"},
			{Name: "cb", Type: "#(x Int, String)"},
			{Name: "opt", Type: "Option(Int)"},
		},
	}
	got := roundTrip(t, p)
	if !reflect.DeepEqual(p, got) {
		t.Errorf("round trip mismatch:\nwant %+v\ngot  %+v", p, got)
	}
}

func TestDecodeHandWritten(t *testing.T) {
	src := `
subject: {
    name: "Person"
    fields: [
        {
            name: "name"
            type: "String"
        },
        {
            name: "age"
            type: "Int"
        }
    ]
}
config: { pretty: true }
`
	p, err := DecodePayload([]byte(src))
	if err != nil {
		t.Fatalf("DecodePayload: %v", err)
	}
	if p.SubjectName != "Person" {
		t.Errorf("subject name = %q", p.SubjectName)
	}
	want := []FieldSpec{{Name: "name", Type: "String"}, {Name: "age", Type: "Int"}}
	if !reflect.DeepEqual(p.Fields, want) {
		t.Errorf("fields = %+v", p.Fields)
	}
	if len(p.Config) != 1 || p.Config[0].Key != "pretty" || !p.Config[0].Value.Bool {
		t.Errorf("config = %+v", p.Config)
	}
}

func TestDecodeErrors(t *testing.T) {
	subject := "subject: {\n    name: \"A\"\n    fields: []\n}"
	cases := map[string]string{
		"missing config":  subject,
		"missing subject": `config: {}`,
		"unknown key":     subject + "\nconfig: {}\nextra: {}",
		"non-string name": "subject: {\n    name: 42\n    fields: []\n}\nconfig: {}",
	}
	for label, src := range cases {
		if _, err := DecodePayload([]byte(src)); err == nil {
			t.Errorf("%s: expected error, got none", label)
		}
	}
}
