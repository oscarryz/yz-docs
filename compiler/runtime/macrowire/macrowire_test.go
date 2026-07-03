package macrowire_test

import (
	"testing"

	"yz/runtime/macrowire"
)

func TestRoundTrip_empty(t *testing.T) {
	p := &macrowire.Payload{SubjectName: "Empty"}
	enc := p.Encode()
	got, err := macrowire.DecodePayload([]byte(enc))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.SubjectName != "Empty" {
		t.Errorf("name: got %q, want %q", got.SubjectName, "Empty")
	}
	if len(got.Fields) != 0 {
		t.Errorf("fields: got %d, want 0", len(got.Fields))
	}
	if len(got.Config) != 0 {
		t.Errorf("config: got %d, want 0", len(got.Config))
	}
}

func TestRoundTrip_scalars(t *testing.T) {
	p := &macrowire.Payload{
		SubjectName: "Person",
		Fields: []macrowire.FieldSpec{
			{Name: "name", Type: "String"},
			{Name: "age", Type: "Int"},
		},
		Config: []macrowire.ConfigEntry{
			{Key: "pretty", Value: macrowire.ConfigBool{V: true}},
			{Key: "prefix", Value: macrowire.ConfigString{V: "debug"}},
			{Key: "indent", Value: macrowire.ConfigInt{V: 4}},
		},
	}
	enc := p.Encode()
	got, err := macrowire.DecodePayload([]byte(enc))
	if err != nil {
		t.Fatalf("decode: %v\nencoded:\n%s", err, enc)
	}
	if got.SubjectName != "Person" {
		t.Errorf("name: got %q, want %q", got.SubjectName, "Person")
	}
	if len(got.Fields) != 2 {
		t.Fatalf("fields len: got %d, want 2", len(got.Fields))
	}
	if got.Fields[0].Name != "name" || got.Fields[0].Type != "String" {
		t.Errorf("field[0]: got {%q,%q}", got.Fields[0].Name, got.Fields[0].Type)
	}
	if got.Fields[1].Name != "age" || got.Fields[1].Type != "Int" {
		t.Errorf("field[1]: got {%q,%q}", got.Fields[1].Name, got.Fields[1].Type)
	}
	if len(got.Config) != 3 {
		t.Fatalf("config len: got %d, want 3", len(got.Config))
	}
	if b, ok := got.Config[0].Value.(macrowire.ConfigBool); !ok || !b.V {
		t.Errorf("config[0]: want true bool")
	}
	if s, ok := got.Config[1].Value.(macrowire.ConfigString); !ok || s.V != "debug" {
		t.Errorf("config[1]: want string 'debug'")
	}
	if i, ok := got.Config[2].Value.(macrowire.ConfigInt); !ok || i.V != 4 {
		t.Errorf("config[2]: want int 4")
	}
}

func TestRoundTrip_quoteEscaping(t *testing.T) {
	p := &macrowire.Payload{
		SubjectName: `Say "hi"`,
		Fields:      []macrowire.FieldSpec{{Name: "msg", Type: "String"}},
		Config:      []macrowire.ConfigEntry{{Key: "sep", Value: macrowire.ConfigString{V: `a"b`}}},
	}
	enc := p.Encode()
	got, err := macrowire.DecodePayload([]byte(enc))
	if err != nil {
		t.Fatalf("decode: %v\nencoded:\n%s", err, enc)
	}
	if got.SubjectName != `Say "hi"` {
		t.Errorf("name: got %q", got.SubjectName)
	}
	if s, ok := got.Config[0].Value.(macrowire.ConfigString); !ok || s.V != `a"b` {
		t.Errorf("config sep: got %v", got.Config[0].Value)
	}
}
