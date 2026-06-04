package rt

// OptionVariant is the discriminant for the built-in Option type.
type OptionVariant int

const (
	OptionSome OptionVariant = iota
	OptionNone
)

// Option represents a value that may be present (Some) or absent (None).
// It is the built-in optional type returned by dict access d[k].
type Option[V any] struct {
	Variant OptionVariant
	Value   V
}

func NewOptionSome[V any](v V) *Option[V] {
	return &Option[V]{Variant: OptionSome, Value: v}
}

func NewOptionNone[V any]() *Option[V] {
	return &Option[V]{Variant: OptionNone}
}

func (o *Option[V]) String() string {
	if o.Variant == OptionNone {
		return "None()"
	}
	return "Some(value: " + StringifyRepr(o.Value) + ")"
}

func (o *Option[V]) ToStr() String {
	return NewString(o.String())
}
